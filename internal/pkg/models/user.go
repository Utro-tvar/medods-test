package models

import "net"

type User struct {
	GUID string
	IP   net.IP
}
