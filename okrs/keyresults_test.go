package okrs

import (
	"context"
	"reflect"
	"testing"
)

func TestKeyResultsCreate(t *testing.T) {
	RunOkrsTest(t, func(ctx context.Context, e *OkrsEnv) {
		o := e.createObjective(ctx, "")
		kr := e.createKeyResult(ctx, o.ID)

		raw := e.krData(ctx, kr.ID)
		t.Log(raw)

		want := map[string]interface{}{
			"ObjectiveID": o.ID,
			"Summary":     kr.Summary,
			"Description": "",
			"Deleted":     false,
		}

		if !reflect.DeepEqual(raw, want) {
			t.Fatal()
		}
	})
}

func TestKeyResultsRead(t *testing.T) {
	RunOkrsTest(t, func(ctx context.Context, e *OkrsEnv) {
		o := e.createObjective(ctx, "")
		kr := e.createKeyResult(ctx, o.ID)
		kr2, err := e.krs.Read(ctx, kr.ID)
		t.Log(kr2, err)

		if err != nil {
			t.Fatal()
		}

		if !reflect.DeepEqual(kr2, kr) {
			t.Fatal()
		}
	})
}

func TestKeyResultUpdate(t *testing.T) {
	RunOkrsTest(t, func(ctx context.Context, e *OkrsEnv) {
		o := e.createObjective(ctx, "")
		kr := e.createKeyResult(ctx, o.ID)
		kr.Summary = e.String("a new summary")

		err := e.krs.Update(ctx, kr)
		t.Log(err)
		if err != nil {
			t.Fatal()
		}

		raw := e.krData(ctx, kr.ID)
		t.Log(raw)

		if sum, ok := raw["Summary"]; !ok || sum.(string) != kr.Summary {
			t.Fatal()
		}
	})
}

func TestKeyResultDelete(t *testing.T) {
	RunOkrsTest(t, func(ctx context.Context, e *OkrsEnv) {
		o := e.createObjective(ctx, "")
		kr := e.createKeyResult(ctx, o.ID)

		err := e.krs.Delete(ctx, kr.ID)
		if err != nil {
			t.Fatal(err)
		}

		raw := e.krData(ctx, kr.ID)
		t.Log(raw)

		if deleted, ok := raw["Deleted"]; !ok || !deleted.(bool) {
			t.Fatal()
		}
	})
}
