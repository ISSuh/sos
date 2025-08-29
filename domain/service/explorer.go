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
	"io"

	"github.com/ISSuh/sos/domain/model/dto"
	"github.com/ISSuh/sos/domain/model/entity"
	"github.com/ISSuh/sos/domain/model/message"
	"github.com/ISSuh/sos/domain/service/object"
	"github.com/ISSuh/sos/infrastructure/transport/rpc"
	rpcmessage "github.com/ISSuh/sos/infrastructure/transport/rpc/message"
	"github.com/ISSuh/sos/internal/empty"
	soserror "github.com/ISSuh/sos/internal/error"
	"github.com/ISSuh/sos/internal/http"
	"github.com/ISSuh/sos/internal/validation"
)

type Explorer interface {
	GetObjectMetadata(c context.Context, req dto.Request) (dto.Item, error)
	FindObjectMetadataOnPath(c context.Context, req dto.Request) (dto.Items, error)
	Upload(c context.Context, req dto.Request, bodyStream io.ReadCloser) (dto.Item, error)
	Download(
		c context.Context, req dto.Request, writer http.Writer, lastVersion bool,
	) error
	Delete(c context.Context, req dto.Request, deleteVersion bool) error
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
		return nil, errors.New("MetadataRegistry requestor is nil")
	case validation.IsNil(storageRequestor):
		return nil, errors.New("BlockStorage requestor is nil")
	}

	return &explorer{
		metadataRequestor: metadataRequestor,
		storageRequestor:  storageRequestor,
	}, nil
}

func (s *explorer) GetObjectMetadata(c context.Context, req dto.Request) (dto.Item, error) {
	switch {
	case !req.ObjectID.IsValid():
		return empty.Struct[dto.Item](), errors.New("object id is invalid")
	case validation.IsEmpty(req.Group):
		return empty.Struct[dto.Item](), errors.New("group is empty")
	case validation.IsEmpty(req.Partition):
		return empty.Struct[dto.Item](), errors.New("partition is empty")
	case validation.IsEmpty(req.Path):
		return empty.Struct[dto.Item](), errors.New("path is empty")
	}

	metadata, err :=
		s.getObjectMetadataByObjectID(c, req.ObjectID, req.Group, req.Partition, req.Path)
	if err != nil {
		return empty.Struct[dto.Item](), err
	}

	return dto.NewItemFromMetadata(*metadata), nil
}

func (s *explorer) FindObjectMetadataOnPath(c context.Context, req dto.Request) (dto.Items, error) {
	switch {
	case validation.IsEmpty(req.Group):
		return nil, errors.New("group is empty")
	case validation.IsEmpty(req.Partition):
		return nil, errors.New("partition is empty")
	case validation.IsEmpty(req.Path):
		return nil, errors.New("path is empty")
	}

	msg := rpcmessage.ObjectMetadataRequest{
		Group:     req.Group,
		Partition: req.Partition,
		Path:      req.Path,
	}

	resp, err := s.metadataRequestor.FindMetadataOnPath(c, &msg)
	if err != nil {
		return nil, err
	}

	return message.ToItemsDTO(resp), nil
}

func (s *explorer) Upload(c context.Context, req dto.Request, bodyStream io.ReadCloser) (dto.Item, error) {
	switch {
	case validation.IsEmpty(req.Group):
		return empty.Struct[dto.Item](), errors.New("group is empty")
	case validation.IsEmpty(req.Partition):
		return empty.Struct[dto.Item](), errors.New("partition is empty")
	case validation.IsEmpty(req.Path):
		return empty.Struct[dto.Item](), errors.New("path is empty")
	case validation.IsEmpty(req.Name):
		return empty.Struct[dto.Item](), errors.New("name is empty")
	case req.Size <= 0:
		return empty.Struct[dto.Item](), errors.New("size is invalid")
	case validation.IsNil(bodyStream):
		return empty.Struct[dto.Item](), errors.New("body stream is nil")
	}

	metadata, err := s.getObjectMetadataByNameOnPath(c, req.Group, req.Partition, req.Path, req.Name)
	if err != nil && !errors.Is(err, soserror.NotFound) {
		return empty.Struct[dto.Item](), err
	}

	objectID := entity.NewObjectID()
	if metadata != nil && metadata.ID.IsValid() {
		objectID = metadata.ID
	}

	uploader := object.NewUploader(s.storageRequestor)
	blockheaders, err := uploader.Upload(c, objectID, bodyStream)
	if err != nil {
		return empty.Struct[dto.Item](), err
	}

	newObject := &dto.Object{
		ID:           objectID,
		Group:        req.Group,
		Partition:    req.Partition,
		Name:         req.Name,
		Path:         req.Path,
		Size:         req.Size,
		BlockHeaders: blockheaders,
	}

	resp, err := s.upsertObjectMetadata(c, newObject)
	if err != nil {
		return empty.Struct[dto.Item](), err
	}

	return dto.NewItemFromMetadata(*resp), nil
}

func (s *explorer) Download(c context.Context, req dto.Request, writer http.Writer, lastVersion bool) error {
	switch {
	case validation.IsEmpty(req.Group):
		return errors.New("group is empty")
	case validation.IsEmpty(req.Partition):
		return errors.New("partition is empty")
	case validation.IsEmpty(req.Path):
		return errors.New("path is empty")
	case !req.ObjectID.IsValid():
		return errors.New("object id is invalid")
	case !lastVersion && req.Version < 0:
		return errors.New("version is invalid")
	}

	metadata, err :=
		s.getObjectMetadataByObjectID(c, req.ObjectID, req.Group, req.Partition, req.Path)
	if err != nil {
		return err
	}

	var version dto.Version
	if lastVersion {
		version, err = metadata.Versions.LastVersion()
		if err != nil {
			return err
		}
	} else {
		version, err = metadata.Versions.Version(req.Version)
		if err != nil {
			return err
		}
	}

	writer.Header(metadata.Name, version.Size)

	downloader := object.NewDownloader(s.storageRequestor)
	err = downloader.Download(c, version, writer.Body)
	if err != nil {
		return err
	}

	return nil
}

func (s *explorer) Delete(c context.Context, req dto.Request, deleteVersion bool) error {
	switch {
	case !req.ObjectID.IsValid():
		return errors.New("object id is invalid")
	case validation.IsEmpty(req.Group):
		return errors.New("group is empty")
	case validation.IsEmpty(req.Partition):
		return errors.New("partition is empty")
	case validation.IsEmpty(req.Path):
		return errors.New("path is empty")
	case deleteVersion && req.Version < 0:
		return errors.New("version is invalid")
	}

	metadata, err :=
		s.getObjectMetadataByObjectID(c, req.ObjectID, req.Group, req.Partition, req.Path)
	if err != nil {
		return err
	}

	if !metadata.ID.IsValid() {
		return errors.New("object not exist")
	}

	if deleteVersion {
		if !metadata.Versions.HasVersion(req.Version) {
			return errors.New("version not exist")
		}
	}

	deleter := object.NewDeleter(s.metadataRequestor, s.storageRequestor)
	if deleteVersion {
		if err := deleter.DeleteVersion(c, *metadata, req.Version); err != nil {
			return err
		}
	} else {
		if err := deleter.Delete(c, *metadata); err != nil {
			return err
		}
	}

	return nil
}

func (s *explorer) getObjectMetadataByObjectID(
	c context.Context, objectID entity.ObjectID, group, partition, path string,
) (*dto.Metadata, error) {
	msg := rpcmessage.ObjectMetadataRequest{
		ObjectID:  objectID.ToInt64(),
		Group:     group,
		Partition: partition,
		Path:      path,
	}

	resp, err := s.metadataRequestor.GetByObjectID(c, &msg)
	if err != nil {
		return nil, err
	}

	//  need not found
	if resp.Id.Id <= 0 {
		return nil, nil
	}

	return message.ToObjectMetadataDTO(resp), nil
}

func (s *explorer) getObjectMetadataByNameOnPath(
	c context.Context, group, partition, path, name string,
) (*dto.Metadata, error) {
	switch {
	case validation.IsEmpty(name):
		return nil, errors.New("name is empty")
	}

	msg := rpcmessage.ObjectMetadataRequest{
		Group:     group,
		Partition: partition,
		Path:      path,
		Name:      name,
	}

	resp, err := s.metadataRequestor.GetByObjectName(c, &msg)
	if err != nil {
		return nil, err
	}

	return message.ToObjectMetadataDTO(resp), nil
}

func (s *explorer) upsertObjectMetadata(c context.Context, object *dto.Object) (*dto.Metadata, error) {
	msg := message.FromObjectDTO(object)
	resp, err := s.metadataRequestor.Put(c, msg)
	if err != nil {
		return nil, err
	}

	return message.ToObjectMetadataDTO(resp), nil
}
