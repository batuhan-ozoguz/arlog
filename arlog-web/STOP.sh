#!/bin/bash

# ArLOG Docker Stop Script
echo "🛑 Stopping ArLOG Docker services..."
echo ""

# Stop all containers
echo "Stopping containers..."
docker stop arlog-postgres arlog-backend arlog-dashboard 2>/dev/null || true

# Remove containers
echo "Removing containers..."
docker rm arlog-postgres arlog-backend arlog-dashboard 2>/dev/null || true

echo ""
echo "✅ All services stopped"
echo ""
echo "To start again, run:"
echo "  ./START_DOCKER.sh"
echo ""
echo "To remove database volume (reset data):"
echo "  docker volume rm arlog_postgres_data"
echo ""

