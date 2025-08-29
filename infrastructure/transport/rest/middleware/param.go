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

package middleware

import (
	"context"
	gohttp "net/http"
	"strconv"

	"github.com/ISSuh/sos/domain/model/dto"
	"github.com/ISSuh/sos/domain/model/entity"
	"github.com/ISSuh/sos/internal/http"
	"github.com/ISSuh/sos/internal/validation"
)

func ParseDefaultParam(next gohttp.HandlerFunc) gohttp.HandlerFunc {
	return gohttp.HandlerFunc(func(w gohttp.ResponseWriter, r *gohttp.Request) {
		params := http.ParseParm(r)

		group := params[http.GroupParamName]
		if validation.IsEmpty(group) {
			return
		}

		partition := params[http.PartitionParamName]
		if validation.IsEmpty(partition) {
			return
		}

		path := params[http.ObjectPathParamName]
		if validation.IsEmpty(path) {
			return
		}

		versionStr := params[http.VersionName]
		if validation.IsEmpty(path) {
			return
		}

		version, err := strconv.Atoi(versionStr)
		if err != nil {
			version = -1
		}

		req := dto.Request{
			Group:     group,
			Partition: partition,
			Path:      path,
			Version:   version,
		}

		ctx := context.WithValue(r.Context(), http.RequestContextKey, req)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ParseObjectIDParam(next gohttp.HandlerFunc) gohttp.HandlerFunc {
	return gohttp.HandlerFunc(func(w gohttp.ResponseWriter, r *gohttp.Request) {
		params := http.ParseParm(r)

		objectID := params[http.ObjectIDParamName]
		if validation.IsEmpty(objectID) {
			return
		}

		req := dto.RequestFromContext(r.Context(), http.RequestContextKey)

		id, err := strconv.ParseInt(objectID, 10, 64)
		if err != nil {
			return
		}

		req.ObjectID = entity.NewObjectIDFrom(id)

		ctx := context.WithValue(r.Context(), http.RequestContextKey, req)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ParseQueryParam(next gohttp.HandlerFunc) gohttp.HandlerFunc {
	return gohttp.HandlerFunc(func(w gohttp.ResponseWriter, r *gohttp.Request) {
		name := r.URL.Query().Get(http.ObjectName)
		if validation.IsEmpty(name) {
			return
		}

		sizeStr := r.URL.Query().Get(http.ObjectSizeName)
		size, err := strconv.Atoi(sizeStr)
		if err != nil {
			return
		}

		req := dto.RequestFromContext(r.Context(), http.RequestContextKey)
		req.Name = name
		req.Size = size

		ctx := context.WithValue(r.Context(), http.RequestContextKey, req)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
