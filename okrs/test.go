package okrs

import (
	"context"
	"github.com/gaffatape-io/p6/crud"
	. "github.com/gaffatape-io/p6/test"
	"testing"
)

type OkrsEnv struct {
	*Env
	store *crud.Store
}

func (e *OkrsEnv) objectives() *Objectives {
	return &Objectives{e.store, e.store.RunTx}
}

func RunOkrsTest(t *testing.T, tf func(ctx context.Context, e *OkrsEnv)) {
	RunTest(t, func(ctx context.Context, e *Env) {
		store := &crud.Store{e.Firestore}
		tf(ctx, &OkrsEnv{e, store})
	})
}
