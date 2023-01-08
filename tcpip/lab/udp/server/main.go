// +build linux
package main

import (
	"flag"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"tcpip/netstack/tcpip"
	"tcpip/netstack/tcpip/link/fdbased"
	"tcpip/netstack/tcpip/link/tuntap"
	"tcpip/netstack/tcpip/network/arp"
	"tcpip/netstack/tcpip/network/ipv4"
	"tcpip/netstack/tcpip/network/ipv6"
	"tcpip/netstack/tcpip/stack"
	"tcpip/netstack/tcpip/transport/udp"
	"tcpip/netstack/waiter"
)

var mac = flag.String("mac", "aa:00:01:01:01:01", "mac address to use in tap device")

func main() {
	flag.Parse()
	if len(flag.Args()) != 4 {
		log.Fatal("Usage: ", os.Args[0], " <tap-device> <local-address/mask> <ip-address> <local-port>")
	}

	log.SetFlags(log.Lshortfile | log.LstdFlags)

	tapName := flag.Arg(0)
	cidrName := flag.Arg(1)
	addrName := flag.Arg(2)
	portName := flag.Arg(3)

	log.Printf("tap: %v, addr: %v, port: %v", tapName, addrName, portName)

	maddr, err := net.ParseMAC(*mac)
	if err != nil {
		log.Fatalf("Bad MAC address: %v", *mac)
	}

	parsedAddr := net.ParseIP(addrName)
	if err != nil {
		log.Fatalf("Bad addrress: %v", addrName)
	}

	// 解析地址ip地址，ipv4或者ipv6地址都支持
	var addr tcpip.Address
	var proto tcpip.NetworkProtocolNumber
	if parsedAddr.To4() != nil {
		addr = tcpip.Address(parsedAddr.To4())
		proto = ipv4.ProtocolNumber
	} else if parsedAddr.To16() != nil {
		addr = tcpip.Address(parsedAddr.To16())
		proto = ipv6.ProtocolNumber
	} else {
		log.Fatalf("Unknown IP type: %v", parsedAddr)
	}

	localPort, err := strconv.Atoi(portName)
	if err != nil {
		log.Fatalf("Unable to convert port %v: %v", portName, err)
	}

	// 虚拟网卡配置
	conf := &tuntap.Config{
		Name: tapName,
		Mode: tuntap.TAP,
	}

	var fd int
	// 新建虚拟网卡
	fd, err = tuntap.NewNetDev(conf)
	if err != nil {
		log.Fatal(err)
	}

	// 启动tap网卡
	_ = tuntap.SetLinkUp(tapName)
	// 设置路由
	_ = tuntap.SetRoute(tapName, cidrName)

	// 抽象的文件接口
	linkID := fdbased.New(&fdbased.Options{
		FD:                 fd,
		MTU:                1500,
		Address:            tcpip.LinkAddress(maddr),
		ResolutionRequired: true,
	})

	// 新建相关协议的协议栈
	s := stack.New([]string{ipv4.ProtocolName, ipv6.ProtocolName, arp.ProtocolName},
		[]string{udp.ProtocolName}, stack.Options{})

	// 新建抽象的网卡
	if err := s.CreateNamedNIC(1, "vnic1", linkID); err != nil {
		log.Fatal(err)
	}

	// 在该网卡上添加和注册相应的网络层
	if err := s.AddAddress(1, proto, addr); err != nil {
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

	var wq waiter.Queue
	// 新建一个UDP端
	ep, e := s.NewEndpoint(udp.ProtocolNumber, proto, &wq)
	if err != nil {
		log.Fatal(e)
	}

	// 绑定本地端口
	if err := ep.Bind(tcpip.FullAddress{1, addr, uint16(localPort)}, nil); err != nil {
		log.Fatal("Bind failed: ", err)
	}

	echo(&wq, ep)
}

func echo(wq *waiter.Queue, ep tcpip.Endpoint) {
	defer ep.Close()

	// Create wait queue entry that notifies a channel.
	waitEntry, notifyCh := waiter.NewChannelEntry(nil)

	wq.EventRegister(&waitEntry, waiter.EventIn)
	defer wq.EventUnregister(&waitEntry)

	var saddr tcpip.FullAddress

	for {
		v, _, err := ep.Read(&saddr)
		if err != nil {
			if err == tcpip.ErrWouldBlock {
				<-notifyCh
				continue
			}

			return
		}

		log.Printf("read and write data: %s", string(v))
		_, _, err = ep.Write(tcpip.SlicePayload(v), tcpip.WriteOptions{To: &saddr})
		if err != nil {
			log.Fatal(err)
		}
	}
}
