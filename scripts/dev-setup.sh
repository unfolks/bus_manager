#!/bin/bash

# Development Setup Script for Bus Manager
# This script sets up the environment for local development with Redis & PostgreSQL in Docker

echo "🚌 Setting up Bus Manager Development Environment..."

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker first."
    exit 1
fi

# Check if Docker Compose is available
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed. Please install Docker Compose first."
    exit 1
fi

# Start PostgreSQL and Redis containers
echo "🐘 Starting PostgreSQL and Redis containers..."
docker-compose -f docker-compose.dev.yml up -d

# Wait for databases to be ready
echo "⏳ Waiting for databases to be ready..."
sleep 10

# Check if containers are running
if ! docker ps | grep -q "bus-manager-db"; then
    echo "❌ PostgreSQL container failed to start."
    exit 1
fi

if ! docker ps | grep -q "bus-manager-redis"; then
    echo "❌ Redis container failed to start."
    exit 1
fi

echo "✅ PostgreSQL and Redis are running!"

# Setup backend
echo "🔧 Setting up backend..."
cd backend

# Copy environment file if it doesn't exist
if [ ! -f .env ]; then
    cp .env.local .env
    echo "📝 Created .env file from .env.local"
fi

# Install Go dependencies
echo "📦 Installing Go dependencies..."
go mod download

echo "✅ Backend setup complete!"

# Setup frontend
echo "🎨 Setting up frontend..."
cd ../frontend

# Install Node.js dependencies
echo "📦 Installing Node.js dependencies..."
npm install

echo "✅ Frontend setup complete!"

cd ..

echo ""
echo "🎉 Development environment setup complete!"
echo ""
echo "📋 Next steps:"
echo "1. Start the backend server:"
echo "   cd backend && go run cmd/api/main.go"
echo ""
echo "2. In another terminal, start the frontend:"
echo "   cd frontend && npm start"
echo ""
echo "3. Access the application:"
echo "   Frontend: http://localhost:3000"
echo "   Backend API: http://localhost:8080"
echo ""
echo "🛑 To stop the databases:"
echo "   docker-compose -f docker-compose.dev.yml down"
echo ""
