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

package requestor

import (
	"context"

	"github.com/ISSuh/sos/internal/domain/model/message"
	"github.com/ISSuh/sos/internal/infrastructure/transport/rpc"
	"github.com/ISSuh/sos/pkg/log"
	sosrpc "github.com/ISSuh/sos/pkg/rpc"
)

type blockStorage struct {
	logger log.Logger
	engine rpc.BlockStorageClient
}

func NewBlockStorage(l log.Logger, address string) (rpc.BlockStorageRequestor, error) {
	conn, err := sosrpc.NewClientConnection(address)
	if err != nil {
		return nil, err
	}

	return &blockStorage{
		logger: l,
		engine: rpc.NewBlockStorageClient(conn),
	}, nil
}

func (r *blockStorage) Put(ctx context.Context, block *message.Block) (*rpc.StorageResponse, error) {
	r.logger.Debugf("[BlockStorage.Put]")
	return r.engine.Put(ctx, block)
}

func (r *blockStorage) Get(ctx context.Context, header *message.BlockHeader) (*message.Block, error) {
	r.logger.Debugf("[BlockStorage.Get]")
	return r.engine.Get(ctx, header)
}

func (r *blockStorage) Delete(ctx context.Context, header *message.BlockHeader) (*rpc.StorageResponse, error) {
	r.logger.Debugf("[BlockStorage.Delete]")
	return r.engine.Delete(ctx, header)
}
