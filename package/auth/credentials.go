package auth

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	jwtgo "github.com/dgrijalva/jwt-go"
)

func ValidationToken(token string, ctx context.Context) (string, error) {
	token = strings.Replace(token, "Bearer ", "", 1)
	povider, err := oidc.NewProvider(ctx, os.Getenv("KEYCLOAK_HOST"))

	if err != nil {
		return "", errors.New("Error on create provider")
	}

	verifier := povider.Verifier(&oidc.Config{ClientID: "loan_app"})
	_, err = verifier.Verify(ctx, token)

	if err != nil {
		return "", errors.New("Invalid token")
	}

	tokenstring, _ := jwtgo.Parse(token, nil)
	claims := tokenstring.Claims.(jwtgo.MapClaims)
	email := claims["email"].(string)

	return email, nil
}
