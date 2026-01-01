package motion

import "log/slog"

type Motion interface {
	Apply(pos CursorPosition, count int, hasCount bool) CursorPosition
}

type DebugMotion struct{}

func (DebugMotion) Apply(pos CursorPosition, count int, hasCount bool) CursorPosition {
	slog.Debug("DebugMotion", slog.Any("pos", pos), slog.Int("count", count), slog.Bool("hasCount", hasCount))
	return pos
}
