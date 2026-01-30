import { Folder, ArrowRight } from 'lucide-react'

interface Namespace {
  name: string
  status: string
  creationTimestamp?: Date
  labels: Record<string, string>
}

interface NamespaceGridProps {
  namespaces: Namespace[]
  currentContext: string | null
  onNamespaceSelect: (namespace: string) => void
}

export default function NamespaceGrid({ namespaces, currentContext, onNamespaceSelect }: NamespaceGridProps) {
  const getStatusColor = (status: string) => {
    if (status === 'Active') return 'text-green-600 bg-green-50'
    return 'text-gray-600 bg-gray-50'
  }

  if (namespaces.length === 0) {
    return (
      <div className="max-w-7xl mx-auto px-6 py-16">
        <div className="text-center">
          <Folder size={64} className="mx-auto text-gray-300 mb-4" />
          <h3 className="text-lg font-semibold text-gray-900 mb-2">No Namespaces Found</h3>
          <p className="text-gray-500">No namespaces available in the current context.</p>
        </div>
      </div>
    )
  }

  return (
    <div className="max-w-7xl mx-auto px-6 py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">
          Kubernetes Namespaces
        </h1>
        <p className="text-gray-600">
          Monitor and manage your cluster namespaces
          {currentContext && <span className="ml-2">• <span className="font-medium">{currentContext}</span></span>}
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {namespaces.map((ns) => (
          <div
            key={ns.name}
            onClick={() => onNamespaceSelect(ns.name)}
            className="bg-white border border-gray-200 rounded-lg p-6 hover:shadow-lg hover:border-primary transition-all cursor-pointer group"
          >
            <div className="flex items-start justify-between mb-4">
              <div className="w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center">
                <Folder size={24} className="text-primary" />
              </div>
              <span className={`px-3 py-1 rounded-full text-xs font-medium ${getStatusColor(ns.status)}`}>
                {ns.status}
              </span>
            </div>

            <h3 className="text-lg font-semibold text-gray-900 mb-2 font-mono">
              {ns.name}
            </h3>

            <div className="flex items-center justify-between mt-4 pt-4 border-t border-gray-100">
              <span className="text-sm text-primary font-medium group-hover:gap-2 flex items-center gap-1 transition-all">
                View Pods
                <ArrowRight size={16} className="group-hover:translate-x-1 transition-transform" />
              </span>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}



