package middlewares

import (
	"context"
	"github.com/Jonattas-21/loan-engine/package/auth"
	"net/http"

	"github.com/go-chi/render"
)

type ValidationFunc func(token string, ctx context.Context) (string, error)

var ValidateToken ValidationFunc = auth.ValidationToken

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("Authorization")
		if token == "" {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, map[string]string{"error": "Token not found in header Authorization"})
			return
		}

		email, err := ValidateToken(token, r.Context())
		if err != nil {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, map[string]string{"error": err.Error()})
			return
		}

		ctx := context.WithValue(r.Context(), "email", email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
