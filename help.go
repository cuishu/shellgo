package shellgo

import (
	"fmt"
	"strings"
	"time"
)

type Help struct {
	Env            *Env
	recentCallTime time.Time
}

func (help *Help) Call(args []string) int {
	defer func() { help.recentCallTime = time.Now() }()
	if help.recentCallTime.Add(time.Second).After(time.Now()) && len(args) == 0 {
		fmt.Println(`|   \_____/   |
/  |\/     \/|  \
\_/ | /\ /\ | \_/
    |_\/ \/_|
   /   \o/   \
   \___/"\___/`)
		fmt.Println("You really need help")
		return 0
	}
	if len(args) == 0 {
		for _, v := range help.Env.BuiltinCmd {
			msg := v.Help()
			if msg != "" {
				fmt.Println(msg)
			}
		}
	} else {
		cmd, ok := help.Env.BuiltinCmd[args[0]]
		if !ok {
			fmt.Println(args[0], "not found")
			return 1
		}
		fmt.Println(cmd.Help())
	}
	return 0
}

func (help *Help) Help() string {
	return "help: help [命令名]\n" +
		"\t打印帮助信息"
}

func (help *Help) AutoComplete(line []rune, pos int) (newLine [][]rune, length int) {
	for k := range help.Env.BuiltinCmd {
		if n := strings.Index(k, string(line)); n == 0 {
			newLine = append(newLine, []rune(k[pos:]+" "))
		}
	}
	return
}
