package auth

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/benc-uk/dapr-store/pkg/env"
	"github.com/dgrijalva/jwt-go"
	"github.com/lestrrat-go/jwx/jwk"
)

const jwksURL = `https://login.microsoftonline.com/common/discovery/v2.0/keys`

var jwkSet *jwk.Set

//
// AuthMiddleware added around any route will protect it
//
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Disable if call is internal from another Dapr service (localhost) or running on dev machine
		if fwdHost := r.Header.Get("X-Forwarded-Host"); strings.Contains(fwdHost, "localhost") {
			log.Printf("### Auth (%s): Bypassing validation for host: %s\n", r.URL, fwdHost)
			next(w, r)
			return
		}

		// Disable if client id is not set, this is running in demo / unsecured mode
		clientID := env.GetEnvString("AUTH_CLIENT_ID", "")
		if len(clientID) == 0 {
			log.Printf("### Auth (%s): No validation as AUTH_CLIENT_ID is not set\n", r.URL)
			next(w, r)
			return
		}

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

		// Decode the token
		token, err := jwt.Parse(tokenString, getKeyFromJWKS)
		if err != nil {
			w.WriteHeader(401)
			return
		}

		// Now validate the decoded claims
		claims := token.Claims.(jwt.MapClaims)
		if claims["scp"] != "store-api" {
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
	// TODO: cache response so we don't have to make a request every time we want to verify a JWT
	if jwkSet == nil {
		var err error
		jwkSet, err = jwk.FetchHTTP(jwksURL)
		if err != nil {
			return nil, err
		}
	}

	keyID, ok := token.Header["kid"].(string)
	if !ok {
		return nil, errors.New("Expecting JWT header to have kid")
	}

	if key := jwkSet.LookupKeyID(keyID); len(key) == 1 {
		return key[0].Materialize()
	}

	return nil, fmt.Errorf("Unable to find key: %q", keyID)
}
