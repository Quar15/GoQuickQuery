package config

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	background = rl.NewColor(30, 30, 46, 255)
	text       = rl.NewColor(205, 214, 244, 255)
	overlay0   = rl.NewColor(108, 112, 134, 255)
	overlay1   = rl.NewColor(127, 132, 156, 255)
	surface0   = rl.NewColor(49, 50, 68, 255)
	surface1   = rl.NewColor(69, 71, 90, 255)
	mantle     = rl.NewColor(24, 24, 37, 255)
	crust      = rl.NewColor(17, 17, 27, 255)
	blue       = rl.NewColor(137, 180, 250, 255)
	yellow     = rl.NewColor(249, 226, 175, 255)
	green      = rl.NewColor(166, 227, 161, 255)
	peach      = rl.NewColor(250, 179, 135, 255)
	pink       = rl.NewColor(245, 194, 231, 255)
	mauve      = rl.NewColor(203, 166, 247, 255)
)

type Color struct {
	rl.Color
}

type colors struct {
	cfg *ColorsConfig
}

type ColorsConfig struct {
	Background  Color `yaml:"background,omitempty"`
	Text        Color `yaml:"text,omitempty"`
	Accent      Color `yaml:"accent,omitempty"`
	Mantle      Color `yaml:"mantle,omitempty"`
	Crust       Color `yaml:"crust,omitempty"`
	Overlay0    Color `yaml:"overlay0,omitempty"`
	Overlay1    Color `yaml:"overlay1,omitempty"`
	Surface0    Color `yaml:"surface0,omitempty"`
	Surface1    Color `yaml:"surface1,omitempty"`
	NormalMode  Color `yaml:"normal_mode,omitempty"`
	InsertMode  Color `yaml:"insert_mode,omitempty"`
	VisualMode  Color `yaml:"visual_mode,omitempty"`
	CommandMode Color `yaml:"command_mode,omitempty"`
}

var predefinedColors = map[string]rl.Color{
	"background": background,
	"text":       text,
	"overlay0":   overlay0,
	"overlay1":   overlay1,
	"surface0":   surface0,
	"surface1":   surface1,
	"mantle":     mantle,
	"crust":      crust,
	"blue":       blue,
	"yellow":     yellow,
	"green":      green,
	"peach":      peach,
	"pink":       pink,
	"mauve":      mauve,
}

// UnmarshalYAML supports either hex or predefined name
func (c *Color) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}

	// lowercase & trim
	s = strings.ToLower(strings.TrimSpace(s))

	// check predefined colors
	if col, ok := predefinedColors[s]; ok {
		c.Color = col
		return nil
	}

	// check hex
	if strings.HasPrefix(s, "#") {
		rgba, err := parseHexColor(s)
		if err != nil {
			return err
		}
		c.Color = rgba
		return nil
	}

	return fmt.Errorf("unknown color: %s", s)
}

func parseHexColor(hex string) (rl.Color, error) {
	slog.Debug("Parsing hex color", slog.String("hex", hex))
	hex = strings.TrimPrefix(hex, "#")
	var r, g, b, a uint8 = 0, 0, 0, 255

	switch len(hex) {
	case 6:
		val, err := strconv.ParseUint(hex, 16, 32)
		if err != nil {
			return rl.Color{}, err
		}
		r = uint8(val >> 16)
		g = uint8((val >> 8) & 0xFF)
		b = uint8(val & 0xFF)
	case 8:
		val, err := strconv.ParseUint(hex, 16, 32)
		if err != nil {
			return rl.Color{}, err
		}
		r = uint8(val >> 24)
		g = uint8((val >> 16) & 0xFF)
		b = uint8((val >> 8) & 0xFF)
		a = uint8(val & 0xFF)
	default:
		return rl.Color{}, fmt.Errorf("invalid hex color length")
	}

	return rl.Color{R: r, G: g, B: b, A: a}, nil
}

func (c *ColorsConfig) MergeDefaults() {
	defaultColor := rl.NewColor(0, 0, 0, 0)
	if c.Background.Color == defaultColor {
		c.Background.Color = background
	}
	if c.Text.Color == defaultColor {
		c.Text.Color = text
	}
	if c.Accent.Color == defaultColor {
		c.Accent.Color = blue
	}
	if c.Mantle.Color == defaultColor {
		c.Mantle.Color = mantle
	}
	if c.Crust.Color == defaultColor {
		c.Crust.Color = crust
	}
	if c.Overlay0.Color == defaultColor {
		c.Overlay0.Color = overlay0
	}
	if c.Overlay1.Color == defaultColor {
		c.Overlay1.Color = overlay1
	}
	if c.Surface0.Color == defaultColor {
		c.Surface0.Color = surface0
	}
	if c.Surface1.Color == defaultColor {
		c.Surface1.Color = surface1
	}
	if c.NormalMode.Color == defaultColor {
		c.NormalMode.Color = blue
	}
	if c.InsertMode.Color == defaultColor {
		c.InsertMode.Color = green
	}
	if c.VisualMode.Color == defaultColor {
		c.VisualMode.Color = mauve
	}
	if c.CommandMode.Color == defaultColor {
		c.CommandMode.Color = blue
	}
}

func (c *colors) Background() rl.Color {
	return c.cfg.Background.Color
}

func (c *colors) Text() rl.Color {
	return c.cfg.Text.Color
}

func (c *colors) Accent() rl.Color {
	return c.cfg.Accent.Color
}

func (c *colors) NormalMode() rl.Color {
	return c.cfg.NormalMode.Color
}

func (c *colors) InsertMode() rl.Color {
	return c.cfg.InsertMode.Color
}

func (c *colors) VisualMode() rl.Color {
	return c.cfg.VisualMode.Color
}

func (c *colors) CommandMode() rl.Color {
	return c.cfg.CommandMode.Color
}

func (c *colors) Mantle() rl.Color {
	return c.cfg.Mantle.Color
}

func (c *colors) Crust() rl.Color {
	return c.cfg.Crust.Color
}

func (c *colors) Overlay0() rl.Color {
	return c.cfg.Overlay0.Color
}

func (c *colors) Overlay1() rl.Color {
	return c.cfg.Overlay1.Color
}

func (c *colors) Surface0() rl.Color {
	return c.cfg.Surface0.Color
}

func (c *colors) Surface1() rl.Color {
	return c.cfg.Surface1.Color
}

// @TODO: Consider moving to config static functions
func (c *colors) Blue() rl.Color {
	return blue
}

func (c *colors) Green() rl.Color {
	return green
}

func (c *colors) Mauve() rl.Color {
	return mauve
}

func (c *colors) Yellow() rl.Color {
	return yellow
}

func (c *colors) Peach() rl.Color {
	return peach
}
