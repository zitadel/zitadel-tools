package keycloak

import (
	"encoding/json"
	"fmt"

	"github.com/zitadel/passwap/pbkdf2"
)

/*
 Models are generated using https://mholt.github.io/json-to-go/.
 Some fields have type `any` because they where empty in the example data.
 There is a big chance these need to be changed if we start to need those.
*/

type realm struct {
	Realm string `json:"realm,omitempty"`
	Users []user `json:"users,omitempty"`
}

type user struct {
	ID                         string       `json:"id,omitempty"`
	CreatedTimestamp           int64        `json:"createdTimestamp,omitempty"`
	Username                   string       `json:"username,omitempty"`
	Enabled                    bool         `json:"enabled,omitempty"`
	Totp                       bool         `json:"totp,omitempty"`
	EmailVerified              bool         `json:"emailVerified,omitempty"`
	FirstName                  string       `json:"firstName,omitempty"`
	LastName                   string       `json:"lastName,omitempty"`
	Email                      string       `json:"email,omitempty"`
	Credentials                []credential `json:"credentials,omitempty"`
	DisableableCredentialTypes []any        `json:"disableableCredentialTypes,omitempty"`
	RequiredActions            []any        `json:"requiredActions,omitempty"`
	RealmRoles                 []string     `json:"realmRoles,omitempty"`
	NotBefore                  int          `json:"notBefore,omitempty"`
	Groups                     []any        `json:"groups,omitempty"`
}

type credential struct {
	ID             string `json:"id,omitempty"`
	Type           string `json:"type,omitempty"`
	UserLabel      string `json:"userLabel,omitempty"`
	CreatedDate    int64  `json:"createdDate,omitempty"`
	SecretData     string `json:"secretData,omitempty"`
	CredentialData string `json:"credentialData,omitempty"`
}

type secretData struct {
	Value                string `json:"value,omitempty"`
	Salt                 string `json:"salt,omitempty"`
	AdditionalParameters any    `json:"additionalParameters,omitempty"`
}

type credentialData struct {
	HashIterations       int    `json:"hashIterations,omitempty"`
	Algorithm            string `json:"algorithm,omitempty"`
	AdditionalParameters any    `json:"additionalParameters,omitempty"`
}

func (u *user) getPassword() (string, error) {
	passwordCredential := u.getPasswordCredential()
	if passwordCredential.SecretData == "" || passwordCredential.CredentialData == "" {
		return "", nil
	}

	var (
		sd secretData
		cd credentialData
	)
	if err := json.Unmarshal([]byte(passwordCredential.SecretData), &sd); err != nil {
		return "", fmt.Errorf("secret data: %w", err)
	}
	if err := json.Unmarshal([]byte(passwordCredential.CredentialData), &cd); err != nil {
		return "", fmt.Errorf("credential data: %w", err)
	}

	return encodePassword(sd, cd)
}

func (u *user) getPasswordCredential() credential {
	for _, c := range u.Credentials {
		if c.Type == "password" {
			return c
		}
	}
	return credential{}
}

func encodePassword(sd secretData, cd credentialData) (string, error) {
	switch cd.Algorithm {
	case pbkdf2.IdentifierSHA1, pbkdf2.IdentifierSHA224, pbkdf2.IdentifierSHA256, pbkdf2.IdentifierSHA384, pbkdf2.IdentifierSHA512:
		return encodePbkdf2(sd, cd), nil
	default:
		return "", fmt.Errorf("unsupported password algorithm: %q", cd.Algorithm)
	}
}

func encodePbkdf2(sd secretData, cd credentialData) string {
	return fmt.Sprintf(pbkdf2.Format, cd.Algorithm, cd.HashIterations, sd.Salt, sd.Value)
}
