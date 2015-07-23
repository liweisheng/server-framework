package agent

import (
// "fmt"
// "net"
)

const (
	AG_INIT  = iota
	AG_START = iota
	AG_CLOSE = iota
)

type Agent interface {
	IsMaster() bool
}
