package shellgo

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/chzyer/readline"
)

var env *Env

func doExec(conf *Config, slice []string) {
	if slice[0] == "exit" {
		os.Exit(0)
	}
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
	env.cid, _, _ = syscall.RawSyscall(syscall.SYS_FORK, 0, 0, 0)
	if env.cid == 0 {
		os.Exit(cmd.Call(slice[1:]))
	} else {
		var ws syscall.WaitStatus
		syscall.Wait4(int(env.cid), &ws, 0, nil)
		if ws.ExitStatus() != 0 {
			env.ErrMesg = strconv.FormatInt(int64(ws.ExitStatus()), 10)
		}
		env.cid = 0
	}
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
	signal.Notify(env.interrupt, os.Interrupt)
	go func() {
		for {
			<-env.interrupt
			if env.cid != 0 {
				syscall.Kill(int(env.cid), syscall.SIGKILL)
			}
		}
	}()
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
			} else {
				select {
				case env.interrupt <- os.Kill:
				default:
				}
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
