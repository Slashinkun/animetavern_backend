package middleware

import (
	"context"
	"main/services"
	"net/http"
)

// intercepte les requetes afin de verifier que le demandeur est autorisé a obtenir la page
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")

		if err != nil {
			// pas connecté, on continue mais sans user
			next.ServeHTTP(w, r)
			return
		}

		username, err := services.ValidateToken(cookie.Value)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		//on met le nom d'utilisateur dans le contexte
		ctx := context.WithValue(r.Context(), "username", username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
