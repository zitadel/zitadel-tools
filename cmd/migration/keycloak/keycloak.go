package keycloak

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/zitadel/zitadel-tools/internal/migration"
)

// Cmd represents the auth0 migration command
var Cmd = &cobra.Command{
	Use:   "keycloak",
	Short: "Transform the exported Keycloak users and passwords to a ZITADEL import JSON",
	RunE: func(cmd *cobra.Command, args []string) error {
		return migrate()
	},
}

var (
	realmPath string
)

func init() {
	Cmd.Flags().StringVar(&realmPath, "realm", "./realm.json", "realm export in json format")
}

func migrate() error {
	realm, err := migration.ReadJSONFile[realm](realmPath)
	if err != nil {
		return fmt.Errorf("read realm: %w", err)
	}
	users, err := createHumanUsers(realm.Users)
	if err != nil {
		return err
	}

	importData := migration.CreateV1Migration(users)
	err = migration.WriteProtoToFile(importData)
	if err != nil {
		return err
	}
	log.Println("Import file done")
	return nil
}

// Currently ignored fields:
// - CreatedTimestamp
// - Enabled
// - Totp (bool flag and credential type)
// - DisableableCredentialTypes
// - RequiredActions
// - RealmRoles
// - NotBefore
// - Groups
//
// also note that credentials seems to be able to contain more
// than just passwords.
func createHumanUsers(users []user) ([]migration.User, error) {
	result := make([]migration.User, len(users))
	for i, u := range users {
		password, err := u.getPassword()
		if err != nil {
			return nil, fmt.Errorf("create users[%d] ID %q: %w", i, u.ID, err)
		}

		result[i] = migration.User{
			UserId:        u.ID,
			UserName:      u.Username,
			FirstName:     u.FirstName,
			LastName:      u.LastName,
			Email:         u.Email,
			EmailVerified: u.EmailVerified,
			PasswordHash:  password,
		}

	}
	return result, nil
}
