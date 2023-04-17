The auth0 migration tool does create an json file which represents the body for an import request to the ZITADEL API.
With this example an organization in ZITADEL has to be existing and only the users with passwords will be imported.

1. Replace the users.json file with your exported Auth0 users
2. Replace the password.json file with the exported Auth0 bcrypt passwords
3. Replace 