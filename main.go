package main

import (
	"flag"
	"log"
	"path"
	"strings"

	"github.com/cheggaaa/pb/v3"
)

func main() {
	var url string
	flag.StringVar(&url, "u", "", "manga page to download from")
	var outputDir string
	flag.StringVar(&outputDir, "o", "./output", "output directory")
	var language string
	flag.StringVar(&language, "l", "English", "language to download")
	var delay int64
	flag.Int64Var(&delay, "d", 1000, "delay in milliseconds between requests")
	var chapterRange string
	flag.StringVar(&chapterRange, "r", "", "chapter range string, specify list to list chapters available or nothing to download everything. Refer to the README for more info.")
	flag.Parse()

	if url == "" {
		log.Fatal("please specify a url!")
	}

	log.Printf("downloading from %s to %s in %s with a delay of %dms...", url, outputDir, language, delay)

	d, err := NewDownloader(url, outputDir, delay)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("getting chapter metadata...")
	chaptersData, err := d.GetChapterData()
	if err != nil {
		log.Fatal(err)
	}

	languageCode, err := chaptersData.GetLanguageCode(language)
	if err != nil {
		log.Fatal(err)
	}

	if chapterRange == "list" {
		log.Println("chapter numbers:")
		log.Println(strings.Join(chaptersData.GetChapterNumbers(), ", "))
		return
	}

	var ranges []*Range
	masterRange, err := chaptersData.GetChapterRange()
	if err != nil {
		log.Fatal(err)
	}

	if chapterRange == "" {
		ranges = append(ranges, masterRange)
	} else {
		r, err := ParseRange(chapterRange, masterRange)
		if err != nil {
			log.Fatal(err)
		}

		ranges = r
	}

	downloadChapters := d.SectionChaptersDownloadData(chaptersData.Metadata, languageCode, ranges)
	if len(downloadChapters) == 0 {
		log.Fatal("no chapters in the specified language found!")
	}

	mangaDir := GetDirName(chaptersData.Title)
	err = DirExists(path.Join(d.OutputDir, mangaDir))
	if err != nil {
		log.Fatal(err)
	}

	for _, chapter := range downloadChapters {
		log.Printf("Retrieving %s metadata...", chapter.FormattedTitle())

		downloadData, err := d.GetChapterDownloadData(chapter)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Retrieving %s...\n", downloadData.Name)

		bar := pb.StartNew(len(downloadData.Urls))
		err = d.GetChapterImages(downloadData, mangaDir, func() {
			bar.Increment()
		})
		bar.Finish()

		if err != nil {
			log.Fatal(err)
		}
	}
	log.Println("done!")
}
