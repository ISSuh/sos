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

package service

import (
	"context"
	"fmt"

	"github.com/ISSuh/sos/internal/domain/model/dto"
	"github.com/ISSuh/sos/internal/domain/model/entity"
	"github.com/ISSuh/sos/internal/domain/model/message"
	"github.com/ISSuh/sos/internal/infrastructure/transport/rpc"
	"github.com/ISSuh/sos/pkg/log"
	"github.com/ISSuh/sos/pkg/validation"
)

type Explorer interface {
	FindObjectMetadata(c context.Context, req dto.Request) (dto.Metadata, error)
	IsObjectExist(c context.Context, req dto.Request) (bool, error)
	UpsertObjectMetadata(c context.Context, object entity.Object) error
}

type explorer struct {
	logger log.Logger

	metadataRequestor rpc.MetadataRegistryRequestor
}

func NewExplorer(
	l log.Logger, metadataRequestor rpc.MetadataRegistryRequestor,
) (Explorer, error) {
	switch {
	case validation.IsNil(l):
		return nil, fmt.Errorf("logger is nil")
	case validation.IsNil(metadataRequestor):
		return nil, fmt.Errorf("MetadataRegistry requestor is nil")
	}

	return &explorer{
		logger:            l,
		metadataRequestor: metadataRequestor,
	}, nil
}

func (s *explorer) FindObjectMetadata(c context.Context, req dto.Request) (dto.Metadata, error) {
	message := &message.MetadataFindRequest{
		Group:     req.Group,
		Partition: req.Partition,
		Path:      req.Path,
		Name:      req.Name,
	}

	metadata, err := s.metadataRequestor.GetByObjectName(c, message)
	if err != nil {
		return dto.NewEmptyMetadata(), err
	}

	return dto.NewMetadataFromMessage(metadata), nil
}

func (s *explorer) IsObjectExist(c context.Context, req dto.Request) (bool, error) {
	metadata, err := s.FindObjectMetadata(c, req)
	if err != nil {
		return false, err
	}
	return !metadata.IsEmpty(), nil
}

func (s *explorer) UpsertObjectMetadata(c context.Context, object entity.Object) error {
	metadata := object.Metadata()
	message := &message.ObjectMetadata{
		Id: &message.ObjectID{
			Id: metadata.ID().ToInt64(),
		},
		Group:     metadata.Group(),
		Partition: metadata.Partition(),
		Path:      metadata.Path(),
		Name:      metadata.Name(),
		Size:      int32(metadata.Size()),
	}

	if _, err := s.metadataRequestor.Create(c, message); err != nil {
		return err
	}
	return nil
}
