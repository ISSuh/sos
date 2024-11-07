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
	"github.com/ISSuh/sos/internal/domain/model/entity"
	"github.com/ISSuh/sos/internal/domain/model/message"
	"github.com/ISSuh/sos/pkg/empty"
	"github.com/ISSuh/sos/pkg/validation"
)

type Object struct {
	ID           entity.ObjectID `json:"object_id"`
	Group        string          `json:"group"`
	Partition    string          `json:"partition"`
	Name         string          `json:"name"`
	Path         string          `json:"path"`
	Size         int             `json:"size"`
	VersionNum   int             `json:"version"`
	BlockHeaders BlockHeaders    `json:"block_headers"`
}

func NewObjectFromMessage(m *message.Object) Object {
	switch {
	case validation.IsNil(m):
		return empty.Struct[Object]()
	case validation.IsNil(m.GetId()):
		return empty.Struct[Object]()
	}

	headers := make(BlockHeaders, 0, len(m.BlockHeaders))
	for _, h := range m.BlockHeaders {
		headers = append(headers, NewBlockHeaderFromMessage(h))
	}

	return Object{
		ID:           entity.ObjectID(m.GetId().Id),
		Group:        m.Group,
		Partition:    m.Partition,
		Name:         m.Name,
		Path:         m.Path,
		Size:         int(m.GetSize()),
		BlockHeaders: headers,
	}
}

func (o *Object) ToEntity() entity.ObjectMetadata {
	headers := o.BlockHeaders.ToEntity()
	return entity.NewObjectMetadataBuilder().
		ID(o.ID).
		Group(o.Group).
		Partition(o.Partition).
		Name(o.Name).
		Path(o.Path).
		Size(o.Size).
		BlockHeaders(headers).
		Build()
}

func (o *Object) ToMessage() *message.Object {
	headers := o.BlockHeaders.ToMessage()
	return &message.Object{
		Id:           &message.ObjectID{Id: o.ID.ToInt64()},
		Group:        o.Group,
		Partition:    o.Partition,
		Name:         o.Name,
		Path:         o.Path,
		Size:         int32(o.Size),
		BlockHeaders: headers,
	}
}
