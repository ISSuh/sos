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

	"github.com/ISSuh/sos/internal/infrastructure/transport/rpc"
	"github.com/ISSuh/sos/internal/infrastructure/transport/rpc/message"
	sosrpc "github.com/ISSuh/sos/pkg/rpc"

	"github.com/golang/protobuf/ptypes/empty"
)

type metadataRegistry struct {
	rpc.UnimplementedMetadataRegistryServer
}

func NewMetadataRegistry() rpc.MetadataRegistryServer {
	return &metadataRegistry{}
}

func (h *metadataRegistry) Create(context.Context, *message.Metadata) (*message.Metadata, error) {
	return nil, nil
}

func (h *metadataRegistry) GetByObjectName(context.Context, *message.MetadataFindRequest) (*message.Metadata, error) {
	return nil, nil
}

func (h *metadataRegistry) GenerateNewObjectID(context.Context, *empty.Empty) (*message.Metadata_ObjectID, error) {
	return nil, nil
}

func RegistMetadataRegistry(handler rpc.MetadataRegistryServer) sosrpc.RegisterFunc {
	return func(engine *sosrpc.Engine) {
		rpc.RegisterMetadataRegistryServer(engine.Server, handler)
	}
}
