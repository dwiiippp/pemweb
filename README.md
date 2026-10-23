# Project Latihan Pemrograman Web — Golang

Project ini menggunakan Go `net/http`, `html/template`, CSS, dan MySQL. Struktur template mengikuti contoh folder yang diberikan.

## Fitur

- Home: contoh data statis, `range`, dan penggunaan template.
- Announcement: contoh kondisi `if` di Go dan HTML template.
- Products: daftar produk dari MySQL.
- Tambah, edit, dan hapus produk (CRUD).
- Validasi input sederhana untuk nama, harga, dan stok.
- Layout bersama: header, sidebar, footer, dan layout utama.

## Persiapan

1. Instal Go 1.22 atau yang lebih baru.
2. Instal dan jalankan MySQL.
3. Buat database dan tabel dengan menjalankan `database.sql`, misalnya melalui MySQL CLI:

   ```sh
   mysql -u root -p < database.sql
   ```

4. Pastikan file `.env` berisi user/password MySQL yang benar. Contoh format `DB_URL`:

   ```env
   DB_URL="root:password@tcp(127.0.0.1:3306)/pemweb?parseTime=true&charset=utf8mb4&loc=Local"
   ADDR=":8000"
   ```

   Ganti `password` dengan password MySQL lokal kamu. Jika MySQL tidak menggunakan port `3306`, sesuaikan portnya.

5. Dari folder project, unduh dependensi dan jalankan server:

   ```sh
   go mod tidy
   go run .
   ```

6. Buka <http://localhost:8000> di browser.

## Route

| Method | URL | Keterangan |
|---|---|---|
| GET | `/` | Halaman Home |
| GET | `/announcements` | Contoh kondisi / pengumuman |
| GET | `/products` | Daftar produk |
| GET | `/products/form` | Form tambah produk |
| POST | `/products` | Menyimpan produk baru |
| GET | `/products/{id}/edit` | Form edit produk |
| POST | `/products/{id}/edit` | Menyimpan perubahan produk |
| POST | `/products/{id}/hapus` | Menghapus produk |

## Struktur folder

```text
pemweb/
├── static/css/             # File CSS
├── template/pages/         # Isi setiap halaman
├── template/layout.html    # Kerangka HTML utama
├── template/header.html
├── template/sidebar.html
├── template/footer.html
├── .env                    # Konfigurasi lokal (jangan commit)
├── .env.example            # Contoh konfigurasi
├── database.sql            # Skema database
├── go.mod
└── main.go
```

## Catatan belajar

- Template di halaman menggunakan `{{define "content"}}` agar isi halaman dapat dimasukkan ke `layout.html`.
- File `app.css` digunakan oleh layout utama. `test.css` disediakan sebagai file latihan tambahan.
- Jika muncul error koneksi database, periksa MySQL sedang berjalan, database `pemweb` sudah dibuat, dan nilai `DB_URL` benar.
