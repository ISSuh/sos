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
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/alexflint/go-arg"
)

const (
	CrcSize    = 4
	BodySize   = 4092
	BufferSize = BodySize + CrcSize

	// chunkSize     = 1 * 1024 * 1024 // 1MB
	chunkSize     = 30 * 1024 // 30KB
	ChunkEncoding = "chunked"

	ObjectNameQueryName = "name"
	ObjectSizeQueryName = "size"
	ChunkSizeQueryName  = "chunk_size"

	ObjectNameHeader = "X-SOS-Object-Name"
	ObjectSizeHeader = "X-SOS-Object-Size"
	ChunkSizeHeader  = "X-SOS-Chunk-Size"
)

var args struct {
	Server     string `arg:"-s,--server,required"`
	Group      string `arg:"-g,--group,required"`
	Partition  string `arg:"-p,--partition,required"`
	ObjectPath string `arg:"-o,--path,required"`
	File       string `arg:"-f,--file,required"`
}

func uploadChunked(url string, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file stat: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.TransferEncoding = []string{"chunked"}
	req.Header.Set("Content-Type", "application/octet-stream")

	q := req.URL.Query()
	q.Add(ObjectNameQueryName, path.Base(filePath))
	q.Add(ObjectSizeQueryName, strconv.FormatInt(fileStat.Size(), 10))
	q.Add(ChunkSizeQueryName, strconv.Itoa(chunkSize))
	req.URL.RawQuery = q.Encode()

	pr, pw := io.Pipe()
	req.Body = pr

	go func() {
		defer pw.Close()
		buf := make([]byte, chunkSize)
		for {
			n, err := file.Read(buf)
			if err != nil && err != io.EOF {
				pw.CloseWithError(fmt.Errorf("failed to read chunk: %v", err))
				return
			}
			if n == 0 {
				break
			}

			_, err = pw.Write(buf[:n])
			if err != nil {
				pw.CloseWithError(fmt.Errorf("failed to write chunk: %v", err))
				return
			}
		}
	}()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned non-200 status: %v", resp.Status)
	}

	return nil
}

func main() {
	arg.MustParse(&args)

	// set url
	url := fmt.Sprintf("http://%s/v1/%s/%s/%s", args.Server, args.Group, args.Partition, args.ObjectPath)

	err := uploadChunked(url, args.File)
	if err != nil {
		fmt.Printf("Error uploading file: %v\n", err)
	} else {
		fmt.Println("File uploaded successfully")
	}
	fmt.Println("File uploaded successfully")
}
