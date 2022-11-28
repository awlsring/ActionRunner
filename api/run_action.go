package api

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	model "github.com/awlsring/action-runner-model"
	"github.com/awlsring/action-runner/api/exceptions"
	"github.com/awlsring/action-runner/data"
	"github.com/awlsring/action-runner/inventory"
	"github.com/awlsring/action-runner/runner"
	"github.com/gin-gonic/gin"
)

func machinesToInventory(m []model.Machine) ([]*inventory.Host, error) {
	hosts := []*inventory.Host{}
	for _, machine := range m {
		
		if !machine.HasId() {
			machine.Id = &machine.Ip
		}
		
		h := &inventory.Host{
			Port: 22,
			IpAddress: machine.Ip,
			ID: *machine.Id,
		}

		if machine.HasAnsiblePort() {
			h.Port = int(*machine.AnsiblePort)
		}

		if machine.HasAnsibleUser() {
			h.User = *machine.AnsibleUser
		}

		if machine.HasAnsiblePassword() {
			h.Password = *machine.AnsiblePassword
		}

		if machine.HasAnsibleSudoPassword() {
			h.SudoPassword = *machine.AnsibleSudoPassword
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
			return 
		}
		req, _ := request.MarshalJSON()
		log.Debug("Request: ", string(req))
		log.Debug("Request validated")

		log.Debug("Getting inventory")
		h, err := machinesToInventory(request.Machines)
		if err != nil {
			log.Println(err)
			exceptions.InternalErrorResponse(c, "Failed to build inventory")
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
		log.Errorf("Action not found: %v", actionName)
		exceptions.ActionNotFoundResponse(c, fmt.Sprintf("Action %v not found", actionName))
		return "", nil, errors.New("")
	}
	
	if err := c.BindJSON(&request); err != nil {
		log.Errorf("Invalid input: %v", err)
		exceptions.InvalidInputResponse(c, fmt.Sprintf("Invalid input: %v", err))
		return "", nil, errors.New("")
	}
	if request.Machines == nil || len(request.Machines) == 0{
		log.Error("No machines provided")
		exceptions.InvalidInputResponse(c, "No machines specified")
		return "", nil, errors.New("")
	}

	for _, machine := range request.Machines {
		if !validIPAddress(machine.Ip){
			log.Errorf("Invalid IP address: %v", machine.Ip)
			exceptions.InvalidInputResponse(c, fmt.Sprintf("Invalid IP address: '%v'", machine.Ip))
			return "", nil, errors.New("")
		}
	}
	return actionFile, request, nil
}

func validIPAddress(ip string) bool {
    return net.ParseIP(ip) != nil
}