package mode

import (
	"log/slog"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/quar15/qq-go/internal/cursor"
	"github.com/quar15/qq-go/internal/motion"
)

type WindowManagementMode struct{}

func (WindowManagementMode) Handle(ctx *Context, k motion.Key) {
	slog.Debug("Window Management Mode | Handling key", slog.Any("key", k))
	switch k.Rune {
	case rl.KeyEscape, rl.KeyCapsLock:
		ctx.Parser.Reset()
		ctx.Cursor.TransitionMode(cursor.ModeNormal)
	case rl.KeyW:
	case rl.KeyN:
	case rl.KeyP:
	case rl.KeyE:
		if k.Modifiers == motion.ModCtrl {
		}
	case '=':
	case '+':
	case rl.KeyMinus:
	}
}
