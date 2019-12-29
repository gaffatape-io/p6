package test

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"k8s.io/klog"
	"math/rand"
	"testing"
	"time"
)

const max63n = 1<<63 - 1

type Env struct {
	*testing.T
	Firestore *firestore.Client
	T0        time.Time
	IID       int64
	Rnd       *rand.Rand
}

func (e *Env) Name() string {
	return e.T.Name()
}

func (e *Env) String(txt string) string {
	return fmt.Sprintf("%s-%d-%d-%s", e.Name(), e.IID, e.T0.Nanosecond(), txt)
}

func Firestore(ctx context.Context, t *testing.T) *firestore.Client {
	fs, err := firestore.NewClient(ctx, "dev-p6")
	if err != nil {
		t.Fatal(err)
	}
	return fs
}

type TestFunc func(context.Context, *Env)

func RunTest(t *testing.T, f TestFunc) {
	ctx := context.Background()
	fs := Firestore(ctx, t)
	t0 := time.Now()
	rnd := rand.New(rand.NewSource(t0.Unix()))
	env := &Env{t, fs, t0, rnd.Int63n(max63n), rnd}
	klog.Info(">>>", t.Name(), " start")
	defer klog.Info("<<<", t.Name(), " end")
	f(ctx, env)
}
