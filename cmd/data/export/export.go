package data_export

import (
	"context"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/zitadel/zitadel-go/v2/pkg/client/admin"
	"github.com/zitadel/zitadel-go/v2/pkg/client/middleware"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel"
	pb "github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/admin"
	"google.golang.org/protobuf/encoding/protojson"
)

// Cmd represents the export command
var Cmd = func(issuer *string, api *string, insecure *bool, keyPath *string) *cobra.Command {
	var dataPath string

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export data from an instance",
		Run: func(cmd *cobra.Command, args []string) {
			exportData(*issuer, *api, *insecure, *keyPath, dataPath)
		},
	}

	cmd.Flags().StringVar(&dataPath, "data", "", "path to the file where exported data will be written")
	cmd.MarkFlagRequired("data")

	return cmd
}

func exportData(issuer string, api string, insecure bool, keyPath string, dataPath string) {
	opts := []zitadel.Option{
		zitadel.WithJWTProfileTokenSource(middleware.JWTProfileFromPath(keyPath)),
	}

	if insecure {
		opts = append(opts, zitadel.WithInsecure())
	}

	client, err := admin.NewClient(
		issuer,
		api,
		[]string{zitadel.ScopeZitadelAPI()},
		opts...,
	)

	if err != nil {
		log.Fatalln("failed to create admin client:", err)
		return
	}

	defer func() {
		if err := client.Connection.Close(); err != nil {
			log.Fatalln("failed to close client connection:", err)
		}
	}()

	resp, err := client.ExportData(context.Background(), &pb.ExportDataRequest{})

	if err != nil {
		log.Fatalln("failed to export data:", err)
		return
	}

	data, err := protojson.Marshal(resp)

	if err != nil {
		log.Fatalln("failed to marshal data:", err)
		return
	}

	err = os.WriteFile(dataPath, data, 0644)

	if err != nil {
		log.Fatalln("failed to write data file:", err)
		return
	}
}
