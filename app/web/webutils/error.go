package webutils

import (
	"github.com/gin-gonic/gin"
	"github.com/vonmutinda/organono/app/utils"
)

func HandleError(c *gin.Context, wrappedError *utils.Error) {
	wrappedError.LogErrorMessages()
	c.JSON(wrappedError.HttpStatus(), wrappedError.JsonResponse())
}
