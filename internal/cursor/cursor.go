package cursor

import "github.com/quar15/qq-go/internal/motion"

type Common struct {
	Mode      Mode
	CmdBuf    string
	MotionBuf string
	Logs      CommandLogs
}

type Type int8

const (
	TypeEditor Type = iota
	TypeSpreadsheet
	TypeConnections
)

type Cursor struct {
	Common   *Common
	Position motion.CursorPosition
	Type     Type
	isActive bool
}
