package main

import (
	"github.com/simonnik/GB_PG_GO_1/hw5/internal/apiserver"
	"github.com/simonnik/GB_PG_GO_1/hw5/internal/config"
	"github.com/simonnik/GB_PG_GO_1/hw5/pkg/post/service"
	"github.com/simonnik/GB_PG_GO_1/hw5/pkg/post/storage"
	"log"
)

func main() {
	cfg, err := config.BuildConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := storage.NewDB(cfg)
	if err != nil {
		log.Fatalf("[ERR]: failed to initialize the DB connector: %v", err)
	}

	defer db.Close()

	// init services
	postService := service.NewPostService(db)

	// start api server
	srv := apiserver.NewAPIServer(cfg.Port, postService)

	log.Println("Let's Go!")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("[ERR]: %v", err)
	}
}
