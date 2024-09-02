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
	"github.com/ISSuh/sos/pkg/log"
	"github.com/ISSuh/sos/pkg/validation"
)

type metadataRegistry struct {
	logger log.Logger

	objectMetadata service.ObjectMetadata
}

func NewMetadataRegistry(l log.Logger, objectMetadata service.ObjectMetadata) (rpc.MetadataRegistryRequestor, error) {
	switch {
	case validation.IsNil(l):
		return nil, fmt.Errorf("logger is nil")
	}

	return &metadataRegistry{
		logger:         l,
		objectMetadata: objectMetadata,
	}, nil
}

func (s *metadataRegistry) Create(c context.Context, metadata *message.ObjectMetadata) (*message.ObjectMetadata, error) {
	d := dto.Request{
		ObjectID:  entity.NewObjectIDFrom(metadata.GetId().Id),
		Name:      metadata.GetName(),
		Group:     metadata.GetGroup(),
		Partition: metadata.GetPartition(),
		Path:      metadata.GetPath(),
		Size:      int(metadata.GetSize()),
	}

	err := s.objectMetadata.Create(c, d)
	if err != nil {
		return nil, err
	}

	return metadata, nil
}

func (s *metadataRegistry) GetByObjectName(c context.Context, req *message.MetadataFindRequest) (*message.ObjectMetadata, error) {
	d := dto.Request{
		Group:     req.GetGroup(),
		Partition: req.GetPartition(),
		Path:      req.GetPath(),
		Name:      req.GetName(),
	}

	metadata, err := s.objectMetadata.MetadataByObjectName(c, d)
	if err != nil {
		return nil, err
	}

	return &message.ObjectMetadata{
		Id: &message.ObjectID{
			Id: metadata.ID().ToInt64(),
		},
		Name:      metadata.Name(),
		Group:     metadata.Group(),
		Partition: metadata.Partition(),
		Path:      metadata.Path(),
		Size:      int32(metadata.Size()),
	}, nil
}
