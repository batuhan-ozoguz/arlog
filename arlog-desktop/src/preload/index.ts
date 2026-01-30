import { contextBridge, ipcRenderer } from 'electron'
import { electronAPI } from '@electron-toolkit/preload'

// Expose Kubernetes API to renderer process
const k8sApi = {
  // Initialize client
  init: () => ipcRenderer.invoke('k8s:init'),
  
  // Context management
  getContexts: () => ipcRenderer.invoke('k8s:getContexts'),
  getCurrentContext: () => ipcRenderer.invoke('k8s:getCurrentContext'),
  switchContext: (contextName: string) => ipcRenderer.invoke('k8s:switchContext', contextName),
  
  // Resource listing
  listNamespaces: () => ipcRenderer.invoke('k8s:listNamespaces'),
  listPods: (namespace: string) => ipcRenderer.invoke('k8s:listPods', namespace),
  
  // Log streaming
  startLogStream: (params: { namespace: string; podName: string; container?: string }) => 
    ipcRenderer.invoke('k8s:startLogStream', params),
  stopLogStream: (streamId: string) => 
    ipcRenderer.invoke('k8s:stopLogStream', streamId),
  
  // Log event listeners
  onLogData: (callback: (data: { streamId: string; data: string }) => void) => {
    const listener = (_: any, data: any) => callback(data)
    ipcRenderer.on('k8s:logData', listener)
    return () => ipcRenderer.removeListener('k8s:logData', listener)
  },
  
  onLogError: (callback: (error: { streamId: string; error: string }) => void) => {
    const listener = (_: any, error: any) => callback(error)
    ipcRenderer.on('k8s:logError', listener)
    return () => ipcRenderer.removeListener('k8s:logError', listener)
  },
  
  onLogEnd: (callback: (data: { streamId: string }) => void) => {
    const listener = (_: any, data: any) => callback(data)
    ipcRenderer.on('k8s:logEnd', listener)
    return () => ipcRenderer.removeListener('k8s:logEnd', listener)
  },
  
  removeAllLogListeners: () => {
    ipcRenderer.removeAllListeners('k8s:logData')
    ipcRenderer.removeAllListeners('k8s:logError')
    ipcRenderer.removeAllListeners('k8s:logEnd')
  }
}

// Use `contextBridge` APIs to expose Electron APIs to renderer only if context isolation is enabled
if (process.contextIsolated) {
  try {
    contextBridge.exposeInMainWorld('electron', electronAPI)
    contextBridge.exposeInMainWorld('k8s', k8sApi)
  } catch (error) {
    console.error('Error exposing APIs:', error)
  }
} else {
  // @ts-ignore (define in dts)
  window.electron = electronAPI
  // @ts-ignore (define in dts)
  window.k8s = k8sApi
}



