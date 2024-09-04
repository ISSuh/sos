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
	"io"

	"github.com/ISSuh/sos/internal/domain/model/dto"
	"github.com/ISSuh/sos/internal/domain/model/entity"
	"github.com/ISSuh/sos/internal/domain/model/message"
	"github.com/ISSuh/sos/internal/domain/service/object"
	"github.com/ISSuh/sos/internal/infrastructure/transport/rpc"
	"github.com/ISSuh/sos/pkg/validation"
)

type Explorer interface {
	GetObjectMetadataByID(c context.Context, req dto.Request) (dto.Metadata, error)
	FindObjectMetadataOnPath(c context.Context, req dto.Request) (dto.MetadataList, error)
	Upload(c context.Context, req dto.Request, bodyStream io.ReadCloser) (dto.Metadata, error)
}

type explorer struct {
	metadataRequestor rpc.MetadataRegistryRequestor
	storageRequestor  rpc.BlockStorageRequestor
}

func NewExplorer(
	metadataRequestor rpc.MetadataRegistryRequestor, storageRequestor rpc.BlockStorageRequestor,
) (Explorer, error) {
	switch {
	case validation.IsNil(metadataRequestor):
		return nil, fmt.Errorf("MetadataRegistry requestor is nil")
	case validation.IsNil(storageRequestor):
		return nil, fmt.Errorf("BlockStorage requestor is nil")
	}

	return &explorer{
		metadataRequestor: metadataRequestor,
		storageRequestor:  storageRequestor,
	}, nil
}

func (s *explorer) GetObjectMetadataByID(c context.Context, req dto.Request) (dto.Metadata, error) {
	switch {
	case !req.ObjectID.IsValid():
		return dto.NewEmptyMetadata(), fmt.Errorf("object id is invalid")
	case validation.IsEmpty(req.Group):
		return dto.NewEmptyMetadata(), fmt.Errorf("group is empty")
	case validation.IsEmpty(req.Partition):
		return dto.NewEmptyMetadata(), fmt.Errorf("partition is empty")
	case validation.IsEmpty(req.Path):
		return dto.NewEmptyMetadata(), fmt.Errorf("path is empty")
	}

	message := rpc.ObjectMetadataRequest{
		Group:     req.Group,
		Partition: req.Partition,
		Path:      req.Path,
		Name:      req.Name,
	}

	metadata, err := s.metadataRequestor.GetByObjectName(c, &message)
	if err != nil {
		return dto.NewEmptyMetadata(), err
	}

	return dto.NewMetadataFromMessage(metadata), nil
}

func (s *explorer) FindObjectMetadataOnPath(c context.Context, req dto.Request) (dto.MetadataList, error) {
	switch {
	case validation.IsEmpty(req.Group):
		return nil, fmt.Errorf("group is empty")
	case validation.IsEmpty(req.Partition):
		return nil, fmt.Errorf("partition is empty")
	case validation.IsEmpty(req.Path):
		return nil, fmt.Errorf("path is empty")
	}

	message := rpc.ObjectMetadataRequest{
		Group:     req.Group,
		Partition: req.Partition,
		Path:      req.Path,
		Name:      req.Name,
	}

	resp, err := s.metadataRequestor.FindMetadataOnPath(c, &message)
	if err != nil {
		return nil, err
	}

	list := resp.GetMetadata()
	metadataList := make(dto.MetadataList, len(list))
	for i, item := range list {
		metadataList[i] = dto.NewMetadataFromMessage(item)
	}

	return metadataList, nil
}

func (s *explorer) Upload(c context.Context, req dto.Request, bodyStream io.ReadCloser) (dto.Metadata, error) {
	switch {
	case validation.IsEmpty(req.Group):
		return dto.NewEmptyMetadata(), fmt.Errorf("group is empty")
	case validation.IsEmpty(req.Partition):
		return dto.NewEmptyMetadata(), fmt.Errorf("partition is empty")
	case validation.IsEmpty(req.Path):
		return dto.NewEmptyMetadata(), fmt.Errorf("path is empty")
	case validation.IsEmpty(req.Name):
		return dto.NewEmptyMetadata(), fmt.Errorf("name is empty")
	case req.Size == 0:
		return dto.NewEmptyMetadata(), fmt.Errorf("size is 0")
	case validation.IsNil(bodyStream):
		return dto.NewEmptyMetadata(), fmt.Errorf("body stream is nil")
	case req.ChunkSize == 0:
		return dto.NewEmptyMetadata(), fmt.Errorf("chunk size is 0")
	}

	exist, err := s.isObjectExist(c, req)
	if err != nil {
		return dto.NewEmptyMetadata(), err
	}

	if exist {
		return dto.NewEmptyMetadata(), fmt.Errorf("object already exist")
	}

	objectID := entity.NewObjectID()
	uploader := object.NewUploader(s.storageRequestor)
	blockheaders, err := uploader.Upload(c, objectID, bodyStream)
	if err != nil {
		return dto.NewEmptyMetadata(), err
	}

	metadataBuilder := entity.NewObjectMetadataBuilder()
	metadataBuilder.
		ID(objectID).
		Group(req.Group).
		Partition(req.Partition).
		Name(req.Name).
		Path(req.Path).
		Size(req.Size).
		BlockHeaders(blockheaders)

	metadata := metadataBuilder.Build()
	if err := s.upsertObjectMetadata(c, metadata); err != nil {
		return dto.NewEmptyMetadata(), err
	}

	return dto.NewMetadataFromModel(&metadata), nil
}

func (s *explorer) getObjectMetadataByNameOnPath(c context.Context, req dto.Request) (entity.ObjectMetadata, error) {
	switch {
	case validation.IsEmpty(req.Name):
		return entity.NewEmptyObjectMetadata(), fmt.Errorf("name is empty")
	}

	message := rpc.ObjectMetadataRequest{
		Group:     req.Group,
		Partition: req.Partition,
		Path:      req.Path,
		Name:      req.Name,
	}

	resp, err := s.metadataRequestor.GetByObjectName(c, &message)
	if err != nil {
		return entity.NewEmptyObjectMetadata(), err
	}

	builder := entity.NewObjectMetadataBuilder()
	builder.ID(entity.NewObjectIDFrom(resp.Id.Id)).
		Group(resp.Group).
		Partition(resp.Partition).
		Name(resp.Name).
		Path(resp.Path).
		Size(int(resp.Size))

	return builder.Build(), nil
}

func (s *explorer) isObjectExist(c context.Context, req dto.Request) (bool, error) {
	metadata, err := s.getObjectMetadataByNameOnPath(c, req)
	if err != nil {
		return false, err
	}
	return !metadata.IsEmpty(), nil
}

func (s *explorer) upsertObjectMetadata(c context.Context, metadata entity.ObjectMetadata) error {
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

	if _, err := s.metadataRequestor.Put(c, message); err != nil {
		return err
	}
	return nil
}
