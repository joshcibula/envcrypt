package commands

import (
	"github.com/spf13/cobra"
)

// NewEditCmdForTest returns a standalone edit command for use in tests.
func NewEditCmdForTest() *cobra.Command {
	cmd := newEditCmd()
	return cmd
}
