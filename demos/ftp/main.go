package main

import (
	"log"

	"goftp.io/server/v2"
	"goftp.io/server/v2/driver/file"
)

func main() {
	username := "foobar"
	password := "foobar"
	path := "E:/musicx"

	auth := &server.SimpleAuth{Name: username, Password: password}
	driver := &file.Driver{RootPath: path}

	opts := &server.Options{
		Name:     "foobar2000",
		Hostname: "192.168.1.7",
		Port:     10086,
		Auth:     auth,
		Driver:   driver,
		Perm:     server.NewSimplePerm("root", "root"),
	}

	ser, err := server.NewServer(opts)
	if err != nil {
		log.Println("server.NewServer(opts):", err)
	}

	err = ser.ListenAndServe()
	if err != nil {
		log.Println("ser.ListenAndServe:", err)
	}

}
