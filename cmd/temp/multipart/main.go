package main

import (
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
)

const boundary = "MachliJalKiRaniHaiJeevanUskaPaaniHai"

func handle(w http.ResponseWriter, req *http.Request) {
	partReader := multipart.NewReader(req.Body, boundary)
	buf := make([]byte, 256)
	for {
		part, err := partReader.NextPart()
		if err == io.EOF {
			break
		}
		var n int
		for {
			n, err = part.Read(buf)
			if err == io.EOF {
				break
			}
			fmt.Printf(string(buf[:n]))
		}
		fmt.Printf(string(buf[:n]))
	}
}

func main() {
	/* Net listener */
	n := "tcp"
	addr := "127.0.0.1:9094"
	l, err := net.Listen(n, addr)
	if err != nil {
		panic("AAAAH")
	}

	/* HTTP server */
	server := http.Server{
		Handler: http.HandlerFunc(handle),
	}
	server.Serve(l)
}
