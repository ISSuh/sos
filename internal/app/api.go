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
	"github.com/ISSuh/sos/internal/controller/rest/router"
	"github.com/ISSuh/sos/internal/factory"
	"github.com/gin-gonic/gin"
)

type Api struct {
	config *config.Config

	handlers *factory.Handlers
	engine   *gin.Engine
}

func NewApi(configPath string) (*Api, error) {
	config, err := config.NewConfig(configPath)
	if err != nil {
		return nil, err
	}

	a := &Api{
		config: config,
		engine: gin.Default(),
	}
	return a, nil
}

func (a *Api) Run() error {
	if err := a.init(); err != nil {
		return err
	}
	return a.engine.Run(a.config.Api.Address.String())
}

func (a *Api) init() error {
	var err error

	if err = a.initHandler(); err != nil {
		return err
	}

	if err = router.Route(a.engine, a.handlers); err != nil {
		return err
	}

	return nil
}

func (a *Api) initHandler() error {
	var err error
	var handlers *factory.Handlers

	if handlers, err = factory.NewHandlers(); err != nil {
		return err
	}

	a.handlers = handlers
	return nil
}
