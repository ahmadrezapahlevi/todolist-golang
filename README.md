# Todo List App - Belajar Golang

Aplikasi Todo List sederhana berbasis web menggunakan **Golang**, **MongoDB**, dan **Bootstrap**. Proyek ini dirancang sebagai sarana pembelajaran bagi pemula yang ingin memahami pengembangan web dengan Golang dan integrasi database NoSQL.

---

## Daftar Isi

- [Fitur](#fitur)
- [Teknologi yang Digunakan](#teknologi-yang-digunakan)
- [Persyaratan](#persyaratan)
- [Instalasi & Menjalankan Proyek](#instalasi--menjalankan-proyek)
- [Struktur Proyek](#struktur-proyek)
- [Lisensi](#lisensi)

---

## Fitur

- Tambah todo baru
- Tampilkan daftar todo
- Hapus todo
- Antarmuka pengguna sederhana menggunakan Bootstrap

---

## Teknologi yang Digunakan

- **Backend:** Go (Golang)
- **Database:** MongoDB
- **Frontend:** HTML + Bootstrap
- **Template Engine:** `html/template` bawaan Golang

---

## Persyaratan

Pastikan Anda memiliki:
- Go versi terbaru [Download Go](https://golang.org/dl/)
- MongoDB berjalan di lokal atau URI Atlas [MongoDB Community](https://www.mongodb.com/try/download/community)
- Git (opsional)

---

## Instalasi & Menjalankan Proyek

### 1. Clone repository

```bash
git clone https://github.com/ahmadrezapahlevi/todolist-golang.git
cd todolist-golang

2. Konfigurasi MongoDB

Edit file main.go dan masukkan URI MongoDB Anda:

clientOptions := options.Client().ApplyURI("mongodb://localhost:27017") // ganti dengan URI MongoDB Anda jika perlu

3. Jalankan proyek

go build
go run main.go

4. Akses di browser

Buka browser dan kunjungi: http://localhost:8080


---

Struktur Proyek

todolist-golang/
│
├── main.go               # Entry point aplikasi
├── templates/            # File HTML template
│   ├── index.html
│
├── static/               # File CSS / JS jika ada (opsional)
│
├── go.mod                # Module file Go
└── README.md             # Dokumentasi proyek


---

Lisensi

Proyek ini bersifat open-source dan bebas digunakan untuk keperluan belajar.


---

Dibuat dengan semangat belajar Golang dan pengembangan web oleh Ahmad Reza Pahlevi.

Silakan sesuaikan bagian `clientOptions` di `main.go` dengan URI MongoDB Anda. Jika Anda memerlukan bantuan lebih lanjut atau ingin menambahkan fitur tambahan, jangan ragu untuk bertanya!0

