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
 - email verified (--email-verified) When the `--email-verified` flag is set, all user emails are considered verified. With this flag unset, all users need to verify their email. 

```bash
zitadel-tools migrate auth0 --org=<organisation id> --users=./users.json --passwords=./passwords.json --output=./importBody.json --timeout=1h --multiline --email-verified
```

You will now get a new file importBody.json
Copy the content from the file and send it as body in the import to ZITADEL

For a more detailed description of the whole migration steps from Auth0 to ZITADEL please visit out Documentation:
https://zitadel.com/docs/guides/migrate/sources/auth0

## Data transformation

Data is currently transformed as such:

### User Fields Mapping

| Source (Auth0)              | Destination (ZITADEL)     | Notes |
| --------------------------- | ------------------------- | ----- |
| `user_id`                   | `userId`                  | Unique identifier |
| `email`                     | `email`                   | Primary email address |
| `email`                     | `userName`                | Used as fallback if `username` is empty |
| `username`                  | `userName`                | Preferred if available |
| `given_name`                | `firstName`               | Falls back to `name` or `userName` if empty |
| `family_name`               | `lastName`                | Falls back to `name` or `userName` if empty |
| `name`                      | `displayName`             | Full display name |
| `nickname`                  | `nickName`                | User's nickname |
| `locale`                    | `preferredLanguage`       | Mapped from Auth0 locale to ZITADEL language codes¹ |
| `phone_number`              | `phone`                   | User's phone number |
| `phone_verified`            | `isPhoneVerified`         | Phone verification status |
| `--email-verified` flag     | `isEmailVerified`         | Email verification (CLI flag only, default: false) |
| Password hash (from passwords.json) | `hashedPassword` | Bcrypt password hash |

### Locale Mapping

¹ Auth0 locales are mapped to ZITADEL's supported language codes:

| Auth0 Locale Examples       | ZITADEL Language Code |
| --------------------------- | --------------------- |
| `en`, `en-US`, `en-GB`      | `en`                  |
| `es`, `es-AR`, `es-MX`      | `es`                  |
| `fr`, `fr-FR`, `fr-CA`      | `fr`                  |
| `de`, `de-DE`, `de-AT`      | `de`                  |
| `pt`, `pt-BR`, `pt-PT`      | `pt`                  |
| `it`, `it-IT`              | `it`                  |
| `ja`, `ja-JP`              | `ja`                  |
| `pl`, `pl-PL`              | `pl`                  |
| `ru`, `ru-RU`              | `ru`                  |
| `zh`, `zh-CN`, `zh-TW`      | `zh`                  |
| Unsupported locales         | *(ignored)*           |

### Fallback Logic

- **firstName**: `given_name` → `name` → `userName` (email or username)
- **lastName**: `family_name` → `name` → `userName` (email or username) → `firstName`
- **userName**: `username` → `email`
