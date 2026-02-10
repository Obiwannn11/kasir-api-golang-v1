package repositories

import (
	"database/sql"
	"errors"
	"kasir-api-golang-v1/models"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// GetAll dengan JOIN dan Search by Name
func (r *ProductRepository) GetAll(nameFilter string) ([]models.Product, error) {
	query := `
		SELECT p.id, p.name, p.price, p.stock, p.category_id, IFNULL(c.name, 'No Category') 
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id`
	
	args := []interface{}{}
	if nameFilter != "" {
		query += " WHERE p.name LIKE ?"
		args = append(args, "%"+nameFilter+"%")
	}
	
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID, &p.CategoryName); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

// GetByID dengan JOIN juga
func (r *ProductRepository) GetByID(id int) (*models.Product, error) {
	query := `
		SELECT p.id, p.name, p.price, p.stock, p.category_id, IFNULL(c.name, 'No Category') 
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.id = ?`

	var p models.Product
	err := r.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID, &p.CategoryName)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) Create(p *models.Product) error {
	// Logic default category ID 1 jika kosong
	if p.CategoryID == 0 {
		p.CategoryID = 1
	}
	result, err := r.db.Exec("INSERT INTO products (name, price, stock, category_id) VALUES (?, ?, ?, ?)", 
		p.Name, p.Price, p.Stock, p.CategoryID)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	p.ID = int(id)
	return nil
}

func (r *ProductRepository) Update(p *models.Product) error {
	query := "UPDATE products SET name = ?, price = ?, stock = ?, category_id = ? WHERE id = ?"
	result, err := r.db.Exec(query, p.Name, p.Price, p.Stock, p.CategoryID, p.ID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("product not found")
	}
	return nil
}

func (r *ProductRepository) Delete(id int) error {
	result, err := r.db.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("product not found")
	}
	return nil
}

// BulkUpdateCategory untuk Safe Delete logic
func (r *ProductRepository) BulkUpdateCategory(oldCatID, newCatID int) error {
	_, err := r.db.Exec("UPDATE products SET category_id = ? WHERE category_id = ?", newCatID, oldCatID)
	return err
}