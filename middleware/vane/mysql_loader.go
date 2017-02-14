package vane

import (
	"fmt"
	"sync"

	"github.com/miekg/coredns/middleware/pkg/dmtree"
	"github.com/miekg/coredns/middleware/pkg/iptree"
	"github.com/miekg/coredns/middleware/proxy"
	"github.com/miekg/coredns/middleware/vane/models"

	"github.com/astaxie/beego/orm"
)

type MySQLoader struct{}

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
	go func() {
		LoadClientSet(o, engine, errChan)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		LoadNetlink(o, engine, errChan)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		LoadDomain(o, engine, errChan)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		LoadPolicy(o, engine, errChan)
		wg.Done()
	}()

	wg.Wait()

	select {
	case err = <-errChan:
		if err != nil {
			close(errChan)
			return
		}
	default:
	}

	return
}

func LoadClientSet(o orm.Ormer, engine *Engine, errChan chan<- error) (err error) {
	var (
		cs        models.ClientSetView
		clientSet []interface{}
		//clientWLSet []interface{}
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

	/*
		clientWLSet, err = models.GetAllClientSetWLView(o, nil, nil, nil, nil, 0, -1)
		if err != nil {
			return
		}
	*/

	ipt := iptree.New()

	for _, c := range clientSet {
		cs = c.(models.ClientSetView)
		err = ipt.AddByString(fmt.Sprintf("%s/%d", cs.Ipnet, int(cs.Mask)), cs.IpnetId)
		if err != nil {
			return
		}
	}

	/*
		for _, c := range clientWLSet {
			cs = c.(models.ClientSetWLView)
			err = ipt.AddByString(fmt.Sprintf("%s/%d", cs.Ipnet, int(cs.Mask)), cs.IpnetId)
			if err != nil {
				return
			}
		}
	*/

	engine.ClientSet = ipt
	return
}

func LoadNetlink(o orm.Ormer, engine *Engine, errChan chan<- error) (err error) {
	var (
		nl       models.NetLinkView
		netlinks []interface{}
	)

	defer func() {
		if err != nil && errChan != nil {
			errChan <- err
		}
	}()

	netlinks, err = models.GetAllNetLinkView(o, nil, nil, nil, nil, 0, -1)
	if err != nil {
		return
	}

	ipt := iptree.New()
	for _, n := range netlinks {
		nl = n.(models.NetLinkView)
		err = ipt.AddByString(fmt.Sprintf("%s/%d", nl.Ipnet, int(nl.Mask)), nl.NetLinkId)
		if err != nil {
			return
		}
	}

	engine.NetLink = ipt
	return
}

func LoadDomain(o orm.Ormer, engine *Engine, errChan chan<- error) (err error) {
	var (
		ds        models.DomainView
		domainSet []interface{}
	)

	defer func() {
		if err != nil && errChan != nil {
			errChan <- err
		}
	}()

	domainSet, err = models.GetAllDomainView(o, nil, nil, nil, nil, 0, -1)
	if err != nil {
		return
	}

	dmt := new(dmtree.DmTree)
	dmt.Insert(".", 1) //TODO: remove this ugly hard code
	for _, d := range domainSet {
		ds = d.(models.DomainView)
		dmt.Insert(ds.Domain, ds.DomainId)
	}

	engine.DomainSet = dmt
	return
}

func LoadPolicy(o orm.Ormer, engine *Engine, errChan chan<- error) (err error) {
	var (
		pv          models.PolicyView
		policyViews []interface{}
	)

	defer func() {
		if err != nil && errChan != nil {
			errChan <- err
		}
	}()

	policyViews, err = models.GetAllPolicyView(o, nil, nil, nil, nil, 0, -1)
	if err != nil {
		return
	}

	policySet := models.PolicySet{}
	for _, v := range policyViews {
		pv = v.(models.PolicyView)
		policySet.Add(pv.PolicyId, &pv)
	}

	ps := make(map[int]*proxy.Proxy, 16)
	for k, s := range policySet {
		hosts := s.Hosts()
		if hosts == nil {
			continue
		}

		proxy := proxy.NewLookup2(hosts)
		if k == 1 {
			fmt.Println(proxy)
		}
		ps[k] = &proxy
	}

	engine.Proxy = ps
	return
}
