package main

import (
	"context"
	"database/sql"
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type User struct {
	Name  string
	Email string
	Age   int
}

type Product struct {
	ID          int
	Name        string
	Price       int
	Stock       int
	Description string
}

type ProductForm struct {
	Judul   string
	Action  string
	Product Product
	Errors  []string
}

// Data statis untuk latihan Go template pada halaman Home.
var UserData = []User{
	{Name: "Budi Test", Email: "budi@mail.com", Age: 23},
	{Name: "Bagas Test", Email: "bagas@mail.com", Age: 20},
	{Name: "Putri Test", Email: "putri@mail.com", Age: 25},
}

var funcs = template.FuncMap{
	"rupiah": formatRupiah,
}

func parsePage(file string) *template.Template {
	return template.Must(template.New("layout").Funcs(funcs).ParseFiles(
		"template/layout.html",
		"template/header.html",
		"template/sidebar.html",
		"template/footer.html",
		"template/pages/"+file,
	))
}

var (
	tplHome         = parsePage("home.html")
	tplAnnouncement = parsePage("announcement.html")
	tplProduct      = parsePage("product.html")
	tplProductForm  = parsePage("product_form.html")
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("File .env tidak ditemukan; membaca konfigurasi dari environment variable")
	}

	dsn := strings.TrimSpace(os.Getenv("DB_URL"))
	if dsn == "" {
		log.Fatal("DB_URL belum diatur. Salin .env.example menjadi .env lalu isi konfigurasi MySQL.")
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Gagal membuka konfigurasi database: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Gagal terhubung ke database MySQL: %v", err)
	}
	log.Println("Database MySQL berhasil terhubung")

	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("GET /{$}", homeHandler)
	mux.HandleFunc("GET /announcements", announcementHandler)
	mux.HandleFunc("GET /products", productHandler(db))
	mux.HandleFunc("GET /products/form", productFormHandler)
	mux.HandleFunc("POST /products", productSimpanHandler(db))
	mux.HandleFunc("GET /products/{id}/edit", productEditHandler(db))
	mux.HandleFunc("POST /products/{id}/edit", productUpdateHandler(db))
	mux.HandleFunc("POST /products/{id}/hapus", productSoftDeleteHandler(db))

	addr := strings.TrimSpace(os.Getenv("ADDR"))
	if addr == "" {
		addr = ":8000"
	}
	log.Printf("Server berjalan di http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"title": "Home",
		"judul": "Belajar Pemrograman Web",
		"hewan": []string{"cat", "bird", "snake"},
		"users": UserData,
	}
	if err := tplHome.ExecuteTemplate(w, "layout", data); err != nil {
		log.Printf("Gagal menampilkan Home: %v", err)
		http.Error(w, "Terjadi kesalahan saat menampilkan halaman Home.", http.StatusInternalServerError)
	}
}

func announcementHandler(w http.ResponseWriter, r *http.Request) {
	nilai := 90
	data := map[string]any{
		"title":      "Announcement",
		"nilai":      nilai,
		"is_kkm_oke": nilai > 70,
	}
	if err := tplAnnouncement.ExecuteTemplate(w, "layout", data); err != nil {
		log.Printf("Gagal menampilkan Announcement: %v", err)
		http.Error(w, "Terjadi kesalahan saat menampilkan pengumuman.", http.StatusInternalServerError)
	}
}

func productHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		products, err := getAllProduct(r.Context(), db)
		if err != nil {
			log.Printf("Gagal mengambil data produk: %v", err)
			http.Error(w, "Gagal mengambil data produk dari database.", http.StatusInternalServerError)
			return
		}

		data := map[string]any{
			"title":    "Products",
			"products": products,
		}
		if err := tplProduct.ExecuteTemplate(w, "layout", data); err != nil {
			log.Printf("Gagal menampilkan Products: %v", err)
			http.Error(w, "Terjadi kesalahan saat menampilkan produk.", http.StatusInternalServerError)
		}
	}
}

func productFormHandler(w http.ResponseWriter, r *http.Request) {
	productForm := ProductForm{
		Judul:  "Tambah Produk",
		Action: "/products",
	}
	renderProductForm(w, productForm)
}

func productEditHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil || id < 1 {
			http.Error(w, "ID produk tidak valid.", http.StatusBadRequest)
			return
		}

		product, err := getProductByID(r.Context(), db, id)
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		if err != nil {
			log.Printf("Gagal mengambil produk dengan ID %d: %v", id, err)
			http.Error(w, "Gagal mengambil data produk.", http.StatusInternalServerError)
			return
		}

		productForm := ProductForm{
			Judul:   "Edit Produk",
			Action:  "/products/" + strconv.Itoa(id) + "/edit",
			Product: *product,
		}
		renderProductForm(w, productForm)
	}
}

func productSimpanHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Form tidak valid.", http.StatusBadRequest)
			return
		}

		product, validationErrors := productFromForm(r)
		if len(validationErrors) > 0 {
			renderProductForm(w, ProductForm{
				Judul:   "Tambah Produk",
				Action:  "/products",
				Product: product,
				Errors:  validationErrors,
			})
			return
		}

		if _, err := addProduct(r.Context(), db, product); err != nil {
			log.Printf("Gagal menyimpan produk: %v", err)
			http.Error(w, "Gagal menyimpan produk ke database.", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/products", http.StatusSeeOther)
	}
}

func productUpdateHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil || id < 1 {
			http.Error(w, "ID produk tidak valid.", http.StatusBadRequest)
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Form tidak valid.", http.StatusBadRequest)
			return
		}

		product, validationErrors := productFromForm(r)
		product.ID = id
		if len(validationErrors) > 0 {
			renderProductForm(w, ProductForm{
				Judul:   "Edit Produk",
				Action:  "/products/" + strconv.Itoa(id) + "/edit",
				Product: product,
				Errors:  validationErrors,
			})
			return
		}

		if _, err := updateProduct(r.Context(), db, product); err != nil {
			log.Printf("Gagal memperbarui produk ID %d: %v", id, err)
			http.Error(w, "Gagal memperbarui produk di database.", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/products", http.StatusSeeOther)
	}
}

func productSoftDeleteHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil || id < 1 {
			http.Error(w, "ID produk tidak valid.", http.StatusBadRequest)
			return
		}
		err = deleteProduct(r.Context(), db, id)

		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		if err != nil {
			log.Printf("Gagal soft delete: %v", err)
			http.Error(w, "Gagal menghapus produk dari database.", 500)
			return
		}
		http.Redirect(w, r, "/products", http.StatusSeeOther)
	}
}

func productFromForm(r *http.Request) (Product, []string) {
	var validationErrors []string
	name := strings.TrimSpace(r.PostFormValue("name"))
	price, priceErr := strconv.Atoi(strings.TrimSpace(r.PostFormValue("price")))
	stock, stockErr := strconv.Atoi(strings.TrimSpace(r.PostFormValue("stock")))

	if name == "" {
		validationErrors = append(validationErrors, "Nama produk tidak boleh kosong.")
	}
	if priceErr != nil || price <= 0 {
		validationErrors = append(validationErrors, "Harga harus berupa angka lebih dari 0.")
	}
	if stockErr != nil || stock < 0 {
		validationErrors = append(validationErrors, "Stok harus berupa angka 0 atau lebih.")
	}

	return Product{
		Name:        name,
		Price:       price,
		Stock:       stock,
		Description: strings.TrimSpace(r.PostFormValue("description")),
	}, validationErrors
}

func renderProductForm(w http.ResponseWriter, productForm ProductForm) {
	data := map[string]any{
		"title":       "Products",
		"ProductForm": productForm,
	}
	if err := tplProductForm.ExecuteTemplate(w, "layout", data); err != nil {
		log.Printf("Gagal menampilkan form produk: %v", err)
		http.Error(w, "Terjadi kesalahan saat menampilkan form produk.", http.StatusInternalServerError)
	}
}

func getProductByID(ctx context.Context, db *sql.DB, id int) (*Product, error) {
	const query = `SELECT id, nama, harga, stok, COALESCE(keterangan, '') FROM products WHERE id = ? AND deleted_at IS NULL`
	var product Product
	if err := db.QueryRowContext(ctx, query, id).Scan(
		&product.ID, &product.Name, &product.Price, &product.Stock, &product.Description,
	); err != nil {
		return nil, err
	}
	return &product, nil
}

func getAllProduct(ctx context.Context, db *sql.DB) ([]Product, error) {
	const query = `SELECT id, nama, harga, stok, COALESCE(keterangan, '') FROM products WHERE deleted_at IS NULL ORDER BY id`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]Product, 0)
	for rows.Next() {
		var product Product
		if err := rows.Scan(
			&product.ID, &product.Name, &product.Price, &product.Stock, &product.Description,
		); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, rows.Err()
}

func addProduct(ctx context.Context, db *sql.DB, product Product) (int64, error) {
	const query = `INSERT INTO products (nama, harga, stok, keterangan) VALUES (?, ?, ?, ?)`
	result, err := db.ExecContext(ctx, query, product.Name, product.Price, product.Stock, product.Description)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func updateProduct(ctx context.Context, db *sql.DB, product Product) (int64, error) {
	const query = `UPDATE products SET nama = ?, harga = ?, stok = ?, keterangan = ? WHERE id = ? AND deleted_at IS NULL`
	result, err := db.ExecContext(ctx, query, product.Name, product.Price, product.Stock, product.Description, product.ID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func deleteProduct(ctx context.Context, db *sql.DB, id int) error {
	query := `UPDATE products SET deleted_at = NOW() WHERE id = ? AND deleted_at IS NULL`
	result, err := db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if affected == 0 {
		return err
	}

	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func formatRupiah(value int) string {
	sign := ""
	if value < 0 {
		sign = "-"
		value = -value
	}
	digits := strconv.Itoa(value)
	for i := len(digits) - 3; i > 0; i -= 3 {
		digits = digits[:i] + "." + digits[i:]
	}
	return "Rp " + sign + digits
}
