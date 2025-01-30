package dto

type TokenResponse_dto struct {
    AccessToken string `json:"access_token"`
    ExpiresIn   int    `json:"expires_in"`
    RefreshToken string `json:"refresh_token"`
    TokenType    string `json:"token_type"`
}