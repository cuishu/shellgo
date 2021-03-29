package shellgo

import (
	"os"
)

type Exit struct{}

func (exit *Exit) Call(args []string) string {
	os.Exit(0)
	return ""
}

func (exit *Exit) Help() string {
	return "exit: exit\n" +
		"\t退出 shell"
}

func (exit *Exit) AutoComplete(line []rune, pos int) (newLine [][]rune, length int) {
	return
}
