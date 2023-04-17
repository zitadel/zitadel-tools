package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/admin"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
	v1 "github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	ORG_ID    = "189964076015681793"
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

func main() {
	users, err := ReadAuth0Users("users.json")
	if err != nil {
		fmt.Printf("ERROR: %v", err)
		return
	}

	passwords, pwerr := ReadAuth0UPasswords("passwords.json")
	if err != nil {
		fmt.Printf("ERROR: %v", pwerr)
		return
	}

	importData := CreateZITADELMigration(ORG_ID, users, passwords)

	err = WriteProtoToFile("importBody.json", importData)
	if err != nil {
		fmt.Printf("ERROR: %v", pwerr)
		return
	}
	fmt.Println("Import file done")
}

func ReadAuth0Users(filename string) ([]User, error) {
	file, fileScanner, err := ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var result []User
	for fileScanner.Scan() {
		data := User{}
		err = json.Unmarshal(fileScanner.Bytes(), &data)
		if err != nil {
			return nil, err
		}
		result = append(result, data)
	}
	file.Close()
	return result, nil
}

func ReadAuth0UPasswords(filename string) ([]Password, error) {
	file, fileScanner, err := ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var result []Password
	for fileScanner.Scan() {
		data := Password{}
		err = json.Unmarshal(fileScanner.Bytes(), &data)
		if err != nil {
			return nil, err
		}
		result = append(result, data)
	}
	file.Close()
	return result, nil
}

func ReadFile(filename string) (*os.File, *bufio.Scanner, error) {
	readFile, err := os.Open(filename)

	if err != nil {
		fmt.Println(err)
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

	jsonpb := &runtime.JSONPb{
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
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
		HumanUsers: createHumanUsers(true, users, passwords),
	}
	return []*admin.DataOrg{org}
}

func createHumanUsers(emailVerified bool, users []User, passwords []Password) []*v1.DataHumanUser {
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
					IsEmailVerified: emailVerified,
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
