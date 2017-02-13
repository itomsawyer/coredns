package models

import (
	"fmt"
	"os"
	"testing"

	"github.com/astaxie/beego/orm"
)

func TestMain(m *testing.M) {
	err := InitDB("root:root@tcp(127.0.0.1:3306)/iwg")
	orm.Debug = true
	if err != nil {
		fmt.Println("cannot launch testing due to init error:", err)
		os.Exit(-1)
	}

	os.Exit(m.Run())
}
