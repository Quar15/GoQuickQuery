package assets

import (
	"errors"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Assets struct {
	MainFont               rl.Font
	MainFontSpacing        float32
	MainFontSize           float32
	MainFontCharacterWidth float32
}

func (a *Assets) LoadAssets() error {
	a.MainFontSize = 18
	a.MainFontSpacing = 1.0
	a.MainFont = rl.LoadFontEx("assets/fonts/FiraCodeNerdFontMono-Regular.ttf", int32(a.MainFontSize), nil, 0)
	if !rl.IsFontValid(a.MainFont) {
		return errors.New("Failed to load font")
	}
	a.MainFontCharacterWidth = rl.MeasureTextEx(a.MainFont, "X", a.MainFontSize, a.MainFontSpacing).X
	return nil
}

func (a *Assets) DrawTextMainFont(text string, position rl.Vector2, color rl.Color) {
	rl.DrawTextEx(a.MainFont, text, position, a.MainFontSize, a.MainFontSpacing, color)
}

func (a *Assets) MeasureTextMainFont(text string) rl.Vector2 {
	return rl.MeasureTextEx(a.MainFont, text, a.MainFontSize, a.MainFontSpacing)
}

func (a *Assets) UnloadAssets() {
	rl.UnloadFont(a.MainFont)
}
