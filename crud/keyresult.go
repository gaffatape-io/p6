package crud

import (
	"cloud.google.com/go/firestore"
	"context"
)

type KeyResult struct {
	Item	
}

type KeyResultEntity struct {	
	Entity
	KeyResult
}

type KeyResultStore interface {
	CreateKeyResult(ctx context.Context, tx *firestore.Transaction, o KeyResult) (KeyResultEntity, error)
	ReadKeyResult(ctx context.Context, tx *firestore.Transaction, id string) (KeyResultEntity, error)
	UpdateKeyResult(ctx context.Context, tx *firestore.Transaction, o KeyResultEntity) error
	DeleteKeyResult(ctx context.Context, tx *firestore.Transaction, o KeyResultEntity) error
}

func (s *Store) keyresults() *firestore.CollectionRef {
	return s.Client.Collection("key_results")
}

func (s *Store) CreateKeyResult(ctx context.Context, tx *firestore.Transaction, o KeyResult) (KeyResultEntity, error) {
	doc := s.keyresults().NewDoc()
	err := create(ctx, tx, doc, o)
	return KeyResultEntity{Entity{doc.ID}, o}, err
}

func (s *Store) ReadKeyResult(ctx context.Context, tx *firestore.Transaction, id string) (KeyResultEntity, error) {
	doc := s.keyresults().Doc(id)
	o := KeyResult{}
	err := read(ctx, tx, doc, &o)
	return KeyResultEntity{Entity{doc.ID}, o}, err
}

func (s *Store) UpdateKeyResult(ctx context.Context, tx *firestore.Transaction, o KeyResultEntity) error {
	doc := s.keyresults().Doc(o.ID)
	return update(ctx, tx, doc, o.KeyResult)
}

func (s *Store) DeleteKeyResult(ctx context.Context, tx *firestore.Transaction, o KeyResultEntity) error {
	doc := s.keyresults().Doc(o.ID)
	return delete(ctx, tx, doc)
}
