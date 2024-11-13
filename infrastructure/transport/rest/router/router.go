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

	"github.com/ISSuh/sos/infrastructure/transport/rest"
	"github.com/ISSuh/sos/infrastructure/transport/rest/middleware"
	"github.com/ISSuh/sos/internal/http"
	"github.com/ISSuh/sos/internal/log"
)

const (
	URLVersion1 = "/v1"

	URLGroup      = "/{" + http.GroupParamName + "}"
	URLPartition  = "/{" + http.PartitionParamName + "}"
	URLObjectPath = "/{" + http.ObjectPathParamName + "}"
	URLObjectID   = "/{" + http.ObjectIDParamName + "}"
	URLMetadata   = "/metadata"
	URLVersion    = "/version"
	URLVersionNum = "/{" + http.VersionName + "}"

	URLDefault        = URLVersion1 + URLGroup + URLPartition + URLObjectPath
	URLObject         = URLDefault + URLObjectID
	URLObjectMetadata = URLObject + URLMetadata
	URLObjectVersion  = URLObject + URLVersion + URLVersionNum
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
			Method:      gohttp.MethodPut,
			Handler:     h.Upload(),
			Middlewares: []http.MiddlewareFunc{},
		},
		// Download latest version
		http.RouteItem{
			URL:     URLObject,
			Method:  gohttp.MethodGet,
			Handler: h.Download(true),
			Middlewares: []http.MiddlewareFunc{
				middleware.ParseObjectIDParam,
			},
		},
		// Download specific version
		http.RouteItem{
			URL:     URLObjectVersion,
			Method:  gohttp.MethodGet,
			Handler: h.Download(false),
			Middlewares: []http.MiddlewareFunc{
				middleware.ParseObjectIDParam,
			},
		},
		// Delete
		http.RouteItem{
			URL:     URLObject,
			Method:  gohttp.MethodDelete,
			Handler: h.Delete(false),
			Middlewares: []http.MiddlewareFunc{
				middleware.ParseObjectIDParam,
			},
		},
		// Delete specific version
		http.RouteItem{
			URL:     URLObjectVersion,
			Method:  gohttp.MethodDelete,
			Handler: h.Delete(true),
			Middlewares: []http.MiddlewareFunc{
				middleware.ParseObjectIDParam,
			},
		},
		// metadata
		http.RouteItem{
			URL:     URLObjectMetadata,
			Method:  gohttp.MethodGet,
			Handler: h.Find(),
			Middlewares: []http.MiddlewareFunc{
				middleware.ParseObjectIDParam,
			},
		},
		// list
		http.RouteItem{
			URL:     URLDefault,
			Method:  gohttp.MethodGet,
			Handler: h.List(),
		},
	}

	s.MuxAll(routes)
}
