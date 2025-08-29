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
	"bytes"
	"context"
	"encoding/gob"
	"fmt"

	"github.com/ISSuh/sos/domain/model/entity"
	"github.com/ISSuh/sos/domain/repository"
	"github.com/ISSuh/sos/internal/log"
	"github.com/ISSuh/sos/internal/persistence"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

type LevelDBObjectStorage struct {
	storage *persistence.LevelDB
}

func NewLevelDBObjectStorage(storage *persistence.LevelDB) (repository.ObjectStorage, error) {
	return &LevelDBObjectStorage{
		storage: storage,
	}, nil
}

func (s *LevelDBObjectStorage) Put(c context.Context, block *entity.Block) error {
	log.FromContext(c).Debugf("[LevelDBObjectStorage.Put] block header: %+v", block.Header())
	switch {
	case c == nil:
		return fmt.Errorf("context is nil")
	case block == nil:
		return fmt.Errorf("block is nil")
	}

	if err := block.Validate(); err != nil {
		return err
	}

	storage, err := s.storage.Engin()
	if err != nil {
		return err
	}

	key := s.makeKey(block.ObjectID(), block.BlockID(), block.Index())
	data, err := s.encodeBlock(block)
	if err != nil {
		return err
	}

	if err := storage.Put(key, data, &opt.WriteOptions{Sync: true}); err != nil {
		return err
	}

	return nil
}

func (s *LevelDBObjectStorage) GetBlock(c context.Context, objectID entity.ObjectID, blockID entity.BlockID, index int) (*entity.Block, error) {
	log.FromContext(c).Debugf("[LevelDBObjectStorage.GetBlock] objectID: %s, blockID : %d, index: %d", objectID, blockID, index)
	switch {
	case c == nil:
		return nil, fmt.Errorf("context is nil")
	case !objectID.IsValid():
		return nil, fmt.Errorf("objectID is invalid")
	case blockID < 0:
		return nil, fmt.Errorf("blockID is invalid")
	case index < 0:
		return nil, fmt.Errorf("index is invalid")
	}

	storage, err := s.storage.Engin()
	if err != nil {
		return nil, err
	}

	key := s.makeKey(objectID, blockID, index)
	data, err := storage.Get(key, nil)
	if err != nil {
		return nil, err
	}

	block, err := s.decodeBlock(data)
	if err != nil {
		return nil, err
	}

	return block, nil
}

func (s *LevelDBObjectStorage) GetBlockHeader(c context.Context, objectID entity.ObjectID, blockID entity.BlockID, index int) (*entity.BlockHeader, error) {
	log.FromContext(c).Debugf("[LevelDBObjectStorage.GetBlockHeader] objectID: %s, blockID : %d, index: %d", objectID, blockID, index)
	switch {
	case c == nil:
		return nil, fmt.Errorf("context is nil")
	case !objectID.IsValid():
		return nil, fmt.Errorf("objectID is invalid")
	case blockID < 0:
		return nil, fmt.Errorf("blockID is invalid")
	case index < 0:
		return nil, fmt.Errorf("index is invalid")
	}

	block, err := s.GetBlock(c, objectID, blockID, index)
	if err != nil {
		return nil, err
	}

	header := block.Header()
	return &header, nil
}

func (s *LevelDBObjectStorage) Delete(c context.Context, objectID entity.ObjectID, blockID entity.BlockID, index int) error {
	log.FromContext(c).Debugf("[LevelDBObjectStorage.Delete] objectID: %s, blockID : %d, index: %d", objectID, blockID, index)
	switch {
	case c == nil:
		return fmt.Errorf("context is nil")
	case !objectID.IsValid():
		return fmt.Errorf("objectID is invalid")
	case blockID < 0:
		return fmt.Errorf("blockID is invalid")
	case index < 0:
		return fmt.Errorf("index is invalid")
	}

	storage, err := s.storage.Engin()
	if err != nil {
		return err
	}

	key := s.makeKey(objectID, blockID, index)
	if err := storage.Delete(key, &opt.WriteOptions{Sync: true}); err != nil {
		return err
	}

	return nil
}

func (s *LevelDBObjectStorage) makeKey(objectID entity.ObjectID, blockID entity.BlockID, index int) []byte {
	key := fmt.Sprintf("%s:%d:%d", objectID, blockID, index)
	return []byte(key)
}

func (s *LevelDBObjectStorage) encodeBlock(block *entity.Block) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	if err := encoder.Encode(block); err != nil {
		return nil, err
	}

	test, _ := s.decodeBlock(buffer.Bytes())
	fmt.Println(test.Header())

	return buffer.Bytes(), nil
}

func (s *LevelDBObjectStorage) decodeBlock(data []byte) (*entity.Block, error) {
	var block entity.Block
	reader := bytes.NewReader(data)
	decoder := gob.NewDecoder(reader)

	if err := decoder.Decode(&block); err != nil {
		return nil, err
	}

	return &block, nil
}
