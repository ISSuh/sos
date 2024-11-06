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

type Versions []Version

func NewVersionsFromModel(v entity.Versions) Versions {
	versions := make(Versions, 0, len(v))
	for _, version := range v {
		versions = append(versions, NewVersionFromModel(version))
	}
	return versions
}

func NewVersionsFromMessage(v []*message.Version) Versions {
	versions := make(Versions, 0, len(v))
	for _, version := range v {
		versions = append(versions, NewVersionFromMessage(version))
	}
	return versions
}

func (v Versions) ToEntity() entity.Versions {
	versions := make(entity.Versions, 0, len(v))
	for _, version := range v {
		versions = append(versions, version.ToEntity())
	}
	return versions
}

func (v Versions) ToMessage() []*message.Version {
	versions := make([]*message.Version, 0, len(v))
	for _, version := range v {
		versions = append(versions, version.ToMessage())
	}
	return versions
}

func (v Versions) Empty() bool {
	return len(v) == 0
}

func (v Versions) LastVersion() Version {
	if v.Empty() {
		return empty.Struct[Version]()
	}
	return v[len(v)-1]
}

type Version struct {
	Number     int       `json:"number"`
	Size       int       `json:"size"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`

	BlockHeaders BlockHeaders
}

func NewVersionFromModel(v entity.Version) Version {
	headers := make(BlockHeaders, 0, len(v.BlockHeaders()))
	for _, h := range v.BlockHeaders() {
		headers = append(headers, NewBlockHeaderFromModel(h))
	}

	return Version{
		Number:       v.Number(),
		Size:         v.Size(),
		BlockHeaders: headers,
		CreatedAt:    v.CreatedAt,
		ModifiedAt:   v.ModifiedAt,
	}
}

func NewVersionFromMessage(m *message.Version) Version {
	switch {
	case validation.IsNil(m):
		return empty.Struct[Version]()
	}

	headers := make(BlockHeaders, 0, len(m.BlockHeaders))
	for _, h := range m.BlockHeaders {
		headers = append(headers, NewBlockHeaderFromMessage(h))
	}

	return Version{
		Number:       int(m.Number),
		Size:         int(m.Size),
		BlockHeaders: headers,
		CreatedAt:    m.CreatedAt.AsTime(),
		ModifiedAt:   m.ModifiedAt.AsTime(),
	}
}

func (v *Version) ToEntity() entity.Version {
	headers := v.BlockHeaders.ToEntity()
	return entity.NewVersionBuilder().
		Number(v.Number).
		Size(v.Size).
		BlockHeaders(headers).
		Build()
}

func (v *Version) ToMessage() *message.Version {
	headers := v.BlockHeaders.ToMessage()
	return &message.Version{
		Number:       int32(v.Number),
		Size:         int32(v.Size),
		BlockHeaders: headers,
	}
}

func (v *Version) IsValid() bool {
	return v.Number > -1
}
