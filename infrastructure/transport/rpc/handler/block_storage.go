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

func (h *blockStorage) Put(c context.Context, dto *message.Block) (*rpcmessage.StorageResponse, error) {
	log.FromContext(c).Debugf("[BlockStorage.Put]")
	switch {
	case validation.IsNil(c):
		return nil, fmt.Errorf("Context is nil")
	case validation.IsNil(dto):
		return nil, fmt.Errorf("Block is nil")
	case validation.IsNil(dto.Header):
		return nil, fmt.Errorf("BlockHeader is nil")
	case validation.IsNil(dto.Header.ObjectID):
		return nil, fmt.Errorf("ObjectID is nil")
	case validation.IsNil(dto.Header.BlockID):
		return nil, fmt.Errorf("BlockID is nil")
	case validation.IsNil(dto.Data):
		return nil, fmt.Errorf("Data is nil")
	case validation.IsNil(dto.Header.Checksum):
		return nil, fmt.Errorf("Checksum is nil")
	}

	block := message.ToBlock(dto)
	if err := h.objectStorage.Put(c, &block); err != nil {
		return nil, err
	}

	return &rpcmessage.StorageResponse{
		Success: true,
	}, nil
}

func (h *blockStorage) GetBlock(c context.Context, dto *message.BlockHeader) (*message.Block, error) {
	log.FromContext(c).Debugf("[BlockStorage.GetBlock]")
	switch {
	case validation.IsNil(c):
		return nil, fmt.Errorf("Context is nil")
	case validation.IsNil(dto):
		return nil, fmt.Errorf("Block is nil")
	case validation.IsNil(dto.ObjectID):
		return nil, fmt.Errorf("ObjectID is nil")
	case validation.IsNil(dto.BlockID):
		return nil, fmt.Errorf("BlockID is nil")
	case validation.IsNil(dto.Checksum):
		return nil, fmt.Errorf("Checksum is nil")
	}

	header := message.ToBlockHeader(dto)
	block, err := h.objectStorage.GetBlock(c, header.ObjectID(), header.BlockID(), header.Index())
	if err != nil {
		return nil, err
	}

	return message.FromBlock(block), nil
}

func (h *blockStorage) GetBlockHeader(c context.Context, dto *message.BlockHeader) (*message.BlockHeader, error) {
	log.FromContext(c).Debugf("[BlockStorage.GetBlockHeader]")
	switch {
	case validation.IsNil(c):
		return nil, fmt.Errorf("Context is nil")
	case validation.IsNil(dto):
		return nil, fmt.Errorf("Block is nil")
	case validation.IsNil(dto.ObjectID):
		return nil, fmt.Errorf("ObjectID is nil")
	case validation.IsNil(dto.BlockID):
		return nil, fmt.Errorf("BlockID is nil")
	case validation.IsNil(dto.Checksum):
		return nil, fmt.Errorf("Checksum is nil")
	}

	header := message.ToBlockHeader(dto)
	blockHeader, err := h.objectStorage.GetBlockHeader(c, header.ObjectID(), header.BlockID(), header.Index())
	if err != nil {
		return nil, err
	}

	return message.FromBlockHeader(blockHeader), nil
}

func (h *blockStorage) Delete(c context.Context, dto *message.BlockHeader) (*rpcmessage.StorageResponse, error) {
	log.FromContext(c).Debugf("[BlockStorage.Delete]")
	switch {
	case validation.IsNil(c):
		return nil, fmt.Errorf("Context is nil")
	case validation.IsNil(dto):
		return nil, fmt.Errorf("Block is nil")
	case validation.IsNil(dto.ObjectID):
		return nil, fmt.Errorf("ObjectID is nil")
	case validation.IsNil(dto.BlockID):
		return nil, fmt.Errorf("BlockID is nil")
	case validation.IsNil(dto.Checksum):
		return nil, fmt.Errorf("Checksum is nil")
	}

	header := message.ToBlockHeader(dto)
	if err := h.objectStorage.Delete(c, header.ObjectID(), header.BlockID(), header.Index()); err != nil {
		return nil, err
	}

	return &rpcmessage.StorageResponse{
		Success: true,
	}, nil
}
