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
	RecipePath []string `json:"recipe_path"`
}

