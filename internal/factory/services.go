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

	"github.com/ISSuh/sos/internal/domain/repository"
	"github.com/ISSuh/sos/internal/domain/service"
	"github.com/ISSuh/sos/internal/infrastructure/transport/rpc"
	"github.com/ISSuh/sos/pkg/log"
	"github.com/ISSuh/sos/pkg/validation"
)

type APIServices struct {
	Explorer   service.Explorer
	Uploader   service.Uploader
	Downloader service.Downloader
	Eraser     service.Eraser
}

func NewAPIServices(
	l log.Logger, metadataRequestor rpc.MetadataRegistryRequestor, storageRequestor rpc.BlockStorageRequestor,
) (*APIServices, error) {
	switch {
	case validation.IsNil(l):
		return nil, fmt.Errorf("logger is nil")
	case validation.IsNil(metadataRequestor):
		return nil, fmt.Errorf("MetadataRegistry requestor is nil")
	case validation.IsNil(storageRequestor):
		return nil, fmt.Errorf("BlockStorage requestor is nil")
	}

	explorer, err := service.NewExplorer(l, metadataRequestor)
	if err != nil {
		return nil, err
	}

	uploader, err := service.NewUploader(l, explorer, storageRequestor)
	if err != nil {
		return nil, err
	}

	downloader, err := service.NewDownloader(l, explorer, storageRequestor)
	if err != nil {
		return nil, err
	}

	eraser, err := service.NewEraser(l, explorer, storageRequestor)
	if err != nil {
		return nil, err
	}

	return &APIServices{
		Explorer:   explorer,
		Uploader:   uploader,
		Downloader: downloader,
		Eraser:     eraser,
	}, nil
}

func NewObjectMetadataService(l log.Logger, repo repository.ObjectMetadata) (service.ObjectMetadata, error) {
	switch {
	case validation.IsNil(l):
		return nil, fmt.Errorf("logger is nil")
	case validation.IsNil(repo):
		return nil, fmt.Errorf("ObjectMetadata repository is nil")
	}

	objectMetadata, err := service.NewObjectMetadata(l, repo)
	if err != nil {
		return nil, err
	}

	return objectMetadata, nil
}

func NewObjectStorageService(l log.Logger, repo repository.ObjectStorage) (service.ObjectStorage, error) {
	switch {
	case validation.IsNil(l):
		return nil, fmt.Errorf("logger is nil")
	case validation.IsNil(repo):
		return nil, fmt.Errorf("ObjectStorage repository is nil")
	}

	objectStorage, err := service.NewObjectStorage(l, repo)
	if err != nil {
		return nil, err
	}

	return objectStorage, nil
}
