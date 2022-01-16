package cli

import (
	"github.com/spf13/cobra"
	"go-scaffold/internal/app/cli/command/greet"
	"go-scaffold/internal/app/global"
)

func Setup() error {
	greetHandler := greet.NewHandler()

	greetCommand := &cobra.Command{Use: "greet", Run: greetHandler.Default}
	greetCommand.AddCommand(&cobra.Command{Use: "to", Run: greetHandler.To})

	global.Command().AddCommand(
		greetCommand,
	)

	return nil
}
