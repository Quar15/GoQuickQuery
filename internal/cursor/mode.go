package cursor

import (
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/quar15/qq-go/internal/config"
)

type Mode int8

const (
	ModeNormal Mode = iota
	ModeInsert
	ModeVisual
	ModeVLine
	ModeVBlock
	ModeCommand
	ModeWindowManagement
)

var modeName = map[Mode]string{
	ModeNormal:           "NORMAL",
	ModeInsert:           "INSERT",
	ModeVisual:           "VISUAL",
	ModeVLine:            "V-LINE",
	ModeVBlock:           "V-BLOCK",
	ModeCommand:          "COMMAND",
	ModeWindowManagement: "WINDOW",
}

func (cm Mode) String() string {
	return modeName[cm]
}

var setupColorOnce sync.Once
var modeColor = map[Mode]rl.Color{}

func (cm Mode) Color() rl.Color {
	setupColorOnce.Do(func() {
		modeColor = map[Mode]rl.Color{
			ModeNormal:           config.Get().Colors.NormalMode(),
			ModeInsert:           config.Get().Colors.InsertMode(),
			ModeVisual:           config.Get().Colors.VisualMode(),
			ModeVLine:            config.Get().Colors.VisualMode(),
			ModeVBlock:           config.Get().Colors.VisualMode(),
			ModeCommand:          config.Get().Colors.CommandMode(),
			ModeWindowManagement: config.Get().Colors.CommandMode(),
		}
	})
	return modeColor[cm]
}
