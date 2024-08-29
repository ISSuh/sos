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

type BlockHeader struct {
	id        BlockID
	objectID  ObjectID
	index     uint64
	size      uint32
	node      Node
	timestamp time.Time
	checksum  string
}

func NewEmptyBlockHeader() BlockHeader {
	return BlockHeader{}
}

func (b *BlockHeader) ID() BlockID {
	return b.id
}

func (b *BlockHeader) ObjectID() ObjectID {
	return b.objectID
}

func (b *BlockHeader) Index() uint64 {
	return b.index
}

func (b *BlockHeader) Size() uint32 {
	return b.size
}

func (b *BlockHeader) Node() Node {
	return b.node
}

func (b *BlockHeader) Timestamp() time.Time {
	return b.timestamp
}

func (b *BlockHeader) Checksum() string {
	return b.checksum
}

func (b *BlockHeader) IsEmpty() bool {
	return b.id == 0
}
