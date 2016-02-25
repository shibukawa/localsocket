package main

import (
	"bufio"
	"fmt"
	"github.com/shibukawa/localsocket"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

func main() {
	server := localsocket.NewLocalServer("testlocalsocket")
	server.SetOnConnectionCallback(func(socket net.Conn) {
		defer socket.Close()
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			writer := bufio.NewWriter(socket)
			count := rand.Intn(20)
			fmt.Printf("send %d items\n", count)
			for i := 0; i < count; i++ {
				fmt.Printf("write %d\n", i)
				fmt.Fprintf(writer, "%d\n", i)
			}
			writer.WriteString("end\n")
			writer.Flush()
			fmt.Println("done write")
			wg.Done()
		}()
		go func() {
			reader := bufio.NewReader(socket)
			counter := 1
			for {
				content, err := reader.ReadString('\n')
				if err != nil {
					fmt.Println(err)
					break
				}
				fmt.Printf("read %d: '%s'\n", counter, strings.TrimRight(content, "\n"))
				if content == "end\n" {
					break
				}
				counter++
			}
			fmt.Println("done read")
			wg.Done()
		}()
		wg.Wait()
	})

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGINT)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		server.ListenAndServe()
	}()

	fmt.Printf("Serving Server at %s\n", server.Path())
	select {
	case signal := <-stop:
		fmt.Printf("Got signal: %v\n", signal)
	}
	fmt.Printf("Stopping listener\n")
	server.Close()
	fmt.Printf("Waiting on server\n")
	wg.Wait()
}
