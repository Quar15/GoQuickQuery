package mode

import "github.com/quar15/qq-go/internal/motion"

type Command interface {
	Execute(ctx *Context) error
}

type CommandRegistry struct {
	bindings map[motion.Key]Command
}

func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{
		bindings: make(map[motion.Key]Command),
	}
}

func (r *CommandRegistry) Bind(k motion.Key, cmd Command) {
	r.bindings[k] = cmd
}

func (r *CommandRegistry) Lookup(k motion.Key) (Command, bool) {
	cmd, ok := r.bindings[k]
	return cmd, ok
}
