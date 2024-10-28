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
)

type BlockHeaders []BlockHeader

func (h BlockHeaders) Empty() bool {
	return len(h) == 0
}

type BlockHeader struct {
	BlockID   entity.BlockID  `json:"block_id"`
	ObjectID  entity.ObjectID `json:"object_id"`
	Index     int             `json:"index"`
	Size      int             `json:"size"`
	Timestamp time.Time       `json:"timestamp"`
}

func NewBlockHeaderFromModel(h entity.BlockHeader) BlockHeader {
	return BlockHeader{
		BlockID:   h.BlockID(),
		ObjectID:  h.ObjectID(),
		Index:     h.Index(),
		Size:      h.Size(),
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
		Timestamp: h.Timestamp.AsTime(),
	}
}

func NewEmptyBlockHeader() BlockHeader {
	return BlockHeader{}
}

func (d BlockHeader) IsEmpty() bool {
	return d == BlockHeader{}
}
