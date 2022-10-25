package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/awlsring/action-runner/data"
	"github.com/awlsring/action-runner/inventory"
	"github.com/awlsring/action-runner/runner"
	model "github.com/awlsring/dws-action-runner"
	"github.com/gin-gonic/gin"
)

func machinesToInventory(m []model.Machine) ([]*inventory.Host, error) {
	hosts := []*inventory.Host{}
	for _, machine := range m {
		if !machine.HasAnsiblePort() {
			port := float32(22)
			machine.AnsiblePort = &port
		}

		if !machine.HasId() {
			machine.Id = &machine.Ip
		}

		h := &inventory.Host{
			Port: int(machine.GetAnsiblePort()),
			IpAddress: machine.Ip,
			ID: *machine.Id,
		}
		hosts = append(hosts, h)
	}

	return hosts, nil
}

func runAction(a *map[string]string, ch chan *runner.RunRequest, s *data.ExecutionRespository) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		log.Debug("Running runAction")
		actionName := c.Param("actionName")
		action, request, err := validateRunActionRequest(a, c)
		if err != nil {
			log.Println(err)
			return 
		}
		log.Debug("Request validated")

		log.Debug("Getting inventory")
		h, err := machinesToInventory(request.Machines)
		if err != nil {
			log.Println(err)
			return
		}
		targets := []string{}
		for _, target := range h {
			targets = append(targets, target.ID)
		}

		log.Debug("Creating execution object")
		id := s.ExecutionDao.GenerateId()
		ex := &model.ExecutionSummary{
			Id: id,
			Action: actionName,
			StartTime: float32(time.Now().Unix()),
			Status: model.RUNNING,
			Targets: targets,
		}
		log.Debug("Saving execution object")
		s.ExecutionDao.Create(ex)

		log.Debug("Sending run request")
		ch <- &runner.RunRequest{
			Action: action,
			Execution: ex,
			Request: request,
			Hosts: h,
		}

		log.Debug("Returning response")
		response := model.RunActionResponseContent{
			ExecutionId: id,
		}
		log.Debugf("Response: %v", response)
		c.IndentedJSON(http.StatusCreated, response)
    }

    return gin.HandlerFunc(fn)
}

func validateRunActionRequest(a *map[string]string, c *gin.Context) (string, *model.RunActionRequestContent, error) {
	actionName := c.Param("actionName")
	var request *model.RunActionRequestContent
	var actionFile string
	ac := *a
	if file, ok := ac[actionName]; ok {
		actionFile = file
	} else {
		resp := model.InvalidInputErrorResponseContent{
			Message: "Specified action does not exist",
		}
		c.IndentedJSON(http.StatusBadRequest, resp)
		return "", nil, errors.New("")
	}
	
	if err := c.BindJSON(&request); err != nil {
		fmt.Println(err)
	}
	if request.Machines == nil || len(request.Machines) == 0{
		resp := model.InvalidInputErrorResponseContent{
			Message: "At least one machine must be passed",
		}
		c.IndentedJSON(http.StatusBadRequest, resp)
		return "", nil, errors.New("")
	}
	return actionFile, request, nil
}