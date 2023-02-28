package snippeter

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const text = `
Sentence 1 marker, other words.
Sentence 2 marker, other words.
Sentence 3 marker, other words.
Sentence 4 marker, other words.
Sentence 5 marker, other words.`

func TestComposeSnippetWithContext(t *testing.T) {
	cases := []struct {
		markerStart int
		markerEnd   int
		count       int
		want        string
	}{
		{
			markerStart: strings.Index(text, "3 marker"),
			markerEnd:   strings.Index(text, "3 marker") + len("3 marker"),
			count:       1,
			want:        `Sentence 3 marker, other words.`,
		},
		{
			markerStart: strings.Index(text, "3 marker"),
			markerEnd:   strings.Index(text, "3 marker") + len("3 marker"),
			count:       2,
			want: `Sentence 2 marker, other words.
Sentence 3 marker, other words.
Sentence 4 marker, other words.`,
		},
		{
			markerStart: strings.Index(text, "1 marker"),
			markerEnd:   strings.Index(text, "1 marker") + len("1 marker"),
			count:       1,
			want:        `Sentence 1 marker, other words.`,
		},
		{
			markerStart: strings.Index(text, "1 marker"),
			markerEnd:   strings.Index(text, "1 marker") + len("1 marker"),
			count:       2,
			want: `Sentence 1 marker, other words.
Sentence 2 marker, other words.`,
		},
		{
			markerStart: strings.Index(text, "5 marker"),
			markerEnd:   strings.Index(text, "5 marker") + len("5 marker"),
			count:       1,
			want:        `Sentence 5 marker, other words.`,
		},
		{
			markerStart: strings.Index(text, "5 marker"),
			markerEnd:   strings.Index(text, "5 marker") + len("5 marker"),
			count:       2,
			want: `Sentence 4 marker, other words.
Sentence 5 marker, other words.`,
		},
		{
			markerStart: strings.Index(text, "5 marker"),
			markerEnd:   strings.Index(text, "5 marker") + len("5 marker"),
			count:       10,
			want:        text,
		},
		{
			markerStart: strings.Index(text, "1 marker"),
			markerEnd:   strings.Index(text, "1 marker") + len("1 marker"),
			count:       10,
			want:        text,
		},
	}

	for _, tc := range cases {
		t.Run("", func(t *testing.T) {
			// act
			start, end := composeSnippetWithContext([]byte(text), tc.markerStart, tc.markerEnd, tc.count)

			// assert
			require.Equal(t, strings.TrimSpace(tc.want), strings.TrimSpace(text[start:end]))
		})
	}
}
