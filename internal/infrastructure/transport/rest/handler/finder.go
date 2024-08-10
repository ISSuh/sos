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

package handler

import (
	"fmt"
	"net/http"

	"github.com/ISSuh/sos/internal/domain/service"
	"github.com/ISSuh/sos/internal/infrastructure/transport/rest"
	"github.com/ISSuh/sos/pkg/logger"
	"github.com/ISSuh/sos/pkg/validation"
)

type finder struct {
	logger logger.Logger

	findService service.Finder
}

func NewFinder(l logger.Logger, findService service.Finder) (rest.Finder, error) {
	switch {
	case validation.IsNil(l):
		return nil, fmt.Errorf("logger is nil")
	case validation.IsNil(findService):
		return nil, fmt.Errorf("find service is nil")
	}

	return &finder{
		logger:      l,
		findService: findService,
	}, nil
}

func (h *finder) Find(w http.ResponseWriter, r *http.Request) {
	h.logger.Debugf("[finder.Find]")
}

func (h *finder) List(w http.ResponseWriter, r *http.Request) {
	h.logger.Debugf("[finder.List]")
}
