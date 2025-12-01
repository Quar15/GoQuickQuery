package colors

import rl "github.com/gen2brain/raylib-go/raylib"

var (
	background = rl.NewColor(30, 30, 46, 255)
	text       = rl.NewColor(205, 214, 244, 255)
	overlay0   = rl.NewColor(108, 112, 134, 255)
	overlay1   = rl.NewColor(127, 132, 156, 255)
	surface1   = rl.NewColor(69, 71, 90, 255)
	surface0   = rl.NewColor(49, 50, 68, 255)
	mantle     = rl.NewColor(24, 24, 37, 255)
	crust      = rl.NewColor(17, 17, 27, 255)
	blue       = rl.NewColor(137, 180, 250, 255)
	green      = rl.NewColor(166, 227, 161, 255)
	peach      = rl.NewColor(250, 179, 135, 255)
	pink       = rl.NewColor(245, 194, 231, 255)
	mauve      = rl.NewColor(203, 166, 247, 255)
)

func Background() rl.Color { return background }
func Text() rl.Color       { return text }
func Overlay0() rl.Color   { return overlay0 }
func Overlay1() rl.Color   { return overlay1 }
func Surface1() rl.Color   { return surface1 }
func Surface0() rl.Color   { return surface0 }
func Mantle() rl.Color     { return mantle }
func Crust() rl.Color      { return crust }
func Blue() rl.Color       { return blue }
func Green() rl.Color      { return green }
func Peach() rl.Color      { return peach }
func Pink() rl.Color       { return pink }
func Mauve() rl.Color      { return mauve }
