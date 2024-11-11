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

package dto

import (
	"time"

	"github.com/ISSuh/sos/internal/domain/model/entity"
	"github.com/ISSuh/sos/internal/domain/model/message"
	"github.com/ISSuh/sos/pkg/empty"
	"github.com/ISSuh/sos/pkg/validation"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type BlockHeaders []BlockHeader

func (h BlockHeaders) Empty() bool {
	return len(h) == 0
}

func (h BlockHeaders) ToEntity() entity.BlockHeaders {
	headers := make(entity.BlockHeaders, 0, len(h))
	for _, header := range h {
		headers = append(headers, header.ToEntity())
	}
	return headers
}

func (h BlockHeaders) ToMessage() []*message.BlockHeader {
	headers := make([]*message.BlockHeader, 0, len(h))
	for _, header := range h {
		headers = append(headers, header.ToMessage())
	}
	return headers
}

type BlockHeader struct {
	BlockID   entity.BlockID  `json:"block_id"`
	ObjectID  entity.ObjectID `json:"object_id"`
	Index     int             `json:"index"`
	Size      int             `json:"size"`
	Timestamp time.Time       `json:"timestamp"`
	Checksum  uint32          `json:"-"`
}

func NewBlockHeaderFromModel(h entity.BlockHeader) BlockHeader {
	return BlockHeader{
		BlockID:   h.BlockID(),
		ObjectID:  h.ObjectID(),
		Index:     h.Index(),
		Size:      h.Size(),
		Checksum:  h.Checksum(),
		Timestamp: h.Timestamp(),
	}
}

func NewBlockHeaderFromMessage(h *message.BlockHeader) BlockHeader {
	switch {
	case validation.IsNil(h):
		return empty.Struct[BlockHeader]()
	case validation.IsNil(h.GetObjectID()):
		return empty.Struct[BlockHeader]()
	case validation.IsNil(h.GetBlockID()):
		return empty.Struct[BlockHeader]()
	}

	return BlockHeader{
		BlockID:   entity.BlockID(h.BlockID.Id),
		ObjectID:  entity.ObjectID(h.ObjectID.Id),
		Index:     int(h.Index),
		Size:      int(h.Size),
		Checksum:  h.Checksum,
		Timestamp: h.Timestamp.AsTime(),
	}
}

func NewEmptyBlockHeader() BlockHeader {
	return BlockHeader{}
}

func (d BlockHeader) Empty() bool {
	return d == BlockHeader{}
}

func (d BlockHeader) ToEntity() entity.BlockHeader {
	return entity.NewBlockHeaderBuilder().
		BlockID(d.BlockID).
		ObjectID(d.ObjectID).
		Index(d.Index).
		Size(d.Size).
		Timestamp(d.Timestamp).
		Checksum(d.Checksum).
		Build()
}

func (d BlockHeader) ToMessage() *message.BlockHeader {
	return &message.BlockHeader{
		ObjectID:  message.FromObjectID(d.ObjectID),
		BlockID:   message.FromBlockID(d.BlockID),
		Index:     int32(d.Index),
		Size:      int32(d.Size),
		Checksum:  d.Checksum,
		Timestamp: timestamppb.New(d.Timestamp),
	}
}

type Block struct {
	Header BlockHeader
	Data   []byte
}

func (d Block) ToMessage() *message.Block {
	return &message.Block{
		Header: d.Header.ToMessage(),
		Data:   d.Data,
	}
}
