package main

import (
	"database/sql"
	"daya-listrik-api/internal/db"
	"daya-listrik-api/internal/handlers"
	"daya-listrik-api/internal/repository"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	dbConn, err := db.Connect()
	if err != nil {
		log.Fatal("Database connection error: ", err)
	}
	defer dbConn.Close()

	r := initializeRouter(dbConn)

	startServer(r)
}

func initializeRouter(dbConn *sql.DB) *mux.Router {
	r := mux.NewRouter()
	repo := &repository.EnergyRecordRepository{DB: dbConn}
	handlers.InitializeRoutes(r, repo)
	return r
}

func startServer(router http.Handler) {
	const addr = ":8080"

	// Bungkus router dengan middleware CORS
	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	}).Handler(router)

	fmt.Printf("Server is running on http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, handler))
}