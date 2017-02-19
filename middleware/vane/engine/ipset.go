package engine

import (
	"net"
)

type ClientSet struct {
	Id    int
	IPNet *net.IPNet
	Name  string
}

type NetLink struct {
	Id     int
	IPNet  *net.IPNet
	Isp    string
	Region string
}
