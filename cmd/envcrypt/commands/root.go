package commands

import (
	"github.com/spf13/cobra"
)

// Root is the top-level cobra command for envcrypt.
var Root = &cobra.Command{
	Use:   "envcrypt",
	Short: "Lightweight .env file encryption and sharing tool",
	Long: `envcrypt encrypts your .env files with age encryption so they can be
safely committed to version control or shared with team members.

Use 'envcrypt init' to create an encrypted vault from an existing .env file,
'envcrypt lock' to re-encrypt after editing, 'envcrypt unlock' to decrypt for
local use, 'envcrypt share' to add recipients, and 'envcrypt rotate' to cycle
the encryption key.`,
}

func init() {
	Root.AddCommand(
		newInitCmd(),
		newLockCmd(),
		newUnlockCmd(),
		newShareCmd(),
		newRotateCmd(),
	)
}
