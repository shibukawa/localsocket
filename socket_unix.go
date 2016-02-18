// +build !windows

package localsocket

import (
	"net"
	"os"
	"path/filepath"
	"time"
)

// LocalSocket is a socket to communicate with other processes in the same box.
// LocalSocket satisfies io.Reader, io.Writer, net.Conn interfaces.
type LocalSocket struct {
	conn *net.UnixConn
}

// Read reads data from the connection.
func (s *LocalSocket) Read(data []byte) (int, error) {
	return s.conn.Read(data)
}

// Write writes data to the connection.
func (s *LocalSocket) Write(data []byte) (int, error) {
	return s.conn.Write(data)
}

// Close closes the connection.
func (s *LocalSocket) Close() error {
	return s.conn.Close()
}

// LocalAddr returns the local network address.
// The Addr returned is shared by all invocations of LocalAddr, so
// do not modify it.
func (s *LocalSocket) LocalAddr() net.Addr {
	return s.conn.LocalAddr()
}

// RemoteAddr returns the remote network address.
// The Addr returned is shared by all invocations of RemoteAddr, so
// do not modify it.
func (s *LocalSocket) RemoteAddr() net.Addr {
	return s.conn.RemoteAddr()
}

// SetDeadline implements the Conn SetDeadline method.
func (s *LocalSocket) SetDeadline(t time.Time) error {
	return s.conn.SetDeadline(t)
}

// SetReadDeadline implements the Conn SetReadDeadline method.
func (s *LocalSocket) SetReadDeadline(t time.Time) error {
	return s.conn.SetReadDeadline(t)
}

// SetWriteDeadline implements the Conn SetWriteDeadline method.
func (s *LocalSocket) SetWriteDeadline(t time.Time) error {
	return s.conn.SetWriteDeadline(t)
}

// NewLocalSocket creates LocalSocket instance
func NewLocalSocket(pathName string) (*LocalSocket, error) {
	socketType := "unix" // SOCK_STREAM
	conn, err := net.DialUnix(socketType, nil, &net.UnixAddr{filepath.Join(os.TempDir(), pathName), socketType})
	if err != nil {
		return nil, err
	}
	return newLocalSocket(conn), nil
}

func newLocalSocket(conn *net.UnixConn) *LocalSocket {
	return &LocalSocket{
		conn: conn,
	}
}
