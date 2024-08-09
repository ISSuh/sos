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
	"github.com/ISSuh/sos/internal/infrastructure/transport/rest/router"
	"github.com/ISSuh/sos/pkg/http"
	"github.com/ISSuh/sos/pkg/logger"
)

type Api struct {
	config   *config.SosConfig
	logger   logger.Logger
	handlers *factory.Handlers
	server   *http.Server
}

func NewApi(c *config.SosConfig, l logger.Logger) (*Api, error) {
	a := &Api{
		config: c,
		logger: l,
		server: http.NewServer(),
	}
	return a, nil
}

func (a *Api) Run() error {
	a.logger.Infof("[Api.Run]")
	if err := a.init(); err != nil {
		return err
	}
	return a.server.Run(a.config.Api.Address.String())
}

func (a *Api) init() error {
	var err error

	if err = a.initHandler(); err != nil {
		return err
	}

	router.Route(a.server, a.handlers)
	return nil
}

func (a *Api) initHandler() error {
	var err error
	var handlers *factory.Handlers

	if handlers, err = factory.NewHandlers(a.logger); err != nil {
		return err
	}

	a.handlers = handlers
	return nil
}
