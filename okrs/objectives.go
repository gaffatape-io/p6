package okrs

import (
	"cloud.google.com/go/firestore"
	"context"
	errs "github.com/gaffatape-io/gopherrs"
	"github.com/gaffatape-io/p6/crud"
)

type Objectives struct {
	Objectives crud.ObjectiveStore
	RunTx      crud.TxRun
}

func isValid(data crud.Objective) bool {
	return data.Summary != ""
}

func (o *Objectives) Create(ctx context.Context, data crud.Objective) (crud.Objective, error) {
	var entity crud.Objective

	if data.ID != "" {
		return entity, errs.InvalidArgument(nil)
	}

	if !isValid(data) {
		return entity, errs.InvalidArgument(nil)
	}

	if data.ParentID == "" {
		return o.Objectives.CreateObjective(ctx, nil, data)
	}

	err := o.RunTx(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		_, err := o.Objectives.ReadObjective(ctx, tx, data.ParentID)
		if errs.IsNotFound(err) {
			return errs.FailedPrecondition(err)
		} else if err != nil {
			return err
		}

		entity, err = o.Objectives.CreateObjective(ctx, tx, data)
		return err
	})

	return entity, err
}

func (o *Objectives) Read(ctx context.Context, id string) (crud.Objective, error) {
	return o.Objectives.ReadObjective(ctx, nil, id)
}

func (o *Objectives) Update(ctx context.Context, data crud.Objective) error {
	if !isValid(data) {
		return errs.InvalidArgument(nil)
	}

	if data.ParentID == "" {
		return o.Objectives.UpdateObjective(ctx, nil, data)
	}

	err := o.RunTx(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		_, err := o.Objectives.ReadObjective(ctx, tx, data.ParentID)
		if errs.IsNotFound(err) {
			return errs.FailedPrecondition(err)
		} else if err != nil {
			return err
		}

		return o.Objectives.UpdateObjective(ctx, tx, data)
	})

	return err
}

func (o *Objectives) Delete(ctx context.Context, id string) error {
	return o.RunTx(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		oe, err := o.Objectives.ReadObjective(ctx, tx, id)
		if err != nil {
			return err
		}

		oe.Deleted = true
		err = o.Objectives.UpdateObjective(ctx, tx, oe)
		return err
	})
}
