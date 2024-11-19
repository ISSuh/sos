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
	"github.com/ISSuh/sos/domain/service"
	"github.com/ISSuh/sos/infrastructure/transport/rest/router"
	"github.com/ISSuh/sos/internal/app/standalone"
	"github.com/ISSuh/sos/internal/config"
	"github.com/ISSuh/sos/internal/factory"
	"github.com/ISSuh/sos/internal/http"
	"github.com/ISSuh/sos/internal/log"
)

type Standalone struct {
	logger log.Logger

	config config.SosConfig
	server http.Server
}

func NewStandalone(c config.SosConfig, l log.Logger) (Standalone, error) {
	a := Standalone{
		config: c,
		logger: l,
		server: http.NewServer(),
	}
	return a, nil
}

func (a *Standalone) Run() error {
	a.logger.Infof("[Standalone.Run]")
	if err := a.init(); err != nil {
		return err
	}
	return a.server.Run(a.config.Api.Address.String())
}

func (a *Standalone) init() error {
	service, err := a.initService()
	if err != nil {
		return err
	}

	handler, err := factory.NewExplorerHandler(service)
	if err != nil {
		return err
	}

	router.Route(a.logger, &a.server, handler)
	return nil
}

func (a *Standalone) initService() (service.Explorer, error) {
	metadataRepo, err := factory.NewObjectMetadataRepository(a.logger, a.config.MetadataRegistry.Database)
	if err != nil {
		return nil, err
	}

	metadataService, err := factory.NewObjectMetadataService(metadataRepo)
	if err != nil {
		return nil, err
	}

	storageRepo, err := factory.NewObjectStorageRepository(a.logger, a.config.BlockStorage.Database)
	if err != nil {
		return nil, err
	}

	storageService, err := factory.NewObjectStorageService(storageRepo)
	if err != nil {
		return nil, err
	}

	metadataRegistry, err := standalone.NewMetadataRegistry(metadataService)
	if err != nil {
		return nil, err
	}

	blockStorage, err := standalone.NewBlockStorage(storageService)
	if err != nil {
		return nil, err
	}

	explorer, err := factory.NewExplorerService(metadataRegistry, blockStorage)
	if err != nil {
		return nil, err
	}

	return explorer, nil
}
