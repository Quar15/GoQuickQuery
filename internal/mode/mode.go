package mode

import (
	"log/slog"

	"github.com/quar15/qq-go/internal/cursor"
	"github.com/quar15/qq-go/internal/database"
	"github.com/quar15/qq-go/internal/editor"
	"github.com/quar15/qq-go/internal/motion"
)

type Handler interface {
	Handle(ctx *Context, k motion.Key)
}

type Context struct {
	Cursor        *cursor.Cursor
	Parser        *motion.Parser
	Commands      *CommandRegistry
	ConnManager   *database.ConnectionManager
	WindowManager *WindowManager
	EditorGrid    *editor.Grid
}

func HandleKey(ctx *Context, k motion.Key) {
	switch k {
	case motion.CtrlW:
		ctx.Cursor.TransitionMode(cursor.ModeWindowManagement)
		return
	case motion.CtrlE:
		ctx.Cursor.TransitionMode(cursor.ModeNormal)
		if ctx.Cursor.Type != cursor.TypeConnections {
			ctx.WindowManager.ChangeWindow(cursor.TypeConnections)
		} else {
			ctx.WindowManager.ChangeWindow(cursor.TypeEditor)
		}
		return
	}

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

	case cursor.ModeWindowManagement:
		WindowManagementMode{}.Handle(ctx, k)

	default:
		slog.Error("Handling of mode failed.", slog.String("mode", ctx.Cursor.Common.Mode.String()))
	}
}
