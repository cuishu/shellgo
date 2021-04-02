package shellgo

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/chzyer/readline"
)

var env *Env

func doExec(conf *Config, slice []string) {
	env.ErrMesg = ""
	cmd, ok := env.BuiltinCmd[slice[0]]
	if !ok {
		if conf.UseSysCmd {
			cmd := exec.Command(slice[0], slice[1:]...)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				env.ErrMesg = err.Error()
			}
		} else {
			env.ErrMesg = "1"
			fmt.Printf("shell-go: %s: command not fond\n", slice[0])
		}
		return
	}
	env.ErrMesg = cmd.Call(slice[1:])
}

type completer struct{}

func (c *completer) Do(line []rune, pos int) (newLine [][]rune, length int) {
	length = pos
	for k, cmd := range env.BuiltinCmd {
		if strings.Index(string(line), k+" ") == 0 {
			s := []rune(strings.TrimLeft(string(line[len(k+" "):]), " "))
			newLine, l := cmd.AutoComplete(s, len(s))
			return newLine, pos + l
		}
		if string(line) == k {
			return [][]rune{[]rune(" ")}, pos
		}
		if n := strings.Index(k, string(line)); n == 0 {
			newLine = append(newLine, []rune(k[pos:]+" "))
		}
	}

	return
}

func Run(conf Config) {
	r, err := readline.NewEx(&readline.Config{
		InterruptPrompt:   "^C",
		EOFPrompt:         "exit",
		HistorySearchFold: true,
		AutoComplete:      &completer{},
	})
	if err != nil {
		panic(err)
	}
	for {
		r.SetPrompt(conf.Prompt.String())

		if line, err := r.Readline(); err != nil {
			if err != readline.ErrInterrupt {
				os.Exit(0)
			}
		} else {
			var slice []string
			for _, s := range strings.Split(line, " ") {
				if s != "" {
					slice = append(slice, s)
				}
			}
			if len(slice) == 0 {
				continue
			}
			doExec(&conf, slice)
		}
	}
}
