package vane

import (
	"fmt"
	"sync"

	"github.com/miekg/coredns/middleware/pkg/iptree"
	"github.com/miekg/coredns/middleware/vane/models"

	"github.com/astaxie/beego/orm"
)

type MySQLoader struct {
	lock sync.Mutex
}

func (l *MySQLoader) Init(nsd string) error {
	return models.InitDB(nsd)
}

func (l *MySQLoader) LoadAll() (engine *Engine, err error) {
	engine = NewEngine(l)
	wg := new(sync.WaitGroup)

	errChan := make(chan error, 16)

	o := orm.NewOrm()
	err = o.Begin()
	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			o.Rollback()
		} else {
			o.Commit()
		}
	}()

	wg.Add(1)
	fmt.Println("loadClientSet")
	go func() {
		l.loadClientSet(o, engine, errChan)
		wg.Done()
	}()

	fmt.Println("wait")
	wg.Wait()
	fmt.Println("loaded")

	select {
	case err = <-errChan:
		if err != nil {
			close(errChan)
			return
		}
	default:
	}

	fmt.Println("return")
	return
}

func (l *MySQLoader) loadClientSet(o orm.Ormer, engine *Engine, errChan chan<- error) (err error) {
	var (
		cs        models.ClientSetView
		clientSet []interface{}
	)

	defer func() {
		if err != nil && errChan != nil {
			errChan <- err
		}
	}()

	clientSet, err = models.GetAllClientSetView(o, nil, nil, nil, nil, 0, -1)
	if err != nil {
		if errChan != nil {
			errChan <- err
		}
		return
	}

	ipt := iptree.New()

	for _, c := range clientSet {
		cs = c.(models.ClientSetView)
		err = ipt.AddByString(fmt.Sprintf("%s/%d", cs.Ipnet, int(cs.Mask)), cs.IpnetId)
		if err != nil {
			return
		}
	}

	l.lock.Lock()
	defer l.lock.Unlock()
	engine.ClientSet = ipt
	return
}
