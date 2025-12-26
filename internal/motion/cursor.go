package motion

type CursorPosition struct {
	Col             int32
	Row             int32
	MaxCol          int32
	MaxRow          int32
	MaxColForRows   []int32
	SelectStartCol  int32
	SelectStartRow  int32
	SelectEndCol    int32
	SelectEndRow    int32
	SelectAnchorCol int32
	SelectAnchorRow int32
}

func (c CursorPosition) Clamp() CursorPosition {
	c.Row = min(max(0, c.Row), c.MaxRow)
	if len(c.MaxColForRows) > 0 {
		c.Col = min(max(0, c.Col), c.MaxColForRows[c.Row]-1)
	} else {
		c.Col = min(max(0, c.Col), c.MaxCol)
	}

	return c
}

func (c *CursorPosition) AnchorSelect() {
	c.ResetSelect()
	c.SelectAnchorCol = c.Col
	c.SelectAnchorRow = c.Row
}

func (c *CursorPosition) ResetSelect() {
	c.SelectStartCol = c.Col
	c.SelectStartRow = c.Row
	c.SelectEndCol = c.Col
	c.SelectEndRow = c.Row
}

func (c *CursorPosition) Init() {
	c.Col = 0
	c.Row = 0
	c.ResetSelect()
}
