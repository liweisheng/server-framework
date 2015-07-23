package main

import (
	"fmt"
	"strings"
)

type S struct {
	name string
}

func (s *S) PrintName() {
	fmt.Println(s.name)
}

func Call(v interface{}, callee interface{}) {
	cb, ok := v.(func())
	if ok {
		cb()
	} else {
		fmt.Println("Fail")
	}

}

func main() {
	s := "connector:connector-1:127.0.0.1:8888::false"
	sss := (string)(nil)
	ss := strings.Split(s, ":")
	for _, v := range ss {
		fmt.Println(v)
	}
}
