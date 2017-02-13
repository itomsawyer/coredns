package vane

import (
	"sync"

	"github.com/miekg/coredns/middleware/vane/models"
)

type MySQLDB struct {
	rwlock sync.RWMutex
	data   interface{}
}

func (db *MySQLDB) Open(nsd string) error {
	return models.InitDB(nsd)
}

func (db *MySQLDB) Load() error {
	db.rwlock.Lock()
	defer db.rwlock.Unlock()

	return nil
}
