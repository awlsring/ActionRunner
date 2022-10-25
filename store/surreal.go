package store

import (
	"log"
	"sync"

	"github.com/awlsring/surreal-db-client/surreal"
)

// singleton may be bad idea depending on how surreal handles the conns

var lock = &sync.Mutex{}

var db *surreal.Surreal

func New(cfg surreal.SurrealConfig) *surreal.Surreal {
	db, err := surreal.New(cfg)
	if err != nil {
		log.Fatalln(err)
	}
	return db
}

func GetSurrealClient(cfg surreal.SurrealConfig) *surreal.Surreal {
    if db == nil {
        lock.Lock()
        defer lock.Unlock()
        if db == nil {
            log.Println("Creating instance of SurrealDB client.")
            db, err := surreal.New(cfg)
			if err != nil {
				log.Fatalln(err)
			}
			return db

		} else {
			log.Println("Returned exisiting SurrealDB client")
		}
	}

    return db
}