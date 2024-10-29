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

type MetadataList []Metadata

func (m MetadataList) Empty() bool {
	return len(m) == 0
}

type Metadata struct {
	ID           entity.ObjectID `json:"object_id"`
	Group        string          `json:"group"`
	Partition    string          `json:"partition"`
	Name         string          `json:"name"`
	Path         string          `json:"path"`
	Size         int             `json:"size"`
	BlockHeaders BlockHeaders    `json:"block_headers"`
	CreatedAt    time.Time       `json:"created_at"`
	ModifiedAt   time.Time       `json:"modified_at"`
}

func NewMetadataFromModel(m entity.ObjectMetadata) Metadata {
	headers := make(BlockHeaders, 0, len(m.BlockHeaders()))
	for _, h := range m.BlockHeaders() {
		headers = append(headers, NewBlockHeaderFromModel(h))
	}

	return Metadata{
		ID:           m.ID(),
		Group:        m.Group(),
		Partition:    m.Partition(),
		Name:         m.Name(),
		Path:         m.Path(),
		Size:         m.Size(),
		BlockHeaders: headers,
		CreatedAt:    m.CreatedAt,
		ModifiedAt:   m.ModifiedAt,
	}
}

func NewMetadataFromMessage(m *message.ObjectMetadata) Metadata {
	switch {
	case validation.IsNil(m):
		return empty.Struct[Metadata]()
	case validation.IsNil(m.GetId()):
		return empty.Struct[Metadata]()
	}

	headers := make(BlockHeaders, 0, len(m.BlockHeaders))
	for _, h := range m.BlockHeaders {
		headers = append(headers, NewBlockHeaderFromMessage(h))
	}

	return Metadata{
		ID:           entity.ObjectID(m.GetId().Id),
		Group:        m.Group,
		Partition:    m.Partition,
		Name:         m.Name,
		Path:         m.Path,
		Size:         int(m.GetSize()),
		BlockHeaders: headers,
		// CreatedAt:  m.CreatedAt,
		// ModifiedAt: m.ModifiedAt,
	}
}

func NewEmptyMetadata() Metadata {
	return Metadata{}
}

func (d Metadata) Empty() bool {
	return d.BlockHeaders.Empty() && d.ID == 0
}
