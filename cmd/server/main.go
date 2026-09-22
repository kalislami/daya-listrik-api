package main

import (
	"context"
	"daya-listrik-api/internal/config"
	"daya-listrik-api/internal/db"
	"daya-listrik-api/internal/handlers"
	"daya-listrik-api/internal/repository"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	c, err := config.Load()
	if err != nil {
		return err
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	dbConn, err := db.Connect(ctx, c)
	if err != nil {
		return err
	}
	defer dbConn.Close()
	log.Print("database connection established")
	if err := db.RunMigrations(dbConn); err != nil {
		return err
	}
	log.Print("database migrations applied")

	app := fiber.New(fiber.Config{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	})
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: c.CORSOrigins,
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Content-Type",
	}))
	app.Use(swagger.New(swagger.Config{FilePath: "./docs/openapi.yaml", Path: "swagger", Title: "Daya Listrik API"}))
	handlers.InitializeRoutes(app, &handlers.EnergyRecordHandler{Repo: &repository.EnergyRecordRepository{DB: dbConn}})

	listenErr := make(chan error, 1)
	go func() { listenErr <- app.Listen(":" + c.ServerPort) }()
	log.Printf("server listening on port %s", c.ServerPort)
	select {
	case err := <-listenErr:
		if err != nil {
			return err
		}
	case <-ctx.Done():
		log.Print("shutting down server")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := app.ShutdownWithContext(shutdownCtx); err != nil && !errors.Is(err, context.Canceled) {
			return err
		}
		if err := <-listenErr; err != nil {
			return err
		}
	}
	log.Print("server stopped; closing database")
	return nil
}
