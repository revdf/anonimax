package main

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"gitlab.com/parrotsec/packages/anonsurf/internal/anonsurf"
	"gitlab.com/parrotsec/packages/anonsurf/internal/cli"
	"gitlab.com/parrotsec/packages/anonsurf/internal/command"
	"gitlab.com/parrotsec/packages/anonsurf/internal/console"
	"gitlab.com/parrotsec/packages/anonsurf/internal/netcheck"
	"gitlab.com/parrotsec/packages/anonsurf/internal/safekill"
	"gitlab.com/parrotsec/packages/anonsurf/internal/status"
	"gitlab.com/parrotsec/packages/anonsurf/internal/systemd"
)

// version is overridden at release time via -ldflags "-X main.version=..." (see Makefile).
var version = "6.0.0"

func main() {
	if err := execute(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		printer := console.New(os.Stdout, os.Stderr)
		printer.Error(fmt.Sprintf("anonsurf: %v", err))
		if hint := cli.Hint(err); hint != "" {
			printer.Info(hint)
		}
		os.Exit(1)
	}
}

func execute(args []string, stdout, stderr io.Writer) error {
	rootCmd := newRootCommand(stdout, stderr)
	rootCmd.SetArgs(args)
	return rootCmd.Execute()
}

func newRootCommand(stdout, stderr io.Writer) *cobra.Command {
	cobra.EnableCommandSorting = false
	var verbose bool
	rootCmd := &cobra.Command{
		Use:           "anonimax",
		Short:         "Manage Anonimax system-wide Tor routing",
		Long:          "Anonimax is a system-wide Tor routing and privacy tool.",
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cli.Subcommands("command"),
		Example: `  anonimax start
  anonimax status
  anonimax changeid
  anonimax autochange 10
  anonimax proxy set socks5://127.0.0.1:1080
  anonimax panic
  anonimax boot enable`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if verbose {
				command.SetVerboseOutput(stderr)
				return
			}
			command.SetVerboseOutput(nil)
		},
	}

	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.SetHelpTemplate(helpTemplate)
	rootCmd.SetVersionTemplate("Anonimax {{.Version}}\n")
	rootCmd.SetFlagErrorFunc(cli.FlagError)
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Show commands and their output")

	service := systemd.NewAnonsurfdController()
	monitor := anonsurf.MonitorClient{}
	safeKill := safekill.Runner{}
	rootCmd.AddCommand(
		serviceCommandWithSafeKill(
			"start", "Route traffic through Tor", "Starting Anonimax...",
			service.Start, service.StatusText, "Anonimax started.", safeKill,
		),
		serviceCommandWithSafeKill(
			"stop", "Restore the original network", "Restoring the original network...",
			service.Stop, service.StatusText, "Anonimax stopped.", safeKill,
		),
		serviceCommandWithSafeKill(
			"restart", "Rebuild Tor routing", "Restarting Anonimax...",
			service.Restart, service.StatusText, "Anonimax restarted.", safeKill,
		),
		statusCommand(status.Reader{}, monitor),
		changeIDCommand(anonsurf.IdentityClient{}),
		autoChangeCommand(anonsurf.AutoChangeClient{}),
		proxyCommand(anonsurf.ProxyClient{}),
		configCommand(),
		ipCommand(netcheck.Checker{}),
		monitorCommand(monitor),
		panicCommand(anonsurf.PanicClient{}),
		bootCommand(service),
	)
	configureRootHelp(rootCmd)

	return rootCmd
}

const helpTemplate = `Anonimax

{{if .Long}}{{.Long}}{{else}}{{.Short}}{{end}}

Usage:
  {{.UseLine}}
{{if .HasAvailableSubCommands}}
Commands:
{{range .Commands}}{{if or .IsAvailableCommand (eq .Name "help")}}  {{rpad .Name .NamePadding }} {{.Short}}
{{end}}{{end}}{{end}}
{{if .HasAvailableFlags}}Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}
{{end}}{{if .HasAvailableInheritedFlags}}
Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}
{{end}}{{if .Example}}
Examples:
{{.Example}}
{{end}}`
