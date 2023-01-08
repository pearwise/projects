package web

import (
	"embed"
	"log"
	"strings"
)

var staticFS embed.FS
var staticPart string
var ContentType = map[string]string{"html":"text/html","css":"text/css","js":"text/javascript","json":"application/json","xml":"application/xml"}

func LoadStaticFile(path string, files embed.FS) {
	staticPart = path
	staticFS = files
}

func LoadStatic(files embed.FS) {
		
}

func static(ctx *Context, fileSubfix string) (data []byte,err error) {
	defer func() {
		if r := recover();r != nil {
			log.Println(0)
			err = r.(error)
		}
	}()
	ctx.Resp = &Response{Header: new(strings.Builder)}
	ctx.Resp.Body, err = staticFS.ReadFile(ctx.Req.URL[1:])
	if err!=nil {
		log.Println(1)
		return nil, err
	}
	ctx.SetStatus(200)
	ctx.SetHeader("Content-Type", ContentType[fileSubfix])
	data = ctx.Resp.ToBytes()
	return
}