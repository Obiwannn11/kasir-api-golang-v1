package services

import (
	"errors"
	"kasir-api-golang-v1/models"
	"kasir-api-golang-v1/repositories"
)

type CategoryService struct {
	catRepo  *repositories.CategoryRepository
	prodRepo *repositories.ProductRepository
}

func NewCategoryService(catRepo *repositories.CategoryRepository, prodRepo *repositories.ProductRepository) *CategoryService {
	return &CategoryService{catRepo: catRepo, prodRepo: prodRepo}
}

func (s *CategoryService) GetAll() ([]models.Category, error) {
	return s.catRepo.GetAll()
}

func (s *CategoryService) GetByID(id int) (*models.Category, error) {
	return s.catRepo.GetByID(id)
}

func (s *CategoryService) Create(category *models.Category) error {
	return s.catRepo.Create(category)
}

func (s *CategoryService) Update(category *models.Category) error {
	return s.catRepo.Update(category)
}

// LOGIC SPESIAL: Safe Delete
func (s *CategoryService) Delete(id int) error {
	// 1. Ambil ID dari "No Category"
	defaultCat, err := s.catRepo.GetByName("No Category")
	if err != nil {
		return errors.New("system error: 'No Category' default data missing")
	}

	// 2. Cegah penghapusan kategori default itu sendiri
	if id == defaultCat.ID {
		return errors.New("cannot delete default category")
	}

	// 3. Pindahkan semua produk di kategori ini ke "No Category"
	err = s.prodRepo.BulkUpdateCategory(id, defaultCat.ID)
	if err != nil {
		return err
	}

	// 4. Hapus kategori
	return s.catRepo.Delete(id)
}