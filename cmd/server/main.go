package main

import (
	"apis/configs"
	"apis/internal/database"
	"apis/internal/dto"
	"apis/internal/entity"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	_, err := configs.LoadConfig(".")
	if err != nil {
		panic("cannot load config: " + err.Error())
	}

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	db.AutoMigrate(&entity.Product{}, &entity.User{})

	// Teste Endpoint
	productDB := database.NewProductDB(db)
	productHandler := NewProductHandler(productDB)
	http.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			productHandler.FindAllProducts(w, r)
		case http.MethodPost:
			productHandler.CreateProduct(w, r)
		default:
			w.Header().Set("Allow", "GET, POST")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/products/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			productHandler.DeleteProduct(w, r)
		case http.MethodGet:
			productHandler.FindProductByID(w, r)
		default:
			w.Header().Set("Allow", "DELETE, GET")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.ListenAndServe(":8080", nil)
}

type ProductHandler struct {
	ProductDB database.ProductInterface
}

func NewProductHandler(productDB database.ProductInterface) *ProductHandler {
	return &ProductHandler{ProductDB: productDB}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product dto.CreateProductInput
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	p, err := entity.NewProduct(product.Name, int(product.Price))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.ProductDB.Create(p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *ProductHandler) FindAllProducts(w http.ResponseWriter, r *http.Request) {
	var page, limit int
	var sort string

	// defaults
	page = 1
	limit = 10

	q := r.URL.Query()
	if ps := q.Get("page"); ps != "" {
		if p, err := strconv.Atoi(ps); err == nil && p > 0 {
			page = p
		}
	}
	if ls := q.Get("limit"); ls != "" {
		if l, err := strconv.Atoi(ls); err == nil && l > 0 {
			limit = l
		}
	}
	sort = q.Get("sort")
	if sort == "" {
		sort = "asc"
	}

	fmt.Print("Fetching products: page ", page, ", limit ", limit, ", sort ", sort, "\n")

	products, err := h.ProductDB.FindAll(page, limit, sort)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Path[len("/products/"):]
	id := vars
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println("Deleting product with ID:", id)
	err := h.ProductDB.Delete(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ProductHandler) FindProductByID(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Path[len("/products/"):]
	id := vars
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println("Fetching product with ID:", id)
	product, err := h.ProductDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}
