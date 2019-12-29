package rest

import (
	"fmt"
	errs "github.com/gaffatape-io/gopherrs"
	"google.golang.org/grpc/codes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGrpc2HttpCode(t *testing.T) {
	tests := []struct {
		code   codes.Code
		status int
	}{
		{codes.InvalidArgument, http.StatusBadRequest},
		{codes.FailedPrecondition, http.StatusPreconditionFailed},
	}

	for _, tc := range tests {
		status := grpc2httpCode(tc.code)
		if status != tc.status {
			t.Fatal(status, tc.status, tc.code)
		}
	}
}

func TestWriteStatus(t *testing.T) {
	tests := []struct {
		err    error
		status int
		txt    string
	}{
		{
			fmt.Errorf("failed"),
			http.StatusInternalServerError,
			"failed\n",
		},
		{
			errs.InvalidArgumentf(nil, "failed"),
			http.StatusBadRequest,
			"desc:failed",
		},
	}

	for _, tc := range tests {
		t.Log(tc)
		rec := httptest.NewRecorder()
		writeStatus(rec, tc.err)

		if rec.Code != tc.status {
			t.Fatal(rec.Code, tc.status)
		}

		body := rec.Body.String()
		t.Log(body)
		if !strings.Contains(body, tc.txt) {
			t.Fatal()
		}
	}
}
