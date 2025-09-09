package main

import (
	"daya-listrik-api/internal/db"
	"daya-listrik-api/internal/handlers"
	"daya-listrik-api/internal/repository"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Koneksi DB
	dbConn, err := db.Connect()
	if err != nil {
		log.Fatal("Database connection error: ", err)
	}
	defer dbConn.Close()

	// Inisialisasi Fiber
	app := fiber.New(fiber.Config{
		Prefork: os.Getenv("PREFORK_ENABLED") == "true", // seperti cluster di nodeJS, untuk memaksimalkan penggunaan core CPU
	})

	// Middleware CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Content-Type",
		AllowCredentials: true,
	}))

	// Inisialisasi repository
	repo := &repository.EnergyRecordRepository{DB: dbConn}

	// Bungkus repository ke handler
	handler := &handlers.EnergyRecordHandler{Repo: repo}

	// Daftarkan routes (ubah InitializeRoutes agar support Fiber)
	handlers.InitializeRoutes(app, handler)

	// Start server
	addr := os.Getenv("PORT")
	log.Printf("Server is running on http://localhost%s\n", addr)
	log.Fatal(app.Listen(addr))
}
