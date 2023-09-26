package auth0

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/zitadel/zitadel-tools/internal/migration"
)

// Cmd represents the auth0 migration command
var Cmd = &cobra.Command{
	Use:   "auth0",
	Short: "Transform the exported Auth0 users and passwords to a ZITADEL import JSON",
	RunE: func(cmd *cobra.Command, args []string) error {
		return migrate()
	},
}

var (
	userPath     string
	passwordPath string
)

func init() {
	Cmd.Flags().StringVar(&userPath, "users", "./users.json", "path to the users.json")
	Cmd.Flags().StringVar(&passwordPath, "passwords", "./passwords.json", "path to the passwords.json")
}

type user struct {
	UserId string `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
}

type password struct {
	Oid          string `json:"oid"`
	Email        string `json:"email"`
	PasswordHash string `json:"passwordHash"`
}

func migrate() error {
	log.Printf("migrate auth0 from users(%s) and passwords(%s) into %s\n", userPath, passwordPath, migration.OutputPath)

	users, err := migration.ReadJSONLinesFile[user](userPath)
	if err != nil {
		return fmt.Errorf("read users: %w", err)
	}

	passwords, err := migration.ReadJSONLinesFile[password](passwordPath)
	if err != nil {
		return fmt.Errorf("read passwords: %w", err)
	}

	importData := migration.CreateV1Migration(createHumanUsers(users, passwords))

	err = migration.WriteProtoToFile(importData)
	if err != nil {
		return err
	}
	log.Println("Import file done")
	return nil
}

func createHumanUsers(users []user, passwords []password) []migration.User {
	result := make([]migration.User, len(users))
	for i, u := range users {
		result[i] = migration.User{
			UserId:        u.UserId,
			UserName:      u.Email,
			FirstName:     u.Name,
			LastName:      u.Name,
			Email:         u.Email,
			EmailVerified: migration.VerifiedEmails,
			PasswordHash:  getPassword(u.Email, passwords),
		}

	}
	return result
}

func getPassword(userEmail string, passwords []password) string {
	for _, p := range passwords {
		if userEmail == p.Email {
			return p.PasswordHash
		}
	}
	return ""
}
