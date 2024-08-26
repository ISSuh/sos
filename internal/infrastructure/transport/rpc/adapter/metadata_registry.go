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
	"github.com/golang/protobuf/ptypes/empty"
)

type MetadataRegistry struct {
	rpc.UnimplementedMetadataRegistryServer
	handler rpc.MetadataRegistryHandler
}

func NewMetadataRegistry(handler rpc.MetadataRegistryHandler) (rpc.Adapter, error) {
	switch {
	case validation.IsNil(handler):
		return nil, fmt.Errorf("handler is nil")
	}

	return &MetadataRegistry{
		handler: handler,
	}, nil
}

func (a *MetadataRegistry) Create(c context.Context, metadata *message.Metadata) (*message.Metadata, error) {
	return a.handler.Create(c, metadata)
}

func (a *MetadataRegistry) GetByObjectName(c context.Context, req *message.MetadataFindRequest) (*message.Metadata, error) {
	return a.handler.GetByObjectName(c, req)
}

func (a *MetadataRegistry) GenerateNewObjectID(c context.Context, _ *empty.Empty) (*message.ObjectID, error) {
	return a.handler.GenerateNewObjectID(c)
}

func (a *MetadataRegistry) Regist() sosrpc.RegisterFunc {
	return func(engine *sosrpc.Engine) {
		rpc.RegisterMetadataRegistryServer(engine.Server, a)
	}
}
