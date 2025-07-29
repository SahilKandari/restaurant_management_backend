# Restaurant Management System - Backend

A comprehensive restaurant management system backend built with Go, featuring user authentication, order management, inventory tracking, and more.

## Features

- **User Management**: Authentication and authorization with JWT tokens
- **Menu Management**: CRUD operations for food items and menus
- **Order Management**: Order creation, tracking, and status updates
- **Table Management**: Table reservation and management
- **Invoice Generation**: PDF invoice generation
- **Email Notifications**: SMTP integration for password reset and notifications
- **Redis Integration**: Caching and session management

## Tech Stack

- **Language**: Go (Golang)
- **Database**: MongoDB
- **Cache**: Redis
- **Authentication**: JWT tokens
- **Email**: SMTP
- **PDF Generation**: Custom PDF helper

## Project Structure

```
backend/
├── config/           # Configuration files
├── controllers/      # Request handlers
├── database/         # Database connection and setup
├── helpers/          # Utility functions
├── middlewares/      # Authentication and other middlewares
├── models/           # Data models
├── routes/           # API routes
└── utils/            # Common utilities
```

## Getting Started

### Prerequisites

- Go 1.19 or higher
- MongoDB
- Redis

### Installation

1. Clone the repository:
```bash
git clone https://github.com/SahilKandari/restaurant-management-backend.git
cd restaurant-management-backend
```

2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables:
```bash
cp .env.example .env
# Edit .env with your configuration
```

4. Run the application:
```bash
go run main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

### Authentication
- `POST /api/auth/signup` - User registration
- `POST /api/auth/login` - User login
- `POST /api/auth/forgot-password` - Password reset request
- `POST /api/auth/reset-password` - Reset password

### Menu Management
- `GET /api/menus` - Get all menus
- `POST /api/menus` - Create menu
- `PUT /api/menus/:id` - Update menu
- `DELETE /api/menus/:id` - Delete menu

### Food Management
- `GET /api/foods` - Get all food items
- `POST /api/foods` - Create food item
- `PUT /api/foods/:id` - Update food item
- `DELETE /api/foods/:id` - Delete food item

### Order Management
- `GET /api/orders` - Get all orders
- `POST /api/orders` - Create order
- `PUT /api/orders/:id` - Update order
- `DELETE /api/orders/:id` - Delete order

### Table Management
- `GET /api/tables` - Get all tables
- `POST /api/tables` - Create table
- `PUT /api/tables/:id` - Update table
- `DELETE /api/tables/:id` - Delete table

## Environment Variables

```
PORT=8080
MONGODB_URL=mongodb://localhost:27017/restaurant
REDIS_URL=redis://localhost:6379
JWT_SECRET=your-jwt-secret
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
