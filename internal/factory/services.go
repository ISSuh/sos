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

package factory

import (
	"fmt"

	"github.com/ISSuh/sos/internal/domain/service"
	"github.com/ISSuh/sos/pkg/log"
	"github.com/ISSuh/sos/pkg/validation"
)

type APIServices struct {
	Finder     service.Finder
	Uploader   service.Uploader
	Downloader service.Downloader
	Eraser     service.Eraser
}

type MetadataService struct {
	ObjectMetadata service.ObjectMetadata
}

type StorageService struct {
	ObjectStorage service.ObjectStorage
}

func NewAPIServices(
	l log.Logger, objectMetadata service.ObjectMetadata, objectStorage service.ObjectStorage,
) (*APIServices, error) {
	switch {
	case validation.IsNil(l):
		return nil, fmt.Errorf("logger is nil")
	case validation.IsNil(objectMetadata):
		return nil, fmt.Errorf("object metadata service is nil")
	case validation.IsNil(objectStorage):
		return nil, fmt.Errorf("object storage service is nil")
	}

	finder, err := service.NewFinder(l, objectMetadata)
	if err != nil {
		return nil, err
	}

	uploader, err := service.NewUploader(l, finder, objectMetadata, objectStorage)
	if err != nil {
		return nil, err
	}

	downloader, err := service.NewDownloader(l, finder, objectStorage)
	if err != nil {
		return nil, err
	}

	eraser, err := service.NewEraser(l, finder, objectStorage)
	if err != nil {
		return nil, err
	}

	return &APIServices{
		Finder:     finder,
		Uploader:   uploader,
		Downloader: downloader,
		Eraser:     eraser,
	}, nil
}

func NewMetadataService(l log.Logger, repositoryFactory *Repositories) (*MetadataService, error) {
	switch {
	case validation.IsNil(l):
		return nil, fmt.Errorf("logger is nil")
	case validation.IsNil(repositoryFactory):
		return nil, fmt.Errorf("repository factory is nil")
	}

	objectMetadata, err := service.NewObjectMetadata(l, repositoryFactory.ObjectMetadata)
	if err != nil {
		return nil, err
	}

	return &MetadataService{
		ObjectMetadata: objectMetadata,
	}, nil
}

func NewStorageService(l log.Logger, repositoryFactory *Repositories) (*StorageService, error) {
	switch {
	case validation.IsNil(l):
		return nil, fmt.Errorf("logger is nil")
	case validation.IsNil(repositoryFactory):
		return nil, fmt.Errorf("repository factory is nil")
	}

	objectStorage, err := service.NewObjectStorage(l, repositoryFactory.ObjectStorage)
	if err != nil {
		return nil, err
	}

	return &StorageService{
		ObjectStorage: objectStorage,
	}, nil
}
