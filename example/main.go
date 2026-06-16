package main

import (
	"fmt"

	"github.com/cuishu/shellgo"
)

type Prompt struct {
	env *shellgo.Env
}

var count int

func (prompt *Prompt) String() string {
	count++
	var errmesg string
	if prompt.env.ErrMesg != "" {
		errmesg = fmt.Sprintf(" \033[0;31m[%s]\033[0m", prompt.env.ErrMesg)
	}
	return fmt.Sprintf("\033[0;35mshell-go \033[1;32m%d\033[0m%s$ ", count, errmesg)
}

func main() {
	shell := shellgo.NewShell()
	shell.SetConfig(shellgo.Config{
		UseSysCmd: true,
		Prompt:    &Prompt{shell.Env()},
	})
	shell.Run()
}
