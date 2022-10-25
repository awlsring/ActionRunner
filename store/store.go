package store

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/awlsring/action-runner/data"
	"github.com/awlsring/surreal-db-client/surreal"
)

var DB *surreal.Surreal

type SurrealStore struct {
	ExecutionDao *ExecutionSurreal
	PlayDao *PlaySurreal
	TaskDao *TaskSurreal
}

func maintainConn(cfg surreal.SurrealConfig) {
	for {
		time.Sleep(5 * time.Second)
		log.Debug("Checking connection to Surreal")
		ctx, _ := context.WithTimeout(context.Background(), time.Second*3)
		_, err := DB.GetItem(ctx, "s")
		if err != nil {
			log.Errorln(err)
			log.Debug("Reconnecting to Surreal")
			db, err := surreal.New(cfg)
			if err != nil {
				log.Errorln(err)
			} else {
				log.Debug("Reconnected to Surreal")
				DB = db
			}
		} else {
			log.Debug("Connection to Surreal is good")
		}
	}
}

func NewSurrealStore(cfg surreal.SurrealConfig) *data.ExecutionRespository {
	db, err := surreal.New(cfg)
	if err != nil {
		log.Fatalln(err)
	}
	DB = db

	go maintainConn(cfg)
	
	return &data.ExecutionRespository{
		ExecutionDao: NewEx(db),
		PlayDao: NewPlayDao(db),
		TaskDao: NewTaskDao(db),
	}
}