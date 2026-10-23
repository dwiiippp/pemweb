# Project Latihan Pemrograman Web dengan Golang

Project ini dibuat untuk mempelajari dasar-dasar pengembangan aplikasi web menggunakan **Go (Golang)**, **HTML Template**, **CSS**, dan **MySQL**. Aplikasi menggunakan paket `net/http` untuk menangani HTTP request dan `html/template` untuk menampilkan halaman.

## Fitur

- **Home:** menampilkan contoh data statis dan perulangan (`range`) pada template.
- **Announcements:** menampilkan contoh percabangan (`if`) pada Go dan template HTML.
- **Products:** menampilkan data produk yang tersimpan di MySQL.
- **CRUD produk:** menambah, mengubah, dan menghapus produk.
- **Validasi input:** pemeriksaan sederhana untuk nama, harga, dan stok produk.
- **Layout bersama:** menggunakan header, sidebar, footer, dan layout utama.

> **Catatan:** Dokumentasi fitur dan route di bawah mengikuti README awal. Cocokkan kembali dengan implementasi di `main.go`, terutama jika fitur *soft delete* sudah ditambahkan.

## Teknologi yang Digunakan

- Go 1.22 atau lebih baru
- Go `net/http` dan `html/template`
- MySQL
- HTML dan CSS

## Struktur Project

```text
pemweb/
├── static/
│   └── css/
│       ├── app.css
│       └── test.css
├── template/
│   ├── pages/
│   ├── footer.html
│   ├── header.html
│   ├── layout.html
│   └── sidebar.html
├── .env                # Konfigurasi lokal (jangan upload ke GitHub)
├── .env.example        # Contoh konfigurasi tanpa password asli
├── .gitignore
├── go.mod
├── go.sum
├── main.go
└── README.md
```

## Cara Menjalankan Project

### 1. Persiapan

Pastikan **Go 1.22+** dan **MySQL** sudah terinstal. Kamu dapat menggunakan XAMPP untuk menjalankan MySQL dan membuka phpMyAdmin.

### 2. Menyiapkan database melalui phpMyAdmin

1. Jalankan layanan MySQL (misalnya melalui XAMPP).
2. Buka `http://localhost/phpmyadmin` di browser.
3. Buat database baru dengan nama **`pemweb`**.
4. Siapkan tabel yang digunakan aplikasi sesuai struktur/query pada `main.go`.

**Penting:** README sebelumnya menyebut file `database.sql`, tetapi file itu tidak terlihat di struktur project yang diberikan. Jika kamu memiliki file SQL tersebut, kamu dapat mengimpornya melalui tab **Import** pada phpMyAdmin. Jika belum ada, skema tabel harus disiapkan terlebih dahulu sebelum fitur produk dapat digunakan.

### 3. Mengatur konfigurasi `.env`

Salin `.env.example` menjadi `.env` (jika belum ada), lalu sesuaikan dengan konfigurasi MySQL lokal. Contoh:

```env
DB_URL="root:password_mysql@tcp(127.0.0.1:3306)/pemweb?parseTime=true&charset=utf8mb4&loc=Local"
ADDR=":8000"
```

Ganti `password_mysql` dengan password MySQL lokal. Sesuaikan juga host, port, dan nama database jika berbeda.

**Jangan upload `.env` ke GitHub.** Pastikan `.gitignore` berisi setidaknya:

```gitignore
.env
.env.local
.env.production
```

Simpan hanya contoh konfigurasi tanpa kredensial pribadi di `.env.example`.

### 4. Menjalankan aplikasi

Buka terminal di folder utama project, kemudian jalankan:

```bash
go mod tidy
go run .
```

Jika server berhasil dijalankan dengan `ADDR=":8000"`, buka:

**http://localhost:8000**

## Daftar Route

| Method | URL | Keterangan |
| --- | --- | --- |
| `GET` | `/` | Menampilkan halaman Home |
| `GET` | `/announcements` | Menampilkan halaman pengumuman |
| `GET` | `/products` | Menampilkan daftar produk |
| `GET` | `/products/form` | Menampilkan form tambah produk |
| `POST` | `/products` | Menyimpan produk baru |
| `GET` | `/products/{id}/edit` | Menampilkan form edit produk |
| `POST` | `/products/{id}/edit` | Menyimpan perubahan produk |
| `POST` | `/products/{id}/hapus` | Menghapus produk |

## Catatan Pembelajaran

- `html/template` digunakan untuk merender HTML dari server.
- Template halaman menggunakan `{{define "content"}}` sesuai pola yang dijelaskan dalam README awal.
- `app.css` digunakan untuk tampilan utama; `test.css` merupakan file latihan tambahan.
- Jika koneksi database gagal, periksa status MySQL, keberadaan database `pemweb`, dan nilai `DB_URL`.
- Jika fitur *soft delete* sudah diterapkan, perbarui deskripsi route hapus agar menjelaskan bahwa data ditandai terhapus, bukan dihapus secara permanen.

## Catatan untuk GitHub

File `main.go`, `go.mod`, `go.sum`, `template/`, `static/`, `.gitignore`, dan `.env.example` dapat disimpan dalam repository. File `.env` harus tetap berada di komputer lokal dan diabaikan oleh Git.
