package snippeter

import (
	"bytes"
	"strings"

	"github.com/denisdubovitskiy/fts/internal/textutil/highlighter"
)

type Options struct {
	// Will be returned to the client
	HighlightStart string
	HighlightEnd   string

	// Returned by the database highlighter
	SourceStart string
	SourceEnd   string
}

func ParseToSnippets(given string, opts Options) []string {
	results := make([]string, 0, 0)
	content := []byte(given)

	// Find all highlights which are highlighted by the datbase
	highlights := highlighter.FindAll(
		opts.SourceStart,
		opts.SourceEnd,
		content,
	)

	for _, highlight := range highlights {
		// Add a context to each highlight, adding sentences before and after.
		startIndex, endIndex := composeSnippetWithContext(
			content,
			highlight.Start,
			highlight.End,
			3,
		)

		sentence := given[startIndex:endIndex]
		sentence = strings.ReplaceAll(sentence, "_hl_start_", opts.HighlightStart)
		sentence = strings.ReplaceAll(sentence, "_hl_end_", opts.HighlightEnd)
		sentence = strings.TrimSpace(sentence)

		if len(results) > 0 {
			if results[len(results)-1] == sentence {
				continue
			}
		}
		results = append(results, sentence)
	}

	return results
}

func composeSnippetWithContext(content []byte, start, end, sentencesBeforeAndAfter int) (int, int) {
	startIndex := findPrevDotIndex(content, start, sentencesBeforeAndAfter)
	endIndex := findNextDotIndex(content, end, sentencesBeforeAndAfter)
	return startIndex, endIndex
}

var dot = []byte(".")

func findNextDotIndex(body []byte, afterIdx, count int) int {
	index := 0
	for count > 0 {
		if afterIdx > len(body)-1 {
			return index
		}
		index = afterIdx + bytes.Index(body[afterIdx:], dot) + 1
		afterIdx = index
		count--
	}
	return index
}

func findPrevDotIndex(content []byte, start, count int) int {
	index := 0

	for count > 0 {
		if start < 0 {
			return index
		}
		index = bytes.LastIndex(content[:start], dot) + 1
		start = index - 1
		count--
	}

	return index
}
