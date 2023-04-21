The auth0 migration tool creates a json file which represents the body for an import request to the ZITADEL API.
With this example an organization in ZITADEL has to be existing and only the users with passwords will be imported.

The migration requires the following input:
 - organisation id (--org)
 - users.json file with your exported Auth0 users (--users; default is ./users.json)
 - password.json file with the exported Auth0 bcrypt passwords (--password; default is ./passwords.json)

Execute the transformation and provide at least the organisation id:
```bash
./zitadel-tools migrate auth0 --org=<organisation id>
```

you can also specify custom path for the users, passwords and output JSON files:
```bash
./zitadel-tools migrate auth0 --org=<organisation id> --users=./users.json --password=./passwords.json --output=./importBody.json
```

You will now get a new file importBody.json
Copy the content from the file and send it as body in the import to ZITADEL

For a more detailed description of the whole migration steps from Auth0 to ZITADEL please visit out Documentation:
https://zitadel.com/docs/guides/migrate/sources/auth0
