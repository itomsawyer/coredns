package vane

import (
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

	data := new(MySQLData)
	err = db.loadClientSet(o, data)
	if err != nil {
		return nil
	}

	db.rwlock.Lock()
	defer db.rwlock.Unlock()
	db.data = data

	return nil
}

func (db *MySQLDB) GetClientSetID(ip string) (int, error) {
	id, _, err := db.data.clientSet.GetByString(ip)
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
		err := data.clientSet.AddByNet(cs.Ipnet, int(cs.Mask), cs.IpnetId)
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
