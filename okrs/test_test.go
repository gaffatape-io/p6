package okrs

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/gaffatape-io/p6/crud"
	. "github.com/gaffatape-io/p6/test"
	"testing"
)

type OkrsEnv struct {
	*Env
	*testing.T
	store      *crud.Store
	objsColRef *firestore.CollectionRef
	objs       *Objectives
	krsColRef  *firestore.CollectionRef
	krs        *KeyResults
}

func (e *OkrsEnv) createObjective(ctx context.Context, parentID string) crud.Objective {
	o := crud.Objective{Summary: e.String("sSsS"), ParentID: parentID}
	oe, err := e.objs.Create(ctx, o)
	e.Log("created objective:", oe, err)
	if err != nil {
		e.Fatal()
	}
	return oe
}

func (e *OkrsEnv) whereObjectives(ctx context.Context, path, op string, value interface{}) []*firestore.DocumentSnapshot {
	matches, err := e.objsColRef.Where(path, op, value).Documents(ctx).GetAll()
	if err != nil {
		e.Fatal()
	}
	return matches
}

func (e *OkrsEnv) findObjectiveBySummary(ctx context.Context, summary string) map[string]interface{} {
	matches := e.whereObjectives(ctx, "Summary", "==", summary)
	e.Log(matches)
	if len(matches) != 1 {
		e.Fatal()
	}

	return matches[0].Data()
}

func (e *OkrsEnv) createKeyResult(ctx context.Context, objectiveID string) crud.KeyResult {
	kr := crud.KeyResult{ObjectiveID: objectiveID, Summary: e.String("sss")}
	kre, err := e.krs.Create(ctx, kr)
	e.Log("created key-result:", kre, err)
	if err != nil {
		e.Fatal()
	}
	return kre
}

func (e *OkrsEnv) krSnapshot(ctx context.Context, id string) *firestore.DocumentSnapshot {
	return e.snapshotByID(ctx, e.krsColRef, id)
}

func (e *OkrsEnv) krData(ctx context.Context, id string) map[string]interface{} {
	return e.krSnapshot(ctx, id).Data()
}

func (e *OkrsEnv) objectiveSnapshot(ctx context.Context, id string) *firestore.DocumentSnapshot {
	return e.snapshotByID(ctx, e.objsColRef, id)
}

func (e *OkrsEnv) objectiveData(ctx context.Context, id string) map[string]interface{} {
	return e.objectiveSnapshot(ctx, id).Data()
}

func (e *OkrsEnv) snapshotByID(ctx context.Context, colRef *firestore.CollectionRef, id string) *firestore.DocumentSnapshot {
	if id == "" {
		e.Fatal("invalid document ID; not set")
	}

	doc := colRef.Doc(id)
	snap, err := doc.Get(ctx)
	if err != nil {
		e.Fatal("failed to get snapshot; err:", err)
	}

	return snap
}

func RunOkrsTest(t *testing.T, tf func(ctx context.Context, e *OkrsEnv)) {
	RunTest(t, func(ctx context.Context, e *Env) {
		store := &crud.Store{e.Firestore}
		
		objsColRef := e.Firestore.Collection("objectives")
		objs := &Objectives{store, store.RunTx}
		krsColRef := e.Firestore.Collection("key_results")
		krs := &KeyResults{store, store, store.RunTx}

		tf(ctx, &OkrsEnv{e, t, store, objsColRef, objs, krsColRef, krs})
	})
}
