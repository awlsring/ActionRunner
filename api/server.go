package api

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/awlsring/action-runner/api/auth"
	"github.com/awlsring/action-runner/data"
	"github.com/awlsring/action-runner/runner"
	"github.com/gin-gonic/gin"
)

type Config struct {
	Port string `mapstructure:"port"`
}

func Start(c Config, a *map[string]string, r *runner.AnsibleRunner, d *data.ExecutionRespository) {
    log.Debug("Starting API")
	router := gin.Default()

	log.Debug("Valid actions for API")
	log.Debug(*a)
	
	q := make(chan *runner.RunRequest)
	log.Debug("Starting background runner")
	go r.BackgroundRunner(q)
	
	log.Debug("Registering routes")

	action := router.Group("/action")
	action.Use(auth.ApiKeyMiddleware())
	execution := router.Group("/execution")
	execution.Use(auth.ApiKeyMiddleware())

	action.POST("/:actionName", runAction(a, q, d))
	action.GET("", listActions(a))
	action.GET("/:actionName", getAction(a, r.Config.PlaybookDir))
	execution.GET("/:executionId", getExecution(d))
	execution.GET("/:executionId/detailed", getDetailedExecution(d))
    router.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
			"message": "healthy",
		})
	})
    router.Run(fmt.Sprintf(":%v", c.Port))
}