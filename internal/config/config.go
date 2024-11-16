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

package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v2"
)

type SosConfig struct {
	Api              ApiConfig              `yaml:"api"`
	MetadataRegistry MetadataRegistryConfig `yaml:"metadata_registry"`
	BlockStorage     BlockStorageConfig     `yaml:"block_storage"`
	Standalone       bool                   `yaml:"standalone"`
}

type Config struct {
	SOS SosConfig `yaml:"sos"`
}

func (c Config) Validate() error {
	if err := c.SOS.Api.Validate(); err != nil {
		return err
	}
	if err := c.SOS.MetadataRegistry.Validate(c.SOS.Standalone); err != nil {
		return err
	}
	if err := c.SOS.BlockStorage.Validate(c.SOS.Standalone); err != nil {
		return err
	}
	return nil
}

func NewConfig(path string) (Config, error) {
	if len(path) == 0 {
		return Config{}, errors.New("can not found config file")
	}

	buffer, err := loadFile(path)
	if err != nil {
		return Config{}, err
	}

	config := Config{}
	if err = yaml.Unmarshal(buffer, &config); err != nil {
		return Config{}, nil
	}

	if err = config.Validate(); err != nil {
		return Config{}, err
	}

	return config, nil
}

func loadFile(path string) (buffer []byte, err error) {
	if buffer, err = os.ReadFile(path); err != nil {
		return nil, err
	}
	return buffer, nil
}
