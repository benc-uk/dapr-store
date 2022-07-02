// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// HandlerFunc middleware for checking JWT token validity
// ----------------------------------------------------------------------------

package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/benc-uk/dapr-store/pkg/env"
	"github.com/dgrijalva/jwt-go"
	"github.com/lestrrat-go/jwx/jwk"
)

// Fix the JWKS URL to be the one for Azure AD / MSIP
const jwksURL = `https://login.microsoftonline.com/common/discovery/v2.0/keys`
const appScopeName = "store-api"

var jwkSet jwk.Set

//
// JWTValidator added around any route will protect it
//
func JWTValidator(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// TODO! Removed as it wasn't doing anything...
		// Disable check if call is internal from another Dapr service (localhost) or running on dev machine
		// fwdHost := r.Header.Get("X-Forwarded-Host")
		// if strings.Contains(fwdHost, "localhost") || r.Host == "example.com" {
		// 	log.Printf("### Auth (%s): Bypassing validation for host: %s %s\n", r.URL, fwdHost, r.Host)
		// 	next(w, r)
		// 	return
		// }

		// Disable check if client id is not set, this is running in demo / unsecured mode
		clientID := env.GetEnvString("AUTH_CLIENT_ID", "")
		if len(clientID) == 0 {
			log.Printf("### Auth (%s): No validation as AUTH_CLIENT_ID is not set\n", r.URL)
			next(w, r)
			return
		}

		// Get auth header & bearer scheme
		authHeader := r.Header.Get("Authorization")
		if len(authHeader) == 0 {
			w.WriteHeader(401)
			return
		}
		authParts := strings.Split(authHeader, " ")
		if len(authParts) != 2 {
			w.WriteHeader(401)
			return
		}
		if strings.ToLower(authParts[0]) != "bearer" {
			w.WriteHeader(401)
			return
		}
		tokenString := authParts[1]

		// Decode the token, using getKeyFromJWKS to get the key
		token, err := jwt.Parse(tokenString, getKeyFromJWKS)
		if err != nil {
			w.WriteHeader(401)
			return
		}

		// Now validate the decoded claims
		claims := token.Claims.(jwt.MapClaims)
		if claims["scp"] != appScopeName {
			w.WriteHeader(401)
			return
		}
		if claims["aud"] != clientID {
			w.WriteHeader(401)
			return
		}

		// Otherwise, we're all good!
		log.Printf("### Auth (%s): token passed validation! [scp:%s] [aud:%s]\n", r.URL, claims["scp"], claims["aud"])
		next(w, r)
	}
}

//
// Get key for given token (from it's kid header)
//
func getKeyFromJWKS(token *jwt.Token) (interface{}, error) {
	// We only support one JWKS, but most identity platforms have just the one, right?
	if jwkSet == nil {
		var err error
		jwkSet, err = jwk.Fetch(context.Background(), jwksURL)
		if err != nil {
			return nil, err
		}
	}

	keyID, ok := token.Header["kid"].(string)
	if !ok {
		return nil, errors.New("Expecting JWT header to have kid")
	}

	if key, found := jwkSet.LookupKeyID(keyID); found {
		// This I *think* gets the key value as raw bytes
		var keyReturn interface{}
		key.Raw(&keyReturn)
		return keyReturn, nil
	}

	return nil, fmt.Errorf("Unable to find key: %q", keyID)
}
