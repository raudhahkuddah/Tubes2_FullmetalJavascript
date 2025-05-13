package search

import (
	"fmt"
	"strings"
	"time"

	"github.com/raudhahkuddah/Tubes2_FullmetalJavascript/backend/models"
	"github.com/raudhahkuddah/Tubes2_FullmetalJavascript/backend/scrapper"
)

// Fungsi untuk menormalisasi nama elemen agar konsisten (huruf kecil semua lalu kapitalisasi awal kata)
func normalisasiNama(nama string) string {
	return strings.Title(strings.ToLower(strings.TrimSpace(nama)))
}

// Fungsi untuk mengecek apakah elemen merupakan elemen dasar
func isBaseElement(element string) bool {
	return scrapper.IsBaseElement(element)
}

// Membuat pohon resep (recipe tree) dari suatu elemen secara rekursif hingga kedalaman tertentu
func buildRecipeTreeDFS(element string, depth int, visited map[string]bool) *models.TreeNode {
	// Batasi kedalaman maksimal rekursi agar tidak terlalu dalam (misal: mencegah stack overflow)
	if depth > 5 {
		return nil
	}

	element = normalisasiNama(element)

	// Hindari elemen yang sudah dikunjungi sebelumnya agar tidak membentuk siklus
	if visited[element] {
		return nil
	}
	visited[element] = true

	// Jika elemen adalah elemen dasar, buat simpul dan kembalikan
	if isBaseElement(element) {
		return &models.TreeNode{
			Name:     element,
			Children: []models.TreeNode{},
		}
	}

	// Ambil data resep dari scrapper
	elementData, err := scrapper.ScrapeElement(element)
	if err != nil || len(elementData.Recipes) == 0 {
		// Jika gagal scraping atau tidak ada resep, kembalikan simpul kosong
		return &models.TreeNode{
			Name:     element,
			Children: []models.TreeNode{},
		}
	}

	// Buat node untuk elemen saat ini
	node := &models.TreeNode{
		Name:     element,
		Children: []models.TreeNode{},
	}

	// Iterasi setiap resep (yang berisi dua bahan)
	for _, recipe := range elementData.Recipes {
		if len(recipe) != 2 {
			continue // Lewati jika format resep tidak sesuai
		}
		// Normalisasi nama bahan
		ing1 := normalisasiNama(recipe[0])
		ing2 := normalisasiNama(recipe[1])

		// Bangun subtree untuk masing-masing bahan
		child1 := buildRecipeTreeDFS(ing1, depth+1, visited)
		child2 := buildRecipeTreeDFS(ing2, depth+1, visited)

		// Buat node kombinasi (nama berupa gabungan dua bahan)
		combinationNode := models.TreeNode{
			Name:     fmt.Sprintf("%s + %s", ing1, ing2),
			Children: []models.TreeNode{},
		}

		// Tambahkan anak-anak ke node kombinasi jika tersedia
		if child1 != nil {
			combinationNode.Children = append(combinationNode.Children, *child1)
		}
		if child2 != nil {
			combinationNode.Children = append(combinationNode.Children, *child2)
		}

		// Tambahkan node kombinasi ke anak dari node utama
		node.Children = append(node.Children, combinationNode)
	}

	return node
}

// BFSNode struktur untuk menyimpan node dalam antrian BFS
type BFSNode struct {
	Element   string
	TreeNode  *models.TreeNode
	Depth     int
	ParentPtr *models.TreeNode // Pointer ke parent node untuk menambahkan anak
}

// buildRecipeTreeBFS menggunakan pendekatan BFS untuk membangun pohon resep
// dengan menampilkan semua resep yang mungkin
func buildRecipeTreeBFS(target string) (*models.TreeNode, int) {
	// Root node for the target element
	root := &models.TreeNode{
		Name:     normalisasiNama(target),
		Children: []models.TreeNode{},
	}

	// Track visited elements to prevent cycles
	visited := make(map[string]bool)
	nodesVisited := 1 // Count root node

	// Get all recipes for the root element
	elementData, err := scrapper.ScrapeElement(target)
	if err != nil || len(elementData.Recipes) == 0 {
		// If no recipes, just return the root node
		return root, nodesVisited
	}

	// Process recipes level by level
	currentLevel := []*models.TreeNode{root}
	nextLevel := []*models.TreeNode{}
	maxDepth := 5 // Limit tree depth

	// BFS loop - process by level
	for depth := 0; depth < maxDepth && len(currentLevel) > 0; depth++ {
		// Process each node at current level
		for _, node := range currentLevel {
			element := node.Name

			// If this is a base element, nothing to expand
			if isBaseElement(element) {
				continue
			}

			// Get all recipes for this element
			elementData, err := scrapper.ScrapeElement(element)
			if err != nil {
				continue
			}

			// Process ALL recipes for this element (this is key)
			for _, recipe := range elementData.Recipes {
				if len(recipe) != 2 {
					continue
				}

				ing1 := normalisasiNama(recipe[0])
				ing2 := normalisasiNama(recipe[1])

				// Create combination node
				combinationNode := models.TreeNode{
					Name:     fmt.Sprintf("%s + %s", ing1, ing2),
					Children: []models.TreeNode{},
				}

				// Create ingredient nodes
				ingredient1Node := models.TreeNode{
					Name:     ing1,
					Children: []models.TreeNode{},
				}

				ingredient2Node := models.TreeNode{
					Name:     ing2,
					Children: []models.TreeNode{},
				}

				// Add ingredient nodes to combination node
				combinationNode.Children = append(combinationNode.Children, ingredient1Node)
				combinationNode.Children = append(combinationNode.Children, ingredient2Node)

				// Add combination to the current element's node
				node.Children = append(node.Children, combinationNode)

				// Count these nodes
				nodesVisited += 3

				// Add non-basic ingredients to next level for further expansion
				// (but avoid cycles)
				if !isBaseElement(ing1) && !visited[ing1] {
					childPtr := &combinationNode.Children[0]
					nextLevel = append(nextLevel, childPtr)
					visited[ing1] = true
				}

				if !isBaseElement(ing2) && !visited[ing2] {
					childPtr := &combinationNode.Children[1]
					nextLevel = append(nextLevel, childPtr)
					visited[ing2] = true
				}
			}
		}

		// Move to next level
		currentLevel = nextLevel
		nextLevel = []*models.TreeNode{}
	}

	return root, nodesVisited
}

// Fungsi pencarian dengan algoritma BFS (Breadth-First Search)
func BFS(target string, numResults int) (*models.TreeResult, error) {
	start := time.Now() // Catat waktu mulai

	tree, visitedCount := buildRecipeTreeBFS(target) // Bangun pohon resep dari target
	if tree == nil {
		return nil, fmt.Errorf("could not build tree for %s", target)
	}

	duration := time.Since(start) // Hitung waktu eksekusi

	// Bungkus hasil dalam TreeResult
	treeResult := &models.TreeResult{
		Tree:         tree,
		Algorithm:    "BFS",
		DurationMs:   duration.Milliseconds(),
		VisitedNodes: visitedCount, // Hitung jumlah node yang dikunjungi
	}

	return treeResult, nil
}

// Fungsi pencarian dengan algoritma DFS (Depth-First Search)
func DFS(target string, numResults int) (*models.TreeResult, error) {
	start := time.Now()
	visited := make(map[string]bool)

	tree := buildRecipeTreeDFS(target, 0, visited)
	if tree == nil {
		return nil, fmt.Errorf("could not build tree for %s", target)
	}

	duration := time.Since(start)

	treeResult := &models.TreeResult{
		Tree:         tree,
		Algorithm:    "DFS",
		DurationMs:   duration.Milliseconds(),
		VisitedNodes: len(visited),
	}

	return treeResult, nil
}

// Fungsi utama untuk melakukan pencarian berdasarkan algoritma (BFS atau DFS)
func Search(request models.SearchRequest) (*models.TreeResult, error) {
	switch strings.ToLower(request.Algorithm) {
	case "bfs":
		return BFS(request.Element, request.NumResults)
	case "dfs":
		return DFS(request.Element, request.NumResults)
	default:
		return nil, fmt.Errorf("invalid algorithm") // Error jika algoritma tidak dikenali
	}
}
