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

type metadataRegistry struct {
	engine rpcmessage.MetadataRegistryClient
}

func NewMetadataRegistry(address string) (rpc.MetadataRegistryRequestor, error) {
	conn, err := sosrpc.NewClientConnection(address)
	if err != nil {
		return nil, err
	}

	return &metadataRegistry{
		engine: rpcmessage.NewMetadataRegistryClient(conn),
	}, nil
}

func (r *metadataRegistry) Put(c context.Context, object *message.Object) (*message.ObjectMetadata, error) {
	log.FromContext(c).Debugf("[MetadataRegistry.Put]")
	return r.engine.Put(c, object)
}

func (r *metadataRegistry) Delete(c context.Context, metadata *message.ObjectMetadata) (bool, error) {
	log.FromContext(c).Debugf("[MetadataRegistry.Delete]")
	res, err := r.engine.Delete(c, metadata)
	if err != nil {
		return false, err
	}

	return res.Value, nil
}

func (r *metadataRegistry) GetByObjectName(c context.Context, req *rpcmessage.ObjectMetadataRequest) (*message.ObjectMetadata, error) {
	log.FromContext(c).Debugf("[MetadataRegistry.GetByObjectName]")
	return r.engine.GetByObjectName(c, req)
}

func (r *metadataRegistry) GetByObjectID(c context.Context, req *rpcmessage.ObjectMetadataRequest) (*message.ObjectMetadata, error) {
	log.FromContext(c).Debugf("[MetadataRegistry.GetByObjectID]")
	return r.engine.GetByObjectID(c, req)
}

func (r *metadataRegistry) FindMetadataOnPath(c context.Context, req *rpcmessage.ObjectMetadataRequest) (*rpcmessage.ObjectMetadataList, error) {
	log.FromContext(c).Debugf("[MetadataRegistry.FindMetadataOnPath]")
	return r.engine.FindMetadataOnPath(c, req)
}
