package display

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
