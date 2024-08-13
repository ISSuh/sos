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
	"github.com/ISSuh/sos/pkg/log"
	"github.com/ISSuh/sos/pkg/validation"
)

type Finder interface {
	FindObjectMetadata(c context.Context, req dto.Request) (entity.ObjectMetadata, error)
	IsObjectExist(c context.Context, req dto.Request) (bool, error)
	CanUploadNewObject(c context.Context, req dto.Request) (bool, uint64, error)
}

type finder struct {
	logger log.Logger

	metadataService ObjectMetadata
}

func NewFinder(
	l log.Logger, metadataService ObjectMetadata,
) (Finder, error) {
	switch {
	case validation.IsNil(l):
		return nil, fmt.Errorf("logger is nil")
	case validation.IsNil(metadataService):
		return nil, fmt.Errorf("ObjectMetadata is nil")
	}

	return &finder{
		logger:          l,
		metadataService: metadataService,
	}, nil
}

func (s *finder) FindObjectMetadata(c context.Context, req dto.Request) (entity.ObjectMetadata, error) {
	metadata, err := s.metadataService.MetadataByObjectName(c, req)
	if err != nil {
		return entity.NewEmptyObjectMetadata(), err
	}
	return metadata, nil
}

func (s *finder) IsObjectExist(c context.Context, req dto.Request) (bool, error) {
	metadata, err := s.FindObjectMetadata(c, req)
	if err != nil {
		return false, err
	}
	return metadata.IsEmpty(), nil
}

func (s *finder) CanUploadNewObject(c context.Context, req dto.Request) (bool, uint64, error) {
	exist, err := s.IsObjectExist(c, req)
	if err != nil {
		return false, 0, err
	}

	if exist {
		return false, 0, nil
	}

	id, err := s.metadataService.GenerateNewObjectID(c)
	if err != nil {
		return false, 0, err
	}

	return true, id, nil
}
