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

package factory

import (
	"fmt"

	"github.com/ISSuh/sos/internal/infrastructure/transport/rest"
	"github.com/ISSuh/sos/internal/infrastructure/transport/rest/handler"
	"github.com/ISSuh/sos/pkg/logger"
	"github.com/ISSuh/sos/pkg/validation"
)

type Handlers struct {
	Uploader   rest.Uploader
	Downloader rest.Downloader
	Finder     rest.Finder
	Eraser     rest.Eraser
}

func NewHandlers(l logger.Logger, serviceFactory *APIServices) (*Handlers, error) {
	switch {
	case validation.IsNil(l):
		return nil, fmt.Errorf("logger is nil")
	case validation.IsNil(serviceFactory):
		return nil, fmt.Errorf("service factory is nil")
	}

	finder, err := handler.NewFinder(l, serviceFactory.Finder)
	if err != nil {
		return nil, err
	}

	uploader, err := handler.NewUploader(l, serviceFactory.Uploader)
	if err != nil {
		return nil, err
	}

	downloader, err := handler.NewDownloader(l, serviceFactory.Downloader)
	if err != nil {
		return nil, err
	}

	eraser, err := handler.NewEraser(l, serviceFactory.Eraser)
	if err != nil {
		return nil, err
	}

	h := &Handlers{
		Finder:     finder,
		Uploader:   uploader,
		Downloader: downloader,
		Eraser:     eraser,
	}
	return h, nil
}
