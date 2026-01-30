# ArLOG Setup Guide

Complete step-by-step guide to set up ArLOG for development and production.

## Table of Contents
1. [Prerequisites](#prerequisites)
2. [Database Setup](#database-setup)
3. [Okta SSO Configuration](#okta-sso-configuration)
4. [Backend Setup](#backend-setup)
5. [Frontend Setup](#frontend-setup)
6. [Kubernetes Configuration](#kubernetes-configuration)
7. [Running the Application](#running-the-application)
8. [Docker Setup](#docker-setup)

---

## Prerequisites

### Required Software
- **Go**: Version 1.21 or higher
  ```bash
  go version
  ```

- **Node.js**: Version 18 or higher
  ```bash
  node --version
  npm --version
  ```

- **PostgreSQL**: Version 13 or higher
  ```bash
  psql --version
  ```

- **Kubernetes**: Access to a Kubernetes cluster
  ```bash
  kubectl version
  ```

- **Docker** (optional): For containerized deployment
  ```bash
  docker --version
  docker-compose --version
  ```

---

## Database Setup

### Option 1: Local PostgreSQL

1. Install PostgreSQL (if not already installed):
   ```bash
   # macOS
   brew install postgresql@15
   brew services start postgresql@15

   # Ubuntu/Debian
   sudo apt-get install postgresql-15

   # Windows
   # Download installer from https://www.postgresql.org/download/windows/
   ```

2. Create database and user:
   ```bash
   psql -U postgres
   ```

   In the PostgreSQL prompt:
   ```sql
   CREATE DATABASE arlog_db;
   CREATE USER arlog WITH PASSWORD 'arlog_password';
   GRANT ALL PRIVILEGES ON DATABASE arlog_db TO arlog;
   \q
   ```

3. Verify connection:
   ```bash
   psql -U arlog -d arlog_db -h localhost
   ```

### Option 2: Docker PostgreSQL

```bash
docker run -d \
  --name arlog-postgres \
  -e POSTGRES_USER=arlog \
  -e POSTGRES_PASSWORD=arlog_password \
  -e POSTGRES_DB=arlog_db \
  -p 5432:5432 \
  postgres:15-alpine
```

---

## Okta SSO Configuration

### 1. Create Okta Developer Account
- Go to https://developer.okta.com/signup/
- Sign up for a free developer account

### 2. Create an Application

1. Log in to Okta Admin Console
2. Navigate to **Applications** → **Applications**
3. Click **Create App Integration**
4. Select:
   - **Sign-in method**: OIDC - OpenID Connect
   - **Application type**: Web Application
5. Click **Next**

### 3. Configure Application

**General Settings:**
- **App integration name**: ArLOG
- **Grant type**: Authorization Code
- **Sign-in redirect URIs**: 
  - `http://localhost:8080/auth/okta/callback` (development)
  - `https://your-domain.com/auth/okta/callback` (production)
- **Sign-out redirect URIs**: 
  - `http://localhost:5173` (development)
  - `https://your-domain.com` (production)

**Assignments:**
- Choose who can access this application
- Assign to groups or everyone

### 4. Note Configuration Details

After creating the application, note:
- **Client ID**
- **Client Secret**
- **Okta Domain** (e.g., `dev-12345.okta.com`)

### 5. Configure Groups (Optional)

1. Navigate to **Directory** → **Groups**
2. Create groups matching your teams (e.g., "Cosmos Team", "Jupiter Team")
3. Add users to groups
4. Note the **Group IDs** for database configuration

---

## Backend Setup

### 1. Navigate to Backend Directory
```bash
cd backend
```

### 2. Install Dependencies
```bash
go mod download
```

### 3. Configure Environment Variables

Copy the example environment file:
```bash
cp .env.example .env
```

Edit `.env` with your configuration:
```env
# Server Configuration
PORT=8080
SERVER_HOST=localhost

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=arlog
DB_PASSWORD=arlog_password
DB_NAME=arlog_db
DB_SSLMODE=disable

# Kubernetes Configuration
KUBE_PROXY_URL=http://localhost:8001

# Okta Configuration
OKTA_DOMAIN=dev-12345.okta.com
OKTA_CLIENT_ID=0oa...your-client-id
OKTA_CLIENT_SECRET=your-client-secret
OKTA_REDIRECT_URI=http://localhost:8080/auth/okta/callback

# JWT Secret
JWT_SECRET=your-super-secret-jwt-key-change-this

# Environment
ENVIRONMENT=development

# Frontend URL
FRONTEND_URL=http://localhost:5173
```

### 4. Run Database Migrations

The application will automatically run migrations on startup. To manually run:
```bash
go run main.go
```

The migrations will create:
- `teams` table
- `permissions` table

### 5. Seed Test Data (Development Only)

The application automatically seeds test data in development mode:
- Cosmos Team with dev/test namespaces
- Jupiter Team with dev namespace

---

## Frontend Setup

### 1. Navigate to Frontend Directory
```bash
cd frontend
```

### 2. Install Dependencies
```bash
npm install
```

### 3. Configure Environment (Optional)

Create `.env` file:
```env
VITE_API_URL=http://localhost:8080
```

---

## Kubernetes Configuration

### Option 1: Local Development with kubectl proxy

1. Start kubectl proxy:
   ```bash
   kubectl proxy --port=8001
   ```

2. The backend will connect via `http://localhost:8001`

3. Ensure you have access to the cluster:
   ```bash
   kubectl get pods --all-namespaces
   ```

### Option 2: Service Account Tokens (Production)

1. Create a service account in Kubernetes:
   ```bash
   kubectl create serviceaccount arlog-viewer -n default
   ```

2. Create a role with read permissions:
   ```bash
   kubectl create clusterrole pod-reader \
     --verb=get,list,watch \
     --resource=pods,pods/log
   ```

3. Bind the role:
   ```bash
   kubectl create clusterrolebinding arlog-viewer-binding \
     --clusterrole=pod-reader \
     --serviceaccount=default:arlog-viewer
   ```

4. Get the service account token:
   ```bash
   # For Kubernetes < 1.24
   kubectl get secret $(kubectl get sa arlog-viewer -o jsonpath='{.secrets[0].name}') -o jsonpath='{.data.token}' | base64 -d

   # For Kubernetes >= 1.24
   kubectl create token arlog-viewer --duration=8760h
   ```

5. Store the token in the `permissions` table for the appropriate namespace

---

## Running the Application

### 1. Start Backend

In the `backend/` directory:
```bash
go run main.go
```

Expected output:
```
🚀 Starting ArLOG Backend Server...
✅ Database connection established successfully
🔄 Running database migrations...
✅ Database migrations completed successfully
🌱 Seeding database with test data...
✅ Database seeded successfully
✅ Server is running on http://localhost:8080
📡 API endpoints available at http://localhost:8080/api
🔌 WebSocket endpoint available at ws://localhost:8080/ws
```

### 2. Start Frontend

In the `frontend/` directory:
```bash
npm run dev
```

Expected output:
```
  VITE v5.0.8  ready in 500 ms

  ➜  Local:   http://localhost:5173/
  ➜  Network: use --host to expose
  ➜  press h to show help
```

### 3. Access the Application

Open your browser and navigate to:
```
http://localhost:5173
```

### 4. Login

Click "Sign in with Okta" and authenticate with your Okta credentials.

---

## Docker Setup

### 1. Using Docker Compose

1. Copy environment file:
   ```bash
   cp .env.example .env
   ```

2. Update `.env` with your Okta credentials

3. Start services:
   ```bash
   docker-compose up -d
   ```

4. View logs:
   ```bash
   docker-compose logs -f
   ```

5. Stop services:
   ```bash
   docker-compose down
   ```

### 2. Building Individual Containers

Backend:
```bash
cd backend
docker build -t arlog-backend .
docker run -p 8080:8080 --env-file .env arlog-backend
```

Frontend:
```bash
cd frontend
docker build -t arlog-frontend .
docker run -p 80:80 arlog-frontend
```

---

## Troubleshooting

### Database Connection Issues
```bash
# Check if PostgreSQL is running
pg_isready -h localhost -p 5432

# Test connection
psql -U arlog -d arlog_db -h localhost
```

### Backend Won't Start
```bash
# Check Go version
go version

# Verify dependencies
go mod verify

# Check environment variables
cat .env
```

### Frontend Won't Start
```bash
# Clear node_modules
rm -rf node_modules package-lock.json
npm install

# Check Node version
node --version
```

### Kubernetes Connection Issues
```bash
# Verify kubectl proxy is running
curl http://localhost:8001/api/v1/namespaces

# Check kubeconfig
kubectl config view
kubectl cluster-info
```

### Okta Authentication Issues
- Verify redirect URIs match exactly
- Check Okta domain format (no https://)
- Ensure client ID and secret are correct
- Check that users are assigned to the application

---

## Next Steps

1. **Configure Team Permissions**: Add teams and their namespace permissions in the database
2. **Test Log Streaming**: Access a namespace and view pod logs
3. **Customize UI**: Modify the frontend styling as needed
4. **Production Deployment**: Follow security best practices for production

---

## Support

For issues or questions:
1. Check the main README.md
2. Review backend and frontend README files
3. Check application logs
4. Verify all environment variables are set correctly

