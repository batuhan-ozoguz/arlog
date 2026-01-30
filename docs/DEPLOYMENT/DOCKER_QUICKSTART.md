# Docker Quick Start for ArLOG

## 🎯 Super Quick Start (3 Steps)

```bash
# 1. Start Docker Desktop
open -a Docker
# Wait 30 seconds for Docker to start...

# 2. Run startup script
cd /Users/bozoguz/ArLog
./START_DOCKER.sh

# 3. Open dashboard
open http://localhost:3000
```

**That's it!** All services running in Docker. 🐳

---

## 📊 What Gets Started

| Service | Container | Port | URL |
|---------|-----------|------|-----|
| PostgreSQL | arlog-postgres | 5432 | localhost:5432 |
| Go Backend | arlog-backend | 8080 | http://localhost:8080 |
| Next.js Dashboard | arlog-dashboard | 3000 | http://localhost:3000 |

---

## 🔧 Common Commands

### Start Services
```bash
./START_DOCKER.sh
```

### Stop Services
```bash
./STOP_DOCKER.sh
```

### View Logs
```bash
# All containers
docker logs -f arlog-backend
docker logs -f arlog-dashboard
docker logs -f arlog-postgres

# Last 50 lines
docker logs --tail 50 arlog-backend
```

### Check Status
```bash
# List running containers
docker ps

# Check specific container
docker ps | grep arlog
```

### Restart Service
```bash
# Restart backend after code changes
docker stop arlog-backend
docker rm arlog-backend
./START_DOCKER.sh  # Rebuilds and restarts
```

### Access Database
```bash
# Connect to PostgreSQL
docker exec -it arlog-postgres psql -U arlog -d arlog_db

# View tables
docker exec -it arlog-postgres psql -U arlog -d arlog_db -c "\dt"

# Query teams
docker exec -it arlog-postgres psql -U arlog -d arlog_db -c "SELECT * FROM teams;"
```

---

## 🐛 Troubleshooting

### Docker Desktop Not Running
```
Error: Cannot connect to the Docker daemon
```

**Fix:**
```bash
open -a Docker
# Wait 30 seconds, then try again
```

### Port Already in Use
```
Error: port 3000 is already allocated
```

**Fix:**
```bash
# Stop the v0 dev server first
# Find the process
lsof -i :3000

# Kill it
kill -9 <PID>

# Or stop all Node processes
pkill -f "next dev"
```

### Container Won't Start
```
Container keeps restarting
```

**Fix:**
```bash
# Check logs
docker logs arlog-backend

# Remove and rebuild
docker stop arlog-backend
docker rm arlog-backend
cd backend
docker build --no-cache -t arlog-backend .
```

### Database Connection Failed
```
failed to connect to postgres
```

**Fix:**
```bash
# Check postgres is running
docker ps | grep postgres

# Check postgres logs
docker logs arlog-postgres

# Restart postgres
docker restart arlog-postgres

# Wait a few seconds, then restart backend
docker restart arlog-backend
```

---

## 🔄 After Code Changes

### Backend Changes
```bash
# Stop and remove container
docker stop arlog-backend
docker rm arlog-backend

# Rebuild
cd backend
docker build -t arlog-backend .

# Restart (or just run START_DOCKER.sh)
docker run -d \
  --name arlog-backend \
  --network arlog-network \
  -e DB_HOST=arlog-postgres \
  -e AUTH_MODE=dev \
  -p 8080:8080 \
  arlog-backend
```

### Dashboard Changes
```bash
# Stop and remove
docker stop arlog-dashboard
docker rm arlog-dashboard

# Rebuild
cd frontend/arlog-dashboard
docker build -t arlog-dashboard .

# Restart
docker run -d \
  --name arlog-dashboard \
  --network arlog-network \
  -e NEXT_PUBLIC_API_URL=http://localhost:8080 \
  -p 3000:3000 \
  arlog-dashboard
```

---

## 📋 First Time Setup Checklist

- [ ] Docker Desktop installed and running
- [ ] Run `./START_DOCKER.sh`
- [ ] Wait for all services to start (30-60 seconds)
- [ ] Check logs: `docker logs arlog-backend`
- [ ] Access dashboard: http://localhost:3000
- [ ] Test health: http://localhost:8080/health

---

## 🎨 What You'll See

1. **Database**: Auto-creates tables and seeds test data
2. **Backend**: Starts with dev auth mode, connects to DB
3. **Dashboard**: Beautiful UI with mock project data

---

## 🚀 Production Deployment

For server deployment, update `.env`:

```env
AUTH_MODE=okta
OKTA_DOMAIN=your-domain.okta.com
OKTA_CLIENT_ID=your-client-id
OKTA_CLIENT_SECRET=your-secret
JWT_SECRET=strong-production-secret
ENVIRONMENT=production
```

Then rebuild:
```bash
./START_DOCKER.sh
```

---

## ⚡ Quick Commands

```bash
# Start everything
./START_DOCKER.sh

# Stop everything
./STOP_DOCKER.sh

# View all logs
docker logs -f arlog-backend &
docker logs -f arlog-dashboard &

# Restart after changes
docker restart arlog-backend

# Reset database
docker volume rm arlog_postgres_data
./START_DOCKER.sh
```

---

## 📱 Access Points

- **Dashboard**: http://localhost:3000 (Main UI)
- **Backend API**: http://localhost:8080/api
- **Health Check**: http://localhost:8080/health
- **WebSocket**: ws://localhost:8080/ws/logs

---

Ready to run! Just start Docker Desktop and execute `./START_DOCKER.sh` 🚀

