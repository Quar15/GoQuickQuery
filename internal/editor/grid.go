package editor

import (
	"bufio"
	"errors"
	"io"
	"log/slog"
	"os"
	"strings"
	"sync"

	"github.com/quar15/qq-go/internal/assets"
)

type Grid struct {
	mu        sync.RWMutex
	Text      []string
	Rows      int32
	Cols      []int32
	Highlight [][]HighlightColorEnum
	MaxCol    int32
}

func NewGrid() *Grid {
	highlight := make([][]HighlightColorEnum, 1)
	highlight = append(highlight, make([]HighlightColorEnum, 0))

	return &Grid{
		Text:      []string{""},
		Rows:      0,
		Cols:      []int32{0},
		Highlight: highlight,
		MaxCol:    0,
	}
}

func (eg *Grid) Lock() {
	eg.mu.Lock()
}

func (eg *Grid) Unlock() {
	eg.mu.Unlock()
}

func (eg *Grid) FakeInit(appAssets *assets.Assets) {
	eg.Text = []string{
		"SELECT * FROM example LIMIT 500;",
		"",
		"UPDATE example SET x = 1 WHERE id = 2;",
		"SELECT * FROM \"public\".\"example\" WHERE name = 'xyz';",
	}
	eg.Rows = int32(len(eg.Text))
	eg.Highlight = make([][]HighlightColorEnum, eg.Rows)
	var maxCol int32 = 0
	for row := int32(0); row < eg.Rows; row++ {
		lineLen := int32(len(eg.Text[row]))
		eg.Cols = append(eg.Cols, int32(lineLen))
		if maxCol < lineLen {
			maxCol = lineLen
		}

		if lineLen > 0 {
			eg.Highlight = append(eg.Highlight, make([]HighlightColorEnum, lineLen))
		} else {
			eg.Highlight = append(eg.Highlight, nil)
		}
	}
	eg.UpdateHighlight(0, eg.Rows-1)
}

func LoadGridFromTextFile(path string, appAssets *assets.Assets) (*Grid, error) {
	var eg *Grid = &Grid{}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	eg.Text = make([]string, 0, 256)
	eg.Rows = 0
	eg.Highlight = make([][]HighlightColorEnum, 0, 256)
	var maxCol int32 = 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			} else {
				return nil, errors.New("Failed to parse provided file")
			}
		}
		line = strings.TrimRight(line, "\n")
		eg.Text = append(eg.Text, line)
		eg.Rows++
		lineLen := int32(len(line))
		eg.Cols = append(eg.Cols, int32(lineLen))
		if maxCol < lineLen {
			maxCol = lineLen
			eg.MaxCol = lineLen
		}

		if lineLen > 0 {
			eg.Highlight = append(eg.Highlight, make([]HighlightColorEnum, lineLen))
		} else {
			eg.Highlight = append(eg.Highlight, nil)
		}
	}

	eg.UpdateHighlight(0, eg.Rows-1)
	return eg, nil
}

func (eg *Grid) DetectQueryRowsBoundaryBasedOnRow(row int32) (start int32, end int32) {
	start = row
	end = row

	for start > 0 {
		if eg.Cols[start-1] == 0 {
			break
		}
		start--
	}

	for end < eg.Rows-1 {
		if eg.Cols[end+1] == 0 {
			break
		}
		end++
	}

	return start, end
}

func (eg *Grid) InsertChar(row, col int32, ch rune) int32 {
	eg.mu.Lock()
	defer eg.mu.Unlock()
	line := eg.Text[row]
	eg.Text[row] = line[:col] + string(ch) + line[col:]
	eg.Cols[row]++
	eg.UpdateHighlight(row, row)
	eg.recalculateMaxCol(row)
	return col + 1
}

func (eg *Grid) DeleteCharBefore(row, col int32) (newRow, newCol int32) {
	eg.mu.Lock()
	slog.Debug("Delecting character before", slog.Int("row", int(row)), slog.Int("col", int(col)))
	defer eg.mu.Unlock()
	if col > 0 {
		line := eg.Text[row]
		eg.Text[row] = line[:col-1] + line[col:]
		eg.Cols[row]--

		eg.UpdateHighlight(row, row)
		eg.recalculateMaxCol(row)
		return row, col - 1
	}

	if row > 0 {
		newRow = row - 1
		eg.joinLines(newRow)
		newCol = eg.Cols[newRow] - 1
		return newRow, newCol
	}

	return row, col
}

func (eg *Grid) joinLines(row int32) {
	if row >= eg.Rows-1 {
		return
	}

	nextRow := row + 1
	eg.Text[row] = eg.Text[row] + eg.Text[nextRow]
	eg.Cols[row] = int32(len(eg.Text[row]))

	// Remove next line
	eg.Text = append(eg.Text[:nextRow], eg.Text[nextRow+1:]...)
	eg.Cols = append(eg.Cols[:nextRow], eg.Cols[nextRow+1:]...)
	eg.Highlight = append(eg.Highlight[:nextRow], eg.Highlight[nextRow+1:]...)
	eg.Rows--

	eg.UpdateHighlight(row, row)
	eg.recalculateMaxCol(row)
}

func (eg *Grid) InsertNewLine(row, col int32) (newRow, newCol int32) {
	eg.mu.Lock()
	defer eg.mu.Unlock()
	line := eg.Text[row]
	before := line[:col]
	after := line[col:]

	eg.Text[row] = before
	eg.Cols[row] = int32(len(before))

	newRowIdx := row + 1

	eg.Text = append(eg.Text, "")
	eg.Cols = append(eg.Cols, 0)
	eg.Highlight = append(eg.Highlight, nil)

	copy(eg.Text[newRowIdx+1:], eg.Text[newRowIdx:])
	copy(eg.Cols[newRowIdx+1:], eg.Cols[newRowIdx:])
	copy(eg.Highlight[newRowIdx+1:], eg.Highlight[newRowIdx:])

	eg.Text[newRowIdx] = after
	eg.Cols[newRowIdx] = int32(len(after))
	eg.Highlight[newRowIdx] = nil

	eg.Rows++
	eg.UpdateHighlight(row, newRowIdx)
	eg.recalculateMaxCol(row)

	return newRowIdx, 0
}

func (eg *Grid) recalculateMaxCol(row int32) {
	if eg.Cols[row] > eg.MaxCol {
		eg.MaxCol = eg.Cols[row]
		return
	}

	eg.MaxCol = row
	for _, col := range eg.Cols {
		if col > eg.MaxCol {
			eg.MaxCol = col
		}
	}
}
