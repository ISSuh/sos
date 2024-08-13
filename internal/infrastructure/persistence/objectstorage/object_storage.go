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

package objectstorage

import (
	"github.com/ISSuh/sos/internal/domain/repository"
	"github.com/ISSuh/sos/pkg/log"
)

type localObjectStorage struct {
	logger log.Logger

	storage map[string][]byte
}

func NewLocalObjectStorage(l log.Logger) (repository.ObjectStorage, error) {
	return &localObjectStorage{
			logger:  l,
			storage: make(map[string][]byte),
		},
		nil
}

func (s *localObjectStorage) Store() {

}

func (s *localObjectStorage) Delete() {

}

func (s *localObjectStorage) Find() {

}
