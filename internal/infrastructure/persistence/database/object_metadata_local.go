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

	db map[string]map[string]entity.ObjectMetadata
}

func NewLocalObjectMetadata(l log.Logger) (repository.ObjectMetadata, error) {
	return &localObjectMetadata{
			logger: l,
			db:     make(map[string]map[string]entity.ObjectMetadata),
		},
		nil
}

func (d *localObjectMetadata) Create(c context.Context, metadata entity.ObjectMetadata) error {
	log.FromContext(c).Debugf("[localObjectMetadata.Create] metadata: %+v", metadata)
	key := d.makeKey(metadata.Group(), metadata.Partition(), metadata.Path())
	_, exist := d.db[key]
	if !exist {
		d.db[key] = make(map[string]entity.ObjectMetadata)
	}

	d.db[key][metadata.Name()] = metadata
	return nil
}

func (d *localObjectMetadata) Update(c context.Context, metadata entity.ObjectMetadata) error {
	log.FromContext(c).Debugf("[localObjectMetadata.Update] metadata: %+v", metadata)
	key := d.makeKey(metadata.Group(), metadata.Partition(), metadata.Path())
	_, exist := d.db[key]
	if !exist {
		return fmt.Errorf("metadata not exist")
	}

	d.db[key][metadata.Name()] = metadata
	return nil
}

func (d *localObjectMetadata) Delete(c context.Context, metadata entity.ObjectMetadata) error {
	log.FromContext(c).Debugf("[localObjectMetadata.Delete] metadata: %+v", metadata)
	key := d.makeKey(metadata.Group(), metadata.Partition(), metadata.Path())
	_, exist := d.db[key]
	if !exist {
		return fmt.Errorf("metadata not exist")
	}

	delete(d.db[key], metadata.Name())
	return nil
}

func (d *localObjectMetadata) MetadataByObjectName(c context.Context, group, partition, path, name string) (entity.ObjectMetadata, error) {
	log.FromContext(c).Debugf("[localObjectMetadata.MetadataByObjectName] group: %s, partition: %s, path: %s, name: %s", group, partition, path, name)
	key := d.makeKey(group, partition, path)
	_, exist := d.db[key]
	if !exist {
		return entity.NewEmptyObjectMetadata(), nil
	}
	return d.db[key][name], nil
}

func (d *localObjectMetadata) FindMetadata(c context.Context, group, partition, path string) ([]entity.ObjectMetadata, error) {
	log.FromContext(c).Debugf("[localObjectMetadata.FindMetadata] group: %s, partition: %s, path: %s", group, partition, path)
	key := d.makeKey(group, partition, path)
	list := d.db[key]

	var metadataList []entity.ObjectMetadata
	for _, v := range list {
		metadataList = append(metadataList, v)
	}
	return metadataList, nil
}

func (d *localObjectMetadata) makeKey(group, partition, path string) string {
	return fmt.Sprintf("%s:%s:%s", group, partition, path)
}
