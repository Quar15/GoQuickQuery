package cursor

import "github.com/quar15/qq-go/internal/motion"

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

type Cursor struct {
	Common   *Common
	Position motion.CursorPosition
	Type     Type
	isActive bool
}

func New(cursorType Type) *Cursor {
	return &Cursor{
		Common:   &Common{},
		Position: motion.CursorPosition{},
		Type:     cursorType,
		isActive: false,
	}
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
