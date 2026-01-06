package mode

import (
	"fmt"
	"log/slog"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/quar15/qq-go/internal/cursor"
	"github.com/quar15/qq-go/internal/motion"
)

type VisualMode struct{}

func (VisualMode) Handle(ctx *Context, k motion.Key) {
	if k.Rune == 'v' || k.Rune == rl.KeyEscape || k.Rune == rl.KeyCapsLock {
		ctx.Cursor.Common.Mode = cursor.ModeNormal
		ctx.Parser.Reset()
		return
	}

	if cmd, ok := ctx.Commands.Lookup(k); ok {
		slog.Debug("Visual Mode | Trying to execute command", slog.String("cmd", fmt.Sprintf("%T", cmd)))
		err := cmd.Execute(ctx)
		if err != nil {
			slog.Error("Visual Mode | Failed to execute command", slog.String("cmd", fmt.Sprintf("%T", cmd)))
		}
		return
	}

	res := ctx.Parser.Feed(k)
	if res.Done && res.Valid {
		slog.Debug("Visual Mode | Applying motion", slog.String("res.Motion", fmt.Sprintf("%T", res.Motion)), slog.Int("res.Count", res.Count))
		ctx.Cursor.Position = res.Motion.Apply(ctx.Cursor.Position, res.Count, res.HasCount)
	}
	ctx.Cursor.UpdateSelectBasedOnPosition()
}
