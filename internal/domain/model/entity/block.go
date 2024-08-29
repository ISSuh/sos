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

const (
	BlockSize = 4 * 1024 * 1024
)

type Block struct {
	id     uint64
	header BlockHeader
	data   []byte

	ModifiedTime
}

type Blocks []Block

func (b *Block) ID() uint64 {
	return b.id
}

func (b *Block) Header() BlockHeader {
	return b.header
}

func (b *Block) Data() []byte {
	return b.data
}

type BlockBuilder struct {
	id     uint64
	header BlockHeader
	data   []byte
}

func NewBlockBuilder() *BlockBuilder {
	return &BlockBuilder{}
}

func (b *BlockBuilder) ID(id uint64) *BlockBuilder {
	b.id = id
	return b
}

func (b *BlockBuilder) Header(header BlockHeader) *BlockBuilder {
	b.header = header
	return b
}

func (b *BlockBuilder) Data(data []byte) *BlockBuilder {
	b.data = data
	return b
}

func (b *BlockBuilder) Build() Block {
	return Block{
		id:     b.id,
		header: b.header,
		data:   b.data,
	}
}
