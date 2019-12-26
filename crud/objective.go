package crud

import (
	"cloud.google.com/go/firestore"
	"context"
)

type Objective struct {
	HItem
}

type ObjectiveEntity struct {	
	Entity
	Objective
}

type ObjectiveStore interface {
	CreateObjective(ctx context.Context, tx *firestore.Transaction, o Objective) (ObjectiveEntity, error)
	ReadObjective(ctx context.Context, tx *firestore.Transaction, id string) (ObjectiveEntity, error)
	UpdateObjective(ctx context.Context, tx *firestore.Transaction, o ObjectiveEntity) error
	DeleteObjective(ctx context.Context, tx *firestore.Transaction, o ObjectiveEntity) error
}

func (s *Store) objectives() *firestore.CollectionRef {
	return s.Client.Collection("objectives")
}

func (s *Store) CreateObjective(ctx context.Context, tx *firestore.Transaction, o Objective) (ObjectiveEntity, error) {
	if o.ParentID != "" && tx == nil {
		panic("CreateObjective with ParentID without transaction")
	}
	
	doc := s.objectives().NewDoc()
	err := create(ctx, tx, doc, o)
	return ObjectiveEntity{Entity{doc.ID}, o}, err
}

func (s *Store) ReadObjective(ctx context.Context, tx *firestore.Transaction, id string) (ObjectiveEntity, error) {	
	doc := s.objectives().Doc(id)
	o := Objective{}
	err := read(ctx, tx, doc, &o)
	return ObjectiveEntity{Entity{doc.ID}, o}, err
}

func (s *Store) UpdateObjective(ctx context.Context, tx *firestore.Transaction, o ObjectiveEntity) error {
	if o.ParentID != "" && tx == nil {
		panic("UpdateObjective with ParentID without transaction")
	}
	
	doc := s.objectives().Doc(o.ID)
	return update(ctx, tx, doc, o.Objective)
}

func (s *Store) DeleteObjective(ctx context.Context, tx *firestore.Transaction, o ObjectiveEntity) error {
	doc := s.objectives().Doc(o.ID)
	return delete(ctx, tx, doc)
}
