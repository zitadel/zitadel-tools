package data

import (
	"github.com/spf13/cobra"
	data_export "github.com/zitadel/zitadel-tools/cmd/data/export"
	data_import "github.com/zitadel/zitadel-tools/cmd/data/import"
)

// Cmd represents the data root command
var Cmd = &cobra.Command{
	Use:   "data",
	Short: "Import/Export data",
}

func init() {
	issuer := Cmd.PersistentFlags().String("issuer", "", "issuer of your ZITADEL instance (in the form: https://<instance>.zitadel.cloud or https://<yourdomain>)")
	api := Cmd.PersistentFlags().String("api", "", "gRPC endpoint of your ZITADEL instance (in the form: <instance>.zitadel.cloud:443 or <yourdomain>:443)")
	insecure := Cmd.PersistentFlags().Bool("insecure", false, "disable TLS to connect to gRPC API (use for local development only)")
	keyPath := Cmd.PersistentFlags().String("key", "", "path to the JSON machine key")

	Cmd.MarkPersistentFlagRequired("issuer")
	Cmd.MarkPersistentFlagRequired("api")
	Cmd.MarkPersistentFlagRequired("key")

	Cmd.AddCommand(data_import.Cmd(issuer, api, insecure, keyPath))
	Cmd.AddCommand(data_export.Cmd(issuer, api, insecure, keyPath))
}
