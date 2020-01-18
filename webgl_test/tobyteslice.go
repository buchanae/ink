package main

import (
	"log"
	"reflect"
	"unsafe"
)

func main() {
	data := []uint32{0, 0, 1}
	b := toByteSlice(data, 12)
	log.Printf("%v", b)
}

func toByteSlice(s interface{}, size int) []byte {
	log.Printf("toByteslice %#v %d", s, size)

	type iface struct {
		Type, Data unsafe.Pointer
	}

	ii := *(*iface)(unsafe.Pointer(&s))
	h := (*reflect.SliceHeader)(ii.Data)
	h.Len = size
	h.Cap = size
	bs := *(*[]byte)(unsafe.Pointer(h))
	return bs
}
