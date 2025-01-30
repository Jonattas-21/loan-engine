package auth

import (
	"context"
	"errors"
	"os"
	"log"
	"strings"
	"net/http"
	"net/url"
	"encoding/json"
	"fmt"
		"github.com/Jonattas-21/loan-engine/internal/api/dto"


	"github.com/coreos/go-oidc/v3/oidc"
	jwtgo "github.com/dgrijalva/jwt-go"
)

func ValidationToken(token string, ctx context.Context) (string, error) {
	token = strings.Replace(token, "Bearer ", "", 1)
	povider, err := oidc.NewProvider(ctx, os.Getenv("KEYCLOAK_HOST"))

	if err != nil {
		return "", errors.New("Error on create provider")
	}

	verifier := povider.Verifier(&oidc.Config{SkipClientIDCheck: true})
	_, err = verifier.Verify(ctx, token)

	log.Println("Error on verify token", err)
	log.Println("Token", token)
	if err != nil {
		return "", errors.New("Invalid token")
	}

	tokenstring, _ := jwtgo.Parse(token, nil)
	claims := tokenstring.Claims.(jwtgo.MapClaims)
	email := claims["email"].(string)

	return email, nil
}

func GetTokenFromKeycloak(username, password string) (*dto.TokenResponse_dto, error) {
    keycloakURL := os.Getenv("KEYCLOAK_HOST") + "/protocol/openid-connect/token"

    data := url.Values{}
    data.Set("grant_type", "password")
    data.Set("client_id", os.Getenv("KEYCLOAK_CLIENT_ID"))
    data.Set("username", username)
    data.Set("password", password)

    req, err := http.NewRequest("POST", keycloakURL, strings.NewReader(data.Encode()))
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }
    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to execute request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
    }

    var tokenResponse dto.TokenResponse_dto
    if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
        return nil, fmt.Errorf("failed to decode response: %w", err)
    }

    return &tokenResponse, nil
}