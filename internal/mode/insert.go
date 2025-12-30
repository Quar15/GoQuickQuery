package mode

import (
	"log/slog"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/quar15/qq-go/internal/cursor"
	"github.com/quar15/qq-go/internal/motion"
)

type InsertMode struct{}

func (InsertMode) Handle(ctx *Context, k motion.Key) {
	row := ctx.Cursor.Position.Row
	col := ctx.Cursor.Position.Col

	if k.Code == motion.KeyEsc {
		ctx.UpdateCursorPositionMax()
		ctx.Cursor.Position.Col = max(col-1, 0)
		ctx.Cursor.TransitionMode(cursor.ModeNormal)

		return
	}

	switch k.Rune {
	case rl.KeyEnter:
		ctx.Cursor.Position.Row, ctx.Cursor.Position.Col = ctx.EditorGrid.InsertNewLine(row, col)
	case rl.KeyDelete:
	case rl.KeyBackspace:
		ctx.Cursor.Position.Row, ctx.Cursor.Position.Col = ctx.EditorGrid.DeleteCharBefore(row, col)
	case rl.KeyLeft, rl.KeyDown, rl.KeyUp, rl.KeyRight:
	default:
		if k.Rune > 31 && k.Rune < 127 {
			slog.Debug("Inserting character", slog.Any("Rune", k.Rune), slog.String("string(Rune)", string(k.Rune)))
			ctx.Cursor.Position.Col = ctx.EditorGrid.InsertChar(row, col, k.Rune)
		} else {
			slog.Debug("SKIP Inserting character", slog.Any("Rune", k.Rune), slog.String("string(Rune)", string(k.Rune)))
		}
	}

	ctx.UpdateCursorPositionMax()
}
