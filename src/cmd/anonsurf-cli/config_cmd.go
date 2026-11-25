package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.com/parrotsec/packages/anonsurf/internal/config"
)

func configCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Read or write Anonimax settings",
		Args:  cobra.NoArgs,
		RunE:  func(cmd *cobra.Command, args []string) error { return cmd.Help() },
	}
	cmd.AddCommand(configListCommand(), configSetCommand())
	return cmd
}

func configListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Print all settings as key=value",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(config.SystemPath)
			if err != nil {
				cfg = config.Default()
			}
			fields := []struct{ key, value string }{
				{"change_mac", fmt.Sprintf("%v", cfg.ChangeMac)},
				{"renew_ip", fmt.Sprintf("%v", cfg.RenewIP)},
				{"force_utc", fmt.Sprintf("%v", cfg.ForceUTC)},
				{"block_webrtc", fmt.Sprintf("%v", cfg.BlockWebRTC)},
				{"block_telemetry", fmt.Sprintf("%v", cfg.BlockTelemetry)},
				{"disable_camera_mic", fmt.Sprintf("%v", cfg.DisableCameraMic)},
				{"clean_on_stop", fmt.Sprintf("%v", cfg.CleanOnStop)},
				{"use_vpn", fmt.Sprintf("%v", cfg.UseVPN)},
				{"vpn_config_path", cfg.VPNConfigPath},
				{"use_proxy", fmt.Sprintf("%v", cfg.UseProxy)},
				{"socks5_proxy_address", cfg.Socks5ProxyAddress},
			}
			for _, f := range fields {
				fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", f.key, f.value)
			}
			return nil
		},
	}
}

func configSetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Update a single setting in the config file (requires root)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return config.UpdateKeys(config.SystemPath, map[string]string{args[0]: args[1]})
		},
	}
}
