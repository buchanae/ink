package main

import (
	"log"
	"os"
	"os/exec"
	"runtime"
	"syscall"

	"github.com/fabiokung/shm"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	log.SetFlags(0)

	file, err := shm.Open("my_region", os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	defer shm.Unlink(file.Name())

	cmd := exec.Command("go", "run", "child.go")
	cmd.ExtraFiles = []*os.File{file}
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		panic(err)
	}

	st, err := file.Stat()
	log.Print(err)
	log.Printf("%#v", st.Size())

	fd := file.Fd()
	log.Printf("fd: %v", fd)

	b, err := syscall.Mmap(
		int(fd), 0, int(st.Size()),
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_SHARED,
	)
	log.Printf("ERR: %v", err)
	log.Printf("LEN: %v", len(b))
	//log.Print(string(b))
}

// func Mmap(fd int, offset int64, length int, prot int, flags int) (data []byte, err error)
