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

	"github.com/ISSuh/sos/domain/model/entity"
)

type MetadataList []Metadata

type Metadata struct {
	ID         entity.ObjectID `json:"object_id"`
	Group      string          `json:"group"`
	Partition  string          `json:"partition"`
	Name       string          `json:"name"`
	Path       string          `json:"path"`
	Versions   Versions        `json:"versions"`
	CreatedAt  time.Time       `json:"created_at"`
	ModifiedAt time.Time       `json:"modified_at"`
}

func NewMetadataFromModel(m *entity.ObjectMetadata) *Metadata {
	return &Metadata{
		ID:         m.ID(),
		Group:      m.Group(),
		Partition:  m.Partition(),
		Name:       m.Name(),
		Path:       m.Path(),
		Versions:   NewVersionsFromModel(m.Versions()),
		CreatedAt:  m.CreatedAt,
		ModifiedAt: m.ModifiedAt,
	}
}

func (d *Metadata) Empty() bool {
	return d.Versions.Empty() && d.ID == 0
}

func (d *Metadata) ToEntity() entity.ObjectMetadata {
	return entity.NewObjectMetadataBuilder().
		ID(d.ID).
		Group(d.Group).
		Partition(d.Partition).
		Name(d.Name).
		Path(d.Path).
		Versions(d.Versions.ToEntity()).
		CreatedAt(d.CreatedAt).
		ModifiedAt(d.ModifiedAt).
		Build()
}
