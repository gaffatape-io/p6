package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"bytes"
	"cloud.google.com/go/firestore"
)

func TestNewMux(t *testing.T) {
	// ctx := context.Background()

	// fs, err := firestore.NewClient(ctx, "dev-p6")
	// t.Log(fs, err)

	// objectives := fs.Collection("objectives")
	// o := &Objective{
	// 	// Results: []KeyResult{
	// 	// 	KeyResult{Name: "worked"},
	// 	// },
	// }

	// doc, wr, err := objectives.Add(ctx, o)
	// fmt.Println(doc, wr, err)

	// d2 := objectives.Doc(doc.ID)
	// snap, err := d2.Get(ctx)
	// t.Log(snap.Data())

	// o2 := fs.Collection("objectives/p0/o1")
	// doc, wr, err = o2.Add(ctx, &Objective{/*Name: "pelle"*/})
	// fmt.Println(">>>>>>>>>>>>>", o2.Path, doc, wr, err)

	// // g.io/objs => objectives filtered per user
	// // g.io/orgs/<id> => root of org with <id>
	// // g.io/objs/<id> =>

	// // the OKR structure is edited by teams but not organized by teans...
	// // teams don't need to match the organizational structure (matrix)
	// // managers want to see their committed work per org structure.
	// // -> manager views is only a reporting view...
	// // -> OKRs are a collaborate editing objective

	// m := NewMux(fs)
	// if m == nil {
	// 	t.Fatal()
	// }
}

func roundTrip(t *testing.T, method, target string, body interface{}) *http.Response {
	ctx := context.Background()
	fs, err := firestore.NewClient(ctx, "dev-p6")
	t.Log(fs, err)

	mux := NewMux(fs)

	buff, ok := body.(*bytes.Buffer)
	if !ok {
		buff = &bytes.Buffer{}
		enc := json.NewEncoder(buff)
		err := enc.Encode(body)
		if err != nil {
			t.Fatal(err)
		}
	}

	req := httptest.NewRequest(method, target, buff)
	resp := httptest.NewRecorder()

	mux.ServeHTTP(resp, req)
	return resp.Result()
}

func roundTripPUT(t *testing.T, target string, body interface{}) *http.Response {
	return roundTrip(t, http.MethodPut, target, body)
}

// checkMethodAllowed issues http roundtrips for all methods defined in the http package.
// All methods passed as supported should return something other than http.StatusMethodNotAllowed.
func checkMethodAllowed(t *testing.T, target string, allowed ...string) {
	expected := map[string]bool{
		http.MethodGet:     false,
		http.MethodHead:    false,
		http.MethodPost:    false,
		http.MethodPut:     false,
		http.MethodPatch:   false,
		http.MethodDelete:  false,
		http.MethodConnect: false,
		http.MethodOptions: false,
		http.MethodTrace:   false,
	}

	for _, a := range allowed {
		expected[a] = true
	}

	for m, allowed := range expected {		
		resp := roundTrip(t, m, target, nil)
		t.Log(m, resp.StatusCode, http.StatusText(resp.StatusCode))
		
		if allowed && resp.StatusCode == http.StatusMethodNotAllowed {
			t.Fatal()
		} else if !allowed && resp.StatusCode != http.StatusMethodNotAllowed {
			t.Fatal()
		}
	}
}

// checkOK fails the test if the response was not a http.StatusOK.
func checkOK(t *testing.T, resp *http.Response) *http.Response {
	return checkResponse(t, http.StatusOK, resp)
}

// checkResponse fails the test if the response is different than code.
func checkResponse(t *testing.T, code int, resp* http.Response) *http.Response {
	t.Log(code, resp.StatusCode, resp.Status, http.StatusText(resp.StatusCode))
	if resp.StatusCode != code {
		t.Fatal()
	}
	return resp
}
