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
	"github.com/quar15/qq-go/internal/database"
	"github.com/quar15/qq-go/internal/display"
	"github.com/quar15/qq-go/internal/format"
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

func handleDropFiles(appAssets *assets.Assets, dg *database.DataGrid, eg *display.EditorGrid) {
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
			newEg, err := display.LoadEditorGridFromTextFile(droppedFilesPaths[0], appAssets)
			if err != nil {
				slog.Error("Failed to load file", slog.Any("path", droppedFilesPaths[0]))
				display.CursorEditor.Common.Logs.Channel <- fmt.Sprintf("ERR: Failed to parse sql file '%s'", droppedFilesPaths[0])
			} else {
				*eg = *newEg
				display.CursorEditor.Handler.Reset(display.CursorSpreadsheet)
				display.CursorEditor.Position.MaxCol = 1
				display.CursorEditor.Position.MaxRow = eg.Rows
				display.CursorEditor.Common.Logs.Channel <- fmt.Sprintf("Loaded sql file '%s'", droppedFilesPaths[0])
			}
		case ".csv":
			slog.Info("Loaded csv file")
			newDg, err := database.LoadDataGridFromCSV(droppedFilesPaths[0], appAssets)
			if err != nil {
				slog.Error("Failed to load file", slog.Any("path", droppedFilesPaths[0]))
				display.CursorSpreadsheet.Common.Logs.Channel <- fmt.Sprintf("ERR: Failed to parse csv file '%s'", droppedFilesPaths[0])
			} else {
				*dg = *newDg
				dg.UpdateColumnsWidth(appAssets)
				display.CursorSpreadsheet.Handler.Reset(display.CursorSpreadsheet)
				display.CursorSpreadsheet.Position.MaxCol = dg.Cols - 1
				display.CursorSpreadsheet.Position.MaxRow = dg.Rows - 1
				display.CursorSpreadsheet.Common.Logs.Channel <- fmt.Sprintf("Loaded csv file '%s'", droppedFilesPaths[0])
			}
		default:
			slog.Warn("Unhandled type of file", slog.String("extension", ext))
		}
	}
}

func handleQuery(appAssets *assets.Assets, cursor *display.Cursor, dg *database.DataGrid, connManager *database.ConnectionManager) {
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
				cursor.Handler.Reset(cursor)
				display.CursorSpreadsheet.Position.MaxCol = dg.Cols
				display.CursorSpreadsheet.Position.MaxRow = dg.Rows
				cursor.Common.Logs.Channel <- fmt.Sprintf("'%s' finished after %s and returned %d result(s)", connData.QueryText, connData.GetQueryRuntimeDynamicString(), dg.Rows)
				// @TODO: Add system notification for query finish
			}
		} else {
			cursor.Common.Logs.Channel <- fmt.Sprintf("'%s' running for %s", connData.QueryText, connData.GetQueryRuntimeDynamicString())
		}
	}
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

	display.CursorEditor.Handler.Init(display.CursorEditor, &topZone)
	display.CursorSpreadsheet.Handler.Init(display.CursorSpreadsheet, &bottomZone)
	display.CursorConnection.Handler.Init(display.CursorConnection, &connectionsZone)
	display.CursorConnection.Position.MaxRow = int32(connMgr.GetNumberOfConnections()) - 1
	display.CurrCursor = display.CursorEditor
	display.CurrCursor.SetActive(true)

	topZone.ContentSize = rl.Vector2{X: 1600, Y: 1200}
	bottomZone.ContentSize = rl.Vector2{X: 2000, Y: 2000}
	commandZone.ContentSize = rl.Vector2{X: 0, Y: 0}
	connectionsZone.ContentSize = rl.Vector2{X: 0, Y: 0}
	var commandZoneHeight float32 = (appAssets.MainFontSize*2 + appAssets.MainFontSpacing*2)

	var dg database.DataGrid
	var eg display.EditorGrid = display.NewEditorGrid()

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

		handleDropFiles(&appAssets, &dg, &eg)
		handleQuery(&appAssets, display.CurrCursor, &dg, connMgr)
		display.CurrCursor.Handler.HandleInput(&appAssets, &dg, &eg, display.CurrCursor, connMgr)

		// --- Drawing ---
		rl.BeginDrawing()
		rl.ClearBackground(cfg.Colors.Background())

		editorIsFocused := (display.CurrCursor.Type == display.CursorTypeEditor)
		topZone.DrawEditor(&appAssets, &eg, display.CursorEditor, editorIsFocused)
		bottomZone.DrawSpreadsheetZone(&appAssets, &dg, display.CursorSpreadsheet)
		commandZone.DrawCommandZone(&appAssets, display.CurrCursor, connMgr.GetCurrentConnectionName())

		splitter.Draw()

		if display.CurrCursor.Type == display.CursorTypeConnections {
			connectionsZone.DrawConnectionSelector(&appAssets, cfg, display.CursorConnection, int32(screenWidth), int32(screenHeight), connMgr)
		}

		rl.EndDrawing()
	}
}
