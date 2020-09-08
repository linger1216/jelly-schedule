package core

import (
	"errors"
	"strconv"
	"strings"
)

type Pattern struct {
	segs               []string
	holders            []*PlaceHolder
	replaces           []string
	arrangeHolderIdxes []int
	splitHolderIdxes   []int
	key                string
	keys               []string
}

const (
	blank = ' '
)

var (
	EParse       = errors.New("parse pattern failed")
	EPlaceHolder = errors.New("unknown placeholder format")
	EDate        = errors.New("not a date")
	ENumber      = errors.New("not a number")
	EAlpha       = errors.New("not a alpha")
)

func defaultKeyGen(s ...string) string {
	return strings.Join(s, "-")
}

func ParsePattern(pattern string) (*Pattern, error) {
	var pat Pattern

	prev := blank
	prevIdx := 0
	bs := make([]rune, 0, 128)
	for idx, c := range pattern {
		switch c {
		case '[', '<':
			if prev != blank {
				return nil, EParse
			}
			pat.segs = append(pat.segs, string(bs))
			bs = bs[:0]
			prev = c
			prevIdx = idx
		case ']':
			if prev != '[' {
				return nil, EParse
			}
			holder, err := NewPlaceHolder(pattern, prevIdx, idx)
			if err == nil {
				pat.arrangeHolderIdxes = append(pat.arrangeHolderIdxes, len(pat.holders))
				pat.holders = append(pat.holders, holder)
				//arrangeHolders = append(arrangeHolders, holder)
			}
			prev = blank
		case '>':
			if prev != '<' {
				return nil, EParse
			}
			holder, err := NewPlaceHolder(pattern, prevIdx, idx)
			if err == nil {
				pat.splitHolderIdxes = append(pat.splitHolderIdxes, len(pat.holders))
				pat.holders = append(pat.holders, holder)
				//splitHolders = append(splitHolders, holder)
			}
			prev = blank
		default:
			if prev == blank {
				bs = append(bs, c)
			}
		}
	}
	pat.segs = append(pat.segs, string(bs))
	return &pat, nil
}

func (p *Pattern) Arrange() []string {
	patterns := p.arrange(nil)

	ret := make([]string, 0, len(patterns))
	for _, pattern := range patterns {
		ret = append(ret, pattern.ToString())
	}

	return ret
}

func (p *Pattern) Map(keyGen func(...string) string) map[string][]string {
	patterns := p.arrange(keyGen)

	ret := map[string][]string{}
	for _, pattern := range patterns {
		key := pattern.key
		ret[key] = append(ret[key], pattern.ToString())
	}
	return ret
}

func (p *Pattern) ToString() string {
	ret := p.segs[0]

	for i := 0; i < len(p.replaces); i++ {
		ret += p.replaces[i] + p.segs[i+1]
	}
	return ret
}

func (p *Pattern) arrange(keyGen func(...string) string) []*Pattern {
	patterns := []*Pattern{p}

	for _, idx := range p.splitHolderIdxes {
		tmp := make([]*Pattern, 0, len(patterns)*len(p.holders[idx].elements))
		for _, pattern := range patterns {
			tmp = append(tmp, pattern.arrangeAt(idx)...)
		}
		patterns = tmp
	}

	for _, pattern := range patterns {
		if keyGen == nil {
			pattern.key = strings.Join(pattern.keys, "")
		} else {
			pattern.key = keyGen(pattern.keys...)
		}

	}

	for _, idx := range p.arrangeHolderIdxes {
		tmp := make([]*Pattern, 0, len(patterns)*len(p.holders[idx].elements))
		for _, pattern := range patterns {
			tmp = append(tmp, pattern.arrangeAt(idx)...)
		}
		patterns = tmp
	}
	return patterns
}

func (p *Pattern) arrangeAt(i int) []*Pattern {
	holder := p.holders[i]
	patterns := make([]*Pattern, 0, len(holder.elements))
	for _, e := range holder.elements {
		pattern := &Pattern{
			segs:     p.segs,
			holders:  p.holders,
			replaces: make([]string, len(p.holders)),
			keys:     make([]string, len(p.keys)),
			key:      p.key,
		}
		copy(pattern.keys, p.keys)
		copy(pattern.replaces, p.replaces)
		pattern.keys = append(pattern.keys, e)
		pattern.replaces[i] = e
		patterns = append(patterns, pattern)
	}
	return patterns
}

type PlaceHolder struct {
	elements []string
}

func NewPlaceHolder(str string, begin, end int) (*PlaceHolder, error) {
	str = str[begin+1 : end]
	elements := strings.Split(str, ",")
	if len(elements) > 1 {
		return &PlaceHolder{
			elements: elements,
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
	}, nil
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
