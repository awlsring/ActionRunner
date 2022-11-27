package store

import (
	"context"
	"strings"
	"time"

	model "github.com/awlsring/action-runner-model"
	"github.com/awlsring/surreal-db-client/surreal"
)

type TaskSurreal struct {
	db *surreal.Surreal
}

func NewTaskDao(s *surreal.Surreal) *TaskSurreal {
	return &TaskSurreal{
		db: s,
	}
}

func surrealIdFromTaskId(id string) string {
	return strings.Replace(id, "t-", "task:", -1)
}

func (dao *TaskSurreal) Create(t *model.Task) error {
	sid := surrealIdFromTaskId(t.Id)
	t.Id = sid
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	return dao.db.CreateItem(ctx, sid, ResourceToEntry(t))
}

