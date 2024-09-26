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

func (r *blockStorage) Put(ctx context.Context, msg *message.Block) (*rpc.StorageResponse, error) {
	blockHeader := r.blockHeaderFromMessage(msg)
	blockBuilder := entity.NewBlockBuilder()
	blockBuilder.
		Buffer(msg.Data).
		Header(blockHeader)

	err := r.objectStorage.Put(ctx, blockBuilder.Build())
	if err != nil {
		return nil, err
	}

	resp := &rpc.StorageResponse{
		Success: true,
	}
	return resp, nil
}

func (r *blockStorage) Get(ctx context.Context, msg *message.BlockHeader) (*message.Block, error) {

	return nil, nil
}

func (r *blockStorage) Delete(ctx context.Context, msg *message.BlockHeader) (*rpc.StorageResponse, error) {
	return nil, nil
}

func (r *blockStorage) blockHeaderFromMessage(msg *message.Block) entity.BlockHeader {
	header := msg.GetHeader()
	headerBuilder := entity.NewBlockHeaderBuilder()
	headerBuilder.
		BlockID(entity.BlockID(header.BlockID.Id)).
		ObjectID(entity.ObjectID(header.ObjectID.Id)).
		Index(int(header.Index)).
		Size(int(header.Size)).
		Timestamp(header.Timestamp.AsTime()).
		Checksum(header.Checksum)

	return headerBuilder.Build()

}
