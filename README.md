# ArLOG - Enterprise Kubernetes Log Viewer

<div align="center">

![Version](https://img.shields.io/badge/version-1.1.0-blue.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)
![TypeScript](https://img.shields.io/badge/TypeScript-5.5+-3178C6?logo=typescript)
![Kubernetes](https://img.shields.io/badge/Kubernetes-1.29+-326CE5?logo=kubernetes)

**Production-ready Kubernetes log management platform with dual deployment options: Desktop application for developers and Web application for teams.**

[Features](#-features) • [Architecture](#-architecture) • [Quick Start](#-quick-start) • [Documentation](#-documentation) • [Contributing](#-contributing)

</div>

---

## 🎯 Overview

ArLOG is a comprehensive Kubernetes log viewing solution designed for both individual developers and enterprise teams. It provides real-time log streaming, multi-cluster support, and enterprise-grade authentication with a beautiful, minimal UI.

### Two Deployment Models

| **Desktop Application (V1)** | **Web Application (V2)** |
|-------------------------------|--------------------------|
| Standalone Electron app | Server-based multi-user platform |
| Uses local `~/.kube/config` | Service account token management |
| Single-user, offline capable | Team collaboration with SSO |
| Cross-platform installers | Docker containerized deployment |
| Zero server setup | PostgreSQL + Okta integration |

---

## ✨ Features

### Core Capabilities
- 🔄 **Real-time Log Streaming** - WebSocket-based live log streaming with pause/resume
- 🌐 **Multi-Cluster Support** - Seamless context switching between Kubernetes clusters
- 📊 **Namespace & Pod Management** - Browse namespaces, view pod status, and stream logs
- 🎨 **Modern UI** - Beautiful, minimal interface built with shadcn/ui and Tailwind CSS
- 🔐 **Enterprise Authentication** - Okta SSO integration with team-based permissions
- 📦 **Container-Ready** - Full Docker Compose setup for production deployment
- 🚀 **CI/CD Integration** - Azure Pipelines for automated builds and releases

### Desktop Application
- ✅ Cross-platform installers (macOS, Windows, Linux)
- ✅ Offline operation with local kubeconfig
- ✅ Multi-container pod support
- ✅ Log download and export
- ✅ Auto-scroll and connection status indicators

### Web Application
- ✅ Multi-project/multi-cluster management
- ✅ PostgreSQL database for permissions and configuration
- ✅ Team-based access control
- ✅ WebSocket log streaming
- ✅ RESTful API with JWT authentication

---

## 🏗️ Architecture

### Desktop Application Architecture

```
┌─────────────────────────────────────────┐
│         Electron Main Process           │
│  ┌───────────────────────────────────┐ │
│  │  Kubernetes Client (TypeScript)   │ │
│  │  - @kubernetes/client-node        │ │
│  │  - Context Management             │ │
│  │  - Log Streaming                  │ │
│  └───────────────────────────────────┘ │
│              ↕ IPC Bridge                │
│  ┌───────────────────────────────────┐ │
│  │  React Renderer Process           │ │
│  │  - React 18 + TypeScript          │ │
│  │  - Tailwind CSS + shadcn/ui       │ │
│  │  - Real-time UI Updates           │ │
│  └───────────────────────────────────┘ │
└─────────────────────────────────────────┘
              ↕
    ~/.kube/config → Kubernetes API
```

### Web Application Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    Browser (Next.js)                     │
│              React 18 + TypeScript + Tailwind            │
└──────────────────────┬──────────────────────────────────┘
                       │ HTTPS/REST API
┌──────────────────────▼──────────────────────────────────┐
│              Go Backend API Server                       │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐ │
│  │   Handlers   │  │  Middleware  │  │   Services   │ │
│  │  - REST API  │  │  - JWT Auth  │  │  - K8s Client│ │
│  │  - WebSocket │  │  - Okta SSO  │  │  - Log Stream│ │
│  └──────────────┘  └──────────────┘  └──────────────┘ │
└──────────┬──────────────────────┬───────────────────────┘
           │                      │
    ┌──────▼──────┐      ┌───────▼────────┐
    │ PostgreSQL  │      │  Kubernetes API │
    │ - Teams      │      │  - client-go   │
    │ - Permissions│      │  - Service Acct│
    │ - Projects   │      │  - Multi-cluster│
    └─────────────┘      └────────────────┘
```

---

## 🛠️ Technology Stack

### Desktop Application
| Layer | Technology | Purpose |
|-------|-----------|---------|
| **Framework** | Electron 32 | Cross-platform desktop runtime |
| **UI Library** | React 18 | Component-based UI |
| **Language** | TypeScript 5.5 | Type-safe development |
| **Build Tool** | Vite + electron-vite | Fast builds and HMR |
| **Styling** | Tailwind CSS 3.4 | Utility-first CSS |
| **UI Components** | shadcn/ui + Radix UI | Accessible component library |
| **K8s Client** | @kubernetes/client-node | Kubernetes API integration |
| **Package Manager** | npm/pnpm | Dependency management |

### Web Application - Backend
| Layer | Technology | Purpose |
|-------|-----------|---------|
| **Language** | Go 1.21+ | High-performance backend |
| **HTTP Router** | Gorilla Mux | RESTful API routing |
| **Database ORM** | GORM | PostgreSQL abstraction |
| **Database** | PostgreSQL 15 | Relational data storage |
| **K8s Client** | k8s.io/client-go v0.29 | Kubernetes API client |
| **WebSocket** | Gorilla WebSocket | Real-time log streaming |
| **Authentication** | Okta OIDC + JWT | Enterprise SSO |
| **Container** | Docker | Containerization |

### Web Application - Frontend
| Layer | Technology | Purpose |
|-------|-----------|---------|
| **Framework** | Next.js 15 | React framework with SSR |
| **Language** | TypeScript 5 | Type-safe frontend |
| **Styling** | Tailwind CSS 4 | Utility-first styling |
| **UI Components** | shadcn/ui | Accessible components |
| **State Management** | React Context API | Global state |
| **HTTP Client** | Fetch API | REST API communication |

### Infrastructure & DevOps
| Tool | Purpose |
|------|---------|
| **Docker** | Containerization |
| **Docker Compose** | Multi-container orchestration |
| **Azure Pipelines** | CI/CD automation |
| **PostgreSQL** | Database persistence |
| **Kubernetes** | Container orchestration (target platform) |

---

## 🚀 Quick Start

### Desktop Application

```bash
# Clone repository
git clone https://github.com/batuhan-ozoguz/arlog.git
cd ArLog/arlog-desktop

# Install dependencies
npm install

# Run in development mode
npm run dev

# Build production installer
npm run build:mac    # macOS .dmg
npm run build:win    # Windows .exe
npm run build:linux  # Linux .AppImage
```

**Prerequisites:**
- Node.js 18+ and npm/pnpm
- kubectl configured with at least one context
- Access to Kubernetes cluster

### Web Application (Docker)

```bash
# Navigate to web application directory
cd arlog-web

# Start all services (PostgreSQL + Backend + Frontend)
./START.sh

# Access application
# Frontend: http://localhost:3000
# Backend API: http://localhost:8080
# Health Check: http://localhost:8080/health

# Stop services
./STOP.sh
```

**Prerequisites:**
- Docker and Docker Compose
- kubectl proxy (for local development) or Kubernetes service account tokens

### Manual Setup (Web Application)

#### Backend Setup
```bash
cd arlog-web/backend

# Create .env file
cat > .env << EOF
DB_HOST=localhost
DB_PORT=5432
DB_USER=arlog
DB_PASSWORD=arlog_password
DB_NAME=arlog_db
DB_SSLMODE=disable
PORT=8080
AUTH_MODE=dev
JWT_SECRET=your-secret-key-change-in-production
OKTA_DOMAIN=${OKTA_DOMAIN}
OKTA_CLIENT_ID=${OKTA_CLIENT_ID}
OKTA_CLIENT_SECRET=${OKTA_CLIENT_SECRET}
EOF

# Install dependencies and run
go mod download
go run main.go
```

#### Frontend Setup
```bash
cd arlog-web/frontend

# Install dependencies
npm install

# Run development server
npm run dev
```

---

## 🔧 Kubernetes Integration

### Desktop Application
The desktop application uses the local Kubernetes configuration file (`~/.kube/config`) to authenticate and connect to clusters.

**Features:**
- Automatic context detection
- Dynamic context switching
- Direct API communication using `@kubernetes/client-node`
- No additional authentication required

**Example Usage:**
```typescript
import { KubernetesClient } from './kubernetes'

const k8s = new KubernetesClient()
const contexts = k8s.getContexts()
k8s.setCurrentContext('production-cluster')
const namespaces = await k8s.listNamespaces()
```

### Web Application
The web application uses Kubernetes service account tokens for secure, scalable access to clusters.

**Service Account Setup:**
```bash
# Create service account
kubectl create serviceaccount arlog-viewer -n default

# Create cluster role with read permissions
kubectl create clusterrole pod-reader \
  --verb=get,list,watch \
  --resource=pods,pods/log,namespaces

# Bind role to service account
kubectl create clusterrolebinding arlog-viewer-binding \
  --clusterrole=pod-reader \
  --serviceaccount=default:arlog-viewer

# Get service account token
kubectl create token arlog-viewer --duration=8760h
```

**Go Implementation:**
```go
import (
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
)

config, _ := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
clientset, _ := kubernetes.NewForConfig(config)
```

---

## 📦 Docker Deployment

### Production Deployment

```bash
# Build and start services
cd arlog-web
docker-compose up -d --build

# View logs
docker-compose logs -f backend
docker-compose logs -f frontend

# Scale backend (if needed)
docker-compose up -d --scale backend=3

# Stop services
docker-compose down
```

### Docker Compose Services

- **postgres**: PostgreSQL 15 database with health checks
- **backend**: Go API server with auto-migration
- **frontend**: Next.js application with SSR
- **dashboard**: Alternative frontend (Vite + React)

### Environment Variables

Key environment variables for production:

```env
# Database
DB_HOST=postgres
DB_PORT=5432
DB_USER=arlog
DB_PASSWORD=<secure-password>
DB_NAME=arlog_db
DB_SSLMODE=require

# Authentication
AUTH_MODE=okta
OKTA_DOMAIN=your-domain.okta.com
OKTA_CLIENT_ID=<client-id>
OKTA_CLIENT_SECRET=<client-secret>
JWT_SECRET=<strong-random-secret>

# Kubernetes
KUBE_PROXY_URL=http://kube-proxy:8001
# OR use service account tokens stored in database
```

---

## 🔐 Authentication & Security

### Desktop Application
- Uses local `~/.kube/config` credentials
- No additional authentication layer
- Secure credential storage via OS keychain

### Web Application

#### Development Mode
```env
AUTH_MODE=dev
```
- Automatic dummy user creation
- No authentication required
- **⚠️ For development only**

#### Production Mode (Okta SSO)
```env
AUTH_MODE=okta
OKTA_DOMAIN=your-domain.okta.com
OKTA_CLIENT_ID=<client-id>
OKTA_CLIENT_SECRET=<client-secret>
```

**Flow:**
1. User clicks "Login with Okta"
2. Redirected to Okta OIDC endpoint
3. After authentication, redirected back with authorization code
4. Backend exchanges code for tokens
5. JWT token issued and stored in HTTP-only cookie
6. Subsequent requests authenticated via JWT middleware

**Security Features:**
- JWT token-based authentication
- HTTP-only cookies for token storage
- CORS protection
- Service account token encryption in database
- Team-based permission system

---

## 📊 API Documentation

### REST Endpoints

#### Authentication
- `GET /auth/okta/login` - Initiate Okta SSO login
- `GET /auth/okta/callback` - Okta callback handler
- `POST /auth/logout` - Logout and clear session

#### Projects & Permissions
- `GET /api/user/projects` - List user's accessible projects
- `GET /api/user/permissions` - Get user's namespace permissions
- `GET /api/projects/:id/namespaces` - List namespaces for project

#### Kubernetes Resources
- `GET /api/pods?namespace=<ns>` - List pods in namespace
- `GET /api/namespaces` - List accessible namespaces
- `GET /api/pods/:name/logs?namespace=<ns>` - Get pod logs (REST)

#### WebSocket
- `WS /ws/logs?namespace=<ns>&podName=<pod>&container=<container>` - Stream logs in real-time

#### Health & Status
- `GET /health` - Health check endpoint
- `GET /api/status` - Application status

### Example API Usage

```bash
# Get user projects
curl -H "Authorization: Bearer <jwt-token>" \
  http://localhost:8080/api/user/projects

# List pods in namespace
curl -H "Authorization: Bearer <jwt-token>" \
  "http://localhost:8080/api/pods?namespace=default"

# Stream logs via WebSocket
wscat -c "ws://localhost:8080/ws/logs?namespace=default&podName=my-pod"
```

---

## 🏭 CI/CD Pipeline

### Azure Pipelines

The project includes Azure Pipelines configuration for automated builds and releases.

**Pipeline Stages:**
1. **Build** - Compile and test applications
2. **Package** - Create distributable artifacts
3. **Release** - Deploy to staging/production

**Configuration:** [`arlog-desktop/azure-pipelines.yml`](./arlog-desktop/azure-pipelines.yml)

**Build Artifacts:**
- macOS: `.dmg` installer
- Windows: `.exe` installer
- Linux: `.AppImage` or `.deb` package

---

## 📁 Project Structure

```
ArLog/
├── arlog-desktop/              # Desktop Application (Electron)
│   ├── src/
│   │   ├── main/               # Electron main process
│   │   │   ├── index.ts        # Main entry point
│   │   │   ├── kubernetes.ts   # K8s client implementation
│   │   │   └── ipc-handlers.ts # IPC communication
│   │   ├── preload/            # IPC bridge
│   │   └── renderer/           # React UI
│   │       └── src/
│   │           ├── components/ # UI components
│   │           ├── contexts/   # React contexts
│   │           ├── pages/      # Page components
│   │           └── lib/        # Utilities
│   ├── build/                  # Build resources
│   ├── dist/                   # Distribution artifacts
│   ├── package.json
│   └── electron.vite.config.ts
│
├── arlog-web/                  # Web Application
│   ├── backend/                # Go API Server
│   │   ├── main.go             # Entry point
│   │   ├── handlers/           # HTTP handlers
│   │   ├── services/           # Business logic
│   │   ├── models/             # Database models
│   │   ├── middleware/         # Auth middleware
│   │   ├── database/           # DB connection
│   │   ├── utils/              # Utilities
│   │   ├── Dockerfile
│   │   └── go.mod
│   │
│   ├── frontend/               # Next.js Frontend
│   │   ├── app/                # Next.js app router
│   │   ├── components/         # React components
│   │   ├── contexts/           # State management
│   │   ├── lib/                # Utilities
│   │   ├── Dockerfile
│   │   └── package.json
│   │
│   ├── docker-compose.yml      # Docker orchestration
│   ├── START.sh                # Start script
│   └── STOP.sh                 # Stop script
│
└── docs/                       # Documentation
    ├── DEPLOYMENT/             # Deployment guides
    └── UI_GUIDES/              # UI development guides
```

---

## 🧪 Development

### Code Quality

**TypeScript/JavaScript:**
- ESLint for code linting
- TypeScript strict mode enabled
- Prettier for code formatting

**Go:**
- `gofmt` for code formatting
- `go vet` for static analysis
- GORM for type-safe database operations

### Testing

```bash
# Desktop application
cd arlog-desktop
npm run typecheck
npm run lint

# Web backend
cd arlog-web/backend
go test ./...
go vet ./...

# Web frontend
cd arlog-web/frontend
npm run lint
npm run typecheck
```

### Building for Production

**Desktop:**
```bash
cd arlog-desktop
npm run build:mac    # macOS
npm run build:win    # Windows
npm run build:linux  # Linux
```

**Web:**
```bash
cd arlog-web
docker-compose build
docker-compose up -d
```

---

## 🐛 Troubleshooting

### Desktop Application

**Issue: No contexts found**
```bash
# Verify kubeconfig exists
kubectl config get-contexts

# Check cluster access
kubectl cluster-info
```

**Issue: Logs not streaming**
```bash
# Test with kubectl
kubectl logs <pod-name> -n <namespace>

# Check pod status
kubectl get pods -n <namespace>
```

### Web Application

**Issue: Backend won't start**
```bash
# Check database connection
docker logs arlog-postgres

# Verify environment variables
docker exec arlog-backend env | grep DB_
```

**Issue: Frontend can't connect to backend**
```bash
# Check backend health
curl http://localhost:8080/health

# Verify network connectivity
docker network inspect arlog-network
```

---

## 📚 Documentation

### Application Documentation
- [Desktop README](./arlog-desktop/README.md) - Desktop application details
- [Web README](./arlog-web/README.md) - Web application details
- [Backend README](./arlog-web/backend/README.md) - Go backend documentation

### Guides
- [Quick Start Guide](./QUICKSTART.md) - Get started quickly
- [Deployment Guide](./docs/DEPLOYMENT/DOCKER_SETUP.md) - Production deployment
- [UI Development Guide](./docs/UI_GUIDES/) - UI component development

---

## 🤝 Contributing

Contributions are welcome! Please follow these guidelines:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines
- Follow existing code style and conventions
- Write meaningful commit messages
- Add tests for new features
- Update documentation as needed

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](./arlog-desktop/LICENSE) file for details.

---

## 👤 Author

**Batuhan Bozoguz**  
Senior DevOps Engineer & Backend Solutions Developer

- **Expertise**: Kubernetes, Cloud Infrastructure, CI/CD, Python, Go, TypeScript
- **Email**: batu@bozoguz.cloud

### Skills Demonstrated

This project showcases expertise in:

- **Kubernetes**: Multi-cluster management, service account integration, real-time log streaming
- **Infrastructure**: Docker containerization, PostgreSQL, microservices architecture
- **Backend Development**: Go (Gorilla Mux, GORM), RESTful APIs, WebSocket implementation
- **Frontend Development**: React, Next.js, TypeScript, modern UI/UX design
- **DevOps**: CI/CD pipelines, automated builds, cross-platform distribution
- **Security**: JWT authentication, Okta SSO integration, secure credential management

---

## 🗺️ Roadmap

### v1.2 (Planned)
- [ ] Log search and filtering
- [ ] Syntax highlighting for logs
- [ ] Save favorite namespaces/pods
- [ ] Pod events viewer
- [ ] Metrics integration (Prometheus)

### v2.0 (Future)
- [ ] Multi-tenant support
- [ ] Advanced permission management
- [ ] Audit logging
- [ ] Alerting and notifications
- [ ] Log aggregation and analysis

---

## 🙏 Acknowledgments

- [Kubernetes Client Libraries](https://kubernetes.io/docs/reference/using-api/client-libraries/)
- [shadcn/ui](https://ui.shadcn.com/) for beautiful UI components
- [Electron](https://www.electronjs.org/) for cross-platform desktop development
- [Next.js](https://nextjs.org/) for the web framework

---

<div align="center">

**Built with ❤️ for Kubernetes practitioners**

[⭐ Star this repo](https://github.com/batuhan-ozoguz/arlog) • [📖 Documentation](./docs/) • [🐛 Report Bug](https://github.com/batuhan-ozoguz/arlog/issues) • [💡 Request Feature](https://github.com/batuhan-ozoguz/arlog/issues)

</div>
