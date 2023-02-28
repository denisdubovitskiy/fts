package highlighter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindAll(t *testing.T) {
	const textWithMarkers = `
Sentence 1 _hl_start_marker1_hl_end_, other words.
Sentence 2 _hl_start_marker2_hl_end_, other words.
Sentence 3 _hl_start_marker3_hl_end_, other words.
Sentence 4 _hl_start_marker4_hl_end_, other words.
Sentence 5 _hl_start_marker5_hl_end_, other words.
Sentence 6 _hl_start_marker6_hl_end_, other words.
Sentence 7 _hl_start_marker7_hl_end_, other words.`

	want := []string{
		"marker1",
		"marker2",
		"marker3",
		"marker4",
		"marker5",
		"marker6",
		"marker7",
	}

	// act
	highlights := FindAll("_hl_start_", "_hl_end_", []byte(textWithMarkers))

	// assert
	got := make([]string, len(highlights))
	for i, hl := range highlights {
		got[i] = textWithMarkers[hl.Start+len("_hl_start_") : hl.End]
	}

	require.Equal(t, want, got)
}
