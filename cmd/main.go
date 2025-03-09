package main

import (
	"fmt"
	"net/http"
	"urlshortener/configs"
	"urlshortener/internal/auth"
	"urlshortener/internal/link"
	"urlshortener/internal/stat"
	"urlshortener/internal/user"
	"urlshortener/pkg/db"
	event "urlshortener/pkg/eventbus"
	"urlshortener/pkg/middleware"
)

func App() http.Handler {
	conf := configs.LoadConfig()
	database := db.NewDb(conf)
	router := http.NewServeMux()
	eventBus := event.NewEventBus()

	/// repositoryies
	linkRepository := link.NewLinkRepository(database)
	userRepository := user.NewUserRepository(database)
	statRepository := stat.NewStatRepository(database)

	/// services
	authService := auth.NewAuthService(userRepository)
	statService := stat.NewStatService(&stat.StatServiceDeps{
		EventBus:       eventBus,
		StatRepository: statRepository,
	})

	/// handler
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{Config: conf, AuthService: authService})
	link.NewLinkHandler(router, link.LinkHandlerDeps{LinkRepository: linkRepository, EventBus: eventBus, Config: conf})
	stat.NewStatHandler(router, stat.StatHandlerDeps{StatRepository: statRepository, Config: conf})

	go statService.AddClick()

	// middlewares
	stack := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)

	return stack(router)
}

func main() {

	app := App()
	server := http.Server{
		Addr:    ":8081",
		Handler: app,
	}

	fmt.Println("Server is listening on port 8081")
	server.ListenAndServe()
}
