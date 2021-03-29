# shellgo

shellgo是用来快速开发命令交互软件的工具库, 支持自动补全功能

## 安装
```go get -u github.com/cuishu/shellgo```

## 使用:

**可参考example/main.go**

使用下面的示例即可运行shell, 如果不想使用系统命令, 修改:```UseSysCmd: false```
```go
import "github.com/cuishu/shellgo"

type Prompt struct {}

func (prompt *Prompt)String() string {return ""}

func main() {
	shellgo.Run(shellgo.Config{
		UseSysCmd: true,
		Prompt:    &Prompt{},
	})
}
```

**添加自定义命令**

命令需要实现 Command 接口

命令被调用时执行 Call 方法, 执行 ```help [命令]```时调用 Help 方法, 按 Tab 键使用自动补全时, 调用 AutoComplete 方法

```go
type Command interface {
	Call([]string) string
	Help() string
	AutoComplete(line []rune, pos int) (newLine [][]rune, length int)
}
```

然后注册到 BuiltinCmd 里

```go
	env.AddBuiltinCmd("help", &Help{})
```

shellgo 只实现了 help 和 exit, 其他命令都需自行实现