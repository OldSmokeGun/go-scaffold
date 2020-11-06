package main

import (
	"gin-scaffold/kernel"
	"gin-scaffold/kernel/components"
	"gin-scaffold/kernel/db"
	"gin-scaffold/kernel/global"
)

func main() {
	if err := components.LoadConfig(); err != nil {
		panic(err)
	}

	if err := components.InitLogger(); err != nil {
		panic(err)
	}

	if err := db.Init(); err != nil {
		panic(err)
	}

	defer func() {
		if err := global.DB.Close(); err != nil {
			panic(err)
		}
	}()

	kernel.Bootstrap()
}
