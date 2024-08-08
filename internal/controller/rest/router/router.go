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

package router

import (
	gohttp "net/http"

	"github.com/ISSuh/sos/internal/controller/rest/middleware"
	"github.com/ISSuh/sos/internal/factory"
	"github.com/ISSuh/sos/internal/http"
)

func Route(s *http.Server, h *factory.Handlers) {
	s.Use(middleware.ParseParam)

	routes := http.RouteList{
		http.RouteItem{
			URL:     "/v1/{group}/{partition}/{object}",
			Method:  gohttp.MethodGet,
			Handler: h.Downloader.Download,
		},
		http.RouteItem{
			URL:     "/v1/{group}/{partition}/{object}/meta",
			Method:  gohttp.MethodGet,
			Handler: h.Downloader.Download,
		},
		http.RouteItem{
			URL:     "/v1/{group}/{partition}/{object}",
			Method:  gohttp.MethodPost,
			Handler: h.Uploader.Upload,
		},
		http.RouteItem{
			URL:     "/v1/{group}/{partition}/{object}",
			Method:  gohttp.MethodDelete,
			Handler: h.Downloader.Download,
		},
	}

	s.MuxAll(routes)
}
