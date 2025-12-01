package display

import (
	"fmt"
	"regexp"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/quar15/qq-go/internal/colors"
	"github.com/quar15/qq-go/internal/utilities"
)

type CursorMode int8

const (
	ModeNormal CursorMode = iota
	ModeInsert
	ModeVisual
	ModeVLine
	ModeCommand
)

var modeName = map[CursorMode]string{
	ModeNormal:  "NORMAL",
	ModeInsert:  "INSERT",
	ModeVisual:  "VISUAL",
	ModeVLine:   "V-LINE",
	ModeCommand: "COMMAND",
}

func (cm CursorMode) String() string {
	return modeName[cm]
}

var modeColor = map[CursorMode]rl.Color{
	ModeNormal:  colors.Blue(),
	ModeInsert:  colors.Green(),
	ModeVisual:  colors.Mauve(),
	ModeVLine:   colors.Mauve(),
	ModeCommand: colors.Peach(),
}

func (cm CursorMode) Color() rl.Color {
	return modeColor[cm]
}

type SpreadsheetCursor struct {
	Mode            CursorMode
	CmdBuf          string
	MotionBuf       string
	Logs            CommandLogs
	Col             int8
	Row             int32
	MaxCol          int8
	MaxRow          int32
	SelectStartCol  int8
	SelectStartRow  int32
	SelectEndCol    int8
	SelectEndRow    int32
	SelectAnchorCol int8
	SelectAnchorRow int32
}

func (c *SpreadsheetCursor) TransitionMode(newMode CursorMode) {
	c.Mode = newMode
}

func (c *SpreadsheetCursor) Reset() {
	c.TransitionMode(ModeNormal)
	c.SelectStartCol = -1
	c.SelectStartRow = -1
	c.SelectEndCol = -1
	c.SelectEndRow = -1
	c.CmdBuf = ""
	c.MotionBuf = ""
	c.UpdateCmdLine()
}

func (c *SpreadsheetCursor) Init() {
	c.Col = 0
	c.Row = 0
	c.Logs.Init()
	c.Reset()
}

func (c *SpreadsheetCursor) SetSelect(startCol int8, startRow int32, endCol int8, endRow int32) {
	c.SelectStartCol = startCol
	c.SelectStartRow = startRow
	c.SelectEndCol = endCol
	c.SelectEndRow = endRow
	c.SelectAnchorCol = startCol
	c.SelectAnchorRow = startRow
	if c.Mode != ModeVisual {
		c.TransitionMode(ModeVisual)
	}
}

func (c *SpreadsheetCursor) UpdateSelectBasedOnPosition() {
	switch c.Mode {
	case ModeVisual:
		if c.Col < c.SelectAnchorCol {
			c.SelectStartCol = c.Col
			c.SelectEndCol = c.SelectAnchorCol
		} else if c.Col > c.SelectAnchorCol {
			c.SelectStartCol = c.SelectAnchorCol
			c.SelectEndCol = c.Col
		} else {
			c.SelectStartCol = c.SelectAnchorCol
			c.SelectEndCol = c.SelectAnchorCol
		}

		if c.Row < c.SelectAnchorRow {
			c.SelectStartRow = c.Row
			c.SelectEndRow = c.SelectAnchorRow
		} else if c.Row > c.SelectAnchorRow {
			c.SelectStartRow = c.SelectAnchorRow
			c.SelectEndRow = c.Row
		} else {
			c.SelectStartRow = c.SelectAnchorRow
			c.SelectEndRow = c.SelectAnchorRow
		}
	case ModeVLine:
		if c.Row < c.SelectAnchorRow {
			c.SelectStartRow = c.Row
			c.SelectEndRow = c.SelectAnchorRow
		} else if c.Row > c.SelectAnchorRow {
			c.SelectStartRow = c.SelectAnchorRow
			c.SelectEndRow = c.Row
		} else {
			c.SelectStartRow = c.SelectAnchorRow
			c.SelectEndRow = c.SelectAnchorRow
		}
	}
}

func (c *SpreadsheetCursor) IsFocused(col int8, row int32) bool {
	return c.Col == col && c.Row == row
}

func (c *SpreadsheetCursor) ClampFocus(limitCol int8, limitRow int32) {
	c.MaxCol = limitCol
	c.MaxRow = limitRow
	c.Col = min(max(0, c.Col), c.MaxCol)
	c.Row = min(max(0, c.Row), c.MaxRow)
	if c.Mode == ModeVisual {
		c.SelectStartCol = min(max(0, c.SelectStartCol), c.MaxCol)
		c.SelectEndCol = min(max(0, c.SelectEndCol), c.MaxCol)
		c.SelectStartRow = min(max(0, c.SelectStartRow), c.MaxRow)
		c.SelectEndRow = min(max(0, c.SelectEndRow), c.MaxRow)
	}
}

func (c *SpreadsheetCursor) IsSelected(col int8, row int32) bool {
	if c.Mode != ModeVisual && c.Mode != ModeVLine {
		return false
	}

	return col >= c.SelectStartCol && col <= c.SelectEndCol && row >= c.SelectStartRow && row <= c.SelectEndRow
}

var HANDLED_MOTION_KEY_CODES []int = []int{
	utilities.KeySmallJ, utilities.KeySmallK, utilities.KeySmallH, utilities.KeySmallL,
	rl.KeyZero, rl.KeyOne, rl.KeyTwo, rl.KeyThree, rl.KeyFour, rl.KeyFive, rl.KeySix, rl.KeySeven, rl.KeyEight, rl.KeyNine,
	rl.KeyV,
	rl.KeyG, utilities.KeySmallG,
}

func (c *SpreadsheetCursor) AppendMotion(char rune) {
	c.MotionBuf += string(char)
	c.CheckForMotion()
}

func (c *SpreadsheetCursor) CheckForMotion() {
	motionExecuted := false
	switch c.MotionBuf {
	case "j":
		c.Row++
		motionExecuted = true
	case "k":
		c.Row--
		motionExecuted = true
	case "h":
		c.Col--
		motionExecuted = true
	case "l":
		c.Col++
		motionExecuted = true
	case "G":
		c.Row = c.MaxRow
		motionExecuted = true
	case "gg":
		c.Row = 0
		motionExecuted = true
	case "V":
		c.SetSelect(0, c.Row, c.MaxCol, c.Row)
		c.TransitionMode(ModeVLine)
		motionExecuted = true
	case "v":
		c.SetSelect(c.Col, c.Row, c.Col, c.Row)
		motionExecuted = true
	default:
		var motionRe = regexp.MustCompile(`^([0-9]+)([hjklG])$`)
		match := motionRe.FindStringSubmatch(c.MotionBuf)
		if match != nil {
			num, _ := strconv.Atoi(match[1])
			cmd := match[2]
			switch cmd {
			case "j":
				c.Row += int32(num)
				motionExecuted = true
			case "k":
				c.Row -= int32(num)
				motionExecuted = true
			case "h":
				c.Col -= int8(num)
				motionExecuted = true
			case "l":
				c.Col += int8(num)
				motionExecuted = true
			case "G":
				c.Row = int32(num - 1)
				motionExecuted = true
			}
		}
	}

	if motionExecuted {
		c.MotionBuf = ""
	}
}

func (c *SpreadsheetCursor) UpdateCmdLine() {
	c.Logs.Channel <- c.CmdBuf
}

func (c *SpreadsheetCursor) ExecuteCommand() {
	c.Logs.Channel <- "Executed command"
	c.CmdBuf = "Executed"
	c.Mode = ModeNormal
}
