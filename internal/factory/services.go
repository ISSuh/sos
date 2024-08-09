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
	"github.com/ISSuh/sos/internal/domain/service"
	"github.com/ISSuh/sos/pkg/logger"
)

type Services struct {
	Uploader       service.Uploader
	Downloader     service.Downloader
	Deleter        service.Deleter
	ObjectMetadata service.ObjectMetadata
}

func NewServices(l logger.Logger, repositoryFacory *Repositories) (*Services, error) {
	uploader, err := service.NewUploader(l)
	if err != nil {
		return nil, err
	}

	downloader, err := service.NewDownloader(l)
	if err != nil {
		return nil, err
	}

	deleter, err := service.NewDeleter(l)
	if err != nil {
		return nil, err
	}

	objectMetadata, err := service.NewObjectMetadata(l)
	if err != nil {
		return nil, err
	}

	return &Services{
		Uploader:       uploader,
		Downloader:     downloader,
		Deleter:        deleter,
		ObjectMetadata: objectMetadata,
	}, nil
}
