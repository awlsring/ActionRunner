package api

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/awlsring/action-runner/api/exceptions"
	"github.com/awlsring/action-runner/data"
	model "github.com/awlsring/dws-action-runner"
	"github.com/gin-gonic/gin"
)

func getExecution(s *data.ExecutionRespository) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		exId := c.Param("executionId")
		log.Info("Given id ", exId)
		ex, err := s.ExecutionDao.Get(exId)
		if err != nil {
			log.Error(err)
			exceptions.ExecutionNotFoundResponse(c, fmt.Sprintf("Execution %s does not exist", exId))
			return
		}

		response := model.GetExecutionResponseContent{
			Summary: *ex,
		}

		c.JSON(http.StatusOK, response)
    }

    return gin.HandlerFunc(fn)
}

func getDetailedExecution(s *data.ExecutionRespository) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		exId := c.Param("executionId")
		log.Info("Given id ", exId)
		ex, err := s.ExecutionDao.GetDetailed(exId)
		if err != nil {
			log.Println(err)
			exceptions.ExecutionNotFoundResponse(c, fmt.Sprintf("Execution %s does not exist", exId))
			return
		}

		response := model.GetDetailedExecutionResponseContent{
			Summary: *ex,
		}

		c.JSON(http.StatusOK, response)
    }

    return gin.HandlerFunc(fn)
}