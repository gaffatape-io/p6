package rest

import (
	"context"
	"net/http"
	"reflect"
	"testing"
)

func TestObjectiveHandlerPUT(t *testing.T) {
	RunRestTest(t, func(ctx context.Context, e *RestEnv) {
		sum := e.String("SsSs")
		desc := e.String("DdDd")
		o := &Objective{HItem{Item{Summary: sum, Description: desc}, ""}}
		resp := checkOK(t, e.roundTripPUT(ctx, "/o", o))
		t.Log(resp)

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
	RunRestTest(t, func(ctx context.Context, e *RestEnv) {
		desc := e.String("DDDD")
		o := &Objective{HItem{Item{Description: desc}, ""}}
		resp := checkResponse(t, http.StatusBadRequest, e.roundTripPUT(ctx, "/o", o))
		t.Log(resp)

		objs := e.Firestore.Collection("objectives")
		matches, err := objs.Where("Description", "==", desc).Documents(ctx).GetAll()
		t.Log(matches, err)

		if err != nil || len(matches) != 0 {
			t.Fatal()
		}
	})
}

func TestObjectiveHandlerPUT_missingParent(t *testing.T) {
	RunRestTest(t, func(ctx context.Context, e *RestEnv) {
		sum := e.String("SsSs")
		parentID := e.String("n0_5uch_p@ren1")
		o := &Objective{HItem{Item{Summary: sum}, parentID}}

		resp := checkResponse(t, http.StatusBadRequest, e.roundTripPUT(ctx, "/o", o))
		t.Log(resp)

		objs := e.Firestore.Collection("objectives")
		matches, err := objs.Where("Summary", "==", sum).Documents(ctx).GetAll()
		t.Log(matches, err)

		if err != nil || len(matches) != 0 {
			t.Fatal()
		}
	})
}

func TestObjectiveHandlerPUT_withParent(t *testing.T) {
	RunRestTest(t, func(ctx context.Context, e *RestEnv) {
		sum := e.String("SsSs")
		o := &Objective{HItem{Item{Summary: sum}, ""}}

		resp := checkOK(t, e.roundTripPUT(ctx, "/o", o))
		var oResp Objective
		err := readJson(resp.Body, &oResp)
		if err != nil {
			t.Fatal(err)
		}

		t.Log(oResp.ID)
		sum2 := e.String("Ss2Ss2")
		o2 := &Objective{HItem{Item{Summary: sum2}, oResp.ID}}
		resp2 := checkOK(t, e.roundTripPUT(ctx, "/o", o2))
		var oResp2 Objective
		err = readJson(resp2.Body, &oResp2)
		if err != nil {
			t.Fatal(err)
		}

		t.Log(oResp2.ID)
		objs := e.Firestore.Collection("objectives")
		doc := objs.Doc(oResp2.ID)
		_, err = doc.Get(ctx)		
		t.Log(err)

		if err != nil {
			t.Fatal()
		}
	})
}
