package web

import "testing"


func TestMap(t *testing.T) {
	m := make(map[string]string)
	m["hello"] = "world"
	s := m["world"]
	if s=="" {
		println("error-----------")
	} else {
		println("ok",s)
	}
}