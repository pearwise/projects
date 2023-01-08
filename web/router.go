package web

import (
	"errors"
	"log"
	"strings"
	//"time"
	"web/utils"
)

type Router struct {
	part       string
	Funcations HandleFuncs
	childs     []*Router
}

var root *Router

func (r *Router) method(path string, fs ...func(c *Context)) {
	count := len(path)
	if strings.ContainsRune(path, ':') {
		panic("the path must not have ':'")
	}
	count = len(r.childs)
	for i := 0; i < count; i++ {
		if r.childs[i].part == path {
			panic("repeat registration")
		}
	}
	router := &Router{
		part:       path,
		Funcations: fs,
		childs:     make([]*Router, 0),
	}
	r.childs = append(r.childs, router)
}

func (r *Router) Get(path string, f ...func(c *Context)) {
	r.method(utils.MethodToChar("GET")+path, f...)
}

func (r *Router) Post(path string, f ...func(c *Context)) {
	r.method(utils.MethodToChar("POST")+path, f...)
}

func (r *Router) Put(path string, f ...func(c *Context)) {
	r.method(utils.MethodToChar("PUT")+path, f...)
}

func (r *Router) Delete(path string, f ...func(c *Context)) {
	r.method(utils.MethodToChar("DELETE")+path, f...)
}

func (r *Router) Head(path string, f ...func(c *Context)) {
	r.method(utils.MethodToChar("HEAD")+path, f...)
}

func (r *Router) Options(path string, f ...func(c *Context)) {
	r.method(utils.MethodToChar("OPTIONS")+path, f...)
}

func (r *Router) Patch(path string, f ...func(c *Context)) {
	r.method(utils.MethodToChar("PATCH")+path, f...)
}

func (r *Router) Trace(path string, f ...func(c *Context)) {
	r.method(utils.MethodToChar("TRACE")+path, f...)
}

func (r *Router) Connect(path string, f ...func(c *Context)) {
	r.method(utils.MethodToChar("CONNECT")+path, f...)
}

func (r *Router) Any(path string, f ...func(c *Context)) {
	r.method(utils.MethodToChar("ANY")+path, f...)
}

func (r *Router) New(groupName string, fs ...func(c *Context)) *Router {
	if len(groupName) > 0 {
		if groupName[0] <= '9' || groupName[0] >= '0' {
			panic("the group name must not begin with a number")
		}
	}
	count := len(r.childs)
	for i := 0; i < count; i++ {
		if r.childs[i].part == groupName {
			panic("repeat registration")
		}
	}
	router := &Router{
		part:       groupName,
		Funcations: fs,
		childs:     make([]*Router, 0),
	}
	r.childs = append(r.childs, router)
	return router
}

func Do(ctx *Context) error {
	parts := strings.Split(ctx.Req.URL, "/")[1:]
	
	// process static file
	if parts[0] == staticPart {
		count := len(parts) - 1
		data, err := static(ctx, utils.GetSubfix(parts[count]))
		if err!=nil {
			ctx.conn.Write(BuildRespErr.ToBytes())
			return err
		}
		return ctx.conn.Write(data)
	}

	// process router
	temp := root
	length := len(parts) - 1
	if temp.Funcations != nil {
		ctx.Funcations = make([]HandleFuncs, 1, length+2)
		ctx.Funcations[0] = temp.Funcations
	} else {
		ctx.Funcations = make([]HandleFuncs, 0, length+1)
	}
	for i := 0; i < length; i++ {
		count := len(temp.childs) - 1
		for count > 0 {
			if temp.childs[count].part == parts[i] {
				// match succ
				if temp.childs[count].Funcations != nil {
					ctx.Funcations = append(ctx.Funcations, temp.childs[count].Funcations)
				}
				temp = temp.childs[count]
				goto NEXT
			}
			count--
		}
		if temp.childs[0].part[0] == ':' {
			ctx.Param = append(ctx.Param, [2]string{temp.childs[0].part[1:], parts[i]})
			// dynamic router, match succ
			if temp.childs[0].Funcations != nil {
				ctx.Funcations = append(ctx.Funcations, temp.childs[0].Funcations)
			}
			temp = temp.childs[0]
		} else if temp.childs[0].part == parts[i] {
			// match succ
			if temp.childs[0].Funcations != nil {
				ctx.Funcations = append(ctx.Funcations, temp.childs[0].Funcations)
			}
			temp = temp.childs[0]
		} else {
			// match fail
			return errors.New("match failed, request url:" + ctx.Req.URL)
		}
	NEXT:
	}
	for i := len(temp.childs)-1; i > -1; i-- {
		if utils.CheckMethodPath(temp.childs[i].part, parts[length], ctx.Req.Method) {
			// match succ
			if temp.childs[i].Funcations != nil {
				ctx.Funcations = append(ctx.Funcations, temp.childs[i].Funcations)
			}
			goto NEXT2
		}
	}
	// match fail
	log.Println("match failed, left =", temp.childs[0].part, ",right =", parts[length], "method =",ctx.Req.Method)
	return errors.New("match failed, request url: " + ctx.Req.URL)
	NEXT2:
	ctx.Resp = &Response{
		Header: new(strings.Builder),
	}
	ctx.Funcations[0][0](ctx)
	err := ctx.conn.rw.Writer.Flush()
	if err != nil {
		log.Println("=============")
		return err
	}
	_, err = ctx.conn.rw.Writer.Write(ctx.Resp.ToBytes())
	if err != nil {
		return err
	}
	err = ctx.conn.rw.Writer.Flush()
	if err != nil {
		log.Println("=============")
		return err
	}
	return nil
}