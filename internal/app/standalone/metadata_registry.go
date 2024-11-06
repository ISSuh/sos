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

	"github.com/ISSuh/sos/internal/domain/model/dto"
	"github.com/ISSuh/sos/internal/domain/model/message"
	"github.com/ISSuh/sos/internal/domain/service"
	"github.com/ISSuh/sos/internal/infrastructure/transport/rpc"
	"github.com/ISSuh/sos/pkg/validation"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func (s *metadataRegistry) Put(c context.Context, object *message.Object) (*message.ObjectMetadata, error) {
	dto := dto.NewMObjectFromMessage(object)
	resp, err := s.objectMetadata.Put(c, dto)
	if err != nil {
		return nil, err
	}

	return resp.ToMessage(), nil
}

func (r *metadataRegistry) Delete(c context.Context, object *message.Object) (bool, error) {
	dto := dto.NewMObjectFromMessage(object)
	if err := r.objectMetadata.Delete(c, dto); err != nil {
		return false, err
	}
	return true, nil
}

func (s *metadataRegistry) GetByObjectName(c context.Context, req *rpc.ObjectMetadataRequest) (*message.ObjectMetadata, error) {
	item, err := s.objectMetadata.MetadataByObjectName(c, req.Group, req.Partition, req.Path, req.Name)
	if err != nil {
		return nil, err
	}

	return &message.ObjectMetadata{
		Id: &message.ObjectID{
			Id: item.ID.ToInt64(),
		},
		Name:       item.Name,
		Group:      item.Group,
		Partition:  item.Partition,
		Path:       item.Path,
		Versions:   item.Versions.ToMessage(),
		CreatedAt:  timestamppb.New(item.CreatedAt),
		ModifiedAt: timestamppb.New(item.ModifiedAt),
	}, nil
}

func (s *metadataRegistry) GetByObjectID(c context.Context, req *rpc.ObjectMetadataRequest) (*message.ObjectMetadata, error) {
	item, err := s.objectMetadata.MetadataByObjectID(c, req.Group, req.Partition, req.Path, req.GetObjectID())
	if err != nil {
		return nil, err
	}

	return &message.ObjectMetadata{
		Id: &message.ObjectID{
			Id: item.ID.ToInt64(),
		},
		Name:       item.Name,
		Group:      item.Group,
		Partition:  item.Partition,
		Path:       item.Path,
		Versions:   item.Versions.ToMessage(),
		CreatedAt:  timestamppb.New(item.CreatedAt),
		ModifiedAt: timestamppb.New(item.ModifiedAt),
	}, nil
}

func (s *metadataRegistry) FindMetadataOnPath(c context.Context, req *rpc.ObjectMetadataRequest) (*rpc.ObjectMetadataList, error) {
	items, err := s.objectMetadata.MetadataListOnPath(c, req.Group, req.Partition, req.Path)
	if err != nil {
		return nil, err
	}

	list := make([]*message.ObjectMetadata, len(items))
	for i, item := range items {
		list[i] = &message.ObjectMetadata{
			Id: &message.ObjectID{
				Id: item.ID.ToInt64(),
			},
			Name:       item.Name,
			Group:      item.Group,
			Partition:  item.Partition,
			Path:       item.Path,
			Versions:   item.Versions.ToMessage(),
			CreatedAt:  timestamppb.New(item.CreatedAt),
			ModifiedAt: timestamppb.New(item.ModifiedAt),
		}
	}

	return &rpc.ObjectMetadataList{
		Metadata: list,
	}, err
}
