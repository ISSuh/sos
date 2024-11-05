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
	"github.com/ISSuh/sos/internal/domain/service/object"
	"github.com/ISSuh/sos/internal/infrastructure/transport/rpc"
	"github.com/ISSuh/sos/pkg/empty"
	"github.com/ISSuh/sos/pkg/http"
	"github.com/ISSuh/sos/pkg/validation"
)

type Explorer interface {
	GetObjectMetadata(c context.Context, req dto.Request) (dto.Item, error)
	FindObjectMetadataOnPath(c context.Context, req dto.Request) (dto.Items, error)
	Upload(c context.Context, req dto.Request, bodyStream io.ReadCloser) (dto.Item, error)
	Download(c context.Context, req dto.Request, headerWriter http.DownloadHeaderWriter, bodyWriter http.DownloadBodyWriter) error
	Delete(c context.Context, req dto.Request) error
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

func (s *explorer) GetObjectMetadata(c context.Context, req dto.Request) (dto.Item, error) {
	switch {
	case !req.ObjectID.IsValid():
		return empty.Struct[dto.Item](), fmt.Errorf("object id is invalid")
	case validation.IsEmpty(req.Group):
		return empty.Struct[dto.Item](), fmt.Errorf("group is empty")
	case validation.IsEmpty(req.Partition):
		return empty.Struct[dto.Item](), fmt.Errorf("partition is empty")
	case validation.IsEmpty(req.Path):
		return empty.Struct[dto.Item](), fmt.Errorf("path is empty")
	}

	metadata, err :=
		s.getObjectMetadataByObjectID(c, req.ObjectID, req.Group, req.Partition, req.Path)
	if err != nil {
		return empty.Struct[dto.Item](), err
	}

	return dto.NewItemFromMetadata(metadata), nil
}

func (s *explorer) FindObjectMetadataOnPath(c context.Context, req dto.Request) (dto.Items, error) {
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
	}

	resp, err := s.metadataRequestor.FindMetadataOnPath(c, &message)
	if err != nil {
		return nil, err
	}

	list := resp.GetMetadata()
	return dto.NewItemFromMetadataListMessage(list), nil
}

func (s *explorer) Upload(c context.Context, req dto.Request, bodyStream io.ReadCloser) (dto.Item, error) {
	switch {
	case validation.IsEmpty(req.Group):
		return empty.Struct[dto.Item](), fmt.Errorf("group is empty")
	case validation.IsEmpty(req.Partition):
		return empty.Struct[dto.Item](), fmt.Errorf("partition is empty")
	case validation.IsEmpty(req.Path):
		return empty.Struct[dto.Item](), fmt.Errorf("path is empty")
	case validation.IsEmpty(req.Name):
		return empty.Struct[dto.Item](), fmt.Errorf("name is empty")
	case req.Size <= 0:
		return empty.Struct[dto.Item](), fmt.Errorf("size is invalid")
	case validation.IsNil(bodyStream):
		return empty.Struct[dto.Item](), fmt.Errorf("body stream is nil")
	}

	exist, err :=
		s.isObjectNameExist(c, req.Group, req.Partition, req.Path, req.Name)
	if err != nil {
		return empty.Struct[dto.Item](), err
	}

	if exist {
		return empty.Struct[dto.Item](), fmt.Errorf("object already exist")
	}

	objectID := entity.NewObjectID()
	uploader := object.NewUploader(s.storageRequestor)
	blockheaders, err := uploader.Upload(c, objectID, bodyStream)
	if err != nil {
		return empty.Struct[dto.Item](), err
	}

	metadata := dto.Metadata{
		ID:           objectID,
		Group:        req.Group,
		Partition:    req.Partition,
		Name:         req.Name,
		Path:         req.Path,
		Size:         req.Size,
		BlockHeaders: blockheaders,
	}

	if err := s.upsertObjectMetadata(c, metadata); err != nil {
		return empty.Struct[dto.Item](), err
	}

	return dto.NewItemFromMetadata(metadata), nil
}

func (s *explorer) Download(
	c context.Context, req dto.Request, headerWriter http.DownloadHeaderWriter, bodyWriter http.DownloadBodyWriter,
) error {
	switch {
	case validation.IsEmpty(req.Group):
		return fmt.Errorf("group is empty")
	case validation.IsEmpty(req.Partition):
		return fmt.Errorf("partition is empty")
	case validation.IsEmpty(req.Path):
		return fmt.Errorf("path is empty")
	case !req.ObjectID.IsValid():
		return fmt.Errorf("object id is invalid")
	}

	metadata, err :=
		s.getObjectMetadataByObjectID(c, req.ObjectID, req.Group, req.Partition, req.Path)
	if err != nil {
		return err
	}

	headerWriter(metadata.Name, metadata.Size)

	downloader := object.NewDownloader(s.storageRequestor)
	err = downloader.Download(c, metadata, bodyWriter)
	if err != nil {
		return err
	}

	return nil
}

func (s *explorer) Delete(c context.Context, req dto.Request) error {
	switch {
	case !req.ObjectID.IsValid():
		return fmt.Errorf("object id is invalid")
	case validation.IsEmpty(req.Group):
		return fmt.Errorf("group is empty")
	case validation.IsEmpty(req.Partition):
		return fmt.Errorf("partition is empty")
	case validation.IsEmpty(req.Path):
		return fmt.Errorf("path is empty")
	}

	metadata, err :=
		s.getObjectMetadataByObjectID(c, req.ObjectID, req.Group, req.Partition, req.Path)
	if err != nil {
		return err
	}

	if !metadata.ID.IsValid() {
		return fmt.Errorf("object not exist")
	}

	deleter := object.NewDeleter(s.metadataRequestor, s.storageRequestor)
	if err := deleter.Delete(c, metadata); err != nil {
		return err
	}

	return nil
}

func (s *explorer) getObjectMetadataByObjectID(
	c context.Context, objectID entity.ObjectID, group, partition, path string,
) (dto.Metadata, error) {
	message := rpc.ObjectMetadataRequest{
		ObjectID:  objectID.ToInt64(),
		Group:     group,
		Partition: partition,
		Path:      path,
	}

	resp, err := s.metadataRequestor.GetByObjectID(c, &message)
	if err != nil {
		return empty.Struct[dto.Metadata](), err
	}

	//  nned not found
	if resp.Id.Id <= 0 {
		return empty.Struct[dto.Metadata](), nil
	}

	return dto.NewMetadataFromMessage(resp), nil
}

func (s *explorer) getObjectMetadataByNameOnPath(
	c context.Context, group, partition, path, name string,
) (dto.Metadata, error) {
	switch {
	case validation.IsEmpty(name):
		return empty.Struct[dto.Metadata](), fmt.Errorf("name is empty")
	}

	message := rpc.ObjectMetadataRequest{
		Group:     group,
		Partition: partition,
		Path:      path,
		Name:      name,
	}

	resp, err := s.metadataRequestor.GetByObjectName(c, &message)
	if err != nil {
		return empty.Struct[dto.Metadata](), err
	}
	return dto.NewMetadataFromMessage(resp), nil
}

func (s *explorer) isObjectNameExist(c context.Context, group, partition, path, name string) (bool, error) {
	metadata, err := s.getObjectMetadataByNameOnPath(c, group, partition, path, name)
	if err != nil {
		return false, err
	}
	return metadata.ID.IsValid(), nil
}

func (s *explorer) upsertObjectMetadata(c context.Context, metadata dto.Metadata) error {
	message := metadata.ToMessage()
	if _, err := s.metadataRequestor.Put(c, message); err != nil {
		return err
	}
	return nil
}
