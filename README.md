# mangadex-downloader

A simple program to download comics from [MangaDex](https://mangadex.org), because I am a compulsive hoarder. Also, [Polite Mangadex Downloader](http://sizer99.com/pmd) is more complicated than it is worth and [mangadex-dl](https://github.com/frozenpandaman/mangadex-dl) has the annoying tendency not to render chapter titles properly.

Try not to set the delay too low or run too many instances at once or read things at the same thing and risk getting IP-banned. Don't be a jerk to the nice people providing this service.

## Usage

The program is written in [Go](https://golang.org), so just set that up and run `go build`. Alternatively, for Windows binaries, look in releases.

In CLI, run:

```
./mangadex-downloader -u <url> -o comic
```

URL in this case is the url of the title page of the comic. Ranges are something like 1, 3-4. Type `-h` for the (rather sparse) help.
