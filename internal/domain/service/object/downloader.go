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
)

type Downloader struct {
	storageRequestor rpc.BlockStorageRequestor
}

func NewDownloader(storageRequestor rpc.BlockStorageRequestor) Uploader {
	return Uploader{
		storageRequestor: storageRequestor,
	}
}

func (o *Downloader) Download(c context.Context, metadata entity.ObjectMetadata) (entity.Blocks, error) {
	blocks := make(entity.Blocks, len(metadata.BlockHeaders()))
	for _, blockHeader := range metadata.BlockHeaders() {
		block, err := o.downloadBlock(c, blockHeader)
		if err != nil {
			return nil, err
		}

		blocks[block.Index()] = block
	}

	return blocks, nil
}

func (o *Downloader) downloadBlock(c context.Context, blockHeader entity.BlockHeader) (entity.Block, error) {
	msg := message.FromBlockHeader(blockHeader)
	resp, err := o.storageRequestor.Get(c, msg)
	if err != nil {
		return empty.Struct[entity.Block](), err
	}

	return message.ToBlock(resp), nil
}
