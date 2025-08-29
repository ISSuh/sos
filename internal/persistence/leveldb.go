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

package persistence

import (
	"fmt"

	"github.com/ISSuh/sos/internal/config"
	"github.com/syndtr/goleveldb/leveldb"
)

type LevelDB struct {
	engin *leveldb.DB
}

func NewLevelDB(dbConfig config.Database) (*LevelDB, error) {
	db, err := leveldb.OpenFile(dbConfig.Path, nil)
	if err != nil {
		return nil, err
	}

	return &LevelDB{
		engin: db,
	}, nil
}

func CloseLevelDB(db *leveldb.DB) error {
	return db.Close()
}

func (l *LevelDB) Engin() (*leveldb.DB, error) {
	if l.engin == nil {
		return nil, fmt.Errorf("leveldb is nil")
	}
	return l.engin, nil
}
