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

	"github.com/ISSuh/sos/internal/domain/model/entity"
	"github.com/ISSuh/sos/internal/domain/repository"
	"github.com/ISSuh/sos/pkg/log"
	"github.com/ISSuh/sos/pkg/validation"
)

type ObjectStorage interface {
	Put(c context.Context, block entity.Block) error

	GetBlock(
		c context.Context, objectID entity.ObjectID, blockID entity.BlockID, index uint64,
	) (entity.Block, error)

	GetBlockHeader(
		c context.Context, objectID entity.ObjectID, blockID entity.BlockID, index uint64,
	) (entity.BlockHeader, error)

	Delete(c context.Context, objectID entity.ObjectID, blockID entity.BlockID, index uint64) error
}

type objectStorage struct {
	logger log.Logger

	storageRepository repository.ObjectStorage
}

func NewObjectStorage(
	l log.Logger, storageRepository repository.ObjectStorage,
) (ObjectStorage, error) {
	switch {
	case validation.IsNil(l):
		return nil, fmt.Errorf("logger is nil")
	case validation.IsNil(storageRepository):
		return nil, fmt.Errorf("StorageRepository is nil")
	}

	return &objectStorage{
		logger:            l,
		storageRepository: storageRepository,
	}, nil
}

func (s *objectStorage) Put(c context.Context, block entity.Block) error {
	if err := s.storageRepository.Put(c, block); err != nil {
		return err
	}
	return nil
}

func (s *objectStorage) GetBlock(c context.Context, objectID entity.ObjectID, blockID entity.BlockID, index uint64) (entity.Block, error) {
	block, err := s.storageRepository.GetBlock(c, objectID, blockID, index)
	if err != nil {
		return entity.Block{}, err
	}
	return block, nil
}

func (s *objectStorage) GetBlockHeader(c context.Context, objectID entity.ObjectID, blockID entity.BlockID, index uint64) (entity.BlockHeader, error) {
	header, err := s.storageRepository.GetBlockHeader(c, objectID, blockID, index)
	if err != nil {
		return entity.BlockHeader{}, err
	}
	return header, nil
}

func (s *objectStorage) Delete(c context.Context, objectID entity.ObjectID, blockID entity.BlockID, index uint64) error {
	if err := s.storageRepository.Delete(c, objectID, blockID, index); err != nil {
		return err
	}
	return nil
}
