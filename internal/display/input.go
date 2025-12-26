package display

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/quar15/qq-go/internal/assets"
	"github.com/quar15/qq-go/internal/config"
	"github.com/quar15/qq-go/internal/cursor"
	"github.com/quar15/qq-go/internal/database"
	"github.com/quar15/qq-go/internal/format"
	"github.com/quar15/qq-go/internal/mode"
	"github.com/quar15/qq-go/internal/motion"
	"golang.design/x/clipboard"
)

const keySmallG int = 103
const keySmallH int = 104
const keySmallJ int = 106
const keySmallK int = 107
const keySmallL int = 108
const keySmallV int = 118
const keySmallW int = 119

func HandleInput(ctx *mode.Context) {
	var (
		keyCharPressed int32          = rl.GetCharPressed()
		keyPressed     int32          = rl.GetKeyPressed()
		code           motion.KeyCode = motion.KeyRune
		modifiers      motion.Modifiers
		keyRune        rune
	)

	arrowKeys := []int32{rl.KeyLeft, rl.KeyDown, rl.KeyUp, rl.KeyRight}
	if slices.Contains(arrowKeys, keyPressed) {
		code = motion.KeyArrow
		keyCharPressed = keyPressed
	} else if keyPressed == rl.KeyEscape || keyPressed == rl.KeyCapsLock {
		code = motion.KeyEsc
		keyCharPressed = keyPressed
	}

	switch {
	case rl.IsKeyDown(rl.KeyLeftControl):
		modifiers = motion.ModCtrl
		keyCharPressed = keyPressed
	}

	if keyCharPressed == 0 {
		return
	}

	keyRune = rune(keyCharPressed)
	key := motion.Key{Code: code, Rune: keyRune, Modifiers: modifiers}

	slog.Debug("Handling key input", slog.Any("key", key))
	mode.HandleKey(ctx, key)
}

func (SpreadsheetCursorStateHandler) HandleInput(appAssets *assets.Assets, dg *database.DataGrid, eg *EditorGrid, c *Cursor, connManager *database.ConnectionManager) {
	const cellHeight int32 = 30 // @TODO: Pass configuration with cellHeight
	var keyPressed int32 = rl.GetCharPressed()

	var mouse rl.Vector2 = rl.GetMousePosition()
	var mouseWheelStep int32 = 1
	// Only scroll if mouse inside the zone
	if rl.CheckCollisionPointRec(mouse, c.Zone.Bounds) {
		if rl.IsKeyDown(rl.KeyLeftShift) {
			// Mouse wheel scroll (horizontal)
			c.Position.Col -= int32(rl.GetMouseWheelMove()) * int32(mouseWheelStep)
		} else {
			// Mouse wheel scroll (vertical)
			c.Position.Row -= int32(rl.GetMouseWheelMove()) * mouseWheelStep
		}
	}

	switch c.Common.Mode {
	case cursor.ModeVLine:
		fallthrough
	case cursor.ModeVBlock:
		fallthrough
	case cursor.ModeVisual:
		fallthrough
	case cursor.ModeNormal:
		if rl.IsKeyDown(rl.KeyLeftShift) {
			switch {
			case rl.IsKeyPressed(rl.KeyHome):
				c.Position.Col = 0
			case rl.IsKeyPressed(rl.KeyEnd):
				c.Position.Col = dg.Cols - 1
			case slices.Contains(HANDLED_MOTION_KEY_CODES, int(keyPressed)):
				c.AppendMotion(rune(keyPressed))
			}
		} else if rl.IsKeyDown(rl.KeyLeftControl) {
			switch {
			case rl.IsKeyPressed(rl.KeyC):
				if c.Position.Col >= 0 && dg.Cols > 0 {
					// @TODO: Add type specific formatting
					var dataString string = ""
					switch c.Common.Mode {
					case cursor.ModeVisual:
						for row := c.Position.SelectStartRow; row <= c.Position.SelectEndRow; row++ {
							for col := int32(0); col < dg.Cols; col++ {
								if c.IsSelected(col, row) {
									dataString += format.GetValueAsString(dg.Data[row][dg.Headers[col]]) + ","
								}
							}
							dataString += "\n"
						}
					case cursor.ModeVLine:
						fallthrough
					case cursor.ModeVBlock:
						for row := c.Position.SelectStartRow; row <= c.Position.SelectEndRow; row++ {
							for col := c.Position.SelectStartCol; col < c.Position.SelectEndCol; col++ {
								dataString += format.GetValueAsString(dg.Data[row][dg.Headers[col]]) + ","
							}
							dataString += format.GetValueAsString(dg.Data[row][dg.Headers[c.Position.SelectEndCol]]) + "\n"
						}
					case cursor.ModeNormal:
						dataString = format.GetValueAsString(dg.Data[c.Position.Row][dg.Headers[c.Position.Col]])
					}
					slog.Debug("Copied to clipboard from spreadsheet", slog.String("dataString", dataString))
					clipboard.Write(clipboard.FmtText, []byte(dataString))
				}
			case rl.IsKeyPressed(rl.KeyE):
				CursorSpreadsheet.SetActive(false)
				CursorConnection.SetActive(true)
				CurrCursor = CursorConnection
			case rl.IsKeyPressed(rl.KeyV):
				c.AppendMotionString("^V")
			case rl.IsKeyPressed(rl.KeyW):
				c.AppendMotionString("^W")
			}
		} else {
			var pageRows int8 = c.Zone.GetNumberOfVisibleRows(int32(cellHeight))
			switch {
			case rl.IsKeyPressed(rl.KeyEscape) || rl.IsKeyPressed(rl.KeyCapsLock): // @TODO: Remove personal preference CapsLock
				c.Handler.Reset(c)
			case slices.Contains(HANDLED_MOTION_KEY_CODES, int(keyPressed)):
				c.AppendMotion(rune(keyPressed))
			case rl.IsKeyPressed(rl.KeySlash):
				c.Common.Mode = cursor.ModeCommand
				c.Common.CmdBuf = "/"
				c.Common.MotionBuf = "/"
				c.UpdateCmdLine()
			case rl.IsKeyPressed(rl.KeyHome):
				c.Position.Row = 0
			case rl.IsKeyPressed(rl.KeyEnd):
				c.Position.Row = dg.Rows - 1
			case rl.IsKeyPressed(rl.KeyPageUp):
				c.Position.Row -= int32(pageRows)
			case rl.IsKeyPressed(rl.KeyPageDown):
				c.Position.Row += int32(pageRows)
			case rl.IsKeyPressed(rl.KeyDown):
				c.Position.Row++
			case rl.IsKeyPressed(rl.KeyUp):
				c.Position.Row--
			case rl.IsKeyPressed(rl.KeyLeft):
				c.Position.Col--
			case rl.IsKeyPressed(rl.KeyRight):
				c.Position.Col++
			}
		}
	case cursor.ModeCommand:
		switch {
		case rl.IsKeyPressed(rl.KeyEnter):
			c.ExecuteCommand()
		case rl.IsKeyPressed(rl.KeyBackspace):
			c.Common.CmdBuf = c.Common.CmdBuf[:len(c.Common.CmdBuf)-1]
			c.UpdateCmdLine()
			if len(c.Common.CmdBuf) <= 0 {
				c.Handler.Reset(c)
			}
		case rl.IsKeyPressed(rl.KeyEscape) || rl.IsKeyPressed(rl.KeyCapsLock): // @TODO: Remove personal preference CapsLock
			c.Handler.Reset(c)
		case keyPressed >= 32 && keyPressed <= 125:
			c.Common.CmdBuf += string(rune(keyPressed))
			c.UpdateCmdLine()
		}
	}

	c.Zone.ClampScrollsToZoneSize()
	c.ClampFocus(dg.Cols-1, dg.Rows-1)
	c.UpdateSelectBasedOnPosition()
}

func (ConnectionsCursorStateHandler) HandleInput(appAssets *assets.Assets, dg *database.DataGrid, eg *EditorGrid, c *Cursor, connManager *database.ConnectionManager) {
	var keyPressed int32 = rl.GetCharPressed()

	switch c.Common.Mode {
	case cursor.ModeNormal:
		if rl.IsKeyDown(rl.KeyLeftShift) {
			switch {
			case slices.Contains(HANDLED_MOTION_KEY_CODES, int(keyPressed)):
				c.AppendMotion(rune(keyPressed))
			case rl.IsKeyPressed(rl.KeyHome):
				c.Position.Col = 0
			case rl.IsKeyPressed(rl.KeyEnd):
				c.Position.Col = dg.Cols - 1
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
				err := connManager.SetCurrentConnectionByName(config.Get().Connections[CursorConnection.Position.Row].Name)
				if err != nil {
					slog.Error("Failed to set current connection by name", slog.Any("error", err))
				}
				// @TODO: Go back to last focused c
				CursorConnection.SetActive(false)
				CursorEditor.SetActive(true)
				CurrCursor = CursorEditor
			case rl.IsKeyPressed(rl.KeyEscape) || rl.IsKeyPressed(rl.KeyCapsLock): // @TODO: Remove personal preference CapsLock
				// @TODO: Go back to last focused c
				CursorConnection.SetActive(false)
				CursorEditor.SetActive(true)
				CurrCursor = CursorEditor
				CurrCursor = CursorSpreadsheet
			case slices.Contains(HANDLED_MOTION_KEY_CODES, int(keyPressed)):
				c.AppendMotion(rune(keyPressed)) // @TODO: Disable visual mode for connection selection
			case rl.IsKeyPressed(rl.KeySlash):
				c.Common.Mode = cursor.ModeCommand
				c.Common.CmdBuf = "/"
				c.Common.MotionBuf = "/"
				c.UpdateCmdLine()
			case rl.IsKeyPressed(rl.KeyHome):
				c.Position.Row = 0
			case rl.IsKeyPressed(rl.KeyEnd):
				c.Position.Row = dg.Rows - 1
			case rl.IsKeyPressed(rl.KeyDown):
				c.Position.Row++
			case rl.IsKeyPressed(rl.KeyUp):
				c.Position.Row--
			case rl.IsKeyPressed(rl.KeyLeft):
				c.Position.Col--
			case rl.IsKeyPressed(rl.KeyRight):
				c.Position.Col++
			}
		}
	case cursor.ModeCommand:
		switch {
		case rl.IsKeyPressed(rl.KeyEnter):
			c.ExecuteCommand()
		case rl.IsKeyPressed(rl.KeyBackspace):
			c.Common.CmdBuf = c.Common.CmdBuf[:len(c.Common.CmdBuf)-1]
			c.UpdateCmdLine()
			if len(c.Common.CmdBuf) <= 0 {
				c.Handler.Reset(c)
			}
		case rl.IsKeyPressed(rl.KeyEscape) || rl.IsKeyPressed(rl.KeyCapsLock): // @TODO: Remove personal preference CapsLock
			c.Handler.Reset(c)
		case keyPressed >= 32 && keyPressed <= 125:
			c.Common.CmdBuf += string(rune(keyPressed))
			c.UpdateCmdLine()
		}
	}

	c.Zone.ClampScrollsToZoneSize()
	c.ClampFocus(0, c.Position.MaxRow)
	c.UpdateSelectBasedOnPosition()
}

func (EditorCursorStateHandler) HandleInput(appAssets *assets.Assets, dg *database.DataGrid, eg *EditorGrid, c *Cursor, connManager *database.ConnectionManager) {
	const cellHeight int32 = 30 // @TODO: Pass configuration with cellHeight
	var keyPressed int32 = rl.GetCharPressed()

	var mouse rl.Vector2 = rl.GetMousePosition()
	var mouseWheelStep int32 = 1
	// Only scroll if mouse inside the zone
	if rl.CheckCollisionPointRec(mouse, c.Zone.Bounds) {
		if rl.IsKeyDown(rl.KeyLeftShift) {
			// Mouse wheel scroll (horizontal)
			c.Position.Col -= int32(rl.GetMouseWheelMove()) * int32(mouseWheelStep)
		} else {
			// Mouse wheel scroll (vertical)
			c.Position.Row -= int32(rl.GetMouseWheelMove()) * mouseWheelStep
		}
	}

	switch c.Common.Mode {
	case cursor.ModeVLine:
		fallthrough
	case cursor.ModeVBlock:
		fallthrough
	case cursor.ModeVisual:
		fallthrough
	case cursor.ModeNormal:
		if rl.IsKeyDown(rl.KeyLeftShift) {
			switch {
			case slices.Contains(HANDLED_MOTION_KEY_CODES, int(keyPressed)):
				c.AppendMotion(rune(keyPressed))
			}
		} else if rl.IsKeyDown(rl.KeyLeftControl) {
			switch {
			case rl.IsKeyPressed(rl.KeyEnter):
				detectAndExecuteQuery(c, eg, connManager)
				// @TODO: Remove hard coded query and detect query in editor
				// query := "SELECT 1"
				// err := connManager.ExecuteQuery(context.Background(), connManager.GetCurrentConnectionName(), query)
				//if err != nil {
				//	slog.Error("Failed to execute query", slog.Any("error", err))
				//	c.Common.Logs.Channel <- "Failed to execute query (Something went wrong)"
				//}
			case rl.IsKeyPressed(rl.KeyC):
				if c.Position.Col >= 0 && eg.Rows > 0 {
					var dataString string = ""
					switch c.Common.Mode {
					case cursor.ModeVisual:
						for row := c.Position.SelectStartRow; row <= c.Position.SelectEndRow; row++ {
							if eg.Cols[row] > 0 {
								for col := int32(0); col < eg.Cols[row]; col++ {
									if c.IsSelected(col, row) {
										dataString += string(eg.Text[row][col])
									}
								}
							}
							dataString += "\n"
						}
					case cursor.ModeVLine:
						for row := c.Position.SelectStartRow; row <= c.Position.SelectEndRow; row++ {
							dataString += eg.Text[row] + "\n"
						}
					case cursor.ModeVBlock:
						for row := c.Position.SelectStartRow; row <= c.Position.SelectEndRow; row++ {
							if eg.Cols[row] > 0 {
								endCol := min(c.Position.SelectEndCol, eg.Cols[row])
								for col := c.Position.SelectStartCol; col <= endCol; col++ {
									dataString += string(eg.Text[row][col])
								}
							}
							dataString += "\n"
						}
					default:
						dataString = eg.Text[c.Position.Row]
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
				c.AppendMotionString("^V")
			case rl.IsKeyPressed(rl.KeyW):
				c.AppendMotionString("^W")
			}
		} else {
			var pageRows int8 = c.Zone.GetNumberOfVisibleRows(int32(cellHeight))
			switch {
			case rl.IsKeyPressed(rl.KeyEscape) || rl.IsKeyPressed(rl.KeyCapsLock): // @TODO: Remove personal preference CapsLock
				c.Handler.Reset(c)
			case slices.Contains(HANDLED_MOTION_KEY_CODES, int(keyPressed)):
				c.AppendMotion(rune(keyPressed))
			case rl.IsKeyPressed(rl.KeySlash):
				c.Common.Mode = cursor.ModeCommand
				c.Common.CmdBuf = "/"
				c.Common.MotionBuf = "/"
				c.UpdateCmdLine()
			case rl.IsKeyPressed(rl.KeyHome):
				c.Position.Col = 0
			case rl.IsKeyPressed(rl.KeyEnd):
				c.Position.Col = eg.Cols[c.Position.Row]
			case rl.IsKeyPressed(rl.KeyPageUp):
				c.Position.Row -= int32(pageRows)
			case rl.IsKeyPressed(rl.KeyPageDown):
				c.Position.Row += int32(pageRows)
			case rl.IsKeyPressed(rl.KeyDown):
				c.Position.Row++
			case rl.IsKeyPressed(rl.KeyUp):
				c.Position.Row--
			case rl.IsKeyPressed(rl.KeyLeft):
				c.Position.Col--
			case rl.IsKeyPressed(rl.KeyRight):
				c.Position.Col++
			}
		}
	case cursor.ModeCommand:
		switch {
		case rl.IsKeyPressed(rl.KeyEnter):
			c.ExecuteCommand()
		case rl.IsKeyPressed(rl.KeyBackspace):
			c.Common.CmdBuf = c.Common.CmdBuf[:len(c.Common.CmdBuf)-1]
			c.UpdateCmdLine()
			if len(c.Common.CmdBuf) <= 0 {
				c.Handler.Reset(c)
			}
		case rl.IsKeyPressed(rl.KeyEscape) || rl.IsKeyPressed(rl.KeyCapsLock): // @TODO: Remove personal preference CapsLock
			c.Handler.Reset(c)
		case keyPressed >= 32 && keyPressed <= 125:
			c.Common.CmdBuf += string(rune(keyPressed))
			c.UpdateCmdLine()
		}
	}

	if eg.Rows <= 0 {
		c.Position.Row = 0
		c.Position.Col = 0
	} else {
		c.Zone.ClampScrollsToZoneSize()
		// Additional focus check to be in scope of Cols
		if c.Position.Row >= eg.Rows {
			c.Position.Row = eg.Rows - 1
		} else if c.Position.Row < 0 {
			c.Position.Row = 0
		}
		c.ClampFocus(eg.Cols[c.Position.Row]-1, eg.Rows)
		c.UpdateSelectBasedOnPosition()
	}
}

func detectAndExecuteQuery(c *Cursor, eg *EditorGrid, connManager *database.ConnectionManager) {
	query := ""
	if eg.Rows <= 0 {
		slog.Error("Failed to execute query", slog.String("error", "No query provided"))
		c.Common.Logs.Channel <- "Failed to execute query (No query provided)"
		return
	}
	// @TODO: Implement other modes behavior
	switch c.Common.Mode {
	case cursor.ModeVisual:
		for row := c.Position.SelectStartRow; row <= c.Position.SelectEndRow; row++ {
			if eg.Cols[row] > 0 {
				for col := int32(0); col < eg.Cols[row]; col++ {
					if c.IsSelected(col, row) {
						query += string(eg.Text[row][col])
					}
				}
			}
		}
	case cursor.ModeVLine:
		for row := c.Position.SelectStartRow; row <= c.Position.SelectEndRow; row++ {
			if eg.Cols[row] > 0 {
				for col := int32(0); col < eg.Cols[row]; col++ {
					query += string(eg.Text[row][col])
				}
			}
		}
	case cursor.ModeInsert:
		fallthrough
	case cursor.ModeNormal:
		start, end := eg.DetectQueryRowsBoundaryBasedOnRow(c.Position.Row)
		var sb strings.Builder
		for i := start; i <= end; i++ {
			sb.WriteString(eg.Text[i])
			if i < end {
				sb.WriteString(" ")
			}
		}

		query = strings.TrimSpace(sb.String())
	}
	if query == "" {
		slog.Error("Failed to execute query", slog.String("error", "No query provided"))
		c.Common.Logs.Channel <- "Failed to execute query (No query provided)"
		return
	}

	err := connManager.ExecuteQuery(context.Background(), connManager.GetCurrentConnectionName(), query)
	if err != nil {
		slog.Error("Failed to execute query", slog.Any("error", err))
		c.Common.Logs.Channel <- fmt.Sprintf("Failed to execute query (%s)", err)
	}
}
