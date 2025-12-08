package assets

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Assets struct {
	MainFont               rl.Font
	MainFontSpacing        float32
	MainFontSize           float32
	MainFontCharacterWidth float32
	Icons                  map[string]rl.Texture2D
}

func (a *Assets) loadFont() error {
	a.MainFontSize = 18
	a.MainFontSpacing = 1.0
	a.MainFont = rl.LoadFontEx("assets/fonts/FiraCodeNerdFontMono-Regular.ttf", int32(a.MainFontSize), nil, 0)
	if !rl.IsFontValid(a.MainFont) {
		return errors.New("Failed to load font")
	}
	a.MainFontCharacterWidth = rl.MeasureTextEx(a.MainFont, "X", a.MainFontSize, a.MainFontSpacing).X

	return nil
}

func (a *Assets) loadIcons() error {
	a.Icons = make(map[string]rl.Texture2D)

	folder := "./assets/img/icons/"
	files, err := os.ReadDir(folder) // @TODO: Consider including assets in build
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		ext := filepath.Ext(file.Name())
		switch ext {
		case ".png", ".jpg", ".jpeg", ".bmp", ".tga", ".gif":
			baseName := strings.TrimSuffix(file.Name(), ext)
			fullpath := filepath.Join(folder, file.Name())
			texture := rl.LoadTexture(fullpath)
			a.Icons[baseName] = texture
			slog.Debug("Loaded icon", slog.String("key", baseName), slog.String("path", fullpath))
		default:
			continue
		}

	}

	return nil
}

func (a *Assets) LoadAssets() error {
	if err := a.loadFont(); err != nil {
		return err
	}
	if err := a.loadIcons(); err != nil {
		return err
	}
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
