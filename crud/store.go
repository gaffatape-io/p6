package crud

import (
	"cloud.google.com/go/firestore"
)

type Store struct {
	client *firestore.Client
}

