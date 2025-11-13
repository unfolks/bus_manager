# Bus Manager Game

A web-based bus company management game set in Indonesia, built with Golang backend and React TypeScript frontend.

## Features

### Core Gameplay
- **Company Management**: Create and manage your bus company
- **Depot Management**: Build and expand depots to accommodate more buses
- **Fleet Management**: Purchase, upgrade, and manage buses of different types
- **Route Management**: Select and operate routes across Indonesia
- **Driver Management**: Recruit and manage drivers and conductors
- **Real-time Tracking**: Monitor buses on live map using Leaflet.js
- **Financial System**: Track revenue, expenses, and profits

### Bus Types
1. Normal
2. High Decker
3. Super High Decker
4. High Decker Double Glass
5. Double Decker
6. Ultra High Decker

### Service Types
- Economy
- Business
- Executive
- Night Bus (for long-distance routes)

### Trip Types
- Intercity within province
- Intercity between provinces

## Technology Stack

### Backend
- **Language**: Go 1.21
- **Framework**: Gin
- **Database**: PostgreSQL
- **Cache**: Redis
- **Authentication**: JWT
- **Real-time**: WebSockets

### Frontend
- **Language**: TypeScript
- **Framework**: React
- **Maps**: Leaflet.js
- **State Management**: React Context/Hooks
- **Styling**: CSS Modules

### Infrastructure
- **Containerization**: Docker & Docker Compose
- **Web Server**: Nginx (for frontend)

## Project Structure

```
bus-manager/
├── backend/
│   ├── cmd/api/           # Main application entry point
│   ├── internal/
│   │   ├── database/      # Database connection and migrations
│   │   ├── handlers/      # HTTP handlers
│   │   ├── middleware/    # Middleware (auth, CORS, etc.)
│   │   └── models/        # Data models
│   ├── go.mod
│   ├── go.sum
│   └── Dockerfile
├── frontend/
│   ├── public/
│   ├── src/
│   │   ├── components/    # React components
│   │   ├── pages/         # Page components
│   │   ├── services/      # API services
│   │   ├── types/         # TypeScript types
│   │   └── utils/         # Utility functions
│   ├── package.json
│   ├── Dockerfile
│   └── nginx.conf
├── docker-compose.yml
└── README.md
```

## Setup Instructions

### Prerequisites
- Docker and Docker Compose
- Go 1.21+ (for local development)
- Node.js 18+ (for local development)
- PostgreSQL 15+ (for local development)
- Redis 7+ (for local development)

### Using Docker Compose (Recommended)

1. Clone the repository
2. Copy environment file:
   ```bash
   cp backend/.env.example backend/.env
   ```
3. Start all services:
   ```bash
   docker-compose up --build
   ```
4. Access the application:
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080

### Local Development

#### Backend Setup
1. Navigate to backend directory:
   ```bash
   cd backend
   ```
2. Install dependencies:
   ```bash
   go mod download
   ```
3. Set up environment variables (copy .env.example to .env)
4. Run the application:
   ```bash
   go run cmd/api/main.go
   ```

#### Frontend Setup
1. Navigate to frontend directory:
   ```bash
   cd frontend
   ```
2. Install dependencies:
   ```bash
   npm install
   ```
3. Start development server:
   ```bash
   npm start
   ```

## API Endpoints

### Authentication
- `POST /auth/register` - Register new user
- `POST /auth/login` - User login
- `POST /auth/logout` - User logout
- `POST /auth/refresh` - Refresh JWT token

### Company Management
- `GET /company` - Get user's company
- `POST /company` - Create new company

### Depot Management
- `GET /depots` - Get user's depots
- `POST /depots` - Create new depot

### Bus Management
- `GET /buses` - Get user's buses
- `POST /buses` - Purchase new bus

### Route Management
- `GET /routes` - Get available routes

### Trip Management
- `GET /trips/active` - Get active trips
- `POST /trips` - Create new trip

### WebSocket
- `WS /ws` - Real-time trip updates

## Game Flow

1. **Registration**: Create account with email, username, and password
2. **Company Setup**: Name your company and place initial depot
3. **Starting Fleet**: Receive a small capacity bus with limited range
4. **Route Selection**: Choose initial route within bus range
5. **Operations**: Dispatch buses, manage fuel, monitor trips
6. **Growth**: Expand fleet, upgrade buses, open new routes
7. **Advanced Features**: Hire drivers, upgrade depots, manage reputation

## Database Schema

### Core Tables
- `users` - User accounts
- `companies` - Bus companies
- `depots` - Bus depots
- `buses` - Bus fleet
- `routes` - Available routes
- `trips` - Active/completed trips
- `drivers` - Driver staff
- `transactions` - Financial records

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License.

## Future Enhancements

- Multiplayer support
- Advanced AI competitors
- Weather system affecting routes
- Bus customization options
- Marketing and advertising system
- Staff training and development
- Maintenance and repair system
- Fuel price fluctuations
- Special events and challenges
