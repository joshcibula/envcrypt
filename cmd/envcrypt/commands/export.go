package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/envcrypt/internal/vault"
)

func newExportCmd() *cobra.Command {
	var (
		output string
		force  bool
	)

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Decrypt the vault and write plaintext env vars to a file",
		Long: `Export decrypts the vault using your local key and writes the
plaintext environment variables to an output file.

Unlike 'unlock', export targets an arbitrary path and does not
modify vault state or remove the encrypted vault file.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("get working directory: %w", err)
			}

			opts := vault.ExportOptions{
				OutputPath: output,
				Overwrite:  force,
			}

			if err := vault.Export(dir, opts); err != nil {
				return err
			}

			out := output
			if out == "" {
				out = ".env"
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Exported vault to %s\n", out)
			return nil
		},
	}

	cmd.Flags().StringVarP(&output, "output", "o", "", "Output file path (default: .env)")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Overwrite output file if it exists")

	return cmd
}
