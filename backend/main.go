package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	// Initialize database
	initDB()

	// Create Gin router
	r := gin.Default()

	// Enable CORS for frontend
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Routes
	setupRoutes(r)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func initDB() {
	var err error
	db, err = gorm.Open(sqlite.Open("shopping_cart.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate models
	db.AutoMigrate(&User{}, &Item{}, &Cart{}, &Order{}, &CartItem{}, &OrderItem{})

	// Create some sample items if database is empty
	var itemCount int64
	db.Model(&Item{}).Count(&itemCount)
	if itemCount == 0 {
		createSampleItems()
	}
}

func createSampleItems() {
	items := []Item{
		{Name: "iPhone 14", Description: "Latest Apple smartphone", Price: 999.99},
		{Name: "Samsung Galaxy S23", Description: "Android flagship phone", Price: 899.99},
		{Name: "MacBook Pro", Description: "Professional laptop from Apple", Price: 1999.99},
		{Name: "Dell XPS 13", Description: "Ultrabook for professionals", Price: 1299.99},
		{Name: "Nike Air Max", Description: "Comfortable running shoes", Price: 129.99},
		{Name: "Adidas Ultraboost", Description: "Premium athletic shoes", Price: 149.99},
	}

	for _, item := range items {
		db.Create(&item)
	}
	log.Println("Sample items created")
}

func setupRoutes(r *gin.Engine) {
	// User routes
	r.POST("/users", createUser)
	r.GET("/users", getUsers)
	r.POST("/users/login", loginUser)

	// Item routes
	r.POST("/items", createItem)
	r.GET("/items", getItems)

	// Cart routes (protected)
	r.POST("/carts", authMiddleware(), addToCart)
	r.GET("/carts", authMiddleware(), getCarts)

	// Order routes (protected)
	r.POST("/orders", authMiddleware(), createOrder)
	r.GET("/orders", authMiddleware(), getOrders)
}
