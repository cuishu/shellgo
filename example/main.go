package main

import (
	"fmt"

	"github.com/cuishu/shellgo"
)

type Prompt struct{}

var count int

func (prompt *Prompt) String() string {
	env := shellgo.GetEnv()
	count++
	var errmesg string
	if env.ErrMesg != "" {
		errmesg = fmt.Sprintf(" \033[0;31m[%s]\033[0m", env.ErrMesg)
	}
	return fmt.Sprintf("\033[0;35mshell-go \033[1;32m%d\033[0m%s$ ", count, errmesg)
}

func main() {
	shellgo.Run(shellgo.Config{
		UseSysCmd: true,
		Prompt:    &Prompt{},
	})
}
