package shellgo

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
	Prompt    Prompt
}

type Env struct {
	ErrMesg    string
	BuiltinCmd map[string]Command
}

func (env *Env) AddBuiltinCmd(name string, cmd Command) {
	env.BuiltinCmd[name] = cmd
}

func GetEnv() *Env {
	return env
}

func init() {
	env = &Env{}
	env.BuiltinCmd = make(map[string]Command)
	env.AddBuiltinCmd("exit", &Exit{})
	env.AddBuiltinCmd("help", &Help{Env: env})
}
