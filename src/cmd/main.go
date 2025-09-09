package main

import (
	"context"
	"fmt"
	"goAuth/internal/database/model"
	"goAuth/internal/server"
	"goAuth/internal/service/auth"
	inmemory "goAuth/internal/service/in-memory"
	"goAuth/internal/service/user"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

func gracefulShutdown(fiberServer *server.FiberServer, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")
	stop() // Allow Ctrl+C to force shutdown

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := fiberServer.ShutdownWithContext(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}

func main() {

	server := server.New()
	dbInstance := server.DB.GetDBInstance()

	makeMigration(server)

	inMemoService := inmemory.NewInMemoryStore()
	authService := auth.NewAuthenticationService(dbInstance, inMemoService)
	userService := user.NewUserService(dbInstance)

	server.SetupRoutes(authService, userService)

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	go func() {
		port, _ := strconv.Atoi(os.Getenv("PORT"))
		host := os.Getenv("HOST")
		if host == "" {
			host = "0.0.0.0"
		}
		err := server.Listen(fmt.Sprintf("%s:%d", host, port))
		if err != nil {
			panic(fmt.Sprintf("http server error: %s", err))
		}
	}()

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(server, done)

	// Wait for the graceful shutdown to complete
	<-done
	log.Println("Graceful shutdown complete.")
}

func makeMigration(server *server.FiberServer) {
	server.DB.AutoMigrate(&model.User{})
}
