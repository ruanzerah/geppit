package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ruanzerah/geppit/internal/database"
	"github.com/ruanzerah/geppit/internal/repository"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port int

	db      database.Service
	queries *repository.Queries
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	databaseService := database.New()
	NewServer := &Server{
		port: port,

		db:      databaseService,
		queries: databaseService.GetQueries(),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
