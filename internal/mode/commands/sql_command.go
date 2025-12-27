package commands

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/quar15/qq-go/internal/mode"
)

type ExecuteSQLCommand struct{}

func (ExecuteSQLCommand) Execute(ctx *mode.Context) error {
	sql, err := ctx.Cursor.DetectQuery(ctx.EditorGrid)
	if err != nil {
		slog.Warn("Failed to execute query", slog.Any("error", err))
		ctx.Cursor.Common.Logs.Channel <- fmt.Sprintf("Failed to execute query (%s)", err)
		return err
	}

	err = ctx.ConnManager.ExecuteQuery(context.Background(), ctx.ConnManager.GetCurrentConnectionName(), sql)
	if err != nil {
		slog.Error("Failed to execute query", slog.Any("error", err))
		ctx.Cursor.Common.Logs.Channel <- fmt.Sprintf("Failed to execute query (%s)", err)
	}
	return err
}
