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

	"github.com/ISSuh/sos/domain/model/entity"
	"github.com/ISSuh/sos/domain/repository"
	"github.com/ISSuh/sos/internal/empty"
	"github.com/ISSuh/sos/internal/log"
)

type localObjectStorage struct {
	storage map[string]entity.Block
}

func NewLocalObjectStorage() (repository.ObjectStorage, error) {
	return &localObjectStorage{
			storage: make(map[string]entity.Block),
		},
		nil
}

func (s *localObjectStorage) Put(c context.Context, block entity.Block) error {
	log.FromContext(c).Debugf("[localObjectStorage.Put] block header: %+v", block.Header())
	header := block.Header()
	key := s.makeKey(header.ObjectID(), header.BlockID(), header.Index())
	s.storage[key] = block
	return nil
}

func (s *localObjectStorage) GetBlock(c context.Context, objectID entity.ObjectID, blockID entity.BlockID, index int) (entity.Block, error) {
	log.FromContext(c).Debugf("[localObjectStorage.GetBlock] objectID: %s, blockID : %d, index: %d", objectID, blockID, index)
	key := s.makeKey(objectID, blockID, index)
	block, exist := s.storage[key]
	if !exist {
		return empty.Struct[entity.Block](), fmt.Errorf("block not found")
	}
	return block, nil
}

func (s *localObjectStorage) GetBlockHeader(c context.Context, objectID entity.ObjectID, blockID entity.BlockID, index int) (entity.BlockHeader, error) {
	log.FromContext(c).Debugf("[localObjectStorage.GetBlockHeader] objectID: %s, blockID : %d, index: %d", objectID, blockID, index)
	key := s.makeKey(objectID, blockID, index)
	block, exist := s.storage[key]
	if !exist {
		return empty.Struct[entity.BlockHeader](), fmt.Errorf("block not found")
	}
	return block.Header(), nil
}

func (s *localObjectStorage) Delete(c context.Context, objectID entity.ObjectID, blockID entity.BlockID, index int) error {
	log.FromContext(c).Debugf("[localObjectStorage.Delete] objectID: %s, blockID : %d, index: %d", objectID, blockID, index)
	key := s.makeKey(objectID, blockID, index)
	delete(s.storage, key)
	return nil
}

func (s *localObjectStorage) makeKey(objectID entity.ObjectID, blockID entity.BlockID, index int) string {
	return objectID.String() + ":" + blockID.String() + ":" + strconv.Itoa(index)
}
