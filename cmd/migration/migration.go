package migration

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/zitadel/zitadel-tools/cmd/migration/auth0"
	"github.com/zitadel/zitadel-tools/cmd/migration/keycloak"
	"github.com/zitadel/zitadel-tools/internal/migration"
)

// Cmd represents the migration root command
var Cmd = &cobra.Command{
	Use:   "migrate",
	Short: "Transform data from other providers (like Auth0) to ZITADEL import data",
}

func init() {
	Cmd.PersistentFlags().StringVar(&migration.OrganizationID, "org", "", "id of the ZITADEL organization, where the users will be imported")
	Cmd.MarkPersistentFlagRequired("org")

	Cmd.PersistentFlags().StringVar(&migration.OutputPath, "output", "./importBody.json", "path where the generated json will be saved")
	Cmd.PersistentFlags().DurationVar(&migration.Timeout, "timeout", 30*time.Minute, "maximum duration to be used for the import")
	Cmd.PersistentFlags().BoolVar(&migration.VerifiedEmails, "email-verified", true, "specify if imported emails are automatically verified")
	Cmd.PersistentFlags().BoolVar(&migration.MultiLine, "multiline", false, "print the JSON output in multiple lines")

	Cmd.AddCommand(auth0.Cmd)
	Cmd.AddCommand(keycloak.Cmd)
}
