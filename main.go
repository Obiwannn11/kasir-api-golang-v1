package main

import (
	"fmt"
	"kasir-api-golang-v1/database"
	"kasir-api-golang-v1/handlers"
	"kasir-api-golang-v1/repositories"
	"kasir-api-golang-v1/services"
	"log"
	"net/http"
	"github.com/spf13/viper"
)

// Config struct untuk mapping .env
type Config struct {
	Port   string `mapstructure:"PORT"`
	DBUser string `mapstructure:"DB_USER"`
	DBPass string `mapstructure:"DB_PASS"`
	DBHost string `mapstructure:"DB_HOST"`
	DBPort string `mapstructure:"DB_PORT"`
	DBName string `mapstructure:"DB_NAME"`
}

func main() {
	// 1. Setup Viper
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Println("No .env file found, using system env")
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal("Error loading config:", err)
	}

	// 2. Connect Database MySQL
	db, err := database.InitDB(config.DBUser, config.DBPass, config.DBHost, config.DBPort, config.DBName)
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}
	defer db.Close()

	// 3. Dependency Injection (Wiring)
	// Repositories
	productRepo := repositories.NewProductRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)

	// Services
	// Perhatikan: CategoryService butuh akses ke productRepo juga untuk mindahin barang
	categoryService := services.NewCategoryService(categoryRepo, productRepo)
	productService := services.NewProductService(productRepo) 

	// Handlers
	categoryHandler := handlers.NewCategoryHandler(categoryService) 
	productHandler := handlers.NewProductHandler(productService)

	// 4. Routes
	mux := http.NewServeMux()
	mux.HandleFunc("/api/categories", categoryHandler.HandleCategories)
	mux.HandleFunc("/api/categories/", categoryHandler.HandleCategoryDelete) // Handle delete by ID
	mux.HandleFunc("/api/products", productHandler.HandleProducts)
	mux.HandleFunc("/api/products/", productHandler.HandleProductByID)

	addr := ":" + config.Port
	fmt.Println("Server running on MySQL at", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}