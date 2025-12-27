package display

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/quar15/qq-go/internal/assets"
	"github.com/quar15/qq-go/internal/config"
	"github.com/quar15/qq-go/internal/cursor"
	"github.com/quar15/qq-go/internal/database"
)

func (z *Zone) DrawConnectionSelector(appAssets *assets.Assets, config *config.Config, cursor *cursor.Cursor, screenWidth int32, screenHeight int32, connManager *database.ConnectionManager) {
	const boxWidth = 300
	const maxVisibleConnections int = 5
	const textPadding int32 = 6
	var cellHeight = appAssets.MainFont.BaseSize + textPadding*2
	var bgColor rl.Color = config.Colors.Mantle()

	var x int32 = (screenWidth - boxWidth) / 2
	renderedConnectionsN := len(config.Connections)
	if renderedConnectionsN > int(maxVisibleConnections) {
		renderedConnectionsN = maxVisibleConnections
	}
	var boxHeight int32 = cellHeight*int32(renderedConnectionsN) + appAssets.MainFont.BaseSize/2 + cellHeight/2
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
		config.Colors.Accent(),
	)
	const boxHeaderText string = "Connections:"
	var boxHeaderTextWidth int32 = int32(appAssets.MeasureTextMainFont(boxHeaderText).X) + textPadding*2
	var boxHeaderTextX int32 = x + boxWidth/2 - boxHeaderTextWidth/2
	rl.DrawRectangle(boxHeaderTextX-textPadding, y, boxHeaderTextWidth, appAssets.MainFont.BaseSize, bgColor)
	appAssets.DrawTextMainFont(
		"Connections:",
		rl.Vector2{X: float32(boxHeaderTextX), Y: float32(y)},
		config.Colors.Accent(),
	)

	const initialSelectionsTopPadding float32 = 20
	z.Bounds = boxRectangle.ToFloat32()
	z.Bounds.Y += initialSelectionsTopPadding
	z.Bounds.Height = float32(maxVisibleConnections * int(cellHeight))
	z.Scroll.Y = float32(cellHeight * cursor.Position.Row)
	z.Scroll.X = 0
	z.ContentSize.Y = float32(cellHeight * (cursor.Position.MaxRow + 1))
	z.ContentSize.X = 0
	z.ClampScrollsToZoneSize()

	rl.DrawRectangle(
		int32(z.Bounds.X),
		int32(z.Bounds.Y)+(cursor.Position.Row*cellHeight)-int32(z.Scroll.Y)-textPadding,
		boxWidth,
		cellHeight,
		config.Colors.Surface0(),
	)

	rl.BeginScissorMode(int32(z.Bounds.X), int32(z.Bounds.Y), int32(z.Bounds.Width), int32(z.Bounds.Height))

	const iconWidth int32 = 16
	const iconHeight int32 = 16
	const iconPadding int32 = textPadding * 2
	const connNamePadding int32 = textPadding * 3
	const connStatusCircleRadius float32 = 2
	var maxNumberOfCharacters int32 = (boxWidth - iconPadding*2 - iconWidth - connNamePadding*2 - int32(connStatusCircleRadius)*2) / int32(appAssets.MainFontCharacterWidth)
	var firstRowToRender int32 = max(int32(z.Scroll.Y)/int32(cellHeight), 0)
	var lastRowToRender int32 = min(cursor.Position.Row+int32(renderedConnectionsN), cursor.Position.MaxRow)
	for i := firstRowToRender; i <= lastRowToRender; i++ {
		conn := config.Connections[i]
		var displayName = conn.Name
		if len(displayName) > int(maxNumberOfCharacters) {
			displayName = displayName[:maxNumberOfCharacters]
		}
		var connTextColor rl.Color = config.Colors.Text()
		if conn.Name == connManager.GetCurrentConnectionName() {
			connTextColor = config.Colors.Accent()
		}
		var cellY float32 = z.Bounds.Y + float32(i*cellHeight) - z.Scroll.Y
		appAssets.DrawTextMainFont(
			displayName,
			rl.Vector2{
				X: z.Bounds.X + float32(textPadding*3+iconWidth),
				Y: cellY,
			},
			connTextColor,
		)
		rl.DrawTexturePro(
			appAssets.Icons[conn.Driver],
			rl.Rectangle{X: 0, Y: 0, Width: float32(iconWidth), Height: float32(iconHeight)},
			rl.Rectangle{X: z.Bounds.X + float32(iconPadding), Y: cellY, Width: float32(iconWidth), Height: float32(iconHeight)},
			rl.Vector2{X: 0, Y: 0},
			0,
			rl.White,
		)
		if connManager.IsConnectionAlive(conn.Name) {
			rl.DrawCircle(boxRectangle.X+iconPadding+iconWidth, int32(cellY)+iconHeight, connStatusCircleRadius, config.Colors.Green())
		}
	}
	rl.EndScissorMode()
}
