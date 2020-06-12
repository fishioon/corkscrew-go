package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

const (
	httpVersion = "HTTP/1.1"
)

func main() {
	var (
		conn            net.Conn
		n              int
		err, err0, err1 error
	)
	if len(os.Args) < 3 {
		return
	}
	host := os.Args[1]
	dst := os.Args[2]
	auth := os.Getenv("CORKSCREW_AUTH")
	if len(os.Args) == 4 {
		auth = os.Args[3]
	}
	uri := fmt.Sprintf("CONNECT %s %s\r\n", dst, httpVersion)
	if auth != "" {
		uri += "Proxy-Authorization: Basic " + auth + "\r\n"
	}
	uri += "\r\n"

	if conn, err = net.Dial("tcp", host); err != nil {
		log.Fatalf("conn %s %s", host, err.Error())
	}
	if n, err = conn.Write([]byte(uri)); err != nil {
		log.Fatalf("write %s", err.Error())
	}
	data := make([]byte, 256)
	if n, err = conn.Read(data); err != nil {
		log.Fatalf("read %s", err.Error())
	}
	if code := string(data[len(httpVersion)+1 : len(httpVersion)+4]); code != "200" {
		log.Fatalf("connect response %s", data[:n])
	}

	go func() {
		for err1 == nil {
			_, err = io.Copy(os.Stdout, conn)
		}
	}()
	for err0 == nil {
		_, err = io.Copy(conn, os.Stdin)
	}
}
