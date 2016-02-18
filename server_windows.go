package localsocket

import (
	"errors"
	"fmt"
	"gopkg.in/natefinch/npipe.v2"
	"sync"
)

var StoppedError = errors.New("Listener stopped")

type LocalServer struct {
	path     string
	callback func(socket *LocalSocket)
	listener *npipe.PipeListener
	wg       sync.WaitGroup
}

func (s *LocalServer) SetOnConnectionCallback(callback func(socket *LocalSocket)) {
	s.callback = callback
}

func (s *LocalServer) Path() string {
	return s.path
}

func NewLocalServer(pathName string) *LocalServer {
	return &LocalServer{
		path: `\\.\pipe\` + pathName,
	}
}

func (s *LocalServer) ListenAndServe() error {
	if s.callback == nil {
		return errors.New("no callback specified")
	}
	s.wg.Wait()
	var err error
	s.listener, err = npipe.Listen(s.path)
	if err != nil {
		return err
	}
	defer s.listener.Close()
	s.wg.Add(1)
	defer func() {
		s.wg.Done()
		s.listener.Close()
	}()
	for {
		conn, err := s.listener.AcceptPipe()
		if conn == nil {
			return StoppedError
		}
		if err == nil {
			s.wg.Add(1)
			go func() {
				s.callback(newLocalSocket(conn))
				s.wg.Done()
			}()
			continue
		}
		println(err.Error())
		s.wg.Done()
		return err
	}
}

func (s *LocalServer) Listen() error {
	s.wg.Wait()
	if s.callback == nil {
		return errors.New("no callback specified")
	}
	var err error
	s.listener, err = npipe.Listen(s.path)
	if err != nil {
		return err
	}
	s.wg.Add(1)
	go func() {
		defer func() {
			s.wg.Done()
			err := s.listener.Close()
			fmt.Println(err)
		}()
		for {
			conn, err := s.listener.AcceptPipe()
			// check channel is still opening
			if conn == nil {
				return
			}
			if err == nil {
				s.wg.Add(1)
				go func() {
					s.callback(newLocalSocket(conn))
					s.wg.Done()
				}()
				continue
			}
			println("error")
			panic(err)
		}
	}()
	return nil
}

func (s *LocalServer) Close() {
	s.listener.Close()
}
