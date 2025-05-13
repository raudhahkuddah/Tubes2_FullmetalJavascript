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
func buildRecipeTree(element string, depth int, visited map[string]bool) *models.TreeNode {
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
		child1 := buildRecipeTree(ing1, depth+1, visited)
		child2 := buildRecipeTree(ing2, depth+1, visited)

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

// Fungsi pencarian dengan algoritma BFS (Breadth-First Search)
func BFS(target string, numResults int) (*models.TreeResult, error) {
	start := time.Now()              // Catat waktu mulai
	visited := make(map[string]bool) // Map untuk mencatat elemen yang sudah dikunjungi

	tree := buildRecipeTree(target, 0, visited) // Bangun pohon resep dari target
	if tree == nil {
		return nil, fmt.Errorf("could not build tree for %s", target)
	}

	duration := time.Since(start) // Hitung waktu eksekusi

	// Bungkus hasil dalam TreeResult
	treeResult := &models.TreeResult{
		Tree:         tree,
		Algorithm:    "BFS",
		DurationMs:   duration.Milliseconds(),
		VisitedNodes: len(visited), // Hitung jumlah node yang dikunjungi
	}

	return treeResult, nil
}

// Fungsi pencarian dengan algoritma DFS (Depth-First Search)
func DFS(target string, numResults int) (*models.TreeResult, error) {
	start := time.Now()
	visited := make(map[string]bool)

	tree := buildRecipeTree(target, 0, visited)
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
