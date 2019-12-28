package rest

import (
	"google.golang.org/grpc/codes"
	"net/http"	
	errs "github.com/gaffatape-io/gopherrs"
)

func grpc2httpCode(code codes.Code) int {
	status := http.StatusInternalServerError
	switch code {
	case codes.InvalidArgument:
		status = http.StatusBadRequest
	}
	return status
}

func writeStatus(w http.ResponseWriter, err error) {
	code := errs.Code(err)
	status := grpc2httpCode(code)
	http.Error(w, err.Error(), status)
}
