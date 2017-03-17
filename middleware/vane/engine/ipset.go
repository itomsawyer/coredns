package engine

import (
	"net"
)

type ClientSet struct {
	ID    int
	IPNet *net.IPNet
	Name  string
}

type NetLink struct {
	ID     int
	IPNet  *net.IPNet
	Isp    string
	Region string
}
