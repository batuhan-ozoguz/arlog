import { useState } from 'react'
import { useKubernetes } from '../contexts/KubernetesContext'
import Navbar from '../components/Navbar'
import NamespaceGrid from '../components/NamespaceGrid'
import PodViewer from '../components/PodViewer'
import LogViewer from '../components/LogViewer'

function Dashboard() {
  const { namespaces, currentContext, loading, error } = useKubernetes()
  const [selectedNamespace, setSelectedNamespace] = useState<string | null>(null)
  const [selectedPod, setSelectedPod] = useState<{ name: string; namespace: string; containers: string[] } | null>(null)

  const handleNamespaceSelect = (namespace: string) => {
    setSelectedNamespace(namespace)
    setSelectedPod(null)
  }

  const handleBackToNamespaces = () => {
    setSelectedNamespace(null)
    setSelectedPod(null)
  }

  const handlePodSelect = (pod: { name: string; namespace: string; containers: string[] }) => {
    setSelectedPod(pod)
  }

  const handleBackToPods = () => {
    setSelectedPod(null)
  }

  if (loading) {
    return (
      <div className="h-screen flex items-center justify-center bg-gray-50">
        <div className="text-center">
          <div className="animate-spin rounded-full h-16 w-16 border-b-4 border-primary mx-auto mb-4"></div>
          <p className="text-lg font-semibold">Loading Kubernetes resources...</p>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="h-screen flex flex-col">
        <Navbar />
        <div className="flex-1 flex items-center justify-center">
          <div className="text-center max-w-md px-8">
            <div className="text-6xl mb-4">⚠️</div>
            <h2 className="text-2xl font-bold text-gray-900 mb-2">Error</h2>
            <p className="text-gray-600">{error}</p>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="h-screen flex flex-col overflow-hidden">
      <Navbar />
      
      <main className="flex-1 overflow-auto">
        {!selectedNamespace && !selectedPod && (
          <NamespaceGrid
            namespaces={namespaces}
            currentContext={currentContext}
            onNamespaceSelect={handleNamespaceSelect}
          />
        )}

        {selectedNamespace && !selectedPod && (
          <PodViewer
            namespace={selectedNamespace}
            onBack={handleBackToNamespaces}
            onPodSelect={handlePodSelect}
          />
        )}

        {selectedPod && (
          <LogViewer
            namespace={selectedPod.namespace}
            podName={selectedPod.name}
            containers={selectedPod.containers}
            onBack={handleBackToPods}
          />
        )}
      </main>
    </div>
  )
}

export default Dashboard



