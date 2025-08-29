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

	"github.com/ISSuh/sos/domain/repository"
	"github.com/ISSuh/sos/domain/service"
	"github.com/ISSuh/sos/infrastructure/transport/rpc"
	"github.com/ISSuh/sos/internal/validation"
)

func NewExplorerService(metadataRequestor rpc.MetadataRegistryRequestor, storageRequestor rpc.BlockStorageRequestor,
) (service.Explorer, error) {
	switch {
	case validation.IsNil(metadataRequestor):
		return nil, fmt.Errorf("MetadataRegistry requestor is nil")
	case validation.IsNil(storageRequestor):
		return nil, fmt.Errorf("BlockStorage requestor is nil")
	}

	explorer, err := service.NewExplorer(metadataRequestor, storageRequestor)
	if err != nil {
		return nil, err
	}

	return explorer, nil
}

func NewObjectMetadataService(repo repository.ObjectMetadata) (service.ObjectMetadata, error) {
	switch {
	case validation.IsNil(repo):
		return nil, fmt.Errorf("ObjectMetadata repository is nil")
	}

	objectMetadata, err := service.NewObjectMetadata(repo)
	if err != nil {
		return nil, err
	}

	return objectMetadata, nil
}

func NewObjectStorageService(repo repository.ObjectStorage) (service.ObjectStorage, error) {
	switch {
	case validation.IsNil(repo):
		return nil, fmt.Errorf("ObjectStorage repository is nil")
	}

	objectStorage, err := service.NewObjectStorage(repo)
	if err != nil {
		return nil, err
	}

	return objectStorage, nil
}
