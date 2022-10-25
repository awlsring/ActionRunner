package watcher

import (
	"fmt"
	"os"
	"strings"

	"github.com/awlsring/action-runner/utils"
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
)

func inActions(path string, f string) bool {

	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if file.Name() == f {
			return true
		}
	}
	return true
}

func WatchActions(path string, actions *map[string]string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Debug("Event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Debug("modified file:", event.Name)
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					if strings.HasSuffix(event.Name, ".yaml") || strings.HasSuffix(event.Name, ".yml") {
						addToMap(actions, event.Name, path)
					}	
				}
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					if strings.HasSuffix(event.Name, ".yaml") || strings.HasSuffix(event.Name, ".yml") {
						removeFromMap(actions, event.Name, path)
					}	
				}
				if event.Op&fsnotify.Rename == fsnotify.Rename {
					if strings.HasSuffix(event.Name, ".yaml") || strings.HasSuffix(event.Name, ".yml") {
						removeFromMap(actions, event.Name, path)
					}	
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Debug("error:", err)
			}
		}
	}()

	err = watcher.Add(path)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func removeFromMap(a *map[string]string, name string, path string) {
	log.Debug("removed file:", name)
	file := strings.Replace(name, fmt.Sprintf("%s/", path), "", 1)
	action := utils.ToCamelCase(strings.Split(file, ".")[0])
	old := *a
	delete(old, action)
	a = &old
	log.Debugf("removed action: %s - %s", action, file)
}

func addToMap(a *map[string]string, name string, path string) {
	log.Debug("created file:", name)
	file := strings.Replace(name, fmt.Sprintf("%s/", path), "", 1)
	action := utils.ToCamelCase(strings.Split(file, ".")[0])
	old := *a
	old[action] = file
	a = &old
	log.Debugf("added action: %s - %s", action, file)
}

func renameInMap(a *map[string]string, name string, path string) {
	log.Debug("renamed file:", name)
	file := strings.Replace(name, fmt.Sprintf("%s/", path), "", 1)
	action := utils.ToCamelCase(strings.Split(file, ".")[0])
	old := *a
	old[action] = file
	a = &old
	log.Debugf("added action: %s - %s", action, file)
}