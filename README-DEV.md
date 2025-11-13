# Local Development Setup

This guide will help you set up the Bus Manager project for local development with Redis and PostgreSQL running in Docker containers, while the backend and frontend run locally on your machine.

## Prerequisites

- Docker and Docker Compose
- Go 1.23+ (for backend)
- Node.js 18+ (for frontend)
- Git

## Quick Start

### Option 1: Automated Setup (Recommended)

Run the setup script to configure everything automatically:

```bash
./scripts/dev-setup.sh
```

This script will:
- Start PostgreSQL and Redis containers
- Install Go and Node.js dependencies
- Set up environment files
- Provide instructions for starting the services

### Option 2: Manual Setup

#### 1. Start Database Services

Start PostgreSQL and Redis in Docker:

```bash
docker-compose -f docker-compose.dev.yml up -d
```

#### 2. Backend Setup

```bash
cd backend

# Copy environment file
cp .env.local .env

# Install dependencies
go mod download

# Start the backend server
go run cmd/api/main.go
```

The backend will be available at http://localhost:8080

#### 3. Frontend Setup

In a new terminal:

```bash
cd frontend

# Install dependencies
npm install

# Start the development server
npm start
```

The frontend will be available at http://localhost:3000

## Development Workflow

### Running the Application

1. **Start databases** (if not already running):
   ```bash
   docker-compose -f docker-compose.dev.yml up -d
   ```

2. **Start backend** (in terminal 1):
   ```bash
   cd backend
   go run cmd/api/main.go
   ```

3. **Start frontend** (in terminal 2):
   ```bash
   cd frontend
   npm start
   ```

### Stopping the Application

1. Stop the backend and frontend with `Ctrl+C` in their respective terminals
2. Stop the databases:
   ```bash
   docker-compose -f docker-compose.dev.yml down
   ```

## Project Structure

```
bus-manager/
├── docker-compose.dev.yml      # Development Docker setup (databases only)
├── docker-compose.yml          # Production Docker setup (full stack)
├── scripts/
│   └── dev-setup.sh           # Automated setup script
├── backend/
│   ├── .env.local             # Local development environment config
│   ├── .env.example           # Environment template
│   ├── cmd/api/main.go        # Backend entry point
│   └── internal/              # Backend source code
└── frontend/
    ├── package.json           # Frontend dependencies
    └── src/                   # Frontend source code
```

## Environment Configuration

### Backend Environment (.env)

The backend uses the following configuration for local development:

```env
# Database (PostgreSQL in Docker)
DB_HOST=localhost
DB_PORT=5433
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=bus_manager
DB_SSLMODE=disable
DB_TIMEZONE=Asia/Jakarta

# Redis (in Docker)
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=

# JWT
JWT_SECRET=your-secret-key-change-in-production

# Server
PORT=8080
GIN_MODE=debug
```

### Frontend Environment

The frontend automatically connects to the local backend at `http://localhost:8080`.

## Database Access

### PostgreSQL

- **Host**: localhost
- **Port**: 5433
- **Database**: bus_manager
- **Username**: postgres
- **Password**: password

You can connect using:
```bash
psql -h localhost -p 5433 -U postgres -d bus_manager
```

### Redis

- **Host**: localhost
- **Port**: 6379

You can connect using:
```bash
redis-cli -p 6379
```

## Troubleshooting

### Port Conflicts

If you encounter port conflicts, make sure:
- PostgreSQL is not running on port 5433
- Redis is not running on port 6379
- Backend is not running on port 8080
- Frontend is not running on port 3000

### Database Connection Issues

1. Check if PostgreSQL container is running:
   ```bash
   docker ps | grep bus-manager-db
   ```

2. Check container logs:
   ```bash
   docker logs bus-manager-db
   ```

3. Restart the database:
   ```bash
   docker-compose -f docker-compose.dev.yml restart postgres
   ```

### Redis Connection Issues

1. Check if Redis container is running:
   ```bash
   docker ps | grep bus-manager-redis
   ```

2. Test Redis connection:
   ```bash
   redis-cli -p 6379 ping
   ```

### Backend Issues

1. Check Go version:
   ```bash
   go version
   ```

2. Clean and rebuild:
   ```bash
   cd backend
   go clean -cache
   go mod tidy
   go run cmd/api/main.go
   ```

### Frontend Issues

1. Clear node modules and reinstall:
   ```bash
   cd frontend
   rm -rf node_modules package-lock.json
   npm install
   npm start
   ```

## Development Tips

### Hot Reload

- **Backend**: The Go server will automatically restart when you save changes
- **Frontend**: React development server supports hot reload out of the box

### Database Migrations

Database migrations are automatically applied when PostgreSQL starts. If you need to re-run migrations:

```bash
docker-compose -f docker-compose.dev.yml down -v
docker-compose -f docker-compose.dev.yml up -d
```

### Debugging

- **Backend**: Set `GIN_MODE=debug` in `.env` for detailed logging
- **Frontend**: Use browser developer tools for React debugging

## Production Deployment

For production deployment, use the full Docker setup:

```bash
docker-compose up --build
```

This will run the entire application stack in containers, including the backend and frontend.

## Contributing

1. Create a feature branch
2. Make your changes
3. Test locally using this development setup
4. Submit a pull request
