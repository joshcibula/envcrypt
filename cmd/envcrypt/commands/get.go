package commands

import (
	"fmt"

	"github.com/nicholasgasior/envcrypt/internal/config"
	"github.com/nicholasgasior/envcrypt/internal/vault"
	"github.com/spf13/cobra"
)

func newGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <key>",
		Short: "Get the value of a key from the encrypted vault",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.LoadOrDefault("")
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			key := args[0]

			val, err := vault.Get(cfg.VaultPath, cfg.KeyPath, key)
			if err != nil {
				return fmt.Errorf("get %q: %w", key, err)
			}

			fmt.Fprintln(cmd.OutOrStdout(), val)
			return nil
		},
	}
	return cmd
}

func init() {
	rootCmd.AddCommand(newGetCmd())
}
