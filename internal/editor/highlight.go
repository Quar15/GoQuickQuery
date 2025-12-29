package editor

import (
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
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

var highlightColor = map[HighlightColorEnum]rl.Color{}

func InitHighlightColors(cfg *config.Config) {
	highlightColor = map[HighlightColorEnum]rl.Color{
		// @TODO: Consider moving to color config
		HighlightKeyword:     cfg.Colors.Mauve(),
		HighlightFunction:    cfg.Colors.Blue(),
		HighlightDatabaseVar: cfg.Colors.Yellow(),
		HighlightText:        cfg.Colors.Green(),
		HighlightNumber:      cfg.Colors.Peach(),
		HighlightNormal:      cfg.Colors.Text(),
	}
}

func (hc HighlightColorEnum) Color() rl.Color {
	return highlightColor[hc]
}

func (eg *Grid) UpdateHighlight(fromRow int32, toRow int32) {
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

func isWordChar(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9') ||
		c == '_'
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}
