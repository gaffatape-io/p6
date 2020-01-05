package crud

import (
	"cloud.google.com/go/firestore"
	"context"
	errs "github.com/gaffatape-io/gopherrs"
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
		err := tx.Create(doc, data)
		errs.WrapError(&err)
		return err
	}

	_, err := doc.Create(ctx, data)
	errs.WrapError(&err)
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
		errs.WrapError(&err)
		return err
	}

	err = snap.DataTo(dest)
	errs.WrapError(&err)
	return err
}

func readAll(ctx context.Context, tx *firestore.Transaction, docs []*firestore.DocumentRef) ([]*firestore.DocumentRef, error) {
	return nil, nil
}

func update(ctx context.Context, tx *firestore.Transaction, doc *firestore.DocumentRef, data interface{}) error {
	if tx != nil {
		err := tx.Set(doc, data)
		errs.WrapError(&err)
		return err
	}

	_, err := doc.Set(ctx, data)
	errs.WrapError(&err)
	return err
}

func delete(ctx context.Context, tx *firestore.Transaction, doc *firestore.DocumentRef) error {
	if tx != nil {
		err := tx.Delete(doc)
		//errs.WrapError(&err)
		return err
	}

	_, err := doc.Delete(ctx)	
	errs.WrapErrorf(&err, "doh")
	errs.Code(nil)	
	return err
}
