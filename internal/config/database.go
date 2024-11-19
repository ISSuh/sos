/*
MIT License

Copyright (c) 2024 ISSuh

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package config

import "fmt"

type DatabaseType string

const (
	DatabaseTypeLocal   DatabaseType = "local"
	DatabaseTypeMongoDB DatabaseType = "mongodb"
	DatabaseTypeLevelDB DatabaseType = "leveldb"
)

type Database struct {
	Type         DatabaseType `yaml:"type"`
	Host         string       `yaml:"host"`
	DatabaseName string       `yaml:"database"`
	Credentials  Credentials  `yaml:"credentials"`
	Options      Options      `yaml:"options"`
	LogLevel     string       `yaml:"log_level"`
	Path         string       `yaml:"path"`
}

func (d Database) Validate() error {
	switch d.Type {
	case DatabaseTypeLocal:
		return nil
	case DatabaseTypeMongoDB:
		return d.validateMogoDBConfig()
	case DatabaseTypeLevelDB:
		return d.validateLevelDBConfig()
	default:
		return fmt.Errorf("invalid database type. %s", d.Type)
	}
}

func (d Database) validateMogoDBConfig() error {
	if d.Host == "" {
		return fmt.Errorf("database host is empty")
	}
	if d.DatabaseName == "" {
		return fmt.Errorf("database name is empty")
	}
	return nil
}

func (d Database) validateLevelDBConfig() error {
	if d.Path == "" {
		return fmt.Errorf("database path is empty")
	}
	return nil
}

type Credentials struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Options struct {
	MaxPoolSize uint64 `yaml:"max_pool"`
	MinPoolSize uint64 `yaml:"min_pool"`
	MaxConnIdle uint64 `yaml:"max_conn_idle"`
}
