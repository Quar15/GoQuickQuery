package display

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/quar15/qq-go/internal/config"
)

type Scrollbar struct {
	Track      rl.Rectangle
	Thumb      rl.Rectangle
	Dragging   bool
	GrabOffset float32
	Value      float32
}

func (s *Scrollbar) Draw() {
	rl.DrawRectangleRec(s.Track, config.Get().Colors.Mantle())
	var scrollbarColor rl.Color = config.Get().Colors.Overlay0()
	if s.Dragging {
		scrollbarColor = config.Get().Colors.Surface1()
	}
	rl.DrawRectangleRec(s.Track, scrollbarColor)
}

