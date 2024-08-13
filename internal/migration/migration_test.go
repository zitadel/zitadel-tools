// Package migration provides common tools for the data migrations.
package migration

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"
	v1 "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/v1"
)

var outDir string

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

func TestReadJSONFile(t *testing.T) {
	const jsonStream = `{"Name": "Ed", "Text": "Knock knock."}`

	filename := filepath.Join(outDir, "test.json")
	err := os.WriteFile(filename, []byte(jsonStream), 0666)
	require.NoError(t, err)
	t.Cleanup(func() {
		err := os.Remove(filename)
		require.NoError(t, err)
	})

	type m struct {
		Name, Text string
	}

	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    *m
		wantErr bool
	}{
		{
			name: "invalid file",
			args: args{
				name: "foo",
			},
			wantErr: true,
		},
		{
			name: "decoding error",
			args: args{
				name: "/dev/urandom",
			},
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				name: filename,
			},
			want: &m{Name: "Ed", Text: "Knock knock."},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadJSONFile[*m](tt.args.name)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestReadJSONLinesFile(t *testing.T) {
	const jsonStream = `
	{"Name": "Ed", "Text": "Knock knock."}
	{"Name": "Sam", "Text": "Who's there?"}
	{"Name": "Ed", "Text": "Go fmt."}
	{"Name": "Sam", "Text": "Go fmt who?"}
	{"Name": "Ed", "Text": "Go fmt yourself!"}`

	filename := filepath.Join(outDir, "lines.json")
	err := os.WriteFile(filename, []byte(jsonStream), 0666)
	require.NoError(t, err)
	t.Cleanup(func() {
		err := os.Remove(filename)
		require.NoError(t, err)
	})

	type m struct {
		Name, Text string
	}

	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    []m
		wantErr bool
	}{
		{
			name: "invalid file",
			args: args{
				name: "foo",
			},
			wantErr: true,
		},
		{
			name: "decoding error",
			args: args{
				name: "/dev/urandom",
			},
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				name: filename,
			},
			want: []m{
				{Name: "Ed", Text: "Knock knock."},
				{Name: "Sam", Text: "Who's there?"},
				{Name: "Ed", Text: "Go fmt."},
				{Name: "Sam", Text: "Go fmt who?"},
				{Name: "Ed", Text: "Go fmt yourself!"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadJSONLinesFile[m](tt.args.name)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCreateV1Migration(t *testing.T) {
	users := []User{{
		UserId:        "user1",
		UserName:      "foobar",
		FirstName:     "foo",
		LastName:      "bar",
		Email:         "foo@bar.com",
		EmailVerified: true,
		PasswordHash:  "xxx",
	}}
	OrganizationID = "123"
	Timeout = time.Minute

	want := &admin.ImportDataRequest{
		Timeout: "1m0s",
		Data: &admin.ImportDataRequest_DataOrgs{
			DataOrgs: &admin.ImportDataOrg{
				Orgs: []*admin.DataOrg{
					{
						OrgId: "123",
						HumanUsers: []*v1.DataHumanUser{
							{
								UserId: "user1",
								User: &management.ImportHumanUserRequest{
									UserName: "foobar",
									Profile: &management.ImportHumanUserRequest_Profile{
										FirstName: "foo",
										LastName:  "bar",
									},
									Email: &management.ImportHumanUserRequest_Email{
										Email:           "foo@bar.com",
										IsEmailVerified: true,
									},
									HashedPassword: &management.ImportHumanUserRequest_HashedPassword{
										Value: "xxx",
									},
								},
							},
						},
					},
				},
			},
		},
	}
	got := CreateV1Migration(users)
	assert.Equal(t, want, got)
}

func TestWriteProtoToFile(t *testing.T) {
	type args struct {
		importData *admin.ImportDataRequest
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "no multiLine",
			args: args{
				importData: &admin.ImportDataRequest{
					Timeout: "123",
					Data: &admin.ImportDataRequest_DataOrgs{
						DataOrgs: &admin.ImportDataOrg{
							Orgs: []*admin.DataOrg{
								{
									OrgId: "org1",
									HumanUsers: []*v1.DataHumanUser{
										{
											UserId: "user11",
											User: &management.ImportHumanUserRequest{
												UserName: "name11",
												Profile: &management.ImportHumanUserRequest_Profile{
													FirstName: "first11",
													LastName:  "last11",
												},
												Email: &management.ImportHumanUserRequest_Email{
													Email:           "first11@zitadel.com",
													IsEmailVerified: true,
												},
												HashedPassword: &management.ImportHumanUserRequest_HashedPassword{
													Value: "hash11",
												},
											},
										},
										{
											UserId: "user12",
											User: &management.ImportHumanUserRequest{
												UserName: "name12",
												Profile: &management.ImportHumanUserRequest_Profile{
													FirstName: "first12",
													LastName:  "last12",
												},
												Email: &management.ImportHumanUserRequest_Email{
													Email:           "first12@zitadel.com",
													IsEmailVerified: true,
												},
												HashedPassword: &management.ImportHumanUserRequest_HashedPassword{
													Value: "hash12",
												},
											},
										},
									},
								},
								{
									OrgId: "org2",
									HumanUsers: []*v1.DataHumanUser{
										{
											UserId: "user21",
											User: &management.ImportHumanUserRequest{
												UserName: "name21",
												Profile: &management.ImportHumanUserRequest_Profile{
													FirstName: "first21",
													LastName:  "last21",
												},
												Email: &management.ImportHumanUserRequest_Email{
													Email:           "first11@zitadel.com",
													IsEmailVerified: true,
												},
												HashedPassword: &management.ImportHumanUserRequest_HashedPassword{
													Value: "hash21",
												},
											},
										},
										{
											UserId: "user22",
											User: &management.ImportHumanUserRequest{
												UserName: "name22",
												Profile: &management.ImportHumanUserRequest_Profile{
													FirstName: "first22",
													LastName:  "last22",
												},
												Email: &management.ImportHumanUserRequest_Email{
													Email:           "first22@zitadel.com",
													IsEmailVerified: true,
												},
												HashedPassword: &management.ImportHumanUserRequest_HashedPassword{
													Value: "hash22",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			want: `{"dataOrgs":{"orgs":[{"orgId":"org1","humanUsers":[{"userId":"user11","user":{"userName":"name11","profile":{"firstName":"first11","lastName":"last11"},"email":{"email":"first11@zitadel.com","isEmailVerified":true},"hashedPassword":{"value":"hash11"}}},{"userId":"user12","user":{"userName":"name12","profile":{"firstName":"first12","lastName":"last12"},"email":{"email":"first12@zitadel.com","isEmailVerified":true},"hashedPassword":{"value":"hash12"}}}]},{"orgId":"org2","humanUsers":[{"userId":"user21","user":{"userName":"name21","profile":{"firstName":"first21","lastName":"last21"},"email":{"email":"first11@zitadel.com","isEmailVerified":true},"hashedPassword":{"value":"hash21"}}},{"userId":"user22","user":{"userName":"name22","profile":{"firstName":"first22","lastName":"last22"},"email":{"email":"first22@zitadel.com","isEmailVerified":true},"hashedPassword":{"value":"hash22"}}}]}]},"timeout":"123"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			OutputPath = filepath.Join(outDir, "importBody.json")
			t.Cleanup(func() {
				require.NoError(t, os.Remove(OutputPath))
			})
			err := WriteProtoToFile(tt.args.importData)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			got, err := os.ReadFile(OutputPath)
			require.NoError(t, err)
			assert.JSONEq(t, tt.want, string(got))
		})
	}
}
