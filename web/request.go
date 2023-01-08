package web

import (
	"bufio"
	"log"
	"strings"
	"web/utils"
)

type Request struct {
	Method  string
	URL     string
	Version string
	Header  map[string][]string
	Body    []byte
}

func BuildRequest(reader *bufio.Reader) (*Request, error) {
	line, err := utils.ReadLine(reader)
	
	if err != nil {
		return nil, err
	}
	req := &Request{
		Header: make(map[string][]string),
	}
	lineParts := strings.Split(string(line), " ")
	// build request line
	req.Method, req.URL, req.Version = lineParts[0], lineParts[1], lineParts[2]
	for {
		// build header
		line, err = utils.ReadLine(reader)
		if err != nil {
			return nil, err
		}
		if len(line) == 0 {
			break
		}
		kv := strings.SplitN(string(line), ": ", 2)
		req.Header[kv[0]] = strings.Split(kv[1], ",")
	}
	// build body
	if _, ok := req.Header["Content-Type"]; !ok {
		return req, nil
	}
	line, err = utils.ReadLine(reader)
	if err != nil {
		return nil, err
	}
	req.Body = line
	log.Println(string(line))
	return req, nil
}
