package main

import (
	"io/ioutil"
	"net/http"
	"time"
)

type Downloader struct {
	Client *http.Client
	OutputDir string
	Url string
	Delay time.Duration
}

func NewDownloader(url string, outputDir string, delay int64) (*Downloader, error) {
	client := &http.Client{}

	err := DirExists(outputDir)
	if err != nil {
		return nil, err
	}

	d := &Downloader {
		Client: client,
		Url: url,
		OutputDir: outputDir,
		Delay: time.Duration(delay) * time.Millisecond,
	}

	return d, nil
}

func (d *Downloader) GetPageRaw(url string, headers map[string] string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	res, err := d.Client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (d *Downloader) GetPage(url string, headers map[string] string) (string, error) {
	res, err := d.GetPageRaw(url, headers)
	if err != nil {
		return "", err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
