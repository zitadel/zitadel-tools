# Auth0 migration

The auth0 migration tool creates a json file which represents the body for an import request to the ZITADEL API.
With this example an organization in ZITADEL has to be existing and only the users with passwords will be imported.

## Basic usage

The migration requires the following input:
 - organisation id (--org)
 - users.json file with your exported Auth0 users (--users; default is ./users.json)
 - password.json file with the exported Auth0 bcrypt passwords (--passwords; default is ./passwords.json)

Execute the transformation and provide at least the organisation id:
```bash
zitadel-tools migrate auth0 --org=<organisation id>
```

## Advanced usage

You can specify additional parameters:
 - output path (--output; default is ./importBody.json)
 - timeout for the import data request (--timeout; default is 30m)
 - pretty print the output JSON (--multiline)

```bash
zitadel-tools migrate auth0 --org=<organisation id> --users=./users.json --passwords=./passwords.json --output=./importBody.json --timeout=1h --multiline
```

You will now get a new file importBody.json
Copy the content from the file and send it as body in the import to ZITADEL

For a more detailed description of the whole migration steps from Auth0 to ZITADEL please visit out Documentation:
https://zitadel.com/docs/guides/migrate/sources/auth0

## Data transformation

Data is currently transformed as such:

<!-- TODO: https://github.com/zitadel/zitadel-tools/issues/99 -->


| Source (Auth0)              | Destination (ZITADEL) |
| --------------------------- | --------------------- |
| `email`                     | `userName`            |
| `name`                      | `firstName`           |
| `name`                      | `lastName`            |
| `--email-verified` flag[^1] | `isEmailVerified`     |

[^1]: When the `--email-verified` flag is set, all user emails are considered verified. With this flag unset, all users need to verify their email.
