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

type Versions []Version

type Version struct {
	number       int
	size         int
	node         Node
	blockHeaders BlockHeaders

	ModifiedTime
}

func (e Version) Number() int {
	return e.number
}

func (e Version) Size() int {
	return e.size
}

func (e Version) Node() Node {
	return e.node
}

func (e Version) BlockHeaders() BlockHeaders {
	return e.blockHeaders
}

type VersionBuilder struct {
	number       int
	size         int
	node         Node
	blockHeaders BlockHeaders
	createdAt    time.Time
	modifiedAt   time.Time
}

func NewVersionBuilder() *VersionBuilder {
	return &VersionBuilder{}
}

func (b *VersionBuilder) Number(number int) *VersionBuilder {
	b.number = number
	return b
}

func (b *VersionBuilder) Size(size int) *VersionBuilder {
	b.size = size
	return b
}

func (b *VersionBuilder) Node(node Node) *VersionBuilder {
	b.node = node
	return b
}

func (b *VersionBuilder) BlockHeaders(blockHeaders BlockHeaders) *VersionBuilder {
	b.blockHeaders = blockHeaders
	return b
}

func (b *VersionBuilder) CreatedAt(createdAt time.Time) *VersionBuilder {
	b.createdAt = createdAt
	return b
}

func (b *VersionBuilder) ModifiedAt(modifiedAt time.Time) *VersionBuilder {
	b.modifiedAt = modifiedAt
	return b
}

func (b *VersionBuilder) Build() Version {
	return Version{
		number:       b.number,
		size:         b.size,
		node:         b.node,
		blockHeaders: b.blockHeaders,
		ModifiedTime: ModifiedTime{
			CreatedAt:  b.createdAt,
			ModifiedAt: b.modifiedAt,
		},
	}
}
