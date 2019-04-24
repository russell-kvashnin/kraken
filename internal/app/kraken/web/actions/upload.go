package actions

import (
	"context"
	"fmt"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/model/config"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/model/file"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/web/presenter"
	kerr "github.com/russell-kvashnin/kraken/internal/pkg/error"
	"net/http"
)

// Upload action specific errors
const (
	FormFileErrorCode      = "UPLOAD_FORM_FILE_ERROR"
	UploadContextErrorCode = "REQUEST_CONTEXT_ERROR"
)

// Prepare upload request middleware
type PrepareUploadRequest struct {
	cfg config.NodeConfig
}

// Prepare upload request middleware constructor
func NewPrepareUploadRequest(cfg config.NodeConfig) *PrepareUploadRequest {
	middleware := new(PrepareUploadRequest)
	middleware.cfg = cfg

	return middleware
}

// Prepare upload request middleware
func (middleware *PrepareUploadRequest) Execute(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, header, err := r.FormFile("file")
		if err != nil {
			e := kerr.NewErr(
				kerr.ErrLvlError,
				ErrorDomain,
				FormFileErrorCode,
				fmt.Errorf("request form error"),
				nil)

			_ = presenter.JsonError(w, e, 402)
			return
		}

		tags := r.MultipartForm.Value["tags"]

		cmd := file.UploadFileCommand{
			File:   f,
			Header: header,
			Tags:   tags,
			NodeId: middleware.cfg.Id,
		}

		ctx := context.WithValue(r.Context(), "command", cmd)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Upload action handler
type UploadAction struct {
	handler *file.UploadFileHandler
}

// Constructor
func NewUploadAction(handler *file.UploadFileHandler) *UploadAction {
	action := new(UploadAction)
	action.handler = handler

	return action
}

// Handle request
func (action *UploadAction) Execute(w http.ResponseWriter, r *http.Request) error {
	cmd, ok := r.Context().Value("command").(file.UploadFileCommand)
	if !ok {
		e := kerr.NewErr(
			kerr.ErrLvlError,
			ErrorDomain,
			UploadContextErrorCode,
			fmt.Errorf("bad request context"),
			nil)

		return presenter.JsonError(w, e, 402)
	}

	doc, err := action.handler.Handle(cmd)
	if err != nil {
		return presenter.JsonError(w, err, 500)
	}

	return presenter.JsonResponse(w, doc)
}
