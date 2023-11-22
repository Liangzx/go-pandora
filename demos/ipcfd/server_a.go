package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const sock = "/tmp/sock.sock"

func main() {
	l, _ := net.Listen("tcp", ":9087")
	defer l.Close()
	// 监听信号
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	// 开启迁移链接的goroutine
	go transferConn(ch, l.(*net.TCPListener))
	for {
		c, err := l.Accept()
		if err != nil {
			log.Println(err)
			break
		}
		go func() {
			for {
				var buf = make([]byte, 1)
				c.Read(buf)
				c.Write([]byte("a process response\n"))
			}
		}()

	}
}

func transferConn(ch <-chan os.Signal, l *net.TCPListener) {
	<-ch
	log.Print("start transfer listener")
	c, err := net.Dial("unix", sock)
	if err != nil {
		log.Println(err)
	}

	unixConn := c.(*net.UnixConn)
	listenFile, err := l.File()
	rights := syscall.UnixRights(int(listenFile.Fd()))
	var buf = make([]byte, 1)

	unixConn.WriteMsgUnix(buf, rights, nil)
	fmt.Println("ending transfer listener")
	l.Close()
}
