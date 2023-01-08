package was

import (
	"strconv"
	"syscall/js"
)

var (
	document = js.Global().Get("document")
	btn = js.Global().Get("btn")
)

func Document() js.Value {
	return document
}

func Btn() js.Value {
	return btn
}

func GetElementById(id string) js.Value {
	return js.Global().Call("getElementById", id)
}

type Object struct {
	Propertys map[string]any
	Get func(property string)any
	Set func(new any, old any)
}

func NewObj(get func(property string)any, set func(new any, old any)) *Object {
	return &Object{
		Propertys: make(map[string]any),
		Get: get,
		Set: set,
	}
}

func FibFunc(this js.Value, args []js.Value) interface{} {
	v := GetElementById("num").Get("value")
	if num, err := strconv.Atoi(v.String()); err == nil {
		GetElementById("ans").Set("innerHTML", js.ValueOf(fib(num)))
	}
	return nil
}

func fib(i int) int {
	if i == 0 || i == 1 {
		return 1
	}
	return fib(i-1) + fib(i-2)
}
