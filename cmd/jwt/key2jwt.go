package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/caos/oidc/pkg/client"
	"github.com/caos/oidc/pkg/oidc"
)

var (
	keyPath    = flag.String("key", "", "path to the key.json")
	audience   = flag.String("audience", "", "audience where the token will be used (e.g. the issuer of zitadel.ch - https://issuer.zitadel.ch)")
	outputPath = flag.String("output", "", "path where the generated jwt will be saved; will print to stdout if empty")
)

func main() {
	flag.Parse()

	if *keyPath == "" || *audience == "" {
		fmt.Println("Please provide at least an audience and key param:")
		flag.PrintDefaults()
		return
	}

	key, err := ioutil.ReadFile(*keyPath)
	if err != nil {
		fmt.Printf("error reading key file: %v", err.Error())
		return
	}
	jwt, err := generateJWT(key)
	if err != nil {
		fmt.Printf("error generating jwt: %v", err.Error())
		return
	}
	f := os.Stdout
	if *outputPath != "" {
		f, err = os.OpenFile(*outputPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
		if err != nil {
			fmt.Printf("error reading key file: %v", err.Error())
			return
		}
	}
	_, err = f.Write([]byte(jwt))
	if errClose := f.Close(); err == nil {
		err = errClose
	}
	if err != nil {
		fmt.Printf("error writing key: %v", err.Error())
		return
	}
}

func generateJWT(key []byte) (string, error) {
	keyType, err := getType(key)
	if err != nil {
		return "", err
	}
	switch keyType {
	case "application":
		keyData, err := client.ConfigFromKeyFile(*keyPath)
		if err != nil {
			return "", err
		}
		signer, err := client.NewSignerFromPrivateKeyByte([]byte(keyData.Key), keyData.KeyID)
		if err != nil {
			return "", err
		}
		return client.SignedJWTProfileAssertion(keyData.ClientID, []string{*audience}, time.Hour, signer)
	case "serviceaccount":
		jwta, err := oidc.NewJWTProfileAssertionFromFileData(key, []string{*audience})
		if err != nil {
			return "", err
		}
		return oidc.GenerateJWTProfileToken(jwta)
	default:
		return "", fmt.Errorf("unsupported key type")
	}
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
