package models

import (
	"github.com/astaxie/beego/orm"

	_ "github.com/go-sql-driver/mysql"
)

func InitDB(nsd string) error {
	//nsd = nsd + "?charset=utf8&loc=Asia%20FShanghai&tx_isolation=%27REPEATABLE-READ%27"
	return orm.RegisterDataBase("default", "mysql", nsd)
}

func RegisterDB(dbname, driver string, nsd string) error {
	if dbname == "" {
		dbname = "default"
	}
	//nsd = nsd + "?charset=utf8&loc=Asia%20FShanghai&tx_isolation=%27REPEATABLE-READ%27"
	return orm.RegisterDataBase(dbname, driver, nsd)
}
