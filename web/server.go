package web

import (
	"log"
	"net"
	"time"
)

var connectTimeout time.Duration
var rwTimeout time.Duration
var tcpAaddr *net.TCPAddr
var tempWorker *TempWorker

func New(addr string, connectTimeOut time.Duration, rwTimeOut time.Duration, workerNum uint, maxTempWorkerNum uint, maxIdleTime time.Duration, fs ...func(c *Context)) *Router {
	tcpaddr, err := net.ResolveTCPAddr("tcp", addr)
	if err!=nil {
		panic(err)
	}
	tcpAaddr = tcpaddr
	if maxTempWorkerNum > 0 {
		tempWorker = &TempWorker{
			maxTempWorkerNum: maxTempWorkerNum,
            maxIdleTime:      maxIdleTime,
		}
	}
	if workerNum < 1 {
		panic("workerNum must be greater than zero")
	} else {
		for i := uint(0); i < workerNum; i++ {
			go Worker()
		}
	}
	connectTimeout = connectTimeOut
	rwTimeout = rwTimeOut
	root = &Router{}
	root.Funcations = fs
	ConnChan = make(chan *Connection, workerNum)
	return root
}

func Default(addr string, fs ...func(c *Context)) *Router {
	return New(addr, 5*time.Second, 5*time.Second, 5, 3, 5*time.Second, fs...)
}

func Run() error {
	lis, err := net.ListenTCP("tcp", tcpAaddr)
	if err!=nil {
		return err
	}

	if tempWorker == nil {
		// static
		for {
			conn, err := lis.AcceptTCP()
			log.Println("connection established, remote addr:", conn.RemoteAddr().String())
			if err!=nil {
				return err
			}
			ConnChan <- NewConn(conn)
		}
	} else {
		// dynamic
		var curTempNum uint
		var reqConn *Connection
		for {
			conn, err := lis.AcceptTCP()
			log.Println("connection established, remote addr:", conn.RemoteAddr().String())
			if err!=nil {
				return err
			}
			reqConn = NewConn(conn)
			select {
				case ConnChan <- reqConn:
				case <- time.After(connectTimeout):
					if tempWorker.maxTempWorkerNum > curTempNum {
						go tempWorker.Work(reqConn)
						curTempNum++
					}
			}
		}
	}
}