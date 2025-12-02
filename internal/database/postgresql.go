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
		fmt.Fprintf(os.Stderr, "Query failed: %v\n", err)
		os.Exit(1)
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

func QueryRows(conn *pgx.Conn) (dg *DataGrid, err error) {
	data, headers, cols, rows, err := queryRows(conn, "SELECT * FROM example LIMIT 500;")
	if err != nil {
		slog.Error("Failed to execute query", slog.Any("err", err))
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

func ConnectToPostgres() (*pgx.Conn, error) {
	dbUrl := "postgres://postgres@127.0.0.1:5432/tmp"
	conn, err := pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return nil, err
	}

	return conn, nil
}
