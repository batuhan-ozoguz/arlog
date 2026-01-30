#!/bin/bash

# ArLOG Web - Docker Startup Script
echo "🚀 Starting ArLOG Web Application..."
echo ""

# Check if Docker is running
if ! docker ps > /dev/null 2>&1; then
    echo "❌ Docker daemon is not running!"
    echo ""
    echo "Please start Docker Desktop:"
    echo "  1. Open Docker Desktop application"
    echo "  2. Wait for it to start (20-30 seconds)"
    echo "  3. Run this script again"
    echo ""
    echo "Or run manually:"
    echo "  open -a Docker"
    exit 1
fi

echo "✅ Docker is running"
echo ""

# Create network
echo "📡 Creating Docker network..."
docker network create arlog-network 2>/dev/null || echo "   Network already exists"

# Stop any existing containers
echo "🛑 Stopping existing containers..."
docker stop arlog-postgres arlog-backend arlog-dashboard 2>/dev/null || true
docker rm arlog-postgres arlog-backend arlog-dashboard 2>/dev/null || true

# Start PostgreSQL
echo "🗄️  Starting PostgreSQL..."
docker run -d \
  --name arlog-postgres \
  --network arlog-network \
  -e POSTGRES_USER=arlog \
  -e POSTGRES_PASSWORD=arlog_password \
  -e POSTGRES_DB=arlog_db \
  -p 5432:5432 \
  -v arlog_postgres_data:/var/lib/postgresql/data \
  postgres:15-alpine

echo "   Waiting for PostgreSQL to be ready..."
sleep 5

# Build and start Backend
echo "🔧 Building and starting Go Backend..."
cd "$(dirname "$0")/backend"
docker build -t arlog-backend .
docker run -d \
  --name arlog-backend \
  --network arlog-network \
  -e PORT=8080 \
  -e DB_HOST=arlog-postgres \
  -e DB_PORT=5432 \
  -e DB_USER=arlog \
  -e DB_PASSWORD=arlog_password \
  -e DB_NAME=arlog_db \
  -e DB_SSLMODE=disable \
  -e ENVIRONMENT=development \
  -e AUTH_MODE=dev \
  -e JWT_SECRET=dev-secret-key \
  -e FRONTEND_URL=http://localhost:3000 \
  -p 8080:8080 \
  arlog-backend

cd ..

echo "   Waiting for backend to start..."
sleep 3

# Build and start Next.js Dashboard
echo "🎨 Building and starting Next.js Dashboard..."
cd "$(dirname "$0")/frontend"
docker build -t arlog-dashboard .
docker run -d \
  --name arlog-dashboard \
  --network arlog-network \
  -e NEXT_PUBLIC_API_URL=http://localhost:8080 \
  -e NEXT_PUBLIC_ENVIRONMENT=development \
  -e NEXT_PUBLIC_AUTH_MODE=dev \
  -p 3000:3000 \
  arlog-dashboard

cd ../..

echo ""
echo "✅ All services started!"
echo ""
echo "📊 Service URLs:"
echo "   Dashboard:  http://localhost:3000"
echo "   Backend:    http://localhost:8080"
echo "   Health:     http://localhost:8080/health"
echo "   PostgreSQL: localhost:5432"
echo ""
echo "📝 Useful commands:"
echo "   docker logs -f arlog-backend     # View backend logs"
echo "   docker logs -f arlog-dashboard   # View dashboard logs"
echo "   docker ps                        # View running containers"
echo "   docker stop arlog-backend        # Stop backend"
echo ""
echo "To stop all services:"
echo "   docker stop arlog-postgres arlog-backend arlog-dashboard"
echo ""

