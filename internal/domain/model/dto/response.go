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
)

type Items []Item

func NewItemsFromMetadataList(m MetadataList) Items {
	items := make(Items, 0, len(m))
	for _, metadata := range m {
		items = append(items, NewItemFromMetadata(metadata))
	}
	return items
}

type Item struct {
	ID         entity.ObjectID `json:"object_id"`
	Group      string          `json:"group"`
	Partition  string          `json:"partition"`
	Path       string          `json:"path"`
	Name       string          `json:"name"`
	Versions   Versions        `json:"versions"`
	CreatedAt  time.Time       `json:"created_at"`
	ModifiedAt time.Time       `json:"modified_at"`
}

func NewItemFromMetadata(m Metadata) Item {
	return Item{
		ID:         m.ID,
		Group:      m.Group,
		Partition:  m.Partition,
		Path:       m.Path,
		Name:       m.Name,
		Versions:   m.Versions,
		CreatedAt:  m.CreatedAt,
		ModifiedAt: m.ModifiedAt,
	}
}

func (i Item) Empty() bool {
	return i.ID == 0
}
