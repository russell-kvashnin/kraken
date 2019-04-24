package actions

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/model/file"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/web/presenter"
	kerr "github.com/russell-kvashnin/kraken/internal/pkg/error"
	"net/http"
)

// Download action specific errors
const (
	ValidationErrorCode      = "VALIDATION_ERROR"
	DownloadContextErrorCode = "DOWNLOAD_CONTEXT_ERROR"
)

// Prepare download request middleware
type PrepareDownloadRequest struct {
}

// Prepare download request middleware constructor
func NewPrepareDownloadRequest() *PrepareDownloadRequest {
	return new(PrepareDownloadRequest)
}

// Prepare upload request middleware
func (middleware *PrepareDownloadRequest) Execute(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		shortUrl := chi.URLParam(r, "shortUrl")

		if shortUrl == "" {
			e := kerr.NewErr(
				kerr.ErrLvlError,
				ErrorDomain,
				ValidationErrorCode,
				fmt.Errorf("field 'shortUrl' must be set"),
				nil)

			_ = presenter.JsonError(w, e, 402)
			return
		}

		query := file.DownloadFileQuery{
			ShortUrl: shortUrl,
		}

		ctx := context.WithValue(r.Context(), "query", query)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Download file action
type DownloadAction struct {
	handler *file.DownloadFileQueryHandler
}

// Download file action constructor
func NewDownloadAction(handler *file.DownloadFileQueryHandler) *DownloadAction {
	action := new(DownloadAction)
	action.handler = handler

	return action
}

// Execute file download
func (action *DownloadAction) Execute(w http.ResponseWriter, r *http.Request) error {
	query, ok := r.Context().Value("query").(file.DownloadFileQuery)
	if !ok {
		e := kerr.NewErr(
			kerr.ErrLvlError,
			ErrorDomain,
			DownloadContextErrorCode,
			fmt.Errorf("bad request context"),
			nil)

		return presenter.JsonError(w, e, 402)
	}

	f, meta, err := action.handler.Handle(query)
	if err != nil {
		return presenter.JsonError(w, err, 500)
	}

	return presenter.FileResponse(w, f, meta)
}
