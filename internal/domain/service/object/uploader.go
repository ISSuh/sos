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
	"context"
	"io"
	"time"

	"github.com/ISSuh/sos/internal/domain/model/entity"
	"github.com/ISSuh/sos/internal/domain/model/message"
	"github.com/ISSuh/sos/internal/infrastructure/transport/rpc"
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

func (o *Uploader) Upload(c context.Context, objectID entity.ObjectID, bodyStream io.ReadCloser) (entity.BlockHeaders, error) {
	var blockheaders entity.BlockHeaders
	var totalReadSize uint64
	var blockIndex int
	blockBuilder := entity.NewBlockBuilder()

	for {
		buf := make([]byte, defaultBufferSize)
		n, err := bodyStream.Read(buf)
		if err != nil && err != io.EOF {
			return nil, err
		}

		if n == 0 {
			blockBuilder.ReSizeBuffer(totalReadSize)

			block := o.buildBlock(objectID, blockIndex, blockBuilder)
			if err := o.uploadBlock(c, block); err != nil {
				return nil, err
			}

			blockheaders = append(blockheaders, block.Header())
			break
		}

		totalReadSize += uint64(n)
		blockBuilder.AppendBuffer(buf)

		if totalReadSize >= entity.BlockSize {
			block := o.buildBlock(objectID, blockIndex, blockBuilder)
			if err := o.uploadBlock(c, block); err != nil {
				return nil, err
			}

			blockheaders = append(blockheaders, block.Header())

			blockBuilder = entity.NewBlockBuilder()
			blockIndex++
			totalReadSize = 0
		}

	}

	log.FromContext(c).Infof("Upload objectID: %s, block count: %d", objectID.String(), len(blockheaders))
	return blockheaders, nil
}

func (o *Uploader) buildBlock(objectID entity.ObjectID, index int, blockBuilder *entity.BlockBuilder) entity.Block {
	c := blockBuilder.CalculateChecksum()
	header := entity.NewBlockHeaderBuilder().
		ObjectID(objectID).
		BlockID(entity.NewBlockID()).
		Index(index).
		Size(blockBuilder.BufferSize()).
		Node(entity.Node{}).
		Timestamp(time.Now()).
		Checksum(c).
		Build()

	blockBuilder.Header(header)
	return blockBuilder.Build()
}

func (o *Uploader) uploadBlock(c context.Context, block entity.Block) error {
	msg := message.FromBlock(block)
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
