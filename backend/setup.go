package main

import (
	"log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Connect to database
	db, err := gorm.Open(sqlite.Open("shopping_cart.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate models
	db.AutoMigrate(&User{}, &Item{}, &Cart{}, &Order{}, &CartItem{}, &OrderItem{})

	// Create sample user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	sampleUser := User{
		Username: "admin",
		Password: string(hashedPassword),
	}

	// Check if user already exists
	var existingUser User
	if err := db.Where("username = ?", "admin").First(&existingUser).Error; err != nil {
		// User doesn't exist, create it
		if err := db.Create(&sampleUser).Error; err != nil {
			log.Printf("Error creating sample user: %v", err)
		} else {
			log.Println("Sample user created: username=admin, password=admin123")
		}
	} else {
		log.Println("Sample user already exists")
	}

	// Create sample items if they don't exist
	var itemCount int64
	db.Model(&Item{}).Count(&itemCount)
	if itemCount == 0 {
		items := []Item{
			{Name: "iPhone 14", Description: "Latest Apple smartphone with advanced features", Price: 999.99},
			{Name: "Samsung Galaxy S23", Description: "Android flagship phone with excellent camera", Price: 899.99},
			{Name: "MacBook Pro", Description: "Professional laptop from Apple for developers", Price: 1999.99},
			{Name: "Dell XPS 13", Description: "Ultrabook perfect for students and professionals", Price: 1299.99},
			{Name: "Nike Air Max", Description: "Comfortable running shoes for daily use", Price: 129.99},
			{Name: "Adidas Ultraboost", Description: "Premium athletic shoes for serious runners", Price: 149.99},
			{Name: "Sony WH-1000XM4", Description: "Noise-canceling wireless headphones", Price: 349.99},
			{Name: "Apple Watch Series 8", Description: "Smartwatch with health monitoring features", Price: 399.99},
		}

		for _, item := range items {
			if err := db.Create(&item).Error; err != nil {
				log.Printf("Error creating item %s: %v", item.Name, err)
			}
		}
		log.Println("Sample items created successfully")
	} else {
		log.Println("Sample items already exist")
	}

	log.Println("Database setup completed!")
}
