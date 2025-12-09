package display

import (
	"log/slog"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/quar15/qq-go/internal/assets"
	"github.com/quar15/qq-go/internal/database"
	"github.com/quar15/qq-go/internal/utilities"
	"golang.design/x/clipboard"
)

func (SpreadsheetCursorStateHandler) HandleInput(appAssets *assets.Assets, dg *database.DataGrid, cursor *Cursor) {
	const cellHeight int32 = 30 // @TODO: Pass configuration with cellHeight
	var keyPressed int32 = rl.GetCharPressed()

	var mouse rl.Vector2 = rl.GetMousePosition()
	var mouseWheelStep int32 = 1
	// Only scroll if mouse inside the zone
	if rl.CheckCollisionPointRec(mouse, cursor.Zone.Bounds) {
		if rl.IsKeyDown(rl.KeyLeftShift) {
			// Mouse wheel scroll (horizontal)
			cursor.Position.Col -= int8(rl.GetMouseWheelMove()) * int8(mouseWheelStep)
		} else {
			// Mouse wheel scroll (vertical)
			cursor.Position.Row -= int32(rl.GetMouseWheelMove()) * mouseWheelStep
		}
	}

	switch cursor.Common.Mode {
	case ModeVLine:
		fallthrough
	case ModeVisual:
		fallthrough
	case ModeNormal:
		if rl.IsKeyDown(rl.KeyLeftShift) {
			switch {
			case rl.IsKeyPressed(rl.KeyHome):
				cursor.Position.Col = 0
			case rl.IsKeyPressed(rl.KeyEnd):
				cursor.Position.Col = dg.Cols - 1
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
				err := database.QueryData(database.CurrDBConnection.Name, query)
				if err != nil {
					slog.Error("Failed to execute query", slog.Any("error", err))
					cursor.Common.Logs.Channel <- "Failed to execute query (Something went wrong)"
				}
			case rl.IsKeyPressed(rl.KeyC):
				if cursor.Position.Col >= 0 && dg.Cols > 0 {
					// @TODO: Add type specific formatting
					var dataString string = ""
					if cursor.Common.Mode == ModeVisual || cursor.Common.Mode == ModeVLine {
						for row := cursor.Position.SelectStartRow; row <= cursor.Position.SelectEndRow; row++ {
							for col := cursor.Position.SelectStartCol; col < cursor.Position.SelectEndCol; col++ {
								dataString += utilities.GetValueAsString(dg.Data[row][dg.Headers[col]]) + ","
							}
							dataString += utilities.GetValueAsString(dg.Data[row][dg.Headers[cursor.Position.SelectEndCol]]) + "\n"
						}
					} else {
						dataString = utilities.GetValueAsString(dg.Data[cursor.Position.Row][dg.Headers[cursor.Position.Col]])
					}
					clipboard.Write(clipboard.FmtText, []byte(dataString))
				}
			case rl.IsKeyPressed(rl.KeyE):
				CursorSpreadsheet.TransitionMode(ModeNormal)
				CurrCursor = CursorConnection
			case rl.IsKeyPressed(rl.KeyD):
				slog.Debug("Connections debug", slog.Any("Position", CursorConnection.Position))
			}
		} else {
			var pageRows int8 = cursor.Zone.GetNumberOfVisibleRows(int32(cellHeight))
			switch {
			case rl.IsKeyPressed(rl.KeyEscape) || rl.IsKeyPressed(rl.KeyCapsLock): // @TODO: Remove personal preference CapsLock
				cursor.Handler.Reset(cursor)
			case utilities.Contains(HANDLED_MOTION_KEY_CODES, int(keyPressed)):
				cursor.AppendMotion(rune(keyPressed))
			case rl.IsKeyPressed(rl.KeySlash):
				cursor.Common.Mode = ModeCommand
				cursor.Common.CmdBuf = "/"
				cursor.Common.MotionBuf = "/"
				cursor.UpdateCmdLine()
			case rl.IsKeyPressed(rl.KeyHome):
				cursor.Position.Row = 0
			case rl.IsKeyPressed(rl.KeyEnd):
				cursor.Position.Row = dg.Rows - 1
			case rl.IsKeyPressed(rl.KeyPageUp):
				cursor.Position.Row -= int32(pageRows)
			case rl.IsKeyPressed(rl.KeyPageDown):
				cursor.Position.Row += int32(pageRows)
			case rl.IsKeyPressed(rl.KeyDown):
				cursor.Position.Row++
			case rl.IsKeyPressed(rl.KeyUp):
				cursor.Position.Row--
			case rl.IsKeyPressed(rl.KeyLeft):
				cursor.Position.Col--
			case rl.IsKeyPressed(rl.KeyRight):
				cursor.Position.Col++
			}
		}
	case ModeCommand:
		switch {
		case rl.IsKeyPressed(rl.KeyEnter):
			cursor.ExecuteCommand()
		case rl.IsKeyPressed(rl.KeyBackspace):
			cursor.Common.CmdBuf = cursor.Common.CmdBuf[:len(cursor.Common.CmdBuf)-1]
			cursor.UpdateCmdLine()
			if len(cursor.Common.CmdBuf) <= 0 {
				cursor.Handler.Reset(cursor)
			}
		case rl.IsKeyPressed(rl.KeyEscape) || rl.IsKeyPressed(rl.KeyCapsLock): // @TODO: Remove personal preference CapsLock
			cursor.Handler.Reset(cursor)
		case keyPressed >= 32 && keyPressed <= 125:
			cursor.Common.CmdBuf += string(rune(keyPressed))
			cursor.UpdateCmdLine()
		}
	}

	cursor.Zone.ClampScrollsToZoneSize()
	cursor.ClampFocus(dg.Cols-1, dg.Rows-1)
	cursor.UpdateSelectBasedOnPosition()
}

func (ConnectionsCursorStateHandler) HandleInput(appAssets *assets.Assets, dg *database.DataGrid, cursor *Cursor) {
	var keyPressed int32 = rl.GetCharPressed()

	switch cursor.Common.Mode {
	case ModeNormal:
		if rl.IsKeyDown(rl.KeyLeftShift) {
			switch {
			case rl.IsKeyPressed(rl.KeyHome):
				cursor.Position.Col = 0
			case rl.IsKeyPressed(rl.KeyEnd):
				cursor.Position.Col = dg.Cols - 1
			case utilities.Contains(HANDLED_MOTION_KEY_CODES, int(keyPressed)):
				cursor.AppendMotion(rune(keyPressed))
			}
		} else if rl.IsKeyDown(rl.KeyLeftControl) {
			switch {
			case rl.IsKeyPressed(rl.KeyE):
				CursorConnection.TransitionMode(ModeNormal)
				CurrCursor = CursorSpreadsheet
			case rl.IsKeyPressed(rl.KeyD):
				slog.Debug("Connections debug", slog.Any("Position", CursorConnection.Position))
			}
		} else {
			switch {
			case rl.IsKeyPressed(rl.KeyEscape) || rl.IsKeyPressed(rl.KeyCapsLock): // @TODO: Remove personal preference CapsLock
				cursor.Handler.Reset(cursor)
			case utilities.Contains(HANDLED_MOTION_KEY_CODES, int(keyPressed)):
				cursor.AppendMotion(rune(keyPressed)) // @TODO: Disable visual mode for connection selection
			case rl.IsKeyPressed(rl.KeySlash):
				cursor.Common.Mode = ModeCommand
				cursor.Common.CmdBuf = "/"
				cursor.Common.MotionBuf = "/"
				cursor.UpdateCmdLine()
			case rl.IsKeyPressed(rl.KeyHome):
				cursor.Position.Row = 0
			case rl.IsKeyPressed(rl.KeyEnd):
				cursor.Position.Row = dg.Rows - 1
			case rl.IsKeyPressed(rl.KeyDown):
				cursor.Position.Row++
			case rl.IsKeyPressed(rl.KeyUp):
				cursor.Position.Row--
			case rl.IsKeyPressed(rl.KeyLeft):
				cursor.Position.Col--
			case rl.IsKeyPressed(rl.KeyRight):
				cursor.Position.Col++
			}
		}
	case ModeCommand:
		switch {
		case rl.IsKeyPressed(rl.KeyEnter):
			cursor.ExecuteCommand()
		case rl.IsKeyPressed(rl.KeyBackspace):
			cursor.Common.CmdBuf = cursor.Common.CmdBuf[:len(cursor.Common.CmdBuf)-1]
			cursor.UpdateCmdLine()
			if len(cursor.Common.CmdBuf) <= 0 {
				cursor.Handler.Reset(cursor)
			}
		case rl.IsKeyPressed(rl.KeyEscape) || rl.IsKeyPressed(rl.KeyCapsLock): // @TODO: Remove personal preference CapsLock
			cursor.Handler.Reset(cursor)
		case keyPressed >= 32 && keyPressed <= 125:
			cursor.Common.CmdBuf += string(rune(keyPressed))
			cursor.UpdateCmdLine()
		}
	}

	// cursor.Zone.ClampScrollsToZoneSize()
	cursor.ClampFocus(0, cursor.Position.MaxRow)
	cursor.UpdateSelectBasedOnPosition()
}

func (EditorCursorStateHandler) HandleInput(appAssets *assets.Assets, dg *database.DataGrid, cursor *Cursor) {
	slog.Debug("Handling input by EditorCursorStateHandler")
}
