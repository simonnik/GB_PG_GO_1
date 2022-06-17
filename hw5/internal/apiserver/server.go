package apiserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	postHttp "github.com/simonnik/GB_PG_GO_1/hw5/pkg/post/http"
	"github.com/simonnik/GB_PG_GO_1/hw5/pkg/post/service"
)

func NewAPIServer(addr string, postSvc *service.PostService) (srv *http.Server) {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	postHandler := postHttp.NewPostApiHandler(*postSvc)

	// endpoints
	api := router.Group("/api")

	posts := api.Group("/posts")
	posts.GET("", postHandler.RetrievePosts)

	srv = &http.Server{
		Addr:    addr,
		Handler: router,
	}
	return srv
}
