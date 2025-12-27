package cursor

import (
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

var modeColor = map[Mode]rl.Color{}

func InitModeColors(cfg *config.Config) {
	modeColor = map[Mode]rl.Color{
		ModeNormal:           cfg.Colors.NormalMode(),
		ModeInsert:           cfg.Colors.InsertMode(),
		ModeVisual:           cfg.Colors.VisualMode(),
		ModeVLine:            cfg.Colors.VisualMode(),
		ModeVBlock:           cfg.Colors.VisualMode(),
		ModeCommand:          cfg.Colors.CommandMode(),
		ModeWindowManagement: cfg.Colors.CommandMode(),
	}
}

func (cm Mode) Color() rl.Color {
	return modeColor[cm]
}
