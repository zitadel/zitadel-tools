package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/zitadel/zitadel-tools/cmd/basicauth"
	"github.com/zitadel/zitadel-tools/cmd/data"
	"github.com/zitadel/zitadel-tools/cmd/jwt"
	"github.com/zitadel/zitadel-tools/cmd/migration"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "zitadel-tools",
	Short: "ZITADEL tools provides you with some helper tools",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
func SetVersionInfo(version, commit, date string) {
	rootCmd.Version = fmt.Sprintf("%s (Built on %s from Git SHA %s)", version, date, commit)
}

func init() {
	rootCmd.AddCommand(jwt.Cmd)
	rootCmd.AddCommand(basicauth.Cmd)
	rootCmd.AddCommand(migration.Cmd)
	rootCmd.AddCommand(data.Cmd)
}
