import { ipcMain, IpcMainInvokeEvent } from 'electron'
import KubernetesClient from './kubernetes'

let k8sClient: KubernetesClient | null = null

export function setupIpcHandlers() {
  console.log('🔧 Setting up IPC handlers...')

  // Initialize Kubernetes client
  ipcMain.handle('k8s:init', async () => {
    try {
      k8sClient = new KubernetesClient()
      return { success: true }
    } catch (error: any) {
      console.error('Failed to initialize K8s client:', error)
      return { success: false, error: error.message }
    }
  })

  // Get contexts (clusters)
  ipcMain.handle('k8s:getContexts', async () => {
    try {
      if (!k8sClient) throw new Error('Kubernetes client not initialized')
      return k8sClient.getContexts()
    } catch (error: any) {
      throw new Error(error.message)
    }
  })

  // Get current context
  ipcMain.handle('k8s:getCurrentContext', async () => {
    try {
      if (!k8sClient) throw new Error('Kubernetes client not initialized')
      return k8sClient.getCurrentContext()
    } catch (error: any) {
      throw new Error(error.message)
    }
  })

  // Switch context
  ipcMain.handle('k8s:switchContext', async (_, contextName: string) => {
    try {
      if (!k8sClient) throw new Error('Kubernetes client not initialized')
      k8sClient.setCurrentContext(contextName)
      return { success: true }
    } catch (error: any) {
      return { success: false, error: error.message }
    }
  })

  // List namespaces
  ipcMain.handle('k8s:listNamespaces', async () => {
    try {
      if (!k8sClient) throw new Error('Kubernetes client not initialized')
      return await k8sClient.listNamespaces()
    } catch (error: any) {
      throw new Error(error.message)
    }
  })

  // List pods
  ipcMain.handle('k8s:listPods', async (_, namespace: string) => {
    try {
      if (!k8sClient) throw new Error('Kubernetes client not initialized')
      return await k8sClient.listPods(namespace)
    } catch (error: any) {
      throw new Error(error.message)
    }
  })

  // Start log stream
  ipcMain.handle('k8s:startLogStream', async (event: IpcMainInvokeEvent, { namespace, podName, container }) => {
    try {
      if (!k8sClient) throw new Error('Kubernetes client not initialized')
      
      console.log('📡 IPC: Starting log stream request:', { namespace, podName, container })
      
      const streamId = await k8sClient.streamLogs(
        namespace,
        podName,
        container,
        (data) => {
          // Send log data to renderer (don't log to avoid EPIPE)
          event.sender.send('k8s:logData', { streamId: `${namespace}/${podName}`, data })
        },
        (error) => {
          console.error('❌ Log error in IPC:', error.message)
          event.sender.send('k8s:logError', { streamId: `${namespace}/${podName}`, error: error.message })
        },
        () => {
          console.log('🏁 Log stream ended in IPC')
          event.sender.send('k8s:logEnd', { streamId: `${namespace}/${podName}` })
        }
      )

      console.log('✅ IPC: Log stream started successfully, streamId:', streamId)
      return { success: true, streamId }
    } catch (error: any) {
      console.error('❌ IPC: Failed to start log stream:', error)
      console.error('Error details:', error.message, error.stack)
      return { success: false, error: error.message }
    }
  })

  // Stop log stream
  ipcMain.handle('k8s:stopLogStream', async (_, streamId: string) => {
    try {
      if (!k8sClient) throw new Error('Kubernetes client not initialized')
      k8sClient.stopLogStream(streamId)
      return { success: true }
    } catch (error: any) {
      return { success: false, error: error.message }
    }
  })

  console.log('✅ IPC handlers registered')
}

