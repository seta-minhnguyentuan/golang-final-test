package main

import (
	"golang-final-test/internal/config"
	"golang-final-test/internal/database"
	httpserver "golang-final-test/internal/http"
	"golang-final-test/internal/repository"
	"golang-final-test/internal/service"
	"log"
	"net/http"
	"time"
)

func main() {
	dbCfg := config.LoadDB()
	srvCfg := config.LoadServerConfig()

	db, err := database.Connect(*dbCfg)

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := database.Migrate(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	postRepo := repository.NewPostRepository(db)
	postService := service.NewPostService(postRepo)

	engine := httpserver.NewRouter(httpserver.RouterDependencies{
		PostService: postService,
	})

	srv := &http.Server{
		Addr:           ":" + srvCfg.Port,
		Handler:        engine,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	srv.ListenAndServe()
}
