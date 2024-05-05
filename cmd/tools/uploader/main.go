// MIT License

// Copyright (c) 2024 ISSuh

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/ISSuh/sos/internal/object"
	"github.com/alexflint/go-arg"
)

const (
	CrcSize       = 4
	BodySize      = 4092
	BufferSize    = BodySize + CrcSize
	ChunkEncoding = "chunked"
)

var args struct {
	Server    string `arg:"-s,--server,required"`
	Group     string `arg:"-g,--group,required"`
	Partition string `arg:"-p,--partision,required"`
	File      string `arg:"-f,--file,required"`
}

func fileReadAndWriteToChan(f *os.File, writer *io.PipeWriter) (bool, int, error) {
	buf := make([]byte, BodySize)
	n, err := f.Read(buf)
	if err == io.EOF {
		fmt.Printf("end of file")
		return true, n, nil
	}

	if err != nil {
		return false, 0, err
	}

	crc := crc32.ChecksumIEEE(buf)
	buf = binary.LittleEndian.AppendUint32(buf, crc)

	fmt.Printf("b : %d / %+v\n", len(buf), buf[BodySize:])

	n, err = writer.Write(buf)
	return false, n, err
}

func fileTask(file string, writer *io.PipeWriter, q <-chan struct{}, wg *sync.WaitGroup) {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}

	defer f.Close()
	defer wg.Done()
	defer writer.Close()

	totalSize := 0
	isEOF := false
	for !isEOF {
		select {
		case <-q:
			break
		default:
			n := 0
			isEOF, n, err = fileReadAndWriteToChan(f, writer)
			if err != nil {
				panic(err)
			}

			totalSize += n
		}
	}

	fmt.Printf("total size : %d", totalSize)
	time.Sleep(1 * time.Second)
}

func main() {
	arg.MustParse(&args)

	// set url
	fileName := filepath.Base(args.File)
	urlStr := fmt.Sprintf("http://%s/v1/%s/%s/%s", args.Server, args.Group, args.Partition, fileName)
	u, err := url.Parse(urlStr)
	if err != nil {
		fmt.Printf("err : %s\n", err.Error())
		return
	}

	// open pipe
	reader, writer := io.Pipe()

	req := &http.Request{
		Method:           http.MethodPut,
		URL:              u,
		TransferEncoding: []string{ChunkEncoding},
		Body:             reader,
		Header:           make(map[string][]string),
	}

	var wg sync.WaitGroup
	wg.Add(1)

	q := make(chan struct{})
	// go fileTask(args.File, writer, q, &wg)
	go func() {
		f, err := os.Open(args.File)
		if err != nil {
			panic(err)
		}

		defer f.Close()
		defer wg.Done()
		defer writer.Close()

		time.Sleep(1 * time.Second)

		stat, err := os.Stat(args.File)
		if err != nil {
			panic(err)
		}

		header := object.Metadata{
			Name:      fileName,
			Group:     args.Group,
			Partition: args.Partition,
			Size:      uint32(stat.Size()),
		}

		j, err := json.Marshal(header)
		if err != nil {
			panic(err)
		}

		n, err := writer.Write(j)
		if err != nil {
			panic(err)
		}

		fmt.Printf("write header len : %d, header : %+v\n", n, header)

		totalSize := 0
		isEOF := false
		for !isEOF {
			select {
			case <-q:
				break
			default:
				n := 0
				isEOF, n, err = fileReadAndWriteToChan(f, writer)
				if err != nil {
					panic(err)
				}

				totalSize += n
			}
		}

		fmt.Printf("total size : %d", totalSize)
		time.Sleep(1 * time.Second)
	}()

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	wg.Wait()

	body, err := io.ReadAll(resp.Body)
	if nil != err {
		fmt.Println("error =>", err.Error())
	} else {
		fmt.Println(string(body))
	}
}
