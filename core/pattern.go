package core

import (
	"errors"
	"strconv"
	"strings"
)

const (
	blank = ' '
)

var (
	ESplitFieldName = errors.New("no split field name specified")
	EString         = errors.New("is not string")
	ESplitField     = errors.New("specified split field does not exist")
	EParse          = errors.New("parse pattern failed")
	EPlaceHolder    = errors.New("unknown placeholder format")
	EDate           = errors.New("not a date")
	ENumber         = errors.New("not a number")
	EAlpha          = errors.New("not a alpha")
)

type SplitConfig struct {
	SplitFieldName string `json:"splitFieldName,omitempty" yaml:"splitFieldName" `
}

type SplitJob struct {
	splitFieldName string
}

type PlaceHolder struct {
	elements []string
	begin    int
	end      int
}

func NewPlaceHolder(str string, begin, end int) (*PlaceHolder, error) {
	str = str[begin+1 : end]
	elements := strings.Split(str, ",")
	if len(elements) > 1 {
		return &PlaceHolder{
			elements: elements,
			begin:    begin,
			end:      end,
		}, nil
	}
	elements = strings.Split(str, "-")
	if len(elements) != 2 {
		return nil, EPlaceHolder
	}
	eBegin := elements[0]
	eEnd := elements[1]
	elements = generateElements(eBegin, eEnd)
	if len(elements) == 0 {
		return nil, EPlaceHolder
	}
	return &PlaceHolder{
		elements: elements,
		begin:    begin,
		end:      end,
	}, nil
}

func (h *PlaceHolder) Replace(pattern string) ([]string, error) {
	ret := make([]string, 0, len(h.elements))

	for _, e := range h.elements {
		ret = append(ret, pattern[:h.begin]+e+pattern[h.end+1:])
	}
	return ret, nil
}

func (h *PlaceHolder) Size() int {
	return len(h.elements)
}

func generateElements(begin, end string) []string {
	var ret []string
	var err error

	ret, err = generateElementsDate(begin, end)
	if err == nil {
		return ret
	}

	ret, err = generateElementsNumber(begin, end)
	if err == nil {
		return ret
	}

	ret, err = generateElementsAlpha(begin, end)
	if err == nil {
		return ret
	}

	return nil
}

func generateElementsDate(begin, end string) ([]string, error) {
	return nil, EDate
}

func generateElementsNumber(begin, end string) ([]string, error) {
	beginNum, err := strconv.ParseInt(begin, 10, 64)
	if err != nil {
		return nil, ENumber
	}

	endNum, err := strconv.ParseInt(end, 10, 64)
	if err != nil {
		return nil, ENumber
	}

	ret := make([]string, 0, endNum-beginNum+1)
	for i := beginNum; i <= endNum; i++ {
		str := strconv.FormatInt(i, 10)
		ret = append(ret, str)
	}
	return ret, nil
}

func generateElementsAlpha(begin, end string) ([]string, error) {
	return nil, EAlpha
}

func generateStringsFromPattern(pattern string, n int) ([][]string, error) {
	splitHolders, arrangeHolders, err := parsePattern(pattern)
	if err != nil {
		return nil, err
	}

	splitPatterns, err := arrangePattern([]string{pattern}, splitHolders)
	if err != nil {
		return nil, err
	}

	ret := make([][]string, n)
	sliceSplitPatterns := splitPattern(splitPatterns, n)
	for i, patterns := range sliceSplitPatterns {
		if arrangePatterns, err := arrangePattern(patterns, arrangeHolders); err == nil {
			ret[i] = arrangePatterns
		}
	}
	return ret, nil
}

func parsePattern(pattern string) ([]*PlaceHolder, []*PlaceHolder, error) {
	var splitHolders []*PlaceHolder
	var arrangeHolders []*PlaceHolder

	prev := blank
	prevIdx := 0
	for idx, c := range pattern {
		switch c {
		case '[', '<':
			if prev != blank {
				return nil, nil, EParse
			}
			prev = c
			prevIdx = idx
		case ']':
			if prev != '[' {
				return nil, nil, EParse
			}
			holder, err := NewPlaceHolder(pattern, prevIdx, idx)
			if err == nil {
				arrangeHolders = append(arrangeHolders, holder)
			}
			prev = blank
		case '>':
			if prev != '<' {
				return nil, nil, EParse
			}
			holder, err := NewPlaceHolder(pattern, prevIdx, idx)
			if err == nil {
				splitHolders = append(splitHolders, holder)
			}
			prev = blank
		}
	}
	return splitHolders, arrangeHolders, nil
}

func arrangePattern(patterns []string, holds []*PlaceHolder) ([]string, error) {
	ret := patterns
	var tmp []string
	for _, hold := range holds {
		tmp = make([]string, 0, len(ret)*hold.Size())
		for _, str := range ret {
			str2, err := hold.Replace(str)
			if err != nil {
				continue
			}
			tmp = append(tmp, str2...)
		}
		ret = tmp
	}
	return ret, nil
}

func splitPattern(patterns []string, n int) [][]string {
	ret := make([][]string, n)

	for idx, pattern := range patterns {
		part := idx % n
		ret[part] = append(ret[part], pattern)
	}

	return ret
}
