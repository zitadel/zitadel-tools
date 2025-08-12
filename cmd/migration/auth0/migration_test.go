package auth0

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/zitadel-tools/internal/migration"
)

const dataDir = "example-data"

var (
	userFile      = filepath.Join(dataDir, "users.json")
	passwordFile  = filepath.Join(dataDir, "passwords.json")
	referenceFile = filepath.Join(dataDir, "referenceOutput.json")
	outDir        string
)

func TestMain(m *testing.M) {
	var err error
	outDir, err = os.MkdirTemp("", "go_test_out")
	if err != nil {
		panic(err)
	}
	res := m.Run()
	err = os.RemoveAll(outDir)
	if err != nil {
		panic(err)
	}
	os.Exit(res)
}

func Test_migrate(t *testing.T) {
	type args struct {
		// global
		OutputPath string

		// package
		userPath     string
		passwordPath string
	}
	tests := []struct {
		name          string
		args          args
		referenceFile string
		wantErr       bool
	}{
		{
			name: "user file error",
			args: args{
				userPath: "foo",
			},
			wantErr: true,
		},
		{
			name: "password file error",
			args: args{
				userPath:     userFile,
				passwordPath: "foo",
			},
			wantErr: true,
		},
		{
			name: "write error",
			args: args{
				OutputPath:   "/foo/bar/xxx/out.json",
				userPath:     userFile,
				passwordPath: passwordFile,
			},
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				OutputPath:   filepath.Join(outDir, "importData.json"),
				userPath:     userFile,
				passwordPath: passwordFile,
			},
			referenceFile: referenceFile,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			migration.OutputPath = tt.args.OutputPath
			migration.OrganizationID = "123"
			migration.Timeout = 30 * time.Minute
			verifiedEmails = true
			migration.MultiLine = true
			userPath = tt.args.userPath
			passwordPath = tt.args.passwordPath

			t.Cleanup(func() {
				err := os.RemoveAll(tt.args.OutputPath)
				require.NoError(t, err)
			})

			err := migrate()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			if tt.referenceFile != "" {
				want, err := os.ReadFile(tt.referenceFile)
				require.NoError(t, err)
				got, err := os.ReadFile(tt.args.OutputPath)
				require.NoError(t, err)
				assert.JSONEq(t, string(want), string(got))
			}
		})
	}
}

func Test_createHumanUsers(t *testing.T) {
	tests := []struct {
		name      string
		users     []user
		passwords []password
		want      []migration.User
	}{
		{
			name: "user with empty lastName should use userName as fallback",
			users: []user{
				{
					UserId:     "test1",
					Email:      "test@example.com",
					GivenName:  "John",
					FamilyName: "", // Empty lastName
					Name:       "", // Empty name
				},
			},
			passwords: []password{},
			want: []migration.User{
				{
					UserId:    "test1",
					UserName:  "test@example.com",
					FirstName: "John",
					LastName:  "test@example.com", // Should fallback to userName (email)
					Email:     "test@example.com",
				},
			},
		},
		{
			name: "user with all empty names and lastName should use firstName as final fallback",
			users: []user{
				{
					UserId:     "test2",
					Email:      "empty@example.com",
					GivenName:  "",
					FamilyName: "",
					Name:       "",
					Username:   "",
				},
			},
			passwords: []password{},
			want: []migration.User{
				{
					UserId:    "test2",
					UserName:  "empty@example.com",
					FirstName: "empty@example.com",
					LastName:  "empty@example.com",
					Email:     "empty@example.com",
				},
			},
		},
		{
			name: "edge case - completely empty lastName should use firstName as final fallback",
			users: []user{
				{
					UserId:     "test3",
					Email:      "test3@example.com",
					GivenName:  "Jane",
					FamilyName: "",
					Name:       "",
					Username:   "",
				},
			},
			passwords: []password{
				{
					Email:        "test3@example.com",
					PasswordHash: "hash123",
				},
			},
			want: []migration.User{
				{
					UserId:       "test3",
					UserName:     "test3@example.com",
					FirstName:    "Jane",
					LastName:     "test3@example.com", // Falls back to userName, not firstName
					Email:        "test3@example.com",
					PasswordHash: "hash123",
				},
			},
		},
		{
			name: "lastName fallback to firstName when all other fields are empty",
			users: []user{
				{
					UserId:     "test4",
					Email:      "",
					GivenName:  "OnlyFirst",
					FamilyName: "",
					Name:       "",
					Username:   "",
				},
			},
			passwords: []password{},
			want: []migration.User{
				{
					UserId:    "test4",
					UserName:  "",
					FirstName: "OnlyFirst",
					LastName:  "OnlyFirst", // Should use firstName as fallback
					Email:     "",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := createHumanUsers(tt.users, tt.passwords)
			assert.Equal(t, len(tt.want), len(got))
			for i, want := range tt.want {
				assert.Equal(t, want.UserId, got[i].UserId)
				assert.Equal(t, want.UserName, got[i].UserName)
				assert.Equal(t, want.FirstName, got[i].FirstName)
				assert.Equal(t, want.LastName, got[i].LastName)
				assert.Equal(t, want.Email, got[i].Email)
			}
		})
	}
}

func Test_mapAuth0LocaleToZitadelLanguage(t *testing.T) {
	tests := []struct {
		name        string
		auth0Locale string
		want        string
	}{
		{
			name:        "empty locale",
			auth0Locale: "",
			want:        "",
		},
		{
			name:        "supported simple locale",
			auth0Locale: "en",
			want:        "en",
		},
		{
			name:        "supported complex locale",
			auth0Locale: "en-US",
			want:        "en",
		},
		{
			name:        "unsupported locale",
			auth0Locale: "xyz",
			want:        "",
		},
		{
			name:        "unsupported complex locale",
			auth0Locale: "xyz-ABC",
			want:        "abc",
		},
		{
			name:        "case insensitive matching",
			auth0Locale: "EN-US",
			want:        "en",
		},
		{
			name:        "multiple hyphens",
			auth0Locale: "es-419-MX",
			want:        "es",
		},
		{
			name:        "all supported languages",
			auth0Locale: "de-DE",
			want:        "de",
		},
		{
			name:        "edge case - single hyphen",
			auth0Locale: "-",
			want:        "",
		},
		{
			name:        "edge case - starts with hyphen",
			auth0Locale: "-en",
			want:        "en",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapAuth0LocaleToZitadelLanguage(tt.auth0Locale)
			assert.Equal(t, tt.want, got)
		})
	}
}
