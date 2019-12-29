package rest

import (
	errs "github.com/gaffatape-io/gopherrs"
	"google.golang.org/grpc/codes"
	"net/http"
)

// grpc2httpCode converts a grpc status code to a http status code
func grpc2httpCode(code codes.Code) int {
	status := http.StatusInternalServerError
	switch code {
	case codes.InvalidArgument:
		status = http.StatusBadRequest
	}
	return status
}

// writeStatus writes the error text and status code to the ResponseWriter.
func writeStatus(w http.ResponseWriter, err error) {
	code := errs.Code(err)
	status := grpc2httpCode(code)
	http.Error(w, err.Error(), status)
}

