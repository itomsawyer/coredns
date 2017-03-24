package supermq

import (
	"fmt"
	nsqlib "github.com/nsqio/go-nsq"
)

type MProducer struct {
	poller chan *nsqlib.Producer
	pdlist []*nsqlib.Producer
}

func NewMProducer() *MProducer {
	return &MProducer{
		poller: make(chan *nsqlib.Producer, 10),
		pdlist: make([]*nsqlib.Producer, 0, 10),
	}
}

func (mp *MProducer) SetLogger(l logger, ll int) {
	for _, p := range mp.pdlist {
		p.SetLogger(l, nsqlib.LogLevel(ll))
	}
}

func (mp *MProducer) Stop() {
	for _, p := range mp.pdlist {
		p.Stop()
	}
}

func (m *MProducer) AddRx(addr string) error {
	config := nsqlib.NewConfig()
	w, err := nsqlib.NewProducer(addr, config)
	if err != nil {
		return err
	}
	if err := w.Ping(); err != nil {
		return err
	}
	m.pdlist = append(m.pdlist, w)
	m.poller <- w
	return nil
}

func (m *MProducer) MultiPublish(topic string, body [][]byte) (err error) {
	if len(m.pdlist) == 0 {
		return fmt.Errorf("no available producer")
	}

	tryCount := 0
	for onePder := range m.poller {
		tryCount++
		m.poller <- onePder
		err = onePder.MultiPublish(topic, body)
		if err == nil {
			return
		}
		if tryCount > len(m.pdlist) {
			break
		}
	}
	return
}
