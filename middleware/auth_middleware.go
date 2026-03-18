package middleware

import (
	"net/http"
	"strings"
	"context"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//read the authorization header
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		//remove the bearer prefix
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		//verify the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("my_secret_key"), nil
		})

		//if parsing fails or token is invalid, return unauthorized
		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		//read the data inside the token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "invalid token claims", http.StatusUnauthorized)
			return
		}

		//extract the user id
		userID := int(claims["user_id"].(float64))

		//add the user id to the request context
		ctx := context.WithValue(r.Context(), "user_id", userID)
		r = r.WithContext(ctx)

		//call the next handler
		next.ServeHTTP(w, r)
	})
}