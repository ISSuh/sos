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

func NewExplorer(explorerService service.Explorer) (rest.Explorer, error) {
	switch {
	case validation.IsNil(explorerService):
		return nil, fmt.Errorf("explorer service is nil")
	}

	return &explorer{
		explorerService: explorerService,
	}, nil
}

func (h *explorer) Find(w gohttp.ResponseWriter, r *gohttp.Request) {
	c := r.Context()
	log.FromContext(c).Debugf("[explorer.Find]")

	dto := dto.RequestFromContext(c, http.RequestContextKey)
	log.FromContext(c).Debugf("Request: %+v\n", dto)

	metadata, err := h.explorerService.GetObjectMetadata(c, dto)
	if err != nil {
		log.FromContext(c).Errorf("Find Error: %s\n", err.Error())
		gohttp.Error(w, err.Error(), gohttp.StatusInternalServerError)
		return
	}

	if err := http.Json(w, metadata); err != nil {
		log.FromContext(c).Errorf("Find Error: %s\n", err.Error())
		gohttp.Error(w, err.Error(), gohttp.StatusInternalServerError)
		return
	}
}

func (h *explorer) List(w gohttp.ResponseWriter, r *gohttp.Request) {
	c := r.Context()
	log.FromContext(c).Debugf("[explorer.List]")

	dto := dto.RequestFromContext(c, http.RequestContextKey)
	log.FromContext(c).Debugf("Request: %+v\n", dto)

	metadata, err := h.explorerService.FindObjectMetadataOnPath(c, dto)
	if err != nil {
		log.FromContext(c).Errorf("List Error: %s\n", err.Error())
		gohttp.Error(w, err.Error(), gohttp.StatusInternalServerError)
		return
	}

	if err := http.Json(w, metadata); err != nil {
		log.FromContext(c).Errorf("List Error: %s\n", err.Error())
		gohttp.Error(w, err.Error(), gohttp.StatusInternalServerError)
		return
	}
}

func (h *explorer) Upload(w gohttp.ResponseWriter, r *gohttp.Request) {
	h.chunkedUpload(w, r)
}

func (h *explorer) multipartUpload(w gohttp.ResponseWriter, r *gohttp.Request) {
	c := r.Context()

	// 최대 메모리 사용량 설정 (32MB)
	r.ParseMultipartForm(32 << 20)

	// 파일을 가져옵니다.
	file, handler, err := r.FormFile(http.MultiPartUploadKey)
	if err != nil {
		gohttp.Error(w, err.Error(), gohttp.StatusInternalServerError)
		return
	}
	defer file.Close()

	tempFile := r.MultipartForm.File
	for k, v := range tempFile {
		log.FromContext(c).Debugf("tempFile => Key: %s, Value: %+v\n", k, v)
	}

	tempValue := r.MultipartForm.Value
	for k, v := range tempValue {
		log.FromContext(c).Debugf("tempValue => Key: %s, Value: %+v\n", k, v)
	}

	// 업로드된 파일 정보를 출력합니다.
	log.FromContext(c).Debugf("Uploaded File: %+v\n", handler.Filename)
	log.FromContext(c).Debugf("File Size: %+v\n", handler.Size)
	log.FromContext(c).Debugf("MIME Header: %+v\n", handler.Header)

	// // 파일을 서버에 저장합니다.
	// dst, err := os.Create("./uploads/" + handler.Filename)
	// if err != nil {
	// 	gohttp.Error(w, err.Error(), gohttp.StatusInternalServerError)
	// 	return
	// }
	// defer dst.Close()

	// // 파일을 복사합니다.
	// if _, err := io.Copy(dst, file); err != nil {
	// 	gohttp.Error(w, err.Error(), gohttp.StatusInternalServerError)
	// 	return
	// }

	log.FromContext(c).Debugf("Successfully Uploaded File\n")
}

func (h *explorer) chunkedUpload(w gohttp.ResponseWriter, r *gohttp.Request) {
	c := r.Context()
	log.FromContext(c).Debugf("[explorer.chunkedUpload]")

	dto := dto.RequestFromContext(c, http.RequestContextKey)
	log.FromContext(c).Debugf("Request: %+v\n", dto)
	log.FromContext(c).Debugf("content type: %s\n", r.Header.Get("Content-Type"))

	metadata, err := h.explorerService.Upload(c, dto, r.Body)
	if err != nil {
		log.FromContext(c).Errorf("Upload Error: %s\n", err.Error())
		gohttp.Error(w, err.Error(), gohttp.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if err := http.Json(w, metadata); err != nil {
		log.FromContext(c).Errorf("Upload Error: %s\n", err.Error())
		gohttp.Error(w, err.Error(), gohttp.StatusInternalServerError)
		return
	}
}

func (h *explorer) Update(w gohttp.ResponseWriter, r *gohttp.Request) {
	c := r.Context()
	log.FromContext(c).Debugf("[explorer.Update]")

	dto := dto.RequestFromContext(c, http.RequestContextKey)
	log.FromContext(c).Debugf("Request: %+v\n", dto)
}

func (h *explorer) Download(w gohttp.ResponseWriter, r *gohttp.Request) {
	c := r.Context()

	dto := dto.RequestFromContext(c, http.RequestContextKey)
	log.FromContext(c).Debugf("[explorer.Download]")
	log.FromContext(c).Debugf("Request: %+v\n", dto)

	err := h.explorerService.Download(c, dto, h.headerWriter(w), h.bodyWriter(w))
	if err != nil {
		log.FromContext(c).Errorf("Delete Error: %s\n", err.Error())
		gohttp.Error(w, err.Error(), gohttp.StatusInternalServerError)
		return
	}
}

func (h *explorer) headerWriter(w gohttp.ResponseWriter) http.DownloadHeaderWriter {
	return func(name string, size int) {
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", name))
		// w.Header().Set("Content-Type", "multipart/form-data; boundary=boundary")
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

func (h *explorer) Delete(w gohttp.ResponseWriter, r *gohttp.Request) {
	c := r.Context()
	log.FromContext(c).Debugf("[explorer.chunkedUpload]")

	dto := dto.RequestFromContext(c, http.RequestContextKey)
	err := h.explorerService.Delete(c, dto)
	if err != nil {
		log.FromContext(c).Errorf("Delete Error: %s\n", err.Error())
		gohttp.Error(w, err.Error(), gohttp.StatusInternalServerError)
		return
	}

	http.NoContent(w)
}
