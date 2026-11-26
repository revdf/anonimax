package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.com/parrotsec/packages/anonsurf/internal/anonsurf"
	"gitlab.com/parrotsec/packages/anonsurf/internal/cli"
	"gitlab.com/parrotsec/packages/anonsurf/internal/console"
)

func proxyCommand(client anonsurf.ProxyClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proxy",
		Short: "Configure upstream HTTP/HTTPS/SOCKS5 proxy",
		Long: "Configure or remove the upstream proxy used with Anonimax.\n" +
			"  anonimax proxy set socks5://127.0.0.1:1080\n" +
			"  anonimax proxy set http://user:pass@127.0.0.1:8080\n" +
			"  anonimax proxy unset",
		Args: cli.Subcommands("action"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(
		&cobra.Command{
			Use:   "set",
			Short: "Configure proxy (e.g.: socks5://127.0.0.1:1080)",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				printer := console.New(cmd.OutOrStdout(), cmd.ErrOrStderr())
				printer.Running("Saving proxy configuration...")
				if err := client.Set(args[0]); err != nil {
					return fmt.Errorf("proxy set: %w", err)
				}
				printer.Success(fmt.Sprintf("Proxy configured: %s", args[0]))
				return nil
			},
		},
		&cobra.Command{
			Use:   "unset",
			Short: "Remove proxy configuration",
			Args:  cli.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				printer := console.New(cmd.OutOrStdout(), cmd.ErrOrStderr())
				printer.Running("Removing proxy configuration...")
				if err := client.Unset(); err != nil {
					return fmt.Errorf("proxy unset: %w", err)
				}
				printer.Success("Proxy configuration removed.")
				return nil
			},
		},
	)
	return cmd
}
