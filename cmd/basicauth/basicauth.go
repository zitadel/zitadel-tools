package main

import (
	b64 "encoding/base64"
	"flag"
	"fmt"
	"net/url"
)

var (
	clientId     = flag.String("id", "", "Client ID as string")
	clientSecret = flag.String("secret", "", "Client secret as string")
)

func main() {
	flag.Parse()

	if *clientId == "" || *clientSecret == "" {
		flag.PrintDefaults()
		panic("please provide a client ID and secret")
	}

	sEscaped := url.QueryEscape(*clientId) + ":" + url.QueryEscape(*clientSecret)

	sEnc := b64.StdEncoding.EncodeToString([]byte(sEscaped))

	fmt.Print(sEnc)
}
