# zitadel-tools

## key2jwt 

Convert a key file to jwt token

### Usage

key2jwt requires two flags (and will print result in stdout):

- audience: where the assertion is going to be used (e.g. https://issuer.zitadel.ch)
- key: the path to the key.json

```
./key2jwt -audience=https://issuer.zitadel.ch -key=key.json
```

Optionally you can pass an `output` flag. This will save the jwt in the provided file path:

```
./key2jwt -audience=https://issuer.zitadel.ch -key=key.json -output=jwt.txt
```

## basicauth

Convert Client ID and secret to be used in Authentication Header for [Client Secret Basic](https://docs.zitadel.com/docs/apis/openidoauth/authn-methods#client-secret-basic)

### Usage

basicauth requires two flags (and will print result in stdout):

- id: client id
- secret: client secret

```zsh
./basicauth -id $CLIENT_ID -secret $CLIENT_SECRET
```
