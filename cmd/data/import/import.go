package data_import

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

// Cmd represents the import command
var Cmd = func(issuer *string, api *string, insecure *bool, keyPath *string) *cobra.Command {
	var dataPath string

	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import data to an instance",
		Run: func(cmd *cobra.Command, args []string) {
			importData(*issuer, *api, *insecure, *keyPath, dataPath)
		},
	}

	cmd.Flags().StringVar(&dataPath, "data", "", "path to the file containing data to import")
	cmd.MarkFlagRequired("data")

	return cmd
}

func importData(issuer string, api string, insecure bool, keyPath string, dataPath string) {
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

	data, err := os.ReadFile(dataPath)

	if err != nil {
		log.Fatalln("failed to read data file:", err)
		return
	}

	var req pb.ImportDataRequest

	err = protojson.Unmarshal(data, &req)

	if err != nil {
		log.Fatalln("failed to unmarshal data:", err)
		return
	}

	resp, err := client.ImportData(context.Background(), &req)

	if err != nil {
		log.Fatalln("failed to import data:", err)
		return
	}

	log.Println("Success: ", resp.Success)
	log.Println("Errors: ", resp.Errors)
}
