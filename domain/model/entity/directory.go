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

type Directories []Directory

type Directory struct {
	id           ObjectID
	group        string
	partition    string
	name         string
	objects      []ObjectID
	subDirectory []ObjectID

	ModifiedTime
}

func (e *Directory) ID() ObjectID {
	return e.id
}

func (e *Directory) Group() string {
	return e.group
}

func (e *Directory) Partition() string {
	return e.partition
}

func (e *Directory) Name() string {
	return e.name
}

func (e *Directory) Objects() []ObjectID {
	return e.objects
}

func (e *Directory) SubDirectory() []ObjectID {
	return e.subDirectory
}

func (e *Directory) AddObjectID(id ObjectID) {
	e.objects = append(e.objects, id)
}

func (e *Directory) AddChild(id ObjectID) {
	e.subDirectory = append(e.subDirectory, id)
}

type DirectoryBuilder struct {
	id           ObjectID
	group        string
	partition    string
	name         string
	objects      []ObjectID
	subDirectory []ObjectID
	createdAt    time.Time
	modifiedAt   time.Time
}

func NewDirectoryBuilder() *DirectoryBuilder {
	return &DirectoryBuilder{}
}

func (b *DirectoryBuilder) ID(id ObjectID) *DirectoryBuilder {
	b.id = id
	return b
}

func (b *DirectoryBuilder) Group(group string) *DirectoryBuilder {
	b.group = group
	return b
}

func (b *DirectoryBuilder) Partition(partition string) *DirectoryBuilder {
	b.partition = partition
	return b
}

func (b *DirectoryBuilder) Name(name string) *DirectoryBuilder {
	b.name = name
	return b
}

func (b *DirectoryBuilder) Objects(objects []ObjectID) *DirectoryBuilder {
	b.objects = objects
	return b
}

func (b *DirectoryBuilder) SubDirectory(subDirectory []ObjectID) *DirectoryBuilder {
	b.subDirectory = subDirectory
	return b
}

func (b *DirectoryBuilder) CreatedAt(createAt time.Time) *DirectoryBuilder {
	b.createdAt = createAt
	return b
}

func (b *DirectoryBuilder) ModifiedAt(modifiedAt time.Time) *DirectoryBuilder {
	b.modifiedAt = modifiedAt
	return b
}

func (b *DirectoryBuilder) Build() *Directory {
	return &Directory{
		id:           b.id,
		group:        b.group,
		partition:    b.partition,
		name:         b.name,
		objects:      b.objects,
		subDirectory: b.subDirectory,
		ModifiedTime: ModifiedTime{
			CreatedAt:  b.createdAt,
			ModifiedAt: b.modifiedAt,
		},
	}
}
