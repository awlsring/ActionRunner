package data

import (
	model "github.com/awlsring/dws-action-runner"
)

type ExecutionDao interface {
	GenerateId() string 
	Create(e *model.ExecutionSummary) error
	Get(id string) (*model.ExecutionSummary, error)
	GetDetailed(id string) (*model.DetailedExecutionSummary, error)
	Update(e *model.ExecutionSummary) error
	Delete(id string) error
	List() ([]*model.ExecutionSummary, error) 
}

type PlayDao interface {
	Create(e *PlayEntity) error
}

type TaskDao interface {
	Create(e *model.Task) error
}

type PlayEntity struct {
	Name string `json:"name"`
	ID string `json:"id"`
	StartTime float32 `json:"startTime"`
	EndTime float32 `json:"endTime"`
	Tasks []string `json:"tasks,omitempty"`
}