package models

// TreeNode mewakili node dalam pohon komposisi elemen
type TreeNode struct {
	Name     string     `json:"name"`
	Children []TreeNode `json:"children"`
}

// TreeResult mewakili hasil pencarian pohon komposisi elemen
type TreeResult struct {
	Tree         *TreeNode `json:"tree"`
	Algorithm    string    `json:"algorithm"`
	DurationMs   int64     `json:"duration_ms"`
	VisitedNodes int       `json:"visited_nodes"`
}

// SearchRequest mewakili permintaan pencarian elemen pada server
type SearchRequest struct {
	Element    string `json:"element"`
	Algorithm  string `json:"algorithm"`
	NumResults int    `json:"num_results"`
}
