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
	"github.com/ISSuh/sos/internal/infrastructure/transport/rpc"
	"github.com/ISSuh/sos/pkg/log"
	"github.com/ISSuh/sos/pkg/validation"
)

type Uploader interface {
	Upload(c context.Context, req dto.Request, bodyStream io.ReadCloser) error
}

type uploader struct {
	logger log.Logger

	findService Finder

	metadataRequestor rpc.MetadataRegistryRequestor
	storageRequestor  rpc.BlockStorageRequestor
}

func NewUploader(
	l log.Logger, findService Finder, metadataRequestor rpc.MetadataRegistryRequestor, storageRequestor rpc.BlockStorageRequestor,
) (Uploader, error) {
	switch {
	case validation.IsNil(l):
		return nil, fmt.Errorf("logger is nil")
	case validation.IsNil(findService):
		return nil, fmt.Errorf("find service is nil")
	case validation.IsNil(metadataRequestor):
		return nil, fmt.Errorf("MetadataRegistry requestor is nil")
	case validation.IsNil(storageRequestor):
		return nil, fmt.Errorf("BlockStorage requestor is nil")
	}

	return &uploader{
		logger:            l,
		findService:       findService,
		metadataRequestor: metadataRequestor,
		storageRequestor:  storageRequestor,
	}, nil
}

func (s *uploader) Upload(c context.Context, req dto.Request, bodyStream io.ReadCloser) error {
	defer bodyStream.Close()

	exist, id, err := s.findService.CanUploadNewObject(c, req)
	if err != nil {
		return err
	}

	if exist {
		return fmt.Errorf("object already exist")
	}

	totalReadSize := uint64(0)
	for {
		buf := make([]byte, 4096)
		n, err := bodyStream.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}

		if n == 0 {
			break
		}

		// need data handling
		totalReadSize += uint64(n)
		s.logger.Debugf("[uploader.chunkedUpload] Read %d bytes\n", n)
	}

	if err := s.uploadNewObjectMetadata(c, req, id); err != nil {
		return err
	}

	return nil
}

func (s *uploader) uploadNewObjectMetadata(c context.Context, req dto.Request, id uint64) error {
	// req.ID = id
	// if err := s.metadataService.Create(c, req); err != nil {
	// 	return err
	// }
	return nil
}
