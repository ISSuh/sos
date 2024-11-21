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

	"github.com/ISSuh/sos/domain/model/message"
	"github.com/ISSuh/sos/domain/service"
	"github.com/ISSuh/sos/infrastructure/transport/rpc"
	rpcmessage "github.com/ISSuh/sos/infrastructure/transport/rpc/message"
	"github.com/ISSuh/sos/internal/validation"
)

type metadataRegistry struct {
	objectMetadata service.ObjectMetadata
}

func NewMetadataRegistry(objectMetadata service.ObjectMetadata) (rpc.MetadataRegistryRequestor, error) {
	switch {
	case validation.IsNil(objectMetadata):
		return nil, fmt.Errorf("ObjectMetadata service is nil")
	}

	return &metadataRegistry{
		objectMetadata: objectMetadata,
	}, nil
}

func (s *metadataRegistry) Put(c context.Context, dto *message.Object) (*message.ObjectMetadata, error) {
	object := message.ToObjectDTO(dto)
	metadata, err := s.objectMetadata.Put(c, object)
	if err != nil {
		return nil, err
	}

	return message.FromObjectMetadataDTO(metadata), nil
}

func (r *metadataRegistry) Delete(c context.Context, dto *message.ObjectMetadata) error {
	metadata := message.ToObjectMetadataDTO(dto)
	if err := r.objectMetadata.Delete(c, metadata); err != nil {
		return err
	}
	return nil
}

func (s *metadataRegistry) GetByObjectName(c context.Context, req *rpcmessage.ObjectMetadataRequest) (*message.ObjectMetadata, error) {
	item, err := s.objectMetadata.MetadataByObjectName(c, req.Group, req.Partition, req.Path, req.Name)
	if err != nil {
		return nil, err
	}

	return message.FromObjectMetadataDTO(item), nil
}

func (s *metadataRegistry) GetByObjectID(c context.Context, req *rpcmessage.ObjectMetadataRequest) (*message.ObjectMetadata, error) {
	item, err := s.objectMetadata.MetadataByObjectID(c, req.Group, req.Partition, req.Path, req.GetObjectID())
	if err != nil {
		return nil, err
	}

	return message.FromObjectMetadataDTO(item), nil
}

func (s *metadataRegistry) FindMetadataOnPath(c context.Context, req *rpcmessage.ObjectMetadataRequest) (*message.ObjectMetadataList, error) {
	items, err := s.objectMetadata.MetadataListOnPath(c, req.Group, req.Partition, req.Path)
	if err != nil {
		return nil, err
	}

	return message.FromObjectMetadataListDTO(items), err
}
