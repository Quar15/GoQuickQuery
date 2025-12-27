package database

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

type ConnectionManager struct {
	connections map[string]*ConnectionData
	factory     ConnectionFactory
	current     *ConnectionData
	mu          sync.RWMutex
}

func NewConnectionManager(connConfigs []ConnectionData, factory ConnectionFactory) *ConnectionManager {
	mgr := &ConnectionManager{
		connections: make(map[string]*ConnectionData),
		factory:     factory,
	}

	for _, cfg := range connConfigs {
		c := cfg
		c.Conn = nil
		mgr.connections[c.Name] = &c
	}

	if len(connConfigs) > 0 {
		mgr.current = mgr.connections[connConfigs[0].Name]
	}

	return mgr
}

func (mgr *ConnectionManager) GetNumberOfConnections() int {
	return len(mgr.connections)
}

func (mgr *ConnectionManager) GetAllConnections() map[string]*ConnectionData {
	return mgr.connections
}

func (mgr *ConnectionManager) GetConnection(name string) (DBConnection, error) {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	connData := mgr.connections[name]

	if connData.Conn == nil || !connData.Conn.IsAlive() {
		newConn, err := mgr.factory.Create(connData.Driver, connData.ConnString)
		if err != nil {
			return nil, err
		}
		connData.Conn = newConn
	}

	return connData.Conn, nil
}

func (mgr *ConnectionManager) SetCurrentConnectionByName(name string) error {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	connData, ok := mgr.connections[name]
	if !ok {
		return fmt.Errorf("Current connection cannot be changed. Connection '%s' not found", name)
	}
	mgr.current = connData
	return nil
}

func (mgr *ConnectionManager) GetCurrentConnectionName() string {
	return mgr.current.Name
}

func (mgr *ConnectionManager) IsConnectionAlive(name string) bool {
	connData := mgr.connections[name]
	return connData.Conn != nil && connData.Conn.IsAlive()
}

func (mgr *ConnectionManager) ExecuteQuery(ctx context.Context, connectionKey string, query string) error {
	slog.Debug("Trying to execute query", slog.String("query", query), slog.String("connectionKey", connectionKey))
	// Find conn in map
	mgr.mu.Lock()
	connData, ok := mgr.connections[connectionKey]
	defer mgr.mu.Unlock()
	if !ok {
		return fmt.Errorf("No connection '%s' found", connectionKey)
	}

	if connData.Conn != nil && connData.QueryChannel != nil && connData.ConnCtxCancel != nil {
		connData.ConnCtxCancel()
		connData.ClearConn()
		return fmt.Errorf("Cancelled query")
	}

	if connData.Conn == nil || !connData.Conn.IsAlive() {
		newConn, err := mgr.factory.Create(connData.Driver, connData.ConnString)
		if err != nil {
			return err
		}

		connData.Conn = newConn
	}

	// Setup query timeout
	if connData.QueryTimeout > 0 {
		var cancelCtx context.CancelFunc
		ctx, cancelCtx = context.WithTimeout(context.Background(), time.Duration(connData.QueryTimeout)*time.Second)
		connData.ConnCtxCancel = cancelCtx
	} else {
		ctx = context.Background()
	}
	// Query data depending of type of connection
	connData.QueryText = query
	ch, err := connData.Conn.Query(ctx, query)
	if err != nil {
		connData.ClearConn()
		return err
	}
	connData.ConnCtxDone = ctx.Done()
	connData.QueryChannel = ch
	connData.QueryStartTimestamp = time.Now().UnixMilli()
	slog.Debug("Query started", slog.String("query", query))

	return nil
}

func (mgr *ConnectionManager) Close(ctx context.Context) {
	for _, connData := range mgr.connections {
		if connData.Conn != nil {
			connData.Conn.Close(ctx)
		}
	}
}
