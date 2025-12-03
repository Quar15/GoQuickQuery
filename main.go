package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/gen2brain/raylib-go/raylib"
	"github.com/jackc/pgx/v5"
	"github.com/lmittmann/tint"
	"golang.design/x/clipboard"
	"gopkg.in/yaml.v3"

	"github.com/quar15/qq-go/internal/assets"
	"github.com/quar15/qq-go/internal/colors"
	"github.com/quar15/qq-go/internal/database"
	"github.com/quar15/qq-go/internal/display"
)

type Config struct {
	Connections []database.ConnectionData `yaml:"connections"`
}

func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func initialize(cfg *Config) error {
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
	database.InitializeConnections(cfg.Connections)

	return nil
}

func main() {
	cfg, err := loadConfig("./gqq.yaml")
	if err != nil {
		slog.Error("Failed to read config", slog.Any("error", err))
		os.Exit(1)
	}
	err = initialize(cfg)
	if err != nil {
		panic("Failed to initialize")
	}
	var spreadsheetCursor display.SpreadsheetCursor
	spreadsheetCursor.Init()

	// --- Init Window ---
	var screenWidth int = rl.GetScreenWidth()
	var screenHeight int = rl.GetScreenHeight()
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(int32(screenWidth), int32(screenHeight), "QQ")
	defer rl.CloseWindow()
	rl.SetExitKey(rl.KeyNull)

	var appAssets assets.Assets
	if err := appAssets.LoadAssets(); err != nil {
		slog.Error("Failed to load assets", slog.Any("error", err))
		return
	}
	defer appAssets.UnloadAssets()

	splitter := display.Splitter{
		Ratio:    0.01,
		Height:   6.0,
		Dragging: false,
	}

	var topZone display.Zone
	var bottomZone display.Zone
	var commandZone display.Zone

	topZone.ContentSize = rl.Vector2{X: 1600, Y: 1200}
	bottomZone.ContentSize = rl.Vector2{X: 2000, Y: 2000}
	commandZone.ContentSize = rl.Vector2{X: 0, Y: 0}
	var commandZoneHeight float32 = (appAssets.MainFontSize*2 + appAssets.MainFontSpacing*2)

	var dg database.DataGrid

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

		// @TODO: Fix to get content into consideration
		if int(topZone.ContentSize.X) < screenWidth {
			topZone.ContentSize = rl.Vector2{X: float32(screenWidth), Y: 1200}
		}

		topZone.UpdateZoneScroll()
		bottomZone.UpdateZoneScroll()

		// --- Dropping files ---
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
				// @TODO: Pass file to editor
			case ".csv":
				slog.Info("Loaded csv file")
				dg, err = database.LoadDataGridFromCSV(droppedFilesPaths[0], &appAssets)
				if err != nil {
					slog.Error("Failed to load file", slog.Any("path", droppedFilesPaths[0]))
					dg = database.DataGrid{} // @TODO: Consider using tmp variable to not stain already loaded data
					spreadsheetCursor.Logs.Channel <- fmt.Sprintf("ERR: Failed to parse csv file '%s'", droppedFilesPaths[0])
				} else {
					spreadsheetCursor.Logs.Channel <- fmt.Sprintf("Loaded csv file '%s'", droppedFilesPaths[0])
				}
				spreadsheetCursor.Reset()
			default:
				slog.Warn("Unhandled type of file", slog.String("extension", ext))
			}
		}

		// --- Drawing ---
		rl.BeginDrawing()
		rl.ClearBackground(colors.Background())

		topZone.Draw(&appAssets)
		bottomZone.DrawSpreadsheetZone(&appAssets, &dg, &spreadsheetCursor)
		commandZone.DrawCommandZone(&appAssets, &spreadsheetCursor)

		rl.DrawRectangleRec(splitter.Rect, colors.Crust())

		rl.EndDrawing()
	}

	for _, connData := range database.DBConnections {
		switch c := connData.Conn.(type) {
		case *pgx.Conn:
			c.Close(context.Background())
		}
	}
}
