package crud

import (
	"cloud.google.com/go/firestore"
	"context"
)

type KeyResult struct {
	ID          string `firestore:"-"`
	ObjectiveID string
	Summary     string
	Description string
	Deleted     bool
}

type KeyResultStore interface {
	CreateKeyResult(ctx context.Context, tx *firestore.Transaction, o KeyResult) (KeyResult, error)
	ReadKeyResult(ctx context.Context, tx *firestore.Transaction, id string) (KeyResult, error)
	UpdateKeyResult(ctx context.Context, tx *firestore.Transaction, o KeyResult) error
	DeleteKeyResult(ctx context.Context, tx *firestore.Transaction, id string) error
}

func (s *Store) keyresults() *firestore.CollectionRef {
	return s.Client.Collection("key_results")
}

func (s *Store) CreateKeyResult(ctx context.Context, tx *firestore.Transaction, o KeyResult) (KeyResult, error) {
	doc := s.keyresults().NewDoc()
	err := create(ctx, tx, doc, o)
	o.ID = doc.ID
	return o, err
}

func (s *Store) ReadKeyResult(ctx context.Context, tx *firestore.Transaction, id string) (KeyResult, error) {
	doc := s.keyresults().Doc(id)
	o := KeyResult{}
	err := read(ctx, tx, doc, &o)
	o.ID = doc.ID
	return o, err
}

func (s *Store) UpdateKeyResult(ctx context.Context, tx *firestore.Transaction, o KeyResult) error {
	doc := s.keyresults().Doc(o.ID)
	return update(ctx, tx, doc, o)
}

func (s *Store) DeleteKeyResult(ctx context.Context, tx *firestore.Transaction, id string) error {
	doc := s.keyresults().Doc(id)
	return delete(ctx, tx, doc)
}

