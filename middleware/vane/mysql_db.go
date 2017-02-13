package vane

import (
	"errors"
	"fmt"
	"net"
	"sync"

	"github.com/miekg/coredns/middleware/pkg/iptree"
	"github.com/miekg/coredns/middleware/vane/models"

	"github.com/astaxie/beego/orm"
)

type MySQLDB struct {
	rwlock sync.RWMutex
	wg     sync.WaitGroup
	data   *MySQLData
}

func (db *MySQLDB) Open(nsd string) error {
	return models.InitDB(nsd)
}

func (db *MySQLDB) Load() error {
	var err error

	o := orm.NewOrm()
	err = o.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			o.Rollback()
		} else {
			o.Commit()
		}
	}()

	data := NewMySQLData()
	err = db.loadClientSet(o, data)
	if err != nil {
		return nil
	}

	db.rwlock.Lock()
	defer db.rwlock.Unlock()
	db.data = data

	return nil
}
func (db *MySQLDB) GetClientSetID(ip net.IP) (int, error) {
	if ip == nil {
		return 0, errors.New("client ip is nil")
	}

	if db.data == nil {
		return 0, errors.New("db does not loaded")
	}

	id, found, err := db.data.clientSet.Get(ip)
	if !found {
		return 0, errors.New("ClientSet not found")
	}
	return id, err
}

func (db *MySQLDB) loadClientSet(o orm.Ormer, data *MySQLData) error {
	var cs models.ClientSetView
	clientSet, err := models.GetAllClientSetView(o, nil, nil, nil, nil, 0, -1)
	if err != nil {
		return err
	}

	for _, c := range clientSet {
		cs = c.(models.ClientSetView)
		err := data.clientSet.AddByString(fmt.Sprintf("%s/%d", cs.Ipnet, int(cs.Mask)), cs.IpnetId)
		if err != nil {
			return err
		}
	}

	return nil
}

type MySQLData struct {
	clientSet *iptree.IPTree
}

func NewMySQLData() *MySQLData {
	return &MySQLData{
		clientSet: iptree.New(),
	}
}
