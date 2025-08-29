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
	"bytes"
	"encoding/gob"
	"errors"

	"github.com/ISSuh/sos/internal/crc"
)

const (
	BlockSize = 4 * 1024 * 1024
)

type Blocks []Block

type Block struct {
	header BlockHeader
	buffer []byte

	ModifiedTime
}

func (b *Block) Validate() error {
	if err := b.header.Validate(); err != nil {
		return err
	}

	if len(b.buffer) > BlockSize {
		return errors.New("block size is too large")
	}

	return nil
}

func (b *Block) ObjectID() ObjectID {
	return b.header.ObjectID()
}

func (b *Block) BlockID() BlockID {
	return b.header.BlockID()
}

func (b *Block) Index() int {
	return b.header.Index()
}

func (b *Block) Header() BlockHeader {
	return b.header
}

func (b *Block) Buffer() []byte {
	return b.buffer
}

func (b *Block) GobEncode() ([]byte, error) {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)

	if err := enc.Encode(&b.header); err != nil {
		return nil, err
	}
	if err := enc.Encode(b.buffer); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (b *Block) GobDecode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buffer)

	if err := dec.Decode(&b.header); err != nil {
		return err
	}
	if err := dec.Decode(&b.buffer); err != nil {
		return err
	}

	return nil
}

type BlockBuilder struct {
	header BlockHeader
	buffer []byte
}

func NewBlockBuilder() *BlockBuilder {
	return &BlockBuilder{}
}

func (b *BlockBuilder) Header(header BlockHeader) *BlockBuilder {
	b.header = header
	return b
}

func (b *BlockBuilder) Buffer(buffer []byte) *BlockBuilder {
	b.buffer = buffer
	return b
}

func (b *BlockBuilder) AppendBuffer(buffer []byte) *BlockBuilder {
	b.buffer = append(b.buffer, buffer...)
	return b
}

func (b *BlockBuilder) ReSizeBuffer(size uint64) *BlockBuilder {
	b.buffer = b.buffer[0:size]
	return b
}

func (b *BlockBuilder) BufferSize() int {
	return len(b.buffer)
}

func (b *BlockBuilder) CalculateChecksum() uint32 {
	return crc.Checksum(b.buffer)
}

func (b *BlockBuilder) Build() Block {
	return Block{
		header: b.header,
		buffer: b.buffer,
	}
}
