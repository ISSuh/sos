package factory

import (
	"github.com/ISSuh/sos/internal/controller/rest/handler"
	"github.com/ISSuh/sos/internal/logger"
)

type Handlers struct {
	Uploader   *handler.UploadHandler
	Downloader *handler.DownloadHandler
}

func NewHandlers(l logger.Logger) (*Handlers, error) {
	h := &Handlers{
		Uploader:   handler.NewUploadHandler(l),
		Downloader: handler.NewDownloadHandler(l),
	}
	return h, nil
}
