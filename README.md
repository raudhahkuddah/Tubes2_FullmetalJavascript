# FullmetalJavascript
## Tampilan Form Pencarian

![Form Pencarian](./doc/search-form.png)
Aplikasi pencarian *recipe* elemen dalam permainan **Little Alchemy 2** menggunakan algoritma **BFS** dan **DFS**. Backend dibangun dengan **Golang** dan Frontend menggunakan **Next.js/React**  melalui scraping dari situs Little Alchemy 2.

---

## i. Penjelasan Singkat Algoritma DFS dan BFS

### Breadth-First Search (BFS)
Algoritma pencarian dengan pendekatan **melebar terlebih dahulu**. BFS menjamin pencarian jalur **terpendek** dari elemen dasar ke target karena menyelesaikan semua kombinasi di level saat ini sebelum turun ke level berikutnya. Diimplementasikan menggunakan **queue**.

### Depth-First Search (DFS)
Algoritma pencarian dengan pendekatan **menyelam sedalam mungkin** terlebih dahulu sebelum mundur dan mencoba cabang lainnya. DFS tidak menjamin jalur terpendek, namun efisien untuk eksplorasi dalam. Diimplementasikan secara **rekursif**.

Keduanya mendukung pencarian:
- **Single Recipe** (1 path menuju target)
- **Multiple Recipes** (banyak path hingga jumlah maksimal yang diminta user)

Mode multiple recipes mengaktifkan fitur **multithreading** pada backend untuk mempercepat proses pencarian.

---

## ii. Requirement Program dan Instalasi

### Tools & Software
- Golang versi 1.20 atau lebih baru
- Node.js + npm
- Git
- Koneksi internet (untuk scraping data)

---

## iii. Cara Menjalankan Program

### 🔧 Clone Repositori
```bash
git clone https://github.com/raudhahkuddah/Tubes2_FullmetalJavascript
cd Tubes2_FullmetalJavascript
```

### 🔧 Start Frontend
```bash
cd frontend
npm i
npm run dev
```

### 🔧 Start Backend
```bash
cd backend
go run main.go
```

### 🔧 Akses endpoint
Endpoint pencarian dapat diakses di:
```bash
http://localhost:8080/search
```

## 👥 iv. Author

| Nama Lengkap             | NIM       | 
|--------------------------|-----------|
| Raudhah Yahya Kuddah     | 13523132  |
| Jonathan Levi            | 13523150  |
| Benedictus Nelson        | 13122003  |
