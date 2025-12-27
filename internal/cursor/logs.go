package cursor

import "log/slog"

type CommandLogs struct {
	Channel     chan string
	LastMessage string
}

func (cl *CommandLogs) Init() {
	cl.Channel = make(chan string, 10) // @TODO: This can cause issues if filled
}

func (cl *CommandLogs) CheckForMessage() {
	select {
	case msg := <-cl.Channel:
		cl.LastMessage = msg
	default:
	}
}

func (cl *CommandLogs) Log(msg string) {
	select {
	case cl.Channel <- msg:
	default:
		// Drop message if channel full
		slog.Warn("Log channel full, dropping message")
	}
}
