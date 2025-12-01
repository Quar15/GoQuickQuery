package display

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/quar15/qq-go/internal/colors"
)

type Scrollbar struct {
	Track      rl.Rectangle
	Thumb      rl.Rectangle
	Dragging   bool
	GrabOffset float32
	Value      float32
}

func (s *Scrollbar) Draw() {
	rl.DrawRectangleRec(s.Track, colors.Mantle())
	var scrollbarColor rl.Color = colors.Overlay0()
	if s.Dragging {
		scrollbarColor = colors.Surface1()
	}
	rl.DrawRectangleRec(s.Thumb, scrollbarColor)
}
