package display

import (
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/quar15/qq-go/internal/assets"
	"github.com/quar15/qq-go/internal/colors"
	"github.com/quar15/qq-go/internal/database"
	"github.com/quar15/qq-go/internal/format"
)

func (z *Zone) GetNumberOfVisibleRows(cellHeight int32) int8 {
	var availableContentSpaceHeight int32 = int32(z.Bounds.Height) - int32(cellHeight)
	return int8(availableContentSpaceHeight / cellHeight)
}

func (z *Zone) DrawSpreadsheetZone(appAssets *assets.Assets, dg *database.DataGrid, cursor *Cursor) {
	// @TODO: Consider approach with RenderTexture2D for performance
	var mouse rl.Vector2 = rl.GetMousePosition()

	const cellHeight int = 30
	const textPadding int32 = 6

	var contentWidth int = 0
	var contentHeight int = 0
	var counterColumnCharactersCount int = format.CountDigits(int(dg.Rows))
	var counterColumnWidth int = int(appAssets.MainFontCharacterWidth)*counterColumnCharactersCount + int(textPadding*2)
	contentWidth += counterColumnWidth
	for col := int8(0); col < dg.Cols; col++ {
		contentWidth += int(dg.ColumnsWidth[col])
	}

	// Update scroll based on cursor
	z.Scroll.Y = float32(cellHeight * int(cursor.Position.Row))
	z.Scroll.X = 0
	for col := int8(0); col < cursor.Position.Col; col++ {
		z.Scroll.X += float32(dg.ColumnsWidth[col])
	}
	z.ClampScrollsToZoneSize()

	var rowsToRender int8 = z.GetNumberOfVisibleRows(int32(cellHeight)) + 1
	const linesPadding int8 = 2
	var firstVisibleRowToScrollIndex int32 = min(max(cursor.Position.Row, 0), dg.Rows) // Swapping screens creates weird behavior
	var lastRowToRender = min(dg.Rows, firstVisibleRowToScrollIndex+int32(rowsToRender))
	if rowsToRender > 4 {
		// @TODO: Make window for scrolling
		topPaddingIndex := max(firstVisibleRowToScrollIndex-2, 0)
		firstVisibleRowToScrollIndex = topPaddingIndex
		z.Scroll.Y = float32(cellHeight * int(firstVisibleRowToScrollIndex))
	}

	rl.BeginScissorMode(int32(z.Bounds.X), int32(z.Bounds.Y), int32(z.Bounds.Width), int32(z.Bounds.Height))

	for row := firstVisibleRowToScrollIndex; row < lastRowToRender; row++ {
		for c, key := range dg.Headers {
			// @TODO: Consider limiting draw to only visible columns
			val := dg.Data[row][key]
			var cellX int32 = int32(z.Bounds.X-z.Scroll.X) + int32(counterColumnWidth)
			for i := range c {
				cellX += dg.ColumnsWidth[i]
			}
			var cellY int32 = int32(z.Bounds.Y) + (row+1)*int32(cellHeight) - int32(z.Scroll.Y)

			var cellBackgroundColor rl.Color = colors.Background()
			var cellBorderColor rl.Color = colors.Mantle()
			if cursor.IsFocused(int8(c), row) {
				cellBackgroundColor = colors.Mantle()
				cellBorderColor = colors.Blue()
			}
			if cursor.IsSelected(int8(c), row) {
				cellBackgroundColor = colors.Surface1()
			}
			var cellRect rl.RectangleInt32 = rl.RectangleInt32{X: cellX, Y: cellY, Width: dg.ColumnsWidth[c], Height: int32(cellHeight)}
			rl.DrawRectangleRec(cellRect.ToFloat32(), cellBackgroundColor)
			if rl.CheckCollisionPointRec(mouse, cellRect.ToFloat32()) {
				if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
					cursor.Position.Col = int8(c)
					cursor.Position.Row = row
				}
			}
			cellText := format.GetValueAsString(val)
			var cellTextSliceLimit int = len(cellText)
			var maxNumberOfCharacters int = int(dg.ColumnsWidth[c] / int32(appAssets.MainFontCharacterWidth))
			if cellTextSliceLimit > maxNumberOfCharacters {
				cellTextSliceLimit = maxNumberOfCharacters
			}
			appAssets.DrawTextMainFont(
				cellText[:cellTextSliceLimit],
				rl.Vector2{X: float32(cellX + textPadding), Y: float32(cellY + textPadding)},
				colors.Text(),
			)
			rl.DrawRectangleLinesEx(
				rl.Rectangle{
					X:      float32(cellX),
					Y:      float32(cellY),
					Width:  float32(dg.ColumnsWidth[c]) + 1,
					Height: float32(cellHeight) + 1,
				},
				2,
				cellBorderColor,
			)
		}
	}
	contentHeight += cellHeight * int(dg.Rows)

	// Draw row counter
	for row := firstVisibleRowToScrollIndex; row < lastRowToRender; row++ {
		var cellX int32 = int32(z.Bounds.X)
		var cellY int32 = int32(z.Bounds.Y) + (row+1)*int32(cellHeight) - int32(z.Scroll.Y)

		var bg rl.Color = colors.Surface0()
		if z.MouseInside(mouse) && mouse.Y > float32(cellY) && mouse.Y < float32(cellY)+float32(cellHeight) {
			bg = colors.Mantle()
		}
		var counterColumnLeftPadding float32 = float32(textPadding) + float32(counterColumnCharactersCount-format.CountDigits(int(row)+1))*appAssets.MainFontCharacterWidth
		rl.DrawRectangle(cellX, cellY, int32(counterColumnWidth), int32(cellHeight), bg)
		rl.DrawLineEx(rl.Vector2{X: float32(cellX), Y: float32(cellY)}, rl.Vector2{X: float32(cellX + int32(counterColumnWidth)), Y: float32(cellY)}, 2, colors.Surface1())
		appAssets.DrawTextMainFont(strconv.Itoa(int(row+1)), rl.Vector2{X: float32(cellX) + counterColumnLeftPadding, Y: float32(cellY + textPadding)}, colors.Overlay0())
	}

	rl.EndScissorMode()

	rl.DrawRectangle(int32(z.Bounds.X), int32(z.Bounds.Y), int32(z.Bounds.Width), int32(cellHeight), colors.Surface0())

	// Draw header
	for col := int8(0); col < dg.Cols; col++ {
		var cellX int32 = int32(z.Bounds.X-z.Scroll.X) + int32(counterColumnWidth)
		for c := int8(0); c < col; c++ {
			cellX += dg.ColumnsWidth[c]
		}
		var cellY int32 = int32(z.Bounds.Y)

		var bg rl.Color = colors.Surface0()
		if z.MouseInside(mouse) && mouse.X > float32(cellX) && mouse.X < float32(cellX)+float32(dg.ColumnsWidth[col]) {
			bg = colors.Mantle()
		}
		rl.DrawRectangle(cellX, cellY, dg.ColumnsWidth[col], int32(cellHeight), bg)
		rl.DrawLineEx(rl.Vector2{X: float32(cellX), Y: float32(cellY)}, rl.Vector2{X: float32(cellX), Y: float32(cellY + int32(cellHeight))}, 2, colors.Surface1())
		appAssets.DrawTextMainFont(dg.Headers[col], rl.Vector2{X: float32(cellX + textPadding), Y: float32(cellY + textPadding)}, colors.Text())
	}
	contentHeight += cellHeight * 2

	z.ContentSize.Y = max(float32(contentHeight), z.Bounds.Height)
	z.ContentSize.X = max(float32(contentWidth), z.Bounds.Width)
	z.drawScrollbars()
}
