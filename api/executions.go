package api

import (
	"net/http"

	log "github.com/sirupsen/logrus"

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
			resp := model.InvalidInputErrorResponseContent{
				Message: "Specified execution does not exist",
			}
			c.JSON(http.StatusBadRequest, resp)
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
			resp := model.InvalidInputErrorResponseContent{
				Message: "Specified execution does not exist",
			}
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		response := model.GetDetailedExecutionResponseContent{
			Summary: *ex,
		}

		c.JSON(http.StatusOK, response)
    }

    return gin.HandlerFunc(fn)
}