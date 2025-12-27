package display

import (
	"log/slog"
	"slices"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/quar15/qq-go/internal/mode"
	"github.com/quar15/qq-go/internal/motion"
)

func HandleInput(ctx *mode.Context) {
	var (
		keyCharPressed int32          = rl.GetCharPressed()
		keyPressed     int32          = rl.GetKeyPressed()
		code           motion.KeyCode = motion.KeyRune
		modifiers      motion.Modifiers
		keyRune        rune
	)

	arrowKeys := []int32{rl.KeyLeft, rl.KeyDown, rl.KeyUp, rl.KeyRight}
	if slices.Contains(arrowKeys, keyPressed) {
		code = motion.KeyArrow
		keyCharPressed = keyPressed
	} else if keyPressed == rl.KeyEscape || keyPressed == rl.KeyCapsLock {
		code = motion.KeyEsc
		keyCharPressed = keyPressed
	}

	switch {
	case rl.IsKeyDown(rl.KeyLeftControl):
		modifiers = motion.ModCtrl
		keyCharPressed = keyPressed
	}

	if keyCharPressed == 0 {
		return
	}

	keyRune = rune(keyCharPressed)
	key := motion.Key{Code: code, Rune: keyRune, Modifiers: modifiers}

	slog.Debug("Handling key input", slog.Any("key", key))
	mode.HandleKey(ctx, key)
}
