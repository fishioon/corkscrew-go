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
		err  error
		conn net.Conn
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
	if _, err = conn.Write([]byte(uri)); err != nil {
		log.Fatalf("write %s", err.Error())
	}
	data := make([]byte, 256)
	if _, err = conn.Read(data); err != nil {
		log.Fatalf("read %s", err.Error())
	}
	if code := string(data[len(httpVersion)+1 : len(httpVersion)+4]); code != "200" {
		log.Fatalf("connect response code %s not equal 200", code)
	}

	go func() {
		for {
			if _, err := io.Copy(os.Stdout, conn); err != nil {
			}
		}
	}()
	for {
		if _, err := io.Copy(conn, os.Stdin); err != nil {
		}
	}
}
