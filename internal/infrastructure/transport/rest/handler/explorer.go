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
	gohttp "net/http"

	"github.com/ISSuh/sos/internal/domain/model/dto"
	"github.com/ISSuh/sos/internal/domain/service"
	"github.com/ISSuh/sos/internal/infrastructure/transport/rest"
	"github.com/ISSuh/sos/pkg/http"
	"github.com/ISSuh/sos/pkg/log"
	"github.com/ISSuh/sos/pkg/validation"
)

type explorer struct {
	logger log.Logger

	explorerService service.Explorer
}

func NewExplorer(l log.Logger, explorerService service.Explorer) (rest.Explorer, error) {
	switch {
	case validation.IsNil(l):
		return nil, fmt.Errorf("logger is nil")
	case validation.IsNil(explorerService):
		return nil, fmt.Errorf("explorer service is nil")
	}

	return &explorer{
		logger:          l,
		explorerService: explorerService,
	}, nil
}

func (h *explorer) Find(w gohttp.ResponseWriter, r *gohttp.Request) {
	c := r.Context()
	dto := dto.RequestFromContext(c, http.RequestContextKey)
	log.FromContext(c).Debugf("[explorer.Find]")
	log.FromContext(c).Debugf("Request: %+v\n", dto)

	metadata, err := h.explorerService.FindObjectMetadata(c, dto)
	if err != nil {
		h.logger.Errorf(err.Error())
		gohttp.Error(w, err.Error(), gohttp.StatusInternalServerError)
		return
	}

	if err := http.Json(w, metadata); err != nil {
		h.logger.Errorf(err.Error())
		gohttp.Error(w, err.Error(), gohttp.StatusInternalServerError)
		return
	}
}

func (h *explorer) List(w gohttp.ResponseWriter, r *gohttp.Request) {
}
