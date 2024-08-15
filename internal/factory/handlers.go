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
	"github.com/ISSuh/sos/internal/infrastructure/transport/rpc"
	rpchandler "github.com/ISSuh/sos/internal/infrastructure/transport/rpc/handler"
	"github.com/ISSuh/sos/pkg/log"
	sosrpc "github.com/ISSuh/sos/pkg/rpc"
	"github.com/ISSuh/sos/pkg/validation"
)

type Handlers struct {
	Uploader   rest.Uploader
	Downloader rest.Downloader
	Finder     rest.Finder
	Eraser     rest.Eraser
}

func NewHandlers(l log.Logger, serviceFactory *APIServices) (*Handlers, error) {
	switch {
	case validation.IsNil(l):
		return nil, fmt.Errorf("logger is nil")
	case validation.IsNil(serviceFactory):
		return nil, fmt.Errorf("service factory is nil")
	}

	finder, err := resthandler.NewFinder(l, serviceFactory.Finder)
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

	h := &Handlers{
		Finder:     finder,
		Uploader:   uploader,
		Downloader: downloader,
		Eraser:     eraser,
	}

	return h, nil
}

type RPCHandlers struct {
	MetadataRegistry rpc.MetadataRegistryServer
}

func NewRPCHandlers(l log.Logger, serviceFactory *APIServices) (*RPCHandlers, error) {
	switch {
	case validation.IsNil(l):
		return nil, fmt.Errorf("logger is nil")
	}

	h := &RPCHandlers{
		MetadataRegistry: rpchandler.NewMetadataRegistry(l),
	}
	return h, nil
}

func (f *RPCHandlers) Registers() []sosrpc.RegisterFunc {
	return []sosrpc.RegisterFunc{
		rpchandler.RegistMetadataRegistry(f.MetadataRegistry),
	}
}
