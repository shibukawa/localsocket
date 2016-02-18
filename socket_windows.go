package localsocket

import (
	"gopkg.in/natefinch/npipe.v2"
)

// LocalSocket is a socket to communicate with other processes in the same box.
// LocalSocket satisfies io.Reader, io.Writer, net.Conn interfaces.
type LocalSocket struct {
	*npipe.PipeConn
}

// NewLocalSocket creates LocalSocket instance
func NewLocalSocket(pathName string) (*LocalSocket, error) {
	conn, err := npipe.Dial(`\\.\pipe\` + pathName)
	if err != nil {
		return nil, err
	}
	return newLocalSocket(conn), nil
}

func newLocalSocket(conn *npipe.PipeConn) *LocalSocket {
	return &LocalSocket{
		PipeConn: conn,
	}
}
