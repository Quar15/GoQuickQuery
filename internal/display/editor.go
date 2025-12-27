package display

import (
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/quar15/qq-go/internal/assets"
	"github.com/quar15/qq-go/internal/config"
	"github.com/quar15/qq-go/internal/cursor"
	"github.com/quar15/qq-go/internal/editor"
	"github.com/quar15/qq-go/internal/format"
)

func (z *Zone) DrawEditor(appAssets *assets.Assets, eg *editor.Grid, cursor *cursor.Cursor, shouldDrawCursor bool) {
	if eg.Rows <= 0 {
		return
	}
	const cellHeight int = 28
	const textPadding int32 = 6
	const rowsInitialPadding int32 = 16

	var (
		contentHeight                int = cellHeight*int(eg.Rows+1) + int(rowsInitialPadding)
		counterColumnCharactersCount int = format.CountDigits(int(eg.Rows))
		counterColumnWidth           int = int(appAssets.MainFontCharacterWidth)*counterColumnCharactersCount + int(textPadding*2)
	)
	const linesPadding int8 = 4
	scrollRow, lastRowToRender := updateEditorScrollBasedOnCursor(z, cursor, cellHeight, linesPadding)

	rl.BeginScissorMode(int32(z.Bounds.X), int32(z.Bounds.Y), int32(z.Bounds.Width), int32(z.Bounds.Height))

	for row := scrollRow; row < lastRowToRender; row++ {
		renderEditorTextRow(z, appAssets, eg, cursor, counterColumnWidth, cellHeight, rowsInitialPadding, textPadding, row)
	}
	renderEditorDetectedQueryOutline(z, appAssets, eg, cursor, counterColumnWidth, cellHeight, rowsInitialPadding, textPadding)
	for row := scrollRow; row < lastRowToRender; row++ {
		renderEditorRowCounter(z, appAssets, counterColumnCharactersCount, cellHeight, rowsInitialPadding, textPadding, row)
	}
	if shouldDrawCursor {
		renderEditorCursor(z, appAssets, eg, cursor, counterColumnWidth, cellHeight, rowsInitialPadding, textPadding)
	}

	rl.EndScissorMode()

	z.ContentSize.Y = max(float32(contentHeight), z.Bounds.Height)
	z.ContentSize.X = max(eg.MaxWidth, z.Bounds.Width)
	z.drawScrollbars()
}

func renderEditorTextRow(z *Zone, appAssets *assets.Assets, eg *editor.Grid, cursor *cursor.Cursor, counterColumnWidth int, cellHeight int, rowsInitialPadding int32, textPadding int32, row int32) {
	var cellY float32 = z.Bounds.Y + float32(row*int32(cellHeight)) - z.Scroll.Y + float32(rowsInitialPadding)
	if eg.Highlight[row] != nil {
		for col := range eg.Text[row] {
			var cellX float32 = z.Bounds.X + float32(counterColumnWidth) + float32(textPadding) + float32(col*int(appAssets.MainFontCharacterWidth))
			if cursor.IsActive() && cursor.IsSelected(int32(col), row) {
				rl.DrawRectangleRec(
					rl.Rectangle{
						X:      cellX,
						Y:      cellY,
						Width:  appAssets.MainFontCharacterWidth,
						Height: appAssets.MainFontSize,
					},
					config.Get().Colors.Surface1(),
				)
			}
			rl.DrawTextEx(
				appAssets.MainFont,
				string(eg.Text[row][col]),
				rl.Vector2{X: float32(cellX), Y: float32(cellY)},
				appAssets.MainFontSize,
				appAssets.MainFontSpacing,
				eg.Highlight[row][col].Color(),
			)
		}
	}
}

func renderEditorDetectedQueryOutline(z *Zone, appAssets *assets.Assets, eg *editor.Grid, cursor *cursor.Cursor, counterColumnWidth int, cellHeight int, rowsInitialPadding int32, textPadding int32) {
	start, end := eg.DetectQueryRowsBoundaryBasedOnRow(cursor.Position.Row)
	outlineRect := rl.Rectangle{
		X:      z.Bounds.X + float32(counterColumnWidth),
		Y:      z.Bounds.Y + float32(start*int32(cellHeight)) - z.Scroll.Y + float32(rowsInitialPadding-textPadding),
		Height: 0,
		Width:  0,
	}
	outlineRect.Height = float32(cellHeight * int(end-start+1))
	maxCols := eg.Cols[start]
	for i := start; i <= end; i++ {
		if maxCols < eg.Cols[i] {
			maxCols = eg.Cols[i]
		}
	}
	outlineRect.Width = float32(maxCols+1) * appAssets.MainFontCharacterWidth

	rl.DrawRectangleLinesEx(outlineRect, 2, config.Get().Colors.Accent())
}

func renderEditorRowCounter(z *Zone, appAssets *assets.Assets, counterColumnCharactersCount int, cellHeight int, rowsInitialPadding int32, textPadding int32, row int32) {
	var counterColumnLeftPadding float32 = float32(textPadding) + float32(counterColumnCharactersCount-format.CountDigits(int(row)+1))*appAssets.MainFontCharacterWidth
	var cellX float32 = z.Bounds.X + counterColumnLeftPadding
	var cellY float32 = z.Bounds.Y + float32(row*int32(cellHeight)) - z.Scroll.Y + float32(rowsInitialPadding)
	appAssets.DrawTextMainFont(strconv.Itoa(int(row+1)), rl.Vector2{X: cellX, Y: cellY}, config.Get().Colors.Overlay0())
}

func renderEditorCursor(z *Zone, appAssets *assets.Assets, eg *editor.Grid, cursor *cursor.Cursor, counterColumnWidth int, cellHeight int, rowsInitialPadding int32, textPadding int32) {
	var cellY float32 = z.Bounds.Y + float32(cursor.Position.Row*int32(cellHeight)) - z.Scroll.Y + float32(rowsInitialPadding)
	var cellX float32 = z.Bounds.X + float32(counterColumnWidth) + float32(textPadding)
	if eg.Cols[cursor.Position.Row] > 0 {
		cellX += float32(cursor.Position.Col) * appAssets.MainFontCharacterWidth
		rl.DrawRectangle(
			int32(cellX),
			int32(cellY),
			int32(appAssets.MainFontCharacterWidth),
			appAssets.MainFont.BaseSize,
			config.Get().Colors.Text(),
		)
		rl.DrawTextEx(
			appAssets.MainFont,
			string(eg.Text[cursor.Position.Row][cursor.Position.Col]),
			rl.Vector2{X: float32(cellX), Y: float32(cellY)},
			appAssets.MainFontSize,
			appAssets.MainFontSpacing,
			config.Get().Colors.Background(),
		)
	} else {
		// Empty row
		rl.DrawRectangle(
			int32(cellX),
			int32(cellY),
			int32(appAssets.MainFontCharacterWidth),
			appAssets.MainFont.BaseSize,
			config.Get().Colors.Text(),
		)
	}
}

func updateEditorScrollBasedOnCursor(z *Zone, cursor *cursor.Cursor, cellHeight int, linesPadding int8) (scrollRow int32, lastRowToRender int32) {
	z.Scroll.X = 0
	// @TODO: Horizontal scrolling with leeway

	var rowsToRender int8 = z.GetNumberOfVisibleRows(int32(cellHeight)) + 1
	scrollRow = int32(z.Scroll.Y) / int32(cellHeight)
	if rowsToRender > linesPadding*2 {
		if cursor.Position.Row < scrollRow+int32(linesPadding) {
			scrollRow = max(cursor.Position.Row-int32(linesPadding), 0)
		}

		if cursor.Position.Row >= scrollRow+int32(rowsToRender-linesPadding) {
			scrollRow = cursor.Position.Row - int32(rowsToRender-linesPadding-1)
			if cursor.Position.Row >= cursor.Position.MaxRow-int32(linesPadding) {
				scrollRow -= int32(linesPadding) - (cursor.Position.MaxRow - cursor.Position.Row)
			}
		}

		z.Scroll.Y = float32(cellHeight * int(scrollRow))
	} else {
		z.Scroll.Y = float32(cellHeight * int(cursor.Position.Row))
	}
	scrollRow = min(max(scrollRow, 0), cursor.Position.MaxRow)
	z.ClampScrollsToZoneSize()
	lastRowToRender = min(cursor.Position.MaxRow, scrollRow+int32(rowsToRender))

	return scrollRow, lastRowToRender
}
