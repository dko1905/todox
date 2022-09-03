package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func InternalServerError(c *gin.Context, err error) {
	panic(err)
	// c.AbortWithStatusJSON(http.StatusInternalServerError, err)
}

func ClientError(c *gin.Context, err error) {
	c.String(http.StatusBadRequest, err.Error())
}
