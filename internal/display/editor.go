package display

import (
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/quar15/qq-go/internal/assets"
	"github.com/quar15/qq-go/internal/colors"
	"github.com/quar15/qq-go/internal/utilities"
)

func (z *Zone) DrawEditor(appAssets *assets.Assets, eg *EditorGrid) {
	const cellHeight int = 28
	const textPadding int32 = 6
	const rowsInitialPadding int32 = 16

	var contentHeight int = 0

	var counterColumnCharactersCount int = utilities.CountDigits(int(eg.Rows))
	var counterColumnWidth int = int(appAssets.MainFontCharacterWidth)*counterColumnCharactersCount + int(textPadding*2)

	var rowsToRender int8 = z.GetNumberOfVisibleRows(int32(cellHeight)) + 1
	var firstVisibleRowToScrollIndex int32 = min(max(int32(z.Scroll.Y)/int32(cellHeight), 0), eg.Rows) // Swapping screens creates weird behavior
	var lastRowToRender = min(eg.Rows, firstVisibleRowToScrollIndex+int32(rowsToRender))

	rl.BeginScissorMode(int32(z.Bounds.X), int32(z.Bounds.Y), int32(z.Bounds.Width), int32(z.Bounds.Height))
	rl.ClearBackground(colors.Background())

	for row := firstVisibleRowToScrollIndex; row < lastRowToRender; row++ {
		var cellY float32 = z.Bounds.Y + float32(row*int32(cellHeight)) - z.Scroll.Y + float32(rowsInitialPadding) + float32(textPadding)
		if eg.Highlight[row] != nil {
			for c := range eg.Text[row] {
				var cellX float32 = z.Bounds.X + float32(counterColumnWidth) + float32(textPadding) + float32(c*int(appAssets.MainFontCharacterWidth))
				rl.DrawTextEx(
					appAssets.MainFont,
					string(eg.Text[row][c]),
					rl.Vector2{X: float32(cellX), Y: float32(cellY)},
					appAssets.MainFontSize,
					appAssets.MainFontSpacing,
					eg.Highlight[row][c].Color(),
				)
			}
		}
	}
	// Draw row counter
	for row := firstVisibleRowToScrollIndex; row < lastRowToRender; row++ {
		var counterColumnLeftPadding float32 = float32(textPadding) + float32(counterColumnCharactersCount-utilities.CountDigits(int(row)+1))*appAssets.MainFontCharacterWidth
		var cellX float32 = z.Bounds.X + counterColumnLeftPadding
		var cellY float32 = z.Bounds.Y + float32(row*int32(cellHeight)) - z.Scroll.Y + float32(rowsInitialPadding) + float32(textPadding)
		appAssets.DrawTextMainFont(strconv.Itoa(int(row+1)), rl.Vector2{X: cellX, Y: cellY}, colors.Overlay0())
	}
	// Draw text

	rl.EndScissorMode()

	contentHeight += cellHeight*int(eg.Rows) + int(rowsInitialPadding)

	z.ContentSize.Y = max(float32(contentHeight), z.Bounds.Height)
	z.ContentSize.X = max(eg.MaxWidth, z.Bounds.Width)
	z.drawScrollbars()
}
