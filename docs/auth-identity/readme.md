# Security, Identity & Authentication

The default mode of operation for the Dapr Store is in "demo mode" where there is no identity provided configured, and no security on the APIs. This makes it simple to run and allows us to focus on the Dapr aspects of the project. In this mode a demo/dummy user account can be used to sign-in and place orders in the store.

Optionally Dapr store can be configured utilise the [Microsoft identity platform](https://docs.microsoft.com/en-us/azure/active-directory/develop/) (aka Azure Active Directory v2) as an identity provider, to enable real user sign-in, and securing of the APIs.

# Registering App

Using the Azure CLI create the new app registration

```
az ad app create --display-name="Dapr Store" \
   --available-to-other-tenants=true \
   --query "appId" -o tsv
```

Make a note of the GUID returned, this is the app ID, or client ID

[Follow the guide here to further configure the app](https://docs.microsoft.com/en-us/azure/active-directory/develop/scenario-spa-app-registration#redirect-uri-msaljs-20-with-auth-code-flow), this currently can't be done from the CLI

Quick summary of the steps, from the portal under 'App registrations':

- Click _'Authentication'_
  - UNSELECT the checkbox _'ID tokens (used for implicit and hybrid flows)'_
  - Click _'Add a platform'_
    - Click _'Single page application'_
    - Enter `http://localhost:9000` as the redirect URI

If you are hosting the app anywhere else, add the relevant redirect URIs

### Expose an API - **Important Step**

Once the app is created and authentication is set up, an API scope must be added. Without this access tokens issued for the app will be scoped for the MS Graph API, and not the registered app

From the portal under 'App registrations':

- Click _'Expose an API'_
- Click _'Add a scope'_
- The scope **must** be called `store-api`
- Under _'Who can consent?'_ click _'Admins and users'_
- The other fields are not important but are required.

# Configuration

When working locally copy `.env.sample` to `.env` in the root of the project repo. In the `.env` file uncomment and set the value for `AUTH_CLIENT_ID` to your registered app's client ID (got from the steps above). This will be picked up by all services when running `make run`

Note. If running a services directly from their own directory i.e. `cmd/cart/` the `.env` file will be looked for there.

# Frontend

This library has been used https://github.com/benc-uk/msal-graph-vue to add the auth and graph services to the app.

To enable auth, when working locally - create the following file `web/frontend/.env.development.local` and set `VUE_APP_AUTH_CLIENT_ID` with your client id. Note the `VUE_APP_` prefix, this is important.

When served from the frontend-host, the frontend will try to fetch it's configuration from the `/config` endpoint and try to get `AUTH_CLIENT_ID` this way. This allows dynamic configuration of the auth feature at runtime.

When `AUTH_CLIENT_ID` is set the application behavior changes as follows:

- Login page allows users to register, and sign-in with real user accounts in Azure AD.
- If a user is signed-in, an access token is acquired via the auth service, and used for all API calls made by the frontend to the backend Dapr Store APIs. This token is requested for the scope `store-api`. The fetched access token is then added to the Authorization header of all API calls.

In both cases if `AUTH_CLIENT_ID` is not found at `/config` or if `VUE_APP_AUTH_CLIENT_ID` is not set locally - then the app falls back into "demo user mode". The auth service provided by https://github.com/benc-uk/msal-graph-vue has a demo user feature and this is used.

# Services & Token Validation

Security is enabled server side by the bearer token scheme as part of the OpenID Connect flow. The [go-rest-api auth package is used](https://github.com/benc-uk/go-rest-api#package-auth), and the `JWTValidator()` HTTP middleware function which can be used to secure any route exposed by the services.

If the environmental var `AUTH_CLIENT_ID` is **not** set, then a `PassthroughValidator` is used in place of the `JWTValidator`. This is default mode of operation, i.e. all APIs on all the services are open and can be called without authentication.

The JWTValidator function gets the access token from the authorization header, decodes it using JWK and the JWKS `https://login.microsoftonline.com/common/discovery/v2.0/keys` and verifies the claims as follows:

- `scp` should equal "store-api"
- `aud` should equal the client ID of the app

If the authorization header is missing, the bearer token is missing, or the claims are not validated - then a HTTP 401 is returned.
