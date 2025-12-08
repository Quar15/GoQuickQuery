package display

import (
	"log/slog"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/quar15/qq-go/internal/assets"
	"github.com/quar15/qq-go/internal/database"
	"github.com/quar15/qq-go/internal/utilities"
	"golang.design/x/clipboard"
)

func HandleSpreadsheetInput(z *Zone, dg *database.DataGrid, cursor *SpreadsheetCursor, appAssets *assets.Assets, cellHeight int32) {
	var keyPressed int32 = rl.GetCharPressed()

	var mouse rl.Vector2 = rl.GetMousePosition()
	var mouseWheelStep int32 = 1
	// Only scroll if mouse inside the zone
	if rl.CheckCollisionPointRec(mouse, z.Bounds) {
		if rl.IsKeyDown(rl.KeyLeftShift) {
			// Mouse wheel scroll (horizontal)
			cursor.Col -= int8(rl.GetMouseWheelMove()) * int8(mouseWheelStep)
		} else {
			// Mouse wheel scroll (vertical)
			cursor.Row -= int32(rl.GetMouseWheelMove()) * mouseWheelStep
		}
	}

	switch cursor.Mode {
	case ModeVLine:
		fallthrough
	case ModeVisual:
		fallthrough
	case ModeNormal:
		if rl.IsKeyDown(rl.KeyLeftShift) {
			switch {
			case rl.IsKeyPressed(rl.KeyHome):
				cursor.Col = 0
			case rl.IsKeyPressed(rl.KeyEnd):
				cursor.Col = dg.Cols - 1
			case utilities.Contains(HANDLED_MOTION_KEY_CODES, int(keyPressed)):
				cursor.AppendMotion(rune(keyPressed))
			}
		} else if rl.IsKeyDown(rl.KeyLeftControl) {
			switch {
			case rl.IsKeyPressed(rl.KeyEnter):
				// @TODO: Get query from editor (temp hard code)
				//query := "SELECT 1;"
				//query := "SELECT pg_sleep(20)"
				query := "SELECT * FROM example LIMIT 500;"
				err := database.QueryData("postgres", query)
				if err != nil {
					slog.Error("Failed to execute query", slog.Any("error", err))
					cursor.Logs.Channel <- "Failed to execute query (Something went wrong)"
				}
			case rl.IsKeyDown(rl.KeyC):
				if cursor.Col >= 0 && dg.Cols > 0 {
					// @TODO: Add type specific formatting
					var dataString string = ""
					if cursor.Mode == ModeVisual || cursor.Mode == ModeVLine {
						for row := cursor.SelectStartRow; row <= cursor.SelectEndRow; row++ {
							for col := cursor.SelectStartCol; col < cursor.SelectEndCol; col++ {
								dataString += utilities.GetValueAsString(dg.Data[row][dg.Headers[col]]) + ","
							}
							dataString += utilities.GetValueAsString(dg.Data[row][dg.Headers[cursor.SelectEndCol]]) + "\n"
						}
					} else {
						dataString = utilities.GetValueAsString(dg.Data[cursor.Row][dg.Headers[cursor.Col]])
					}
					clipboard.Write(clipboard.FmtText, []byte(dataString))
				}
			}
		} else {
			var pageRows int8 = z.GetNumberOfVisibleRows(int32(cellHeight))
			switch {
			case rl.IsKeyPressed(rl.KeyEscape) || rl.IsKeyPressed(rl.KeyCapsLock): // @TODO: Remove personal preference CapsLock
				cursor.Reset()
			case utilities.Contains(HANDLED_MOTION_KEY_CODES, int(keyPressed)):
				cursor.AppendMotion(rune(keyPressed))
			case rl.IsKeyPressed(rl.KeySlash):
				cursor.Mode = ModeCommand
				cursor.CmdBuf = "/"
				cursor.MotionBuf = "/"
				cursor.UpdateCmdLine()
			case rl.IsKeyPressed(rl.KeyHome):
				cursor.Row = 0
			case rl.IsKeyPressed(rl.KeyEnd):
				cursor.Row = dg.Rows - 1
			case rl.IsKeyPressed(rl.KeyPageUp):
				cursor.Row -= int32(pageRows)
			case rl.IsKeyPressed(rl.KeyPageDown):
				cursor.Row += int32(pageRows)
			case rl.IsKeyPressed(rl.KeyDown):
				cursor.Row++
			case rl.IsKeyPressed(rl.KeyUp):
				cursor.Row--
			case rl.IsKeyPressed(rl.KeyLeft):
				cursor.Col--
			case rl.IsKeyPressed(rl.KeyRight):
				cursor.Col++
			}
		}
	case ModeCommand:
		switch {
		case rl.IsKeyPressed(rl.KeyEnter):
			cursor.ExecuteCommand()
		case rl.IsKeyPressed(rl.KeyBackspace):
			cursor.CmdBuf = cursor.CmdBuf[:len(cursor.CmdBuf)-1]
			cursor.UpdateCmdLine()
			if len(cursor.CmdBuf) <= 0 {
				cursor.Reset()
			}
		case rl.IsKeyPressed(rl.KeyEscape) || rl.IsKeyPressed(rl.KeyCapsLock): // @TODO: Remove personal preference CapsLock
			cursor.Reset()
		case keyPressed >= 32 && keyPressed <= 125:
			cursor.CmdBuf += string(rune(keyPressed))
			cursor.UpdateCmdLine()
		}
	}

	z.ClampScrollsToZoneSize()
	cursor.ClampFocus(dg.Cols-1, dg.Rows-1)
	cursor.UpdateSelectBasedOnPosition()
}
