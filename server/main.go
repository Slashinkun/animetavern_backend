package main

import (
	"fmt"
	"main/database"
	"main/handlers"
	"main/middleware"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {

	//connection a la DB
	err := database.InitDB()
	if err != nil {
		println("Erreur connection DB")
	} else {
		println("Connecté à la DB")
	}

	// routes
	http.HandleFunc("/anime/", middleware.OptionalMiddleware(handlers.GetAnimePage))

	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)
	http.HandleFunc("/search", handlers.Search)

	fmt.Println("Serveur démarré sur le port 8080...")

	http.ListenAndServe(":8080", nil)
}
