package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gaffatape-io/p6/crud"
	"github.com/gaffatape-io/p6/okrs"
	. "github.com/gaffatape-io/p6/test"
	"k8s.io/klog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	klog.InitFlags(Flags)
	res := m.Run()
	klog.Infof("Test finished: %d", res)
	os.Exit(res)
}

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

func encodeJson(data interface{}) (*bytes.Buffer, error) {
	buff := &bytes.Buffer{}
	enc := json.NewEncoder(buff)
	err := enc.Encode(data)
	return buff, err
}

type RestEnv struct {
	*Env
	mux *http.ServeMux
}

type RestTestFunc func(context.Context, *RestEnv)

func RunRestTest(t *testing.T, rtf RestTestFunc) {
	RunTest(t, func(ctx context.Context, e *Env) {
		store := &crud.Store{e.Firestore}
		mux := NewMux(store, &okrs.Objectives{Objectives:store, RunTx:store.RunTx})
		rtf(ctx, &RestEnv{e, mux})
	})
}

func (e *RestEnv) roundTrip(ctx context.Context, method, target string, body interface{}) *http.Response {
	buff, ok := body.(*bytes.Buffer)
	var err error
	if !ok {
		buff, err = encodeJson(body)
		if err != nil {
			e.Fatal(err)
		}
	}

	req := httptest.NewRequest(method, target, buff)
	e.Logf("RoundTrip: %+v", req)
	resp := httptest.NewRecorder()

	e.mux.ServeHTTP(resp, req)
	return resp.Result()
}

func (e *RestEnv) roundTripPUT(ctx context.Context, target string, body interface{}) *http.Response {
	return e.roundTrip(ctx, http.MethodPut, target, body)
}

// checkMethodAllowed issues http roundtrips for all methods defined in the http package.
// All methods passed as supported should return something other than http.StatusMethodNotAllowed.
func (e *RestEnv) checkMethodAllowed(ctx context.Context, target string, allowed ...string) {
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
		resp := e.roundTrip(ctx, m, target, nil)
		e.Log(m, resp.StatusCode, http.StatusText(resp.StatusCode))

		if allowed && resp.StatusCode == http.StatusMethodNotAllowed {
			e.Fatal()
		} else if !allowed && resp.StatusCode != http.StatusMethodNotAllowed {
			e.Fatal()
		}
	}
}

// checkOK fails the test if the response was not a http.StatusOK.
func checkOK(t *testing.T, resp *http.Response) *http.Response {
	return checkResponse(t, http.StatusOK, resp)
}

// checkResponse fails the test if the response is different than code.
func checkResponse(t *testing.T, code int, resp *http.Response) *http.Response {
	t.Log(code, resp.StatusCode, resp.Status, http.StatusText(resp.StatusCode))
	if resp.StatusCode != code {
		t.Fatalf("wrong code; got:%+v want:%+v", resp.StatusCode, code)
	}
	return resp
}
