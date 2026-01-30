# ArLOG Backend

Backend service for ArLOG - A lightweight, cost-effective, web-based log viewer for Kubernetes pods.

## Technology Stack

- **Language**: Go 1.21+
- **Web Framework**: Gorilla Mux
- **Database**: PostgreSQL with GORM
- **Kubernetes**: Official client-go library
- **Authentication**: Okta SSO (OIDC)
- **Real-time Communication**: WebSockets

## Project Structure

```
backend/
├── main.go              # Application entry point
├── database/            # Database configuration and migrations
│   └── database.go
├── models/              # Database models
│   ├── team.go
│   └── permission.go
├── handlers/            # HTTP handlers
│   ├── health.go
│   ├── permissions.go
│   ├── pods.go
│   ├── logs.go
│   └── auth.go
├── services/            # Business logic
│   └── kubernetes.go
├── middleware/          # HTTP middleware
│   └── auth.go
└── utils/               # Utility functions
```

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 13 or higher
- Kubernetes cluster (for development: kubectl proxy)
- Okta account (for SSO)

## Setup Instructions

### 1. Install Dependencies

```bash
cd backend
go mod download
```

### 2. Configure Environment Variables

Copy the example environment file and update with your values:

```bash
cp .env.example .env
```

Edit `.env` and configure:
- Database connection details
- Okta SSO credentials
- JWT secret
- Kubernetes configuration

### 3. Set Up PostgreSQL Database

Create a PostgreSQL database:

```sql
CREATE DATABASE arlog_db;
CREATE USER arlog WITH PASSWORD 'arlog_password';
GRANT ALL PRIVILEGES ON DATABASE arlog_db TO arlog;
```

### 4. Run Database Migrations

The application will automatically run migrations on startup. To seed test data:

```bash
# Set ENVIRONMENT=development in .env
go run main.go
```

### 5. For Local Development with Kubernetes

Start kubectl proxy to access the Kubernetes API:

```bash
kubectl proxy --port=8001
```

This allows the backend to communicate with your local Kubernetes cluster without authentication for development purposes.

### 6. Run the Application

```bash
go run main.go
```

The server will start on `http://localhost:8080` (or the port specified in `.env`).

## API Endpoints

### Health Check
```
GET /health
```

### User Permissions
```
GET /api/user/permissions
```
Returns the namespaces the authenticated user can access.

### List Pods
```
GET /api/pods?namespace=<namespace>
```
Lists all pods in the specified namespace.

### Stream Logs (WebSocket)
```
WS /ws/logs?namespace=<namespace>&podName=<podName>
```
Establishes a WebSocket connection to stream pod logs in real-time.

### Authentication
```
GET /auth/okta/login
GET /auth/okta/callback
```
Okta SSO authentication endpoints.

## Development

### Database Models

- **Team**: Represents a team/group mapped to an Okta group
- **Permission**: Maps teams to Kubernetes namespaces with service account tokens

### Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...
```

### Building for Production

```bash
# Build binary
go build -o arlog-backend main.go

# Run binary
./arlog-backend
```

## Security Considerations

- Service account tokens are stored in the database (should be encrypted in production)
- JWT tokens are validated for all protected endpoints
- CORS is enabled for development (should be restricted in production)
- Use environment variables for sensitive configuration

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| PORT | Server port | 8080 |
| SERVER_HOST | Server host | localhost |
| DB_HOST | Database host | localhost |
| DB_PORT | Database port | 5432 |
| DB_USER | Database user | arlog |
| DB_PASSWORD | Database password | arlog_password |
| DB_NAME | Database name | arlog_db |
| DB_SSLMODE | Database SSL mode | disable |
| KUBE_PROXY_URL | Kubernetes proxy URL | http://localhost:8001 |
| OKTA_DOMAIN | Okta domain | - |
| OKTA_CLIENT_ID | Okta client ID | - |
| OKTA_CLIENT_SECRET | Okta client secret | - |
| OKTA_REDIRECT_URI | Okta redirect URI | http://localhost:8080/auth/okta/callback |
| JWT_SECRET | JWT signing secret | - |
| ENVIRONMENT | Environment (development/production) | development |

## License

Proprietary

