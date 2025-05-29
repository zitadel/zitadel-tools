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
	userPath       string
	passwordPath   string
	verifiedEmails bool
)

func init() {
	Cmd.Flags().StringVar(&userPath, "users", "./users.json", "path to the users.json")
	Cmd.Flags().StringVar(&passwordPath, "passwords", "./passwords.json", "path to the passwords.json")
	Cmd.Flags().BoolVar(&verifiedEmails, "email-verified", true, "specify if imported emails are automatically verified")
}

type user struct {
	UserId        string `json:"user_id"`        // mandatory
	Email         string `json:"email"`          // mandatory
	Name          string `json:"name"`           // optional, maps to displayName
	Username      string `json:"username"`       // optional
	GivenName     string `json:"given_name"`     // optional
	FamilyName    string `json:"family_name"`    // optional
	Nickname      string `json:"nickname"`       // optional
	Locale        string `json:"locale"`         // optional
	PhoneNumber   string `json:"phone_number"`   // optional
	PhoneVerified bool   `json:"phone_verified"` // optional
	EmailVerified bool   `json:"email_verified"` // optional
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
		// Use username if available, otherwise fall back to email
		userName := u.Username
		if userName == "" {
			userName = u.Email
		}
		
		// Ensure firstName and lastName are always populated (required by ZITADEL)
		firstName := u.GivenName
		lastName := u.FamilyName
		
		// If given_name or family_name are missing, use fallbacks
		if firstName == "" {
			if u.Name != "" {
				firstName = u.Name
			} else {
				// Ultimate fallback: derive from email or username
				firstName = userName
			}
		}
		if lastName == "" {
			if u.Name != "" {
				lastName = u.Name
			} else {
				// Ultimate fallback: derive from email or username
				lastName = userName
			}
		}
		
		// Ensure lastName is never empty (ZITADEL requirement)
		if lastName == "" {
			lastName = firstName // Use firstName as fallback
		}
		
		result[i] = migration.User{
			UserId:        u.UserId,
			UserName:      userName,
			FirstName:     firstName,
			LastName:      lastName,
			Email:         u.Email,
			EmailVerified: u.EmailVerified || verifiedEmails, // Use field value or CLI flag
			PasswordHash:  getPassword(u.Email, passwords),
			Nickname:      u.Nickname,
			Name:          u.Name,
			Locale:        u.Locale,
			PhoneNumber:   u.PhoneNumber,
			PhoneVerified: u.PhoneVerified,
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
