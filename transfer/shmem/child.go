package main

import (
	"reflect"
	"syscall"
	"unsafe"

	"github.com/buchanae/ink/trac"
	"github.com/buchanae/ink/transfer/flowfield"
)

func main() {
	trac.Enabled = true

	trac.Log("doc")
	doc := flowfield.Doc()

	trac.Log("plan")
	plan := doc.Plan()

	trac.Log("make buffer")
	data := make([]byte, 0, 51_000_000)

	trac.Log("encode")
	data = append(data, uintsToBytes(plan.FaceData)...)
	data = append(data, floatsToBytes(plan.AttrData)...)

	length := len(data)

	trac.Log("truncate")
	fd := 3
	err := syscall.Ftruncate(fd, int64(length))
	if err != nil {
		panic(err)
	}

	trac.Log("mmap")
	b, err := syscall.Mmap(
		fd, 0, length,
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_SHARED,
	)

	trac.Log("copy")
	copy(b, data)

	trac.Log("done: %v", err)
}

func uintsToBytes(in []uint32) []byte {
	hdr := *(*reflect.SliceHeader)(unsafe.Pointer(&in))
	hdr.Len *= 4
	hdr.Cap *= 4
	return *(*[]byte)(unsafe.Pointer(&hdr))
}

func floatsToBytes(in []float32) []byte {
	hdr := *(*reflect.SliceHeader)(unsafe.Pointer(&in))
	hdr.Len *= 4
	hdr.Cap *= 4
	return *(*[]byte)(unsafe.Pointer(&hdr))
}
