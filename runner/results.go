package runner

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/apenella/go-ansible/pkg/stdoutcallback/results"
	model "github.com/awlsring/action-runner-model"
	"github.com/awlsring/action-runner/data"
)

var AnsibleTimestamp = "2006-01-02T15:04:05.000000Z"

func formId(t string, i string) string {
	id := strings.Replace(i, "-", "", -1)
	return fmt.Sprintf("%v-%v", t, id)
}

func ansibleStartAndEndToUnix(s string, e string) (float32, float32, error) {
	start, err := time.Parse(AnsibleTimestamp, s)
	if err != nil {
		return 0, 0, err
	}
	end, err := time.Parse(AnsibleTimestamp, e)
	if err != nil {
		return 0, 0, err
	}
	return float32(start.UTC().Unix()), float32(end.UTC().Unix()), nil
}

func (a *AnsibleRunner) storeResults(plays []results.AnsiblePlaybookJSONResultsPlay) ([]string, error) {
	playGroup := new(sync.WaitGroup)
	playGroup.Add(len(plays))
	pids := []string{}
	for _, play := range plays {
		pid := formId("p", play.Play.Id)
		pids = append(pids, pid)
		go a.storePlay(playGroup, pid, play)
	}

	playGroup.Wait()

	return pids, nil
}

func (a *AnsibleRunner) storePlay(wg *sync.WaitGroup, pid string, play results.AnsiblePlaybookJSONResultsPlay) {
	defer wg.Done()

	taskGroup := new(sync.WaitGroup)
	taskGroup.Add(len(play.Tasks))
	tids := []string{}
	for _, task := range play.Tasks {
		tid := formId("t", task.Task.Id)
		tids = append(tids, tid)
		go a.storeTask(taskGroup, tid, task)
	}

	start, end, err := ansibleStartAndEndToUnix(play.Play.Duration.Start, play.Play.Duration.End)
	if err != nil {
		return
	}
	p := &data.PlayEntity{
		Name: play.Play.Name,
		ID: pid,
		StartTime: start,
		EndTime: end,
		Tasks: tids,
	}
	a.Repository.PlayDao.Create(p)
	taskGroup.Wait()
}

func (a *AnsibleRunner) storeTask(wg *sync.WaitGroup, tid string, task results.AnsiblePlaybookJSONResultsPlayTask) {
	defer wg.Done()
	start, end, err := ansibleStartAndEndToUnix(task.Task.Duration.Start, task.Task.Duration.End)
	if err != nil {
		return
	}

	targets := []model.TargetResult{}
	for target, stats := range task.Hosts {
		targets = append(targets, model.TargetResult{
			Machine: target,
			Action: stats.Action,
			Changed: &stats.Changed,
			Stdout: stats.StdoutLines,
			Stderr: stats.StderrLines,
			Failed: &stats.Failed,
			FailedWhenResult: &stats.FailedWhenResult,
			Skipped: &stats.Skipped,
			Unreachable: &stats.Unreachable,
		})
	}

	t := &model.Task{
		Name: task.Task.Name,
		Id: tid,
		StartTime: start,
		EndTime: end,
		Targets: targets,
	}

	a.Repository.TaskDao.Create(t)
}
