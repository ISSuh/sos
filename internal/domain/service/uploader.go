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
	"time"

	"github.com/ISSuh/sos/internal/domain/model/dto"
	"github.com/ISSuh/sos/internal/domain/model/entity"
	"github.com/ISSuh/sos/internal/infrastructure/transport/rpc"
	"github.com/ISSuh/sos/pkg/log"
	"github.com/ISSuh/sos/pkg/validation"
)

type Uploader interface {
	Upload(c context.Context, req dto.Request, bodyStream io.ReadCloser) error
}

type uploader struct {
	logger log.Logger

	explorerService Explorer

	storageRequestor rpc.BlockStorageRequestor
}

func NewUploader(
	l log.Logger, explorerService Explorer, storageRequestor rpc.BlockStorageRequestor,
) (Uploader, error) {
	switch {
	case validation.IsNil(l):
		return nil, fmt.Errorf("logger is nil")
	case validation.IsNil(explorerService):
		return nil, fmt.Errorf("find service is nil")
	case validation.IsNil(storageRequestor):
		return nil, fmt.Errorf("BlockStorage requestor is nil")
	}

	return &uploader{
		logger:           l,
		explorerService:  explorerService,
		storageRequestor: storageRequestor,
	}, nil
}

func (s *uploader) Upload(c context.Context, req dto.Request, bodyStream io.ReadCloser) error {
	exist, err := s.explorerService.IsObjectExist(c, req)
	if err != nil {
		return err
	}

	if exist {
		return fmt.Errorf("object already exist")
	}

	objectID := entity.NewObjectID()
	metadataBuilder := entity.NewObjectMetadataBuilder()
	metadataBuilder.
		ID(objectID).
		Group(req.Group).
		Partition(req.Partition).
		Name(req.Name).
		Path(req.Path).
		Size(req.Size)

	var totalReadSize uint64
	var blockheaders entity.BlockHeaders
	blockBuilder := entity.NewBlockBuilder()
	blockIndex := 0
	for {
		buf := make([]byte, 4096)
		n, err := bodyStream.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}

		if n == 0 {
			block := s.buildBlock(objectID, blockIndex, blockBuilder)
			blockheaders = append(blockheaders, block.Header())
			break
		}

		// need data handling
		totalReadSize += uint64(n)
		if totalReadSize >= entity.BlockSize {
			block := s.buildBlock(objectID, blockIndex, blockBuilder)

			blockheaders = append(blockheaders, block.Header())

			blockBuilder = entity.NewBlockBuilder()
			blockIndex++
			totalReadSize = 0
		}

		blockBuilder.AppendBuffer(buf)
	}

	metadata := metadataBuilder.Build()
	object := entity.NewObject(metadata, blockheaders)
	if err := s.explorerService.UpsertObjectMetadata(c, object); err != nil {
		return err
	}

	return nil
}

func (s *uploader) buildBlock(objectID entity.ObjectID, index int, blockBuilder *entity.BlockBuilder) entity.Block {
	c := blockBuilder.CalculateChecksum()

	blockerHeaderBuilder := entity.NewBlockHeaderBuilder(objectID)
	blockerHeaderBuilder.
		BlockID(entity.NewBlockID()).
		Index(index).
		Size(blockBuilder.BufferSize()).
		Node(entity.Node{}).
		Timestamp(time.Now()).
		Checksum(c)

	blockHeader := blockerHeaderBuilder.Build()

	blockBuilder.Header(blockHeader)
	return blockBuilder.Build()
}
