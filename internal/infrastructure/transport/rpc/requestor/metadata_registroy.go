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
	"google.golang.org/protobuf/types/known/emptypb"
)

type metadataRegistry struct {
	logger log.Logger
	engine rpc.MetadataRegistryClient
}

func NewMetadataRegistry(l log.Logger, address string) (rpc.MetadataRegistryRequestor, error) {
	conn, err := sosrpc.NewClientConnection(address)
	if err != nil {
		return nil, err
	}

	return &metadataRegistry{
		logger: l,
		engine: rpc.NewMetadataRegistryClient(conn),
	}, nil
}

func (r *metadataRegistry) Create(c context.Context, in *message.Metadata) (*message.Metadata, error) {
	r.logger.Debugf("[MetadataRegistry.Create]")
	return r.engine.Create(c, in)
}

func (r *metadataRegistry) GetByObjectName(c context.Context, in *message.MetadataFindRequest) (*message.Metadata, error) {
	r.logger.Debugf("[MetadataRegistry.GetByObjectName]")
	return r.engine.GetByObjectName(c, in)
}

func (r *metadataRegistry) GenerateNewObjectID(c context.Context) (*message.ObjectID, error) {
	r.logger.Debugf("[MetadataRegistry.GenerateNewObjectID]")
	e := emptypb.Empty{}
	return r.engine.GenerateNewObjectID(c, &e)
}
