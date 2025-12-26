package mode

import (
	"github.com/quar15/qq-go/internal/cursor"
	"github.com/quar15/qq-go/internal/motion"
)

type InsertMode struct{}

func (InsertMode) Handle(ctx *Context, k motion.Key) {
	if k.Code == motion.KeyEsc {
		ctx.Cursor.Common.Mode = cursor.ModeNormal
		return
	}
	// @TODO: text insertion handled elsewhere
}
