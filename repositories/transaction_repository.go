package repositories

import (
	"database/sql"
	"fmt"
	"kasir-api-golang-v1/models"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repo *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	totalAmount := 0
	details := make([]models.TransactionDetail, 0)

	for _, item := range items {
		var productPrice, stock int
		var productName string

		err := tx.QueryRow("SELECT name, price, stock FROM products WHERE id = ?", item.ProductID).Scan(&productName, &productPrice, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}
		if err != nil {
			return nil, err
		}

		// Cek apakah stock cukup
		if stock < item.Quantity {
			return nil, fmt.Errorf("insufficient stock for product %s (available: %d, requested: %d)", productName, stock, item.Quantity)
		}

		subtotal := productPrice * item.Quantity
		totalAmount += subtotal

		_, err = tx.Exec("UPDATE products SET stock = stock - ? WHERE id = ?", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		details = append(details, models.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	result, err := tx.Exec("INSERT INTO transactions (total_amount) VALUES (?)", totalAmount)
	if err != nil {
		return nil, err
	}

	transactionID64, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	transactionID := int(transactionID64)

	// PERBAIKAN: Gunakan batch insert atau prepared statement untuk efisiensi
	stmt, err := tx.Prepare("INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES (?, ?, ?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	for i := range details {
		details[i].TransactionID = transactionID
		_, err = stmt.Exec(transactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}, nil
}

// GetSalesToday untuk report hari ini
func (repo *TransactionRepository) GetSalesToday() (totalRevenue int, totalTransaksi int, err error) {
	query := `
		SELECT IFNULL(SUM(total_amount), 0), COUNT(*)
		FROM transactions
		WHERE DATE(created_at) = CURDATE()`
	
	err = repo.db.QueryRow(query).Scan(&totalRevenue, &totalTransaksi)
	return
}

// GetTopProductToday untuk produk terlaris hari ini
func (repo *TransactionRepository) GetTopProductToday() (productName string, qtySold int, err error) {
	query := `
		SELECT p.name, SUM(td.quantity) as total_qty
		FROM transaction_details td
		JOIN products p ON td.product_id = p.id
		JOIN transactions t ON td.transaction_id = t.id
		WHERE DATE(t.created_at) = CURDATE()
		GROUP BY td.product_id, p.name
		ORDER BY total_qty DESC
		LIMIT 1`
	
	err = repo.db.QueryRow(query).Scan(&productName, &qtySold)
	if err == sql.ErrNoRows {
		return "", 0, nil
	}
	return
}

// GetSalesInRange untuk report dengan date range
func (repo *TransactionRepository) GetSalesInRange(startDate, endDate string) (totalRevenue int, totalTransaksi int, err error) {
	query := `
		SELECT IFNULL(SUM(total_amount), 0), COUNT(*)
		FROM transactions
		WHERE DATE(created_at) BETWEEN ? AND ?`
	
	err = repo.db.QueryRow(query, startDate, endDate).Scan(&totalRevenue, &totalTransaksi)
	return
}

// GetTopProductInRange untuk produk terlaris dalam date range
func (repo *TransactionRepository) GetTopProductInRange(startDate, endDate string) (productName string, qtySold int, err error) {
	query := `
		SELECT p.name, SUM(td.quantity) as total_qty
		FROM transaction_details td
		JOIN products p ON td.product_id = p.id
		JOIN transactions t ON td.transaction_id = t.id
		WHERE DATE(t.created_at) BETWEEN ? AND ?
		GROUP BY td.product_id, p.name
		ORDER BY total_qty DESC
		LIMIT 1`
	
	err = repo.db.QueryRow(query, startDate, endDate).Scan(&productName, &qtySold)
	if err == sql.ErrNoRows {
		return "", 0, nil
	}
	return
}
