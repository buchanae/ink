package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/buchanae/ink/render"
)

func main() {
	cmd := exec.Command("go", "run", "child.go")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	cmd.Stderr = os.Stderr

	pre := time.Now()

	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	var in io.Reader
	in = stdout
	//in = bufio.NewReaderSize(in, 4096)

	w := &watcher{R: in}
	in = w

	start := time.Now()

	b, err := ioutil.ReadAll(in)
	if err != nil {
		panic(err)
	}

	done := time.Now()

	err = cmd.Wait()
	if err != nil {
		panic(err)
	}

	after := time.Now()

	plan := &render.Plan{}
	_ = b
	_ = plan

	/*
		buf := bytes.NewBuffer(b)
		dec := gob.NewDecoder(buf)
		err = dec.Decode(plan)
		if err != nil {
			panic(err)
		}
	*/

	/*
		err = msgpack.Unmarshal(b, plan)
		if err != nil {
			panic(err)
		}
	*/

	/*
		var mh codec.MsgpackHandle
		mh.MapType = reflect.TypeOf(map[string]interface{}(nil))
		dec := codec.NewDecoderBytes(b, &mh)
		err = dec.Decode(plan)
		if err != nil {
			panic(err)
		}
	*/

	/*
		buf := bytes.NewBuffer(b)
		err = msgp.Decode(buf, plan)
		if err != nil {
			panic(err)
		}
	*/

	log.SetFlags(0)
	log.Printf("bytes read: %.2f MB", float32(w.Count)/1000/1000)
	log.Print("")
	log.Printf("time since start: %s", time.Since(start))
	log.Print("")
	log.Printf("pre until start: %s", start.Sub(pre))
	log.Printf("start until first: %s", w.First.Sub(start))
	log.Printf("first until done: %s", done.Sub(w.First))
	log.Printf("done until after: %s", after.Sub(done))
	log.Printf("decode time: %s", time.Since(after))

}

type watcher struct {
	R     io.Reader
	Count int
	First time.Time
}

func (w *watcher) Read(p []byte) (int, error) {
	n, err := w.R.Read(p)
	if w.Count == 0 {
		w.First = time.Now()
	}
	w.Count += n
	return n, err
}

/*
go run parent.go
encode time: 662.23758ms
bytes read: 52501607

time since start: 1.922489313s
time until first: 1.146882543s

pre until start: 1.246574ms
first until done: 205.914457ms
done until after: 92.812µs
decode time: 569.622345ms
*/
