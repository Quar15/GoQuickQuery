package motion

type MoveLeft struct{}
type MoveRight struct{}
type MoveUp struct{}
type MoveDown struct{}

func (MoveLeft) Apply(p CursorPosition, count int) CursorPosition {
	p.Col -= int32(count)
	return p.Clamp()
}

func (MoveRight) Apply(p CursorPosition, count int) CursorPosition {
	p.Col += int32(count)
	return p.Clamp()
}

func (MoveUp) Apply(p CursorPosition, count int) CursorPosition {
	p.Row -= int32(count)
	return p.Clamp()
}

func (MoveDown) Apply(p CursorPosition, count int) CursorPosition {
	p.Row += int32(count)
	return p.Clamp()
}
