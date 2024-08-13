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

package database

import (
	"context"
	"fmt"

	"github.com/ISSuh/sos/internal/domain/model/entity"
	"github.com/ISSuh/sos/internal/domain/repository"
	"github.com/ISSuh/sos/pkg/log"
)

type localObjectMetadata struct {
	logger log.Logger

	db map[string]map[string]map[string]map[string]entity.ObjectMetadata
}

func NewLocalObjectMetadata(l log.Logger) (repository.ObjectMetadata, error) {
	return &localObjectMetadata{
			logger: l,
			db:     make(map[string]map[string]map[string]map[string]entity.ObjectMetadata),
		},
		nil
}

func (d *localObjectMetadata) Create(c context.Context, metadata entity.ObjectMetadata) error {
	d.db[metadata.Group][metadata.Partition][metadata.Path][metadata.Name] = metadata
	return nil
}

func (d *localObjectMetadata) Update(c context.Context, metadata entity.ObjectMetadata) error {
	d.db[metadata.Group][metadata.Partition][metadata.Path][metadata.Name] = metadata
	return nil
}

func (d *localObjectMetadata) Delete(c context.Context, metadata entity.ObjectMetadata) error {
	_, exist := d.db[metadata.Group][metadata.Partition][metadata.Path][metadata.Name]
	if !exist {
		return fmt.Errorf("metadata not exist")
	}
	delete(d.db[metadata.Group][metadata.Partition][metadata.Name], metadata.Path)
	return nil
}

func (d *localObjectMetadata) MetadataByObjectName(c context.Context, group, partition, path, name string) (entity.ObjectMetadata, error) {
	metadata, exist := d.db[group][partition][path][name]
	if !exist {
		return entity.NewEmptyObjectMetadata(), fmt.Errorf("metadata not exist")
	}
	return metadata, nil
}

func (d *localObjectMetadata) FindMetadata(c context.Context, group, partition, path string) ([]entity.ObjectMetadata, error) {
	list, exist := d.db[group][partition][path]
	if !exist {
		return nil, fmt.Errorf("metadata not exist")
	}

	var metadataList []entity.ObjectMetadata
	for _, v := range list {
		metadataList = append(metadataList, v)
	}
	return metadataList, nil
}
