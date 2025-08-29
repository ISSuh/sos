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

	"github.com/ISSuh/sos/domain/model/dto"
	"github.com/ISSuh/sos/domain/service"
	"github.com/ISSuh/sos/infrastructure/transport/rest"
	"github.com/ISSuh/sos/internal/apm"
	"github.com/ISSuh/sos/internal/http"
	"github.com/ISSuh/sos/internal/log"
	"github.com/ISSuh/sos/internal/validation"
)

type explorer struct {
	explorerService service.Explorer
}

func NewExplorer(explorerService service.Explorer) (rest.Explorer, error) {
	switch {
	case validation.IsNil(explorerService):
		return nil, fmt.Errorf("explorer service is nil")
	}

	return &explorer{
		explorerService: explorerService,
	}, nil
}

func (h *explorer) Find() http.Handler {
	return func(w gohttp.ResponseWriter, r *gohttp.Request) {
		c := r.Context()
		log.FromContext(c).Debugf("[explorer.Find]")

		dto := dto.RequestFromContext(c, http.RequestContextKey)
		log.FromContext(c).Debugf("Request: %+v\n", dto)

		span := apm.SpanStart(c, "List", "explorer", nil)
		defer span.End()

		span.Context.SetLabel("group", dto.Group)
		span.Context.SetLabel("partition", dto.Partition)
		span.Context.SetLabel("path", dto.Path)
		span.Context.SetLabel("objectID", dto.ObjectID)

		item, err := h.explorerService.GetObjectMetadata(c, dto)
		if err != nil {
			log.FromContext(c).Errorf("Find Error: %s\n", err.Error())
			gohttp.Error(w, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		if item.Empty() {
			http.NoContent(w)
			return
		}

		if err := http.Json(w, item); err != nil {
			log.FromContext(c).Errorf("Find Error: %s\n", err.Error())
			gohttp.Error(w, err.Error(), gohttp.StatusInternalServerError)
			return
		}
	}
}

func (h *explorer) List() http.Handler {
	return func(w gohttp.ResponseWriter, r *gohttp.Request) {
		c := r.Context()
		log.FromContext(c).Debugf("[explorer.List]")

		log.Debugf(c, "[explorer.List]")

		req := dto.RequestFromContext(c, http.RequestContextKey)
		log.FromContext(c).Debugf("Request: %+v\n", req)

		span := apm.SpanStart(c, "List", "explorer", nil)
		defer span.End()

		span.Context.SetLabel("group", req.Group)
		span.Context.SetLabel("partition", req.Partition)
		span.Context.SetLabel("path", req.Path)

		items, err := h.explorerService.FindObjectMetadataOnPath(c, req)
		if err != nil {
			log.FromContext(c).Errorf("List Error: %s\n", err.Error())
			gohttp.Error(w, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		if err := http.Json(w, items); err != nil {
			log.FromContext(c).Errorf("List Error: %s\n", err.Error())
			gohttp.Error(w, err.Error(), gohttp.StatusInternalServerError)
			return
		}
	}
}

func (h *explorer) Upload() http.Handler {
	return func(w gohttp.ResponseWriter, r *gohttp.Request) {
		c := r.Context()
		log.FromContext(c).Debugf("[explorer.Upload]")

		req := dto.RequestFromContext(c, http.RequestContextKey)
		log.FromContext(c).Debugf("Request: %+v\n", req)
		log.FromContext(c).Debugf("content type: %s\n", r.Header.Get("Content-Type"))

		span := apm.SpanStart(c, "List", "explorer", nil)
		defer span.End()

		if err := r.ParseMultipartForm(32 << 20); err != nil {
			gohttp.Error(w, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		multipartForm := r.MultipartForm

		var items []dto.Item
		for key, fileHeaders := range multipartForm.File {
			for _, fileHeader := range fileHeaders {
				log.FromContext(c).Debugf("[%s]Uploaded File: %+v\n", key, fileHeader.Filename)
				log.FromContext(c).Debugf("[%s]File Size: %+v\n", key, fileHeader.Size)
				log.FromContext(c).Debugf("[%s]MIME Header: %+v\n", key, fileHeader.Header)

				f, err := fileHeader.Open()
				if err != nil {
					gohttp.Error(w, err.Error(), gohttp.StatusInternalServerError)
					return
				}

				req.Name = fileHeader.Filename
				req.Size = int(fileHeader.Size)

				subSpan := apm.SpanStart(c, "List", "explorer", span)

				subSpan.Context.SetLabel("group", req.Group)
				subSpan.Context.SetLabel("partition", req.Partition)
				subSpan.Context.SetLabel("path", req.Path)
				subSpan.Context.SetLabel("name", req.Name)
				subSpan.Context.SetLabel("size", req.Size)

				item, err := h.explorerService.Upload(c, req, f)
				if err != nil {
					log.FromContext(c).Errorf("Upload Error: %s\n", err.Error())
					gohttp.Error(w, err.Error(), gohttp.StatusInternalServerError)
					return
				}

				items = append(items, item)
				f.Close()

				subSpan.End()
			}
		}

		resp := struct {
			Items dto.Items `json:"items"`
		}{
			Items: items,
		}

		log.FromContext(c).Debugf("Successfully Uploaded File\n")
		if err := http.Json(w, resp); err != nil {
			log.FromContext(c).Errorf("List Error: %s\n", err.Error())
			gohttp.Error(w, err.Error(), gohttp.StatusInternalServerError)
			return
		}
	}
}

func (h *explorer) Download(lastVersion bool) http.Handler {
	return func(w gohttp.ResponseWriter, r *gohttp.Request) {
		c := r.Context()

		dto := dto.RequestFromContext(c, http.RequestContextKey)
		log.FromContext(c).Debugf("[explorer.Download]")
		log.FromContext(c).Debugf("Request: %+v\n", dto)

		writer := http.Writer{
			Header: h.headerWriter(w),
			Body:   h.bodyWriter(w),
		}

		err := h.explorerService.Download(c, dto, writer, lastVersion)
		if err != nil {
			log.FromContext(c).Errorf("Delete Error: %s\n", err.Error())
			gohttp.Error(w, err.Error(), gohttp.StatusInternalServerError)
			return
		}
	}
}

func (h *explorer) Delete(deleteVersion bool) http.Handler {
	return func(w gohttp.ResponseWriter, r *gohttp.Request) {
		c := r.Context()
		log.FromContext(c).Debugf("[explorer.Delete]")

		dto := dto.RequestFromContext(c, http.RequestContextKey)
		err := h.explorerService.Delete(c, dto, deleteVersion)
		if err != nil {
			log.FromContext(c).Errorf("Delete Error: %s\n", err.Error())
			gohttp.Error(w, err.Error(), gohttp.StatusBadRequest)
			return
		}

		http.NoContent(w)
	}
}
func (h *explorer) headerWriter(w gohttp.ResponseWriter) http.DownloadHeaderWriter {
	return func(name string, size int) {
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", name))
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", size))
	}
}

func (h *explorer) bodyWriter(w gohttp.ResponseWriter) http.DownloadBodyWriter {
	return func(buffer []byte) error {
		_, err := w.Write(buffer)
		if err != nil {
			return err
		}
		return nil
	}
}
