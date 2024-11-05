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
	Create(c context.Context, dto dto.Metadata) error
	Update(c context.Context, dto dto.Metadata) error
	Delete(c context.Context, dto dto.Metadata) error
	MetadataByObjectName(c context.Context, group, partition, path, objectName string) (dto.Metadata, error)
	MetadataByObjectID(c context.Context, group, partition, path string, objectID int64) (dto.Metadata, error)
	MetadataListOnPath(c context.Context, group, partition, path string) (dto.MetadataList, error)
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

func (s *objectMetadata) Create(c context.Context, dto dto.Metadata) error {
	log.FromContext(c).Debugf("[objectMetadata.Create] request: %+v", dto)

	version := entity.Versions{
		entity.NewVersionBuilder().
			Number(0).
			Size(dto.Size).
			BlockHeaders(dto.BlockHeaders.ToEntity()).
			Build(),
	}

	metadata := dto.ToEntityWithVersion(version)
	if err := s.metadataRepository.Create(c, metadata); err != nil {
		return err
	}

	return nil
}

func (s *objectMetadata) Update(c context.Context, dto dto.Metadata) error {
	log.FromContext(c).Debugf("[objectMetadata.Delete] Update: %+v", dto)

	// metadata := dto.ToEntity()
	// object, err := s.metadataRepository.MetadataByObjectID(c, metadata.Group(), metadata.Partition(), metadata.Path(), metadata.ID().ToInt64())
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (s *objectMetadata) Delete(c context.Context, dto dto.Metadata) error {
	log.FromContext(c).Debugf("[objectMetadata.Delete] request: %+v", dto)

	metadata := dto.ToEntity()
	if err := s.metadataRepository.Delete(c, metadata); err != nil {
		return err
	}

	return nil
}

func (s *objectMetadata) MetadataByObjectName(
	c context.Context, group, partition, path, objectName string,
) (dto.Metadata, error) {
	metadata, err :=
		s.metadataRepository.MetadataByObjectName(c, group, partition, path, objectName)
	if err != nil {
		return dto.NewEmptyMetadata(), err
	}
	return dto.NewMetadataFromModel(metadata), nil
}

func (s *objectMetadata) MetadataByObjectID(
	c context.Context, group, partition, path string, objectID int64,
) (dto.Metadata, error) {
	metadata, err :=
		s.metadataRepository.MetadataByObjectID(c, group, partition, path, objectID)
	if err != nil {
		return dto.NewEmptyMetadata(), err
	}
	return dto.NewMetadataFromModel(metadata), nil
}

func (s *objectMetadata) MetadataListOnPath(c context.Context, group, partition, path string) (dto.MetadataList, error) {
	items, err :=
		s.metadataRepository.FindMetadata(c, group, partition, path)
	if err != nil {
		return nil, err
	}

	list := make(dto.MetadataList, len(items))
	for i, item := range items {
		list[i] = dto.NewMetadataFromModel(item)
	}
	return list, nil
}
