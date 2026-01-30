.PHONY: help install-backend install-frontend install run-backend run-frontend run dev stop clean build docker-up docker-down

# Default target
help:
	@echo "ArLOG - Kubernetes Log Viewer"
	@echo ""
	@echo "Available commands:"
	@echo "  make install          - Install all dependencies (backend + frontend)"
	@echo "  make install-backend  - Install backend dependencies"
	@echo "  make install-frontend - Install frontend dependencies"
	@echo "  make run             - Run both backend and frontend"
	@echo "  make run-backend     - Run backend server"
	@echo "  make run-frontend    - Run frontend dev server"
	@echo "  make dev             - Run in development mode (parallel)"
	@echo "  make build           - Build both backend and frontend"
	@echo "  make docker-up       - Start Docker containers"
	@echo "  make docker-down     - Stop Docker containers"
	@echo "  make clean           - Clean build artifacts"
	@echo "  make test            - Run tests"
	@echo "  make kubectl-proxy   - Start kubectl proxy"

# Install all dependencies
install: install-backend install-frontend
	@echo "✅ All dependencies installed"

# Install backend dependencies
install-backend:
	@echo "📦 Installing backend dependencies..."
	cd backend && go mod download
	@echo "✅ Backend dependencies installed"

# Install frontend dependencies
install-frontend:
	@echo "📦 Installing frontend dependencies..."
	cd frontend && npm install
	@echo "✅ Frontend dependencies installed"

# Run backend
run-backend:
	@echo "🚀 Starting backend server..."
	cd backend && go run main.go

# Run frontend
run-frontend:
	@echo "🚀 Starting frontend dev server..."
	cd frontend && npm run dev

# Run both (requires separate terminals or use 'make dev')
run:
	@echo "Use 'make dev' to run both backend and frontend in parallel"
	@echo "Or run in separate terminals:"
	@echo "  Terminal 1: make run-backend"
	@echo "  Terminal 2: make run-frontend"

# Development mode (requires tmux or run in separate terminals)
dev:
	@echo "Starting development environment..."
	@echo "Backend: http://localhost:8080"
	@echo "Frontend: http://localhost:5173"
	@echo ""
	@echo "Run these commands in separate terminals:"
	@echo "  1. make run-backend"
	@echo "  2. make run-frontend"

# Build backend
build-backend:
	@echo "🔨 Building backend..."
	cd backend && go build -o arlog-backend main.go
	@echo "✅ Backend built: backend/arlog-backend"

# Build frontend
build-frontend:
	@echo "🔨 Building frontend..."
	cd frontend && npm run build
	@echo "✅ Frontend built: frontend/dist"

# Build both
build: build-backend build-frontend
	@echo "✅ Build complete"

# Docker commands
docker-up:
	@echo "🐳 Starting Docker containers..."
	docker-compose up -d
	@echo "✅ Containers started"
	@echo "Backend: http://localhost:8080"
	@echo "Frontend: http://localhost:5173"

docker-down:
	@echo "🐳 Stopping Docker containers..."
	docker-compose down
	@echo "✅ Containers stopped"

docker-logs:
	docker-compose logs -f

# Start kubectl proxy
kubectl-proxy:
	@echo "🔧 Starting kubectl proxy..."
	kubectl proxy --port=8001

# Run tests
test:
	@echo "🧪 Running backend tests..."
	cd backend && go test ./...
	@echo "✅ Tests complete"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -f backend/arlog-backend
	rm -rf frontend/dist
	rm -rf frontend/node_modules/.vite
	@echo "✅ Clean complete"

# Database setup
db-create:
	@echo "🗄️  Creating database..."
	psql -U postgres -c "CREATE DATABASE arlog_db;"
	psql -U postgres -c "CREATE USER arlog WITH PASSWORD 'arlog_password';"
	psql -U postgres -c "GRANT ALL PRIVILEGES ON DATABASE arlog_db TO arlog;"
	@echo "✅ Database created"

db-drop:
	@echo "🗄️  Dropping database..."
	psql -U postgres -c "DROP DATABASE IF EXISTS arlog_db;"
	psql -U postgres -c "DROP USER IF EXISTS arlog;"
	@echo "✅ Database dropped"

db-reset: db-drop db-create
	@echo "✅ Database reset complete"

# Check prerequisites
check:
	@echo "Checking prerequisites..."
	@command -v go >/dev/null 2>&1 || { echo "❌ Go is not installed"; exit 1; }
	@command -v node >/dev/null 2>&1 || { echo "❌ Node.js is not installed"; exit 1; }
	@command -v psql >/dev/null 2>&1 || { echo "❌ PostgreSQL is not installed"; exit 1; }
	@command -v kubectl >/dev/null 2>&1 || { echo "❌ kubectl is not installed"; exit 1; }
	@echo "✅ All prerequisites installed"
	@go version
	@node --version
	@psql --version
	@kubectl version --client

# Setup environment
setup:
	@echo "🔧 Setting up ArLOG..."
	@if [ ! -f backend/.env ]; then \
		cp backend/.env.example backend/.env; \
		echo "⚠️  Please update backend/.env with your configuration"; \
	fi
	@echo "✅ Setup complete"

