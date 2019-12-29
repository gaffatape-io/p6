package okrs

import (
	"context"
	errs "github.com/gaffatape-io/gopherrs"
	"github.com/gaffatape-io/p6/crud"
	"reflect"
	"testing"
)

func TestObjectiveHandlerPUT(t *testing.T) {
	RunOkrsTest(t, func(ctx context.Context, e *OkrsEnv) {
		sum := e.String("SsSs")
		desc := e.String("DdDd")

		o := crud.Objective{crud.HItem{crud.Item{Summary: sum, Description: desc}, ""}}
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

func TestObjectiveHandlerPUT_missingSummary(t *testing.T) {
	RunOkrsTest(t, func(ctx context.Context, e *OkrsEnv) {
		desc := e.String("DDDD")
		o := crud.Objective{crud.HItem{crud.Item{Description: desc}, ""}}

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

func TestObjectiveHandlerPUT_missingParent(t *testing.T) {
	RunOkrsTest(t, func(ctx context.Context, e *OkrsEnv) {
		sum := e.String("SsSs")
		parentID := e.String("n0_5uch_p@ren1")
		o := crud.Objective{crud.HItem{crud.Item{Summary: sum}, parentID}}

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

// func TestObjectiveHandlerPUT_withParent(t *testing.T) {
// 	RunRestTest(t, func(ctx context.Context, e *RestEnv) {
// 		sum := e.String("SsSs")
// 		o := &Objective{HItem{Item{Summary: sum}, ""}}

// 		resp := checkOK(t, e.roundTripPUT(ctx, "/o", o))
// 		var oResp Objective
// 		err := readJson(resp.Body, &oResp)
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		t.Log(oResp.ID)
// 		sum2 := e.String("Ss2Ss2")
// 		o2 := &Objective{HItem{Item{Summary: sum2}, oResp.ID}}
// 		resp2 := checkOK(t, e.roundTripPUT(ctx, "/o", o2))
// 		var oResp2 Objective
// 		err = readJson(resp2.Body, &oResp2)
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		t.Log(oResp2.ID)
// 		objs := e.Firestore.Collection("objectives")
// 		doc := objs.Doc(oResp2.ID)
// 		_, err = doc.Get(ctx)
// 		t.Log(err)

// 		if err != nil {
// 			t.Fatal()
// 		}
// 	})
// }
