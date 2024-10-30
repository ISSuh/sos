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

type BlockHeaders []BlockHeader

type BlockHeader struct {
	blockID   BlockID
	objectID  ObjectID
	index     int
	size      int
	node      Node
	timestamp time.Time
	checksum  uint32
}

func NewEmptyBlockHeader() BlockHeader {
	return BlockHeader{}
}

func (b *BlockHeader) BlockID() BlockID {
	return b.blockID
}

func (b *BlockHeader) ObjectID() ObjectID {
	return b.objectID
}

func (b *BlockHeader) Index() int {
	return b.index
}

func (b *BlockHeader) Size() int {
	return b.size
}

func (b *BlockHeader) Node() Node {
	return b.node
}

func (b *BlockHeader) Timestamp() time.Time {
	return b.timestamp
}

func (b *BlockHeader) Checksum() uint32 {
	return b.checksum
}

func (b *BlockHeader) Empty() bool {
	return b.blockID == 0
}

type BlockHeaderBuilder struct {
	blockID   BlockID
	objectID  ObjectID
	index     int
	size      int
	node      Node
	timestamp time.Time
	checksum  uint32
}

func NewBlockHeaderBuilder() *BlockHeaderBuilder {
	return &BlockHeaderBuilder{}
}

func (b *BlockHeaderBuilder) ObjectID(objectID ObjectID) *BlockHeaderBuilder {
	b.objectID = objectID
	return b
}

func (b *BlockHeaderBuilder) BlockID(blockID BlockID) *BlockHeaderBuilder {
	b.blockID = blockID
	return b
}

func (b *BlockHeaderBuilder) Index(index int) *BlockHeaderBuilder {
	b.index = index
	return b
}

func (b *BlockHeaderBuilder) Size(size int) *BlockHeaderBuilder {
	b.size = size
	return b
}

func (b *BlockHeaderBuilder) Node(node Node) *BlockHeaderBuilder {
	b.node = node
	return b
}

func (b *BlockHeaderBuilder) Timestamp(timestamp time.Time) *BlockHeaderBuilder {
	b.timestamp = timestamp
	return b
}

func (b *BlockHeaderBuilder) Checksum(checksum uint32) *BlockHeaderBuilder {
	b.checksum = checksum
	return b
}

func (b *BlockHeaderBuilder) Build() BlockHeader {
	return BlockHeader{
		blockID:   b.blockID,
		objectID:  b.objectID,
		index:     b.index,
		size:      b.size,
		node:      b.node,
		timestamp: b.timestamp,
		checksum:  b.checksum,
	}
}
