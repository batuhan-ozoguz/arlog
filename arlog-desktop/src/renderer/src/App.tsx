import { useState, useEffect } from 'react'
import { KubernetesProvider } from './contexts/KubernetesContext'
import Dashboard from './pages/Dashboard'

function App() {
  const [initialized, setInitialized] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    checkKubernetesAvailability()
  }, [])

  const checkKubernetesAvailability = async () => {
    try {
      const result = await window.k8s.init()
      if (result.success) {
        setInitialized(true)
      } else {
        setError(result.error || 'Failed to initialize Kubernetes')
      }
    } catch (err: any) {
      setError(err.message || 'Failed to initialize Kubernetes')
    }
  }

  if (error) {
    return (
      <div className="flex items-center justify-center h-screen bg-gray-50">
        <div className="text-center max-w-md px-8">
          <div className="text-6xl mb-4">⚠️</div>
          <h2 className="text-2xl font-bold text-gray-900 mb-2">Kubernetes Not Available</h2>
          <p className="text-gray-600 mb-4">{error}</p>
          <div className="text-left bg-gray-100 p-4 rounded-lg text-sm">
            <p className="font-semibold mb-2">Please check:</p>
            <ul className="list-disc list-inside space-y-1 text-gray-700">
              <li>~/.kube/config file exists</li>
              <li>kubectl is configured</li>
              <li>At least one cluster context is available</li>
            </ul>
          </div>
          <button
            onClick={() => window.location.reload()}
            className="mt-6 px-6 py-2 bg-primary text-white rounded-lg hover:bg-primary-hover transition-colors"
          >
            Retry
          </button>
        </div>
      </div>
    )
  }

  if (!initialized) {
    return (
      <div className="flex items-center justify-center h-screen bg-gray-50">
        <div className="text-center">
          <div className="animate-spin rounded-full h-16 w-16 border-b-4 border-primary mx-auto mb-4"></div>
          <p className="text-lg font-semibold text-gray-700">Initializing Kubernetes</p>
          <p className="text-sm text-gray-500 mt-2">Reading ~/.kube/config</p>
        </div>
      </div>
    )
  }

  return (
    <KubernetesProvider>
      <Dashboard />
    </KubernetesProvider>
  )
}

export default App



