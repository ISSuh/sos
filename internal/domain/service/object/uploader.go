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

package object

import (
	"bytes"
	"context"
	"errors"
	"io"
	"time"

	"github.com/ISSuh/sos/internal/domain/model/dto"
	"github.com/ISSuh/sos/internal/domain/model/entity"
	"github.com/ISSuh/sos/internal/infrastructure/transport/rpc"
	"github.com/ISSuh/sos/pkg/crc"
	"github.com/ISSuh/sos/pkg/log"
)

const (
	kilobyte          = 1024
	defaultBufferSize = 4 * kilobyte
)

type Uploader struct {
	storageRequestor rpc.BlockStorageRequestor
}

func NewUploader(storageRequestor rpc.BlockStorageRequestor) Uploader {
	return Uploader{
		storageRequestor: storageRequestor,
	}
}

func (o *Uploader) Upload(c context.Context, objectID entity.ObjectID, bodyStream io.ReadCloser) (dto.BlockHeaders, error) {
	var blockheaders dto.BlockHeaders
	var totalReadSize int
	var blockIndex int

	var blockBuffer bytes.Buffer
	for {
		buf := make([]byte, defaultBufferSize)
		n, err := bodyStream.Read(buf)
		if err != nil && err != io.EOF {
			return nil, err
		}

		if n == 0 {
			blockBuffer.Truncate(totalReadSize)

			block := o.buildBlock(objectID, blockIndex, blockBuffer.Bytes())
			if err := o.uploadBlock(c, block); err != nil {
				return nil, err
			}

			blockheaders = append(blockheaders, block.Header)
			break
		}

		// resize buffer
		buf = buf[:n]

		totalReadSize += n
		writeSize, err := blockBuffer.Write(buf)
		if err != nil {
			return nil, err
		}

		if writeSize != n {
			return nil, errors.New("write size is not equal to read size")
		}

		if totalReadSize >= entity.BlockSize {
			block := o.buildBlock(objectID, blockIndex, blockBuffer.Bytes())
			if err := o.uploadBlock(c, block); err != nil {
				return nil, err
			}

			blockheaders = append(blockheaders, block.Header)

			blockIndex++
			totalReadSize = 0
		}
	}
	return blockheaders, nil
}

func (o *Uploader) buildBlock(objectID entity.ObjectID, index int, buffer []byte) dto.Block {
	block := dto.Block{
		Header: dto.BlockHeader{
			ObjectID:  objectID,
			BlockID:   entity.NewBlockID(),
			Index:     index,
			Size:      len(buffer),
			Timestamp: time.Now(),
			Checksum:  crc.Checksum(buffer),
		},
		Data: make([]byte, len(buffer)),
	}

	copy(block.Data, buffer)
	return block
}

func (o *Uploader) uploadBlock(c context.Context, block dto.Block) error {
	msg := block.ToMessage()
	resp, err := o.storageRequestor.Put(c, msg)
	if err != nil {
		log.FromContext(c).Errorf("Upload Error: %s", err.Error())
		return err
	}

	if !resp.Success {
		log.FromContext(c).Errorf("Upload fail. message : %s", resp.Message)
		return err
	}
	return nil
}
