package supermq

import (
	"fmt"
	nsqlib "github.com/nsqio/go-nsq"
)

type MProducer struct {
	pdChan chan *nsqlib.Producer
}

func NewMProducer() *MProducer {
	return &MProducer{
		pdChan: make(chan *nsqlib.Producer, 10),
	}
}

func (mp *MProducer) Stop() {
	for p := range mp.pdChan {
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
	m.pdChan <- w
	return nil
}

func (m *MProducer) MultiPublish(topic string, body [][]byte) (err error) {
	for i := 0; i < len(m.pdChan); i++ {
		onePder := <-m.pdChan
		err = onePder.MultiPublish(topic, body)
		m.pdChan <- onePder
		if err == nil {
			return
		} else {
			fmt.Printf("publish to %s failed, try next\n", onePder.String())
		}
	}
	return
}
