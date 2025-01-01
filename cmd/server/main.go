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
	fmt.Printf("Server is running on http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
