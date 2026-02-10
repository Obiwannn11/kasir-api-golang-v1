package repositories

import (
	"database/sql"
	"errors"
	"kasir-api-golang-v1/models"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetAll() ([]models.Category, error) {
	rows, err := r.db.Query("SELECT id, name FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (r *CategoryRepository) GetByID(id int) (*models.Category, error) {
	var c models.Category
	err := r.db.QueryRow("SELECT id, name FROM categories WHERE id = ?", id).Scan(&c.ID, &c.Name)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// GetByName dipakai untuk mencari ID dari 'No Category'
func (r *CategoryRepository) GetByName(name string) (*models.Category, error) {
	var c models.Category
	err := r.db.QueryRow("SELECT id, name FROM categories WHERE name = ?", name).Scan(&c.ID, &c.Name)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CategoryRepository) Create(category *models.Category) error {
	result, err := r.db.Exec("INSERT INTO categories (name) VALUES (?)", category.Name)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	category.ID = int(id)
	return nil
}

func (r *CategoryRepository) Update(category *models.Category) error {
	query := "UPDATE categories SET name = ? WHERE id = ?"
	result, err := r.db.Exec(query, category.Name, category.ID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("category not found")
	}
	return nil
}

func (r *CategoryRepository) Delete(id int) error {
	result, err := r.db.Exec("DELETE FROM categories WHERE id = ?", id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("category not found")
	}
	return nil
}