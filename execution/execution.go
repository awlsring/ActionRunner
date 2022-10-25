package execution

import (
	"fmt"
	"strings"
	"time"

	"github.com/apenella/go-ansible/pkg/stdoutcallback/results"
	model "github.com/awlsring/dws-action-runner"
	"github.com/google/uuid"
)

var TimestampMicro = "2006-01-02T15:04:05.000000Z"

type Status string

const (
	Success Status = "SUCCESS"
	Failure Status = "FAILURE"
	Running Status = "RUNNING"
	Timeout Status = "TIMEOUT"
	SuccessWithErrors Status = "SUCCESS_WITH_ERRORS"
)

type ExecutionStats struct {
	Machine string `json:"machine"`
	Changed int `json:"changed"`
	Failures int `json:"failures"`
	Ignored int `json:"ignored"`
	Ok int `json:"ok"`
	Rescued int `json:"rescued"`
	Skipped int `json:"skipped"`
	Unreachable int `json:"unreachable"`
}

type Execution struct {
	ID string `json:"id"`
	Stats []ExecutionStats `json:"stats,omitempty"`
	StartTime int64 `json:"startTime"`
	EndTime int64 `json:"endTime"`
	Status model.ExecutionStatus `json:"status"`
	Action string `json:"action"`
	Plays []string `json:"plays,omitempty"`
}

func generateExecutionId() string {
	id := strings.Replace(uuid.New().String(), "-", "", -1)
	return fmt.Sprintf("execution:%v", id)
}

func New(action string) *Execution {
	return &Execution{
		ID: generateExecutionId(),
		StartTime: time.Now().Unix(),
		Status: model.RUNNING,
		Action: action,
	}
}

func (e *Execution) Close(
	r model.ExecutionStatus,
	plays []string,
	stats map[string]*results.AnsiblePlaybookJSONResultsStats,
) {
	e.EndTime = time.Now().Unix()
	e.Status = r
	e.Plays = plays

	statList := []ExecutionStats{}
	for host, stats := range stats {
		s := ExecutionStats{
			Machine: host,
			Ignored: stats.Ignored,
			Changed: stats.Changed,
			Failures: stats.Failures,
			Ok: stats.Ok,
			Rescued: stats.Rescued,
			Skipped: stats.Skipped,
			Unreachable: stats.Unreachable,
		}
		statList = append(statList, s)
	}

	e.Stats = statList
}
