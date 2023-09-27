package keycloak

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
	realmFile     = filepath.Join(dataDir, "realm.json")
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
		realmPath string
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
				realmPath: "foo",
			},
			wantErr: true,
		},
		{
			name: "write error",
			args: args{
				OutputPath: "/foo/bar/xxx/out.json",
				realmPath:  realmFile,
			},
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				OutputPath: filepath.Join(outDir, "importData.json"),
				realmPath:  realmFile,
			},
			referenceFile: referenceFile,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			migration.OutputPath = tt.args.OutputPath
			migration.OrganizationID = "123"
			migration.Timeout = 30 * time.Minute
			migration.MultiLine = true
			realmPath = tt.args.realmPath

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
