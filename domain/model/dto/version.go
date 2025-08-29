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
	"errors"
	"time"

	"github.com/ISSuh/sos/domain/model/entity"
	"github.com/ISSuh/sos/internal/empty"
)

type Versions []Version

func NewVersionsFromModel(v entity.Versions) Versions {
	versions := make(Versions, 0, len(v))
	for _, version := range v {
		versions = append(versions, NewVersionFromModel(version))
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

func (v Versions) Empty() bool {
	return len(v) == 0
}

func (v Versions) Version(versionNum int) (Version, error) {
	for _, version := range v {
		if version.Number == versionNum {
			return version, nil
		}
	}
	return empty.Struct[Version](), errors.New("version not exist")
}

func (v Versions) LastVersion() (Version, error) {
	if v.Empty() {
		return empty.Struct[Version](), errors.New("version not exist")
	}
	return v[len(v)-1], nil
}

func (v Versions) HasVersion(versionNum int) bool {
	for _, version := range v {
		if version.Number == versionNum {
			return true
		}
	}
	return false
}

type Version struct {
	Number       int          `json:"number"`
	Size         int          `json:"size"`
	BlockHeaders BlockHeaders `json:"-"`
	CreatedAt    time.Time    `json:"created_at"`
	ModifiedAt   time.Time    `json:"modified_at"`
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

func (v *Version) ToEntity() entity.Version {
	headers := v.BlockHeaders.ToEntity()
	return entity.NewVersionBuilder().
		Number(v.Number).
		Size(v.Size).
		BlockHeaders(headers).
		CreatedAt(v.CreatedAt).
		ModifiedAt(v.ModifiedAt).
		Build()
}

func (v *Version) IsValid() bool {
	return v.Number > -1
}
