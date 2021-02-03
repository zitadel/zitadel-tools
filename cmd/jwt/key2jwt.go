package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/caos/oidc/pkg/oidc"
	"github.com/caos/oidc/pkg/rp"
)

var (
	keyPath    = flag.String("key", "", "path to the key.json")
	issuer     = flag.String("issuer", "", "issuer where the token will be used (e.g. https://issuer.zitadel.ch)")
	outputPath = flag.String("output", "", "path where the generated jwt will be saved; will print to stdout if empty")
)

func main() {
	flag.Parse()

	if *keyPath == "" || *issuer == "" {
		fmt.Println("Please provide at least an issuer and key param:")
		flag.PrintDefaults()
		return
	}

	key, err := ioutil.ReadFile(*keyPath)
	if err != nil {
		fmt.Printf("error reading key file: %v", err.Error())
		return
	}
	assertion, err := oidc.NewJWTProfileAssertionFromFileData(key, []string{*issuer})
	if err != nil {
		fmt.Printf("error generating assertion: %v", err.Error())
		return
	}
	jwt, err := rp.GenerateJWTProfileToken(assertion)
	if err != nil {
		fmt.Printf("error generating key: %v", err.Error())
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
