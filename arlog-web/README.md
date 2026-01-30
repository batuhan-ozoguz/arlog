# ArLOG Web - Multi-User Web Application

Server-based Kubernetes log viewer with team authentication and multi-project support.

## 🎯 Overview

ArLOG Web is the team-oriented version featuring:
- Multi-user support with Okta SSO
- Team-based access control
- Multi-project/multi-cluster management
- PostgreSQL database for permissions
- WebSocket log streaming
- Docker deployment

## 🏗️ Architecture

```
Browser (Next.js)
    ↓
Backend API (Go)
    ↓
PostgreSQL + Okta SSO
    ↓
Kubernetes API (Service Account Tokens)
```

## 🚀 Quick Start with Docker

```bash
# Start all services (PostgreSQL + Backend + Frontend)
./START.sh

# Access
open http://localhost:3000

# Stop
./STOP.sh
```

## 📁 Structure

```
arlog-web/
├── backend/                 # Go API server
│   ├── main.go
│   ├── models/             # Database models
│   ├── handlers/           # API endpoints
│   ├── services/           # Kubernetes integration
│   └── Dockerfile
│
├── frontend/               # Next.js application
│   ├── app/               # Pages
│   ├── components/        # UI components
│   ├── lib/              # Utilities
│   └── Dockerfile
│
├── docker-compose.yml     # Docker orchestration
├── START.sh              # Start script
└── STOP.sh               # Stop script
```

## 🔧 Manual Setup

### Backend:
```bash
cd backend
cp .env.example .env
# Edit .env with database and Okta config
go mod download
go run main.go
```

### Frontend:
```bash
cd frontend
npm install
npm run dev
```

### Database:
```bash
# With Docker
docker run -d \
  --name arlog-postgres \
  -e POSTGRES_USER=arlog \
  -e POSTGRES_PASSWORD=arlog_password \
  -e POSTGRES_DB=arlog_db \
  -p 5432:5432 \
  postgres:15-alpine
```

## 🔐 Authentication

### Dev Mode (No Okta):
```env
AUTH_MODE=dev
```
Automatically creates a dummy user. Perfect for testing!

### Production Mode (With Okta):
```env
AUTH_MODE=okta
OKTA_DOMAIN=your-domain.okta.com
OKTA_CLIENT_ID=your-client-id
OKTA_CLIENT_SECRET=your-secret
```

See [AUTHENTICATION.md](../AUTHENTICATION.md) for details.

## 🎨 Features

- ✅ Project selection page
- ✅ Namespace browsing
- ✅ Pod listing with status
- ✅ Real-time log streaming via WebSocket
- ✅ Pause/Resume/Download logs
- ✅ Team-based permissions
- ✅ Multi-project support
- ✅ Beautiful minimal UI

## 📊 API Endpoints

- `GET /api/user/projects` - List user's projects
- `GET /api/user/permissions` - List accessible namespaces
- `GET /api/pods?namespace=<ns>` - List pods
- `WS /ws/logs?namespace=<ns>&podName=<pod>` - Stream logs
- `GET /auth/okta/login` - SSO login
- `GET /health` - Health check

## 🐳 Docker Deployment

```bash
# Build and start
docker-compose up -d --build

# View logs
docker-compose logs -f backend
docker-compose logs -f frontend

# Stop
docker-compose down
```

## 🚀 Production Deployment

1. Configure environment:
   ```bash
   vim .env
   # Set AUTH_MODE=okta
   # Add Okta credentials
   # Set strong JWT_SECRET
   ```

2. Deploy with Docker:
   ```bash
   ./START.sh
   ```

3. Set up reverse proxy (nginx/traefik)
4. Configure SSL certificate
5. Set up monitoring

## 🆘 Troubleshooting

**Backend won't start:**
```bash
docker logs arlog-backend
# Check database connection
```

**Frontend won't connect:**
```bash
# Check backend is running
curl http://localhost:8080/health
```

**Database issues:**
```bash
# Reset database
docker-compose down -v
docker-compose up -d
```

## 📚 Documentation

- [`README.md`](./README.md) - This file
- [`backend/README.md`](./backend/README.md) - Backend details
- [`QUICKSTART.md`](../QUICKSTART.md) - Quick start guide
- [`../docs/`](../docs/) - Full documentation

---

For team collaboration and production use! 🌐



