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

package handler

import (
	"context"
	"fmt"

	"github.com/ISSuh/sos/domain/model/message"
	"github.com/ISSuh/sos/domain/service"
	"github.com/ISSuh/sos/infrastructure/transport/rpc"
	rpcmessage "github.com/ISSuh/sos/infrastructure/transport/rpc/message"
	"github.com/ISSuh/sos/internal/log"
	"github.com/ISSuh/sos/internal/validation"
)

type blockStorage struct {
	objectStorage service.ObjectStorage
}

func NewBlockStorage(objectStorage service.ObjectStorage) (rpc.BlockStorageHandler, error) {
	switch {
	case validation.IsNil(objectStorage):
		return nil, fmt.Errorf("ObjectStorage service is nil")
	}

	return &blockStorage{
		objectStorage: objectStorage,
	}, nil
}

func (h *blockStorage) Put(c context.Context, block *message.Block) (*rpcmessage.StorageResponse, error) {
	log.FromContext(c).Debugf("[BlockStorage.Put]")
	return &rpcmessage.StorageResponse{}, nil
}

func (h *blockStorage) GetBlock(c context.Context, header *message.BlockHeader) (*message.Block, error) {
	log.FromContext(c).Debugf("[BlockStorage.GetBlock]")
	return &message.Block{}, nil
}

func (h *blockStorage) GetBlockHeader(c context.Context, header *message.BlockHeader) (*message.BlockHeader, error) {
	log.FromContext(c).Debugf("[BlockStorage.GetBlockHeader]")
	return &message.BlockHeader{}, nil
}

func (h *blockStorage) Delete(c context.Context, header *message.BlockHeader) (*rpcmessage.StorageResponse, error) {
	log.FromContext(c).Debugf("[BlockStorage.Delete]")
	return &rpcmessage.StorageResponse{}, nil
}
