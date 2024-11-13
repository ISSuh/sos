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

import "time"

type ObjectDirectories []ObjectDirectory

type ObjectDirectory struct {
	id        ObjectID
	group     string
	partition string
	name      string
	path      string
	depth     int

	directories ObjectDirectories
	metadata    ObjectMetadataList

	ModifiedTime
}

func (e *ObjectDirectory) ID() ObjectID {
	return e.id
}

func (e *ObjectDirectory) Group() string {
	return e.group
}

func (e *ObjectDirectory) Partition() string {
	return e.partition
}

func (e *ObjectDirectory) Name() string {
	return e.name
}

func (e *ObjectDirectory) Path() string {
	return e.path
}

func (e *ObjectDirectory) Depth() int {
	return e.depth
}

func (e *ObjectDirectory) Metadata() ObjectMetadataList {
	return e.metadata
}

func (e *ObjectDirectory) Directories() ObjectDirectories {
	return e.directories
}

func (e *ObjectDirectory) AddMetadata(metadata ObjectMetadata) {
	e.metadata = append(e.metadata, metadata)
}

func (e *ObjectDirectory) AddDirectory(directory ObjectDirectory) {
	e.directories = append(e.directories, directory)
}

type ObjectDirectoryBuilder struct {
	id          ObjectID
	group       string
	partition   string
	name        string
	path        string
	depth       int
	metadata    ObjectMetadataList
	directories ObjectDirectories
	createdAt   time.Time
	modifiedAt  time.Time
}

func NewObjectDirectoryBuilder() *ObjectDirectoryBuilder {
	return &ObjectDirectoryBuilder{}
}

func (b *ObjectDirectoryBuilder) ID(id ObjectID) *ObjectDirectoryBuilder {
	b.id = id
	return b
}

func (b *ObjectDirectoryBuilder) Group(group string) *ObjectDirectoryBuilder {
	b.group = group
	return b
}

func (b *ObjectDirectoryBuilder) Partition(partition string) *ObjectDirectoryBuilder {
	b.partition = partition
	return b
}

func (b *ObjectDirectoryBuilder) Name(name string) *ObjectDirectoryBuilder {
	b.name = name
	return b
}

func (b *ObjectDirectoryBuilder) Path(path string) *ObjectDirectoryBuilder {
	b.path = path
	return b
}

func (b *ObjectDirectoryBuilder) Depth(depth int) *ObjectDirectoryBuilder {
	b.depth = depth
	return b
}

func (b *ObjectDirectoryBuilder) CreatedAt(createdAt time.Time) *ObjectDirectoryBuilder {
	b.createdAt = createdAt
	return b
}

func (b *ObjectDirectoryBuilder) Metadata(metadata ObjectMetadataList) *ObjectDirectoryBuilder {
	b.metadata = metadata
	return b
}

func (b *ObjectDirectoryBuilder) Directories(directories ObjectDirectories) *ObjectDirectoryBuilder {
	b.directories = directories
	return b
}

func (b *ObjectDirectoryBuilder) ModifiedAt(modifiedAt time.Time) *ObjectDirectoryBuilder {
	b.modifiedAt = modifiedAt
	return b
}

func (b *ObjectDirectoryBuilder) Build() *ObjectDirectory {
	return &ObjectDirectory{
		id:          b.id,
		group:       b.group,
		partition:   b.partition,
		name:        b.name,
		depth:       b.depth,
		metadata:    b.metadata,
		directories: b.directories,
		ModifiedTime: ModifiedTime{
			CreatedAt:  b.createdAt,
			ModifiedAt: b.modifiedAt,
		},
	}
}
