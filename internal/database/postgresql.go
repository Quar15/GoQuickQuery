package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func QueryRows(ctx context.Context, conn *pgx.Conn, query string) ([]map[string]any, error) {
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Query failed: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	fieldDescriptions := rows.FieldDescriptions()
	columns := make([]string, len(fieldDescriptions))
	for i, field := range fieldDescriptions {
		columns[i] = string(field.Name)
	}

	var results []map[string]any
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, err
		}

		rowMap := make(map[string]any)
		for i, col := range values {
			rowMap[columns[i]] = col
		}

		results = append(results, rowMap)
	}

	return results, nil
}

func connectToPostgresTest() {
	dbUrl := "postgres://postgres@127.0.0.1:5432/tmp"
	conn, err := pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	data, err := QueryRows(context.Background(), conn, "SELECT * FROM example LIMIT 500;")
	if err != nil {
		fmt.Println("ERROR:", err)
	} else {
		fmt.Printf("%+v\n", data)
	}
}
