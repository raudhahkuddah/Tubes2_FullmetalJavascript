package scrapper

import (
	"fmt"
	"net/http"
	"strings"

	// Library untuk melakukan web scraping pada HTML
	"github.com/PuerkitoBio/goquery"
)

// Struct untuk menyimpan data hasil scraping berupa resep-resep elemen
type ElementData struct {
	Recipes [][]string // Setiap resep terdiri dari dua string (dua elemen pembentuk)
}

// Daftar elemen dasar yang tidak bisa dibuat dari elemen lain
var baseElements = map[string]bool{
	"Air":   true,
	"Earth": true,
	"Fire":  true,
	"Water": true,
}

// Fungsi untuk memeriksa apakah sebuah elemen merupakan elemen dasar
func IsBaseElement(element string) bool {
	return baseElements[element]
}

// Fungsi utama untuk melakukan scraping halaman elemen di website Little Alchemy
func ScrapeElement(elementName string) (*ElementData, error) {
	// Buat URL dengan mengganti spasi pada nama elemen menjadi garis bawah
	url := fmt.Sprintf("https://little-alchemy.fandom.com/wiki/%s", strings.ReplaceAll(elementName, " ", "_"))

	// Lakukan HTTP GET ke halaman elemen
	resp, err := http.Get(url)
	if err != nil {
		return nil, err // Jika gagal, kembalikan error
	}
	defer resp.Body.Close() // Tutup response body setelah selesai digunakan

	// Jika status HTTP bukan 200 OK, anggap gagal
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	// Parsing isi halaman HTML menjadi dokumen GoQuery
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err // Kembalikan error jika parsing gagal
	}

	var recipes [][]string // Tempat menyimpan semua resep
	found := false         // Penanda apakah bagian "Recipes" sudah ditemukan

	// Telusuri setiap elemen HTML dalam bagian konten utama artikel
	doc.Find("div.mw-content-ltr.mw-parser-output").Children().EachWithBreak(func(i int, s *goquery.Selection) bool {
		// Jika menemukan heading h2 dengan span ber-ID "Recipes", tandai bagian mulai scraping
		if s.Is("h3") && s.Find("span.mw-headline#Little_Alchemy_2").Length() > 0 {
			found = true
			return true // Lanjutkan ke elemen berikutnya
		}

		// Jika bagian Recipes sudah ditemukan dan heading h2 baru ditemukan, artinya sudah keluar dari bagian resep
		if found {
			if goquery.NodeName(s) == "h2" {
				return false // Hentikan iterasi
			}

			// Temukan semua list item <li> dalam bagian Recipes
			s.Find("li").Each(func(_ int, li *goquery.Selection) {
				var pair []string // Menyimpan dua elemen pembentuk

				// Cari semua tag <a> (link ke elemen) dalam list item tersebut
				li.Find("a").Each(func(_ int, a *goquery.Selection) {
					text := strings.TrimSpace(a.Text()) // Ambil teks dari <a> dan hilangkan spasi
					if text != "" {
						pair = append(pair, text) // Tambahkan ke pair
					}
				})

				// Jika pair berisi dua elemen, maka itu resep yang valid
				if len(pair) == 2 {
					recipes = append(recipes, pair) // Tambahkan ke daftar resep
				}
			})
		}

		return true // Lanjutkan iterasi
	})

	// Kembalikan hasil scraping dalam bentuk struct ElementData
	return &ElementData{Recipes: recipes}, nil
}
