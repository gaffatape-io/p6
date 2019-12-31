package okrs

import (
	"cloud.google.com/go/firestore"
	"context"
	errs "github.com/gaffatape-io/gopherrs"
	"github.com/gaffatape-io/p6/crud"
)

type Objectives struct {
	S     crud.ObjectiveStore
	RunTx crud.TxRun
}

func isValid(data crud.Objective) bool {	
	return data.Summary != ""	
}

func (o *Objectives) Create(ctx context.Context, data crud.Objective) (crud.Objective, error) {
	var entity crud.Objective

	if !isValid(data) {
		return entity, errs.InvalidArgument(nil)
	}

	if data.ID != "" {
		return entity, errs.InvalidArgumentf(nil, "%+v", data)
	}

	if data.ParentID == "" {
		return o.S.CreateObjective(ctx, nil, data)
	}

	err := o.RunTx(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		_, err := o.S.ReadObjective(ctx, tx, data.ParentID)
		if errs.IsNotFound(err) {
			return errs.FailedPrecondition(err)
		} else if err != nil {
			return err
		}

		entity, err = o.S.CreateObjective(ctx, tx, data)
		return err
	})

	return entity, err
}

func (o *Objectives) Update(ctx context.Context, data crud.Objective) error {	
	if !isValid(data) {
		return errs.InvalidArgument(nil)
	}

	if data.ParentID == "" {
		return o.S.UpdateObjective(ctx, nil, data)
	}

	err := o.RunTx(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		_, err := o.S.ReadObjective(ctx, tx, data.ParentID)
		if errs.IsNotFound(err) {
			return errs.FailedPrecondition(err)
		} else if err != nil {
			return err
		}

		return o.S.UpdateObjective(ctx, tx, data)
	})

	return err
}
