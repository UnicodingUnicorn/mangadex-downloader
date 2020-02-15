package main

import (
	"errors"
	"math"
	"sort"
	"strconv"
	"strings"
)

type Range struct {
	Start float64
	End float64
}

func InRange(ranges []*Range, n float64) bool {
	for _, r := range ranges {
		if n >= r.Start && n <= r.End {
			return true
		}
	}

	return false
}

func NewRange(start float64, end float64, masterRange *Range) (*Range, error) {
	if start <= 0.0 || end <= 0.0 {
		return nil, errors.New("term cannot be negative or zero!")
	}

	if start < masterRange.Start {
		return nil, errors.New("supplied start cannot be lower than the smallest chapter")
	}

	if end > masterRange.End {
		return nil, errors.New("supplied end cannot be higher than the largest chapter")
	}

	if start > end {
		tmp := end
		end = start
		start = tmp
	}

	r := &Range {
		Start: start,
		End: end,
	}

	return r, nil
}

func ParseRange(rangeString string, masterRange *Range) ([]*Range, error) {
	numbers := strings.Split(rangeString, ",")
	if len(numbers) == 0 {
		return nil, errors.New("no range string specified")
	}

	ranges := make([]*Range, 0)
	for _, number := range numbers {
		parts := strings.Split(strings.TrimSpace(number), "-")
		if len(parts) == 0 {
			continue
		} else if len(parts) == 1 {
			n, err := strconv.ParseFloat(parts[0], 64)
			if err != nil {
				return nil, err
			}

			r, err := NewRange(n, n, masterRange)
			if err != nil {
				return nil, err
			}

			ranges = append(ranges, r)
		} else if len(parts) == 2 {
			start, err := strconv.ParseFloat(parts[0], 64)
			if err != nil {
				return nil, err
			}

			end, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				return nil, err
			}

			r, err := NewRange(start, end, masterRange)
			if err != nil {
				return nil, err
			}

			ranges = append(ranges, r)
		} else {
			return nil, errors.New("ranges can only have two terms!")
		}
	}

	return ranges, nil
}

func (c *ChapterPageData) GetChapterNumbers() []string {
	numbers := c.GetChapterNumbersRaw()
	numberStrings := make([]string, 0)

	for _, n := range numbers {
		numberStrings = append(numberStrings, strconv.FormatFloat(n, 'f', -1, 64))
	}

	return numberStrings
}

func (c *ChapterPageData) GetChapterNumbersRaw() []float64 {
	unique := make(map[float64]bool)
	for _, data := range c.Metadata {
		unique[data.Chapter] = true
	}

	numbers := make([]float64, 0)
	for k, v := range unique {
		if v == true {
			numbers = append(numbers, k)
		}
	}

	sort.Float64s(numbers)

	return numbers
}

func (c *ChapterPageData) GetChapterRange() (*Range, error) {
	start := math.MaxFloat64
	end := -1.0

	numbers := c.GetChapterNumbersRaw()

	for _, chapter := range numbers {
		if chapter > end {
			end = chapter
		}
		if chapter < start {
			start = chapter
		}
	}

	if start > end {
		return nil, errors.New("invalid chapter range, start larger than end")
	}

	r := &Range {
		Start: math.Floor(start),
		End: math.Ceil(end),
	}

	return r, nil
}

