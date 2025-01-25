package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GreeterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}
