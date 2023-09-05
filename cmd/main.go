package main

import (
	klog "github.com/go-kratos/kratos/v2/log"

	"go-scaffold/internal/command"
	_ "go-scaffold/pkg/dump"
	"go-scaffold/pkg/log"
	iklog "go-scaffold/pkg/log/kratos"
)

func main() {
	klog.SetLogger(iklog.NewLogger(log.NewNop()))

	cmd := command.NewCommand()
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
