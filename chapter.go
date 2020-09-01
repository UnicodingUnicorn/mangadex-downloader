package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type ChapterPageData struct {
	Metadata		[]*ChapterMetadata
	Languages		map[string]string
	Title				string
}

type ChapterMetadata struct {
	Id			string
	Title		string
	Chapter	float64
	Volume	string
	Lang		string
}

func (m *ChapterMetadata) FormattedTitle() string {
	titleBits := make([]string, 0)

	if m.Volume != "" {
		titleBits = append(titleBits, fmt.Sprintf("Vol. %s", m.Volume))
	}

	if m.Chapter != 0.0 {
		titleBits = append(titleBits, fmt.Sprintf("Ch. %s", strconv.FormatFloat(m.Chapter, 'f', -1, 64)))
	}

	if m.Title != "" {
		titleBits = append(titleBits, m.Title)
	}

	return strings.Join(titleBits, " ")
}

func (d *Downloader) GetChapterData() (*ChapterPageData, error) {
	metadata := make([]*ChapterMetadata, 0)
	var languages map[string]string
	var title string

	i := 1
	for {
		p, err := d.GetPage(fmt.Sprintf("%s/chapters/%d", d.Url, i), nil)
		if err != nil {
			return nil, err
		}

		page := string(p)

		doc, err := html.Parse(strings.NewReader(page))
		if err != nil {
			return nil, err
		}

		m := GetChaptersData(doc)
		languages = GetLanguages(doc)
		title = GetTitle(doc)

		if len(m) == 0 {
			break
		}

		metadata = append(metadata, m...)
		i++
		time.Sleep(d.Delay)
	}

	if len(metadata) == 0 {
		return nil, errors.New("number of retrieved chapters is 0, are you sure you have input the correct url?")
	}

	data := &ChapterPageData {
		Metadata: metadata,
		Languages: languages,
		Title: title,
	}

	return data, nil
}

func GetChaptersData(n *html.Node) []*ChapterMetadata {
	ChaptersDataSteps := GetChaptersDataSteps()
	ChapterDataSteps := GetChapterDataSteps()

	metadata := make([]*ChapterMetadata, 0)

	step(n, 0, ChaptersDataSteps, func (n2 *html.Node) {
		i := 0
		for c := n2.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode && c.Data == "div" && GetAttr(c, "class") == "row no-gutters" {
				if i > 0 {
					step(c, 0, ChapterDataSteps, func(n3 *html.Node) {
						chapter := GetAttr(n3, "data-chapter")
						chapterFloat, _ := strconv.ParseFloat(chapter, 64)

						m := &ChapterMetadata {
							Id: GetAttr(n3, "data-id"),
							Title: GetAttr(n3, "data-title"),
							Chapter: chapterFloat,
							Volume: GetAttr(n3, "data-volume"),
							Lang: GetAttr(n3, "data-lang"),
						}

						metadata = append(metadata, m)
					})
				}
				i++
			}
		}
	})

	return metadata
}

func GetLanguages(n *html.Node) map[string]string {
	LanguageStepsA := GetLanguageStepsA()
	LanguageStepsB := GetLanguageStepsB()

	languages := make(map[string] string)
	step(n, 0, LanguageStepsA, func (n2 *html.Node) {
		i := 0
		for c := n2.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode && c.Data == "div" {
				if i == 1 {
					step(c, 0, LanguageStepsB, func(n3 *html.Node) {
						for c2 := n3.FirstChild; c2 != nil; c2 = c2.NextSibling {
							if c2.Type == html.ElementNode && c2.Data == "option" {
								code := GetAttr(c2, "value")

								var name strings.Builder
								for c3 := c2.FirstChild; c3 != nil; c3 = c3.NextSibling {
									if c3.Type == html.TextNode {
										name.WriteString(c3.Data)
									}
								}

								languages[name.String()] = code
							}
						}
					})
					break
				}

				i++
			}
		}
	})

	return languages
}

func GetTitle(n *html.Node) string {
	TitleSteps := GetTitleSteps()

	var title strings.Builder

	step(n, 0, TitleSteps, func(n2 *html.Node) {
		for c := n2.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.TextNode {
				title.WriteString(c.Data)
			}
		}
	})

	return title.String()
}

func GetChaptersDataSteps() []*Step {
	steps := make([]*Step, 5)

	steps[0] = &Step {
		Element: "html",
		Id: "",
		Class: "",
	}

	steps[1] = &Step {
		Element: "body",
		Id: "",
		Class: "",
	}

	steps[2] = &Step {
		Element: "div",
		Id: "content",
		Class: "container",
	}

	steps[3] = &Step {
		Element: "div",
		Id: "",
		Class: "edit tab-content",
	}

	steps[4] = &Step {
		Element: "div",
		Id: "",
		Class: "chapter-container",
	}

	return steps
}

func GetChapterDataSteps() []*Step {
	steps := make([]*Step, 2)

	steps[0] = &Step {
		Element: "div",
		Id: "",
		Class: "col",
	}

	steps[1] = &Step {
		Element: "div",
		Id: "",
		Class: "chapter-row",
	}

	return steps
}

func GetLanguageStepsA() []*Step {
	steps := make([]*Step, 7)

	steps[0] = &Step {
		Element: "html",
		Id: "",
		Class: "",
	}

	steps[1] = &Step {
		Element: "body",
		Id: "",
		Class: "",
	}

	steps[2] = &Step {
		Element: "div",
		Id: "homepage_settings_modal",
		Class: "",
	}

	steps[3] = &Step {
		Element: "div",
		Id: "",
		Class: "",
	}

	steps[4] = &Step {
		Element: "div",
		Id: "",
		Class: "",
	}

	steps[5] = &Step {
		Element: "div",
		Id: "",
		Class: "modal-body",
	}

	steps[6] = &Step {
		Element: "form",
		Id: "",
		Class: "",
	}

	return steps
}

func GetLanguageStepsB() []*Step {
	steps := make([]*Step, 2)

	steps[0] = &Step {
		Element: "div",
		Id: "",
		Class: "",
	}

	steps[1] = &Step {
		Element: "select",
		Id: "",
		Class: "",
	}

	return steps
}

func GetTitleSteps() []*Step {
	steps := make([]*Step, 6)

	steps[0] = &Step {
		Element: "html",
		Id: "",
		Class: "",
	}

	steps[1] = &Step {
		Element: "body",
		Id: "",
		Class: "",
	}

	steps[2] = &Step {
		Element: "div",
		Id: "content",
		Class: "container",
	}

	steps[3] = &Step {
		Element: "div",
		Id: "",
		Class: "card mb-3",
	}

	steps[4] = &Step {
		Element: "h6",
		Id: "",
		Class: "",
	}

	steps[5] = &Step {
		Element: "span",
		Id: "",
		Class: "mx-1",
	}

	return steps
}
