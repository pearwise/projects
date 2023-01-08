// +build linux

package main

import (
	"flag"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"tcpip/netstack/tcpip"
	"tcpip/netstack/tcpip/link/loopback"
	"tcpip/netstack/tcpip/network/arp"
	"tcpip/netstack/tcpip/network/ipv4"
	"tcpip/netstack/tcpip/stack"
	"tcpip/netstack/tcpip/transport/tcp"
	"tcpip/netstack/waiter"
)

var mac = flag.String("mac", "aa:00:01:01:01:01", "mac address to use in tap device")

// cd state; go build
// ./state tap0 192.168.1.1 9000
func main() {
	flag.Parse()
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	if len(os.Args) != 3 {
		log.Fatal("Usage: ", os.Args[0], "<ipv4-address> <port>")
	}

	addrName := os.Args[1]
	portName := os.Args[2]

	addr := tcpip.Address(net.ParseIP(addrName).To4())
	port, err := strconv.Atoi(portName)
	if err != nil {
		log.Fatalf("Unable to convert port %v: %v", portName, err)
	}

	s := newStack(addr, port)

	done := make(chan int, 1)
	go tcpServer(s, addr, port, done)
	<-done

	tcpClient(s, addr, port)
}

func newStack(addr tcpip.Address, port int) *stack.Stack {
	// 创建本地环回网卡
	linkID := loopback.New()

	// 新建相关协议的协议栈
	s := stack.New([]string{ipv4.ProtocolName, arp.ProtocolName},
		[]string{tcp.ProtocolName}, stack.Options{})

	// 新建抽象的网卡
	if err := s.CreateNamedNIC(1, "lo0", linkID); err != nil {
		log.Fatal(err)
	}

	// 在该网卡上添加和注册相应的网络层
	if err := s.AddAddress(1, ipv4.ProtocolNumber, addr); err != nil {
		log.Fatal(err)
	}

	// 在该协议栈上添加和注册ARP协议
	if err := s.AddAddress(1, arp.ProtocolNumber, arp.ProtocolAddress); err != nil {
		log.Fatal(err)
	}

	// 添加默认路由
	s.SetRouteTable([]tcpip.Route{
		{
			Destination: tcpip.Address(strings.Repeat("\x00", len(addr))),
			Mask:        tcpip.AddressMask(strings.Repeat("\x00", len(addr))),
			Gateway:     "",
			NIC:         1,
		},
	})

	return s
}

func tcpServer(s *stack.Stack, addr tcpip.Address, port int, done chan int) {
	var wq waiter.Queue
	// 新建一个TCP端
	ep, e := s.NewEndpoint(tcp.ProtocolNumber, ipv4.ProtocolNumber, &wq)
	if e != nil {
		log.Fatal(e)
	}

	// 绑定本地端口
	if err := ep.Bind(tcpip.FullAddress{0, "", uint16(port)}, nil); err != nil {
		log.Fatal("Bind failed: ", err)
	}

	// 监听tcp
	if err := ep.Listen(10); err != nil {
		log.Fatal("Listen failed: ", err)
	}

	// Wait for connections to appear.
	waitEntry, notifyCh := waiter.NewChannelEntry(nil)
	wq.EventRegister(&waitEntry, waiter.EventIn)
	defer wq.EventUnregister(&waitEntry)

	done <- 1

	for {
		n, _, err := ep.Accept()
		if err != nil {
			if err == tcpip.ErrWouldBlock {
				<-notifyCh
				continue
			}

			log.Fatal("Accept() failed:", err)
		}
		ra, err := n.GetRemoteAddress()
		log.Printf("new conn: %v %v", ra, err)
	}
}

func tcpClient(s *stack.Stack, addr tcpip.Address, port int) {
	remote := tcpip.FullAddress{
		Addr: addr,
		Port: uint16(port),
	}

	var wq waiter.Queue
	// 新建一个TCP端
	ep, e := s.NewEndpoint(tcp.ProtocolNumber, ipv4.ProtocolNumber, &wq)
	if e != nil {
		log.Fatal(e)
	}

	// Issue connect request and wait for it to complete.
	waitEntry, notifyCh := waiter.NewChannelEntry(nil)
	wq.EventRegister(&waitEntry, waiter.EventOut)
	terr := ep.Connect(remote)
	if terr == tcpip.ErrConnectStarted {
		log.Println("Connect is pending...")
		<-notifyCh
		terr = ep.GetSockOpt(tcpip.ErrorOption{})
	}
	wq.EventUnregister(&waitEntry)

	if terr != nil {
		log.Fatal("Unable to connect: ", terr)
	}

	log.Println("Connected")
	time.Sleep(1 * time.Second)

	ep.Close()
	log.Println("tcp disconected")

	time.Sleep(3 * time.Second)
}
