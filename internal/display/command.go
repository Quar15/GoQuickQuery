package display

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/quar15/qq-go/internal/assets"
	"github.com/quar15/qq-go/internal/colors"
	"github.com/quar15/qq-go/internal/database"
)

func (z *Zone) DrawCommandZone(appAssets *assets.Assets, c *Cursor) {
	const textSpacing float32 = 4
	var statusLineColor rl.Color = c.Common.Mode.Color()
	// Status Line
	rl.DrawRectangle(int32(z.Bounds.X), int32(z.Bounds.Y), int32(z.Bounds.Width), int32(z.Bounds.Height/2), colors.Mantle())
	var modeStatusText string = c.Common.Mode.String()
	var modeStatusTextWidth float32 = appAssets.MeasureTextMainFont(modeStatusText).X
	rl.DrawRectangle(int32(z.Bounds.X), int32(z.Bounds.Y), int32(modeStatusTextWidth+textSpacing*4), int32(z.Bounds.Height/2), statusLineColor)
	// @TODO: Add horizontal spacing
	appAssets.DrawTextMainFont(modeStatusText, rl.Vector2{X: z.Bounds.X + textSpacing*2, Y: z.Bounds.Y + textSpacing/2}, colors.Mantle())
	var detailsStatusText string = fmt.Sprintf("%d/%d | %d/%d | %d%%", c.Position.Col+1, c.Position.MaxCol+1, c.Position.Row+1, c.Position.MaxRow+1, 0)
	var detailsStatusTextWidth float32 = appAssets.MeasureTextMainFont(detailsStatusText).X
	var detailsStatusWidth float32 = detailsStatusTextWidth + textSpacing*4
	rl.DrawRectangle(int32(z.Bounds.Width-detailsStatusWidth), int32(z.Bounds.Y), int32(detailsStatusWidth), int32(z.Bounds.Height/2), statusLineColor)
	appAssets.DrawTextMainFont(detailsStatusText, rl.Vector2{X: z.Bounds.Width - z.Bounds.X - detailsStatusWidth + textSpacing*2, Y: z.Bounds.Y + textSpacing/2}, colors.Mantle())

	var connectionStatusText = database.CurrDBConnection.Name
	var connectionStatusTextWidth float32 = appAssets.MeasureTextMainFont(connectionStatusText).X
	var connectionStatusTextX float32 = z.Bounds.Width - z.Bounds.X - detailsStatusWidth - connectionStatusTextWidth - textSpacing*2
	appAssets.DrawTextMainFont(
		connectionStatusText,
		rl.Vector2{X: connectionStatusTextX, Y: z.Bounds.Y + textSpacing/2},
		colors.Text(),
	)
	const iconWidth int32 = 16
	const iconHeight int32 = 16
	rl.DrawTexturePro(
		appAssets.Icons[database.CurrDBConnection.Driver],
		rl.Rectangle{X: 0, Y: 0, Width: float32(iconWidth), Height: float32(iconHeight)},
		rl.Rectangle{X: connectionStatusTextX - textSpacing*2 - float32(iconWidth), Y: z.Bounds.Y + textSpacing/2, Width: float32(iconWidth), Height: float32(iconHeight)},
		rl.Vector2{X: 0, Y: 0},
		0,
		rl.White,
	)

	// Command Input
	rl.DrawRectangle(int32(z.Bounds.X), int32(z.Bounds.Y+z.Bounds.Height/2), int32(z.Bounds.Width), int32(z.Bounds.Height/2), colors.Background())
	c.Common.Logs.CheckForMessage()
	appAssets.DrawTextMainFont(c.Common.Logs.LastMessage, rl.Vector2{X: z.Bounds.X + textSpacing, Y: z.Bounds.Y + z.Bounds.Height/2 + textSpacing/2}, colors.Text())

	var motionBufWidth float32 = appAssets.MeasureTextMainFont(c.Common.MotionBuf).X + textSpacing*8
	appAssets.DrawTextMainFont(c.Common.MotionBuf, rl.Vector2{X: z.Bounds.Width - z.Bounds.X - motionBufWidth, Y: z.Bounds.Y + z.Bounds.Height/2 + textSpacing/2}, colors.Text())
}
