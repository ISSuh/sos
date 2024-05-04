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
	"github.com/ISSuh/sos/internal/factory"
	"github.com/gin-gonic/gin"
)

const (
	Version1 = "v1"
)

func Route(e *gin.Engine, h *factory.Handlers) error {
	// TODO: need regist middleware

	// regist handler
	v1 := e.Group(Version1)
	{
		api := v1.Group("")
		api.GET("/:group/:partition/:filename", h.Downloader.Download())
		api.PUT("/:group/:partition/:filename", h.Uploader.Upload())
		// api.POST("/:group/:partition/:filename", h.Uploader.Upload())
		// api.DELETE("/:group/:partition/:filename", h.Downloader.Download())

	}

	return nil
}
