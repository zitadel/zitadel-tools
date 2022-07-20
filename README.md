# zitadel-tools

## key2jwt 

Convert a *key file* to *jwt token*

### Usage

key2jwt requires two flags:

- audience: where the assertion is going to be used (e.g. https://issuer.zitadel.ch)
- key: the path to the key.json

The tool prints the result to standard output.

```zsh
go run ./cmd/jwt/*.go -audience=https://issuer.zitadel.ch -key=key.json
```

Optionally you can pass an `output` flag. This will save the jwt in the provided file path:

```zsh
go run ./cmd/jwt/*.go -audience=https://issuer.zitadel.ch -key=key.json -output=jwt.txt
```

You can also create a JWT by providing a RSA private key (.pem file). You then also need to specify the issuer of the token:
```zsh
go run ./cmd/jwt/*.go -audience=https://issuer.zitadel.ch -key=key.pem -issuer=client_id
```

## basicauth

Convert *client ID* and *client secret* to be used in *Authorization* header for [Client Secret Basic](https://docs.zitadel.com/docs/apis/openidoauth/authn-methods#client-secret-basic)

### Usage

basicauth requires two flags:

- id: client id
- secret: client secret

The tool prints the URL- and Base64 encoded result to standard output

```zsh
go run ./cmd/basicauth/*.go -id $CLIENT_ID -secret $CLIENT_SECRET
```
