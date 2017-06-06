package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/coredns/coredns/utils/mkconf/conf"
)

func main() {
	db := flag.String("db", "root:@tcp(127.0.0.1:3306)/iwg", "iwg db nsd")
	of := flag.String("o", "auto.conf", "output file path")
	flag.Parse()

	c := conf.NewConf()
	if err := c.BuildConfFromDB(*db); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	f, err := os.OpenFile(*of, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(0666))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	data, err := json.Marshal(c)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	f.Write(data)
	f.Close()
}
