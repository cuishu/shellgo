package shellgo

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/chzyer/readline"
)

type Shell struct {
	env  *Env
	conf Config
}

func NewShell() *Shell {
	env := &Env{}
	env.interrupt = make(chan os.Signal)
	env.BuiltinCmd = make(map[string]Command)
	env.AddBuiltinCmd("help", &Help{Env: env})
	return &Shell{env: env}
}

func (s *Shell) AddBuiltinCmd(name string, cmd Command) {
	s.env.AddBuiltinCmd(name, cmd)
}

func (s *Shell) Env() *Env {
	return s.env
}

func (s *Shell) doExec(slice []string) {
	if slice[0] == "exit" {
		os.Exit(0)
	}
	s.env.ErrMesg = ""
	cmd, ok := s.env.BuiltinCmd[slice[0]]
	if !ok {
		if s.conf.UseSysCmd {
			cmd := exec.Command(slice[0], slice[1:]...)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				s.env.ErrMesg = err.Error()
			}
		} else {
			s.env.ErrMesg = "1"
			fmt.Printf("shell-go: %s: command not fond\n", slice[0])
		}
		return
	}
	if s.conf.ForkCmd {
		s.env.cid, _, _ = syscall.RawSyscall(syscall.SYS_CLONE, 0, 0, 0)
		if s.env.cid == 0 {
			os.WriteFile(fmt.Sprintf("%s/%d.out", os.TempDir(), os.Getpid()), []byte(cmd.Call(slice[1:])), 0600)
			os.Exit(0)
		} else {
			var ws syscall.WaitStatus
			syscall.Wait4(int(s.env.cid), &ws, 0, nil)
			filename := fmt.Sprintf("%s/%d.out", os.TempDir(), s.env.cid)
			data, err := os.ReadFile(filename)
			if err == nil {
				os.Remove(filename)
				s.env.ErrMesg = string(data)
			}
			s.env.cid = 0
		}
	} else {
		s.env.ErrMesg = cmd.Call(slice[1:])
	}
}

type completer struct {
	env *Env
}

func (c *completer) Do(line []rune, pos int) (newLine [][]rune, length int) {
	length = pos
	for k, cmd := range c.env.BuiltinCmd {
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

func (s *Shell) SetConfig(conf Config) {
	s.conf = conf
}

func (s *Shell) Run() {
	signal.Notify(s.env.interrupt, os.Interrupt)
	go func() {
		for {
			<-s.env.interrupt
			if s.env.cid != 0 {
				syscall.Kill(int(s.env.cid), syscall.SIGKILL)
			}
		}
	}()
	r, err := readline.NewEx(&readline.Config{
		InterruptPrompt:   "^C",
		EOFPrompt:         "exit",
		HistorySearchFold: true,
		AutoComplete:      &completer{env: s.env},
	})
	if err != nil {
		panic(err)
	}
	for {
		r.SetPrompt(s.conf.Prompt.String())

		if line, err := r.Readline(); err != nil {
			if err != readline.ErrInterrupt {
				os.Exit(0)
			} else {
				select {
				case s.env.interrupt <- os.Kill:
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
			s.doExec(slice)
		}
	}
}
