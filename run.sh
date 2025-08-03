#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}🛒 Shopping Cart Application Startup Script${NC}"
echo "=================================================="

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go is not installed. Please install Go 1.21+ to continue.${NC}"
    exit 1
fi

# Check if Node.js is installed
if ! command -v node &> /dev/null; then
    echo -e "${RED}❌ Node.js is not installed. Please install Node.js 16+ to continue.${NC}"
    exit 1
fi

echo -e "${YELLOW}📋 Setting up backend...${NC}"

# Navigate to backend directory
cd backend

# Install Go dependencies
echo "Installing Go dependencies..."
go mod tidy

# Set up database with sample data
echo "Setting up database with sample data..."
go run setup.go

# Start backend server in background
echo "Starting backend server on port 8080..."
go run main.go models.go handlers.go &
BACKEND_PID=$!

# Wait for backend to start
sleep 3

echo -e "${YELLOW}📋 Setting up frontend...${NC}"

# Navigate to frontend directory
cd ../frontend

# Install npm dependencies
echo "Installing npm dependencies..."
npm install

echo -e "${GREEN}✅ Setup complete!${NC}"
echo ""
echo "🚀 Starting applications:"
echo "   - Backend API: http://localhost:8080"
echo "   - Frontend App: http://localhost:3000"
echo ""
echo "📝 Demo login credentials:"
echo "   - Username: admin"
echo "   - Password: admin123"
echo ""
echo "⏹️  Press Ctrl+C to stop both servers"
echo ""

# Start frontend server
npm start

# Kill backend when frontend stops
kill $BACKEND_PID 2>/dev/null
