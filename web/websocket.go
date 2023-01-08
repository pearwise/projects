package web

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"strings"
	"unsafe"
)

const (
	// Frame header byte 0 bits from Section 5.2 of RFC 6455
	finalBit = 1 << 7
	rsv1Bit  = 1 << 6
	rsv2Bit  = 1 << 5
	rsv3Bit  = 1 << 4

	maskBit = 1 << 7
)

// The message types are defined in RFC 6455, section 11.8.
const (
	// TextMessage denotes a text data message. The text message payload is
	// interpreted as UTF-8 encoded text data.
	TextMessage byte = 0x1

	// BinaryMessage denotes a binary data message.
	BinaryMessage byte = 0x2

	// CloseMessage denotes a close control message. The optional message
	// payload contains a numeric code and text. Use the FormatCloseMessage
	// function to format a close message payload.
	CloseMessage byte = 0x8

	// PingMessage denotes a ping control message. The optional message payload
	// is UTF-8 encoded text.
	PingMessage byte = 0x9

	// PongMessage denotes a pong control message. The optional message payload
	// is UTF-8 encoded text.
	PongMessage byte = 0x10
)

const (
	noFrame = -1
)

type WebsocketConn struct {
	conn *net.TCPConn
	r    *bufio.Reader
}

func NewWebsocketConn(conn *net.TCPConn) *WebsocketConn {
	return &WebsocketConn{conn: conn, r: bufio.NewReader(conn)}
}

func (c *WebsocketConn) Close() error {
	return c.conn.Close()
}

func (c *WebsocketConn) Read() (byte, []byte, error) {
	var p = make([]byte, 2)
	c.conn.SetKeepAlive(true)
	err := binary.Read(c.conn, binary.LittleEndian, p[0])
	if err != nil {
		log.Println("读取错误", err)
		return 0, nil, err
	}
	err = binary.Read(c.conn, binary.LittleEndian, p[1])
	if err != nil {
		log.Println("读取错误", err)
		return 0, nil, err
	}
	log.Println("前2个byte值", p)
	//通过将1进行左移一定的位数后与第一个字节进行与运算，可以判断该值是否为1
	final := p[0]&finalBit == 0
	rsv1 := p[0]&rsv1Bit != 0
	rsv2 := p[0]&rsv2Bit != 0
	rsv3 := p[0]&rsv3Bit != 0
	//通过0xf=15来与后4位进行与运行，可以得到这4位原本的值
	opcode := p[0] & 0xf
	log.Println("第一个字节的值", final, rsv1, rsv2, rsv3, opcode)
	//获取第2个字节的数据
	masked := p[1]&maskBit != 0
	length := uint64(p[1] & 0x7f)
	if length == 126 {
		var lenBytes = make([]byte, 2)
		c.r.Read(lenBytes)
		length = binary.BigEndian.Uint64(lenBytes)
	}
	if length == 127 {
		var lenBytes = make([]byte, 8)
		c.r.Read(lenBytes)
		length = binary.BigEndian.Uint64(lenBytes)
	}
	log.Println("masked and length", masked, length)
	switch opcode {
	case CloseMessage, PingMessage, PongMessage:
		log.Println("控制消息", opcode)
		break
	case TextMessage:
		log.Println("文本消息", opcode)
		break
	case BinaryMessage:
		log.Println("二进制消息", opcode)
		break
	}
	//读取mask掩码
	maskkey := make([]byte, 4)
	if masked {
		c.r.Read(maskkey)
	}
	//读取内容
	payload := make([]byte, length)
	c.r.Read(payload)
	if masked {
		//解码
		maskBytes(maskkey, 0, payload)
	}
	log.Println("内容", string(payload))
	// fin不为1则一直循环下去
	for final {
		_, err := c.r.Read(p)
		if err != nil {
			log.Println("读取错误", err)
			return 0, nil, err
		}
		log.Println("前2个byte值", p)
		//通过将1进行左移一定的位数后与第一个字节进行与运算，可以判断该值是否为1
		final = p[0]&finalBit == 0
		rsv1 = p[0]&rsv1Bit != 0
		rsv2 = p[0]&rsv2Bit != 0
		rsv3 = p[0]&rsv3Bit != 0
		//通过0xf=15来与后4位进行与运行，可以得到这4位原本的值
		opcode := p[0] & 0xf
		log.Println("第一个字节的值", final, rsv1, rsv2, rsv3, opcode)
		//获取第2个字节的数据
		masked := p[1]&maskBit != 0
		length := uint64(p[1] & 0x7f)
		if length == 126 {
			var lenBytes = make([]byte, 2)
			c.r.Read(lenBytes)
			length = binary.BigEndian.Uint64(lenBytes)
		}
		if length == 127 {
			var lenBytes = make([]byte, 8)
			c.r.Read(lenBytes)
			length = binary.BigEndian.Uint64(lenBytes)
		}
		log.Println("masked and length", masked, length)
		switch opcode {
		case CloseMessage, PingMessage, PongMessage:
			log.Println("控制消息", opcode)
		case TextMessage:
			log.Println("文本消息", opcode)
		case BinaryMessage:
			log.Println("二进制消息", opcode)
		}
		//读取mask掩码
		maskkey := make([]byte, 4)
		if masked {
			c.r.Read(maskkey)
		}
		//读取内容
		payloadData := make([]byte, length)
		c.r.Read(payloadData)
		if masked {
			//解码
			maskBytes(maskkey, 0, payloadData)
		}
		payload = append(payload, payloadData...)
	}
	log.Println("内容", string(payload))
	return opcode, payload, nil
}

func (c *WebsocketConn) Write(opCode byte, payload []byte) error {
	defer func() {
		payload = nil
	}()
	if !CheckDataType(opCode) {
		return fmt.Errorf("operate code: %v is not supported", opCode)
	}
	if payload == nil {
		return fmt.Errorf("payload is nil")
	}
	length := len(payload)
	if len(payload) == 0 {
		return fmt.Errorf("payload is empty")
	}
	// data type is 4 bit
	if opCode > 15 {
		return fmt.Errorf("data type max value is 15")
	}
	switch {
	case length < 126:
		header := make([]byte, 2)
		header[0] = finalBit + opCode
		header[1] = byte(length)
		_, err := c.conn.Write(header)
		if err != nil {
			return err
		}
		_, err = c.conn.Write(payload)
		// this err is nil or error
		return err
	case length < 65536:
		header := make([]byte, 4)
		header[0] = finalBit + opCode
		header[1] = 126
		binary.BigEndian.PutUint16(header[2:], uint16(length))
		_, err := c.conn.Write(header)
		if err != nil {
			return err
		}
		_, err = c.conn.Write(payload)
		// this err is nil or error
		return err
	default:
		if length < 9223372036854775807 {
			header := make([]byte, 10)
			header[0] = finalBit + opCode
			header[1] = 127
			binary.BigEndian.PutUint64(header[2:], uint64(length))
			_, err := c.conn.Write(header)
			if err != nil {
				return err
			}
			_, err = c.conn.Write(payload)
			// this err is nil or error
			return err
		} else {
			header := make([]byte, 10)
			header[1] = 127
			if length2 := uint64(len(payload[9223372036854775807:])); length2 < 9223372036854775807 {
				header[0] = finalBit + opCode
				binary.BigEndian.PutUint64(header[2:], uint64(length)+length2)
				_, err := c.conn.Write(header)
				if err != nil {
					return err
				}
				_, err = c.conn.Write(payload[:length])
				if err != nil {
					return err
				}
				payload = payload[length:]
				_, err = c.conn.Write(payload)
				return err
			} else if length2 = uint64(len(payload[9223372036854775807:][9223372036854775807:])); length2 < 3 {
				header[0] = finalBit + opCode
				binary.BigEndian.PutUint64(header[2:], uint64(length)<<1+length2)
				_, err := c.conn.Write(header)
				if err != nil {
					return err
				}
				_, err = c.conn.Write(payload[:length])
				if err != nil {
					return err
				}
				payload = payload[length:]
				_, err = c.conn.Write(payload[:length])
				if err != nil {
					return err
				}
				payload = payload[length:]
				_, err = c.conn.Write(payload)
				// this err is nil or error
				return err
			} else {
				// payload is too large, split it
				header[0] = opCode
				binary.BigEndian.PutUint64(header[2:], 18446744073709551615)
				_, err := c.conn.Write(header)
				if err != nil {
					return err
				}
				_, err = c.conn.Write(payload[:length])
				if err != nil {
					return err
				}
				payload = payload[length:]
				_, err = c.conn.Write(payload[:length])
				if err != nil {
					return err
				}
				_, err = c.conn.Write(payload[:2])
				if err != nil {
					return err
				} else {
					payload = payload[2:]
					err = c.Write(opCode, payload)
					// this err is nil or error
					return err
				}
			}
		}
	}
}

func (c *WebsocketConn) WriteWithMask(opCode byte, payload []byte, mask uint32) error {
	defer func() {
		payload = nil
	}()
	if !CheckDataType(opCode) {
		return fmt.Errorf("operate code: %v is not supported", opCode)
	}
	if payload == nil {
		return fmt.Errorf("payload is nil")
	}
	length := len(payload)
	if len(payload) == 0 {
		return fmt.Errorf("payload is empty")
	}
	// data type is 4 bit
	if opCode > 15 {
		return fmt.Errorf("data type max value is 15")
	}
	switch {
	case length < 126:
		header := make([]byte, 6)
		header[0] = finalBit + opCode
		header[1] = maskBit + byte(length)
		// write mask
		binary.BigEndian.PutUint32(header[2:], mask)
		_, err := c.conn.Write(header)
		if err != nil {
			return err
		}
		// write bady
		_, err = c.conn.Write(payload)
		// this err is nil or error
		return err
	case length < 65536:
		header := make([]byte, 8)
		header[0] = finalBit + opCode
		header[1] = maskBit + 126
		binary.BigEndian.PutUint16(header[2:4], uint16(length))
		binary.BigEndian.PutUint32(header[4:], mask)
		_, err := c.conn.Write(header)
		if err != nil {
			return err
		}
		_, err = c.conn.Write(payload)
		// this err is nil or error
		return err
	default:
		if length < 9223372036854775807 {
			header := make([]byte, 14)
			header[0] = finalBit + opCode
			header[1] = 127
			binary.BigEndian.PutUint64(header[2:10], uint64(length))
			binary.BigEndian.PutUint32(header[10:], mask)
			_, err := c.conn.Write(header)
			if err != nil {
				return err
			}
			_, err = c.conn.Write(payload)
			// this err is nil or error
			return err
		} else {
			header := make([]byte, 14)
			header[1] = 127
			binary.BigEndian.PutUint32(header[10:], mask)
			if length2 := uint64(len(payload[9223372036854775807:])); length2 < 9223372036854775807 {
				header[0] = finalBit + opCode
				binary.BigEndian.PutUint64(header[2:10], uint64(length)+length2)
				_, err := c.conn.Write(header)
				if err != nil {
					return err
				}
				_, err = c.conn.Write(payload[:length])
				if err != nil {
					return err
				}
				payload = payload[length:]
				_, err = c.conn.Write(payload)
				return err
			} else if length2 = uint64(len(payload[9223372036854775807:][9223372036854775807:])); length2 < 3 {
				header[0] = finalBit + opCode
				binary.BigEndian.PutUint64(header[2:10], uint64(length)<<1+length2)
				_, err := c.conn.Write(header)
				if err != nil {
					return err
				}
				_, err = c.conn.Write(payload[:length])
				if err != nil {
					return err
				}
				payload = payload[length:]
				_, err = c.conn.Write(payload[:length])
				if err != nil {
					return err
				}
				payload = payload[length:]
				_, err = c.conn.Write(payload)
				// this err is nil or error
				return err
			} else {
				// payload is too large, split it
				header[0] = opCode
				binary.BigEndian.PutUint64(header[2:10], 18446744073709551615)
				_, err := c.conn.Write(header)
				if err != nil {
					return err
				}
				_, err = c.conn.Write(payload[:length])
				if err != nil {
					return err
				}
				payload = payload[length:]
				_, err = c.conn.Write(payload[:length])
				if err != nil {
					return err
				}
				_, err = c.conn.Write(payload[:2])
				if err != nil {
					return err
				} else {
					payload = payload[2:]
					err = c.Write(opCode, payload)
					// this err is nil or error
					return err
				}
			}
		}
	}
}

func CheckDataType(dataType byte) bool {
	switch dataType {
	case CloseMessage, PingMessage, PongMessage:
		log.Println("控制消息", dataType)
	case TextMessage:
		log.Println("文本消息", dataType)
	case BinaryMessage:
		log.Println("二进制消息", dataType)
	default:
		return false
	}
	return true
}

const wordSize = int(unsafe.Sizeof(uintptr(0)))

func maskBytes(key []byte, pos int, b []byte) int {
	// Mask one byte at a time for small buffers.
	if len(b) < 2*wordSize {
		for i := range b {
			b[i] ^= key[pos&3]
			pos++
		}
		return pos & 3
	}

	// Mask one byte at a time to word boundary.
	if n := int(uintptr(unsafe.Pointer(&b[0]))) % wordSize; n != 0 {
		n = wordSize - n
		for i := range b[:n] {
			b[i] ^= key[pos&3]
			pos++
		}
		b = b[n:]
	}

	// Create aligned word size key.
	var k [wordSize]byte
	for i := range k {
		k[i] = key[(pos+i)&3]
	}
	kw := *(*uintptr)(unsafe.Pointer(&k))

	// Mask one word at a time.
	n := (len(b) / wordSize) * wordSize
	for i := 0; i < n; i += wordSize {
		*(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(&b[0])) + uintptr(i))) ^= kw
	}

	// Mask one byte at a time for remaining bytes.
	b = b[n:]
	for i := range b {
		b[i] ^= key[pos&3]
		pos++
	}

	return pos & 3
}

func SecWebSocketAccept(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// handshake
func Upgrade(c *Context) (*WebsocketConn, error) {
	defer c.conn.conn.CloseWrite()
	
	if ss, ok := c.Req.Header["Connection"]; !ok {
		return nil, fmt.Errorf("missing Connection header")
	} else if ss[0] != "Upgrade" {
		return nil, fmt.Errorf("request Connection header error")
	}
	if ss, ok := c.Req.Header["Upgrade"]; !ok {
		return nil, fmt.Errorf("missing Upgrade header")
	} else if ss[0] != "websocket" {
		return nil, fmt.Errorf("request Upgrade header error")
	}
	if ss, ok := c.Req.Header["Sec-WebSocket-Version"]; !ok {
		return nil, fmt.Errorf("missing Sec-WebSocket-Version header")
	} else if ss[0] != "13" {
		return nil, fmt.Errorf("request websocket version not 13")
	}
	ss, ok := c.Req.Header["Sec-WebSocket-Key"]
	if !ok {
		return nil, fmt.Errorf("missing Sec-WebSocket-Key header")
	}
	res := NewResponse(101, new(strings.Builder), nil)
	res.WriteHeader("Sec-WebSocket-Accept", SecWebSocketAccept(ss[0]))
	res.WriteHeader("Upgrade", "websocket")
	res.WriteHeader("Connection", "upgrade")
	err := c.conn.Write(c.Resp.ToBytes())
	if err != nil {
		return nil, err
	}
	return NewWebsocketConn(c.conn.conn), nil
}