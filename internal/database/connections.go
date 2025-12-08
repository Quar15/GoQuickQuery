package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
)

var DBConnections map[string]*ConnectionData = make(map[string]*ConnectionData)

type ConnectionData struct {
	Name                string             `yaml:"name"`
	Driver              string             `yaml:"driver"`
	ConnString          string             `yaml:"conn"`
	QueryTimeout        int                `yaml:"timeout"`
	Conn                any                `yaml:"-"`
	ConnCtx             *context.Context   `yaml:"-"`
	ConnCtxCancel       context.CancelFunc `yaml:"-"`
	QueryChannel        chan queryResult   `yaml:"-"`
	QueryText           string             `yaml:"-"`
	QueryStartTimestamp int64              `yaml:"-"`
}

func InitializeConnections(connections []ConnectionData) {
	for _, c := range connections {
		c.Conn = false
		DBConnections[c.Name] = &c
	}
}

func QueryData(connectionKey string, query string) error {
	// Find conn in map
	connData, ok := DBConnections[connectionKey]
	if !ok {
		return fmt.Errorf("No connection '%s' found", connectionKey)
	}

	if connData.QueryChannel != nil && connData.ConnCtxCancel != nil {
		connData.ConnCtxCancel()
		return nil
	}

	switch connData.Driver {
	case "postgresql":
		if err := InitPostgresConnection(connData); err != nil {
			return err
		}
	default:
		return fmt.Errorf("Driver of type '%s' is not handled", connData.Driver)
	}
	// Setup query timeout
	var ctx context.Context
	if connData.QueryTimeout > 0 {
		var cancelCtx context.CancelFunc
		ctx, cancelCtx = context.WithTimeout(context.Background(), time.Duration(connData.QueryTimeout)*time.Second)
		connData.ConnCtxCancel = cancelCtx
	} else {
		ctx = context.Background()
	}
	// Query data depending of type of connection
	connData.QueryText = query
	var err error
	switch c := connData.Conn.(type) {
	case *pgx.Conn:
		connData.ConnCtx = &ctx
		connData.QueryChannel = nil
		connData.QueryChannel, err = QueryRows(*connData.ConnCtx, c, query)
		if err != nil {
			connData.ClearConn()
			return err
		}
		connData.QueryStartTimestamp = time.Time.UnixMilli(time.Now())
	}

	return nil
}

func CheckForResult(ctx context.Context, ch chan queryResult, connName string) (dg *DataGrid, done bool, err error) {
	dg = &DataGrid{}
	select {
	// Check if query timed out
	case <-ctx.Done():
		slog.Warn("Timeout:", slog.Any("error", ctx.Err()))
		DBConnections[connName].ClearConn()
		return nil, true, fmt.Errorf("Timeout reached | Cancelled query")
	// Check if query finished
	case res, ok := <-ch:
		if !ok {
			// Thread closed - race condition possible on query finish / cancel
			return nil, false, nil
		}
		if res.Err != nil {
			slog.Error("Failed to execute query", slog.Any("error", res.Err))
			DBConnections[connName].ClearConn()
			return nil, true, err
		}
		slog.Debug("Query result", slog.Any("res", res))
		*dg = *res.Results
		DBConnections[connName].ClearQuery()
		return dg, true, nil
	// Query still running
	default:
		return nil, false, nil
	}
}

func (c *ConnectionData) ClearConn() {
	c.ClearQuery()
	c.Conn = false
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
