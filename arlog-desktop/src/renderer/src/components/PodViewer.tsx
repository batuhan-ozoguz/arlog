import { useState, useEffect } from 'react'
import { useKubernetes } from '../contexts/KubernetesContext'
import { ArrowLeft, RefreshCw, Server, Eye } from 'lucide-react'

interface Pod {
  name: string
  namespace: string
  status: string
  ready: string
  restarts: number
  age: string
  containers: string[]
}

interface PodViewerProps {
  namespace: string
  onBack: () => void
  onPodSelect: (pod: { name: string; namespace: string; containers: string[] }) => void
}

export default function PodViewer({ namespace, onBack, onPodSelect }: PodViewerProps) {
  const { listPods } = useKubernetes()
  const [pods, setPods] = useState<Pod[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    fetchPods()
    // Auto-refresh every 10 seconds
    const interval = setInterval(fetchPods, 10000)
    return () => clearInterval(interval)
  }, [namespace])

  const fetchPods = async () => {
    try {
      setError(null)
      const podList = await listPods(namespace)
      setPods(podList)
    } catch (err: any) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  const getStatusBadge = (status: string) => {
    const statusLower = status.toLowerCase()
    if (statusLower === 'running') return 'bg-green-100 text-green-700 border-green-200'
    if (statusLower === 'pending') return 'bg-yellow-100 text-yellow-700 border-yellow-200'
    if (statusLower.includes('error') || statusLower.includes('failed')) return 'bg-red-100 text-red-700 border-red-200'
    return 'bg-gray-100 text-gray-700 border-gray-200'
  }

  return (
    <div className="max-w-7xl mx-auto px-6 py-8">
      <div className="mb-6">
        <button
          onClick={onBack}
          className="flex items-center gap-2 text-primary hover:text-primary-hover mb-4 transition-colors"
        >
          <ArrowLeft size={20} />
          <span className="font-medium">Back to Namespaces</span>
        </button>

        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-gray-900 mb-2 font-mono">{namespace}</h1>
            <p className="text-gray-600">Pods in this namespace</p>
          </div>
          <button
            onClick={() => fetchPods()}
            disabled={loading}
            className="flex items-center gap-2 px-4 py-2 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors disabled:opacity-50"
          >
            <RefreshCw size={16} className={loading ? 'animate-spin' : ''} />
            <span>Refresh</span>
          </button>
        </div>
      </div>

      {error && (
        <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg text-red-700">
          {error}
        </div>
      )}

      {loading && pods.length === 0 ? (
        <div className="text-center py-16">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto mb-4"></div>
          <p className="text-gray-600">Loading pods...</p>
        </div>
      ) : pods.length === 0 ? (
        <div className="text-center py-16 bg-white rounded-lg border border-gray-200">
          <Server size={48} className="mx-auto text-gray-300 mb-4" />
          <h3 className="text-lg font-semibold text-gray-900 mb-2">No Pods Found</h3>
          <p className="text-gray-500">There are no pods in this namespace</p>
        </div>
      ) : (
        <div className="bg-white rounded-lg border border-gray-200 overflow-hidden">
          <table className="w-full">
            <thead className="bg-gray-50 border-b border-gray-200">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
                  Pod Name
                </th>
                <th className="px-6 py-3 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
                  Status
                </th>
                <th className="px-6 py-3 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
                  Ready
                </th>
                <th className="px-6 py-3 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
                  Restarts
                </th>
                <th className="px-6 py-3 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
                  Age
                </th>
                <th className="px-6 py-3 text-right text-xs font-semibold text-gray-600 uppercase tracking-wider">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200">
              {pods.map((pod) => (
                <tr key={pod.name} className="hover:bg-gray-50 transition-colors">
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="text-sm font-medium text-gray-900 font-mono">{pod.name}</div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <span className={`px-3 py-1 inline-flex text-xs leading-5 font-semibold rounded-full border ${getStatusBadge(pod.status)}`}>
                      {pod.status}
                    </span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-700">
                    {pod.ready}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-700">
                    {pod.restarts}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-700">
                    {pod.age}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-right text-sm">
                    <button
                      onClick={() => onPodSelect({
                        name: pod.name,
                        namespace: pod.namespace,
                        containers: pod.containers
                      })}
                      className="inline-flex items-center gap-2 px-4 py-2 bg-primary text-white rounded-lg hover:bg-primary-hover transition-colors"
                    >
                      <Eye size={16} />
                      View Logs
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}



