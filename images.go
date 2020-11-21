package main

import (
	"fmt"
	"io"
	"log"
	"mime"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/cheggaaa/pb/v3"
)

func (d *Downloader) GetChaptersImages(downloadData []*ChapterDownloadData, mangaDir string) error {
	// Make sure master subfolder exists
	err := DirExists(path.Join(d.OutputDir, mangaDir))
	if err != nil {
		return err
	}

	for _, dd := range downloadData {
		log.Printf("retrieving %s...\n", dd.Name)

		bar := pb.StartNew(len(dd.Urls))

		err = d.GetChapterImages(dd, mangaDir, func() {
			bar.Increment()
		})

		bar.Finish()

		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Downloader) GetChapterImages(downloadData *ChapterDownloadData, mangaDir string, update func()) error {
	// Create padded format string for image names
	nameFmtStr := fmt.Sprintf("%%0%dd%%s", len(strconv.Itoa(len(downloadData.Urls))))

	// Get images
	for i, url := range downloadData.Urls {
		res, err := d.GetPageRaw(url, nil)
		if err != nil {
			return err
		}

		// Get extension (if extension corresponding to content-type exists
		extension := ".png" // Default to .png
		if content_type, exists := res.Header["Content-Type"]; exists {
			if len(content_type) > 0 {
				extensions, err := mime.ExtensionsByType(content_type[0])
				if err != nil {
					return err
				}

				if extensions != nil && len(extensions) > 0 {
					extension = extensions[0]
				}
			}
		}

		// Get os friendly folder name
		chapDir := GetDirName(downloadData.Name)

		// Make sure subfolder exists
		err = DirExists(path.Join(d.OutputDir, mangaDir, chapDir))
		if err != nil {
			return err
		}

		// Create and write file
		file, err := os.Create(path.Join(d.OutputDir, mangaDir, chapDir, fmt.Sprintf(nameFmtStr, i+1, extension)))
		if err != nil {
			return err
		}

		_, err = io.Copy(file, res.Body)
		if err != nil {
			return err
		}

		// Update progress bar
		update()

		// Don't trigger DDOS protection or something
		time.Sleep(d.Delay)
	}

	return nil
}
