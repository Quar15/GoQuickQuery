package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/lmittmann/tint"
	"golang.design/x/clipboard"

	"github.com/quar15/qq-go/internal/assets"
	"github.com/quar15/qq-go/internal/config"
	"github.com/quar15/qq-go/internal/cursor"
	"github.com/quar15/qq-go/internal/database"
	"github.com/quar15/qq-go/internal/display"
	"github.com/quar15/qq-go/internal/editor"
	"github.com/quar15/qq-go/internal/format"
	"github.com/quar15/qq-go/internal/mode"
	"github.com/quar15/qq-go/internal/motion"
	"github.com/quar15/qq-go/internal/setup"
)

// App holds all application state and dependencies.
type App struct {
	cfg      *config.Config
	assets   *assets.Assets
	connMgr  *database.ConnectionManager
	dataGrid *database.DataGrid
	editGrid *editor.Grid
	splitter *display.Splitter
	zones    *zones
	cursors  *cursors
}

type zones struct {
	top         display.Zone
	bottom      display.Zone
	command     display.Zone
	connections display.Zone
}

type cursors struct {
	common      *cursor.Common
	editor      *mode.Context
	spreadsheet *mode.Context
	connections *mode.Context
}

func main() {
	if err := run(); err != nil {
		slog.Error("Application error", slog.Any("error", err))
		os.Exit(1)
	}
}

func run() error {
	initLogger()

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("Load config: %w", err)
	}

	if err := clipboard.Init(); err != nil {
		return fmt.Errorf("Init clipboard: %w", err)
	}

	// Initialize color mappings
	editor.InitHighlightColors(cfg)
	cursor.InitModeColors(cfg)

	app := newApp(cfg)
	defer app.close()

	return app.run()
}

func initLogger() {
	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stdout, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.RFC3339,
		}),
	))
}

func newApp(cfg *config.Config) *App {
	connMgr := database.NewConnectionManager(cfg.Connections, &database.DefaultConnectionFactory{})

	dg := &database.DataGrid{}
	eg := editor.NewGrid()

	cursorCommon := &cursor.Common{}
	cursorCommon.Logs.Init()

	app := &App{
		cfg:      cfg,
		connMgr:  connMgr,
		dataGrid: dg,
		editGrid: eg,
		splitter: &display.Splitter{
			Ratio:    0.6,
			Height:   6.0,
			Dragging: false,
		},
		zones: &zones{},
		cursors: &cursors{
			common:      cursorCommon,
			editor:      initEditorContext(cursorCommon, connMgr, eg),
			spreadsheet: initSpreadsheetContext(cursorCommon),
			connections: initConnectionsContext(cursorCommon),
		},
	}

	app.cursors.editor.Cursor.Activate()

	return app
}

func (a *App) close() {
	a.connMgr.Close(context.Background())
	if a.assets != nil {
		a.assets.UnloadAssets()
	}
	rl.CloseWindow()
}

func (a *App) run() error {
	a.initWindow()

	appAssets := &assets.Assets{}
	if err := appAssets.LoadAssets(); err != nil {
		return fmt.Errorf("load assets: %w", err)
	}
	a.assets = appAssets

	commandZoneHeight := appAssets.MainFontSize*2 + appAssets.MainFontSpacing*2
	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		a.update(commandZoneHeight)
		a.draw()
	}

	return nil
}

func (a *App) initWindow() {
	screenWidth := rl.GetScreenWidth()
	screenHeight := rl.GetScreenHeight()

	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(int32(screenWidth), int32(screenHeight), "GQQ")
	rl.SetExitKey(rl.KeyNull)
}

func (a *App) update(commandZoneHeight float32) {
	screenWidth := max(rl.GetScreenWidth(), 100)
	screenHeight := max(rl.GetScreenHeight(), 100)

	rl.SetMouseCursor(rl.MouseCursorDefault)
	a.splitter.HandleZoneSplit(screenWidth, screenHeight, int(commandZoneHeight))

	a.updateZoneBounds(screenWidth, screenHeight, commandZoneHeight)

	a.zones.top.UpdateZoneScroll()
	a.zones.bottom.UpdateZoneScroll()

	a.handleDroppedFiles()
	display.HandleInput(a.cursors.editor)
	a.handleQueryResults()
}

func (a *App) updateZoneBounds(screenWidth, screenHeight int, commandZoneHeight float32) {
	a.zones.top.Bounds = rl.Rectangle{
		X:      0,
		Y:      0,
		Width:  float32(screenWidth),
		Height: a.splitter.Y - a.splitter.Height/2,
	}

	a.zones.bottom.Bounds = rl.Rectangle{
		X:      0,
		Y:      a.splitter.Y + a.splitter.Height/2,
		Width:  float32(screenWidth),
		Height: float32(screenHeight) - (a.splitter.Y + a.splitter.Height/2) - commandZoneHeight,
	}

	a.zones.command.Bounds = rl.Rectangle{
		X:      0,
		Y:      a.zones.bottom.Bounds.Y + a.zones.bottom.Bounds.Height,
		Width:  float32(screenWidth),
		Height: commandZoneHeight,
	}
}

func (a *App) draw() {
	rl.BeginDrawing()
	defer rl.EndDrawing()

	rl.ClearBackground(a.cfg.Colors.Background())

	editorIsFocused := a.cursors.editor.Cursor.IsActive()
	a.zones.top.DrawEditor(a.assets, a.editGrid, a.cursors.editor.Cursor, editorIsFocused)
	a.zones.bottom.DrawSpreadsheetZone(a.assets, a.dataGrid, a.cursors.spreadsheet.Cursor)
	a.zones.command.DrawCommandZone(a.cfg, a.assets, a.cursors.editor.Cursor, a.connMgr.GetCurrentConnectionName())

	a.splitter.Draw(a.cursors.editor.Cursor.Type)

	if a.cursors.connections.Cursor.IsActive() {
		screenWidth := int32(rl.GetScreenWidth())
		screenHeight := int32(rl.GetScreenHeight())
		a.zones.connections.DrawConnectionSelector(a.assets, a.cfg, a.cursors.connections.Cursor, screenWidth, screenHeight, a.connMgr)
	}
}

func (a *App) handleDroppedFiles() {
	if !rl.IsFileDropped() {
		return
	}

	paths := rl.LoadDroppedFiles()
	if len(paths) == 0 {
		return
	}

	if len(paths) > 1 {
		slog.Warn("Tried to load more than one file")
	}

	path := paths[0]
	ext := filepath.Ext(path)

	switch ext {
	case ".sql":
		a.loadSQLFile(path)
	case ".csv":
		a.loadCSVFile(path)
	default:
		slog.Warn("Unhandled type of file", slog.String("extension", ext))
	}
}

func (a *App) loadSQLFile(path string) {
	newEg, err := editor.LoadGridFromTextFile(path, a.assets)
	if err != nil {
		slog.Error("Failed to load file", slog.String("path", path))
		a.cursors.editor.Cursor.Common.Logs.Log(fmt.Sprintf("ERR: Failed to parse sql file '%s'", path))
		return
	}

	*a.editGrid = *newEg
	cur := a.cursors.editor.Cursor
	cur.Reset()
	cur.Position.MaxCol = a.editGrid.MaxCol
	cur.Position.MaxColForRows = a.editGrid.Cols
	cur.Position.MaxRow = a.editGrid.Rows - 1
	cur.Common.Logs.Log(fmt.Sprintf("Loaded sql file '%s'", path))
	slog.Info("Loaded sql file", slog.String("path", path))
}

func (a *App) loadCSVFile(path string) {
	newDg, err := database.LoadDataGridFromCSV(path, a.assets)
	if err != nil {
		slog.Error("Failed to load file", slog.String("path", path))
		a.cursors.spreadsheet.Cursor.Common.Logs.Log(fmt.Sprintf("ERR: Failed to parse csv file '%s'", path))
		return
	}

	*a.dataGrid = *newDg
	a.dataGrid.UpdateColumnsWidth(a.assets)

	cur := a.cursors.spreadsheet.Cursor
	cur.Position.MaxCol = a.dataGrid.Cols - 1
	cur.Position.MaxRow = a.dataGrid.Rows - 1
	cur.Common.Logs.Log(fmt.Sprintf("Loaded csv file '%s'", path))
	slog.Info("Loaded csv file", slog.String("path", path))
}

func (a *App) handleQueryResults() {
	for _, connData := range a.connMgr.GetAllConnections() {
		if connData.Conn == nil || connData.ConnCtxDone == nil || connData.QueryChannel == nil {
			continue
		}

		newDg, done, err := connData.CheckForQueryResult()
		runtime := connData.GetQueryRuntimeDynamicString()
		logs := a.cursors.editor.Cursor.Common.Logs

		switch {
		case err != nil:
			slog.Error("Something went wrong during query", slog.Any("error", err))
			connData.Conn = nil
			logs.Log(fmt.Sprintf("'%s' cancelled after %s", connData.QueryText, runtime))

		case done && newDg == nil:
			slog.Debug("Query finished with nil result", slog.String("query", connData.QueryText))
			logs.Log(fmt.Sprintf("'%s' failed after %s", connData.QueryText, runtime))

		case done:
			slog.Debug("Query finished", slog.String("query", connData.QueryText))
			*a.dataGrid = *newDg
			a.dataGrid.UpdateColumnsWidth(a.assets)
			format.PrintMap(a.dataGrid.Data)

			a.cursors.editor.Cursor.Reset()
			a.cursors.spreadsheet.Cursor.Position.MaxCol = a.dataGrid.Cols
			a.cursors.spreadsheet.Cursor.Position.MaxRow = a.dataGrid.Rows
			logs.Log(fmt.Sprintf("'%s' finished after %s and returned %d result(s)", connData.QueryText, runtime, a.dataGrid.Rows))

		default:
			logs.Log(fmt.Sprintf("'%s' running for %s", connData.QueryText, runtime))
		}
	}
}

func initEditorContext(common *cursor.Common, connManager *database.ConnectionManager, eg *editor.Grid) *mode.Context {
	motions, commandRegistry := setup.EditorMotionSet()
	parser := motion.NewParser(motions.Root())
	cur := cursor.New(common, cursor.TypeEditor)

	return &mode.Context{
		Cursor:      cur,
		Parser:      parser,
		Commands:    commandRegistry,
		ConnManager: connManager,
		EditorGrid:  eg,
	}
}

func initSpreadsheetContext(common *cursor.Common) *mode.Context {
	motions, commandRegistry := setup.SpreadsheetMotionSet()
	parser := motion.NewParser(motions.Root())
	cur := cursor.New(common, cursor.TypeSpreadsheet)

	return &mode.Context{
		Cursor:   cur,
		Parser:   parser,
		Commands: commandRegistry,
	}
}

func initConnectionsContext(common *cursor.Common) *mode.Context {
	motions := setup.ConnectionsMotionSet()
	parser := motion.NewParser(motions.Root())
	cur := cursor.New(common, cursor.TypeConnections)

	return &mode.Context{
		Cursor: cur,
		Parser: parser,
	}
}
