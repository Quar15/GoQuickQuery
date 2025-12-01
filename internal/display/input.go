package display

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/quar15/qq-go/internal/database"
	"github.com/quar15/qq-go/internal/utilities"
	"golang.design/x/clipboard"
)

func HandleSpreadsheetInput(z *Zone, dg *database.DataGrid, cursor *SpreadsheetCursor, cellHeight int32) {
	var keyPressed int32 = rl.GetCharPressed()
	switch cursor.Mode {
	case ModeVLine:
		fallthrough
	case ModeVisual:
		fallthrough
	case ModeNormal:
		if rl.IsKeyDown(rl.KeyLeftShift) {
			switch {
			case rl.IsKeyPressed(rl.KeyHome):
				z.Scroll.X = 0
				cursor.Col = 0
			case rl.IsKeyPressed(rl.KeyEnd):
				z.Scroll.X = z.ContentSize.X - z.Bounds.Width
				cursor.Col = dg.Cols - 1
			case rl.IsKeyPressed(rl.KeyV):
				cursor.SetSelect(0, cursor.Row, dg.Cols-1, cursor.Row)
				cursor.TransitionMode(ModeVLine)
			case rl.IsKeyPressed(rl.KeyG):
				cursor.Row = dg.Rows - 1
				z.Scroll.Y = z.ContentSize.Y - z.Bounds.Height
			}
		} else if rl.IsKeyDown(rl.KeyLeftControl) {
			switch {
			case rl.IsKeyDown(rl.KeyC):
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
		} else {
			var pageRows int8 = z.GetNumberOfVisibleRows(int32(cellHeight))
			switch {
			case rl.IsKeyPressed(rl.KeyEscape) || rl.IsKeyPressed(rl.KeyCapsLock): // @TODO: Remove personal preference CapsLock
				cursor.Reset()
			case rl.IsKeyPressed(rl.KeyHome):
				cursor.Row = 0
				z.Scroll.Y = 0
			case rl.IsKeyPressed(rl.KeyEnd):
				cursor.Row = dg.Rows - 1
				z.Scroll.Y = z.ContentSize.Y - z.Bounds.Height
			case rl.IsKeyPressed(rl.KeyPageUp):
				cursor.Row -= int32(pageRows)
				z.Scroll.Y -= float32(cellHeight * int32(pageRows))
			case rl.IsKeyPressed(rl.KeyPageDown):
				cursor.Row += int32(pageRows)
				z.Scroll.Y += float32(cellHeight * int32(pageRows))
			case rl.IsKeyPressed(rl.KeyJ) || rl.IsKeyPressed(rl.KeyDown):
				z.Scroll.Y += float32(cellHeight)
				cursor.Row++
			case rl.IsKeyPressed(rl.KeyK) || rl.IsKeyPressed(rl.KeyUp):
				z.Scroll.Y -= float32(cellHeight)
				cursor.Row--
			case rl.IsKeyPressed(rl.KeyH) || rl.IsKeyPressed(rl.KeyLeft):
				z.Scroll.X -= float32(dg.ColumnsWidth[cursor.Col])
				cursor.Col--
			case rl.IsKeyPressed(rl.KeyL) || rl.IsKeyPressed(rl.KeyRight):
				z.Scroll.X += float32(dg.ColumnsWidth[cursor.Col])
				cursor.Col++
			case rl.IsKeyPressed(rl.KeyV):
				cursor.SetSelect(cursor.Col, cursor.Row, cursor.Col, cursor.Row)
			case rl.IsKeyPressed(rl.KeySlash):
				cursor.Mode = ModeCommand
				cursor.CmdBuf = "/"
				cursor.MotionBuf = "/"
				cursor.UpdateCmdLine()
			case rl.IsKeyPressed(rl.KeyG):
				cursor.MotionBuf += "g"
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
