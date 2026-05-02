package middleware

import (
	"context"
	"fmt"
	"main/contextkeys"
	"main/services"
	"net/http"
	"reflect"
)

func OptionalMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// fmt.Println("\n====================")
		// fmt.Println("METHOD:", r.Method)
		// fmt.Println("URL:", r.URL.Path)
		// fmt.Println("Cookie header:", r.Header.Get("Cookie"))
		// fmt.Println("All cookies:", r.Cookies())
		// fmt.Println("====================")

		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		fmt.Println("MIDDLEWARE HIT")

		fmt.Println("COOKIE HEADER RAW:", r.Header.Get("Cookie"))
		cookie, err := r.Cookie("session_token")
		fmt.Println(err)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		userId, err := services.ValidateToken(cookie.Value)
		if err != nil {
			fmt.Println("pas connecté")
			next.ServeHTTP(w, r)
			return
		}

		fmt.Println("connecté")

		ctx := context.WithValue(r.Context(), contextkeys.UserID, userId)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("Cookies:", r.Cookies())

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

		fmt.Println("userId type:", reflect.TypeOf(userId), "value:", userId)

		ctx := context.WithValue(r.Context(), contextkeys.UserID, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
