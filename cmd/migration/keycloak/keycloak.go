package keycloak

import (
	"github.com/spf13/cobra"
)

// Cmd represents the auth0 migration command
var Cmd = &cobra.Command{
	Use:   "auth0",
	Short: "Transform the exported Keycloak users and passwords to a ZITADEL import JSON",
}
