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

	"github.com/ISSuh/sos/domain/model/message"
	"github.com/ISSuh/sos/infrastructure/transport/rpc"
	rpcmessage "github.com/ISSuh/sos/infrastructure/transport/rpc/message"
	"github.com/ISSuh/sos/internal/log"
	sosrpc "github.com/ISSuh/sos/internal/rpc"
)

type blockStorage struct {
	engine rpcmessage.BlockStorageClient
}

func NewBlockStorage(address string) (rpc.BlockStorageRequestor, error) {
	conn, err := sosrpc.NewClientConnection(address)
	if err != nil {
		return nil, err
	}

	return &blockStorage{
		engine: rpcmessage.NewBlockStorageClient(conn),
	}, nil
}

func (r *blockStorage) Put(c context.Context, block *message.Block) (*rpcmessage.StorageResponse, error) {
	log.FromContext(c).Debugf("[BlockStorage.Put]")
	return r.engine.Put(c, block)
}

func (r *blockStorage) GetBlock(c context.Context, header *message.BlockHeader) (*message.Block, error) {
	log.FromContext(c).Debugf("[BlockStorage.Get]")
	return r.engine.GetBlock(c, header)
}

func (r *blockStorage) GetBlockHeader(c context.Context, header *message.BlockHeader) (*message.BlockHeader, error) {
	log.FromContext(c).Debugf("[BlockStorage.Get]")
	return r.engine.GetBlockHeader(c, header)
}

func (r *blockStorage) Delete(c context.Context, header *message.BlockHeader) (*rpcmessage.StorageResponse, error) {
	log.FromContext(c).Debugf("[BlockStorage.Delete]")
	return r.engine.Delete(c, header)
}
