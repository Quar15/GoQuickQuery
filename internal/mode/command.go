package mode

import (
	"github.com/quar15/qq-go/internal/cursor"
	"github.com/quar15/qq-go/internal/motion"
)

type CommandMode struct{}

func (CommandMode) Handle(ctx *Context, k motion.Key) {
	if k.Code == motion.KeyEnter {
		ctx.Cursor.Common.Mode = cursor.ModeNormal
		return
	}
}
