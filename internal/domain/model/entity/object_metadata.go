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
	ID        uint64
	Group     string
	Partition string
	Name      string
	Path      string
	Size      uint64

	ModifiedTime
}

func NewEmptyObjectMetadata() ObjectMetadata {
	return ObjectMetadata{}
}

func (e ObjectMetadata) IsEmpty() bool {
	return e == ObjectMetadata{}
}

type ObjectMetadataBuilder struct {
	id        uint64
	group     string
	partition string
	name      string
	path      string
	size      uint64
}

func NewObjectMetadataBuilder() *ObjectMetadataBuilder {
	return &ObjectMetadataBuilder{}
}

func (b *ObjectMetadataBuilder) ID(id uint64) *ObjectMetadataBuilder {
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

func (b *ObjectMetadataBuilder) Size(size uint64) *ObjectMetadataBuilder {
	b.size = size
	return b
}

func (b *ObjectMetadataBuilder) Build() ObjectMetadata {
	return ObjectMetadata{
		ID:        b.id,
		Group:     b.group,
		Partition: b.partition,
		Name:      b.name,
		Path:      b.path,
		Size:      b.size,
	}
}
