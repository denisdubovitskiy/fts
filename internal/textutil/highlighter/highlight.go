package highlighter

import "bytes"

type Highlight struct {
	Start int
	End   int
}

func FindAll(opening, closing string, data []byte) []Highlight {
	var (
		shift   = 0
		prevIdx = 0
		next    = opening
	)

	results := make([]Highlight, 0, 0)

	for shift < len(data) {
		phraseStartIdx := bytes.Index(data[shift:], []byte(next))
		if phraseStartIdx < 0 {
			break
		}

		if next == opening {
			next = closing
			prevIdx = phraseStartIdx + shift
		} else {
			results = append(results, Highlight{
				Start: prevIdx,
				End:   phraseStartIdx + shift,
			})
			next = opening
			prevIdx = 0
		}

		shift += phraseStartIdx + len(next)
	}

	return results
}
