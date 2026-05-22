package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourorg/envcrypt/internal/vault"
)

func newInitCmd() *cobra.Command {
	var force bool
	var envPath string
	var vaultPath string

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new encrypted vault from a .env file",
		RunE: func(cmd *cobra.Command, args []string) error {
			recip, err := vault.Init(envPath, vaultPath, force)
			if err != nil {
				return fmt.Errorf("init failed: %w", err)
			}
			fmt.Printf("Vault initialised at %s\n", vaultPath)
			fmt.Printf("Your public key (share with teammates):\n  %s\n", recip)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Overwrite existing vault")
	cmd.Flags().StringVarP(&envPath, "env", "e", ".env", "Path to source .env file")
	cmd.Flags().StringVarP(&vaultPath, "vault", "v", ".env.age", "Path to output vault file")

	return cmd
}
