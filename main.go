package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/gen2brain/raylib-go/raylib"
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

func initialize() error {
	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stdout, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.RFC3339,
		}),
	))
	err := clipboard.Init()
	if err != nil {
		slog.Error("Failed to initialize clipboard", slog.Any("error", err))
		return err
	}

	return nil
}

func handleDropFiles(appAssets *assets.Assets, editorCursor *cursor.Cursor, spreadsheetCursor *cursor.Cursor, dg *database.DataGrid, eg *editor.Grid) {
	if rl.IsFileDropped() {
		droppedFilesPaths := rl.LoadDroppedFiles()
		// @TODO: Extend to handle every file if tabs are added
		if len(droppedFilesPaths) > 1 {
			slog.Warn("Tried to load more than one file")
		}
		ext := filepath.Ext(droppedFilesPaths[0])
		switch ext {
		case ".sql":
			slog.Info("Loaded sql file")
			newEg, err := editor.LoadGridFromTextFile(droppedFilesPaths[0], appAssets)
			if err != nil {
				slog.Error("Failed to load file", slog.Any("path", droppedFilesPaths[0]))
				editorCursor.Common.Logs.Channel <- fmt.Sprintf("ERR: Failed to parse sql file '%s'", droppedFilesPaths[0])
			} else {
				*eg = *newEg
				editorCursor.Reset()
				editorCursor.Position.MaxCol = eg.MaxCol
				editorCursor.Position.MaxColForRows = eg.Cols
				editorCursor.Position.MaxRow = eg.Rows - 1
				editorCursor.Common.Logs.Channel <- fmt.Sprintf("Loaded sql file '%s'", droppedFilesPaths[0])
				//editorCursor.Common.Logs.Channel <- fmt.Sprintf("Loaded sql file '%s'", droppedFilesPaths[0])
			}
		case ".csv":
			slog.Info("Loaded csv file")
			newDg, err := database.LoadDataGridFromCSV(droppedFilesPaths[0], appAssets)
			if err != nil {
				slog.Error("Failed to load file", slog.Any("path", droppedFilesPaths[0]))
				spreadsheetCursor.Common.Logs.Channel <- fmt.Sprintf("ERR: Failed to parse csv file '%s'", droppedFilesPaths[0])
			} else {
				*dg = *newDg
				dg.UpdateColumnsWidth(appAssets)
				//display.CursorSpreadsheet.Handler.Reset(display.CursorSpreadsheet)
				spreadsheetCursor.Position.MaxCol = dg.Cols - 1
				spreadsheetCursor.Position.MaxRow = dg.Rows - 1
				spreadsheetCursor.Common.Logs.Channel <- fmt.Sprintf("Loaded csv file '%s'", droppedFilesPaths[0])
			}
		default:
			slog.Warn("Unhandled type of file", slog.String("extension", ext))
		}
	}
}

func handleQuery(appAssets *assets.Assets, cursor *cursor.Cursor, spreadsheetCursor *cursor.Cursor, dg *database.DataGrid, connManager *database.ConnectionManager) {
	for _, connData := range connManager.GetAllConnections() {
		if connData.Conn == nil || connData.ConnCtx == nil || connData.QueryChannel == nil {
			continue
		}
		newDg, done, err := connData.CheckForQueryResult()
		if err != nil {
			slog.Error("Something went wrong during query", slog.Any("error", err))
			connData.Conn = nil
			cursor.Common.Logs.Channel <- fmt.Sprintf("'%s' cancelled after %s", connData.QueryText, connData.GetQueryRuntimeDynamicString())
		} else if done == true {
			slog.Debug("Query finished", slog.String("query", connData.QueryText), slog.Any("dg", newDg), slog.Any("error", err))
			if newDg == nil {
				cursor.Common.Logs.Channel <- fmt.Sprintf("'%s' failed after %s", connData.QueryText, connData.GetQueryRuntimeDynamicString())
			} else {
				*dg = *newDg
				dg.UpdateColumnsWidth(appAssets)
				format.PrintMap(dg.Data)
				cursor.Reset()
				spreadsheetCursor.Position.MaxCol = dg.Cols
				spreadsheetCursor.Position.MaxRow = dg.Rows
				cursor.Common.Logs.Channel <- fmt.Sprintf("'%s' finished after %s and returned %d result(s)", connData.QueryText, connData.GetQueryRuntimeDynamicString(), dg.Rows)
				// @TODO: Add system notification for query finish
			}
		} else {
			cursor.Common.Logs.Channel <- fmt.Sprintf("'%s' running for %s", connData.QueryText, connData.GetQueryRuntimeDynamicString())
		}
	}
}

func initEditorCursor(common *cursor.Common, connManager *database.ConnectionManager, eg *editor.Grid) *mode.Context {
	motions, commandRegistry := setup.EditorMotionSet()
	parser := motion.NewParser(motions.Root())

	cur := cursor.New(common, cursor.TypeEditor)

	ctx := &mode.Context{
		Cursor:      cur,
		Parser:      parser,
		Commands:    commandRegistry,
		ConnManager: connManager,
		EditorGrid:  eg,
	}

	return ctx
}

func initSpreadsheetCursor(common *cursor.Common) *mode.Context {
	motions, commandRegistry := setup.SpreadsheetMotionSet()
	parser := motion.NewParser(motions.Root())

	cur := cursor.New(common, cursor.TypeSpreadsheet)

	ctx := &mode.Context{
		Cursor:   cur,
		Parser:   parser,
		Commands: commandRegistry,
	}

	return ctx
}

func initConnectionsCursor(common *cursor.Common) *mode.Context {
	motions := setup.ConnectionsMotionSet()
	parser := motion.NewParser(motions.Root())

	cur := cursor.New(common, cursor.TypeConnections)

	ctx := &mode.Context{
		Cursor: cur,
		Parser: parser,
	}

	return ctx
}

func main() {
	err := initialize()
	if err != nil {
		panic("Failed to initialize")
	}
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to read config", slog.Any("error", err))
		os.Exit(1)
	}
	connMgr := database.NewConnectionManager(cfg.Connections, &database.DefaultConnectionFactory{})
	defer connMgr.Close(context.Background())

	var dg database.DataGrid
	var eg editor.Grid = editor.NewGrid()

	// --- Init Window ---
	var screenWidth int = rl.GetScreenWidth()
	var screenHeight int = rl.GetScreenHeight()
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(int32(screenWidth), int32(screenHeight), "GQQ")
	defer rl.CloseWindow()
	rl.SetExitKey(rl.KeyNull)

	var appAssets assets.Assets
	if err := appAssets.LoadAssets(); err != nil {
		slog.Error("Failed to load assets", slog.Any("error", err))
		return
	}
	defer appAssets.UnloadAssets()

	splitter := display.Splitter{
		Ratio:    0.6,
		Height:   6.0,
		Dragging: false,
	}

	var topZone display.Zone
	var bottomZone display.Zone
	var commandZone display.Zone
	var connectionsZone display.Zone

	cursorCommon := &cursor.Common{}
	cursorCommon.Logs.Init()
	editorCursorCtx := initEditorCursor(cursorCommon, connMgr, &eg)
	editorCursorCtx.Cursor.Activate()
	spreadsheetCursorCtx := initSpreadsheetCursor(cursorCommon)
	connectionsCursorCtx := initConnectionsCursor(cursorCommon)

	topZone.ContentSize = rl.Vector2{X: 1920, Y: 540}
	bottomZone.ContentSize = rl.Vector2{X: 1920, Y: 540}
	commandZone.ContentSize = rl.Vector2{X: 0, Y: 0}
	connectionsZone.ContentSize = rl.Vector2{X: 0, Y: 0}
	var commandZoneHeight float32 = (appAssets.MainFontSize*2 + appAssets.MainFontSpacing*2)

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		// --- Pre Drawing ---
		screenWidth = max(rl.GetScreenWidth(), 100)
		screenHeight = max(rl.GetScreenHeight(), 100)

		rl.SetMouseCursor(rl.MouseCursorDefault)
		splitter.HandleZoneSplit(screenWidth, screenHeight, int(commandZoneHeight))

		topZone.Bounds = rl.Rectangle{X: 0, Y: 0, Width: float32(screenWidth), Height: splitter.Y - splitter.Height/2}
		bottomZone.Bounds = rl.Rectangle{X: 0, Y: splitter.Y + splitter.Height/2, Width: float32(screenWidth), Height: float32(screenHeight) - (splitter.Y + splitter.Height/2) - commandZoneHeight}
		commandZone.Bounds = rl.Rectangle{X: 0, Y: bottomZone.Bounds.Y + bottomZone.Bounds.Height, Width: float32(screenWidth), Height: commandZoneHeight}

		topZone.UpdateZoneScroll()
		bottomZone.UpdateZoneScroll()

		handleDropFiles(&appAssets, editorCursorCtx.Cursor, spreadsheetCursorCtx.Cursor, &dg, &eg)
		handleQuery(&appAssets, editorCursorCtx.Cursor, spreadsheetCursorCtx.Cursor, &dg, connMgr)
		display.HandleInput(editorCursorCtx)

		// --- Drawing ---
		rl.BeginDrawing()
		rl.ClearBackground(cfg.Colors.Background())

		editorIsFocused := editorCursorCtx.Cursor.IsActive()
		topZone.DrawEditor(&appAssets, &eg, editorCursorCtx.Cursor, editorIsFocused)
		bottomZone.DrawSpreadsheetZone(&appAssets, &dg, spreadsheetCursorCtx.Cursor)
		commandZone.DrawCommandZone(&appAssets, editorCursorCtx.Cursor, connMgr.GetCurrentConnectionName())

		splitter.Draw(editorCursorCtx.Cursor.Type)

		if connectionsCursorCtx.Cursor.IsActive() {
			connectionsZone.DrawConnectionSelector(&appAssets, cfg, connectionsCursorCtx.Cursor, int32(screenWidth), int32(screenHeight), connMgr)
		}

		rl.EndDrawing()
	}
}
