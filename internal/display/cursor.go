package display

import (
	"regexp"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/quar15/qq-go/internal/assets"
	"github.com/quar15/qq-go/internal/config"
	"github.com/quar15/qq-go/internal/database"
)

type CursorMode int8

const (
	ModeNormal CursorMode = iota
	ModeInsert
	ModeVisual
	ModeVLine
	ModeVBlock
	ModeCommand
)

var modeName = map[CursorMode]string{
	ModeNormal:  "NORMAL",
	ModeInsert:  "INSERT",
	ModeVisual:  "VISUAL",
	ModeVLine:   "V-LINE",
	ModeVBlock:  "V-BLOCK",
	ModeCommand: "COMMAND",
}

func (cm CursorMode) String() string {
	return modeName[cm]
}

var modeColor = map[CursorMode]rl.Color{
	ModeNormal:  config.Get().Colors.NormalMode(),
	ModeInsert:  config.Get().Colors.InsertMode(),
	ModeVisual:  config.Get().Colors.VisualMode(),
	ModeVLine:   config.Get().Colors.VisualMode(),
	ModeVBlock:  config.Get().Colors.VisualMode(),
	ModeCommand: config.Get().Colors.NormalMode(),
}

func (cm CursorMode) Color() rl.Color {
	return modeColor[cm]
}

type CursorPosition struct {
	Col             int32
	Row             int32
	MaxCol          int32
	MaxRow          int32
	SelectStartCol  int32
	SelectStartRow  int32
	SelectEndCol    int32
	SelectEndRow    int32
	SelectAnchorCol int32
	SelectAnchorRow int32
}

func (c *CursorPosition) Reset() {
	c.SelectStartCol = -1
	c.SelectStartRow = -1
	c.SelectEndCol = -1
	c.SelectEndRow = -1
}

func (c *CursorPosition) Init() {
	c.Col = 0
	c.Row = 0
	c.Reset()
}

type CursorCommon struct {
	Mode      CursorMode
	CmdBuf    string
	MotionBuf string
	Logs      CommandLogs
}

type CursorHandler interface {
	HandleInput(appAssets *assets.Assets, dg *database.DataGrid, eg *EditorGrid, cursor *Cursor, connManager *database.ConnectionManager)
	Reset(c *Cursor)
	Init(c *Cursor, z *Zone)
}

type CursorType int8

const (
	CursorTypeEditor CursorType = iota
	CursorTypeSpreadsheet
	CursorTypeConnections
)

type Cursor struct {
	Common   *CursorCommon
	Position CursorPosition
	Handler  CursorHandler
	Type     CursorType
	Zone     *Zone
	isActive bool
}

var cursorCommon *CursorCommon = &CursorCommon{}
var CursorSpreadsheet *Cursor = &Cursor{
	Common:   cursorCommon,
	Position: CursorPosition{},
	Handler:  SpreadsheetCursorStateHandler{},
	Type:     CursorTypeSpreadsheet,
	isActive: false,
}
var CursorConnection *Cursor = &Cursor{
	Common:   cursorCommon,
	Position: CursorPosition{},
	Handler:  ConnectionsCursorStateHandler{},
	Type:     CursorTypeConnections,
	isActive: false,
}
var CursorEditor *Cursor = &Cursor{
	Common:   cursorCommon,
	Position: CursorPosition{},
	Handler:  EditorCursorStateHandler{},
	Type:     CursorTypeEditor,
	isActive: false,
}
var CurrCursor *Cursor

type BaseCursorHandler struct{}

func (BaseCursorHandler) HandleInput(c *Cursor) {}
func (BaseCursorHandler) Reset(c *Cursor) {
	c.TransitionMode(ModeNormal)
	c.Position.Reset()
	c.Common.CmdBuf = ""
	c.Common.MotionBuf = ""
	c.UpdateCmdLine()
}
func (BaseCursorHandler) Init(c *Cursor, z *Zone) {
	c.Position.Col = 0
	c.Position.Row = 0
	c.Common.Logs.Init()
	c.Handler.Reset(c)
	c.Zone = z
}

func (c *Cursor) TransitionMode(newMode CursorMode) {
	c.Common.Mode = newMode
}

func (c *Cursor) IsActive() bool {
	return c.isActive
}

func (c *Cursor) SetActive(newState bool) {
	c.isActive = newState
	c.TransitionMode(ModeNormal)
}

func (c *Cursor) SetSelect(startCol int32, startRow int32, endCol int32, endRow int32) {
	c.Position.SelectStartCol = startCol
	c.Position.SelectStartRow = startRow
	c.Position.SelectEndCol = endCol
	c.Position.SelectEndRow = endRow
	c.Position.SelectAnchorCol = startCol
	c.Position.SelectAnchorRow = startRow
	if c.Common.Mode != ModeVisual {
		c.TransitionMode(ModeVisual)
	}
}

func (c *Cursor) UpdateSelectBasedOnPosition() {
	switch c.Common.Mode {
	case ModeVBlock:
		fallthrough
	case ModeVisual:
		if c.Position.Col < c.Position.SelectAnchorCol {
			c.Position.SelectStartCol = c.Position.Col
			c.Position.SelectEndCol = c.Position.SelectAnchorCol
		} else if c.Position.Col > c.Position.SelectAnchorCol {
			c.Position.SelectStartCol = c.Position.SelectAnchorCol
			c.Position.SelectEndCol = c.Position.Col
		} else {
			c.Position.SelectStartCol = c.Position.SelectAnchorCol
			c.Position.SelectEndCol = c.Position.SelectAnchorCol
		}

		if c.Position.Row < c.Position.SelectAnchorRow {
			c.Position.SelectStartRow = c.Position.Row
			c.Position.SelectEndRow = c.Position.SelectAnchorRow
		} else if c.Position.Row > c.Position.SelectAnchorRow {
			c.Position.SelectStartRow = c.Position.SelectAnchorRow
			c.Position.SelectEndRow = c.Position.Row
		} else {
			c.Position.SelectStartRow = c.Position.SelectAnchorRow
			c.Position.SelectEndRow = c.Position.SelectAnchorRow
		}
	case ModeVLine:
		if c.Position.Row < c.Position.SelectAnchorRow {
			c.Position.SelectStartRow = c.Position.Row
			c.Position.SelectEndRow = c.Position.SelectAnchorRow
		} else if c.Position.Row > c.Position.SelectAnchorRow {
			c.Position.SelectStartRow = c.Position.SelectAnchorRow
			c.Position.SelectEndRow = c.Position.Row
		} else {
			c.Position.SelectStartRow = c.Position.SelectAnchorRow
			c.Position.SelectEndRow = c.Position.SelectAnchorRow
		}
	}
}

func (c *Cursor) IsFocused(col int32, row int32) bool {
	return c.Position.Col == col && c.Position.Row == row
}

func (c *Cursor) ClampFocus(limitCol int32, limitRow int32) {
	c.Position.MaxCol = limitCol
	c.Position.MaxRow = limitRow
	c.Position.Col = min(max(0, c.Position.Col), c.Position.MaxCol)
	c.Position.Row = min(max(0, c.Position.Row), c.Position.MaxRow)
	if c.Common.Mode == ModeVisual {
		c.Position.SelectStartCol = min(max(0, c.Position.SelectStartCol), c.Position.MaxCol)
		c.Position.SelectEndCol = min(max(0, c.Position.SelectEndCol), c.Position.MaxCol)
		c.Position.SelectStartRow = min(max(0, c.Position.SelectStartRow), c.Position.MaxRow)
		c.Position.SelectEndRow = min(max(0, c.Position.SelectEndRow), c.Position.MaxRow)
	}
}

func (c *Cursor) IsSelected(col int32, row int32) bool {
	switch c.Common.Mode {
	case ModeVisual:
		startRow, startCol := c.Position.SelectAnchorRow, c.Position.SelectAnchorCol
		endRow, endCol := c.Position.Row, c.Position.Col

		// Swap if selection is "backwards"
		if startRow > endRow || (startRow == endRow && startCol > endCol) {
			startRow, endRow = endRow, startRow
			startCol, endCol = endCol, startCol
		}

		// Single-line selection
		if startRow == endRow {
			return row == startRow && col >= startCol && col <= endCol
		}

		// Multi-line selection
		switch {
		case row > startRow && row < endRow:
			// Entire middle rows are selected
			return true
		case row == startRow:
			// Only from startCol to end of this row
			return col >= startCol
		case row == endRow:
			// Only from beginning of this row to endCol
			return col <= endCol
		default:
			return false
		}
	case ModeVLine:
		if row >= c.Position.SelectStartRow && row <= c.Position.SelectEndRow {
			return true
		}
	case ModeVBlock:
		return col >= c.Position.SelectStartCol &&
			col <= c.Position.SelectEndCol &&
			row >= c.Position.SelectStartRow &&
			row <= c.Position.SelectEndRow
	}

	return false
}

var HANDLED_MOTION_KEY_CODES []int = []int{
	keySmallJ, keySmallK, keySmallH, keySmallL,
	rl.KeyZero, rl.KeyOne, rl.KeyTwo, rl.KeyThree, rl.KeyFour, rl.KeyFive, rl.KeySix, rl.KeySeven, rl.KeyEight, rl.KeyNine,
	rl.KeyV, keySmallV,
	rl.KeyG, keySmallG,
	keySmallW,
}

func (c *Cursor) AppendMotion(char rune) {
	c.Common.MotionBuf += string(char)
	c.CheckForMotion()
}

func (c *Cursor) AppendMotionString(chars string) {
	c.Common.MotionBuf += chars
	c.CheckForMotion()
}

func (c *Cursor) CheckForMotion() {
	motionExecuted := false
	switch c.Common.MotionBuf {
	case "j":
		c.Position.Row++
		motionExecuted = true
	case "k":
		c.Position.Row--
		motionExecuted = true
	case "h":
		c.Position.Col--
		motionExecuted = true
	case "l":
		c.Position.Col++
		motionExecuted = true
	case "G":
		c.Position.Row = c.Position.MaxRow
		motionExecuted = true
	case "gg":
		c.Position.Row = 0
		motionExecuted = true
	case "V":
		c.SetSelect(0, c.Position.Row, c.Position.MaxCol, c.Position.Row)
		c.TransitionMode(ModeVLine)
		motionExecuted = true
	case "^V":
		c.SetSelect(c.Position.Col, c.Position.Row, c.Position.Col, c.Position.Row)
		c.TransitionMode(ModeVBlock)
		motionExecuted = true
	case "v":
		c.SetSelect(c.Position.Col, c.Position.Row, c.Position.Col, c.Position.Row)
		motionExecuted = true
	case "^Ww":
		// @TODO: Consider swapping from connection cursor
		// @TODO: Cleanup cursor and add transition functions
		switch c.Type {
		case CursorTypeEditor:
			CurrCursor.SetActive(false)
			CurrCursor = CursorSpreadsheet
			CurrCursor.SetActive(true)
		case CursorTypeSpreadsheet:
			CurrCursor.SetActive(false)
			CurrCursor = CursorEditor
			CurrCursor.SetActive(true)
		}
		motionExecuted = true
	default:
		var motionRe = regexp.MustCompile(`^([0-9]+)([hjklG])$`)
		match := motionRe.FindStringSubmatch(c.Common.MotionBuf)
		if match != nil {
			num, _ := strconv.Atoi(match[1])
			cmd := match[2]
			switch cmd {
			case "j":
				c.Position.Row += int32(num)
				motionExecuted = true
			case "k":
				c.Position.Row -= int32(num)
				motionExecuted = true
			case "h":
				c.Position.Col -= int32(num)
				motionExecuted = true
			case "l":
				c.Position.Col += int32(num)
				motionExecuted = true
			case "G":
				c.Position.Row = int32(num - 1)
				motionExecuted = true
			}
		}
	}

	if motionExecuted {
		c.Common.MotionBuf = ""
	}
}

func (c *Cursor) UpdateCmdLine() {
	c.Common.Logs.Channel <- c.Common.CmdBuf
}

func (c *Cursor) ExecuteCommand() {
	// @TODO
	c.Common.Logs.Channel <- "Executed command"
	c.Common.CmdBuf = "Executed"
	c.Common.Mode = ModeNormal
}

type SpreadsheetCursorStateHandler struct {
	BaseCursorHandler
}

type ConnectionsCursorStateHandler struct {
	BaseCursorHandler
}

func (ConnectionsCursorStateHandler) Init(c *Cursor, z *Zone) {
	c.Position.Col = 0
	c.Position.Row = 0
	c.Common.Logs.Init()
	c.Handler.Reset(c)
	c.Position.MaxCol = 0
	c.Position.MaxRow = 0 // @TODO: int32(connManager.GetNumberOFConnections())
	c.Zone = z
}

type EditorCursorStateHandler struct {
	BaseCursorHandler
}

// Handling input for handlers in internal/display/input.go
