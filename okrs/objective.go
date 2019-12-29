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

func (o *Objectives) Create(ctx context.Context, data crud.Objective) (crud.ObjectiveEntity, error) {
	var entity crud.ObjectiveEntity
	
	if data.Summary == "" {
		return entity, errs.InvalidArgumentf(nil, "summary not set")
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
