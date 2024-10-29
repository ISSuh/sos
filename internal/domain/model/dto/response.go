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

func (m Items) Empty() bool {
	return len(m) == 0
}

func NewItemsFromMetadataList(d MetadataList) Items {
	var items Items
	for _, v := range d {
		items = append(items, NewItemFromMetadata(v))
	}

	return items
}

type Item struct {
	ID         entity.ObjectID `json:"object_id"`
	Group      string          `json:"group"`
	Partition  string          `json:"partition"`
	Name       string          `json:"name"`
	Path       string          `json:"path"`
	Size       int             `json:"size"`
	CreatedAt  time.Time       `json:"created_at"`
	ModifiedAt time.Time       `json:"modified_at"`
}

func NewItemFromMetadata(d Metadata) Item {
	return Item{
		ID:         d.ID,
		Group:      d.Group,
		Partition:  d.Partition,
		Name:       d.Name,
		Path:       d.Path,
		Size:       d.Size,
		CreatedAt:  d.CreatedAt,
		ModifiedAt: d.ModifiedAt,
	}
}

func NewItemFromModel(e entity.ObjectMetadata) Item {
	return Item{
		ID:         e.ID(),
		Group:      e.Group(),
		Partition:  e.Partition(),
		Name:       e.Name(),
		Path:       e.Path(),
		Size:       e.Size(),
		CreatedAt:  e.CreatedAt,
		ModifiedAt: e.ModifiedAt,
	}
}

func (i Item) Empty() bool {
	return i.ID == 0
}
