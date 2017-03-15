package supermq

import (
	"fmt"
	//	nsqlib "github.com/nsqio/go-nsq"
	"testing"
)

type MyHandler struct {
	num int
}

func (h *MyHandler) HandleMessage(message *Message) error {
	fmt.Printf("[M]%s: %s\n", message.NSQDAddress, string(message.Body))
	return nil
}

func TestTconsumer(t *testing.T) {
	h := MyHandler{}
	r, _ := NewTconsumer("dlresult", "ch1", h.HandleMessage)
	fmt.Println(r.ConnectToNSQLookupds(
		[]string{"192.168.30.192:4161", "192.168.30.221:4161"}))
	//fmt.Println(r.ConnectToNSQDs(
	//	[]string{"192.168.30.192:4150", "192.168.30.221:4150"}))

	stop := make(chan int)
	<-stop
}
