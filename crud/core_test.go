package crud

import (
	"cloud.google.com/go/firestore"
	"context"
	"reflect"
	"testing"

	. "github.com/gaffatape-io/p6/test"
)

type CoreEnv struct {
	*Env
	col  *firestore.CollectionRef
	data CoreData
}

type CoreTestFunc func(context.Context, *CoreEnv)

func runCoreEnv(t *testing.T, ct CoreTestFunc) {
	RunTest(t, func(ctx context.Context, e *Env) {
		col := e.Firestore.Collection("core_data")
		data := CoreData{e.String("foo"), 123}
		ct(ctx, &CoreEnv{e, col, data})
	})
}

type CoreData struct {
	Name  string
	Count int
}

type CoreDataEntity struct {
	Entity
	CoreData
}

func (c *CoreEnv) CreateCoreData(ctx context.Context, tx *firestore.Transaction, o CoreData) (CoreDataEntity, error) {
	doc := c.col.NewDoc()
	err := create(ctx, tx, doc, o)
	return CoreDataEntity{Entity{doc.ID}, o}, err
}

func (c *CoreEnv) ReadCoreData(ctx context.Context, tx *firestore.Transaction, id string) (CoreDataEntity, error) {
	doc := c.col.Doc(id)
	o := CoreData{}
	err := read(ctx, tx, doc, &o)
	return CoreDataEntity{Entity{doc.ID}, o}, err
}

func (c *CoreEnv) UpdateCoreData(ctx context.Context, tx *firestore.Transaction, o CoreDataEntity) error {
	doc := c.col.Doc(o.ID)
	return update(ctx, tx, doc, o.CoreData)
}

func (c *CoreEnv) DeleteCoreData(ctx context.Context, tx *firestore.Transaction, o CoreDataEntity) error {
	doc := c.col.Doc(o.ID)
	return delete(ctx, tx, doc)
}

func TestCoreCreate(t *testing.T) {
	runCoreEnv(t, func(ctx context.Context, e *CoreEnv) {
		de, err := e.CreateCoreData(ctx, nil, e.data)
		t.Log(de, err)
		if err != nil {
			t.Fatal()
		}

		if de.ID == "" {
			t.Fatal()
		}
	})
}

func TestCoreCreateTx(t *testing.T) {
	runCoreEnv(t, func(ctx context.Context, e *CoreEnv) {
		e.Firestore.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
			de, err := e.CreateCoreData(ctx, tx, e.data)
			t.Log(de, err)
			if err != nil {
				t.Fatal()
			}

			if de.ID == "" {
				t.Fatal()
			}
			return nil
		})
	})
}

func TestCoreRead(t *testing.T) {
	runCoreEnv(t, func(ctx context.Context, e *CoreEnv) {
		de, err := e.CreateCoreData(ctx, nil, e.data)
		if err != nil {
			t.Fatal(err)
		}

		de2, err := e.ReadCoreData(ctx, nil, de.ID)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(de2, de) {
			t.Fatal(de2, de)
		}
	})
}

func TestCoreReadTx(t *testing.T) {
	runCoreEnv(t, func(ctx context.Context, e *CoreEnv) {
		de, err := e.CreateCoreData(ctx, nil, e.data)
		if err != nil {
			t.Fatal(err)
		}

		e.Firestore.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
			de2, err := e.ReadCoreData(ctx, tx, de.ID)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(de2, de) {
				t.Fatal(de2, de)
			}
			return nil
		})
	})
}

func TestCoreUpdate(t *testing.T) {
	runCoreEnv(t, func(ctx context.Context, e *CoreEnv) {
		de, err := e.CreateCoreData(ctx, nil, e.data)
		if err != nil {
			t.Fatal(err)
		}

		de.CoreData = CoreData{e.String("bar"), 234}
		err = e.UpdateCoreData(ctx, nil, de)
		if err != nil {
			t.Fatal()
		}

		de2, err := e.ReadCoreData(ctx, nil, de.ID)
		if err != nil {
			t.Fatal()
		}

		if !reflect.DeepEqual(de2, de) {
			t.Fatal(de2, de)
		}
	})
}

func TestCoreUpdateTx(t *testing.T) {
	runCoreEnv(t, func(ctx context.Context, e *CoreEnv) {
		de, err := e.CreateCoreData(ctx, nil, e.data)
		if err != nil {
			t.Fatal(err)
		}

		e.Firestore.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
			de.CoreData = CoreData{e.String("bar"), 234}
			err = e.UpdateCoreData(ctx, tx, de)
			if err != nil {
				t.Fatal()
			}
			return nil
		})

		de2, err := e.ReadCoreData(ctx, nil, de.ID)
		if err != nil {
			t.Fatal()
		}

		if !reflect.DeepEqual(de2, de) {
			t.Fatal(de2, de)
		}
	})
}

func TestCoreDelete(t *testing.T) {
	runCoreEnv(t, func(ctx context.Context, e *CoreEnv) {
		de, err := e.CreateCoreData(ctx, nil, e.data)
		if err != nil {
			t.Fatal(err)
		}

		err = e.DeleteCoreData(ctx, nil, de)
		if err != nil {
			t.Fatal(err)
		}

		_, err = e.ReadCoreData(ctx, nil, de.ID)
		if err == nil {
			t.Fatal(err)
		}
	})
}

func TestCoreDeleteTx(t *testing.T) {
	runCoreEnv(t, func(ctx context.Context, e *CoreEnv) {
		de, err := e.CreateCoreData(ctx, nil, e.data)
		if err != nil {
			t.Fatal(err)
		}

		e.Firestore.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
			err = e.DeleteCoreData(ctx, nil, de)
			if err != nil {
				t.Fatal(err)
			}
			return nil
		})

		_, err = e.ReadCoreData(ctx, nil, de.ID)
		if err == nil {
			t.Fatal(err)
		}
	})
}
