package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.com/parrotsec/packages/anonsurf/internal/anonsurf"
	"gitlab.com/parrotsec/packages/anonsurf/internal/cli"
	"gitlab.com/parrotsec/packages/anonsurf/internal/console"
)

func panicCommand(client anonsurf.PanicClient) *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:   "panic",
		Short: "Emergency shutdown: kills all connections, clears traces and powers off",
		Long: "Emergency shutdown for high-risk situations.\n" +
			"Cuts network interfaces, stops Tor, flushes firewall rules, clears traces,\n" +
			"kills sensitive processes and powers the machine off.\n" +
			"Confirmation is required unless --force is set.",
		Args: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			printer := console.New(cmd.OutOrStdout(), cmd.ErrOrStderr())
			if !force {
				fmt.Fprint(cmd.OutOrStdout(), "PANIC MODE: isso vai desligar o computador imediatamente. Confirmar? [s/N]: ")
				var answer string
				if _, err := fmt.Fscanln(cmd.InOrStdin(), &answer); err != nil || (answer != "s" && answer != "S") {
					printer.Result("Cancelado.")
					return nil
				}
			}
			printer.Running("Executing panic shutdown...")
			if err := client.Run(); err != nil {
				return fmt.Errorf("panic: %w", err)
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&force, "force", false, "Skip confirmation prompt")
	return cmd
}
