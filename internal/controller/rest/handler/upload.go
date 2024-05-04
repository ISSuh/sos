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
	"encoding/binary"
	"hash/crc32"
	"io"

	"github.com/ISSuh/sos/internal/logger"
	"github.com/gin-gonic/gin"
)

const (
	CrcSize    = 4
	BodySize   = 4092
	BufferSize = BodySize + CrcSize
)

type UploadHandler struct {
	logger logger.Logger
}

func NewUploadHandler(l logger.Logger) *UploadHandler {
	return &UploadHandler{
		logger: l,
	}
}

func (h *UploadHandler) Upload() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.logger.Debugf("[UploadHandler.Upload]")

		totalSize := 0
		body := c.Request.Body
		for {
			buf := make([]byte, BufferSize)
			n, err := body.Read(buf)
			if err != nil {
				if err == io.EOF {
					h.logger.Debugf("[UploadHandler.Upload] end of file")
					break
				}
				h.logger.Errorf("[UploadHandler.Upload] err : %s", err.Error())
				return
			}

			checksumByte := buf[BodySize:]
			buf = buf[:BodySize]

			checksum := binary.LittleEndian.Uint32(checksumByte)
			crc := crc32.ChecksumIEEE(buf)
			if checksum != crc {
				h.logger.Errorf("[UploadHandler.Upload] checksum : %d/%d", checksum, crc)
				return
			}

			h.logger.Debugf("[UploadHandler.Upload] read n : %d", n)
			totalSize += n
		}

		h.logger.Debugf("[UploadHandler.Upload] total file size : %d", totalSize)
	}
}
