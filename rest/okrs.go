package rest

type Item struct {
	ID          string
	Summary     string
	Description string
}

type HItem struct {
	Item
	ParentID string
}

type KeyResult struct {
	Item
}
