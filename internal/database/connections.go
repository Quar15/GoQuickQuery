package database

import (
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

var DBConnections map[string]*ConnectionData = make(map[string]*ConnectionData)

type ConnectionData struct {
	Name       string `yaml:"name"`
	Driver     string `yaml:"driver"`
	ConnString string `yaml:"conn"`
	Conn       any    `yaml:"-"`
}

func InitializeConnections(connections []ConnectionData) {
	for _, c := range connections {
		c.Conn = false
		DBConnections[c.Name] = &c
	}
}

func QueryData(connectionKey string, query string) (*DataGrid, error) {
	dg := &DataGrid{}
	var err error
	// Find conn in map
	connData, ok := DBConnections[connectionKey]
	if !ok {
		return dg, fmt.Errorf("No connection '%s' found", connectionKey)
	}
	// Create connection if does not exist
	if connData.Conn == false {
		postgresConn, err := ConnectToPostgres(connData.ConnString)
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to initialize postgres connection for conn string '%s'", connData.ConnString), slog.Any("error", err))
			return dg, err
		}
		DBConnections[connectionKey].Conn = postgresConn
	} else {
		slog.Debug(fmt.Sprintf("Using already established connection for '%s'", connectionKey))
	}

	// Query data depending of type of connection
	switch c := connData.Conn.(type) {
	case *pgx.Conn:
		dg, err = QueryRows(c, query)
		if err != nil {
			return &DataGrid{}, nil
		}
	}

	return dg, nil
}
