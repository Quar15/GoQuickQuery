package setup

import (
	"log/slog"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/quar15/qq-go/internal/motion"
)

const keySmallG rune = 103
const keySmallH rune = 104
const keySmallJ rune = 106
const keySmallK rune = 107
const keySmallL rune = 108

func EditorMotionSet() *motion.Set {
	s := motion.NewSet()
	s.AddRune(keySmallH, motion.MoveLeft{})
	s.AddRune(keySmallJ, motion.MoveDown{})
	s.AddRune(keySmallK, motion.MoveUp{})
	s.AddRune(keySmallL, motion.MoveRight{})

	s.AddArrow(rl.KeyLeft, motion.MoveLeft{})
	s.AddArrow(rl.KeyDown, motion.MoveDown{})
	s.AddArrow(rl.KeyUp, motion.MoveUp{})
	s.AddArrow(rl.KeyRight, motion.MoveRight{})

	s.AddRune(rl.KeyG, motion.MoveEndDown{})
	s.AddRune(rl.KeyG, motion.MoveToSpecificLineOrDown{})
	s.Add([]motion.Key{
		{Code: motion.KeyRune, Rune: keySmallG},
		{Code: motion.KeyRune, Rune: keySmallG},
	}, motion.MoveStartUp{})

	slog.Debug("Initialized editor motion set", slog.Any("setTrie", s.Root()))
	return s
}

func SpreadsheetMotionSet() *motion.Set {
	// same as editor for now
	return EditorMotionSet()
}

func ConnectionsMotionSet() *motion.Set {
	s := motion.NewSet()
	s.AddRune('j', motion.MoveDown{})
	s.AddRune('k', motion.MoveUp{})
	return s
}
