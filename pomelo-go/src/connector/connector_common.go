//author:liweisheng date:2015/07/07

package connector

const (
	ST_INITED = iota
	ST_CLOSED = iota
)

type Socket interface {
	ID() int32
	Socket() interface{}
	RemoteAddress() map[string]string
}
