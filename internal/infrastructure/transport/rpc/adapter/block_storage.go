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

package adapter

import (
	"context"
	"fmt"

	"github.com/ISSuh/sos/internal/domain/model/message"
	"github.com/ISSuh/sos/internal/infrastructure/transport/rpc"
	sosrpc "github.com/ISSuh/sos/pkg/rpc"
	"github.com/ISSuh/sos/pkg/validation"
)

type BlockStorage struct {
	rpc.UnimplementedBlockStorageServer
	handler rpc.BlockStorageHandler
}

func NewBlockStorage(handler rpc.BlockStorageHandler) (rpc.Adapter, error) {
	switch {
	case validation.IsNil(handler):
		return nil, fmt.Errorf("handler is nil")
	}

	return &BlockStorage{
		handler: handler,
	}, nil
}

func (a *BlockStorage) Put(c context.Context, block *message.Block) (*rpc.StorageResponse, error) {
	return a.handler.Put(c, block)
}

func (a *BlockStorage) Get(c context.Context, header *message.BlockHeader) (*message.Block, error) {
	return a.handler.Get(c, header)
}

func (a *BlockStorage) Delete(c context.Context, header *message.BlockHeader) (*rpc.StorageResponse, error) {
	return a.handler.Delete(c, header)
}

func (a *BlockStorage) Regist() sosrpc.RegisterFunc {
	return func(engine *sosrpc.Engine) {
		rpc.RegisterBlockStorageServer(engine.Server, a)
	}
}
