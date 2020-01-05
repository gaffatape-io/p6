package okrs

import (
	"context"
	errs "github.com/gaffatape-io/gopherrs"
	"github.com/gaffatape-io/p6/crud"
	"reflect"
	"testing"
)

func TestObjectivesCreate(t *testing.T) {
	RunOkrsTest(t, func(ctx context.Context, e *OkrsEnv) {
		oe := e.createObjective(ctx, "")

		matches := e.whereObjectives(ctx, "Summary", "==", oe.Summary)
		if len(matches) != 1 {
			t.Fatal()
		}

		want := map[string]interface{}{
			"Summary":     oe.Summary,
			"Description": oe.Description,
			"ParentID":    "",
			"Deleted":     false,
		}

		got := matches[0].Data()
		t.Log(got, want)

		if !reflect.DeepEqual(got, want) {
			t.Fatal()
		}
	})
}

func TestObjectivesCreate_missingSummary(t *testing.T) {
	RunOkrsTest(t, func(ctx context.Context, e *OkrsEnv) {
		desc := e.String("DDDD")
		o := crud.Objective{Description: desc}

		oe, err := e.objs.Create(ctx, o)
		t.Log(oe, err)
		if !errs.IsInvalidArgument(err) {
			t.Fatal()
		}

		matches := e.whereObjectives(ctx, "Description", "==", desc)
		if len(matches) != 0 {
			t.Fatal()
		}
	})
}

func TestObjectivesCreate_withID(t *testing.T) {
	RunOkrsTest(t, func(ctx context.Context, e *OkrsEnv) {
		sum := e.String("s")
		o := crud.Objective{ID: "should-not-be-set", Summary: sum}

		oe, err := e.objs.Create(ctx, o)
		t.Log(oe, err)
		if !errs.IsInvalidArgument(err) {
			t.Fatal()
		}

		matches := e.whereObjectives(ctx, "Summary", "==", sum)
		if len(matches) != 0 {
			t.Fatal()
		}
	})
}

func TestObjectivesCreate_missingParent(t *testing.T) {
	RunOkrsTest(t, func(ctx context.Context, e *OkrsEnv) {
		parentID := e.String("n0_5uch_p@ren1")
		sum := e.String("SsSs")
		o := crud.Objective{Summary: sum, ParentID: parentID}

		oe, err := e.objs.Create(ctx, o)
		t.Log(oe, err)
		if !errs.IsFailedPrecondition(err) {
			t.Fatal()
		}

		matches := e.whereObjectives(ctx, "Summary", "==", sum)
		if len(matches) != 0 {
			t.Fatal()
		}
	})
}

func TestObjectivesCreate_withParent(t *testing.T) {
	RunOkrsTest(t, func(ctx context.Context, e *OkrsEnv) {
		oe := e.createObjective(ctx, "")
		oe2 := e.createObjective(ctx, oe.ID)

		objs := e.Firestore.Collection("objectives")
		doc := objs.Doc(oe2.ID)
		_, err := doc.Get(ctx)
		t.Log(err)

		if err != nil {
			t.Fatal()
		}
	})
}

func TestObjectivesUpdate(t *testing.T) {
	RunOkrsTest(t, func(ctx context.Context, e *OkrsEnv) {
		oe := e.createObjective(ctx, "")

		oe.Summary = e.String("changed summary")
		oe.Description = e.String("new description")
		err := e.objs.Update(ctx, oe)
		t.Log(err)
		if err != nil {
			t.Fatal()
		}

		want := map[string]interface{}{
			"Summary":     oe.Summary,
			"Description": oe.Description,
			"ParentID":    "",
			"Deleted":     false,
		}

		got := e.findObjectiveBySummary(ctx, oe.Summary)
		t.Log(got, want)

		if !reflect.DeepEqual(got, want) {
			t.Fatal()
		}
	})
}

func TestObjectivesUpdate_withParentID(t *testing.T) {
	RunOkrsTest(t, func(ctx context.Context, e *OkrsEnv) {
		p := e.createObjective(ctx, "")
		o := e.createObjective(ctx, p.ID)

		o.Summary = e.String("changed summary")
		err := e.objs.Update(ctx, o)
		t.Log(err)
		if err != nil {
			t.Fatal()
		}

		matches := e.whereObjectives(ctx, "Summary", "==", o.Summary)
		got := matches[0].Data()
		want := map[string]interface{}{
			"Summary":     o.Summary,
			"Description": o.Description,
			"ParentID":    o.ParentID,
			"Deleted":     false,
		}

		if !reflect.DeepEqual(got, want) {
			t.Log(got)
			t.Log(want)
			t.Fatal()
		}
	})
}

func TestDeleteObjective(t *testing.T) {
	RunOkrsTest(t, func(ctx context.Context, e *OkrsEnv) {
		o := e.createObjective(ctx, "")

		err := e.objs.Delete(ctx, o.ID)
		if err != nil {
			t.Fatal(err)
		}

		raw := e.objectiveData(ctx, o.ID)

		if deleted, ok := raw["Deleted"]; !(deleted.(bool)) || !ok {
			t.Fatal(deleted, ok)
		}
	})
}
