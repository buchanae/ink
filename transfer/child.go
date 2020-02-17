package main

import (
	"bytes"
	"encoding/gob"
	"io"
	"log"
	"os"
	"reflect"
	"time"
	"unsafe"

	"github.com/buchanae/ink/transfer/flowfield"
)

func main() {
	log.SetFlags(0)

	var out io.Writer
	out = os.Stdout
	//out = bufio.NewWriterSize(out, 4096)

	/*
		plan := render.Plan{
			FaceData: make([]uint32, 7_000_000),
			AttrData: make([]float32, 7_000_000),
			Passes:   make([]render.Pass, 1000),
		}
		for i := range plan.Passes {
			plan.Passes[i].Name = "name"
		}
	*/
	plan := flowfield.Plan()

	b := make([]byte, 0, 251_000_000)
	buf := bytes.NewBuffer(b)
	enc := gob.NewEncoder(buf)

	start := time.Now()

	/*
		b = append(b, uintsToBytes(plan.FaceData)...)
		b = append(b, floatsToBytes(plan.AttrData)...)
	*/

	err := enc.Encode(plan)
	if err != nil {
		panic(err)
	}

	/*
		buf := bytes.NewBuffer(b)
		enc := gob.NewEncoder(buf)
		err := enc.Encode(plan.FaceData)
		if err != nil {
			panic(err)
		}
		err = enc.Encode(plan.AttrData)
		if err != nil {
			panic(err)
		}
	*/

	/*
		buf := bytes.NewBuffer(b)
		err := binary.Write(buf, binary.LittleEndian, plan.FaceData)
		if err != nil {
			panic(err)
		}
		err = binary.Write(buf, binary.LittleEndian, plan.AttrData)
		if err != nil {
			panic(err)
		}
	*/

	/*
		b, err := msgpack.Marshal(plan)
		if err != nil {
			panic(err)
		}
	*/

	/*
		var mh codec.MsgpackHandle
		mh.MapType = reflect.TypeOf(map[string]interface{}(nil))
		var b []byte
		enc := codec.NewEncoderBytes(&b, &mh)
		err := enc.Encode(plan)
		if err != nil {
			panic(err)
		}
	*/

	/*
		b := make([]byte, 0, 51000000)
		buf := bytes.NewBuffer(b)
		err := msgp.Encode(buf, &plan)
		if err != nil {
			panic(err)
		}
	*/

	log.Printf("encode time: %s", time.Since(start))

	out.Write(buf.Bytes())

	//out.Write(b)
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
