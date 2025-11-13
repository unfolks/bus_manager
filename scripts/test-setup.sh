#!/bin/bash

# Test script to verify the local development setup

echo "🧪 Testing Bus Manager Development Setup..."

# Test 1: Check if Docker containers are running
echo "1️⃣ Checking Docker containers..."
if docker ps | grep -q "bus-manager-db"; then
    echo "✅ PostgreSQL container is running"
else
    echo "❌ PostgreSQL container is not running"
    exit 1
fi

if docker ps | grep -q "bus-manager-redis"; then
    echo "✅ Redis container is running"
else
    echo "❌ Redis container is not running"
    exit 1
fi

# Test 2: Test database connections
echo ""
echo "2️⃣ Testing database connections..."

# Test Redis
if docker exec bus-manager-redis redis-cli ping | grep -q "PONG"; then
    echo "✅ Redis connection successful"
else
    echo "❌ Redis connection failed"
    exit 1
fi

# Test PostgreSQL
if docker exec bus-manager-db pg_isready -U postgres | grep -q "accepting connections"; then
    echo "✅ PostgreSQL connection successful"
else
    echo "❌ PostgreSQL connection failed"
    exit 1
fi

# Test 3: Check backend setup
echo ""
echo "3️⃣ Checking backend setup..."
cd backend

if [ -f .env ]; then
    echo "✅ Backend .env file exists"
else
    echo "❌ Backend .env file missing"
    exit 1
fi

if go mod verify > /dev/null 2>&1; then
    echo "✅ Go modules are valid"
else
    echo "❌ Go modules verification failed"
    exit 1
fi

# Test 4: Check frontend setup
echo ""
echo "4️⃣ Checking frontend setup..."
cd ../frontend

if [ -d node_modules ]; then
    echo "✅ Frontend dependencies are installed"
else
    echo "❌ Frontend dependencies not installed"
    exit 1
fi

if npm list --depth=0 > /dev/null 2>&1; then
    echo "✅ Frontend packages are valid"
else
    echo "❌ Frontend packages verification failed"
    exit 1
fi

cd ..

echo ""
echo "🎉 All tests passed! Your development environment is ready to use."
echo ""
echo "📋 To start development:"
echo "1. Start backend: cd backend && go run cmd/api/main.go"
echo "2. Start frontend: cd frontend && npm start"
echo "3. Access app: http://localhost:3000"
echo ""
echo "🛑 To stop databases: docker-compose -f docker-compose.dev.yml down"
