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

	"github.com/ISSuh/sos/internal/domain/model/entity"
	"github.com/ISSuh/sos/internal/domain/model/message"
	"github.com/ISSuh/sos/internal/infrastructure/transport/rpc"
	"github.com/ISSuh/sos/pkg/empty"
	"github.com/ISSuh/sos/pkg/http"
)

type Downloader struct {
	storageRequestor rpc.BlockStorageRequestor
}

func NewDownloader(storageRequestor rpc.BlockStorageRequestor) Downloader {
	return Downloader{
		storageRequestor: storageRequestor,
	}
}

func (o *Downloader) Download(c context.Context, metadata entity.ObjectMetadata, writer http.DownloadBodyWriter) error {
	blockChan := make([]chan entity.Block, len(metadata.BlockHeaders()))
	for i := range blockChan {
		blockChan[i] = make(chan entity.Block)
	}

	for _, blockHeader := range metadata.BlockHeaders() {
		go func(blockHeader entity.BlockHeader) {
			defer func() {
				if r := recover(); r != nil {
					return
				}
			}()

			block, err := o.downloadBlock(c, blockHeader)
			if err != nil {
				return
			}

			blockChan[blockHeader.Index()] <- block
		}(blockHeader)
	}

	blockSize := 0
	for _, ch := range blockChan {
		block := <-ch

		blockSize += len(block.Buffer())
		err := writer(block.Buffer())
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *Downloader) downloadBlock(c context.Context, blockHeader entity.BlockHeader) (entity.Block, error) {
	msg := message.FromBlockHeader(blockHeader)
	resp, err := o.storageRequestor.GetBlock(c, msg)
	if err != nil {
		return empty.Struct[entity.Block](), err
	}

	return message.ToBlock(resp), nil
}
