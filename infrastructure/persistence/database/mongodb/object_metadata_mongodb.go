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

	"github.com/ISSuh/sos/domain/model/entity"
	"github.com/ISSuh/sos/domain/repository"

	"github.com/ISSuh/sos/internal/log"
	"github.com/ISSuh/sos/internal/mongodb"
)

type mongoDBObjectMetadata struct {
	db *mongodb.DB
}

func NewMongoDBObjectMetadata(db *mongodb.DB) (repository.ObjectMetadata, error) {
	return &mongoDBObjectMetadata{
		db: db,
	}, nil
}

func (d *mongoDBObjectMetadata) Create(c context.Context, metadata *entity.ObjectMetadata) error {
	log.FromContext(c).Debugf("[mongoDBObjectMetadata.Create] metadata: %+v", metadata)

	_, err := d.db.Collection(entity.ObjectMetadataCollectionName).
		InsertOne(c, metadata)
	if err != nil {
		return err
	}
	return nil
}

func (d *mongoDBObjectMetadata) Update(c context.Context, metadata *entity.ObjectMetadata) error {
	log.FromContext(c).Debugf("[mongoDBObjectMetadata.Update] metadata: %+v", metadata)

	return nil
}

func (d *mongoDBObjectMetadata) Delete(c context.Context, metadata *entity.ObjectMetadata) error {
	log.FromContext(c).Debugf("[mongoDBObjectMetadata.Delete] metadata: %+v", metadata)
	return nil
}

func (d *mongoDBObjectMetadata) MetadataByObjectName(c context.Context, group, partition, path, name string) (*entity.ObjectMetadata, error) {
	log.FromContext(c).Debugf("[mongoDBObjectMetadata.MetadataByObjectID] group: %s, partition: %s, path: %s, name: %s", group, partition, path, name)
	return nil, nil
}

func (d *mongoDBObjectMetadata) MetadataByObjectID(c context.Context, group, partition, path string, objectID int64) (*entity.ObjectMetadata, error) {
	log.FromContext(c).Debugf("[mongoDBObjectMetadata.MetadataByObjectID] group: %s, partition: %s, path: %s, objectID: %d", group, partition, path, objectID)

	return nil, nil
}

func (d *mongoDBObjectMetadata) FindMetadata(c context.Context, group, partition, path string) (entity.ObjectMetadataList, error) {
	log.FromContext(c).Debugf("[mongoDBObjectMetadata.FindMetadata] group: %s, partition: %s, path: %s", group, partition, path)
	return nil, nil
}
