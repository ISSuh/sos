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

const (
	CrcSize    = 4
	BodySize   = 4092
	BufferSize = BodySize + CrcSize
)

type uploader struct {
	logger log.Logger

	uploadService service.Uploader
}

func NewUploader(l log.Logger, uploadService service.Uploader) (rest.Uploader, error) {
	switch {
	case validation.IsNil(l):
		return nil, fmt.Errorf("logger is nil")
	case validation.IsNil(uploadService):
		return nil, fmt.Errorf("upload service is nil")
	}

	return &uploader{
		logger:        l,
		uploadService: uploadService,
	}, nil
}

func (h *uploader) Upload(w gohttp.ResponseWriter, r *gohttp.Request) {
	h.chunkedUpload(w, r)
}

func (h *uploader) multipartUpload(w gohttp.ResponseWriter, r *gohttp.Request) {
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

func (h *uploader) chunkedUpload(w gohttp.ResponseWriter, r *gohttp.Request) {
	c := r.Context()
	dto := dto.RequestFromContext(c, http.RequestContextKey)
	log.FromContext(c).Debugf("[uploader.chunkedUpload]")
	log.FromContext(c).Debugf("Request: %+v\n", dto)
	log.FromContext(c).Debugf("content type: %s\n", r.Header.Get("Content-Type"))

	err := h.uploadService.Upload(c, dto, r.Body)
	if err != nil {
		gohttp.Error(w, err.Error(), gohttp.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	w.WriteHeader(gohttp.StatusOK)
	w.Write([]byte("Chunk uploaded successfully"))
}

func (h *uploader) Update(w gohttp.ResponseWriter, r *gohttp.Request) {
}
