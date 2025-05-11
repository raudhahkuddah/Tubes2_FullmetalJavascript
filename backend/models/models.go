package models

type SearchNode struct {
	Element  string
	Path     []string
	Children []*SearchNode
}

type SearchRequest struct {
	Element    string `json:"element"`
	Algorithm  string `json:"algorithm"`
	NumResults int    `json:"num_results"`
}

type SearchResult struct {
	Steps []Step `json:"steps"`
}

type Step struct {
	Result      string   `json:"result"`
	Ingredients []string `json:"ingredients"`
}
