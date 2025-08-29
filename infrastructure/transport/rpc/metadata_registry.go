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

package rpc

import (
	"context"

	message "github.com/ISSuh/sos/domain/model/message"
	rpcmessage "github.com/ISSuh/sos/infrastructure/transport/rpc/message"
)

type MetadataRegistryHandler interface {
	Put(c context.Context, object *message.Object) (*message.ObjectMetadata, error)
	Delete(c context.Context, metadata *message.ObjectMetadata) error
	GetByObjectName(c context.Context, req *rpcmessage.ObjectMetadataRequest) (*message.ObjectMetadata, error)
	GetByObjectID(c context.Context, req *rpcmessage.ObjectMetadataRequest) (*message.ObjectMetadata, error)
	FindMetadataOnPath(c context.Context, req *rpcmessage.ObjectMetadataRequest) (*message.ObjectMetadataList, error)
}

type MetadataRegistryRequestor interface {
	Put(c context.Context, object *message.Object) (*message.ObjectMetadata, error)
	Delete(c context.Context, metadata *message.ObjectMetadata) error
	GetByObjectName(c context.Context, req *rpcmessage.ObjectMetadataRequest) (*message.ObjectMetadata, error)
	GetByObjectID(c context.Context, req *rpcmessage.ObjectMetadataRequest) (*message.ObjectMetadata, error)
	FindMetadataOnPath(c context.Context, rew *rpcmessage.ObjectMetadataRequest) (*message.ObjectMetadataList, error)
}
