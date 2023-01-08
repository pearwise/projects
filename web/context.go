package web

import (
	"encoding/json"
	"strings"
)

type HandleFuncs []func(c *Context)

type Context struct {
	Req        *Request
	Resp       *Response
	conn       *Connection
	Funcations []HandleFuncs
	// 0 index is k, 1 index is v
	Param [][2]string
}

func (c *Context) GetParam(k string) (string, bool) {
	for i := len(c.Param) - 1; i > -1; i-- {
		if c.Param[i][0] == k {
			return c.Param[i][1], true
		}
	}
	return "", false
}

func (c *Context) Next() {
	c.Funcations[0] = c.Funcations[0][1:]
	if len(c.Funcations[0]) == 0 {
		c.Funcations = c.Funcations[1:]
	}
	c.Funcations[0][0](c)
}

func (c *Context) SetStatus(code int) {
	c.Resp.Status = code
}

func (c *Context) SetHeader(key string, values ...string) {
	c.Resp.WriteHeader(key, values...)
}

func (c *Context) SetBody(body []byte) {
	c.Resp.Body = body
}

func (c *Context) String(code int, s string) {
	c.Resp.Body = []byte(s)
	c.SetStatus(code)
	c.SetHeader("Content-Type", "text/plain")
}

func (c *Context) Data(code int, data []byte) {
	c.Resp.Body = data
	c.SetStatus(code)
	c.SetHeader("Content-Type", "text/plain")
}

func (c *Context) Json(code int, obj any) {
	var err error
	c.Resp.Body, err = json.Marshal(obj)
	if err != nil {
		c.Resp = BuildRespErr
	}
	c.SetStatus(code)
	c.SetHeader("Content-Type", "application/json")
}

func (c *Context) HTML(code int, filePath string) {
	c.SetStatus(code)
	parts := strings.Split(filePath, "/")
	count := len(parts) - 1
	filename := parts[count]
	code = len(filename)
	for code = len(filename) - 1; code > -1; code-- {
		if filename[code] == '.' {
			break
		}
	}
	if code == -1 {
		panic("not file")
	}
	c.SetHeader("Content-Type", "text/html")
}
