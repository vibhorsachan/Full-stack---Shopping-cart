package main

import (
	"time"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"unique;not null"`
	Password  string         `json:"password" gorm:"not null"`
	Token     string         `json:"token,omitempty" gorm:"index"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Item represents a product in the store
type Item struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	Price       float64        `json:"price" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// Cart represents a shopping cart
type Cart struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null;index"`
	User      User           `json:"user" gorm:"foreignKey:UserID"`
	CartItems []CartItem     `json:"cart_items" gorm:"foreignKey:CartID"`
	Status    string         `json:"status" gorm:"default:'active'"` // active, ordered
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// CartItem represents items in a cart
type CartItem struct {
	ID       uint  `json:"id" gorm:"primaryKey"`
	CartID   uint  `json:"cart_id" gorm:"not null;index"`
	Cart     Cart  `json:"cart" gorm:"foreignKey:CartID"`
	ItemID   uint  `json:"item_id" gorm:"not null;index"`
	Item     Item  `json:"item" gorm:"foreignKey:ItemID"`
	Quantity int   `json:"quantity" gorm:"default:1"`
}

// Order represents a completed order
type Order struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	UserID     uint           `json:"user_id" gorm:"not null;index"`
	User       User           `json:"user" gorm:"foreignKey:UserID"`
	CartID     uint           `json:"cart_id" gorm:"not null;index"`
	Cart       Cart           `json:"cart" gorm:"foreignKey:CartID"`
	OrderItems []OrderItem    `json:"order_items" gorm:"foreignKey:OrderID"`
	TotalPrice float64        `json:"total_price"`
	Status     string         `json:"status" gorm:"default:'completed'"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

// OrderItem represents items in an order
type OrderItem struct {
	ID       uint    `json:"id" gorm:"primaryKey"`
	OrderID  uint    `json:"order_id" gorm:"not null;index"`
	Order    Order   `json:"order" gorm:"foreignKey:OrderID"`
	ItemID   uint    `json:"item_id" gorm:"not null;index"`
	Item     Item    `json:"item" gorm:"foreignKey:ItemID"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"` // Price at time of order
}

// Request/Response structs
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AddToCartRequest struct {
	ItemID   uint `json:"item_id" binding:"required"`
	Quantity int  `json:"quantity"`
}

type CreateOrderRequest struct {
	CartID uint `json:"cart_id" binding:"required"`
}
