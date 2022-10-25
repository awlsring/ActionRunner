package store

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"github.com/awlsring/action-runner/execution"
	model "github.com/awlsring/dws-action-runner"
	"github.com/awlsring/surreal-db-client/surreal"
)

type ExecutionSurreal struct {
	db *surreal.Surreal
}

func NewEx(s *surreal.Surreal) *ExecutionSurreal {
	return &ExecutionSurreal{
		db: s,
	}
}

func ResourceToEntry(resource interface{}) map[string]interface{} {
	var entry map[string]interface{}
    inrec, _ := json.Marshal(resource)
    json.Unmarshal(inrec, &entry)
	return entry
}

func surrealIdFromId(id string) string {
	return strings.Replace(id, "e-", "execution:", -1)
}

func surrealIdToPlayId(id string) string {
	return strings.Replace(id, "play:", "p-", -1)
}

func surrealIdToTaskId(id string) string {
	return strings.Replace(id, "task:", "t-", -1)
}

func SurrealIdToMachineId(id string) string {
	return strings.Replace(id, "machine:", "m-", -1)
}

func SurrealIdFromMachineId(id string) string {
	return strings.Replace(id, "m-", "machine:", -1)
}

func UnmarshalGet(data interface{}, v interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err 
	}
	err = json.Unmarshal(b, &v)
	if err != nil {
		return err 
	}
	return nil
}

func (dao *ExecutionSurreal) GenerateId() string {
	id := strings.Replace(uuid.New().String(), "-", "", -1)
	return fmt.Sprintf("e-%v", id)
}

func (dao *ExecutionSurreal) Create(e *model.ExecutionSummary) error {
	sid := surrealIdFromId(e.Id)
	e.Id = sid
	wg := new(sync.WaitGroup)
	wg.Add(len(e.Targets))
	for _, target := range e.Targets {
		go func(wg *sync.WaitGroup, t string) {
			defer wg.Done()
			go dao.db.RelateRecords(context.Background(), sid, t, "executedOn")
		}(wg, target)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	return dao.db.CreateItem(ctx, sid, ResourceToEntry(e))
}

func (dao *ExecutionSurreal) Get(id string) (*model.ExecutionSummary, error) {
	sid := surrealIdFromId(id)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	blob, err := dao.db.GetItem(ctx, sid)
	if err != nil {
		log.Error(err)
		return nil, err 
	}

	var ex *model.ExecutionSummary
	err = UnmarshalGet(blob, &ex)
	if err != nil {
		log.Error(err)
		return nil, err 
	}
	ex.Id = id

	for i, play := range ex.Plays {
		ex.Plays[i] = surrealIdToPlayId(play)
	}

	stats := []model.ExecutionStats{}
	for _, stat := range ex.Stats {
		stat.Machine = SurrealIdToMachineId(stat.Machine)
		stats = append(stats, stat)
	}
	ex.Stats = stats
	
	return ex, nil
}

func (dao *ExecutionSurreal) getTask(id string, wg *sync.WaitGroup, c chan <- model.Task) {
	defer wg.Done()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	tblob, err := dao.db.GetItem(ctx, id)
	if err != nil {
		log.Error(err)
	}
	var t model.Task
	err = UnmarshalGet(tblob, &t)
	if err != nil {
		log.Error(err)
	}
	for i, target := range t.Targets {
		t.Targets[i].Machine = SurrealIdToMachineId(target.Machine)
	}

	t.Id = surrealIdToTaskId(t.Id)

	c <- t
}

func (dao *ExecutionSurreal) getPlay(id string,  wg *sync.WaitGroup, c chan <- model.PlayExtended) {
	defer wg.Done()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	pblob, err := dao.db.GetItem(ctx, surrealIdFromPlayId(id))
	if err != nil {
		log.Error(err)
	}
	var p *execution.Play
	err = UnmarshalGet(pblob, &p)
	if err != nil {
		log.Error(err)
	}

	pid := surrealIdToPlayId(id)
	playExtended := model.PlayExtended{
		Id: pid,
		Name: p.Name,
		StartTime: p.StartTime,
		EndTime: p.EndTime,
	}

	tasks := []model.Task{}
	taskGroup := new(sync.WaitGroup)
	taskGroup.Add(len(p.Tasks))
	taskChan := make(chan model.Task)
	for _, task := range p.Tasks {
		go dao.getTask(task, taskGroup, taskChan)
	}

	go func() {
		taskGroup.Wait()
		close(taskChan)
	}()

	for task := range taskChan {
		tasks = append(tasks, task)
	}

	playExtended.Tasks = tasks

	c <- playExtended
}

func (dao *ExecutionSurreal) GetDetailed(id string) (*model.DetailedExecutionSummary, error) {
	blob, err := dao.Get(id)
	if err != nil {
		log.Error(err)
		return nil, err 
	}

	var ex *model.ExecutionSummary
	err = UnmarshalGet(blob, &ex)
	if err != nil {
		log.Error(err)
		return nil, err 
	}
	ex.Id = id
	
	plays := []model.PlayExtended{}
	playGroup := new(sync.WaitGroup)
	playGroup.Add(len(ex.Plays))
	playChan := make(chan model.PlayExtended)
	for _, play := range ex.Plays {
		go dao.getPlay(play, playGroup, playChan)
	}

	go func() {
		playGroup.Wait()
		close(playChan)
	}()

	for play := range playChan {
		plays = append(plays, play)
	}

	exDetailed := model.DetailedExecutionSummary{
		Id: id,
		Action: ex.Action,
		Status: ex.Status,
		StartTime: ex.StartTime,
		EndTime: ex.EndTime,
		Stats: ex.Stats,
		Plays: plays,
		Targets: ex.Targets,
	}

	return &exDetailed, nil
}

func (dao *ExecutionSurreal) Update(e *model.ExecutionSummary) error {
	e.Id = surrealIdFromId(e.Id)

	plays := []string{}
	for _, play := range e.Plays {
		p := surrealIdFromPlayId(play)
		plays = append(plays, p)
	}

	e.Plays = plays

	entry := ResourceToEntry(e)
	log.Debugf("Storing entry: %v", entry)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	return dao.db.UpdateItem(ctx, e.Id, entry)
}

func (dao *ExecutionSurreal) Delete(id string) error {
	return nil
}

func (dao *ExecutionSurreal) List() ([]*model.ExecutionSummary, error) {
	return nil, nil
}