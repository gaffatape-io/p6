package crud

import (
	"cloud.google.com/go/firestore"
	"context"
)

type Store struct {
	Client *firestore.Client
}

func (s *Store) Close() error {
	return s.Client.Close()
}

type TxFunc func(context.Context, *firestore.Transaction) error
type TxRun func(context.Context, TxFunc) error

func (s *Store) RunTx(ctx context.Context, txf TxFunc) error {
	return s.Client.RunTransaction(ctx, txf)
}

func create(ctx context.Context, tx *firestore.Transaction, doc *firestore.DocumentRef, data interface{}) error {
	if tx != nil {
		return tx.Create(doc, data)
	}

	_, err := doc.Create(ctx, data)
	return err
}

func read(ctx context.Context, tx *firestore.Transaction, doc *firestore.DocumentRef, dest interface{}) error {
	var snap *firestore.DocumentSnapshot
	var err error
	if tx != nil {
		snap, err = tx.Get(doc)
	} else {
		snap, err = doc.Get(ctx)
	}

	if err != nil {
		return err
	}

	return snap.DataTo(dest)
}

func update(ctx context.Context, tx *firestore.Transaction, doc *firestore.DocumentRef, data interface{}) error {
	if tx != nil {
		return tx.Set(doc, data)
	}

	_, err := doc.Set(ctx, data)
	return err
}

func delete(ctx context.Context, tx *firestore.Transaction, doc *firestore.DocumentRef) error {
	if tx != nil {
		return tx.Delete(doc)
	}

	_, err := doc.Delete(ctx)
	return err
}
