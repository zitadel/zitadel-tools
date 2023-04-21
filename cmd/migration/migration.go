package migration

import (
	"github.com/spf13/cobra"

	"github.com/zitadel/zitadel-tools/cmd/migration/auth0"
)

// Cmd represents the migration root command
var Cmd = &cobra.Command{
	Use:   "migrate",
	Short: "Transform data from other providers (like Auth0) to ZITADEL import data",
}

func init() {
	Cmd.AddCommand(auth0.Cmd)
}
