# Keycloak migration

The Keycloak migration tool creates a json file which represents the body for an import request to the ZITADEL API.
With this example an organization in ZITADEL has to be existing and users with only their passwords will be imported.

## Basic usage

The migration requires the following input:
 - organisation id (--org)
 - realm.json file with your exported Keycloak realm, containing users (--realm; default is ./realm.json)

Execute the transformation and provide at least the organisation id:
```bash
zitadel-tools migrate keycloak --org=<organisation id>
```

## Advanced usage

You can specify additional parameters:
 - output path (--output; default is ./importBody.json)
 - timeout for the import data request (--timeout; default is 30m)
 - pretty print the output JSON (--multiline)

```bash
zitadel-tools migrate keycloak --org=<organisation id> --realm=./realm.json --output=./importBody.json --timeout=1h --multiline
```

You will now get a new file importBody.json
Copy the content from the file and send it as body in the import to ZITADEL

For a more detailed description of the whole migration steps from Auth0 to ZITADEL please visit out Documentation:
https://zitadel.com/docs/guides/migrate/sources/keycloak

## Data transformation

- The user's `enabled` flag cannot be passed to the import. All imported users will be enabled.
- password is the only credential type supported. If there are other credentials, they are silently ignored.
- currently only passwords of the pbkdf2 algorithm family are supported and transformed into a "Modular Crypt Format" string.
- Role and group information are not transferred.
