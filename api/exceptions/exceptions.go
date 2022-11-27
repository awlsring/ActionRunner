package exceptions

import (
	"fmt"

	model "github.com/awlsring/action-runner-model"
	"github.com/gin-gonic/gin"
)

func InvalidInputResponse(c *gin.Context, message string) {
	resp := model.InvalidInputErrorResponseContent{
		Message: fmt.Sprintf("InvalidInputError: %v", message),
	}
	c.JSON(400, resp)
}

func ActionNotFoundResponse(c *gin.Context, message string) {
	resp := model.ActionNotFoundErrorResponseContent{
		Message: fmt.Sprintf("ActionNotFoundError: %v", message),
	}
	c.JSON(404, resp)
}

func ExecutionNotFoundResponse(c *gin.Context, message string) {
	resp := model.ActionNotFoundErrorResponseContent{
		Message: fmt.Sprintf("ExecutionNotFoundError: %v", message),
	}
	c.JSON(404, resp)
}

func InternalErrorResponse(c *gin.Context, message string) {
	resp := model.InternalServerErrorResponseContent{
		Message: fmt.Sprintf("InternalServerError: %v", message),
	}
	c.JSON(500, resp)
}