package display

import (
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/quar15/qq-go/internal/assets"
	"github.com/quar15/qq-go/internal/colors"
	"github.com/quar15/qq-go/internal/database"
	"github.com/quar15/qq-go/internal/format"
)

func (z *Zone) DrawSpreadsheetZone(appAssets *assets.Assets, dg *database.DataGrid, cursor *Cursor) {
	// @TODO: Consider approach with RenderTexture2D for performance
	const cellHeight int = 30
	const textPadding int32 = 6
	var mouse rl.Vector2 = rl.GetMousePosition()
	var (
		counterColumnCharactersCount int = format.CountDigits(int(dg.Rows))
		counterColumnWidth           int = int(appAssets.MainFontCharacterWidth)*counterColumnCharactersCount + int(textPadding*2)
		contentWidth                 int = counterColumnWidth
		contentHeight                int = cellHeight * (int(dg.Rows) + 2)
	)
	for col := int32(0); col < dg.Cols; col++ {
		contentWidth += int(dg.ColumnsWidth[col])
	}

	const linesPadding int8 = 2
	scrollRow, lastRowToRender := updateSpreadsheetScrollBasedOnCursor(z, dg, cursor, cellHeight, linesPadding)

	rl.BeginScissorMode(int32(z.Bounds.X), int32(z.Bounds.Y), int32(z.Bounds.Width), int32(z.Bounds.Height))

	for row := scrollRow; row < lastRowToRender; row++ {
		renderContentRow(z, appAssets, dg, cursor, counterColumnWidth, cellHeight, textPadding, mouse, row)
	}

	for row := scrollRow; row < lastRowToRender; row++ {
		renderSpreadsheetCounterColumnRow(z, appAssets, counterColumnWidth, counterColumnCharactersCount, cellHeight, textPadding, mouse, row)
	}

	rl.EndScissorMode()

	// Draw static header
	rl.DrawRectangle(int32(z.Bounds.X), int32(z.Bounds.Y), int32(z.Bounds.Width), int32(cellHeight), colors.Surface0()) // Left upper corner fill
	renderSpreadsheetHeadersRow(z, appAssets, dg, counterColumnWidth, cellHeight, textPadding, mouse)

	z.ContentSize.Y = max(float32(contentHeight), z.Bounds.Height)
	z.ContentSize.X = max(float32(contentWidth), z.Bounds.Width)
	z.drawScrollbars()
}

func renderContentRow(z *Zone, appAssets *assets.Assets, dg *database.DataGrid, cursor *Cursor, counterColumnWidth int, cellHeight int, textPadding int32, mouse rl.Vector2, row int32) {
	for col, key := range dg.Headers {
		// @TODO: Consider limiting draw to only visible columns
		val := dg.Data[row][key]
		var cellX int32 = int32(z.Bounds.X-z.Scroll.X) + int32(counterColumnWidth)
		for i := range col {
			cellX += dg.ColumnsWidth[i]
		}
		var cellY int32 = int32(z.Bounds.Y) + (row+1)*int32(cellHeight) - int32(z.Scroll.Y)

		var cellBackgroundColor rl.Color = colors.Background()
		var cellBorderColor rl.Color = colors.Mantle()
		if cursor.IsActive() && cursor.IsFocused(int32(col), row) {
			cellBackgroundColor = colors.Mantle()
			cellBorderColor = colors.Blue()
		}
		if cursor.IsActive() && cursor.IsSelected(int32(col), row) {
			cellBackgroundColor = colors.Surface1()
		}
		var cellRect rl.RectangleInt32 = rl.RectangleInt32{X: cellX, Y: cellY, Width: dg.ColumnsWidth[col], Height: int32(cellHeight)}
		rl.DrawRectangleRec(cellRect.ToFloat32(), cellBackgroundColor)
		if rl.CheckCollisionPointRec(mouse, cellRect.ToFloat32()) {
			if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
				cursor.Position.Col = int32(col)
				cursor.Position.Row = row
			}
		}
		cellText := format.GetValueAsString(val)
		var cellTextSliceLimit int = len(cellText)
		var maxNumberOfCharacters int = int(dg.ColumnsWidth[col] / int32(appAssets.MainFontCharacterWidth))
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
				Width:  float32(dg.ColumnsWidth[col]) + 1,
				Height: float32(cellHeight) + 1,
			},
			2,
			cellBorderColor,
		)
	}
}

func renderSpreadsheetCounterColumnRow(z *Zone, appAssets *assets.Assets, counterColumnWidth int, counterColumnCharactersCount int, cellHeight int, textPadding int32, mouse rl.Vector2, row int32) {
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

func renderSpreadsheetHeadersRow(z *Zone, appAssets *assets.Assets, dg *database.DataGrid, counterColumnWidth int, cellHeight int, textPadding int32, mouse rl.Vector2) {
	for col := int32(0); col < dg.Cols; col++ {
		var cellX int32 = int32(z.Bounds.X-z.Scroll.X) + int32(counterColumnWidth)
		for c := int32(0); c < col; c++ {
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
}

func updateSpreadsheetScrollBasedOnCursor(z *Zone, dg *database.DataGrid, cursor *Cursor, cellHeight int, linesPadding int8) (scrollRow int32, lastRowToRender int32) {
	z.Scroll.X = 0
	for col := int32(0); col < cursor.Position.Col; col++ {
		z.Scroll.X += float32(dg.ColumnsWidth[col])
	}

	var rowsToRender int8 = z.GetNumberOfVisibleRows(int32(cellHeight)) + 1
	scrollRow = int32(z.Scroll.Y) / int32(cellHeight)
	if rowsToRender > linesPadding*2 {
		if cursor.Position.Row < scrollRow+int32(linesPadding) {
			scrollRow = max(cursor.Position.Row-int32(linesPadding), 0)
		}

		if cursor.Position.Row >= scrollRow+int32(rowsToRender)-int32(linesPadding) {
			scrollRow = cursor.Position.Row - int32(rowsToRender-linesPadding-1)
			if cursor.Position.Row >= cursor.Position.MaxRow-int32(linesPadding) {
				scrollRow -= int32(linesPadding) - (cursor.Position.MaxRow - cursor.Position.Row + 1)
			}
		}

		z.Scroll.Y = float32(cellHeight * int(scrollRow))
	} else {
		z.Scroll.Y = float32(cellHeight * int(cursor.Position.Row))
	}
	scrollRow = min(max(scrollRow, 0), cursor.Position.MaxRow)
	z.ClampScrollsToZoneSize()
	lastRowToRender = min(dg.Rows, scrollRow+int32(rowsToRender))

	return scrollRow, lastRowToRender
}

func (z *Zone) GetNumberOfVisibleRows(cellHeight int32) int8 {
	var availableContentSpaceHeight int32 = int32(z.Bounds.Height) - int32(cellHeight)
	return int8(availableContentSpaceHeight / cellHeight)
}
