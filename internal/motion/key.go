package motion

type KeyCode int

const (
	KeyRune KeyCode = iota
	KeyEsc
	KeyEnter
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
