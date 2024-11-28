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
	"errors"
	"fmt"
	"time"

	"github.com/ISSuh/sos/domain/model/dto"
	"github.com/ISSuh/sos/domain/model/entity"
	"github.com/ISSuh/sos/domain/repository"
	soserror "github.com/ISSuh/sos/internal/error"
	"github.com/ISSuh/sos/internal/log"
	"github.com/ISSuh/sos/internal/validation"
)

type ObjectMetadata interface {
	Put(c context.Context, objectDTO *dto.Object) (*dto.Metadata, error)
	Delete(c context.Context, metadataDTO *dto.Metadata) error
	MetadataByObjectName(c context.Context, group, partition, path, objectName string) (*dto.Metadata, error)
	MetadataByObjectID(c context.Context, group, partition, path string, objectID int64) (*dto.Metadata, error)
	MetadataListOnPath(c context.Context, group, partition, path string) (dto.MetadataList, error)
}

type objectMetadata struct {
	metadataRepository  repository.ObjectMetadata
	directoryRepository repository.ObjectDirectory
	tempID              uint64
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

func (s *objectMetadata) Put(c context.Context, objectDTO *dto.Object) (*dto.Metadata, error) {
	log.FromContext(c).Debugf("[objectMetadata.Put] request: %+v", objectDTO)
	metadata, err :=
		s.metadataRepository.MetadataByObjectID(
			c, objectDTO.Group, objectDTO.Partition, objectDTO.Path, objectDTO.ID.ToInt64(),
		)
	if err != nil && !errors.Is(err, soserror.NotFound) {
		return nil, err
	}

	now := time.Now()
	if metadata == nil {
		metadata, err = s.createMetadata(c, objectDTO, now)
		if err != nil {
			return nil, err
		}
	} else {
		metadata.ModifiedAt = now

		metadata, err = s.updateMetadata(c, metadata, objectDTO, now)
		if err != nil {
			return nil, err
		}
	}

	resp := dto.NewMetadataFromModel(metadata)
	return resp, nil
}

func (s *objectMetadata) Delete(c context.Context, metadataDTO *dto.Metadata) error {
	log.FromContext(c).Debugf("[objectMetadata.Delete] request: %+v", metadataDTO)
	if metadataDTO.Versions.Empty() {
		deletedMetadata := metadataDTO.ToEntity()
		if err := s.metadataRepository.Delete(c, &deletedMetadata); err != nil {
			return err
		}
		return nil
	}

	metadata, err :=
		s.metadataRepository.MetadataByObjectID(
			c, metadataDTO.Group, metadataDTO.Partition, metadataDTO.Path, metadataDTO.ID.ToInt64(),
		)
	if err != nil {
		return err
	}

	for _, version := range metadataDTO.Versions {
		if err := metadata.DeleteVersion(version.Number); err != nil {
			return err
		}
	}

	if metadata.Versions().Empty() {
		if err := s.metadataRepository.Delete(c, metadata); err != nil {
			return err
		}
	} else {
		if err := s.metadataRepository.Update(c, metadata); err != nil {
			return err
		}
	}

	return nil
}

func (s *objectMetadata) MetadataByObjectName(
	c context.Context, group, partition, path, objectName string,
) (*dto.Metadata, error) {
	metadata, err :=
		s.metadataRepository.MetadataByObjectName(c, group, partition, path, objectName)
	if err != nil {
		return nil, err
	}

	if metadata == nil {
		metadata = &entity.ObjectMetadata{}
	}

	return dto.NewMetadataFromModel(metadata), nil
}

func (s *objectMetadata) MetadataByObjectID(
	c context.Context, group, partition, path string, objectID int64,
) (*dto.Metadata, error) {
	metadata, err :=
		s.metadataRepository.MetadataByObjectID(c, group, partition, path, objectID)
	if err != nil {
		return nil, err
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
		list[i] = *dto.NewMetadataFromModel(&item)
	}
	return list, nil
}

func (s *objectMetadata) createMetadata(
	c context.Context, object *dto.Object, now time.Time,
) (*entity.ObjectMetadata, error) {
	metadata := object.ToEntity()
	metadata.CreatedAt = now
	metadata.ModifiedAt = now

	versionNumber := 0
	blockHeaders := object.BlockHeaders.ToEntity()
	version := s.newVersion(versionNumber, object.Size, blockHeaders, now)
	metadata.AppendVersion(version)

	if err := s.metadataRepository.Create(c, &metadata); err != nil {
		return nil, err
	}

	return &metadata, nil
}

func (s *objectMetadata) updateMetadata(
	c context.Context, metadata *entity.ObjectMetadata, object *dto.Object, now time.Time,
) (*entity.ObjectMetadata, error) {
	versionNumber := metadata.LastVersion() + 1
	blockHeaders := object.BlockHeaders.ToEntity()
	version := s.newVersion(versionNumber, object.Size, blockHeaders, now)
	metadata.AppendVersion(version)

	if err := s.metadataRepository.Update(c, metadata); err != nil {
		return nil, err
	}

	return metadata, nil
}

func (s *objectMetadata) newVersion(
	versionNum int, size int, blockHeaders entity.BlockHeaders, now time.Time,
) entity.Version {
	versionNumber := versionNum
	version := entity.NewVersionBuilder().
		Number(versionNumber).
		Size(size).
		BlockHeaders(blockHeaders).
		Build()

	version.CreatedAt = now
	version.ModifiedAt = now
	return version
}
