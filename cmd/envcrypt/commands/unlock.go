package commands

import (
	"fmt"

	"github.com/nicholasgasior/envcrypt/internal/config"
	"github.com/nicholasgasior/envcrypt/internal/vault"
	"github.com/spf13/cobra"
)

func newUnlockCmd() *cobra.Command {
	var cfgPath string
	var outputFile string

	cmd := &cobra.Command{
		Use:   "unlock",
		Short: "Decrypt the vault and write the .env file",
		Long: `Decrypt the vault file using the local age identity and write the
plaintext environment variables to the configured .env file (or a custom path).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var cfg *config.Config
			var err error
			if cfgPath != "" {
				cfg, err = config.Load(cfgPath)
			} else {
				cfg, err = config.LoadOrDefault("")
			}
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			opts := vault.UnlockOptions{
				OutputFile: outputFile,
			}

			if err := vault.Unlock(cfg, opts); err != nil {
				return fmt.Errorf("unlock: %w", err)
			}

			dest := cfg.EnvFile
			if outputFile != "" {
				dest = outputFile
			}
			fmt.Printf("Vault decrypted → %s\n", dest)
			return nil
		},
	}

	cmd.Flags().StringVar(&cfgPath, "config", "", "Path to envcrypt config file (default: .envcrypt.yaml)")
	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "Write decrypted .env to this path instead of the configured default")

	return cmd
}
