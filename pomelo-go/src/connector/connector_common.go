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
	Send([]byte) (int, error)
	SendBatch(...[]byte)
	Receive([]byte) (int, error)
	Disconnect() int
}

type Connector interface {
	Start()
	Decode([]byte) (interface{}, error)
	Encode(string, string, string) ([]byte, error)
}
