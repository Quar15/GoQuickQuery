package commands

import (
	"github.com/quar15/qq-go/internal/config"
	"github.com/quar15/qq-go/internal/cursor"
	"github.com/quar15/qq-go/internal/mode"
)

type ConnectionsChange struct{}

func (ConnectionsChange) Execute(ctx *mode.Context) error {
	err := ctx.ConnManager.SetCurrentConnectionByName(config.Get().Connections[ctx.Cursor.Position.Row].Name)
	if err != nil {
		return err
	}
	ctx.WindowManager.ChangeWindow(cursor.TypeEditor)
	return nil
}
