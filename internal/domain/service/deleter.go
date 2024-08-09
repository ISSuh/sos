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
	"fmt"

	"github.com/ISSuh/sos/internal/domain/repository"
	"github.com/ISSuh/sos/pkg/logger"
	"github.com/ISSuh/sos/pkg/validation"
)

type Deleter interface {
}

type deleter struct {
	logger logger.Logger

	metadataRepository repository.ObjectMetadata
	storageRepository  repository.ObjectStorage
}

func NewDeleter(
	l logger.Logger, metadataRepository repository.ObjectMetadata, storageRepository repository.ObjectStorage,
) (Deleter, error) {
	switch {
	case validation.IsNil(l):
		return nil, fmt.Errorf("logger is nil")
	case validation.IsNil(metadataRepository):
		return nil, fmt.Errorf("MetadataRepository is nil")
	case validation.IsNil(storageRepository):
		return nil, fmt.Errorf("StorageRepository is nil")
	}

	return &deleter{
		logger:             l,
		metadataRepository: metadataRepository,
		storageRepository:  storageRepository,
	}, nil
}
