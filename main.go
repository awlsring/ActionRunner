package main

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/awlsring/action-runner/api"
	"github.com/awlsring/action-runner/config"
	"github.com/awlsring/action-runner/runner"
	"github.com/awlsring/action-runner/store"
	"github.com/awlsring/action-runner/utils"
	"github.com/awlsring/action-runner/watcher"
)

func getActions(path string) (*map[string]string, error) {
	actions := map[string]string{}

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if !file.IsDir() {
			if strings.Contains(file.Name(), ".yaml") || strings.Contains(file.Name(), ".yml") {
				action := strings.Split(file.Name(), ".")[0]
				actions[utils.ToCamelCase(action)] = file.Name()
			}
		}
	}
	return &actions, nil
}

func main() {
	log.SetLevel(log.DebugLevel)
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalln(err)
	}
	
	a, err := getActions(cfg.Runner.PlaybookDir)
	if err != nil {
		log.Fatalln(err)
	}

	conn := store.NewSurrealStore(cfg.Database.Surreal)

	r, err := runner.New(cfg.Runner, conn)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(a, r)

	go watcher.WatchActions(cfg.Runner.PlaybookDir, a)
	
	api.Start(cfg.Api, a, r, conn)

}