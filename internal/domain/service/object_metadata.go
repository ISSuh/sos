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
	"github.com/ISSuh/sos/internal/domain/repository"
	"github.com/ISSuh/sos/pkg/log"
	"github.com/ISSuh/sos/pkg/validation"
)

type ObjectMetadata interface {
	Create(c context.Context, req dto.Request) error
	MetadataByObjectName(c context.Context, req dto.Request) (dto.Metadata, error)
	MetadataListOnPath(c context.Context, req dto.Request) (dto.MetadataList, error)
}

type objectMetadata struct {
	metadataRepository repository.ObjectMetadata
	tempID             uint64
}

func NewObjectMetadata(metadataRepository repository.ObjectMetadata) (ObjectMetadata, error) {
	switch {
	case validation.IsNil(metadataRepository):
		return nil, fmt.Errorf("MetadataRepository is nil")
	}

	return &objectMetadata{
		metadataRepository: metadataRepository,
		tempID:             0,
	}, nil
}

func (s *objectMetadata) Create(c context.Context, req dto.Request) error {
	log.FromContext(c).Debugf("[objectMetadata.Create] request: %+v", req)

	builder := entity.NewObjectMetadataBuilder()
	metadata :=
		builder.ID(req.ObjectID).
			Group(req.Group).
			Partition(req.Partition).
			Path(req.Path).
			Name(req.Name).
			Size(req.Size).
			Build()

	if err := s.metadataRepository.Create(c, metadata); err != nil {
		return err
	}

	return nil
}

func (s *objectMetadata) MetadataByObjectName(c context.Context, req dto.Request) (dto.Metadata, error) {
	metadata, err := s.metadataRepository.MetadataByObjectName(c, req.Group, req.Partition, req.Path, req.Name)
	if err != nil {
		return dto.NewEmptyMetadata(), err
	}
	return dto.NewMetadataFromModel(&metadata), nil
}

func (s *objectMetadata) MetadataListOnPath(c context.Context, req dto.Request) (dto.MetadataList, error) {
	return nil, nil
}
