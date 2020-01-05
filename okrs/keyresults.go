package okrs

import (
	"cloud.google.com/go/firestore"
	"context"
	errs "github.com/gaffatape-io/gopherrs"
	"github.com/gaffatape-io/p6/crud"
)

type KeyResults struct {
	Objectives crud.ObjectiveStore
	KeyResults crud.KeyResultStore
	RunTx      crud.TxRun
}

func (k *KeyResults) Create(ctx context.Context, kr crud.KeyResult) (crud.KeyResult, error) {
	if kr.ObjectiveID == "" {
		return crud.KeyResult{}, errs.InvalidArgumentf(nil, "ObjectiveID not set")
	}

	err := k.RunTx(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		_, err := k.Objectives.ReadObjective(ctx, tx, kr.ObjectiveID)
		if err != nil {
			return err
		}

		oe, err := k.KeyResults.CreateKeyResult(ctx, tx, kr)
		kr.ID = oe.ID
		return err
	})

	return kr, err
}

func (k *KeyResults) Read(ctx context.Context, id string) (crud.KeyResult, error) {
	return k.KeyResults.ReadKeyResult(ctx, nil, id)
}

func (k *KeyResults) Update(ctx context.Context, kr crud.KeyResult) error {
	if kr.ObjectiveID == "" {
		return errs.InvalidArgumentf(nil, "ObjectiveID not set")
	}

	return k.RunTx(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		_, err := k.Objectives.ReadObjective(ctx, tx, kr.ObjectiveID)
		if err != nil {
			return err
		}

		return k.KeyResults.UpdateKeyResult(ctx, tx, kr)
	})
}

func (k *KeyResults) Delete(ctx context.Context, id string) error {
	return k.RunTx(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		kr, err := k.KeyResults.ReadKeyResult(ctx, nil, id)
		if err != nil {
			return err
		}

		kr.Deleted = true
		return k.KeyResults.UpdateKeyResult(ctx, nil, kr)
	})
}
