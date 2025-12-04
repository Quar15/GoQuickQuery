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
	Name          string              `yaml:"name"`
	Driver        string              `yaml:"driver"`
	ConnString    string              `yaml:"conn"`
	QueryTimeout  int                 `yaml:"timeout"`
	Conn          any                 `yaml:"-"`
	ConnCtx       *context.Context    `yaml:"-"`
	ConnCtxCancel *context.CancelFunc `yaml:"-"`
	QueryChannel  chan queryResult    `yaml:"-"`
	QueryText     string              `yaml:"-"`
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
	// Create connection if does not exist
	if connData.Conn == false {
		postgresConn, err := ConnectToPostgres(connData.ConnString)
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to initialize postgres connection for conn string '%s'", connData.ConnString), slog.Any("error", err))
			return err
		}
		DBConnections[connectionKey].Conn = postgresConn
	} else {
		slog.Debug(fmt.Sprintf("Using already established connection for '%s'", connectionKey))
	}
	// Setup query timeout
	var ctx context.Context
	if connData.QueryTimeout > 0 {
		var cancelCtx context.CancelFunc
		ctx, cancelCtx = context.WithTimeout(context.Background(), time.Duration(connData.QueryTimeout)*time.Second)
		connData.ConnCtxCancel = &cancelCtx
	} else {
		ctx = context.Background()
	}
	// Query data depending of type of connection
	connData.QueryText = query
	switch c := connData.Conn.(type) {
	case *pgx.Conn:
		connData.ConnCtx = &ctx
		connData.QueryChannel = QueryRows(*connData.ConnCtx, c, query)
	}

	return nil
}

func CheckForResult(ctx context.Context, ch chan queryResult, connName string) (dg *DataGrid, done bool, err error) {
	dg = &DataGrid{}
	select {
	// Check if query timed out
	case <-ctx.Done():
		slog.Warn("Timeout:", slog.Any("error", ctx.Err()))
		DBConnections[connName].Conn = false // pgx closes the connection on context timeout
		return nil, true, fmt.Errorf("Timeout reached | Cancelled query")
	// Check if query finished
	case res := <-ch:
		if res.Err != nil {
			slog.Error("Failed to execute query", slog.Any("error", err))
			return nil, true, err
		}
		*dg = *res.Results
		DBConnections[connName].ConnCtxCancel = nil
		DBConnections[connName].ConnCtx = nil
		return dg, true, nil
	// Query still running
	default:
		return nil, false, nil
	}
}
