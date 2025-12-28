package cursor

import (
	"errors"
	"strings"

	"github.com/quar15/qq-go/internal/editor"
	"github.com/quar15/qq-go/internal/motion"
)

type Common struct {
	Mode      Mode
	CmdBuf    string
	MotionBuf string
	Logs      CommandLogs
}

type Type int8

const (
	TypeEditor Type = iota
	TypeSpreadsheet
	TypeConnections
)

var typeName = map[Type]string{
	TypeEditor:      "TypeEditor",
	TypeSpreadsheet: "TypeSpreadsheet",
	TypeConnections: "TypeConnections",
}

func (t Type) String() string {
	return typeName[t]
}

type Cursor struct {
	Common   *Common
	Position motion.CursorPosition
	Type     Type
	isActive bool
}

func New(common *Common, cursorType Type) *Cursor {
	c := &Cursor{
		Common:   common,
		Position: motion.CursorPosition{},
		Type:     cursorType,
		isActive: false,
	}
	c.Common.Logs.Init()

	return c
}

func (c *Cursor) Reset() {
	c.TransitionMode(ModeNormal)
	c.Position.ResetSelect()
	c.Common.CmdBuf = ""
	c.Common.MotionBuf = ""
	c.UpdateCmdLine()
}

func (c *Cursor) TransitionMode(newMode Mode) {
	c.Common.Mode = newMode
}

func (c *Cursor) UpdateCmdLine() {
	c.Common.Logs.Log(c.Common.CmdBuf)
}

func (c *Cursor) IsActive() bool {
	return c.isActive
}

func (c *Cursor) Activate() {
	c.isActive = true
}

func (c *Cursor) Deactivate() {
	c.isActive = false
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

func (c *Cursor) DetectQuery(eg *editor.Grid) (string, error) {
	query := ""
	if eg.Rows <= 0 {
		errMsg := "No query provided"
		return "", errors.New(errMsg)
	}
	// @TODO: Implement other modes behavior
	switch c.Common.Mode {
	case ModeVisual:
		for row := c.Position.SelectStartRow; row <= c.Position.SelectEndRow; row++ {
			if eg.Cols[row] > 0 {
				for col := int32(0); col < eg.Cols[row]; col++ {
					if c.IsSelected(col, row) {
						query += string(eg.Text[row][col])
					}
				}
			}
		}
	case ModeVLine:
		for row := c.Position.SelectStartRow; row <= c.Position.SelectEndRow; row++ {
			if eg.Cols[row] > 0 {
				for col := int32(0); col < eg.Cols[row]; col++ {
					query += string(eg.Text[row][col])
				}
			}
		}
	case ModeInsert:
		fallthrough
	case ModeNormal:
		start, end := eg.DetectQueryRowsBoundaryBasedOnRow(c.Position.Row)
		var sb strings.Builder
		sb.Grow(128)

		for i := start; i <= end; i++ {
			sb.WriteString(eg.Text[i])
			if i < end {
				sb.WriteString(" ")
			}
		}

		query = strings.TrimSpace(sb.String())
	}
	if query == "" {
		errMsg := "No query provided/found"
		return "", errors.New(errMsg)
	}

	return query, nil
}
