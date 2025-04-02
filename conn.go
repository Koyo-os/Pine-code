package main

import "net"

type Conn struct {
	*net.Conn
}
