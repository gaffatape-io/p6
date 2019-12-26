package rest

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestDecodeRequest(t *testing.T) {
	type X struct {
		Name string
		Cnt  int
	}

	buff, err := encodeJson(&X{"alice", 123})
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, "/foo", buff)
	t.Logf("%+v", req)
	var body X
	err = decodeRequestBody(req, &body)
	t.Log(body, err)
	if err != nil {
		t.Fatal()
	}

	if !reflect.DeepEqual(body, X{"alice", 123}) {
		t.Fatal()
	}
}
