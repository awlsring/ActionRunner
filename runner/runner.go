package runner

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/apenella/go-ansible/pkg/execute"
	"github.com/apenella/go-ansible/pkg/options"
	"github.com/apenella/go-ansible/pkg/playbook"
	"github.com/apenella/go-ansible/pkg/stdoutcallback/results"
	"github.com/awlsring/action-runner/data"
	"github.com/awlsring/action-runner/inventory"
	"github.com/awlsring/action-runner/store"
	model "github.com/awlsring/dws-action-runner"
	"gopkg.in/yaml.v2"
)

type Config struct {
	User string `mapstructure:"ansibleUser"`
	ConnectionType string `mapstructure:"connectionType"`
	PlaybookDir string `mapstructure:"playbookSource"`
}

type AnsibleRunner struct {
	Config Config
	Queue chan *model.RunActionRequestContent
	ConnectionOptions *options.AnsibleConnectionOptions
	Repository *data.ExecutionRespository
}

func New(c Config, d *data.ExecutionRespository) (*AnsibleRunner, error) {
	connectionOptions := &options.AnsibleConnectionOptions{
		Connection: c.ConnectionType,
		User:       c.User,
	}
	
	q := make(chan *model.RunActionRequestContent)

	return &AnsibleRunner{
		Config: c,
		Queue: q,
		ConnectionOptions: connectionOptions,
		Repository: d,
	}, nil
}

type RunRequest struct {
	Action string
	Execution *model.ExecutionSummary
	Request *model.RunActionRequestContent
	Hosts []*inventory.Host
}

func (a *AnsibleRunner) BackgroundRunner(ch <-chan *RunRequest) {
    for request := range ch {
        go a.executeAction(request)
    }
}

func intToFloat(i int) *float32 {
	f := float32(i)
	return &f
}

func closeExecution(
	ex model.ExecutionSummary,
	r model.ExecutionStatus,
	plays []string,
	stats map[string]*results.AnsiblePlaybookJSONResultsStats,
) *model.ExecutionSummary {
	close := float32(time.Now().Unix())
	ex.EndTime = &close
	ex.Status = r
	ex.Plays = plays

	statList := []model.ExecutionStats{}
	for host, stats := range stats {
		s := model.ExecutionStats{
			Machine: host,
			Ignored: intToFloat(stats.Ignored),
			Changed: intToFloat(stats.Changed),
			Failures: intToFloat(stats.Failures),
			Ok: intToFloat(stats.Ok),
			Rescued: intToFloat(stats.Rescued),
			Skipped: intToFloat(stats.Skipped),
			Unreachable: intToFloat(stats.Unreachable),
		}
		statList = append(statList, s)
	}

	ex.Stats = statList
	return &ex
}

func (a *AnsibleRunner) executeAction(request *RunRequest) {
	inv, err := a.CreateEphemeralInventory(request.Hosts, request.Execution.Id)
	if err != nil {
		log.Println(err)
		return
	}

	playbookOptions := &playbook.AnsiblePlaybookOptions{
		Inventory: inv,
	}
	buff := new(bytes.Buffer)

	// Investigate a custom executor to build execution objects as they are run
	// Currently, the writer is used to read each line and build a json blob that can be parsed after execution.
	// Custom executor would likely need to read each line and build task objects as they are run.
	execute := execute.NewDefaultExecute(
		execute.WithWrite(io.Writer(buff)),
	)

	playbook := &playbook.AnsiblePlaybookCmd{
		Playbooks:         []string{fmt.Sprintf("%v/%v", a.Config.PlaybookDir, request.Action)},
		Exec:              execute,
		ConnectionOptions: a.ConnectionOptions,
		Options:           playbookOptions,
		StdoutCallback:    "json",
	}

	os.Setenv("ANSIBLE_HOST_KEY_CHECKING", "false")

	log.Debug("Running playbook")
	err = playbook.Run(context.TODO())
	if err != nil {
		fmt.Println(err.Error())
	}
	log.Debug("Playbook done")

	log.Debug("Parsing results")
	res, err := results.ParseJSONResultsStream(io.Reader(buff))
	if err != nil {
		log.Println(err)
		return
	}

	result := model.SUCCESS
	for _, stats := range res.Stats {
		if stats.Failures != 0 {
			result = model.FAILED
		}
	}

	a.removeEphemeralInventory(inv)
	playIds, err := a.storeResults(res.Plays)
	if err != nil {
		log.Println(err)
		return
	}

	ex := closeExecution(*request.Execution, result, playIds, res.Stats)

	a.Repository.ExecutionDao.Update(ex)
	log.Debug("Execution done for ", request.Execution.Id)
}



func (a AnsibleRunner) CreateEphemeralInventory(hosts []*inventory.Host, executionId string) (string, error) {
	hostMap := make(map[string]inventory.HostOptions)
	for _, host := range hosts {
		// vvv Probably make a function to regex entries from x- to name instead of doing in runner
		hostMap[store.SurrealIdFromMachineId(host.ID)] = inventory.HostOptions{
			AnsiblePort: host.Port,
			AsibleHost: host.IpAddress,
			AnsibleOsFamily: host.OsFamily,
		}
	}
	
	i := inventory.Inventory{
		Inventory: inventory.HostList{
			Hosts: hostMap,
		},
	}

	yamlData, err := yaml.Marshal(&i)
	if err != nil {
        return "", err
    }

	fileName := fmt.Sprintf("%v.yaml", executionId)
    err = os.WriteFile(fileName, yamlData, 0644)
    if err != nil {
        return "", err
    }

	return fileName, nil
}

func (a AnsibleRunner) removeEphemeralInventory(fileName string) error {
	return os.Remove(fileName)
}
