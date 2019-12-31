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
		sum := e.String("SsSs")
		desc := e.String("DdDd")

		o := crud.Objective{Summary: sum, Description: desc}
		oe, err := e.objectives().Create(ctx, o)
		t.Log(oe, err)
		if err != nil {
			t.Fatal()
		}

		objs := e.Firestore.Collection("objectives")
		matches, err := objs.Where("Summary", "==", sum).Documents(ctx).GetAll()
		t.Log(matches, err)
		if err != nil || len(matches) != 1 {
			t.Fatal()
		}

		want := map[string]interface{}{
			"Summary":     sum,
			"Description": desc,
			"ParentID":    "",
		}

		got := matches[0].Data()
		t.Log(got, want)
		
		if !reflect.DeepEqual(matches[0].Data(), want) {
			t.Fatal()
		}
	})
}

func TestObjectivesCreate_missingSummary(t *testing.T) {
	RunOkrsTest(t, func(ctx context.Context, e *OkrsEnv) {
		desc := e.String("DDDD")
		o := crud.Objective{Description: desc}

		oe, err := e.objectives().Create(ctx, o)
		t.Log(oe, err)
		if !errs.IsInvalidArgument(err) {
			t.Fatal()
		}

		objs := e.Firestore.Collection("objectives")
		matches, err := objs.Where("Description", "==", desc).Documents(ctx).GetAll()
		t.Log(matches, err)

		if err != nil || len(matches) != 0 {
			t.Fatal()
		}
	})
}

func TestObjectivesCreate_withID(t *testing.T) {
	RunOkrsTest(t, func(ctx context.Context, e *OkrsEnv) {
		sum := e.String("s")
		o := crud.Objective{ID: "should-not-be-set", Summary:sum}

		oe, err := e.objectives().Create(ctx, o)
		t.Log(oe, err)
		if !errs.IsInvalidArgument(err) {
			t.Fatal()
		}

		objs := e.Firestore.Collection("objectives")
		matches, err := objs.Where("Summary", "==", sum).Documents(ctx).GetAll()
		t.Log(matches, err)

		if err != nil || len(matches) != 0 {
			t.Fatal()
		}
	})
}

func TestObjectivesCreate_missingParent(t *testing.T) {
	RunOkrsTest(t, func(ctx context.Context, e *OkrsEnv) {
		sum := e.String("SsSs")
		parentID := e.String("n0_5uch_p@ren1")
		o := crud.Objective{Summary: sum, ParentID:parentID}

		oe, err := e.objectives().Create(ctx, o)
		t.Log(oe, err)
		if !errs.IsFailedPrecondition(err) {
			t.Fatal()
		}
		
		objs := e.Firestore.Collection("objectives")
		matches, err := objs.Where("Summary", "==", sum).Documents(ctx).GetAll()
		t.Log(matches, err)

		if err != nil || len(matches) != 0 {
			t.Fatal()
		}
	})
}

func TestObjectivesCreate_withParent(t *testing.T) {
	RunOkrsTest(t, func(ctx context.Context, e *OkrsEnv) {
		sum := e.String("SsSs")
		o := crud.Objective{Summary:sum}

		objectives := e.objectives()
		oe, err := objectives.Create(ctx, o)
		t.Log(oe, err)
		if err != nil {
			t.Fatal()
		}	

		sum2 := e.String("Ss2Ss2")
		o2 := crud.Objective{Summary: sum2, ParentID: oe.ID}
		oe2, err := objectives.Create(ctx, o2)
		t.Log(oe2, err)
		if err != nil {
			t.Fatal()
		}		

		objs := e.Firestore.Collection("objectives")
		doc := objs.Doc(oe2.ID)
		_, err = doc.Get(ctx)
		t.Log(err)

		if err != nil {
			t.Fatal()
		}
	})
}
