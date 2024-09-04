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
	"fmt"

	"github.com/ISSuh/sos/internal/domain/model/message"
	"github.com/ISSuh/sos/internal/domain/service"
	"github.com/ISSuh/sos/internal/infrastructure/transport/rpc"
	"github.com/ISSuh/sos/pkg/log"
	"github.com/ISSuh/sos/pkg/validation"
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

func (h *metadataRegistry) Put(c context.Context, metadata *message.ObjectMetadata) (*message.ObjectMetadata, error) {
	log.FromContext(c).Debugf("[MetadataRegistry.Put]")
	return &message.ObjectMetadata{}, nil
}

func (h *metadataRegistry) Delete(c context.Context, metadata *message.ObjectMetadata) (bool, error) {
	log.FromContext(c).Debugf("[MetadataRegistry.Delete]")
	return true, nil
}

func (h *metadataRegistry) GetByObjectName(c context.Context, rew *rpc.ObjectMetadataRequest) (*message.ObjectMetadata, error) {
	log.FromContext(c).Debugf("[MetadataRegistry.GetByObjectName]")
	return &message.ObjectMetadata{
		Id: &message.ObjectID{
			Id: 1,
		},
	}, nil
}

func (h *metadataRegistry) FindMetadataOnPath(c context.Context, req *rpc.ObjectMetadataRequest) (*rpc.ObjectMetadataList, error) {
	log.FromContext(c).Debugf("[MetadataRegistry.FindMetadataOnPath]")
	return &rpc.ObjectMetadataList{}, nil
}
