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

package standalone

import (
	"context"
	"fmt"

	"github.com/ISSuh/sos/internal/domain/model/entity"
	"github.com/ISSuh/sos/internal/domain/model/message"
	"github.com/ISSuh/sos/internal/domain/service"
	"github.com/ISSuh/sos/internal/infrastructure/transport/rpc"
	"github.com/ISSuh/sos/pkg/validation"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type blockStorage struct {
	objectStorage service.ObjectStorage
}

func NewBlockStorage(objectStorage service.ObjectStorage) (rpc.BlockStorageRequestor, error) {
	switch {
	case validation.IsNil(objectStorage):
		return nil, fmt.Errorf("ObjectStorage service is nil")
	}

	return &blockStorage{
		objectStorage: objectStorage,
	}, nil
}

func (s *blockStorage) Put(
	ctx context.Context, blockMessage *message.Block,
) (*rpc.StorageResponse, error) {
	header :=
		entity.NewBlockHeaderBuilder().
			BlockID(entity.BlockID(blockMessage.Header.BlockID.Id)).
			ObjectID(entity.ObjectID(blockMessage.Header.ObjectID.Id)).
			Index(int(blockMessage.Header.Index)).
			Checksum(blockMessage.Header.Checksum).
			Timestamp(blockMessage.Header.Timestamp.AsTime()).
			Build()

	block :=
		entity.NewBlockBuilder().
			Buffer(blockMessage.Data).
			Header(header).
			Build()

	if err := s.objectStorage.Put(ctx, block); err != nil {
		return &rpc.StorageResponse{
			Success: false,
			Message: err.Error(),
		}, err
	}

	return &rpc.StorageResponse{
		Success: true,
	}, nil
}

func (s *blockStorage) GetBlock(
	ctx context.Context, headerMessage *message.BlockHeader,
) (*message.Block, error) {
	block, err := s.objectStorage.GetBlock(
		ctx, entity.ObjectID(headerMessage.ObjectID.Id),
		entity.BlockID(headerMessage.BlockID.Id), int(headerMessage.Index),
	)

	if err != nil {
		return nil, err
	}

	blockHeader := block.Header()
	return &message.Block{
		Header: &message.BlockHeader{
			ObjectID: &message.ObjectID{
				Id: blockHeader.ObjectID().ToInt64(),
			},
			BlockID: &message.BlockID{
				Id: blockHeader.ObjectID().ToInt64(),
			},
			Index:     int32(blockHeader.Index()),
			Timestamp: timestamppb.New(blockHeader.Timestamp()),
		},
		Data: block.Buffer(),
	}, nil
}

func (s *blockStorage) GetBlockHeader(
	ctx context.Context, headerMessage *message.BlockHeader,
) (*message.BlockHeader, error) {
	blockHeader, err := s.objectStorage.GetBlockHeader(
		ctx, entity.ObjectID(headerMessage.ObjectID.Id),
		entity.BlockID(headerMessage.BlockID.Id), int(headerMessage.Index),
	)

	if err != nil {
		return nil, err
	}

	return &message.BlockHeader{
		ObjectID: &message.ObjectID{
			Id: blockHeader.ObjectID().ToInt64(),
		},
		BlockID: &message.BlockID{
			Id: blockHeader.BlockID().ToInt64(),
		},
		Index:     int32(blockHeader.Index()),
		Size:      int32(blockHeader.Size()),
		Timestamp: timestamppb.New(blockHeader.Timestamp()),
	}, nil
}

func (s *blockStorage) Delete(
	ctx context.Context, headerMessage *message.BlockHeader,
) (*rpc.StorageResponse, error) {
	err := s.objectStorage.Delete(
		ctx, entity.ObjectID(headerMessage.ObjectID.Id),
		entity.BlockID(headerMessage.BlockID.Id), int(headerMessage.Index),
	)

	if err != nil {
		return &rpc.StorageResponse{
			Success: false,
			Message: err.Error(),
		}, err
	}
	return &rpc.StorageResponse{
		Success: true,
	}, nil
}
