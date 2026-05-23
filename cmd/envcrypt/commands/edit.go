package commands

import (
	"fmt"
	"os"

	"github.com/nicholasgasior/envcrypt/internal/config"
	"github.com/nicholasgasior/envcrypt/internal/vault"
	"github.com/spf13/cobra"
)

func newEditCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit KEY VALUE",
		Short: "Edit or add a key-value pair in the encrypted vault",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.LoadOrDefault("")
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			key := args[0]
			value := args[1]

			if err := vault.Edit(cfg.VaultPath, cfg.KeyPath, key, value); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				return err
			}

			fmt.Printf("Key %q updated in vault.\n", key)
			return nil
		},
	}
	return cmd
}

func init() {
	rootCmd.AddCommand(newEditCmd())
}
