package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//Auth is used to verify HTTP header
func Auth(context *gin.Context) {
	authKey := context.GetHeader("Authorization")

	if authKey != "November 10, 2009" {
		context.JSON(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		context.Abort()
		return
	}

	context.Next()
}
