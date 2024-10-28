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
	"github.com/ISSuh/sos/internal/domain/model/entity"
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

func (s *metadataRegistry) Put(c context.Context, metadata *message.ObjectMetadata) (*message.ObjectMetadata, error) {
	blockHeaders := entity.BlockHeaders{}
	for _, headers := range metadata.BlockHeaders {
		builder :=
			entity.NewBlockHeaderBuilder().
				BlockID(entity.BlockID(headers.BlockID.Id)).
				ObjectID(entity.ObjectID(headers.ObjectID.Id)).
				Index(int(headers.Index)).
				Size(int(headers.Size)).
				Timestamp(headers.Timestamp.AsTime())

		blockHeaders = append(blockHeaders, builder.Build())
	}

	d := dto.Request{
		ObjectID:     entity.NewObjectIDFrom(metadata.GetId().Id),
		Name:         metadata.GetName(),
		Group:        metadata.GetGroup(),
		Partition:    metadata.GetPartition(),
		Path:         metadata.GetPath(),
		Size:         int(metadata.GetSize()),
		BlockHeaders: blockHeaders,
	}

	err := s.objectMetadata.Create(c, d)
	if err != nil {
		return nil, err
	}

	return metadata, nil
}

func (r *metadataRegistry) Delete(c context.Context, metadata *message.ObjectMetadata) (bool, error) {
	return false, nil
}

func (s *metadataRegistry) GetByObjectName(c context.Context, req *rpc.ObjectMetadataRequest) (*message.ObjectMetadata, error) {
	d := dto.Request{
		Group:     req.GetGroup(),
		Partition: req.GetPartition(),
		Path:      req.GetPath(),
		Name:      req.GetName(),
	}

	item, err := s.objectMetadata.MetadataByObjectName(c, d)
	if err != nil {
		return nil, err
	}

	return &message.ObjectMetadata{
		Id: &message.ObjectID{
			Id: item.ID.ToInt64(),
		},
		Name:      item.Name,
		Group:     item.Group,
		Partition: item.Partition,
		Path:      item.Path,
		Size:      int32(item.Size),
	}, nil
}

func (s *metadataRegistry) GetByObjectID(c context.Context, req *rpc.ObjectMetadataRequest) (*message.ObjectMetadata, error) {
	d := dto.Request{
		ObjectID:  entity.NewObjectIDFrom(req.GetObjectID()),
		Group:     req.GetGroup(),
		Partition: req.GetPartition(),
		Path:      req.GetPath(),
		Name:      req.GetName(),
	}

	item, err := s.objectMetadata.MetadataByObjectID(c, d)
	if err != nil {
		return nil, err
	}

	blockHeaders := make([]*message.BlockHeader, 0)
	for _, header := range item.BlockHeaders {
		blockHeader := &message.BlockHeader{
			ObjectID: &message.ObjectID{
				Id: header.ObjectID.ToInt64(),
			},
			BlockID: &message.BlockID{
				Id: header.BlockID.ToInt64(),
			},
			Index:     int32(header.Index),
			Size:      int32(header.Size),
			Timestamp: timestamppb.New(header.Timestamp),
		}

		blockHeaders = append(blockHeaders, blockHeader)
	}

	return &message.ObjectMetadata{
		Id: &message.ObjectID{
			Id: item.ID.ToInt64(),
		},
		Name:         item.Name,
		Group:        item.Group,
		Partition:    item.Partition,
		Path:         item.Path,
		Size:         int32(item.Size),
		BlockHeaders: blockHeaders,
	}, nil
}

func (s *metadataRegistry) FindMetadataOnPath(c context.Context, req *rpc.ObjectMetadataRequest) (*rpc.ObjectMetadataList, error) {
	d := dto.Request{
		Group:     req.GetGroup(),
		Partition: req.GetPartition(),
		Path:      req.GetPath(),
	}

	items, err := s.objectMetadata.MetadataListOnPath(c, d)
	if err != nil {
		return nil, err
	}

	list := make([]*message.ObjectMetadata, len(items))
	for i, item := range items {
		blockHeaders := make([]*message.BlockHeader, 0)
		for _, header := range item.BlockHeaders {
			blockHeader := &message.BlockHeader{
				ObjectID: &message.ObjectID{
					Id: header.ObjectID.ToInt64(),
				},
				BlockID: &message.BlockID{
					Id: header.BlockID.ToInt64(),
				},
				Index:     int32(header.Index),
				Size:      int32(header.Size),
				Timestamp: timestamppb.New(header.Timestamp),
			}

			blockHeaders = append(blockHeaders, blockHeader)
		}

		list[i] = &message.ObjectMetadata{
			Id: &message.ObjectID{
				Id: item.ID.ToInt64(),
			},
			Name:         item.Name,
			Group:        item.Group,
			Partition:    item.Partition,
			Path:         item.Path,
			Size:         int32(item.Size),
			BlockHeaders: blockHeaders,
		}
	}

	return &rpc.ObjectMetadataList{
		Metadata: list,
	}, err
}
