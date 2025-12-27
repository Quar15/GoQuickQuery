package display

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/quar15/qq-go/internal/config"
	"github.com/quar15/qq-go/internal/cursor"
)

type Splitter struct {
	Rect     rl.Rectangle
	Dragging bool
	Ratio    float32
	Height   float32
	Y        float32
}

func (splitter *Splitter) HandleZoneSplit(screenWidth int, screenHeight int, commandZoneHeight int) {
	splitter.Y = splitter.Ratio * float32(screenHeight)

	var mouse rl.Vector2 = rl.GetMousePosition()

	splitter.Rect = rl.Rectangle{
		X:      0,
		Y:      splitter.Y - splitter.Height/2,
		Width:  float32(screenWidth),
		Height: splitter.Height,
	}
	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) && rl.CheckCollisionPointRec(mouse, splitter.Rect) {
		splitter.Dragging = true
	}
	if rl.IsMouseButtonReleased(rl.MouseButtonLeft) {
		splitter.Dragging = false
	}
	if splitter.Dragging {
		rl.SetMouseCursor(rl.MouseCursorResizeNS)
		splitter.Y = rl.Clamp(mouse.Y, 10, float32(screenHeight)-float32(commandZoneHeight)-10)
		splitter.Ratio = splitter.Y / float32(screenHeight)
	}
}

func (splitter *Splitter) Draw(currCursorType cursor.Type) {
	var defaultColor rl.Color = config.Get().Colors.Crust()
	var focusColor rl.Color = config.Get().Colors.Accent()
	rl.DrawRectangleRec(splitter.Rect, defaultColor)
	switch currCursorType {
	case cursor.TypeEditor:
		var newRect rl.Rectangle = splitter.Rect
		newRect.Width /= 2
		newRect.Height /= 2
		rl.DrawRectangleRec(newRect, focusColor)
	case cursor.TypeSpreadsheet:
		var newRect rl.Rectangle = splitter.Rect
		newRect.Width /= 2
		newRect.Height /= 2
		newRect.Y += newRect.Height
		newRect.X += newRect.Width
		rl.DrawRectangleRec(newRect, focusColor)
	}
}
