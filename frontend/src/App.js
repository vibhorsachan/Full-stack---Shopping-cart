import React, { useState, useEffect } from 'react';
import axios from 'axios';
import './App.css';

// API Base URL
const API_BASE_URL = 'http://localhost:8080';

function App() {
  // State management
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [user, setUser] = useState(null);
  const [token, setToken] = useState('');
  const [items, setItems] = useState([]);
  const [cart, setCart] = useState(null);

  // Form states
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');

  // Check if user is already logged in on component mount
  useEffect(() => {
    const savedToken = localStorage.getItem('token');
    const savedUser = localStorage.getItem('user');

    if (savedToken && savedUser) {
      setToken(savedToken);
      setUser(JSON.parse(savedUser));
      setIsLoggedIn(true);
      fetchItems();
    }
  }, []);

  // Set up axios interceptor for authentication
  useEffect(() => {
    if (token) {
      axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
    }
  }, [token]);

  // Login function
  const handleLogin = async (e) => {
    e.preventDefault();

    try {
      const response = await axios.post(`${API_BASE_URL}/users/login`, {
        username,
        password
      });

      const { token: userToken, user: userData } = response.data;

      // Save to localStorage
      localStorage.setItem('token', userToken);
      localStorage.setItem('user', JSON.stringify(userData));

      // Update state
      setToken(userToken);
      setUser(userData);
      setIsLoggedIn(true);

      // Clear form
      setUsername('');
      setPassword('');

      // Fetch items after login
      fetchItems();

    } catch (error) {
      console.error('Login error:', error);
      if (error.response && error.response.status === 401) {
        window.alert('Invalid username/password');
      } else {
        window.alert('Login failed. Please try again.');
      }
    }
  };

  // Logout function
  const handleLogout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    setToken('');
    setUser(null);
    setIsLoggedIn(false);
    setItems([]);
    setCart(null);
    delete axios.defaults.headers.common['Authorization'];
  };

  // Fetch all items
  const fetchItems = async () => {
    try {
      const response = await axios.get(`${API_BASE_URL}/items`);
      setItems(response.data);
    } catch (error) {
      console.error('Error fetching items:', error);
      window.alert('Failed to fetch items');
    }
  };

  // Add item to cart
  const addToCart = async (itemId) => {
    try {
      const response = await axios.post(`${API_BASE_URL}/carts`, {
        item_id: itemId,
        quantity: 1
      });

      setCart(response.data);
      window.alert('Item added to cart successfully!');

    } catch (error) {
      console.error('Error adding to cart:', error);
      window.alert('Failed to add item to cart');
    }
  };

  // Show cart items
  const showCart = async () => {
    try {
      const response = await axios.get(`${API_BASE_URL}/carts`);
      const carts = response.data;

      if (carts.length === 0) {
        window.alert('Your cart is empty');
        return;
      }

      // Find active cart
      const activeCart = carts.find(cart => cart.status === 'active');

      if (!activeCart || !activeCart.cart_items || activeCart.cart_items.length === 0) {
        window.alert('Your cart is empty');
        return;
      }

      // Prepare cart display message
      let cartMessage = 'Your Cart Items:\n\n';
      activeCart.cart_items.forEach(cartItem => {
        cartMessage += `• ${cartItem.item.name} - Quantity: ${cartItem.quantity} - $${cartItem.item.price}\n`;
      });

      window.alert(cartMessage);

    } catch (error) {
      console.error('Error fetching cart:', error);
      window.alert('Failed to fetch cart');
    }
  };

  // Show order history
  const showOrderHistory = async () => {
    try {
      const response = await axios.get(`${API_BASE_URL}/orders`);
      const orders = response.data;

      if (orders.length === 0) {
        window.alert('You have no order history');
        return;
      }

      // Prepare order history message
      let orderMessage = 'Your Order History:\n\n';
      orders.forEach(order => {
        orderMessage += `Order ID: ${order.id} - Total: $${order.total_price.toFixed(2)} - Date: ${new Date(order.created_at).toLocaleDateString()}\n`;
      });

      window.alert(orderMessage);

    } catch (error) {
      console.error('Error fetching orders:', error);
      window.alert('Failed to fetch order history');
    }
  };

  // Checkout function
  const handleCheckout = async () => {
    try {
      // First get the active cart
      const cartResponse = await axios.get(`${API_BASE_URL}/carts`);
      const carts = cartResponse.data;

      const activeCart = carts.find(cart => cart.status === 'active');

      if (!activeCart || !activeCart.cart_items || activeCart.cart_items.length === 0) {
        window.alert('Your cart is empty. Add some items before checkout.');
        return;
      }

      // Create order
      const orderResponse = await axios.post(`${API_BASE_URL}/orders`, {
        cart_id: activeCart.id
      });

      if (orderResponse.status === 201) {
        window.alert('Order successful! Thank you for your purchase.');
        setCart(null);
      }

    } catch (error) {
      console.error('Error during checkout:', error);
      window.alert('Checkout failed. Please try again.');
    }
  };

  // Login Screen Component
  if (!isLoggedIn) {
    return (
      <div className="App">
        <div className="login-container">
          <h1>Shopping Cart - Login</h1>
          <form onSubmit={handleLogin} className="login-form">
            <div className="form-group">
              <label htmlFor="username">Username:</label>
              <input
                type="text"
                id="username"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                required
              />
            </div>
            <div className="form-group">
              <label htmlFor="password">Password:</label>
              <input
                type="password"
                id="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
              />
            </div>
            <button type="submit" className="login-btn">Login</button>
          </form>
          <div className="signup-info">
            <p>Demo users: Use username "admin" with password "admin123"</p>
            <p>Or create a new user via POST /users endpoint</p>
          </div>
        </div>
      </div>
    );
  }

  // Items List Screen Component
  return (
    <div className="App">
      <header className="app-header">
        <h1>Shopping Cart</h1>
        <div className="user-info">
          <span>Welcome, {user?.username}!</span>
          <button onClick={handleLogout} className="logout-btn">Logout</button>
        </div>
      </header>

      <div className="action-buttons">
        <button onClick={handleCheckout} className="checkout-btn">Checkout</button>
        <button onClick={showCart} className="cart-btn">Cart</button>
        <button onClick={showOrderHistory} className="history-btn">Order History</button>
      </div>

      <div className="items-container">
        <h2>Available Items</h2>
        <div className="items-grid">
          {items.map(item => (
            <div key={item.id} className="item-card">
              <h3>{item.name}</h3>
              <p className="item-description">{item.description}</p>
              <p className="item-price">${item.price}</p>
              <button 
                onClick={() => addToCart(item.id)} 
                className="add-to-cart-btn"
              >
                Add to Cart
              </button>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}

export default App;
