package main

import (
	"log"
	"main/database"
	"main/handlers"
	"main/middleware"
	"net/http"
	"os"

	"github.com/gorilla/mux"
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

	database.InitTables()

	router := mux.NewRouter()

	router.Use(middleware.CORSMiddleware)

	// routes

	router.HandleFunc("/register", handlers.RegisterHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/login", handlers.LoginHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/logout", handlers.LogoutHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/search", handlers.Search).Methods("GET")

	//anime
	router.HandleFunc(
		"/anime/{anime_id}",
		middleware.OptionalMiddleware(handlers.GetAnimePage),
	).Methods("GET", "OPTIONS")

	router.HandleFunc("/user/anime", middleware.AuthMiddleware(handlers.AddAnimeToUserList)).Methods("POST")

	router.HandleFunc("/anime/{id}/reviews", middleware.AuthMiddleware(handlers.AddReview)).Methods("POST", "OPTIONS")

	router.HandleFunc("/user/anime/{id}", middleware.AuthMiddleware(handlers.RemoveAnimeFromList)).Methods("DELETE", "OPTIONS")

	router.HandleFunc("/reviews/{id}", middleware.AuthMiddleware(handlers.RemoveReview)).Methods("DELETE", "OPTIONS")

	//user
	router.HandleFunc(
		"/user/{id}",
		middleware.OptionalMiddleware(handlers.UserPageHandler),
	).Methods("GET", "OPTIONS")

	router.HandleFunc("/user/{id}/reviews", middleware.OptionalMiddleware(handlers.UserReviewsHandler)).Methods("GET", "OPTIONS")

	router.HandleFunc("/user/anime/{id}/episodes", middleware.AuthMiddleware(handlers.UpdateEpisodes)).Methods("PATCH", "OPTIONS")

	router.HandleFunc("/user/anime/{id}/status", middleware.AuthMiddleware(handlers.UpdateAnimeStatus)).Methods("PATCH", "OPTIONS")

	router.HandleFunc("/user/anime/{id}/favorite", middleware.AuthMiddleware(handlers.UpdateAnimeFavorite)).Methods("PATCH", "OPTIONS")

	//auth
	router.HandleFunc("/me", handlers.MeHandler).Methods("GET", "OPTIONS")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("0.0.0.0:" + port)
	http.ListenAndServe("0.0.0.0:"+port, router)
	log.Println("Serveur démarré sur le port " + "0.0.0.0:" + port)
}
