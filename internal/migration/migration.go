// Package migration provides common tools for the data migrations.
package migration

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"
	v1 "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

var (
	OutputPath     string
	OrganizationID string
	Timeout        time.Duration
	MultiLine      bool
)

func ReadJSONFile[T any](name string) (out T, err error) {
	file, err := os.Open(name)
	if err != nil {
		return out, fmt.Errorf("json file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&out)
	if err != nil {
		return out, fmt.Errorf("json file: %w", err)
	}
	return out, nil
}

func ReadJSONLinesFile[T any](name string) (out []T, err error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("json lines file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	for {
		var dst T
		err = decoder.Decode(&dst)
		if errors.Is(err, io.EOF) {
			return out, nil
		} else if err != nil {
			return nil, fmt.Errorf("json lines file: %w", err)
		}
		out = append(out, dst)
	}
}

type User struct {
	UserId        string // ignored in V1
	UserName      string
	FirstName     string
	LastName      string
	Email         string
	EmailVerified bool
	PasswordHash  string
	// Additional optional Auth0 fields
	Nickname      string // maps to nickName
	Name          string // maps to displayName
	Locale        string // maps to preferredLanguage
	PhoneNumber   string // maps to phone
	PhoneVerified bool   // maps to isPhoneVerified
}

func CreateV1Migration(users []User) *admin.ImportDataRequest {
	importDataOrg := &admin.ImportDataOrg{
		Orgs: createOrgs(OrganizationID, users),
	}
	importData := &admin.ImportDataRequest{
		Timeout: Timeout.String(),
		Data: &admin.ImportDataRequest_DataOrgs{
			DataOrgs: importDataOrg,
		},
	}

	return importData
}

func createOrgs(id string, users []User) []*admin.DataOrg {
	org := &admin.DataOrg{
		OrgId:      id,
		HumanUsers: createHumanUsers(users),
	}
	return []*admin.DataOrg{org}
}

func createHumanUsers(users []User) []*v1.DataHumanUser {
	result := make([]*v1.DataHumanUser, len(users))
	for i, u := range users {
		result[i] = &v1.DataHumanUser{
			UserId: u.UserId,
			User: &management.ImportHumanUserRequest{
				UserName: u.UserName,
				Profile: &management.ImportHumanUserRequest_Profile{
					FirstName:         u.FirstName,
					LastName:          u.LastName,
					NickName:          u.Nickname,
					DisplayName:       u.Name,
					PreferredLanguage: u.Locale,
				},
				Email: &management.ImportHumanUserRequest_Email{
					Email:           u.Email,
					IsEmailVerified: u.EmailVerified,
				},
			},
		}
		
		// Add phone if present
		if u.PhoneNumber != "" {
			result[i].User.Phone = &management.ImportHumanUserRequest_Phone{
				Phone:           u.PhoneNumber,
				IsPhoneVerified: u.PhoneVerified,
			}
		}
		
		if u.PasswordHash != "" {
			result[i].User.HashedPassword = &management.ImportHumanUserRequest_HashedPassword{
				Value: u.PasswordHash,
			}
		}
	}
	return result
}

func WriteProtoToFile(importData *admin.ImportDataRequest) error {
	opts := protojson.MarshalOptions{
		Multiline: MultiLine,
	}
	encodedData, err := opts.Marshal(importData)
	if err != nil {
		return err
	}

	return os.WriteFile(OutputPath, encodedData, 0666)
}
