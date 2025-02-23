package main

import (
	"fmt"
	"net/http"
	"urlshortener/configs"
	"urlshortener/internal/auth"
	"urlshortener/internal/link"
	"urlshortener/internal/user"
	"urlshortener/pkg/db"
	"urlshortener/pkg/middleware"
)

func main() {
	conf := configs.LoadConfig()
	database := db.NewDb(conf)
	router := http.NewServeMux()

	/// repositoryies
	linkRepository := link.NewLinkRepository(database)
	userRepository := user.NewUserRepository(database)

	/// services
	authService := auth.NewAuthService(userRepository)

	/// handler
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{Config: conf, AuthService: authService})
	link.NewLinkHandler(router, link.LinkHandlerDeps{LinkRepository: linkRepository, Config: conf})

	// middlewares
	stack := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)

	server := http.Server{
		Addr:    ":8081",
		Handler: stack(router),
	}

	fmt.Println("Server is listening on port 8081")
	server.ListenAndServe()
}
