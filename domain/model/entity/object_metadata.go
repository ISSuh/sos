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

package entity

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type ObjectMetadataList []ObjectMetadata

type ObjectMetadata struct {
	id        ObjectID `bson:"object_id"`
	group     string   `bson:"group"`
	partition string   `bson:"partition"`
	name      string   `bson:"name"`
	path      string   `bson:"path"`
	versions  Versions `bson:"versions"`

	ModifiedTime
}

func (e *ObjectMetadata) ID() ObjectID {
	return e.id
}

func (e *ObjectMetadata) Group() string {
	return e.group
}

func (e *ObjectMetadata) Partition() string {
	return e.partition
}

func (e *ObjectMetadata) Name() string {
	return e.name
}

func (e *ObjectMetadata) Path() string {
	return e.path
}

func (e *ObjectMetadata) Versions() Versions {
	return e.versions
}

func (e *ObjectMetadata) IsValid() bool {
	return e.id.IsValid()
}

func (e *ObjectMetadata) AppendVersion(version Version) {
	e.versions = append(e.versions, version)
}

func (e *ObjectMetadata) DeleteVersion(versionNum int) error {
	for i, version := range e.versions {
		if version.Number() == versionNum {
			e.versions = append(e.versions[:i], e.versions[i+1:]...)
			return nil
		}
	}
	return errors.New("version not exist")
}

func (e *ObjectMetadata) LastVersion() int {
	if len(e.versions) == 0 {
		return -1
	}
	return e.versions[len(e.versions)-1].Number()
}

func (e *ObjectMetadata) MarshalBSON() ([]byte, error) {
	dto := struct {
		ID         ObjectID  `bson:"object_id"`
		Group      string    `bson:"group"`
		Partition  string    `bson:"partition"`
		Name       string    `bson:"name"`
		Path       string    `bson:"path"`
		Size       int       `bson:"size"`
		Node       Node      `bson:"node"`
		Versions   Versions  `bson:"versions"`
		CreatedAt  time.Time `bson:"created_at"`
		ModifiedAt time.Time `bson:"modified_at"`
	}{
		ID:         e.id,
		Group:      e.group,
		Partition:  e.partition,
		Name:       e.name,
		Path:       e.path,
		Versions:   e.versions,
		CreatedAt:  e.CreatedAt,
		ModifiedAt: e.ModifiedAt,
	}

	return bson.Marshal(dto)
}

func (e *ObjectMetadata) UnmarshalBSON(data []byte) error {
	dto := struct {
		ID         ObjectID  `bson:"object_id"`
		Group      string    `bson:"group"`
		Partition  string    `bson:"partition"`
		Name       string    `bson:"name"`
		Path       string    `bson:"path"`
		Versions   Versions  `bson:"versions"`
		CreatedAt  time.Time `bson:"created_at"`
		ModifiedAt time.Time `bson:"modified_at"`
	}{}

	if err := bson.Unmarshal(data, &dto); err != nil {
		return err
	}

	e.id = dto.ID
	e.group = dto.Group
	e.partition = dto.Partition
	e.name = dto.Name
	e.path = dto.Path
	e.versions = dto.Versions
	e.CreatedAt = dto.CreatedAt
	e.ModifiedAt = dto.ModifiedAt
	return nil
}

type ObjectMetadataBuilder struct {
	id         ObjectID
	group      string
	partition  string
	name       string
	path       string
	versions   Versions
	createdAt  time.Time
	modifiedAt time.Time
}

func NewObjectMetadataBuilder() *ObjectMetadataBuilder {
	return &ObjectMetadataBuilder{}
}

func (b *ObjectMetadataBuilder) ID(id ObjectID) *ObjectMetadataBuilder {
	b.id = id
	return b
}

func (b *ObjectMetadataBuilder) Group(group string) *ObjectMetadataBuilder {
	b.group = group
	return b
}

func (b *ObjectMetadataBuilder) Partition(partition string) *ObjectMetadataBuilder {
	b.partition = partition
	return b
}

func (b *ObjectMetadataBuilder) Name(name string) *ObjectMetadataBuilder {
	b.name = name
	return b
}

func (b *ObjectMetadataBuilder) Path(path string) *ObjectMetadataBuilder {
	b.path = path
	return b
}

func (b *ObjectMetadataBuilder) Versions(versions Versions) *ObjectMetadataBuilder {
	b.versions = versions
	return b
}

func (b *ObjectMetadataBuilder) CreatedAt(createAt time.Time) *ObjectMetadataBuilder {
	b.createdAt = createAt
	return b
}

func (b *ObjectMetadataBuilder) ModifiedAt(modifiedAt time.Time) *ObjectMetadataBuilder {
	b.modifiedAt = modifiedAt
	return b
}

func (b *ObjectMetadataBuilder) Build() ObjectMetadata {
	return ObjectMetadata{
		id:        b.id,
		group:     b.group,
		partition: b.partition,
		name:      b.name,
		path:      b.path,
		versions:  b.versions,
		ModifiedTime: ModifiedTime{
			CreatedAt:  b.createdAt,
			ModifiedAt: b.modifiedAt,
		},
	}
}
