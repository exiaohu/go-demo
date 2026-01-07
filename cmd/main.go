package main

import (
	"github.com/exiaohu/go-demo/cmd/commands"
)

// @title Go Demo API
// @version 1.0
// @description This is a sample server for Go Demo.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

var (
	// 编译时注入
	GitCommit = "unknown"
	BuildTime = "unknown"
)

func main() {
	commands.Execute(GitCommit, BuildTime)
}
