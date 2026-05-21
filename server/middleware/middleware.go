package middleware

import (
	"context"
	"fmt"
	"main/contextkeys"
	"main/services"
	"net/http"
)

// filtre invité/connecté mais pas proprio de la page/ proprio de la page
func OptionalMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie("session_token")
		fmt.Println(err)

		// non connecté : guest
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		// connecté mais non prorio
		userId, err := services.ValidateToken(cookie.Value)
		if err != nil {
			fmt.Println("not connected")
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), contextkeys.UserID, userId)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// filtre page/action protégé seul le proprio peut acceder/faire à la page/ l'action
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie("session_token")

		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userId, err := services.ValidateToken(cookie.Value)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), contextkeys.UserID, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
