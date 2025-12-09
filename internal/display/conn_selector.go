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
	const maxVisibleConnections int = 5
	const textPadding int32 = 6
	var bgColor rl.Color = colors.Mantle()

	var x int32 = (screenWidth - boxWidth) / 2
	renderedConnectionsN := len(config.Connections)
	if renderedConnectionsN > int(maxVisibleConnections) {
		renderedConnectionsN = maxVisibleConnections
	}
	var boxHeight int32 = int32(selectionHeight*renderedConnectionsN) + appAssets.MainFont.BaseSize/2 + selectionHeight/2
	var y int32 = (screenHeight - boxHeight) / 2

	var boxRectangle rl.RectangleInt32 = rl.RectangleInt32{
		X: x, Y: y, Width: boxWidth, Height: boxHeight,
	}
	const boxRoundness float32 = 0.05

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

	const iconWidth int32 = 16
	const iconHeight int32 = 16
	const iconPadding int32 = textPadding * 2
	const connNamePadding int32 = textPadding * 3
	const connStatusCircleRadius float32 = 2
	var maxNumberOfCharacters int32 = (boxWidth - iconPadding*2 - iconWidth - connNamePadding*2 - int32(connStatusCircleRadius)*2) / int32(appAssets.MainFontCharacterWidth)
	var textTopPadding int32 = boxRectangle.Y + textPadding*2 + appAssets.MainFont.BaseSize/2
	for i := CursorConnection.Position.Row; i < CursorConnection.Position.MaxRow; i++ {
		if i > CursorConnection.Position.Row+int32(renderedConnectionsN-1) {
			break
		}
		conn := config.Connections[i]
		var realConnection *database.ConnectionData = database.DBConnections[conn.Name]

		var displayName = conn.Name
		if len(displayName) > int(maxNumberOfCharacters) {
			displayName = displayName[:maxNumberOfCharacters]
		}
		var connTextColor rl.Color = colors.Text()
		if conn.Name == database.CurrDBConnection.Name {
			connTextColor = colors.Blue()
		}
		appAssets.DrawTextMainFont(
			displayName,
			rl.Vector2{
				X: float32(boxRectangle.X + textPadding*3 + iconWidth),
				Y: float32(textTopPadding),
			},
			connTextColor,
		)
		rl.DrawTexturePro(
			appAssets.Icons[conn.Driver],
			rl.Rectangle{X: 0, Y: 0, Width: float32(iconWidth), Height: float32(iconHeight)},
			rl.Rectangle{X: float32(boxRectangle.X + iconPadding), Y: float32(textTopPadding), Width: float32(iconWidth), Height: float32(iconHeight)},
			rl.Vector2{X: 0, Y: 0},
			0,
			rl.White,
		)
		if realConnection.Conn != false {
			rl.DrawCircle(boxRectangle.X+iconPadding+iconWidth, textTopPadding+iconHeight, connStatusCircleRadius, colors.Green())
		}
		textTopPadding += appAssets.MainFont.BaseSize + textPadding*2
	}
	// @TODO: Highlight current connection
}
