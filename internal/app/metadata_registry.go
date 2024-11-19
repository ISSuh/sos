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

package app

import (
	"github.com/ISSuh/sos/internal/config"
	"github.com/ISSuh/sos/internal/factory"
	"github.com/ISSuh/sos/internal/log"
	"github.com/ISSuh/sos/internal/rpc"
)

type MetadataRegistry struct {
	logger log.Logger

	config config.SosConfig
	server rpc.Server
}

func NewMetadata(c config.SosConfig, l log.Logger) (MetadataRegistry, error) {
	a := MetadataRegistry{
		config: c,
		logger: l,
		server: rpc.NewServer(),
	}
	return a, nil
}

func (a *MetadataRegistry) Run() error {
	a.logger.Infof("[MetadataRegistry.Run]")
	if err := a.init(); err != nil {
		return err
	}
	return a.server.Run(a.config.MetadataRegistry.Address.String())
}

func (a *MetadataRegistry) init() error {
	repository, err := factory.NewObjectMetadataRepository(a.logger, a.config.MetadataRegistry.Database)
	if err != nil {
		return err
	}

	service, err := factory.NewObjectMetadataService(repository)
	if err != nil {
		return err
	}

	registers, err := factory.MetadataRegistryHandler(service)
	if err != nil {
		return err
	}

	a.server.Regist(registers)
	return nil
}
