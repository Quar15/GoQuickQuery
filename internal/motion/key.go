package motion

import rl "github.com/gen2brain/raylib-go/raylib"

type KeyCode int

const (
	KeyRune KeyCode = iota
	KeyEsc
	KeyEnter
	KeyArrow
	KeySpecial
)

type Modifiers uint8

const (
	ModCtrl Modifiers = 1 << iota
)

type Key struct {
	Code      KeyCode
	Rune      rune
	Modifiers Modifiers
}

var CtrlW Key = Key{Code: KeyRune, Rune: rl.KeyW, Modifiers: ModCtrl}
var CtrlE Key = Key{Code: KeyRune, Rune: rl.KeyE, Modifiers: ModCtrl}
