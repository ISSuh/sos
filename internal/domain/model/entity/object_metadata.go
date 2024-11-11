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
)

type ObjectMetadataList []ObjectMetadata

type ObjectMetadata struct {
	id        ObjectID
	group     string
	partition string
	name      string
	path      string
	size      int
	node      Node
	versions  Versions

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

func (e *ObjectMetadata) Size() int {
	return e.size
}

func (e *ObjectMetadata) Node() Node {
	return e.node
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

type ObjectMetadataBuilder struct {
	id         ObjectID
	group      string
	partition  string
	name       string
	path       string
	versions   Versions
	size       int
	node       Node
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

func (b *ObjectMetadataBuilder) Size(size int) *ObjectMetadataBuilder {
	b.size = size
	return b
}

func (b *ObjectMetadataBuilder) Node(node Node) *ObjectMetadataBuilder {
	b.node = node
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
		size:      b.size,
		node:      b.node,
		versions:  b.versions,
		ModifiedTime: ModifiedTime{
			CreatedAt:  b.createdAt,
			ModifiedAt: b.modifiedAt,
		},
	}
}
