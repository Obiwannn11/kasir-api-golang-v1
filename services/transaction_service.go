package services

import (
	"kasir-api-golang-v1/models"
	"kasir-api-golang-v1/repositories"
)

type TransactionService struct {
	repo *repositories.TransactionRepository
}

func NewTransactionService(repo *repositories.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) Checkout(items []models.CheckoutItem) (*models.Transaction, error) {
	return s.repo.CreateTransaction(items)
}

// GetTodayReport untuk sales summary hari ini
func (s *TransactionService) GetTodayReport() (map[string]interface{}, error) {
	totalRevenue, totalTransaksi, err := s.repo.GetSalesToday()
	if err != nil {
		return nil, err
	}

	productName, qtySold, err := s.repo.GetTopProductToday()
	if err != nil {
		return nil, err
	}

	report := map[string]interface{}{
		"total_revenue":   totalRevenue,
		"total_transaksi": totalTransaksi,
	}

	if productName != "" {
		report["produk_terlaris"] = map[string]interface{}{
			"nama":        productName,
			"qty_terjual": qtySold,
		}
	} else {
		report["produk_terlaris"] = nil
	}

	return report, nil
}

// GetRangeReport untuk sales summary dengan date range
func (s *TransactionService) GetRangeReport(startDate, endDate string) (map[string]interface{}, error) {
	totalRevenue, totalTransaksi, err := s.repo.GetSalesInRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	productName, qtySold, err := s.repo.GetTopProductInRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	report := map[string]interface{}{
		"total_revenue":   totalRevenue,
		"total_transaksi": totalTransaksi,
		"start_date":      startDate,
		"end_date":        endDate,
	}

	if productName != "" {
		report["produk_terlaris"] = map[string]interface{}{
			"nama":        productName,
			"qty_terjual": qtySold,
		}
	} else {
		report["produk_terlaris"] = nil
	}

	return report, nil
}
