package auth0

import (
	"bufio"
	"encoding/json"
	"log"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/cobra"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/admin"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
	v1 "github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

// Cmd represents the auth0 migration command
var Cmd = &cobra.Command{
	Use:   "auth0",
	Short: "Transform the exported Auth0 users and passwords to a ZITADEL import JSON",
	Run: func(cmd *cobra.Command, args []string) {
		migrate()
	},
}

var (
	userPath       string
	passwordPath   string
	outputPath     string
	organisationID string
	verifiedEmails bool
	multiLine      bool
)

func init() {
	Cmd.Flags().StringVar(&userPath, "users", "./users.json", "path to the users.json")
	Cmd.Flags().StringVar(&passwordPath, "passwords", "./passwords.json", "path to the passwords.json")
	Cmd.Flags().StringVar(&outputPath, "output", "./importBody.json", "path where the generated json will be saved")
	Cmd.Flags().StringVar(&organisationID, "org", "", "id of the ZITADEL organisation, where the users will be imported")
	Cmd.Flags().BoolVar(&verifiedEmails, "email-verified", true, "specify if imported emails are automatically verified (default)")
	Cmd.Flags().BoolVar(&multiLine, "multiline", false, "print the JSON output in multiple lines")
}

const (
	Algorithm = "bcrypt"
)

type User struct {
	UserId string `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
}

type Password struct {
	Oid          string `json:"oid"`
	Email        string `json:"email"`
	PasswordHash string `json:"passwordHash"`
}

func migrate() {
	if organisationID == "" {
		log.Fatal("Please provide the organisation id")
	}
	log.Printf("migrate auth0 from users(%s) and passwords(%s) into %s\n", userPath, passwordPath, outputPath)

	users, err := ReadAuth0Users(userPath)
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}

	passwords, pwerr := ReadAuth0UPasswords(passwordPath)
	if err != nil {
		log.Fatalf("ERROR: %v", pwerr)
	}

	importData := CreateZITADELMigration(organisationID, users, passwords)

	err = WriteProtoToFile(outputPath, importData)
	if err != nil {
		log.Fatalf("ERROR: %v", pwerr)
	}
	log.Println("Import file done")
}

func ReadAuth0Users(filename string) ([]User, error) {
	file, fileScanner, err := ReadFile(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var result []User
	for fileScanner.Scan() {
		data := User{}
		err = json.Unmarshal(fileScanner.Bytes(), &data)
		if err != nil {
			return nil, err
		}
		result = append(result, data)
	}
	return result, nil
}

func ReadAuth0UPasswords(filename string) ([]Password, error) {
	file, fileScanner, err := ReadFile(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var result []Password
	for fileScanner.Scan() {
		data := Password{}
		err = json.Unmarshal(fileScanner.Bytes(), &data)
		if err != nil {
			return nil, err
		}
		result = append(result, data)
	}
	return result, nil
}

func ReadFile(filename string) (*os.File, *bufio.Scanner, error) {
	readFile, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	return readFile, fileScanner, nil
}

func CreateZITADELMigration(orgID string, users []User, passwords []Password) *admin.ImportDataRequest {
	importDataOrg := &admin.ImportDataOrg{
		Orgs: createOrgs(orgID, users, passwords),
	}
	importData := &admin.ImportDataRequest{
		Timeout: "30m",
		Data: &admin.ImportDataRequest_DataOrgs{
			DataOrgs: importDataOrg,
		},
	}

	return importData
}

func WriteProtoToFile(filepath string, importData *admin.ImportDataRequest) error {
	outFile, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer outFile.Close()

	jsonpb := &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			Multiline: multiLine,
		},
	}
	encodedData, err := jsonpb.Marshal(importData)
	if err != nil {
		return err
	}

	// writing the actual transaction item to the file
	if _, err := outFile.Write(encodedData); err != nil {
		return err
	}

	return nil
}

func createOrgs(id string, users []User, passwords []Password) []*admin.DataOrg {
	org := &admin.DataOrg{
		OrgId:      id,
		HumanUsers: createHumanUsers(users, passwords),
	}
	return []*admin.DataOrg{org}
}

func createHumanUsers(users []User, passwords []Password) []*v1.DataHumanUser {
	result := make([]*v1.DataHumanUser, 0)
	for _, u := range users {
		user := &v1.DataHumanUser{
			User: &management.ImportHumanUserRequest{
				UserName: u.Email,
				Profile: &management.ImportHumanUserRequest_Profile{
					FirstName: u.Name,
					LastName:  u.Name,
				},
				Email: &management.ImportHumanUserRequest_Email{
					Email:           u.Email,
					IsEmailVerified: verifiedEmails,
				},
			},
		}
		passwordHash := getPassword(u.Email, passwords)
		if passwordHash != "" {
			user.User.HashedPassword = &management.ImportHumanUserRequest_HashedPassword{
				Value:     passwordHash,
				Algorithm: Algorithm,
			}
		}
		result = append(result, user)
	}
	return result
}

func getPassword(userEmail string, passwords []Password) string {
	for _, p := range passwords {
		if userEmail == p.Email {
			return p.PasswordHash
		}
	}
	return ""
}
