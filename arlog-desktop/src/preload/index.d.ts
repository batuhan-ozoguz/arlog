import { ElectronAPI } from '@electron-toolkit/preload'

declare global {
  interface Window {
    electron: ElectronAPI
    k8s: {
      init: () => Promise<{ success: boolean; error?: string }>
      getContexts: () => Promise<Array<{
        name: string
        cluster: string
        namespace: string
        user: string
      }>>
      getCurrentContext: () => Promise<string>
      switchContext: (contextName: string) => Promise<{ success: boolean; error?: string }>
      listNamespaces: () => Promise<Array<{
        name: string
        status: string
        creationTimestamp?: Date
        labels: Record<string, string>
      }>>
      listPods: (namespace: string) => Promise<Array<{
        name: string
        namespace: string
        status: string
        ready: string
        restarts: number
        age: string
        containers: string[]
        labels: Record<string, string>
      }>>
      startLogStream: (params: {
        namespace: string
        podName: string
        container?: string
      }) => Promise<{ success: boolean; streamId: string; error?: string }>
      stopLogStream: (streamId: string) => Promise<{ success: boolean; error?: string }>
      onLogData: (callback: (data: { streamId: string; data: string }) => void) => () => void
      onLogError: (callback: (error: { streamId: string; error: string }) => void) => () => void
      onLogEnd: (callback: (data: { streamId: string }) => void) => () => void
      removeAllLogListeners: () => void
    }
  }
}



