package main

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"gitlab.com/parrotsec/packages/anonsurf/internal/anonsurf"
	"gitlab.com/parrotsec/packages/anonsurf/internal/cli"
	"gitlab.com/parrotsec/packages/anonsurf/internal/console"
)

func autoChangeCommand(client anonsurf.AutoChangeClient) *cobra.Command {
	return &cobra.Command{
		Use:   "autochange",
		Short: "Auto-rotate Tor identity every N minutes",
		Long: "Start or stop automatic Tor identity rotation.\n" +
			"  anonsurf autochange <N>     rotate every N minutes (N >= 1)\n" +
			"  anonsurf autochange stop    stop automatic rotation",
		Args: cobra.ExactArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) == 0 {
				return []string{"stop"}, cobra.ShellCompDirectiveNoFileComp
			}
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			printer := console.New(cmd.OutOrStdout(), cmd.ErrOrStderr())
			if args[0] == "stop" {
				printer.Running("Stopping auto identity change...")
				if err := client.Stop(); err != nil {
					return fmt.Errorf("autochange stop: %w", err)
				}
				printer.Success("Auto identity change stopped.")
				return nil
			}
			minutes, err := strconv.Atoi(args[0])
			if err != nil || minutes < 1 {
				return cli.NewError(
					"interval must be a number >= 1",
					"Example: anonsurf autochange 10",
				)
			}
			printer.Running(fmt.Sprintf("Starting auto identity change every %d minute(s)...", minutes))
			if err := client.Start(minutes); err != nil {
				return fmt.Errorf("autochange: %w", err)
			}
			printer.Success(fmt.Sprintf("Auto identity change every %d minute(s).", minutes))
			return nil
		},
	}
}
