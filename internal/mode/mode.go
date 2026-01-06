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
	DataGrid      *database.DataGrid
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

	case cursor.ModeWindowManagement:
		WindowManagementMode{}.Handle(ctx, k)

	default:
		slog.Error("Handling of mode failed.", slog.String("mode", ctx.Cursor.Common.Mode.String()))
	}
}

func (ctx *Context) UpdateCursorPositionMax() {
	ctx.Cursor.Position.UpdateMax(
		max(0, ctx.EditorGrid.Cols[ctx.Cursor.Position.Row]-1),
		max(0, ctx.EditorGrid.Rows-1),
		ctx.EditorGrid.Cols,
	)
	//slog.Debug("Cursor max positions updated", slog.Any("ctx.Cursor.Position", ctx.Cursor.Position))
}
