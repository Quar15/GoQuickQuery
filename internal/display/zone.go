package display

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/quar15/qq-go/internal/assets"
	"github.com/quar15/qq-go/internal/colors"
)

type Zone struct {
	Bounds      rl.Rectangle
	Scroll      rl.Vector2
	ContentSize rl.Vector2
	vScrollbar  Scrollbar
	hScrollbar  Scrollbar
}

func (z *Zone) MouseInside(mouse rl.Vector2) bool {
	return rl.CheckCollisionPointRec(mouse, z.Bounds)
}

func (z *Zone) drawScrollbars() {
	if z.ContentSize.Y > z.Bounds.Height {
		z.vScrollbar.Draw()
	}
	if z.ContentSize.X > z.Bounds.Width {
		z.hScrollbar.Draw()
	}
}

func (z *Zone) Draw(appAssets *assets.Assets) {
	rl.BeginScissorMode(int32(z.Bounds.X), int32(z.Bounds.Y), int32(z.Bounds.Width), int32(z.Bounds.Height))
	rl.ClearBackground(colors.Background())

	rl.DrawTextEx(appAssets.MainFont, "TESTING ABCDEFGHIJKLMNOPRSTUWXYZ", rl.Vector2{X: z.Bounds.X + 20 - z.Scroll.X, Y: z.Bounds.Y + 20 - z.Scroll.Y}, appAssets.MainFontSize, appAssets.MainFontSpacing, colors.Text())

	rl.EndScissorMode()

	z.drawScrollbars()
}

func (z *Zone) ClampScrollsToZoneSize() {
	if z.Scroll.X < 0 {
		z.Scroll.X = 0
	}
	if z.Scroll.Y < 0 {
		z.Scroll.Y = 0
	}
	if z.Scroll.X > z.ContentSize.X-float32(z.Bounds.Width) {
		z.Scroll.X = z.ContentSize.X - float32(z.Bounds.Width)
	}
	if z.Scroll.Y > z.ContentSize.Y-float32(z.Bounds.Height) {
		z.Scroll.Y = z.ContentSize.Y - float32(z.Bounds.Height)
	}
}

func (z *Zone) InitZoneScrollbars() {
	// Vertical scrollbar
	z.vScrollbar.Value = float32(z.Bounds.Height) / z.ContentSize.Y
	z.vScrollbar.Track = rl.Rectangle{X: (z.Bounds.X + z.Bounds.Width - 10), Y: z.Bounds.Y, Width: 10, Height: z.Bounds.Height - 10}
	z.vScrollbar.Thumb = rl.Rectangle{
		X:      z.vScrollbar.Track.X,
		Y:      z.vScrollbar.Track.Y + (z.Scroll.Y/(z.ContentSize.Y-z.Bounds.Height))*(z.vScrollbar.Track.Height-z.vScrollbar.Value*z.vScrollbar.Track.Height),
		Width:  10,
		Height: z.vScrollbar.Value * z.vScrollbar.Track.Height,
	}

	// Horizontal scrollbar
	z.hScrollbar.Value = float32(z.Bounds.Width) / z.ContentSize.X
	z.hScrollbar.Track = rl.Rectangle{X: z.Bounds.X, Y: z.Bounds.Y + z.Bounds.Height - 10, Width: z.Bounds.Width - 10, Height: 10}
	z.hScrollbar.Thumb = rl.Rectangle{
		X:      z.hScrollbar.Track.X + (z.Scroll.X/(z.ContentSize.X-z.Bounds.Width))*(z.hScrollbar.Track.Width-z.hScrollbar.Value*z.hScrollbar.Track.Width),
		Y:      z.hScrollbar.Track.Y,
		Width:  z.hScrollbar.Value * z.hScrollbar.Track.Width,
		Height: 10,
	}
}

func (z *Zone) UpdateZoneScroll() {
	z.InitZoneScrollbars()

	var mouse rl.Vector2 = rl.GetMousePosition()
	var mouseWheelStep float32 = 40

	// Only scroll if mouse inside the zone
	if rl.CheckCollisionPointRec(mouse, z.Bounds) {
		if rl.IsKeyDown(rl.KeyLeftShift) {
			// Mouse wheel scroll (horizontal)
			z.Scroll.X -= rl.GetMouseWheelMove() * mouseWheelStep
		} else {
			// Mouse wheel scroll (vertical)
			z.Scroll.Y -= rl.GetMouseWheelMove() * mouseWheelStep
		}
	}

	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) && rl.CheckCollisionPointRec(mouse, z.vScrollbar.Thumb) {
		z.vScrollbar.Dragging = true
		z.vScrollbar.GrabOffset = mouse.Y - z.vScrollbar.Thumb.Y
	}
	if rl.IsMouseButtonReleased(rl.MouseButtonLeft) {
		z.vScrollbar.Dragging = false
	}

	if z.vScrollbar.Dragging {
		rl.SetMouseCursor(rl.MouseCursorPointingHand)
		var newY float32 = mouse.Y - z.vScrollbar.GrabOffset
		var posRatio float32 = (newY - z.vScrollbar.Track.Y) / (z.vScrollbar.Track.Height - z.vScrollbar.Thumb.Height)
		posRatio = rl.Clamp(posRatio, 0, 1)
		z.vScrollbar.Value = posRatio
		z.Scroll.Y = posRatio * (z.ContentSize.Y - z.Bounds.Height)
	}

	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) && rl.CheckCollisionPointRec(mouse, z.hScrollbar.Thumb) {
		z.hScrollbar.Dragging = true
		z.hScrollbar.GrabOffset = mouse.X - z.hScrollbar.Thumb.X
	}
	if rl.IsMouseButtonReleased(rl.MouseButtonLeft) {
		z.hScrollbar.Dragging = false
	}

	if z.hScrollbar.Dragging {
		rl.SetMouseCursor(rl.MouseCursorPointingHand)
		var newX float32 = mouse.X - z.hScrollbar.GrabOffset
		var posRatio float32 = (newX - z.hScrollbar.Track.X) / (z.hScrollbar.Track.Width - z.hScrollbar.Thumb.Width)
		posRatio = rl.Clamp(posRatio, 0, 1)
		z.hScrollbar.Value = posRatio
		z.Scroll.X = posRatio * (z.ContentSize.X - z.Bounds.Width)
	}

	z.ClampScrollsToZoneSize()
}
