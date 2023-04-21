package basicauth

import (
	b64 "encoding/base64"
	"fmt"
	"log"
	"net/url"

	"github.com/spf13/cobra"
)

// Cmd represents the basicauth command
var Cmd = &cobra.Command{
	Use:   "basicauth",
	Short: "Convert <client ID> and <client secret> to be used in Authorization header for Client Secret Basic",
	Run: func(cmd *cobra.Command, args []string) {
		basicAuth(cmd)
	},
}

var (
	clientId     string
	clientSecret string
)

func init() {
	Cmd.Flags().StringVar(&clientId, "id", "", "Client ID as string")
	Cmd.Flags().StringVar(&clientSecret, "secret", "", "Client secret as string")
}

func basicAuth(cmd *cobra.Command) {
	if clientId == "" || clientSecret == "" {
		log.Println("please provide a client ID and secret")
		fmt.Println(cmd.Flags().FlagUsages())
		return
	}

	sEscaped := url.QueryEscape(clientId) + ":" + url.QueryEscape(clientSecret)

	sEnc := b64.StdEncoding.EncodeToString([]byte(sEscaped))

	fmt.Println(sEnc)
}
