package crud

type Item struct {
	Summary     string
	Description string
}

	type HItem struct {
		Item
	ParentID    string
}

type Entity struct {
	ID string
}
