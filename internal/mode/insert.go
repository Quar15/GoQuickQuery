package mode

import (
	"log/slog"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/quar15/qq-go/internal/cursor"
	"github.com/quar15/qq-go/internal/motion"
)

type InsertMode struct{}

func (InsertMode) Handle(ctx *Context, k motion.Key) {
	if k.Code == motion.KeyEsc {
		ctx.Cursor.Common.Mode = cursor.ModeNormal
		return
	}
	row := ctx.Cursor.Position.Row
	col := ctx.Cursor.Position.Col

	switch k.Rune {
	case rl.KeyEnter:
		ctx.Cursor.Position.Row, ctx.Cursor.Position.Col = ctx.EditorGrid.InsertNewLine(row, col)
	case rl.KeyDelete:
	case rl.KeyBackspace:
		ctx.Cursor.Position.Row, ctx.Cursor.Position.Col = ctx.EditorGrid.DeleteCharBefore(row, col)
	case rl.KeyLeft, rl.KeyDown, rl.KeyUp, rl.KeyRight:
	default:
		slog.Debug("Inserting character", slog.Any("Rune", k.Rune))
		ctx.Cursor.Position.Col = ctx.EditorGrid.InsertChar(row, col, k.Rune)
	}

	ctx.Cursor.Position.MaxRow = ctx.EditorGrid.Rows - 1
}
