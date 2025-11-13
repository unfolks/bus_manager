# Quick Start Guide

## 🚀 One-Command Setup

```bash
./scripts/dev-setup.sh
```

## 📋 Manual Steps

### 1. Start Databases
```bash
docker-compose -f docker-compose.dev.yml up -d
```

### 2. Start Backend
```bash
cd backend
go run cmd/api/main.go
```

### 3. Start Frontend (in new terminal)
```bash
cd frontend
npm start
```

### 4. Access Application
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080

## 🛑 Stop Services

```bash
# Stop databases
docker-compose -f docker-compose.dev.yml down

# Stop backend/frontend with Ctrl+C in their terminals
```

## ✅ Verify Setup

```bash
./scripts/test-setup.sh
```

## 📁 Project Structure

```
bus-manager/
├── docker-compose.dev.yml      # Development databases only
├── scripts/
│   ├── dev-setup.sh           # Automated setup
│   └── test-setup.sh          # Verify setup
├── backend/
│   ├── .env.local             # Local dev config
│   └── cmd/api/main.go        # Backend entry point
└── frontend/
    └── src/                   # React app
```

## 🔧 Development Workflow

1. **Make changes** to backend/frontend code
2. **Backend** auto-restarts on file changes
3. **Frontend** hot-reloads in browser
4. **Databases** persist data between sessions

## 📚 Documentation

- `README-DEV.md` - Detailed development guide
- `README.md` - Project overview and API docs
