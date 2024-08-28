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

package objectstorage

import (
	"context"
	"fmt"
	"strconv"

	"github.com/ISSuh/sos/internal/domain/model/entity"
	"github.com/ISSuh/sos/internal/domain/repository"
	"github.com/ISSuh/sos/pkg/log"
)

type localObjectStorage struct {
	logger log.Logger

	storage map[string]entity.Block
}

func NewLocalObjectStorage(l log.Logger) (repository.ObjectStorage, error) {
	return &localObjectStorage{
			logger:  l,
			storage: make(map[string]entity.Block),
		},
		nil
}

func (s *localObjectStorage) Put(c context.Context, block entity.Block) error {
	log.FromContext(c).Debugf("[localObjectStorage.Put] block: %+v", block)
	header := block.Header()
	key := s.makeKey(header.ObjectID(), header.ID(), header.Index())
	s.storage[key] = block
	return nil
}

func (s *localObjectStorage) GetBlock(c context.Context, objectId string, blockID, index uint64) (entity.Block, error) {
	log.FromContext(c).Debugf("[localObjectStorage.GetBlock] objectID: %s, blockID : %d, index: %d", objectId, blockID, index)
	key := s.makeKey(objectId, blockID, index)
	block, exist := s.storage[key]
	if !exist {
		return entity.Empty[entity.Block](), fmt.Errorf("block not found")
	}
	return block, nil
}

func (s *localObjectStorage) GetBlockHeader(c context.Context, objectId string, blockID, index uint64) (entity.BlockHeader, error) {
	log.FromContext(c).Debugf("[localObjectStorage.GetBlockHeader] objectID: %s, blockID : %d, index: %d", objectId, blockID, index)
	key := s.makeKey(objectId, blockID, index)
	block, exist := s.storage[key]
	if !exist {
		return entity.Empty[entity.BlockHeader](), fmt.Errorf("block not found")
	}
	return block.Header(), nil
}

func (s *localObjectStorage) Delete(c context.Context, objectId string, blockID, index uint64) error {
	log.FromContext(c).Debugf("[localObjectStorage.Delete] objectID: %s, blockID : %d, index: %d", objectId, blockID, index)
	key := s.makeKey(objectId, blockID, index)
	delete(s.storage, key)
	return nil
}

func (s *localObjectStorage) makeKey(objectId string, blockID, index uint64) string {
	return objectId + ":" + strconv.FormatUint(blockID, 10) + ":" + strconv.FormatUint(index, 10)
}
