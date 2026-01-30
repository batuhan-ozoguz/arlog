# ArLOG Desktop - Kubernetes Log Viewer

Standalone desktop application for viewing Kubernetes pod logs using your local kubeconfig.

## 🚀 Quick Start for Users

### Download & Install (No Development Setup Needed)

**Windows Users**:
1. Download `ArLOG-Setup-1.1.0.exe` from [Releases](https://dev.azure.com/arcelikdevops/ArLog/_git/ArLog?path=/&version=GBv1.1-production-build&_a=releases)
2. Run the installer
3. Launch ArLOG from Start Menu
4. Done! ✅

**macOS Users**:
1. Download `ArLOG-1.1.0.dmg` from Releases
2. Open DMG and drag to Applications
3. Launch ArLOG
4. Done! ✅

---

## 🎯 Features

- ✅ **No Server Required** - Uses local ~/.kube/config
- ✅ **Multi-Cluster Support** - Switch between Kubernetes contexts
- ✅ **Real-time Log Streaming** - Live pod logs
- ✅ **Beautiful Minimal UI** - Clean, professional interface
- ✅ **Cross-Platform** - macOS, Windows, Linux
- ✅ **Offline Capable** - Works with local clusters
- ✅ **No Authentication** - Uses your kubectl config

---

## 🛠️ For Developers

### Prerequisites

- Node.js 18+ and npm/pnpm
- kubectl configured with at least one context
- Access to Kubernetes cluster

### Installation from Source

#### macOS / Linux:

```bash
# Clone repository
git clone https://dev.azure.com/arcelikdevops/ArLog/_git/ArLog
cd ArLog

# Install dependencies
npm install

# Run application
npm run dev
```

#### 🪟 Windows (RECOMMENDED METHOD):

Windows'ta npm install sorunları yaşanabilir. **pnpm kullanın** ve **OneDrive dışı bir klasörde** çalışın:

```powershell
# 1. OneDrive DIŞI klasöre git
cd C:\
mkdir Dev
cd Dev

# 2. Clone repository
git clone https://dev.azure.com/arcelikdevops/ArLog/_git/ArLog arlog-desktop
cd arlog-desktop

# 3. pnpm'i global yükle (bir kez)
npm install -g pnpm

# 4. Dependencies yükle
pnpm install

# 5. Uygulamayı çalıştır
pnpm run dev
```

**⚠️ Önemli**: OneDrive klasöründe (`D:\Users\...\OneDrive\`) ÇALIŞTIRMAYIN! node_modules sorun yaratır.

---

## 🏗️ Building Windows Installer (.exe)

### For Windows Developers:

```powershell
# Build Windows installer
cd C:\Dev\arlog-desktop
pnpm run build:win

# Output: dist\ArLOG Setup 1.1.0.exe
```

**Süre**: 10-15 dakika  
**Çıktı**: NSIS installer (.exe)

Upload the .exe to Azure DevOps Releases for distribution.

---

## 💡 How It Works

1. **Reads your kubeconfig** - Uses ~/.kube/config automatically
2. **Auto-detects contexts** - Shows all available clusters
3. **Switch contexts** - Seamlessly change between clusters
4. **View namespaces** - See all namespaces in current context
5. **List pods** - View pods with status, ready state, restarts
6. **Stream logs** - Real-time log viewing with pause/resume

---

## 🎨 Features

### Context Switching
- Dropdown showing all Kubernetes contexts
- One-click switching between clusters
- Auto-refresh namespaces on switch

### Namespace Viewer
- Grid view of all namespaces
- Status indicators
- Click to view pods

### Pod Viewer
- Table of pods with status, ready, restarts, age
- Auto-refresh every 10 seconds
- Click "View Logs" to stream logs

### Log Viewer
- Real-time log streaming
- Pause/Resume streaming
- Auto-scroll toggle
- Clear logs
- Download logs as .txt
- Container selector (for multi-container pods)
- Connection status indicator
- Dark terminal theme

---

## 🔧 Development Commands

```bash
# Run in development mode (hot reload)
npm run dev         # or: pnpm run dev

# Build for production
npm run build

# Build installers
npm run build:mac      # macOS .dmg
npm run build:win      # Windows .exe
npm run build:linux    # Linux .AppImage
```

---

## 🛠️ Tech Stack

- **Electron** 32 - Desktop app framework
- **Vite** - Build tool (via electron-vite)
- **React** 18 - UI framework
- **TypeScript** - Type safety
- **Tailwind CSS** - Styling
- **@kubernetes/client-node** - Kubernetes API client
- **shadcn/ui** - UI components

---

## 🐛 Troubleshooting

### "No contexts found"
**Problem**: Can't read kubeconfig
```bash
# Check if ~/.kube/config exists
kubectl config get-contexts

# Verify cluster access
kubectl cluster-info
```

### "Failed to connect to cluster"
**Problem**: Cluster not accessible
```bash
# Test connection
kubectl get pods --all-namespaces
```

### Logs not streaming
**Problem**: Pod or permissions issue
```bash
# Test with kubectl
kubectl logs <pod-name> -n <namespace>
```

### Windows: npm install error
**Solution**: Use pnpm and avoid OneDrive paths
```powershell
cd C:\Dev
git clone <repo> arlog-desktop
cd arlog-desktop
npm install -g pnpm
pnpm install
pnpm run dev
```

---

## 📄 License

MIT License - Copyright (c) 2025 Batuhan Bozoguz

See [LICENSE](./LICENSE) file for details.

---

## 👥 Author

**Batuhan Bozoguz**  
Senior DevOps Engineer  
Email: batuhan.ozoguz@beko.com

### Built With
- Expertise in Kubernetes, Electron, React, TypeScript
- DevOps best practices and cloud solutions
- Beautiful UI/UX design with v0.dev

---

## 🔮 Roadmap

**v1.1** (Current):
- ✅ Production-ready installers
- ✅ Windows optimization
- ✅ pnpm support

**v1.2** (Planned):
- [ ] Log search and filtering
- [ ] Syntax highlighting
- [ ] Save favorite namespaces
- [ ] Pod events viewer

**v2.0** (Future):
- [ ] Web application
- [ ] Multi-user support
- [ ] Centralized authentication

---

**Note**: This is V1 (Desktop). V2 will be a web application with multi-user support and centralized authentication.
