package store

import (
	"context"
	"strings"
	"time"

	"github.com/awlsring/action-runner/data"
	"github.com/awlsring/surreal-db-client/surreal"
)

type PlaySurreal struct {
	db *surreal.Surreal
}

func NewPlayDao(s *surreal.Surreal) *PlaySurreal {
	return &PlaySurreal{
		db: s,
	}
}

func surrealIdFromPlayId(id string) string {
	return strings.Replace(id, "p-", "play:", -1)
}

func (dao *PlaySurreal) Create(p *data.PlayEntity) error {
	sid := surrealIdFromPlayId(p.ID)
	p.ID = sid

	for i, task := range p.Tasks {
		p.Tasks[i] = surrealIdFromTaskId(task)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	return dao.db.CreateItem(ctx, sid, ResourceToEntry(p))
}