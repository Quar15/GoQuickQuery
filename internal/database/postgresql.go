package database

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5"
)

type queryResult struct {
	Results *DataGrid
	Err     error
}

func queryRows(ctx context.Context, conn *pgx.Conn, query string) (ch chan queryResult) {
	ch = make(chan queryResult, 1)
	go func() {
		defer close(ch)

		rows, err := conn.Query(ctx, query)
		if err != nil {
			slog.Error(fmt.Sprintf("Query failed: %s", query), slog.Any("error", err))
			ch <- queryResult{nil, err}
			return
		}
		defer rows.Close()

		fieldDescriptions := rows.FieldDescriptions()
		dg := &DataGrid{}
		dg.Cols = 0
		dg.Headers = make([]string, len(fieldDescriptions))
		for i, field := range fieldDescriptions {
			dg.Headers[i] = string(field.Name)
			dg.Cols++
		}

		dg.Rows = 0
		for rows.Next() {
			select {
			case <-ctx.Done():
				return
			default:
			}
			values, err := rows.Values()
			if err != nil {
				ch <- queryResult{nil, err}
				return
			}

			rowMap := make(map[string]any)
			for i, col := range values {
				rowMap[dg.Headers[i]] = col
			}

			dg.Data = append(dg.Data, rowMap)
			dg.Rows++
		}

		select {
		case <-ctx.Done():
			return
		default:
		}
		ch <- queryResult{dg, nil}
	}()

	return ch
}

func QueryRows(ctx context.Context, conn *pgx.Conn, query string) (chan queryResult, error) {
	if conn.IsClosed() {
		return nil, fmt.Errorf("Broken connection")
	}
	// @TODO: Sanitize query to always limit number of results
	return queryRows(ctx, conn, query), nil
}

func connectToPostgres(connString string) (*pgx.Conn, error) {
	slog.Debug("Trying to connect with postgres", slog.String("connString", connString))
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return nil, err
	}

	return conn, nil
}

func InitPostgresConnection(connData *ConnectionData) error {
	// Create connection if does not exist
	if connData.Conn == false {
		slog.Debug("Establishing fresh connection")
		postgresConn, err := connectToPostgres(connData.ConnString)
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to initialize postgres connection for conn string '%s'", connData.ConnString), slog.Any("error", err))
			connData.ClearConn()
			return err
		}
		connData.Conn = postgresConn
	} else if conn, ok := connData.Conn.(*pgx.Conn); !ok {
		panic("Non postgres connection passed to init as postgres")
	} else if conn.IsClosed() {
		slog.Warn("Detected closed connection, reestablishing")
		postgresConn, err := connectToPostgres(connData.ConnString)
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to initialize postgres connection for conn string '%s'", connData.ConnString), slog.Any("error", err))
			connData.ClearConn()
			return err
		}
		connData.Conn = postgresConn
	} else {
		slog.Debug(fmt.Sprintf("Using already established connection for '%s'", connData.Name))
	}

	return nil
}

func ShutdownPostgresConnection(connData *ConnectionData) {
	if connData.Conn == false {
		return
	}
	if conn, ok := connData.Conn.(*pgx.Conn); ok {
		conn.Close(*connData.ConnCtx)
	}
}
