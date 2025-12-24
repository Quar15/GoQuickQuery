package display

import (
	"bufio"
	"errors"
	"os"
	"strings"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/quar15/qq-go/internal/assets"
	"github.com/quar15/qq-go/internal/config"
	"github.com/quar15/qq-go/internal/database"
)

type HighlightColorEnum int8

const (
	HighlightKeyword HighlightColorEnum = iota
	HighlightFunction
	HighlightDatabaseVar // tablespaces, tables
	HighlightText
	HighlightNumber
	HighlightNormal
)

var setupHighlightColorOnce sync.Once
var highlightColor = map[HighlightColorEnum]rl.Color{}

func (hc HighlightColorEnum) Color() rl.Color {
	setupHighlightColorOnce.Do(func() {
		highlightColor = map[HighlightColorEnum]rl.Color{
			// @TODO: Consider moving to color config
			HighlightKeyword:     config.Get().Colors.Mauve(),
			HighlightFunction:    config.Get().Colors.Blue(),
			HighlightDatabaseVar: config.Get().Colors.Yellow(),
			HighlightText:        config.Get().Colors.Green(),
			HighlightNumber:      config.Get().Colors.Peach(),
			HighlightNormal:      config.Get().Colors.Text(),
		}

	})
	return highlightColor[hc]
}

type EditorGrid struct {
	Text      []string
	Rows      int32
	Cols      []int32
	Highlight [][]HighlightColorEnum
	MaxCol    int8
	MaxWidth  float32
}

func NewEditorGrid() EditorGrid {
	return EditorGrid{
		Text:      []string{},
		Rows:      0,
		Cols:      []int32{0},
		Highlight: [][]HighlightColorEnum{},
		MaxCol:    0,
		MaxWidth:  0,
	}
}

func (eg *EditorGrid) UpdateHighlight(fromRow int32, toRow int32) {
	for row := fromRow; row <= toRow; row++ {
		if eg.Cols[row] == 0 {
			eg.Highlight[row] = nil
			continue
		}

		if cap(eg.Highlight[row]) < int(eg.Cols[row]) {
			eg.Highlight[row] = make([]HighlightColorEnum, eg.Cols[row])
		} else {
			eg.Highlight[row] = eg.Highlight[row][:eg.Cols[row]]
		}

		// @TODO: Highlight
		for i := range eg.Highlight[row] {
			eg.Highlight[row][i] = HighlightNormal
		}

		i := int32(0)
		line := eg.Text[row]
		for i < eg.Cols[row] {
			c := line[i]
			// Text
			if c == '\'' {
				eg.Highlight[row][i] = HighlightText
				i++
				for i < eg.Cols[row] {
					eg.Highlight[row][i] = HighlightText
					if line[i] == '\'' {
						i++
						break
					}
					i++
				}
			}

			// DB Variables
			if c == '"' {
				eg.Highlight[row][i] = HighlightDatabaseVar
				i++
				for i < eg.Cols[row] {
					eg.Highlight[row][i] = HighlightDatabaseVar
					if line[i] == '"' {
						i++
						break
					}
					i++
				}
			}

			// Digits
			if isDigit(c) {
				start := i
				for i < eg.Cols[row] && isDigit(line[i]) {
					i++
				}
				for j := start; j < i; j++ {
					eg.Highlight[row][j] = HighlightNumber
				}
				continue
			}

			if isWordChar(c) {
				start := i
				for i < eg.Cols[row] && isWordChar(line[i]) {
					i++
				}

				word := strings.ToLower(line[start:i])

				var color HighlightColorEnum = HighlightNormal

				if _, ok := database.SqlKeywords[word]; ok {
					color = HighlightKeyword
				} else if _, ok := database.PostgresqlFunctionsKeywords[word]; ok {
					color = HighlightFunction
				}

				if color != HighlightNormal {
					for j := start; j < i; j++ {
						eg.Highlight[row][j] = color
					}
				}

				continue
			}

			i++
		}

	}
}

func (eg *EditorGrid) FakeInit(appAssets *assets.Assets) {
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
			eg.MaxWidth = appAssets.MeasureTextMainFont(eg.Text[row]).X
		}

		if lineLen > 0 {
			eg.Highlight = append(eg.Highlight, make([]HighlightColorEnum, lineLen))
		} else {
			eg.Highlight = append(eg.Highlight, nil)
		}
	}
	eg.UpdateHighlight(0, eg.Rows-1)
}

func LoadEditorGridFromTextFile(path string, appAssets *assets.Assets) (*EditorGrid, error) {
	var eg *EditorGrid = &EditorGrid{}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	eg.Text = make([]string, 0, 256)
	eg.Rows = 0
	eg.MaxWidth = 0
	eg.Highlight = make([][]HighlightColorEnum, 0, 256)
	var eof bool = false
	var maxCol int32 = 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				eof = true
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
			eg.MaxCol = int8(lineLen)
			eg.MaxWidth = appAssets.MeasureTextMainFont(line).X
		}

		if lineLen > 0 {
			eg.Highlight = append(eg.Highlight, make([]HighlightColorEnum, lineLen))
		} else {
			eg.Highlight = append(eg.Highlight, nil)
		}

		if eof {
			break
		}
	}

	eg.UpdateHighlight(0, eg.Rows-1)
	return eg, nil
}

func (eg *EditorGrid) DetectQueryRowsBoundaryBasedOnRow(row int32) (start int32, end int32) {
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

func isWordChar(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9') ||
		c == '_'
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}
