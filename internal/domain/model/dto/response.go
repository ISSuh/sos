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
)

type Items []Item

func NewItemFromMetadataListMessage(list []*message.ObjectMetadata) Items {
	items := Items{}
	for _, m := range list {
		items = append(items, NewItemFromMetadataMessage(m))
	}
	return items
}

type Item struct {
	ID         entity.ObjectID `json:"object_id"`
	Group      string          `json:"group"`
	Partition  string          `json:"partition"`
	Path       string          `json:"path"`
	Name       string          `json:"name"`
	Size       int             `json:"size"`
	CreatedAt  time.Time       `json:"created_at"`
	ModifiedAt time.Time       `json:"modified_at"`
}

func NewItemFromMetadataMessage(m *message.ObjectMetadata) Item {
	return Item{
		ID:         entity.ObjectID(m.Id.Id),
		Group:      m.Group,
		Partition:  m.Partition,
		Path:       m.Path,
		Name:       m.Name,
		Size:       int(m.Size),
		CreatedAt:  m.CreatedAt.AsTime(),
		ModifiedAt: m.ModifiedAt.AsTime(),
	}
}

func NewItemFromMetadataModel(e entity.ObjectMetadata) Item {
	return Item{
		ID:         e.ID(),
		Group:      e.Group(),
		Partition:  e.Partition(),
		Path:       e.Path(),
		Name:       e.Name(),
		Size:       e.Size(),
		CreatedAt:  e.CreatedAt,
		ModifiedAt: e.ModifiedAt,
	}
}

func (i Item) Empty() bool {
	return i.ID == 0
}
