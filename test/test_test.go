package test

import (
	"context"
	"testing"
	"strings"
	"strconv"
)

func TestRunTest(t *testing.T) {
	RunTest(t, func(ctx context.Context, e *Env) {
		t.Log(e)
		if ctx == nil {
			t.Fatal()
		}

		if e.T != t {
			t.Fatal()
		}

		if e.Firestore == nil {
			t.Fatal()
		}

		str := e.String("hello world")
		t.Log(str)
		if !strings.HasPrefix(str, e.Name()) {
			t.Fatal(e.Name(), str)
		}

		if !strings.Contains(str, strconv.FormatUint(uint64(e.IID), 10)) {
			t.Fatal()
		}

		if !strings.HasSuffix(str, "hello world") {
			t.Fatal()
		}
	})
}
