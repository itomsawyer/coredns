package engine

import (
	"testing"
	"time"

	"github.com/hashicorp/golang-lru"
)

func TestLRU(t *testing.T) {
	lru, _ := lru.New(10)
	lru.Add("1.1.1.1/1", 1)
	if v, ok := lru.Get("1.1.1.1/1"); !ok {
		t.Error("unexpected nout found")
	} else if v.(int) != 1 {
		t.Error("value error")
	}
}

func TestLinkManagerCreate(t *testing.T) {
	lm, err := NewLinkManager(-1)
	if err != nil {
		t.Error(err)
	}

	if ls, ok := lm.GetLink("1.1.1.1", "ot1"); ok || ls != nil {
		t.Errorf("unexpected found")
	}

	if ls, ok := lm.GetLink("1.1.1.1", "ot1"); (!ok) || ls == nil {
		t.Errorf("unexpected not found")
	} else if ls.Status != LinkStatusUnknown {
		t.Errorf("expected link status %d, get %d", LinkStatusUnknown, ls.Status)
	}
}

func TestLinkStatusExpire(t *testing.T) {
	ls := LinkStatus{}
	now := time.Now()
	ls.SetTTL(3*time.Second, now.Add(-1*time.Second))

	if left, ok := ls.IsExpire(now); ok {
		t.Errorf("unexpected expired")
	} else if left != 2*time.Second {
		t.Errorf("expect ttl left %s, get %s", time.Second*2, left)
	}

	ls.SetTTL(1*time.Second, now.Add(-1*time.Second))
	if left, ok := ls.IsExpire(now); ok {
		t.Errorf("unexpected expired")
	} else if left != 0 {
		t.Errorf("expect ttl left 0,  get %s", left)
	}

	ls.SetTTL(1*time.Second, now.Add(-2*time.Second))
	if left, ok := ls.IsExpire(now); !ok {
		t.Errorf("unexpected not expired")
	} else if left != (time.Second * -1) {
		t.Errorf("expect ttl left %s,  get %s", -1*time.Second, left)
	}
}

func TestLinkManagerSend(t *testing.T) {
	lm, err := NewLinkManager(100)
	if err != nil {
		t.Error(err)
		return
	}
	if err := lm.RegisterSender([]string{"192.168.30.192:4150"}, "dst2otlnk"); err != nil {
		t.Error(err)
		return
	}

	if err := lm.RegisterReader([]string{"192.168.30.192:4161"}, "dlresult", "channel"); err != nil {
		t.Error(err)
		return
	}

	//for {
	//	lm.GetLink("1.1.1.1", "ot1")
	//	time.Sleep(1 * time.Second)
	//}
}
