package editor

import "strings"

type Match struct {
	StartRow int
	StartCol int
	EndRow   int
	EndCol   int
}

func (eg *Grid) FindAll(word string) []Match {
	if word == "" {
		return nil
	}

	// Build flat string and index map
	var flat strings.Builder
	indexToPos := make([][2]int, 0)

	for r, line := range eg.Text {
		for c, ch := range line {
			flat.WriteRune(ch)
			indexToPos = append(indexToPos, [2]int{r, c})
		}

		// Add virtual newline (not selectable but used for spanning lines)
		if r < len(eg.Text)-1 {
			flat.WriteByte('\n')
			indexToPos = append(indexToPos, [2]int{r, len(line)})
		}
	}

	data := flat.String()
	var results []Match

	searchPos := 0
	for {
		i := strings.Index(data[searchPos:], word)
		if i == -1 {
			break
		}

		i += searchPos
		start := i
		end := i + len(word) - 1

		sr, sc := indexToPos[start][0], indexToPos[start][1]
		er, ec := indexToPos[end][0], indexToPos[end][1]

		results = append(results, Match{
			StartRow: sr,
			StartCol: sc,
			EndRow:   er,
			EndCol:   ec,
		})

		searchPos = i + 1
	}

	return results
}
