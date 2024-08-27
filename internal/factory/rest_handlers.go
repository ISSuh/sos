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
	resthandler "github.com/ISSuh/sos/internal/infrastructure/transport/rest/handler"
	"github.com/ISSuh/sos/pkg/log"
	"github.com/ISSuh/sos/pkg/validation"
)

type RestHandlers struct {
	Uploader   rest.Uploader
	Downloader rest.Downloader
	Explorer   rest.Explorer
	Eraser     rest.Eraser
}

func NewHandlers(l log.Logger, serviceFactory *APIServices) (*RestHandlers, error) {
	switch {
	case validation.IsNil(l):
		return nil, fmt.Errorf("logger is nil")
	case validation.IsNil(serviceFactory):
		return nil, fmt.Errorf("service factory is nil")
	}

	explorer, err := resthandler.NewExplorer(l, serviceFactory.Explorer)
	if err != nil {
		return nil, err
	}

	uploader, err := resthandler.NewUploader(l, serviceFactory.Uploader)
	if err != nil {
		return nil, err
	}

	downloader, err := resthandler.NewDownloader(l, serviceFactory.Downloader)
	if err != nil {
		return nil, err
	}

	eraser, err := resthandler.NewEraser(l, serviceFactory.Eraser)
	if err != nil {
		return nil, err
	}

	h := &RestHandlers{
		Explorer:   explorer,
		Uploader:   uploader,
		Downloader: downloader,
		Eraser:     eraser,
	}

	return h, nil
}
