package main

import (
	"github.com/exiaohu/go-demo/cmd/commands"
)

var (
	// 编译时注入
	GitCommit = "unknown"
	BuildTime = "unknown"
)

func main() {
	commands.Execute(GitCommit, BuildTime)
}
