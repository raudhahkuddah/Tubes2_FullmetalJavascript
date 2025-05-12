package models

type SearchNode struct {
	Element  string
	Path     []string
	Children []*SearchNode
}

type TreeNode struct {
	Name           string      `json:"name"`
	NodeDiscovered int         `json:"node_discovered"`
	Children       []*TreeNode `json:"children"`
	Recipe         []string    `json:"recipe,omitempty"`
}

type TreeResult struct {
	Tree         *TreeNode   `json:"tree"`
	Algorithm    string      `json:"algorithm"`
	DurationMs   int64       `json:"duration_ms"`
	VisitedNodes int         `json:"visited_nodes"`
	Recipes      [][]string  `json:"recipes,omitempty"` 
}

type SearchRequest struct {
	Element    string `json:"element"`     
	Algorithm  string `json:"algorithm"`   
	NumResults int    `json:"num_results"`
}

type SearchResult struct {
	RecipePath []string `json:"recipe_path"`
}
