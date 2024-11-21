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

package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/ISSuh/sos/domain/model/message"
	"github.com/ISSuh/sos/domain/service"
	"github.com/ISSuh/sos/infrastructure/transport/rpc"
	rpcmessage "github.com/ISSuh/sos/infrastructure/transport/rpc/message"
	soserror "github.com/ISSuh/sos/internal/error"
	"github.com/ISSuh/sos/internal/log"
	"github.com/ISSuh/sos/internal/validation"

	"google.golang.org/grpc/status"
)

type metadataRegistry struct {
	objectMetadata service.ObjectMetadata
}

func NewMetadataRegistry(objectMetadata service.ObjectMetadata) (rpc.MetadataRegistryHandler, error) {
	switch {
	case validation.IsNil(objectMetadata):
		return nil, fmt.Errorf("ObjectMetadata service is nil")
	}

	return &metadataRegistry{
		objectMetadata: objectMetadata,
	}, nil
}

func (h *metadataRegistry) Put(c context.Context, msg *message.Object) (*message.ObjectMetadata, error) {
	log.FromContext(c).Debugf("[MetadataRegistry.Put]")
	switch {
	case validation.IsNil(c):
		return nil, fmt.Errorf("Context is nil")
	case validation.IsNil(msg):
		return nil, fmt.Errorf("Object is nil")
	}

	object := message.ToObjectDTO(msg)
	metadata, err := h.objectMetadata.Put(c, object)
	if err != nil {
		return nil, err
	}

	return message.FromObjectMetadataDTO(metadata), nil
}

func (h *metadataRegistry) Delete(c context.Context, msg *message.ObjectMetadata) error {
	log.FromContext(c).Debugf("[MetadataRegistry.Delete]")
	switch {
	case validation.IsNil(c):
		return fmt.Errorf("Context is nil")
	case validation.IsNil(msg):
		return fmt.Errorf("ObjectMetadata is nil")
	}

	metadata := message.ToObjectMetadataDTO(msg)
	err := h.objectMetadata.Delete(c, metadata)
	if err != nil {
		return err
	}

	return nil
}

func (h *metadataRegistry) GetByObjectName(c context.Context, msg *rpcmessage.ObjectMetadataRequest) (*message.ObjectMetadata, error) {
	log.FromContext(c).Debugf("[MetadataRegistry.GetByObjectName]")
	switch {
	case validation.IsNil(c):
		return nil, fmt.Errorf("Context is nil")
	case validation.IsNil(msg):
		return nil, fmt.Errorf("ObjectMetadataRequest is nil")
	case validation.IsEmpty(msg.Group):
		return nil, fmt.Errorf("Group is empty")
	case validation.IsEmpty(msg.Partition):
		return nil, fmt.Errorf("Partition is empty")
	case validation.IsEmpty(msg.Path):
		return nil, fmt.Errorf("Path is empty")
	case validation.IsEmpty(msg.Name):
		return nil, fmt.Errorf("Name is empty")
	}

	metadata, err := h.objectMetadata.MetadataByObjectName(c, msg.Group, msg.Partition, msg.Name, msg.Path)
	if err != nil {
		if errors.Is(err, soserror.NotFound) {
			return nil, status.Errorf(soserror.NotFoundErrorCode, "%v", err)
		}
		return nil, err
	}

	return message.FromObjectMetadataDTO(metadata), nil
}

func (h *metadataRegistry) GetByObjectID(c context.Context, msg *rpcmessage.ObjectMetadataRequest) (*message.ObjectMetadata, error) {
	log.FromContext(c).Debugf("[MetadataRegistry.GetByObjectID]")
	switch {
	case validation.IsNil(c):
		return nil, fmt.Errorf("Context is nil")
	case validation.IsNil(msg):
		return nil, fmt.Errorf("ObjectMetadataRequest is nil")
	case validation.IsEmpty(msg.Group):
		return nil, fmt.Errorf("Group is empty")
	case validation.IsEmpty(msg.Partition):
		return nil, fmt.Errorf("Partition is empty")
	case validation.IsEmpty(msg.Path):
		return nil, fmt.Errorf("Path is empty")
	case msg.ObjectID <= 0:
		return nil, fmt.Errorf("ObjectID is invalid")
	}

	metadata, err := h.objectMetadata.MetadataByObjectID(c, msg.Group, msg.Partition, msg.Path, msg.ObjectID)
	if err != nil {
		if errors.Is(err, soserror.NotFound) {
			return nil, status.Errorf(soserror.NotFoundErrorCode, "%v", err)
		}
		return nil, err
	}

	return message.FromObjectMetadataDTO(metadata), nil
}

func (h *metadataRegistry) FindMetadataOnPath(c context.Context, msg *rpcmessage.ObjectMetadataRequest) (*message.ObjectMetadataList, error) {
	log.FromContext(c).Debugf("[MetadataRegistry.FindMetadataOnPath]")
	switch {
	case validation.IsNil(c):
		return nil, fmt.Errorf("Context is nil")
	case validation.IsNil(msg):
		return nil, fmt.Errorf("ObjectMetadataRequest is nil")
	case validation.IsEmpty(msg.Group):
		return nil, fmt.Errorf("Group is empty")
	case validation.IsEmpty(msg.Partition):
		return nil, fmt.Errorf("Partition is empty")
	}

	list, err := h.objectMetadata.MetadataListOnPath(c, msg.Group, msg.Partition, msg.Path)
	if err != nil {
		if errors.Is(err, soserror.NotFound) {
			return nil, status.Errorf(soserror.NotFoundErrorCode, "%v", err)
		}
		return nil, err
	}

	return message.FromObjectMetadataListDTO(list), nil
}
