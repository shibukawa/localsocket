// +build !windows

package localsocket

import (
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var StoppedError = errors.New("Listener stopped")

type LocalServer struct {
	path     string
	callback func(socket net.Conn)
	stop     chan int
	wg       sync.WaitGroup
}

func (s *LocalServer) SetOnConnectionCallback(callback func(net.Conn)) {
	s.callback = callback
}

func (s *LocalServer) Path() string {
	return s.path
}

func NewLocalServer(pathName string) *LocalServer {
	return &LocalServer{
		path: filepath.Join(os.TempDir(), pathName),
	}
}

func (s *LocalServer) ListenAndServe() error {
	if s.callback == nil {
		return errors.New("no callback specified")
	}
	listener, err := net.ListenUnix("unix", &net.UnixAddr{s.path, "unix"})
	if err != nil {
		return err
	}
	defer listener.Close()
	s.stop = make(chan int)
	s.wg.Add(1)
	for {
		listener.SetDeadline(time.Now().Add(time.Millisecond * time.Duration(1000)))
		conn, err := listener.AcceptUnix()
		// check channel is still opening
		select {
		case <-s.stop:
			s.wg.Done()
			return StoppedError
		default:
			// opening
		}
		if err == nil {
			s.wg.Add(1)
			go func() {
				s.callback(newLocalSocket(conn))
				s.wg.Done()
			}()
			continue
		}
		netErr, ok := err.(net.Error)
		if ok && netErr.Timeout() && netErr.Temporary() {
			continue
		}
		println(err.Error())
		s.wg.Done()
		return err
	}
}

func (s *LocalServer) Listen() error {
	if s.callback == nil {
		return errors.New("no callback specified")
	}
	listener, err := net.ListenUnix("unix", &net.UnixAddr{s.path, "unix"})
	if err != nil {
		return err
	}
	s.stop = make(chan int)
	s.wg.Add(1)
	go func() {
		defer func() {
			s.wg.Done()
			err := listener.Close()
			fmt.Println(err)
		}()
		for {
			listener.SetDeadline(time.Now().Add(time.Millisecond * time.Duration(333)))
			conn, err := listener.AcceptUnix()
			// check channel is still opening
			select {
			case <-s.stop:
				return
			default:
				// opening
			}
			if err == nil {
				s.wg.Add(1)
				go func() {
					s.callback(newLocalSocket(conn))
					s.wg.Done()
				}()
				continue
			}
			netErr, ok := err.(net.Error)
			if ok && netErr.Timeout() && netErr.Temporary() {
				continue
			}
			println("error")
			panic(err)
		}
	}()
	return nil
}

func (s *LocalServer) Close() {
	if s.stop != nil {
		close(s.stop)
		s.wg.Wait()
		s.stop = nil
	}
}
