package rest

type Item struct {
	Summary string
	Description string
}

type KeyResult struct {
	Item
}

type Objective struct {
	Item
	KeyResults  []KeyResult
}
