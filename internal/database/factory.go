package database

import (
	"fmt"
	"log/slog"
)

type ConnectionFactory interface {
	Create(driver string, connString string) (DBConnection, error)
}

type DefaultConnectionFactory struct{}

func (f *DefaultConnectionFactory) Create(driver string, connString string) (DBConnection, error) {
	switch driver {
	case "postgresql":
		conn, err := connectToPostgres(connString)
		if err != nil {
			return nil, err
		}
		slog.Debug("Created new postgres connection", slog.String("connString", connString))
		return &PostgresConn{Conn: conn}, nil
	default:
		return nil, fmt.Errorf("Unsupported driver: %s", driver)
	}
}
