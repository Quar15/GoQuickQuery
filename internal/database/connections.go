package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

type DBConnection interface {
	// Query executes query in new thread and will return channel that will return result when ready
	Query(ctx context.Context, query string) (chan queryResult, error)

	// Close terminates the connection gracefully
	Close(ctx context.Context) error

	// IsAlive checks if the connection is still valid
	IsAlive() bool
}

type ConnectionData struct {
	Name                string             `yaml:"name"`
	Driver              string             `yaml:"driver"`
	ConnString          string             `yaml:"conn"`
	QueryTimeout        int                `yaml:"timeout"`
	Conn                DBConnection       `yaml:"-"`
	ConnCtx             *context.Context   `yaml:"-"`
	ConnCtxCancel       context.CancelFunc `yaml:"-"`
	QueryChannel        chan queryResult   `yaml:"-"`
	QueryText           string             `yaml:"-"`
	QueryStartTimestamp int64              `yaml:"-"`
}

func (c *ConnectionData) CheckForQueryResult() (dg *DataGrid, done bool, err error) {
	ctx := *c.ConnCtx
	ch := c.QueryChannel
	dg = &DataGrid{}
	select {
	// Check if query timed out
	case <-ctx.Done():
		slog.Warn("Timeout:", slog.Any("error", ctx.Err()))
		c.ClearConn()
		return nil, true, fmt.Errorf("Timeout reached | Cancelled query")
	// Check if query finished
	case res, ok := <-ch:
		if !ok {
			// Thread closed - race condition possible on query finish / cancel
			return nil, false, nil
		}
		if res.Err != nil {
			slog.Error("Failed to execute query", slog.Any("error", res.Err))
			c.ClearConn()
			return nil, true, err
		}
		slog.Debug("Query result", slog.Any("res", res))
		*dg = *res.Results
		c.ClearQuery()
		return dg, true, nil
	// Query still running
	default:
		return nil, false, nil
	}
}

func (c *ConnectionData) ClearConn() {
	c.ClearQuery()
	c.Conn = nil
}

func (c *ConnectionData) ClearQuery() {
	c.QueryChannel = nil
	c.ConnCtx = nil
	c.ConnCtxCancel = nil
}

func (c *ConnectionData) GetQueryRuntime() int64 {
	return time.Time.UnixMilli(time.Now()) - c.QueryStartTimestamp
}

func (c *ConnectionData) GetQueryRuntimeDynamicString() string {
	queryRuntime := c.GetQueryRuntime()
	if queryRuntime < 1000 {
		return fmt.Sprintf("%dms", queryRuntime)
	}
	if queryRuntime < 60000 {
		return fmt.Sprintf("%ds", queryRuntime/1000)
	}
	queryRuntimeMinutes := 0
	for queryRuntime > 1000 {
		queryRuntime /= 1000
		queryRuntimeMinutes++
	}
	return fmt.Sprintf("%dmin %ds", queryRuntimeMinutes, queryRuntime)
}
