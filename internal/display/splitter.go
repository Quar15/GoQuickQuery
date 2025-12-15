package display

import (
	"github.com/gen2brain/raylib-go/raylib"
	"github.com/quar15/qq-go/internal/colors"
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

func (splitter *Splitter) Draw() {
	var defaultColor rl.Color = colors.Crust()
	var focusColor rl.Color = colors.Blue()
	rl.DrawRectangleRec(splitter.Rect, defaultColor)
	switch CurrCursor.Type {
	case CursorTypeEditor:
		var newRect rl.Rectangle = splitter.Rect
		newRect.Width /= 2
		newRect.Height /= 2
		rl.DrawRectangleRec(newRect, focusColor)
	case CursorTypeSpreadsheet:
		var newRect rl.Rectangle = splitter.Rect
		newRect.Width /= 2
		newRect.Height /= 2
		newRect.Y += newRect.Height
		newRect.X += newRect.Width
		rl.DrawRectangleRec(newRect, focusColor)
	}
}
