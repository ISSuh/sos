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

package database

import (
	"context"
	"fmt"

	"github.com/ISSuh/sos/domain/model/entity"
	"github.com/ISSuh/sos/domain/repository"
	"github.com/ISSuh/sos/internal/log"
)

type localObjectDirectory struct {
	db map[string]map[int]map[string]*entity.ObjectDirectory
}

func NewLocalObjectDirectory() (repository.ObjectDirectory, error) {
	return &localObjectDirectory{
		db: make(map[string]map[int]map[string]*entity.ObjectDirectory),
	}, nil
}

func (d *localObjectDirectory) Create(c context.Context, dir *entity.ObjectDirectory) error {
	log.FromContext(c).Debugf("[localObjectDirectory.Create] dir: %+v", dir)
	key := d.makeKey(
		dir.Group(), dir.Partition(),
	)

	_, exist := d.db[key]
	if !exist {
		d.db[key] = make(map[int]map[string]*entity.ObjectDirectory)
	}

	_, exist = d.db[key][dir.Depth()]
	if !exist {
		d.db[key][dir.Depth()] = make(map[string]*entity.ObjectDirectory)
	}

	d.db[key][dir.Depth()][dir.Name()] = dir
	return nil
}

func (d *localObjectDirectory) Update(c context.Context, dir *entity.ObjectDirectory) error {
	log.FromContext(c).Debugf("[localObjectDirectory.Update] dir: %+v", dir)
	key := d.makeKey(
		dir.Group(), dir.Partition(),
	)

	_, exist := d.db[key]
	if !exist {
		return fmt.Errorf("dir key not exist")
	}

	_, exist = d.db[key][dir.Depth()]
	if !exist {
		return fmt.Errorf("not exist parent")
	}

	d.db[key][dir.Depth()][dir.Name()] = dir
	return nil
}

func (d *localObjectDirectory) Delete(c context.Context, dir *entity.ObjectDirectory) error {
	log.FromContext(c).Debugf("[localObjectDirectory.Delete] dir: %+v", dir)
	key := d.makeKey(
		dir.Group(), dir.Partition(),
	)

	_, exist := d.db[key]
	if !exist {
		return fmt.Errorf("dir key not exist")
	}

	_, exist = d.db[key][dir.Depth()]
	if !exist {
		return fmt.Errorf("not exist parent")
	}

	delete(d.db[key][dir.Depth()], dir.Name())
	return nil
}

func (d *localObjectDirectory) Find(c context.Context, group, partition, name string, depth int) (*entity.ObjectDirectory, error) {
	log.FromContext(c).Debugf("[localObjectDirectory.Find] group: %s, partition: %s, name: %s, depth: %d", group, partition, name, depth)
	key := d.makeKey(group, partition)

	_, exist := d.db[key]
	if !exist {
		return nil, fmt.Errorf("dir key not exist")
	}

	_, exist = d.db[key][depth]
	if !exist {
		return nil, fmt.Errorf("not exist parent")
	}

	return d.db[key][depth][name], nil
}

func (d *localObjectDirectory) makeKey(group, partition string) string {
	return fmt.Sprintf("%s:%s", group, partition)
}
