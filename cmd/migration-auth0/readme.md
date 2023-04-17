The auth0 migration tool does create an json file which represents the body for an import request to the ZITADEL API.
With this example an organization in ZITADEL has to be existing and only the users with passwords will be imported.

1. Replace the users.json file with your exported Auth0 users
2. Replace the password.json file with the exported Auth0 bcrypt passwords
3. Replace the ORG_ID const in the migration.go file with the your ZITADEL organization Id where the users should be added
4. Run the migration.go file

You will now get a new file importBody.json
Copy the content from the file and send it as body in the import to ZITADEL