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
	"github.com/ISSuh/sos/pkg/log"
	"github.com/ISSuh/sos/pkg/rpc"
)

type BlockStorage struct {
	logger log.Logger

	config config.SosConfig
	server rpc.Server
}

func NewBlockStorage(c config.SosConfig, l log.Logger) (BlockStorage, error) {
	a := BlockStorage{
		config: c,
		logger: l,
		server: rpc.NewServer(),
	}
	return a, nil
}

func (a *BlockStorage) Run() error {
	a.logger.Infof("[BlockStorage.Run]")
	if err := a.init(); err != nil {
		return err
	}
	return a.server.Run(a.config.BlockStorage.Address.String())
}

func (a *BlockStorage) init() error {
	repository, err := factory.NewObjectStorageRepository()
	if err != nil {
		return err
	}

	service, err := factory.NewObjectStorageService(repository)
	if err != nil {
		return err
	}

	registers, err := factory.BlockStorageHandler(service)
	if err != nil {
		return err
	}

	a.server.Regist(registers)
	return nil
}
