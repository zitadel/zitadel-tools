# zitadel-tools

## Usage

key2jwt requires to flags (and will print result in stdout):
- issuer: where the assertion is going to be used (e.g. https://issuer.zitadel.ch)
- key: the path to the key.json

```
./key2jwt -issuer=https://issuer.zitadel.ch -key=key.json
```

Optinally you can pass an `output` flag. This will save the jwt in the provided file path:

```
./key2jwt -issuer=https://issuer.zitadel.ch -key=key.json -output=jwt.txt
```
