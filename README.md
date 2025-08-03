# Shopping Cart Application

A full-stack e-commerce shopping cart application built with Go (Gin framework) backend and React frontend. This project demonstrates a complete shopping flow from user registration/login to order placement.

## Project Overview

This application implements a simple e-commerce shopping cart system with the following features:

- User registration and authentication
- Item browsing and cart management
- Order placement and history tracking
- RESTful API architecture
- Responsive web interface

## Technology Stack

### Backend
- **Go 1.21+** - Programming language
- **Gin** - Web framework for HTTP routing
- **GORM** - ORM for database operations
- **SQLite** - Database for data persistence
- **bcrypt** - Password hashing for security

### Frontend
- **React 18** - Frontend JavaScript library
- **Axios** - HTTP client for API communication
- **CSS3** - Styling and responsive design
- **HTML5** - Markup structure

## Project Structure

```
shopping-cart-complete/
├── backend/
│   ├── main.go           # Main application entry point
│   ├── models.go         # Database models and structs
│   ├── handlers.go       # API route handlers
│   ├── setup.go          # Database setup and sample data
│   └── go.mod           # Go module dependencies
└── frontend/
    ├── public/
    │   └── index.html    # HTML template
    ├── src/
    │   ├── App.js        # Main React component
    │   ├── App.css       # Application styles
    │   ├── index.js      # React app entry point
    │   └── index.css     # Global styles
    └── package.json      # Node.js dependencies
```

## API Endpoints

### User Management
- `POST /users` - Create a new user account
- `GET /users` - List all users (for admin purposes)
- `POST /users/login` - User authentication and token generation

### Item Management
- `POST /items` - Create new items (admin function)
- `GET /items` - List all available items

### Cart Operations (Protected)
- `POST /carts` - Add items to user's cart
- `GET /carts` - Retrieve user's cart information

### Order Management (Protected)
- `POST /orders` - Convert cart to completed order
- `GET /orders` - View user's order history

## Installation and Setup

### Prerequisites
- Go 1.21 or higher
- Node.js 16+ and npm
- Git for version control

### Backend Setup

1. **Navigate to backend directory:**
   ```bash
   cd shopping-cart-complete/backend
   ```

2. **Initialize Go module and install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Set up database with sample data:**
   ```bash
   go run setup.go
   ```

4. **Start the backend server:**
   ```bash
   go run main.go models.go handlers.go
   ```

   The backend server will start on `http://localhost:8080`

### Frontend Setup

1. **Open a new terminal and navigate to frontend directory:**
   ```bash
   cd shopping-cart-complete/frontend
   ```

2. **Install Node.js dependencies:**
   ```bash
   npm install
   ```

3. **Start the React development server:**
   ```bash
   npm start
   ```

   The frontend will start on `http://localhost:3000` and automatically open in your browser.

## Usage Instructions

### Getting Started

1. **Access the Application:**
   - Open your web browser and go to `http://localhost:3000`
   - You'll see the login screen

2. **Login with Sample Account:**
   - Username: `admin`
   - Password: `admin123`
   - Click "Login" to access the shopping interface

3. **Browse and Shop:**
   - View available items on the main page
   - Click "Add to Cart" on any item to add it to your cart
   - Use the "Cart" button to view current cart contents

4. **Manage Orders:**
   - Click "Checkout" to convert your cart to an order
   - Use "Order History" to view your previous orders
   - Click "Logout" when finished

### Creating New Users

To create additional user accounts, you can use the API endpoint:

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"username": "newuser", "password": "password123"}'
```

### Adding New Items

To add new items to the store (admin function):

```bash
curl -X POST http://localhost:8080/items \
  -H "Content-Type: application/json" \
  -d '{"name": "New Product", "description": "Product description", "price": 99.99}'
```

## Database Schema

### Users Table
- `id` - Primary key
- `username` - Unique username for login
- `password` - Hashed password using bcrypt
- `token` - Authentication token for sessions

### Items Table
- `id` - Primary key
- `name` - Product name
- `description` - Product description
- `price` - Product price in USD

### Carts Table
- `id` - Primary key
- `user_id` - Foreign key to users
- `status` - Cart status (active/ordered)

### Cart Items Table
- `id` - Primary key
- `cart_id` - Foreign key to carts
- `item_id` - Foreign key to items
- `quantity` - Number of items

### Orders Table
- `id` - Primary key
- `user_id` - Foreign key to users
- `cart_id` - Foreign key to carts
- `total_price` - Total order amount
- `status` - Order status

### Order Items Table
- `id` - Primary key
- `order_id` - Foreign key to orders
- `item_id` - Foreign key to items
- `quantity` - Number of items ordered
- `price` - Price at time of order

## Security Features

- **Password Hashing:** User passwords are hashed using bcrypt
- **Token Authentication:** JWT-like tokens for session management
- **CORS Protection:** Configured for frontend-backend communication
- **Input Validation:** Request data validation using Gin binding

## Testing the API

You can test the API endpoints using curl or Postman:

### Login Example
```bash
curl -X POST http://localhost:8080/users/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin123"}'
```

### Get Items Example
```bash
curl -X GET http://localhost:8080/items
```

### Add to Cart Example (requires authentication)
```bash
curl -X POST http://localhost:8080/carts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{"item_id": 1, "quantity": 2}'
```

## Troubleshooting

### Common Issues

1. **Backend won't start:**
   - Ensure Go 1.21+ is installed
   - Run `go mod tidy` to install dependencies
   - Check if port 8080 is available

2. **Frontend won't start:**
   - Ensure Node.js 16+ is installed
   - Delete `node_modules` and run `npm install` again
   - Check if port 3000 is available

3. **Database issues:**
   - Delete `shopping_cart.db` file and run `setup.go` again
   - Ensure write permissions in the backend directory

4. **Authentication problems:**
   - Make sure you're using the correct login credentials
   - Check browser console for any JavaScript errors

## Future Enhancements

Potential improvements for this application:

- **Inventory Management:** Track item stock levels
- **Payment Integration:** Add payment processing
- **User Profiles:** Extended user information and preferences
- **Product Categories:** Organize items into categories
- **Search Functionality:** Search and filter products
- **Admin Dashboard:** Administrative interface for managing items and orders
- **Email Notifications:** Order confirmation emails
- **Mobile App:** React Native mobile application

## Development Notes

This project was designed as a learning exercise to demonstrate:
- Full-stack development skills
- RESTful API design principles
- Database relationship modeling
- Authentication and authorization
- Responsive web design
- Clean code practices and documentation

The code is written in a beginner-friendly style with extensive comments and clear variable names to aid understanding and maintenance.

## License

This project is created for educational purposes and is free to use and modify.

## Contact

For questions or suggestions about this project, please create an issue in the repository or contact the development team.
