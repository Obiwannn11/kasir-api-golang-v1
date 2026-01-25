package main

import (
	"encoding/json"
	"fmt"
	"net/http"	
	"strconv"
	"strings"
)

type Category struct {
	ID 			int		`json:"id"`
	Name 		string	`json:"name"`
	Description string	`json:"description"`
}

var categories = []Category{
	{ID: 1, Name: "Makanan", Description: "Semua Jenis Makanan"},
}

func main() {
	// Route untuk /categories (Semua & Tambah)
	http.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		
		if r.Method == "GET" {
			json.NewEncoder(w).Encode(categories)
		} else if r.Method == "POST" {
			var newCat Category
			json.NewDecoder(r.Body).Decode(&newCat)
			newCat.ID = len(categories) + 1
			categories = append(categories, newCat)
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(newCat)
		}
	})

	// Route untuk /categories/{id} (Detail, Update, Delete)
	http.HandleFunc("/categories/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		idStr := strings.TrimPrefix(r.URL.Path, "/categories/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "ID tidak valid", http.StatusBadRequest)
			return
		}

		for i, c := range categories {
			if c.ID == id {
				if r.Method == "GET" {
					json.NewEncoder(w).Encode(c)
					return
				} else if r.Method == "PUT" {
					var updatedCat Category
					json.NewDecoder(r.Body).Decode(&updatedCat)
					updatedCat.ID = id
					categories[i] = updatedCat
					json.NewEncoder(w).Encode(updatedCat)
					return
				} else if r.Method == "DELETE" {
					categories = append(categories[:i], categories[i+1:]...)
					json.NewEncoder(w).Encode(map[string]string{"message": "Berhasil hapus"})
					return
				}
			}
		}
		http.Error(w, "Kategori tidak ditemukan", http.StatusNotFound)
	})

	fmt.Println("Server jalan di http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}