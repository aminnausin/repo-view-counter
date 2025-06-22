package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"repo-view-counter/internal/badge"
	"repo-view-counter/internal/db"
)

type Server struct {
	port int

	db           db.Database
	badgeService badge.Service
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	databaseDriver := os.Getenv("DATABASE_DRIVER")

	db := db.NewDatabase(databaseDriver)
	badgeService := badge.NewService(db)
	NewServer := &Server{
		port: port,

		db:           db,
		badgeService: badgeService,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("Starting service on port %d...\n", port)
	log.Println(`R E K A`)

	return server
}
