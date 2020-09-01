package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cheggaaa/pb/v3"
)

func (c *ChapterPageData) GetLanguageCode(name string) (string, error) {
	for n, code := range c.Languages {
		if strings.Contains(n, name) {
			return code, nil
		}
	}

	var errorMessage strings.Builder
	errorMessage.WriteString(fmt.Sprintf("language %s is not available! Printing list of available languages:\n", name))

	for k := range c.Languages {
		errorMessage.WriteString(k + "\n")
	}

	return "", errors.New(errorMessage.String())
}

type ChapterDownloadData struct {
	Name	string
	Urls	[]string
}

type RawDownloadData struct {
	Hash			string		`json:"hash"`
	Server		string		`json:"server"`
	PageArray	[]string	`json:"page_array"`
}

func (d *Downloader) GetChaptersDownloadData(metadata []*ChapterMetadata, languageCode string, ranges []*Range) ([]*ChapterDownloadData, error) {
	downloads := d.SectionChaptersDownloadData(metadata, languageCode, ranges)
	if len(downloads) == 0 {
		return nil, errors.New("no chapters in the specified language found!")
	}

	bar := pb.StartNew(len(downloads))
	downloadData := make([]*ChapterDownloadData, 0)
	for _, metadata := range downloads {
		dd, err := d.GetChapterDownloadData(metadata);
		if err != nil {
			return nil, err;
		}

		downloadData = append(downloadData, dd)

		bar.Increment()
		time.Sleep(d.Delay)
	}

	bar.Finish()

	return downloadData, nil
}
func (d *Downloader) SectionChaptersDownloadData(metadata []*ChapterMetadata, languageCode string, ranges []*Range) []*ChapterMetadata {
	downloads := make([]*ChapterMetadata, 0)
	for _, m := range metadata {
		if m.Lang == languageCode && InRange(ranges, m.Chapter) {
			downloads = append(downloads, m)
		}
	}

	return downloads
}

func (d *Downloader) GetChapterDownloadData(metadata *ChapterMetadata) (*ChapterDownloadData, error) {
		rawData, err := d.GetPage(fmt.Sprintf("https://mangadex.org/api/chapter/%s", metadata.Id), nil)
		if err != nil {
			return nil, err
		}

		var rawDownloadData RawDownloadData
		err = json.Unmarshal([]byte(rawData), &rawDownloadData)
		if err != nil {
			return nil, err
		}

		if rawDownloadData.Server[len(rawDownloadData.Server) - 1] != '/' {
			rawDownloadData.Server += "/"
		}

		chapterUrls := make([]string, 0)
		for _, url := range rawDownloadData.PageArray {
			chapterUrl := fmt.Sprintf("%s%s/%s", rawDownloadData.Server, rawDownloadData.Hash, url)
			chapterUrls = append(chapterUrls, chapterUrl)
		}

		dd := &ChapterDownloadData {
			Name: metadata.FormattedTitle(),
			Urls: chapterUrls,
		}

		return dd, nil
}
