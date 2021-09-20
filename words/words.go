package words

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// 参考: https://github.com/polaris1119/wordscount/blob/master/count.go

// WordsCount 字数统计模型
type WordsCount struct {
	// Word 不包括标点符号
	Word int32
	// Punctuation 标点符号数量
	Punctuation int32
}

// Pure 字数，不包括标点符号
func (wc *WordsCount) Pure() int32 {
	return wc.Word
}

// Total 总字数，包括标点符号
func (wc *WordsCount) Total() int32 {
	return wc.Word + wc.Punctuation
}

// Total 总字数

// Count 统计字数
func Count(content string) *WordsCount {

	wc := new(WordsCount)

	c1 := autoSpace(content)
	c2 := strings.Fields(c1)

	for _, w := range c2 {
		words := strings.FieldsFunc(w, func(r rune) bool {
			if unicode.IsPunct(r) {
				wc.Punctuation++
				return true
			}

			return false
		})

		for _, word := range words {
			runeCount := utf8.RuneCountInString(word)
			if len(word) == runeCount {
				wc.Word++
			} else {
				wc.Word += int32(runeCount)
			}
		}
	}

	return wc
}

// autoSpace 自动给中英文之间加上空格
func autoSpace(str string) string {
	out := ""

	for _, r := range str {
		out = addSpaceAtBoundary(out, r)
	}

	return out
}

func addSpaceAtBoundary(prefix string, nextChar rune) string {
	if len(prefix) == 0 {
		return string(nextChar)
	}

	r, size := utf8.DecodeLastRuneInString(prefix)
	if isLatin(size) != isLatin(utf8.RuneLen(nextChar)) &&
		isAllowSpace(nextChar) && isAllowSpace(r) {
		return prefix + " " + string(nextChar)
	}

	return prefix + string(nextChar)
}

func isLatin(size int) bool {
	return size == 1
}

func isAllowSpace(r rune) bool {
	return !unicode.IsSpace(r) && !unicode.IsPunct(r)
}
