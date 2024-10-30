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

	"github.com/ISSuh/sos/internal/infrastructure/transport/rest"
	"github.com/ISSuh/sos/internal/infrastructure/transport/rest/middleware"
	"github.com/ISSuh/sos/pkg/http"
	"github.com/ISSuh/sos/pkg/log"
)

const (
	URLVersion1 = "/v1"

	URLGroup      = "/{" + http.GroupParamName + "}"
	URLPartition  = "/{" + http.PartitionParamName + "}"
	URLObjectPath = "/{" + http.ObjectPathParamName + "}"
	URLObjectID   = "/{" + http.ObjectIDParamName + "}"
	URLMetadata   = "/metadata"

	URLDefault            = URLVersion1 + URLGroup + URLPartition + URLObjectPath
	URLObject             = URLDefault + URLObjectID
	URLObjectMetadata     = URLObject + URLMetadata
	URLObjectMetadataList = URLDefault + URLMetadata
)

func Route(logger log.Logger, s *http.Server, h rest.Explorer) {
	s.Use(middleware.Recover)
	s.Use(middleware.WithLog(logger))
	s.Use(middleware.GenerateRequestID)
	s.Use(middleware.ParseDefaultParam)
	s.Use(middleware.ErrorHandler)

	routes := http.RouteList{
		// Upload
		http.RouteItem{
			URL:         URLDefault,
			Method:      gohttp.MethodPost,
			Handler:     h.Upload,
			Middlewares: []http.MiddlewareFunc{},
		},
		// Download
		http.RouteItem{
			URL:     URLObject,
			Method:  gohttp.MethodGet,
			Handler: h.Download,
			Middlewares: []http.MiddlewareFunc{
				middleware.ParseObjectIDParam,
			},
		},
		// Update
		http.RouteItem{
			URL:     URLObject,
			Method:  gohttp.MethodPut,
			Handler: h.Update,
			Middlewares: []http.MiddlewareFunc{
				middleware.ParseObjectIDParam,
			},
		},
		// Delete
		http.RouteItem{
			URL:     URLObject,
			Method:  gohttp.MethodDelete,
			Handler: h.Delete,
			Middlewares: []http.MiddlewareFunc{
				middleware.ParseObjectIDParam,
			},
		},
		// metadata
		http.RouteItem{
			URL:     URLObjectMetadata,
			Method:  gohttp.MethodGet,
			Handler: h.Find,
			Middlewares: []http.MiddlewareFunc{
				middleware.ParseObjectIDParam,
			},
		},
		// list
		http.RouteItem{
			URL:     URLDefault,
			Method:  gohttp.MethodGet,
			Handler: h.List,
		},
	}

	s.MuxAll(routes)
}
