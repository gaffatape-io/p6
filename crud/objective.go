package crud

import (
	"cloud.google.com/go/firestore"
	"context"
)

type Objective struct {
	ID          string `firestore:"-"`
	ParentID    string
	Summary     string
	Description string
	Deleted     bool
}

type ObjectiveStore interface {
	CreateObjective(ctx context.Context, tx *firestore.Transaction, o Objective) (Objective, error)
	ReadObjective(ctx context.Context, tx *firestore.Transaction, id string) (Objective, error)
	UpdateObjective(ctx context.Context, tx *firestore.Transaction, o Objective) error
	DeleteObjective(ctx context.Context, tx *firestore.Transaction, id string) error
}

func (s *Store) objectives() *firestore.CollectionRef {
	return s.Client.Collection("objectives")
}

func (s *Store) CreateObjective(ctx context.Context, tx *firestore.Transaction, o Objective) (Objective, error) {
	if o.ParentID != "" && tx == nil {
		panic("CreateObjective with ParentID without transaction")
	}

	doc := s.objectives().NewDoc()
	err := create(ctx, tx, doc, o)
	o.ID = doc.ID
	return o, err
}

func (s *Store) ReadObjective(ctx context.Context, tx *firestore.Transaction, id string) (Objective, error) {
	doc := s.objectives().Doc(id)
	o := Objective{}
	err := read(ctx, tx, doc, &o)
	o.ID = doc.ID
	return o, err
}

func (s *Store) UpdateObjective(ctx context.Context, tx *firestore.Transaction, o Objective) error {
	if o.ParentID != "" && tx == nil {
		panic("UpdateObjective with ParentID without transaction")
	}

	doc := s.objectives().Doc(o.ID)
	return update(ctx, tx, doc, o)
}

func (s *Store) DeleteObjective(ctx context.Context, tx *firestore.Transaction, id string) error {
	doc := s.objectives().Doc(id)
	return delete(ctx, tx, doc)
}
