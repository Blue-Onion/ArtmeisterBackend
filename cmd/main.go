package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Blue-Onion/RestApi-Go/config"
	"github.com/Blue-Onion/RestApi-Go/handler"
	"github.com/Blue-Onion/RestApi-Go/middleware"

	"github.com/Blue-Onion/RestApi-Go/handler/user"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func main() {
	//Load Env
	cfg := config.LoadConfig()

	//DB
	apiCfg, err := config.DbQuries()
	if err != nil {
		log.Fatalf("Couldn't connect to database: %v", err)
	}

	//Handlers
	userHandler := &user.Handler{
		Repo: apiCfg.UserRepo,
	}
	middlewareHandler := &middleware.Handler{
		Repo: apiCfg.UserRepo,
	}

	//Server
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	router.Use(middleware.MiddlewareRateLimit)
	router.Get("/health", handler.Health)
	router.Get("/", handler.MainPage)

	// User Routes
	userRoute := chi.NewRouter()
	userRoute.Post("/users", userHandler.HandleCreateUser)
	userRoute.Post("/login", userHandler.HandleLogin)
	userRoute.Post("/logOut", middlewareHandler.MiddlewareAuth(http.HandlerFunc(userHandler.HandleLogOut)))

	router.Mount("/api", userRoute)

	server := http.Server{
		Handler: router,
		Addr:    ":" + cfg.Port,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Listening on http://localhost:%s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error occurred: %v", err)
		}
	}()

	<-stop
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Error occurred in Shutdown: %v", err)
	}
	log.Println("Server Shutdown gracefully")
}
