package main

import (
	"embed"
	"log"
	"strings"
	"time"
	"web"
)

//go:embed template
var Template embed.FS
//go:embed favicon.ico
var ico []byte

func main() {
	r := web.Default(":9999")
	web.BuildRespErr = web.NewResponse(500, new(strings.Builder),nil)
	r.Get("websocket",func(c *web.Context) {
		conn, err := web.Upgrade(c)
		if err != nil {
			panic(err)
		}
		go func()  {
			var opCode byte
			var readBuf []byte
			for {
				opCode, readBuf, err = conn.Read()
				if err != nil {
					panic(err)
				}
				log.Println(opCode)
				log.Println(string(readBuf))
			}
		}()	
		for {
			err = conn.Write(web.TextMessage, []byte("hello world"))
			if err!=nil {
				conn.Close()
			}
			println("websocket connect successfully")
			time.Sleep(time.Minute)
		}
	})
	r.Get("favicon.ico", func(c *web.Context) {
		c.Resp = &web.Response{Header: new(strings.Builder)}
		c.Resp.WriteHeader("Content-Type", "image/x-ico")
		c.Resp.Status = 200
		c.Resp.Body = ico
	})
	web.LoadStaticFile("template",Template)
	err := web.Run()
	if err!=nil {
		log.Println(err.Error())
	}
}