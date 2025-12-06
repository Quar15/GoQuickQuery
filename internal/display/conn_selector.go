package display

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/quar15/qq-go/internal/assets"
	"github.com/quar15/qq-go/internal/colors"
	"github.com/quar15/qq-go/internal/config"
	"github.com/quar15/qq-go/internal/database"
)

func DrawConnectionSelector(appAssets *assets.Assets, config *config.Config, screenWidth int32, screenHeight int32) {
	const boxWidth = 300
	const selectionHeight = 30
	const maxVisibleConnections = 5
	const textPadding int32 = 6
	var bgColor rl.Color = colors.Background()

	var x int32 = (screenWidth - boxWidth) / 2
	var boxHeight int32 = int32(selectionHeight*(len(database.DBConnections))) + appAssets.MainFont.BaseSize/2 + selectionHeight/2
	var y int32 = (screenHeight - boxHeight) / 2

	var boxRectangle rl.RectangleInt32 = rl.RectangleInt32{
		X: x, Y: y, Width: boxWidth, Height: boxHeight,
	}
	const boxRoundness float32 = 0.2

	boxRectangle.X -= appAssets.MainFont.BaseSize / 2
	boxRectangle.Height += appAssets.MainFont.BaseSize
	boxRectangle.Width += appAssets.MainFont.BaseSize
	rl.DrawRectangleRounded(
		boxRectangle.ToFloat32(),
		boxRoundness,
		0.0,
		bgColor,
	)
	boxRectangle.Y += appAssets.MainFont.BaseSize / 2
	boxRectangle.X += appAssets.MainFont.BaseSize / 2
	boxRectangle.Height -= appAssets.MainFont.BaseSize
	boxRectangle.Width -= appAssets.MainFont.BaseSize
	rl.DrawRectangleRoundedLinesEx(
		boxRectangle.ToFloat32(),
		boxRoundness,
		0.0,
		2,
		colors.Blue(),
	)
	const boxHeaderText string = "Connections:"
	var boxHeaderTextWidth int32 = int32(appAssets.MeasureTextMainFont(boxHeaderText).X) + textPadding*2
	var boxHeaderTextX int32 = x + boxWidth/2 - boxHeaderTextWidth/2
	rl.DrawRectangle(boxHeaderTextX-textPadding, y, boxHeaderTextWidth, appAssets.MainFont.BaseSize, bgColor)
	appAssets.DrawTextMainFont(
		"Connections:",
		rl.Vector2{X: float32(boxHeaderTextX), Y: float32(y)},
		colors.Blue(),
	)

	var textTopPadding int32 = boxRectangle.Y + textPadding*2 + appAssets.MainFont.BaseSize/2
	for _, conn := range config.Connections {
		appAssets.DrawTextMainFont(
			conn.Name,
			rl.Vector2{
				X: float32(boxRectangle.X + textPadding*2),
				Y: float32(textTopPadding),
			},
			colors.Text(),
		)
		textTopPadding += appAssets.MainFont.BaseSize + textPadding*2
	}
	// @TODO: Highlight current connection
}
