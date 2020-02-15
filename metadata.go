package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
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

	for k, _ := range c.Languages {
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

func (d *Downloader) GetChapterDownloadData(metadata []*ChapterMetadata, languageCode string, ranges []*Range) ([]*ChapterDownloadData, error) {
	downloads := make([]*ChapterMetadata, 0)
	for _, m := range metadata {
		if m.Lang == languageCode && InRange(ranges, m.Chapter) {
			downloads = append(downloads, m)
		}
	}

	if len(downloads) == 0 {
		return nil, errors.New("no chapters in the specified language found!")
	}

	bar := pb.StartNew(len(downloads))
	downloadData := make([]*ChapterDownloadData, 0)
	for _, download := range downloads {
		rawData, err := d.GetPage(fmt.Sprintf("https://mangadex.org/api/chapter/%s", download.Id), nil)
		if err != nil {
			return nil, err
		}

		var rawDownloadData RawDownloadData
		err = json.Unmarshal([]byte(rawData), &rawDownloadData)
		if err != nil {
			return nil, err
		}

		chapterUrls := make([]string, 0)
		for _, url := range rawDownloadData.PageArray {
			chapterUrl := fmt.Sprintf("%s/%s/%s", rawDownloadData.Server, rawDownloadData.Hash, url)
			chapterUrls = append(chapterUrls, chapterUrl)
		}

		titleBits := make([]string, 0)
		if download.Volume != "" {
			titleBits = append(titleBits, fmt.Sprintf("Vol. %s", download.Volume))
		}
		if download.Chapter != 0.0 {
			titleBits = append(titleBits, fmt.Sprintf("Ch. %s", strconv.FormatFloat(download.Chapter, 'f', -1, 64)))
		}
		if download.Title != "" {
			titleBits = append(titleBits, download.Title)
		}

		dd := &ChapterDownloadData {
			Name: strings.Join(titleBits, " "),
			Urls: chapterUrls,
		}

		downloadData = append(downloadData, dd)

		bar.Increment()
		time.Sleep(d.Delay)
	}

	bar.Finish()

	return downloadData, nil
}
