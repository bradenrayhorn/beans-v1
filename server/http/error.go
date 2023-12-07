package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/bradenrayhorn/beans/server/beans"
)

var codeToHTTPStatus = map[string]int{
	beans.EFORBIDDEN:     http.StatusForbidden,
	beans.EINTERNAL:      http.StatusInternalServerError,
	beans.EINVALID:       http.StatusUnprocessableEntity,
	beans.ENOTFOUND:      http.StatusNotFound,
	beans.EUNAUTHORIZED:  http.StatusUnauthorized,
	beans.EUNPROCESSABLE: http.StatusBadRequest,
}

func Error(w http.ResponseWriter, err error) {
	type errorResponse struct {
		Error string `json:"error"`
		Code  string `json:"code"`
	}

	var code = beans.EINTERNAL
	var msg = beans.ErrorInternal.Error()
	var beansError beans.Error
	if errors.As(err, &beansError) {
		code, msg = beansError.BeansError()
	}

	if code == beans.EINTERNAL {
		slog.Error(err.Error())
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(codeToHTTPStatus[code])
	_ = json.NewEncoder(w).Encode(&errorResponse{Code: code, Error: msg})
}
