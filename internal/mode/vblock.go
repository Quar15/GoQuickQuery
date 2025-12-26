package mode

import (
	"fmt"
	"log/slog"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/quar15/qq-go/internal/cursor"
	"github.com/quar15/qq-go/internal/motion"
)

type VBlockMode struct{}

func (VBlockMode) Handle(ctx *Context, k motion.Key) {
	if k.Rune == 'v' || k.Rune == rl.KeyEscape || k.Rune == rl.KeyCapsLock {
		ctx.Cursor.Common.Mode = cursor.ModeNormal
		ctx.Parser.Reset()
		return
	}

	res := ctx.Parser.Feed(k)
	if res.Done && res.Valid {
		slog.Debug("VBlock Mode | Applying motion", slog.String("res.Motion", fmt.Sprintf("%T", res.Motion)), slog.Int("res.Count", res.Count))
		ctx.Cursor.Position = res.Motion.Apply(ctx.Cursor.Position, res.Count, res.HasCount)
	}
	ctx.Cursor.UpdateSelectBasedOnPosition()
}
