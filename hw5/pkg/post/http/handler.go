package http

import "github.com/simonnik/GB_PG_GO_1/hw5/pkg/post/service"

type PostApiHandler struct {
	service *service.PostService
}

func NewPostApiHandler(svc service.PostService) *PostApiHandler {
	return &PostApiHandler{&svc}
}
