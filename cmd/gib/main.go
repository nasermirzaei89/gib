package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/nasermirzaei89/gib"
	"github.com/spf13/cobra"
)

func main() {
	runCmd := &cobra.Command{
		Use:          "run [path]",
		Short:        "Run game from ./main.lua, directory, or .lua file",
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			target := ""
			if len(args) == 1 {
				target = args[0]
			}

			return gib.RunGame(target)
		},
	}

	rootCmd := &cobra.Command{
		Use:           "gib",
		Short:         "Lua-driven 2D game engine",
		SilenceErrors: true,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	rootCmd.AddCommand(runCmd)

	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error:", err)
		if strings.HasPrefix(err.Error(), "unknown command ") {
			rootCmd.SetOut(os.Stderr)
			_ = rootCmd.Help()
		}
		os.Exit(1)
	}
}
