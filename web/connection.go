package web

import (
	"bufio"
	"io"
	"log"
	"net"
)

var ConnChan chan *Connection

type Connection struct {
	conn   *net.TCPConn
	rw *bufio.ReadWriter
}

func NewConn(conn *net.TCPConn) *Connection {
	return &Connection{
		conn:   conn,
		rw: bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)),
	}
}

func (c *Connection) Close() error {
	return c.conn.Close()
}

func (c *Connection) ReadLine() ([]byte, error) {
	data, more, err := c.rw.ReadLine()
	if err != nil {
		return nil, err
	}
	if !more {
		return data, nil
	}

	var line []byte
	for more {
		line, more, err = c.rw.ReadLine()
		if err != nil {
			return nil, err
		}
		data = append(data, line...)
	}
	return data, nil
}

func (c *Connection) Read(length int) ([]byte, error) {
	buf := make([]byte, length)
	_, err := io.ReadFull(c.rw, buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (c *Connection) Write(data []byte) error {
	_, err := c.conn.Write(data)
	return err
}

func (c *Connection) WriteHeader(map[string][]string) {

}

func (c *Connection) ErrorResponse(resp []byte) (int, error) {
	return c.conn.Write(resp)
}

func (c *Connection) Do() {
	for {
		log.Println("process connection")
		NEXT:
		req, err := BuildRequest(c.rw.Reader)
		if err != nil {
			if err == io.EOF {
				c.rw.Reader.Reset(c.conn)
				goto NEXT
			} else {
				log.Println(err.Error())
				c.conn.Write(ProcessReqErr)
				return
			}
		}
		ctx := &Context{
			conn: c,
			Req:  req,
		}
		err = Do(ctx)
		if err != nil {
			log.Println(err.Error())
			c.conn.Write(ProcessReqErr)
			return
		}
	}
}
