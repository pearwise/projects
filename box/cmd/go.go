package main

import (
	"fmt"
	"box"
	"reflect"
)

type TwoString struct {
	s1 *string
	s2 *string
}

func main() {
	// type inject
	err := box.Put("java")
	if err!=nil {
		panic(err)
	}
	s, err := box.Get[string]()
	if err!=nil {
		panic(err)
	}
	// name inject
	err = box.PutByName("name","jack")
	if err!=nil {
		panic(err)
	}
	name, err := box.GetByName[string]("name")
	if err!=nil {
		panic(err)
	}
	fmt.Println(s, name)

	// shallow copy
	twoStr := TwoString{}
	twoStr.s1 = new(string)
	twoStr.s2 = new(string)
	*twoStr.s1 = "string1"
	*twoStr.s2 = "string2"
	err = box.Put(twoStr)
	if err!= nil {
        panic(err)
    }
	twoStr2, err := box.Get[TwoString]()
	if err!=nil {
		panic(err)
	}

	fmt.Println("twoStr==twoStr2: ", twoStr==twoStr2)
	// deep copy
	err = box.DeepPut(func() reflect.Value {
		twoString := TwoString{}
		twoString.s1 = new(string)
		twoString.s2 = new(string)
		*twoString.s1 = "string1"
		*twoString.s2 = "string2"
		return reflect.ValueOf(twoString)
	})
	if err!=nil {
		panic(err)
	}
	ts, err := box.DeepGet[TwoString]()
	if err!=nil {
		panic(err)
	}
	twoString, err := box.DeepGet[TwoString]()
	if err!=nil {
		panic(err)
	}
	fmt.Println("ts==twoString: ",ts==twoString)
	fmt.Println(*ts.s1, *ts.s2)
	fmt.Println(*twoString.s1, *twoString.s2)
	*ts.s1, *ts.s2 = "string2", "string1"
	fmt.Println(*ts.s1, *ts.s2)
	fmt.Println(*twoString.s1, *twoString.s2)
}
