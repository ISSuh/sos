package factory

import "github.com/ISSuh/sos/internal/controller/rest/handler"

type Handlers struct {
	Uploader   *handler.UploadHandler
	Downloader *handler.DownloadHandler
}

func NewHandlers() (*Handlers, error) {
	h := &Handlers{
		Uploader:   handler.NewUploadHandler(),
		Downloader: handler.NewDownloadHandler(),
	}
	return h, nil
}
