package main

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// User handlers

func createUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already exists
	var existingUser User
	if err := db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create user
	user := User{
		Username: req.Username,
		Password: string(hashedPassword),
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Remove password from response
	user.Password = ""
	c.JSON(http.StatusCreated, user)
}

func getUsers(c *gin.Context) {
	var users []User
	if err := db.Select("id, username, created_at, updated_at").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func loginUser(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user
	var user User
	if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Generate token
	token, err := generateToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Update user with token
	user.Token = token
	db.Save(&user)

	// Remove password from response
	user.Password = ""

	c.JSON(http.StatusOK, LoginResponse{
		Token: token,
		User:  user,
	})
}

// Item handlers

func createItem(c *gin.Context) {
	var item Item
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create item"})
		return
	}

	c.JSON(http.StatusCreated, item)
}

func getItems(c *gin.Context) {
	var items []Item
	if err := db.Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch items"})
		return
	}

	c.JSON(http.StatusOK, items)
}

// Cart handlers

func addToCart(c *gin.Context) {
	userID := c.GetUint("userID")

	var req AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default quantity if not provided
	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	// Check if item exists
	var item Item
	if err := db.First(&item, req.ItemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	// Find or create active cart for user
	var cart Cart
	if err := db.Where("user_id = ? AND status = ?", userID, "active").First(&cart).Error; err != nil {
		// Create new cart
		cart = Cart{
			UserID: userID,
			Status: "active",
		}
		if err := db.Create(&cart).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cart"})
			return
		}
	}

	// Check if item already exists in cart
	var existingCartItem CartItem
	if err := db.Where("cart_id = ? AND item_id = ?", cart.ID, req.ItemID).First(&existingCartItem).Error; err == nil {
		// Update quantity
		existingCartItem.Quantity += req.Quantity
		db.Save(&existingCartItem)
	} else {
		// Add new item to cart
		cartItem := CartItem{
			CartID:   cart.ID,
			ItemID:   req.ItemID,
			Quantity: req.Quantity,
		}
		if err := db.Create(&cartItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add item to cart"})
			return
		}
	}

	// Return updated cart
	if err := db.Preload("CartItems.Item").First(&cart, cart.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart"})
		return
	}

	c.JSON(http.StatusOK, cart)
}

func getCarts(c *gin.Context) {
	userID := c.GetUint("userID")

	var carts []Cart
	if err := db.Where("user_id = ?", userID).Preload("CartItems.Item").Find(&carts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch carts"})
		return
	}

	c.JSON(http.StatusOK, carts)
}

// Order handlers

func createOrder(c *gin.Context) {
	userID := c.GetUint("userID")

	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find cart
	var cart Cart
	if err := db.Where("id = ? AND user_id = ? AND status = ?", req.CartID, userID, "active").
		Preload("CartItems.Item").First(&cart).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found or already ordered"})
		return
	}

	if len(cart.CartItems) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart is empty"})
		return
	}

	// Calculate total price
	var totalPrice float64
	for _, cartItem := range cart.CartItems {
		totalPrice += cartItem.Item.Price * float64(cartItem.Quantity)
	}

	// Create order
	order := Order{
		UserID:     userID,
		CartID:     cart.ID,
		TotalPrice: totalPrice,
		Status:     "completed",
	}

	if err := db.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	// Create order items
	for _, cartItem := range cart.CartItems {
		orderItem := OrderItem{
			OrderID:  order.ID,
			ItemID:   cartItem.ItemID,
			Quantity: cartItem.Quantity,
			Price:    cartItem.Item.Price,
		}
		db.Create(&orderItem)
	}

	// Mark cart as ordered
	cart.Status = "ordered"
	db.Save(&cart)

	// Load order with items
	if err := db.Preload("OrderItems.Item").First(&order, order.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch order"})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func getOrders(c *gin.Context) {
	userID := c.GetUint("userID")

	var orders []Order
	if err := db.Where("user_id = ?", userID).Preload("OrderItems.Item").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// Utility functions

func generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// Auth middleware
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Remove "Bearer " prefix if present
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		var user User
		if err := db.Where("token = ?", token).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("userID", user.ID)
		c.Next()
	}
}
