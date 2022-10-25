package api

import (
	"fmt"
	"runtime"

	log "github.com/sirupsen/logrus"

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
	go r.BackgroundRunner(q)
	
	log.Debug("Registering routes")
	router.POST("/action/:actionName", runAction(a, q, d))
	router.GET("/action", listActions(a))
	router.GET("/action/:actionName", getAction(a, r.Config.PlaybookDir))
	router.GET("/execution/:executionId", getExecution(d))
	router.GET("/execution/:executionId/detailed", getDetailedExecution(d))
    router.GET("/os", func(c *gin.Context) {
        c.String(200, runtime.GOOS)
    })
    router.Run(fmt.Sprintf(":%v", c.Port))
}