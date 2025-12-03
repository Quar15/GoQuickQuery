package database

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5"
)

func queryRows(conn *pgx.Conn, query string) (data []map[string]any, headers []string, colsN int8, rowsN int32, err error) {
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		slog.Error(fmt.Sprintf("Query failed: %s", query), slog.Any("error", err))
		return nil, nil, 0, 0, err
	}
	defer rows.Close()

	fieldDescriptions := rows.FieldDescriptions()
	colsN = 0
	columns := make([]string, len(fieldDescriptions))
	for i, field := range fieldDescriptions {
		columns[i] = string(field.Name)
		colsN++
	}

	var results []map[string]any
	rowsN = 0
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, nil, -1, -1, err
		}

		rowMap := make(map[string]any)
		for i, col := range values {
			rowMap[columns[i]] = col
		}

		results = append(results, rowMap)
		rowsN++
	}

	return results, columns, colsN, rowsN, nil
}

func QueryRows(conn *pgx.Conn, query string) (dg *DataGrid, err error) {
	// @TODO: Sanitize query to always limit number of results
	data, headers, cols, rows, err := queryRows(conn, query)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to execute query '%s'", query), slog.Any("error", err))
		return nil, err
	} else {
		slog.Info("Query finished")
		fmt.Printf("%+v\n", data)
	}
	dg = &DataGrid{}
	dg.Data = data
	dg.Headers = headers
	dg.Cols = cols
	dg.Rows = rows
	return dg, nil
}

func ConnectToPostgres(connString string) (*pgx.Conn, error) {
	slog.Debug("Trying to connect with postgres", slog.String("connString", connString))
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return nil, err
	}

	return conn, nil
}
