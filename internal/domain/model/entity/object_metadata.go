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

type ObjectMetadata struct {
	id        ObjectID
	group     string
	partition string
	name      string
	path      string
	size      int
	node      Node

	ModifiedTime
}

func NewEmptyObjectMetadata() ObjectMetadata {
	return ObjectMetadata{}
}

func (e ObjectMetadata) ID() ObjectID {
	return e.id
}

func (e ObjectMetadata) Group() string {
	return e.group
}

func (e ObjectMetadata) Partition() string {
	return e.partition
}

func (e ObjectMetadata) Name() string {
	return e.name
}

func (e ObjectMetadata) Path() string {
	return e.path
}

func (e ObjectMetadata) Size() int {
	return e.size
}

func (e ObjectMetadata) Node() Node {
	return e.node
}

func (e ObjectMetadata) IsEmpty() bool {
	return e == ObjectMetadata{}
}

type ObjectMetadataBuilder struct {
	id        ObjectID
	group     string
	partition string
	name      string
	path      string
	size      int
	node      Node
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

func (b *ObjectMetadataBuilder) Build() ObjectMetadata {
	return ObjectMetadata{
		id:        b.id,
		group:     b.group,
		partition: b.partition,
		name:      b.name,
		path:      b.path,
		size:      b.size,
		node:      b.node,
	}
}
