package emutools

import (
	"fmt"
	"reflect"
	"runtime"
)

func Hex(variable interface{}, n int) string {
	h := fmt.Sprintf("%x", variable)
	if len(h) > n {
		h = h[:n]
	}
	return h
}

func GetFunName(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func GetFunNameAddr(f interface{}) string {
	n := GetFunName(f)
	n = n[len(n)-6 : len(n)-3]
	return n
}
