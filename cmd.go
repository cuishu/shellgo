package shellgo

import "os"

type Command interface {
	Call([]string) string
	Help() string
	AutoComplete(line []rune, pos int) (newLine [][]rune, length int)
}

type Prompt interface {
	String() string
}

type Config struct {
	UseSysCmd bool
	ForkCmd   bool
	Prompt    Prompt
}

type Env struct {
	ErrMesg    string
	BuiltinCmd map[string]Command
	interrupt  chan os.Signal
	cid        uintptr
}

func (env *Env) AddBuiltinCmd(name string, cmd Command) {
	env.BuiltinCmd[name] = cmd
}
