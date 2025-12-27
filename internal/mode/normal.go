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
	case 'v':
		ctx.Cursor.Position.AnchorSelect()
		ctx.Cursor.Common.Mode = cursor.ModeVisual
		ctx.Parser.Reset()
		return
	case 'V':
		ctx.Cursor.Position.AnchorSelect()
		if k.Modifiers == motion.ModCtrl {
			ctx.Cursor.Common.Mode = cursor.ModeVBlock
		} else {
			ctx.Cursor.Common.Mode = cursor.ModeVLine
		}
		ctx.Parser.Reset()
		return
	case ':':
		ctx.Cursor.Common.Mode = cursor.ModeCommand
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
