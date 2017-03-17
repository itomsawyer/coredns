package engine

import (
	"fmt"
	"os"
	"testing"

	"github.com/astaxie/beego/orm"
)

func TestMain(m *testing.M) {
	err := orm.RegisterDataBase("default", "mysql", "root:@tcp(127.0.0.1:3306)/iwg")
	orm.Debug = true
	if err != nil {
		fmt.Println("cannot launch testing due to init error:", err)
		fmt.Println("please modify middleware/vane/models/init_test.go for db connection")
		os.Exit(-1)
	}

	os.Exit(m.Run())
}
