package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/coredns/coredns/utils/mkconf/conf"
)

func main() {
	var corednsHost []string
	var corednsDomainPool []string
	db := flag.String("db", "root:@tcp(127.0.0.1:3306)/iwg", "iwg db nsd")
	of := flag.String("o", "auto.conf", "output file path")
	chost := flag.String("coredns", "[]", "coredns hosts formated as json string array. ep. [\"10.0.0.1\",\"10.0.0.2\"]")
	cdp := flag.String("dp", "[]", "coredns domain pool name formated as json string array. ep. [\"httpdns\",\"dm_hn_test\"]")
	flag.Parse()

	if err := json.Unmarshal([]byte(*chost), &corednsHost); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := json.Unmarshal([]byte(*cdp), &corednsDomainPool); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	c := conf.NewConf()

	for _, h := range corednsHost {
		c.AddCorednsHost(h)
	}

	for _, dp := range corednsDomainPool {
		c.AddCorednsDomainPool(dp)
	}

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
