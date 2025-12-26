package motion

type MoveLeft struct{}
type MoveRight struct{}
type MoveUp struct{}
type MoveDown struct{}

type MoveStartLeft struct{}
type MoveEndRight struct{}
type MoveStartUp struct{}
type MoveEndDown struct{}

type MoveToSpecificLineOrDown struct{}

func (MoveLeft) Apply(p CursorPosition, count int, hasCount bool) CursorPosition {
	p.Col -= int32(count)
	return p.Clamp()
}

func (MoveRight) Apply(p CursorPosition, count int, hasCount bool) CursorPosition {
	p.Col += int32(count)
	return p.Clamp()
}

func (MoveUp) Apply(p CursorPosition, count int, hasCount bool) CursorPosition {
	p.Row -= int32(count)
	return p.Clamp()
}

func (MoveDown) Apply(p CursorPosition, count int, hasCount bool) CursorPosition {
	p.Row += int32(count)
	return p.Clamp()
}

func (MoveStartLeft) Apply(p CursorPosition, count int, hasCount bool) CursorPosition {
	p.Col = 0
	return p.Clamp()
}

func (MoveEndRight) Apply(p CursorPosition, count int, hasCount bool) CursorPosition {
	p.Col = p.MaxColForRows[p.Row]
	return p.Clamp()
}

func (MoveStartUp) Apply(p CursorPosition, count int, hasCount bool) CursorPosition {
	p.Row = 0
	return p.Clamp()
}

func (MoveEndDown) Apply(p CursorPosition, count int, hasCount bool) CursorPosition {
	p.Row = p.MaxRow
	return p.Clamp()
}

func (MoveToSpecificLineOrDown) Apply(p CursorPosition, count int, hasCount bool) CursorPosition {
	if !hasCount {
		return MoveEndDown{}.Apply(p, count, hasCount)
	}
	p.Row = int32(count - 1)
	return p.Clamp()
}
