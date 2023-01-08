package web

import (
	"strconv"
	"strings"
	"web/utils"
)

type Response struct {
	Status int
	Header *strings.Builder
	Body   []byte
}

func NewResponse(status int, header *strings.Builder, body []byte) *Response {
	return &Response{
        Status: status,
        Header: header,
        Body:   body,
    }
}

func (r *Response) WriteHeader(key string, values ...string) {
	r.Header.WriteString(key)
	r.Header.WriteString(": ")
	count := len(values)
	if count != 0 {
		r.Header.WriteString(values[0])
		count--
		for count > 0 {
			r.Header.WriteByte(',')
			r.Header.WriteString(values[count])
			count--
		}
	}
	r.Header.WriteString("\r\n")
}

var BuildRespErr *Response

var ProcessReqErr []byte

func (r *Response) ToBytes() []byte {
	resp := new(strings.Builder)
	resp.WriteString("HTTP/1.1 ")
	resp.WriteString(strconv.Itoa(r.Status))
	resp.WriteString(utils.StatusText(r.Status))
	resp.WriteString("\r\n")
	resp.WriteString(r.Header.String())
	resp.WriteString("\r\n")
	resp.Write(r.Body)
	resp.WriteString("\r\n")
	return []byte(resp.String())
}
