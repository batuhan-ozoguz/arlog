import { createContext, useContext, useState, useEffect, ReactNode } from 'react'

interface Context {
  name: string
  cluster: string
  namespace: string
  user: string
}

interface Namespace {
  name: string
  status: string
  creationTimestamp?: Date
  labels: Record<string, string>
}

interface Pod {
  name: string
  namespace: string
  status: string
  ready: string
  restarts: number
  age: string
  containers: string[]
  labels: Record<string, string>
}

interface KubernetesContextType {
  contexts: Context[]
  currentContext: string | null
  namespaces: Namespace[]
  switchContext: (contextName: string) => Promise<void>
  refreshNamespaces: () => Promise<void>
  listPods: (namespace: string) => Promise<Pod[]>
  startLogStream: (namespace: string, podName: string, container?: string) => Promise<string>
  stopLogStream: (streamId: string) => Promise<void>
  loading: boolean
  error: string | null
}

const KubernetesContext = createContext<KubernetesContextType | null>(null)

export function KubernetesProvider({ children }: { children: ReactNode }) {
  const [contexts, setContexts] = useState<Context[]>([])
  const [currentContext, setCurrentContext] = useState<string | null>(null)
  const [namespaces, setNamespaces] = useState<Namespace[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    initializeK8s()
  }, [])

  const initializeK8s = async () => {
    try {
      setLoading(true)
      
      // Initialize Kubernetes client
      const initResult = await window.k8s.init()
      if (!initResult.success) {
        throw new Error(initResult.error || 'Failed to initialize Kubernetes')
      }

      // Get available contexts
      const ctxList = await window.k8s.getContexts()
      setContexts(ctxList)

      // Get current context
      const current = await window.k8s.getCurrentContext()
      setCurrentContext(current)

      // Load namespaces for current context
      await refreshNamespaces()

      setError(null)
    } catch (err: any) {
      console.error('Failed to initialize Kubernetes:', err)
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  const switchContext = async (contextName: string) => {
    try {
      setLoading(true)
      const result = await window.k8s.switchContext(contextName)
      if (result.success) {
        setCurrentContext(contextName)
        // Reload namespaces for new context
        await refreshNamespaces()
      } else {
        throw new Error(result.error)
      }
    } catch (err: any) {
      setError(err.message)
      throw err
    } finally {
      setLoading(false)
    }
  }

  const refreshNamespaces = async () => {
    try {
      const nsList = await window.k8s.listNamespaces()
      setNamespaces(nsList)
      setError(null)
    } catch (err: any) {
      console.error('Failed to load namespaces:', err)
      setError(err.message)
      setNamespaces([])
    }
  }

  const listPods = async (namespace: string): Promise<Pod[]> => {
    try {
      return await window.k8s.listPods(namespace)
    } catch (err: any) {
      console.error('Failed to list pods:', err)
      throw err
    }
  }

  const startLogStream = async (namespace: string, podName: string, container?: string): Promise<string> => {
    try {
      const result = await window.k8s.startLogStream({ namespace, podName, container })
      if (!result.success) {
        throw new Error(result.error)
      }
      return result.streamId
    } catch (err: any) {
      console.error('Failed to start log stream:', err)
      throw err
    }
  }

  const stopLogStream = async (streamId: string) => {
    try {
      await window.k8s.stopLogStream(streamId)
    } catch (err: any) {
      console.error('Failed to stop log stream:', err)
    }
  }

  const value = {
    contexts,
    currentContext,
    namespaces,
    switchContext,
    refreshNamespaces,
    listPods,
    startLogStream,
    stopLogStream,
    loading,
    error
  }

  return (
    <KubernetesContext.Provider value={value}>
      {children}
    </KubernetesContext.Provider>
  )
}

export function useKubernetes() {
  const context = useContext(KubernetesContext)
  if (!context) {
    throw new Error('useKubernetes must be used within KubernetesProvider')
  }
  return context
}



