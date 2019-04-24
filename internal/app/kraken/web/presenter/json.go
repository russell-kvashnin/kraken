package presenter

import (
	"github.com/json-iterator/go"
	kerr "github.com/russell-kvashnin/kraken/internal/pkg/error"
	"net/http"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Upload response
type Response struct {
	IsSuccess bool `json:"isSuccess"`

	Item interface{} `json:"item,omitempty"`

	Error Error `json:"error,omitempty"`
}

// Response error
type Error struct {
	Code    string            `json:"code,omitempty"`
	Message string            `json:"message,omitempty"`
	Details map[string]string `json:"details,omitempty"`
}

// Render success response
func JsonResponse(w http.ResponseWriter, item interface{}) error {
	r := Response{
		IsSuccess: true,
		Item:      item,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	b, _ := json.Marshal(r)

	_, err := w.Write(b)
	if err != nil {
		return err
	}

	return nil
}

// Render error response
func JsonError(w http.ResponseWriter, err error, status int) error {
	var e Error

	switch err.(type) {
	case kerr.Error:
		e = Error{
			Code:    err.(kerr.Error).Code(),
			Message: err.(kerr.Error).Error(),
			Details: err.(kerr.Error).Details(),
		}
	default:
		e = Error{
			Message: err.Error(),
			Details: nil,
		}
	}

	r := Response{
		IsSuccess: false,
		Error:     e,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	b, _ := json.Marshal(r)

	_, err = w.Write(b)
	if err != nil {
		return err
	}

	return nil
}
