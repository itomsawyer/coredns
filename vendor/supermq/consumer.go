package supermq

import (
	nsqlib "github.com/nsqio/go-nsq"
	"strings"
	"time"
)

type Tconsumer struct {
	*nsqlib.Consumer
}

func (t *Tconsumer) Stop() {
	t.Consumer.Stop()
}

func (t *Tconsumer) SetLogger(l logger, ll int) {
	t.Consumer.SetLogger(l, nsqlib.LogLevel(ll))
}

type Message struct {
	*nsqlib.Message
}

type Handle func(message *Message) error

func (h Handle) HandleMessage(message *nsqlib.Message) error {
	return h(&Message{message})
}

func NewTconsumer(topic string, channel string, h Handle) (*Tconsumer, error) {
	if !strings.HasSuffix(channel, "#ephemeral") {
		channel += "#ephemeral"
	}
	config := nsqlib.NewConfig()
	config.LookupdPollInterval = time.Second * 30
	r, err := nsqlib.NewConsumer(topic, channel, config)
	if err != nil {
		return nil, err
	}
	//concurrent handler
	r.ChangeMaxInFlight(2500)
	r.AddConcurrentHandlers(h, 10)
	return &Tconsumer{r}, nil
}
