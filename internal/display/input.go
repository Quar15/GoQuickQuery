package display

import (
	"context"
	"log/slog"
	"slices"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/quar15/qq-go/internal/assets"
	"github.com/quar15/qq-go/internal/config"
	"github.com/quar15/qq-go/internal/database"
	"github.com/quar15/qq-go/internal/format"
	"golang.design/x/clipboard"
)

const keySmallG int = 103
const keySmallH int = 104
const keySmallJ int = 106
const keySmallK int = 107
const keySmallL int = 108
const keySmallV int = 118
const keySmallW int = 119

func (SpreadsheetCursorStateHandler) HandleInput(appAssets *assets.Assets, dg *database.DataGrid, eg *EditorGrid, cursor *Cursor, connManager *database.ConnectionManager) {
	const cellHeight int32 = 30 // @TODO: Pass configuration with cellHeight
	var keyPressed int32 = rl.GetCharPressed()

	var mouse rl.Vector2 = rl.GetMousePosition()
	var mouseWheelStep int32 = 1
	// Only scroll if mouse inside the zone
	if rl.CheckCollisionPointRec(mouse, cursor.Zone.Bounds) {
		if rl.IsKeyDown(rl.KeyLeftShift) {
			// Mouse wheel scroll (horizontal)
			cursor.Position.Col -= int32(rl.GetMouseWheelMove()) * int32(mouseWheelStep)
		} else {
			// Mouse wheel scroll (vertical)
			cursor.Position.Row -= int32(rl.GetMouseWheelMove()) * mouseWheelStep
		}
	}

	switch cursor.Common.Mode {
	case ModeVLine:
		fallthrough
	case ModeVBlock:
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
			case slices.Contains(HANDLED_MOTION_KEY_CODES, int(keyPressed)):
				cursor.AppendMotion(rune(keyPressed))
			}
		} else if rl.IsKeyDown(rl.KeyLeftControl) {
			switch {
			case rl.IsKeyPressed(rl.KeyC):
				if cursor.Position.Col >= 0 && dg.Cols > 0 {
					// @TODO: Add type specific formatting
					var dataString string = ""
					switch cursor.Common.Mode {
					case ModeVisual:
						for row := cursor.Position.SelectStartRow; row <= cursor.Position.SelectEndRow; row++ {
							for col := int32(0); col < dg.Cols; col++ {
								if cursor.IsSelected(col, row) {
									dataString += format.GetValueAsString(dg.Data[row][dg.Headers[col]]) + ","
								}
							}
							dataString += "\n"
						}
					case ModeVLine:
						fallthrough
					case ModeVBlock:
						for row := cursor.Position.SelectStartRow; row <= cursor.Position.SelectEndRow; row++ {
							for col := cursor.Position.SelectStartCol; col < cursor.Position.SelectEndCol; col++ {
								dataString += format.GetValueAsString(dg.Data[row][dg.Headers[col]]) + ","
							}
							dataString += format.GetValueAsString(dg.Data[row][dg.Headers[cursor.Position.SelectEndCol]]) + "\n"
						}
					case ModeNormal:
						dataString = format.GetValueAsString(dg.Data[cursor.Position.Row][dg.Headers[cursor.Position.Col]])
					}
					slog.Debug("Copied to clipboard from spreadsheet", slog.String("dataString", dataString))
					clipboard.Write(clipboard.FmtText, []byte(dataString))
				}
			case rl.IsKeyPressed(rl.KeyE):
				CursorSpreadsheet.SetActive(false)
				CursorConnection.SetActive(true)
				CurrCursor = CursorConnection
			case rl.IsKeyPressed(rl.KeyV):
				cursor.AppendMotionString("^V")
			case rl.IsKeyPressed(rl.KeyW):
				cursor.AppendMotionString("^W")
			}
		} else {
			var pageRows int8 = cursor.Zone.GetNumberOfVisibleRows(int32(cellHeight))
			switch {
			case rl.IsKeyPressed(rl.KeyEscape) || rl.IsKeyPressed(rl.KeyCapsLock): // @TODO: Remove personal preference CapsLock
				cursor.Handler.Reset(cursor)
			case slices.Contains(HANDLED_MOTION_KEY_CODES, int(keyPressed)):
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

func (ConnectionsCursorStateHandler) HandleInput(appAssets *assets.Assets, dg *database.DataGrid, eg *EditorGrid, cursor *Cursor, connManager *database.ConnectionManager) {
	var keyPressed int32 = rl.GetCharPressed()

	switch cursor.Common.Mode {
	case ModeNormal:
		if rl.IsKeyDown(rl.KeyLeftShift) {
			switch {
			case slices.Contains(HANDLED_MOTION_KEY_CODES, int(keyPressed)):
				cursor.AppendMotion(rune(keyPressed))
			case rl.IsKeyPressed(rl.KeyHome):
				cursor.Position.Col = 0
			case rl.IsKeyPressed(rl.KeyEnd):
				cursor.Position.Col = dg.Cols - 1
			}
		} else if rl.IsKeyDown(rl.KeyLeftControl) {
			switch {
			case rl.IsKeyPressed(rl.KeyE):
				CursorConnection.SetActive(false)
				CursorEditor.SetActive(true)
				CurrCursor = CursorEditor
			case rl.IsKeyPressed(rl.KeyD):
				slog.Debug("Connections debug", slog.Any("Position", CursorConnection.Position))
			}
		} else {
			switch {
			case rl.IsKeyPressed(rl.KeyEnter):
				err := connManager.SetCurrentConnectionByName(config.Cfg.Connections[CursorConnection.Position.Row].Name)
				if err != nil {
					slog.Error("Failed to set current connection by name", slog.Any("error", err))
				}
				// @TODO: Go back to last focused cursor
				CursorConnection.SetActive(false)
				CursorEditor.SetActive(true)
				CurrCursor = CursorEditor
			case rl.IsKeyPressed(rl.KeyEscape) || rl.IsKeyPressed(rl.KeyCapsLock): // @TODO: Remove personal preference CapsLock
				// @TODO: Go back to last focused cursor
				CursorConnection.SetActive(false)
				CursorEditor.SetActive(true)
				CurrCursor = CursorEditor
				CurrCursor = CursorSpreadsheet
			case slices.Contains(HANDLED_MOTION_KEY_CODES, int(keyPressed)):
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

	cursor.Zone.ClampScrollsToZoneSize()
	cursor.ClampFocus(0, cursor.Position.MaxRow)
	cursor.UpdateSelectBasedOnPosition()
}

func (EditorCursorStateHandler) HandleInput(appAssets *assets.Assets, dg *database.DataGrid, eg *EditorGrid, cursor *Cursor, connManager *database.ConnectionManager) {
	const cellHeight int32 = 30 // @TODO: Pass configuration with cellHeight
	var keyPressed int32 = rl.GetCharPressed()

	var mouse rl.Vector2 = rl.GetMousePosition()
	var mouseWheelStep int32 = 1
	// Only scroll if mouse inside the zone
	if rl.CheckCollisionPointRec(mouse, cursor.Zone.Bounds) {
		if rl.IsKeyDown(rl.KeyLeftShift) {
			// Mouse wheel scroll (horizontal)
			cursor.Position.Col -= int32(rl.GetMouseWheelMove()) * int32(mouseWheelStep)
		} else {
			// Mouse wheel scroll (vertical)
			cursor.Position.Row -= int32(rl.GetMouseWheelMove()) * mouseWheelStep
		}
	}

	switch cursor.Common.Mode {
	case ModeVLine:
		fallthrough
	case ModeVBlock:
		fallthrough
	case ModeVisual:
		fallthrough
	case ModeNormal:
		if rl.IsKeyDown(rl.KeyLeftShift) {
			switch {
			case slices.Contains(HANDLED_MOTION_KEY_CODES, int(keyPressed)):
				cursor.AppendMotion(rune(keyPressed))
			}
		} else if rl.IsKeyDown(rl.KeyLeftControl) {
			switch {
			case rl.IsKeyPressed(rl.KeyEnter):
				// @TODO: Get query from editor (temp hard code)
				query := "SELECT 1;"
				//query := "SELECT pg_sleep(20)"
				//query := "SELECT * FROM example LIMIT 500;"
				err := connManager.ExecuteQuery(context.Background(), connManager.GetCurrentConnectionName(), query)
				if err != nil {
					slog.Error("Failed to execute query", slog.Any("error", err))
					cursor.Common.Logs.Channel <- "Failed to execute query (Something went wrong)"
				}
			case rl.IsKeyPressed(rl.KeyC):
				if cursor.Position.Col >= 0 && eg.Rows > 0 {
					var dataString string = ""
					switch cursor.Common.Mode {
					case ModeVisual:
						for row := cursor.Position.SelectStartRow; row <= cursor.Position.SelectEndRow; row++ {
							if eg.Cols[row] > 0 {
								for col := int32(0); col < eg.Cols[row]; col++ {
									if cursor.IsSelected(col, row) {
										dataString += string(eg.Text[row][col])
									}
								}
							}
							dataString += "\n"
						}
					case ModeVLine:
						for row := cursor.Position.SelectStartRow; row <= cursor.Position.SelectEndRow; row++ {
							dataString += eg.Text[row] + "\n"
						}
					case ModeVBlock:
						for row := cursor.Position.SelectStartRow; row <= cursor.Position.SelectEndRow; row++ {
							if eg.Cols[row] > 0 {
								endCol := min(cursor.Position.SelectEndCol, eg.Cols[row])
								for col := cursor.Position.SelectStartCol; col <= endCol; col++ {
									dataString += string(eg.Text[row][col])
								}
							}
							dataString += "\n"
						}
					default:
						dataString = eg.Text[cursor.Position.Row]
					}
					slog.Debug("Copied to clipboard from editor", slog.String("dataString", dataString))
					clipboard.Write(clipboard.FmtText, []byte(dataString))
				}
			case rl.IsKeyPressed(rl.KeyE):
				CursorEditor.SetActive(false)
				CursorConnection.SetActive(true)
				CurrCursor = CursorConnection
			case rl.IsKeyPressed(rl.KeyD):
				slog.Debug("Editor debug", slog.Any("Position", CursorConnection.Position), slog.Any("eg", eg))
			case rl.IsKeyPressed(rl.KeyV):
				cursor.AppendMotionString("^V")
			case rl.IsKeyPressed(rl.KeyW):
				cursor.AppendMotionString("^W")
			}
		} else {
			var pageRows int8 = cursor.Zone.GetNumberOfVisibleRows(int32(cellHeight))
			switch {
			case rl.IsKeyPressed(rl.KeyEscape) || rl.IsKeyPressed(rl.KeyCapsLock): // @TODO: Remove personal preference CapsLock
				cursor.Handler.Reset(cursor)
			case slices.Contains(HANDLED_MOTION_KEY_CODES, int(keyPressed)):
				cursor.AppendMotion(rune(keyPressed))
			case rl.IsKeyPressed(rl.KeySlash):
				cursor.Common.Mode = ModeCommand
				cursor.Common.CmdBuf = "/"
				cursor.Common.MotionBuf = "/"
				cursor.UpdateCmdLine()
			case rl.IsKeyPressed(rl.KeyHome):
				cursor.Position.Col = 0
			case rl.IsKeyPressed(rl.KeyEnd):
				cursor.Position.Col = eg.Cols[cursor.Position.Row]
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
	// Additional focus check to be in scope of Cols
	if cursor.Position.Row >= eg.Rows {
		cursor.Position.Row = eg.Rows - 1
	} else if cursor.Position.Row < 0 {
		cursor.Position.Row = 0
	}
	cursor.ClampFocus(eg.Cols[cursor.Position.Row]-1, eg.Rows)
	cursor.UpdateSelectBasedOnPosition()
}
