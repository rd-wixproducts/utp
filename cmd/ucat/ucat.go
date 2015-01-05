package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"bitbucket.org/anacrolix/go.torrent/util"

	"bitbucket.org/anacrolix/utp"
)

func main() {
	util.LoggedHTTPServe("")
	listen := flag.Bool("l", false, "listen")
	port := flag.Int("p", 0, "port to listen on")
	flag.Parse()
	var (
		conn net.Conn
		err  error
	)
	if *listen {
		s, err := utp.NewSocket(fmt.Sprintf(":%d", *port))
		if err != nil {
			log.Fatal(err)
		}
		defer s.Close()
		conn, err = s.Accept()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		conn, err = utp.Dial(net.JoinHostPort(flag.Arg(0), flag.Arg(1)))
		if err != nil {
			log.Fatal(err)
		}
	}
	defer conn.Close()
	writerDone := make(chan struct{})
	go func() {
		defer close(writerDone)
		written, err := io.Copy(conn, os.Stdin)
		if err != nil {
			log.Fatalf("error after writing %d bytes: %s", err)
		}
		log.Printf("wrote %d bytes", written)
		conn.Close()
	}()
	n, err := io.Copy(os.Stdout, conn)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("received %d bytes", n)
	// <-writerDone
}