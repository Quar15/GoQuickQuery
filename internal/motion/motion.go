package motion

type Motion interface {
	Apply(pos CursorPosition, count int) CursorPosition
}
