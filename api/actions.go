package api

import (
	"fmt"
	"os"
	"strings"

	model "github.com/awlsring/action-runner-model"
	"github.com/awlsring/action-runner/api/exceptions"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func getComments(path string) (string, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	var playbook yaml.Node
	err = yaml.Unmarshal(file, &playbook)
	if err != nil {
		return "", err
	}
	if len(playbook.Content) != 0 {
		if len(playbook.Content[0].Content) != 0 {
			comments := playbook.Content[0].Content[0].HeadComment
			comments = strings.Replace(comments, "#", "", 1)
			comments = strings.TrimSpace(comments)
			return comments, nil
		}
	}
	return "", fmt.Errorf("no comments found")
}

func getAction(a *map[string]string, path string) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		actionName := c.Param("actionName")
		action, ok := (*a)[actionName]
		if !ok {
			exceptions.ActionNotFoundResponse(c, fmt.Sprintf("Action %v not found", actionName))
			return
		}

		comments, err := getComments(path + "/" + action)
		if err != nil {
			log.Debug("Error getting comments: ", err)
		} else {
			log.Debug("comments: ", comments)
		}
		response := model.GetActionResponseContent{
			Summary: model.ActionSummary{
				Name:        actionName,
				Description: &comments,
				File: action,
			},
		}
		c.JSON(200, response)
	}
	return fn
}

func listActions(a *map[string]string) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		actions := []string{}
		for action, _ := range *a {
			actions = append(actions, action)
		}
		log.Debug("Actions: ", actions)

		response := model.ListActionResponseContent{
			Actions: actions,
		}
		c.JSON(200, response)
	}
	return fn
}