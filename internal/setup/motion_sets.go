package setup

import (
	"log/slog"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/quar15/qq-go/internal/mode"
	"github.com/quar15/qq-go/internal/mode/commands"
	"github.com/quar15/qq-go/internal/motion"
)

const keySmallG rune = 103
const keySmallH rune = 104
const keySmallJ rune = 106
const keySmallK rune = 107
const keySmallL rune = 108

func baseMotionSet() *motion.Set {
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

	return s
}

func baseCommandRegistry() *mode.CommandRegistry {
	cr := mode.NewCommandRegistry()

	cr.Bind(motion.Key{Code: motion.KeyRune, Rune: 'E', Modifiers: motion.ModCtrl}, commands.ConnectionsSwap{})
	cr.Bind(motion.Key{Code: motion.KeyRune, Rune: 'W', Modifiers: motion.ModCtrl}, mode.WindowManagementModeActivate{})

	return cr
}

func EditorMotionSet() (*motion.Set, *mode.CommandRegistry) {
	s := baseMotionSet()
	cr := baseCommandRegistry()
	cr.Bind(
		motion.Key{Code: motion.KeyEnter, Rune: rl.KeyEnter, Modifiers: motion.ModCtrl},
		commands.ExecuteSQLCommand{},
	)

	slog.Debug("Initialized editor motion set", slog.Any("setTrie", s.Root()))
	return s, cr
}

func SpreadsheetMotionSet() (*motion.Set, *mode.CommandRegistry) {
	s := baseMotionSet()
	cr := baseCommandRegistry()

	slog.Debug("Initialized spreadsheet motion set", slog.Any("setTrie", s.Root()))
	return s, cr
}

func ConnectionsMotionSet() (*motion.Set, *mode.CommandRegistry) {
	s := motion.NewSet()
	s.AddRune('j', motion.MoveDown{})
	s.AddRune('k', motion.MoveUp{})

	s.AddArrow(rl.KeyDown, motion.MoveDown{})
	s.AddArrow(rl.KeyUp, motion.MoveUp{})

	s.AddRune(rl.KeyG, motion.MoveEndDown{})
	s.AddRune(rl.KeyG, motion.MoveToSpecificLineOrDown{})
	s.Add([]motion.Key{
		{Code: motion.KeyRune, Rune: keySmallG},
		{Code: motion.KeyRune, Rune: keySmallG},
	}, motion.MoveStartUp{})

	cr := baseCommandRegistry()
	cr.Bind(motion.Key{Code: motion.KeyEnter, Rune: rl.KeyEnter}, commands.ConnectionsChange{})
	cr.Bind(motion.Key{Code: motion.KeyEsc, Rune: rl.KeyEscape}, commands.ConnectionsExit{})
	cr.Bind(motion.Key{Code: motion.KeyEsc, Rune: rl.KeyCapsLock}, commands.ConnectionsExit{})

	slog.Debug("Initialized connections motion set", slog.Any("setTrie", s.Root()), slog.Any("cr", cr))
	return s, cr
}
