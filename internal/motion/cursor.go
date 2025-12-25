package motion

type CursorPosition struct {
	Col    int32
	Row    int32
	MaxCol int32
	MaxRow int32
	SelectStartCol  int32
	SelectStartRow  int32
	SelectEndCol    int32
	SelectEndRow    int32
	SelectAnchorCol int32
	SelectAnchorRow int32
}

func (c CursorPosition) Clamp() CursorPosition {
	c.Col = min(max(0, c.Col), c.MaxCol)
	c.Row = min(max(0, c.Row), c.MaxRow)

	return c
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
