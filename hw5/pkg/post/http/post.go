package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *PostApiHandler) RetrievePosts(c *gin.Context) {
	posts, err := h.service.GetPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIBaseError{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, posts)
}
