package main

import (
	"github.com/qkja/go-zero/core/load"
	"github.com/qkja/go-zero/core/logx"
	"github.com/qkja/go-zero/tools/goctl/cmd"
)

func main() {
	logx.Disable()
	load.Disable()
	cmd.Execute()
}
