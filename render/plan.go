package render

import (
	"image"
	"reflect"
	"unsafe"
)

type Plan struct {
	// TODO find a way to remove RootLayer
	//      only used for snapshots?
	RootLayer int

	Shaders  map[int]Shader
	Images   map[int]image.Image
	FaceData FaceData
	AttrData AttrData

	Passes []Pass
}

type FaceData []uint32

func (fd FaceData) GobEncode() ([]byte, error) {
	return uintsToBytes(fd), nil
}

func (fd *FaceData) GobDecode(b []byte) error {
	*fd = bytesToUints(b)
	return nil
}

type AttrData []float32

func (ad AttrData) GobEncode() ([]byte, error) {
	return floatsToBytes(ad), nil
}

func (ad *AttrData) GobDecode(b []byte) error {
	*ad = bytesToFloats(b)
	return nil
}

func uintsToBytes(in []uint32) []byte {
	hdr := *(*reflect.SliceHeader)(unsafe.Pointer(&in))
	hdr.Len *= 4
	hdr.Cap *= 4
	return *(*[]byte)(unsafe.Pointer(&hdr))
}

func bytesToUints(in []byte) []uint32 {
	hdr := *(*reflect.SliceHeader)(unsafe.Pointer(&in))
	hdr.Len /= 4
	hdr.Cap /= 4
	return *(*[]uint32)(unsafe.Pointer(&hdr))
}

func floatsToBytes(in []float32) []byte {
	hdr := *(*reflect.SliceHeader)(unsafe.Pointer(&in))
	hdr.Len *= 4
	hdr.Cap *= 4
	return *(*[]byte)(unsafe.Pointer(&hdr))
}

func bytesToFloats(in []byte) []float32 {
	hdr := *(*reflect.SliceHeader)(unsafe.Pointer(&in))
	hdr.Len /= 4
	hdr.Cap /= 4
	return *(*[]float32)(unsafe.Pointer(&hdr))
}

type Shader struct {
	Vert, Frag, Geom, Output string
}

type Pass struct {
	Name      string
	Shader    int
	Layer     int
	Vertices  int
	Instances int

	Faces    Faces
	Attrs    []Attr
	Uniforms map[string]interface{}
}

type Faces struct {
	Offset int
	// TODO doesn't have the same meaning as Attr.Count
	Count int
}

type Attr struct {
	Name    string
	Offset  int
	Count   int
	Divisor int
}
