package main

import (
	"flag"
	"log"
	"strings"
)

func main() {
	var url string
	flag.StringVar(&url, "u", "", "manga page to download from")
	var outputDir string
	flag.StringVar(&outputDir, "o", "./output", "output directory")
	var language string
	flag.StringVar(&language, "l", "English", "language to download")
	var delay int64
	flag.Int64Var(&delay, "d", 5000, "delay in milliseconds between requests")
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

	log.Println("getting download data...")
	downloadData, err := d.GetChapterDownloadData(chaptersData.Metadata, languageCode, ranges)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("getting images...")
	err = d.GetChaptersImages(downloadData, chaptersData.Title)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("done!")
}
