package main

import (
	"bufio"
	"fmt"
	"github.com/shibukawa/localsocket"
	"math/rand"
	"strings"
	"sync"
)

func main() {
	socket, err := localsocket.NewLocalSocket("testlocalsocket")
	if err != nil {
		panic(err)
	}
	defer socket.Close()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		counter := 1
		reader := bufio.NewReader(socket)
		for {
			fmt.Println("reading")
			content, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("err:", err)
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
	go func() {
		count := rand.Intn(30)
		fmt.Printf("send %d items\n", count)
		writer := bufio.NewWriter(socket)
		for i := 0; i < count; i++ {
			fmt.Printf("write: %d\n", i)
			fmt.Fprintf(writer, "%d\n", i)
		}
		writer.WriteString("end\n")
		writer.Flush()
		fmt.Println("done write")
		wg.Done()
	}()
	wg.Wait()
}
