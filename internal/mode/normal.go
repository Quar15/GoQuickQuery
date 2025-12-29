package mode

import (
	"fmt"
	"log/slog"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/quar15/qq-go/internal/cursor"
	"github.com/quar15/qq-go/internal/motion"
)

type NormalMode struct{}

func (NormalMode) Handle(ctx *Context, k motion.Key) {
	switch k.Rune {
	case rl.KeyEscape, rl.KeyCapsLock:
		ctx.Parser.Reset()
	case 'i':
		ctx.Cursor.TransitionMode(cursor.ModeInsert)
	case 'a':
		if ctx.Cursor.Position.Col < ctx.EditorGrid.Cols[ctx.Cursor.Position.Row] {
			ctx.Cursor.Position.Col++
		}
		ctx.Cursor.TransitionMode(cursor.ModeInsert)
	case 'A':
		ctx.Cursor.Position.Col = ctx.EditorGrid.Cols[ctx.Cursor.Position.Row]
		ctx.Cursor.TransitionMode(cursor.ModeInsert)
	case 'o':
		ctx.Cursor.Position.Row, ctx.Cursor.Position.Col = ctx.EditorGrid.InsertNewLine(ctx.Cursor.Position.Row, ctx.Cursor.Position.Col)
		ctx.Cursor.TransitionMode(cursor.ModeInsert)
	case 'v':
		ctx.Cursor.Position.AnchorSelect()
		ctx.Cursor.TransitionMode(cursor.ModeVisual)
		ctx.Parser.Reset()
		return
	case 'V':
		ctx.Cursor.Position.AnchorSelect()
		if k.Modifiers == motion.ModCtrl {
			ctx.Cursor.TransitionMode(cursor.ModeVBlock)
		} else {
			ctx.Cursor.TransitionMode(cursor.ModeVLine)
		}
		ctx.Parser.Reset()
		return
	case ':':
		ctx.Cursor.TransitionMode(cursor.ModeCommand)
		ctx.Parser.Reset()
		return
	}

	if cmd, ok := ctx.Commands.Lookup(k); ok {
		slog.Debug("Normal Mode | Trying to execute command", slog.String("cmd", fmt.Sprintf("%T", cmd)))
		err := cmd.Execute(ctx)
		if err != nil {
			slog.Error("Normal Mode | Failed to execute command", slog.String("cmd", fmt.Sprintf("%T", cmd)))
		}
		return
	}

	res := ctx.Parser.Feed(k)
	if res.Done && res.Valid {
		slog.Debug("Normal Mode | Applying motion", slog.String("res.Motion", fmt.Sprintf("%T", res.Motion)), slog.Int("res.Count", res.Count))
		ctx.Cursor.Position = res.Motion.Apply(ctx.Cursor.Position, res.Count, res.HasCount)
	}
}
