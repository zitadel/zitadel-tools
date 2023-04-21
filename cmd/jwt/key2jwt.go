package jwt

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/zitadel/oidc/pkg/client"
	"github.com/zitadel/oidc/pkg/oidc"
)

// Cmd represents the jwt command
var Cmd = &cobra.Command{
	Use:   "key2jwt",
	Short: "Convert a <key file> to <jwt token>",
	Run: func(cmd *cobra.Command, args []string) {
		key2JWT(cmd)
	},
}

var (
	keyPath    string
	audience   string
	issuer     string
	outputPath string
)

func init() {
	Cmd.Flags().StringVar(&keyPath, "key", "", "path to the key.json / RSA private key.pem")
	Cmd.Flags().StringVar(&audience, "audience", "", "audience where the token will be used (e.g. the issuer of zitadel.cloud - https://zitadel.cloud or from your domain https://<your domain>)")
	Cmd.Flags().StringVar(&issuer, "issuer", "", "issuer of the JWT (e.g. userID / client_id; only needed when generating from RSA private key)")
	Cmd.Flags().StringVar(&outputPath, "output", "", "path where the generated jwt will be saved; will print to stdout if empty")
}

func key2JWT(cmd *cobra.Command) {
	if keyPath == "" || audience == "" {
		log.Println("Please provide at least an audience and key param:")
		fmt.Println(cmd.LocalFlags().FlagUsages())
		return
	}

	key, err := os.ReadFile(keyPath)
	if err != nil {
		log.Fatalf("error reading key file: %v", err.Error())
		return
	}
	var jwt string
	switch ext := filepath.Ext(keyPath); ext {
	case ".json":
		jwt, err = generateJWTFromJSON(key)
	case ".pem":
		if issuer == "" {
			log.Fatal("Please provide the issuer of token when using a pem file")
		}
		jwt, err = generateJWTFromPEM(key, issuer)
	default:
		log.Fatalf("file extension %v is not supported, please provide either a json or pem file\n", ext)
		return
	}
	if err != nil {
		log.Fatalf("error generating jwt: %v", err.Error())
		return
	}
	f := os.Stdout
	if outputPath != "" {
		f, err = os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
		if err != nil {
			log.Fatalf("error reading key file: %v", err.Error())
			return
		}
	}
	_, err = fmt.Fprintln(f, jwt)
	if errClose := f.Close(); err == nil {
		err = errClose
	}
	if err != nil {
		log.Fatalf("error writing key: %v", err.Error())
		return
	}
}

func generateJWTFromJSON(key []byte) (string, error) {
	keyType, err := getType(key)
	if err != nil {
		return "", err
	}
	switch keyType {
	case "application":
		keyData, err := client.ConfigFromKeyFile(keyPath)
		if err != nil {
			return "", err
		}
		signer, err := client.NewSignerFromPrivateKeyByte([]byte(keyData.Key), keyData.KeyID)
		if err != nil {
			return "", err
		}
		return client.SignedJWTProfileAssertion(keyData.ClientID, []string{audience}, time.Hour, signer)
	case "serviceaccount":
		jwta, err := oidc.NewJWTProfileAssertionFromFileData(key, []string{audience})
		if err != nil {
			return "", err
		}
		return oidc.GenerateJWTProfileToken(jwta)
	default:
		return "", fmt.Errorf("unsupported key type")
	}
}

func generateJWTFromPEM(key []byte, issuer string) (string, error) {
	signer, err := client.NewSignerFromPrivateKeyByte(key, "")
	if err != nil {
		return "", err
	}
	return client.SignedJWTProfileAssertion(issuer, []string{audience}, time.Hour, signer)
}

func getType(data []byte) (string, error) {
	keyData := new(struct {
		Type string `json:"type"` // serviceaccount or application
	})
	err := json.Unmarshal(data, keyData)
	if err != nil {
		return "", err
	}
	return keyData.Type, nil
}
