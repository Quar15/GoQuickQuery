package mode

import (
	"log/slog"

	"github.com/quar15/qq-go/internal/cursor"
	"github.com/quar15/qq-go/internal/motion"
)

type Handler interface {
	Handle(ctx *Context, k motion.Key)
}

type Context struct {
	Cursor *cursor.Cursor
	Parser *motion.Parser
}

func HandleKey(ctx *Context, k motion.Key) {
	switch ctx.Cursor.Common.Mode {
	case cursor.ModeNormal:
		NormalMode{}.Handle(ctx, k)

	case cursor.ModeVisual:
		VisualMode{}.Handle(ctx, k)

	case cursor.ModeVLine:
		VLineMode{}.Handle(ctx, k)

	case cursor.ModeVBlock:
		VBlockMode{}.Handle(ctx, k)

	case cursor.ModeInsert:
		InsertMode{}.Handle(ctx, k)

	case cursor.ModeCommand:
		CommandMode{}.Handle(ctx, k)

	default:
		slog.Error("Handling of mode failed.", slog.String("mode", ctx.Cursor.Common.Mode.String()))
	}
}
