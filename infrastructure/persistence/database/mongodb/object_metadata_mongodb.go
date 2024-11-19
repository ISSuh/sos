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

	"github.com/ISSuh/sos/domain/model/entity"
	"github.com/ISSuh/sos/domain/repository"
	soserror "github.com/ISSuh/sos/internal/error"
	"github.com/ISSuh/sos/internal/log"
	"github.com/ISSuh/sos/internal/persistence"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	objectMetadataCollectionName = "object_metadata"
)

type mongoDBObjectMetadata struct {
	db *persistence.MongoDB
}

func NewMongoDBObjectMetadata(db *persistence.MongoDB) (repository.ObjectMetadata, error) {
	return &mongoDBObjectMetadata{
		db: db,
	}, nil
}

func (d *mongoDBObjectMetadata) Create(c context.Context, metadata *entity.ObjectMetadata) error {
	log.FromContext(c).Debugf("[mongoDBObjectMetadata.Create] metadata: %+v", metadata)
	switch {
	case c == nil:
		return fmt.Errorf("context is nil")
	case metadata.Group() == "":
		return fmt.Errorf("group is invalid")
	case metadata.Partition() == "":
		return fmt.Errorf("partition is empty")
	case metadata.Path() == "":
		return fmt.Errorf("path is empty")
	case !metadata.ID().IsValid():
		return fmt.Errorf("objectID is invalid. %d", metadata.ID())
	}

	collection, err := d.db.Collection(objectMetadataCollectionName)
	if err != nil {
		return err
	}

	res, err := collection.InsertOne(c, metadata)
	if err != nil {
		return fmt.Errorf("failed to insert data: %w", err)
	}

	if res.InsertedID == nil {
		return fmt.Errorf("invakid insertedID")
	}

	return nil
}

func (d *mongoDBObjectMetadata) Update(c context.Context, metadata *entity.ObjectMetadata) error {
	log.FromContext(c).Debugf("[mongoDBObjectMetadata.Update] metadata: %+v", metadata)
	switch {
	case c == nil:
		return fmt.Errorf("context is nil")
	case metadata.Group() == "":
		return fmt.Errorf("group is invalid")
	case metadata.Partition() == "":
		return fmt.Errorf("partition is empty")
	case metadata.Path() == "":
		return fmt.Errorf("path is empty")
	case !metadata.ID().IsValid():
		return fmt.Errorf("objectID is invalid. %d", metadata.ID())
	}

	collection, err := d.db.Collection(objectMetadataCollectionName)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "group", Value: metadata.Group()},
		{Key: "partition", Value: metadata.Partition()},
		{Key: "path", Value: metadata.Path()},
		{Key: "object_id", Value: metadata.ID()},
	}

	res, err := collection.ReplaceOne(c, filter, metadata)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return soserror.NewNotFoundError(fmt.Errorf("can not find metadata"))
		}
		return fmt.Errorf("failed to update data: %w", err)
	}

	if res.ModifiedCount == 0 {
		return fmt.Errorf("can not find metadata")
	}

	return nil
}

func (d *mongoDBObjectMetadata) Delete(c context.Context, metadata *entity.ObjectMetadata) error {
	log.FromContext(c).Debugf("[mongoDBObjectMetadata.Delete] metadata: %+v", metadata)
	switch {
	case c == nil:
		return fmt.Errorf("context is nil")
	case metadata.Group() == "":
		return fmt.Errorf("group is invalid")
	case metadata.Partition() == "":
		return fmt.Errorf("partition is empty")
	case metadata.Path() == "":
		return fmt.Errorf("path is empty")
	case !metadata.ID().IsValid():
		return fmt.Errorf("objectID is invalid. %d", metadata.ID())
	}

	collection, err := d.db.Collection(objectMetadataCollectionName)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "group", Value: metadata.Group()},
		{Key: "partition", Value: metadata.Partition()},
		{Key: "path", Value: metadata.Path()},
		{Key: "object_id", Value: metadata.ID()},
	}

	res, err := collection.DeleteOne(c, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return soserror.NewNotFoundError(fmt.Errorf("can not find metadata"))
		}
		return fmt.Errorf("failed to delete data: %w", err)
	}

	if res.DeletedCount == 0 {
		return fmt.Errorf("can not find metadata")
	}

	return nil
}

func (d *mongoDBObjectMetadata) MetadataByObjectName(c context.Context, group, partition, path, name string) (*entity.ObjectMetadata, error) {
	log.FromContext(c).Debugf("[mongoDBObjectMetadata.MetadataByObjectID] group: %s, partition: %s, path: %s, name: %s", group, partition, path, name)
	switch {
	case c == nil:
		return nil, fmt.Errorf("context is nil")
	case group == "":
		return nil, fmt.Errorf("group is invalid")
	case partition == "":
		return nil, fmt.Errorf("partition is empty")
	case path == "":
		return nil, fmt.Errorf("path is empty")
	}

	list, err := d.FindMetadata(c, group, partition, path)
	if err != nil {
		return nil, err
	}

	for _, metadata := range list {
		if metadata.Name() == name {
			return &metadata, nil
		}
	}
	return nil, soserror.NewNotFoundError(fmt.Errorf("can not find metadata"))
}

func (d *mongoDBObjectMetadata) MetadataByObjectID(c context.Context, group, partition, path string, objectID int64) (*entity.ObjectMetadata, error) {
	log.FromContext(c).Debugf("[mongoDBObjectMetadata.MetadataByObjectID] group: %s, partition: %s, path: %s, objectID: %d", group, partition, path, objectID)
	switch {
	case c == nil:
		return nil, fmt.Errorf("context is nil")
	case group == "":
		return nil, fmt.Errorf("group is invalid")
	case partition == "":
		return nil, fmt.Errorf("partition is empty")
	case path == "":
		return nil, fmt.Errorf("path is empty")
	case objectID <= 0:
		return nil, fmt.Errorf("objectID is invalid. %d", objectID)
	}

	collection, err := d.db.Collection(objectMetadataCollectionName)
	if err != nil {
		return nil, err
	}

	filter := bson.D{
		{Key: "group", Value: group},
		{Key: "partition", Value: partition},
		{Key: "path", Value: path},
		{Key: "object_id", Value: objectID},
	}

	res := collection.FindOne(c, filter)
	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			return nil, soserror.NewNotFoundError(fmt.Errorf("can not find metadata"))
		}
		return nil, fmt.Errorf("failed to find metadata: %w", res.Err())
	}

	var metadata entity.ObjectMetadata
	if err := res.Decode(&metadata); err != nil {
		return nil, fmt.Errorf("failed to decode metadata: %w", err)
	}

	return &metadata, nil
}

func (d *mongoDBObjectMetadata) FindMetadata(c context.Context, group, partition, path string) (entity.ObjectMetadataList, error) {
	log.FromContext(c).Debugf("[mongoDBObjectMetadata.FindMetadata] group: %s, partition: %s, path: %s", group, partition, path)
	switch {
	case c == nil:
		return nil, fmt.Errorf("context is nil")
	case group == "":
		return nil, fmt.Errorf("group is invalid")
	case partition == "":
		return nil, fmt.Errorf("partition is empty")
	case path == "":
		return nil, fmt.Errorf("path is empty")
	}

	collection, err := d.db.Collection(objectMetadataCollectionName)
	if err != nil {
		return nil, err
	}

	filter := bson.D{
		{Key: "group", Value: group},
		{Key: "partition", Value: partition},
		{Key: "path", Value: path},
	}

	res, err := collection.Find(c, filter)
	if err != nil {
		if res.Err() == mongo.ErrNoDocuments {
			return nil, soserror.NewNotFoundError(fmt.Errorf("can not find metadata"))
		}
		return nil, fmt.Errorf("failed to find metadata: %w", err)
	}

	var metadataList entity.ObjectMetadataList
	if err := res.All(c, &metadataList); err != nil {
		return nil, fmt.Errorf("failed to decode metadata: %w", err)
	}

	return metadataList, nil
}
