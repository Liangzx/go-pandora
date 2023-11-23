package main

import (
	"log"
	"net"
	"os"
	"syscall"

	"golang.org/x/sys/unix"
)

const sock = "/tmp/sock.sock"

func main() {
	_ = syscall.Unlink(sock)
	l, _ := net.Listen("unix", sock)
	log.Print("listen " + sock)
	defer l.Close()
	recvListener(l.(*net.UnixListener))
}

func recvListener(unixListener *net.UnixListener) {
	unixConn, _ := unixListener.AcceptUnix()
	defer unixConn.Close()
	var (
		buf = make([]byte, 1)
		oob = make([]byte, 1024)
	)
	_, oobn, _, _, _ := unixConn.ReadMsgUnix(buf, oob)
	scms, _ := unix.ParseSocketControlMessage(oob[0:oobn])

	recvFds, _ := unix.ParseUnixRights(&scms[0])

	recvFile := os.NewFile(uintptr(recvFds[0]), "")

	recvFileListener, _ := net.FileListener(recvFile)

	recvListener := recvFileListener.(*net.TCPListener)

	for {
		conn, _ := recvListener.Accept()
		go func() {
			for {
				var buf = make([]byte, 1)
				conn.Read(buf)
				conn.Write([]byte("b process response\n"))
			}
		}()

	}
}
