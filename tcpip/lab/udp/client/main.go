package main

import (
	"flag"
	"log"
	"net"
)

func main() {
	var (
		addr = flag.String("a", "192.168.1.1:9000", "udp dst address")
	)

	log.SetFlags(log.Lshortfile | log.LstdFlags)

	udpAddr, err := net.ResolveUDPAddr("udp", *addr)
	if err != nil {
		panic(err)
	}

	// 建立UDP连接（只是填息了目的IP和端口，并未真正的建立连接）
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		panic(err)
	}

	send := []byte("hello")
	recv := make([]byte, 10)

	_, _ = conn.Write(send)
	log.Printf("send: %s", string(send))
	rn, _, err := conn.ReadFrom(recv)
	if err != nil {
		panic(err)
	}
	log.Printf("recv: %s", string(recv[:rn]))
}
