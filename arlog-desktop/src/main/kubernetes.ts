import * as k8s from '@kubernetes/client-node'
import * as path from 'path'
import * as os from 'os'
import { PassThrough } from 'stream'

export class KubernetesClient {
  private kc: k8s.KubeConfig
  private k8sApi: k8s.CoreV1Api
  private activeStreams: Map<string, any> = new Map()

  constructor() {
    this.kc = new k8s.KubeConfig()
    
    // Load from default kubeconfig location
    const kubeConfigPath = process.env.KUBECONFIG || 
                          path.join(os.homedir(), '.kube', 'config')
    
    try {
      this.kc.loadFromFile(kubeConfigPath)
      this.k8sApi = this.kc.makeApiClient(k8s.CoreV1Api)
      console.log('✅ Kubernetes config loaded from:', kubeConfigPath)
    } catch (error) {
      console.error('❌ Failed to load kubeconfig:', error)
      throw error
    }
  }

  // Get all available contexts (clusters)
  getContexts() {
    return this.kc.getContexts().map(ctx => ({
      name: ctx.name,
      cluster: ctx.cluster,
      namespace: ctx.namespace || 'default',
      user: ctx.user
    }))
  }

  // Get current context
  getCurrentContext() {
    return this.kc.getCurrentContext()
  }

  // Switch context
  setCurrentContext(contextName: string) {
    this.kc.setCurrentContext(contextName)
    this.k8sApi = this.kc.makeApiClient(k8s.CoreV1Api)
    console.log('✅ Switched to context:', contextName)
  }

  // List namespaces
  async listNamespaces() {
    try {
      const response = await this.k8sApi.listNamespace()
      return response.body.items.map(ns => ({
        name: ns.metadata?.name || '',
        status: ns.status?.phase || 'Active',
        creationTimestamp: ns.metadata?.creationTimestamp,
        labels: ns.metadata?.labels || {}
      }))
    } catch (error: any) {
      console.error('❌ Failed to list namespaces:', error.message)
      throw error
    }
  }

  // List pods in namespace
  async listPods(namespace: string) {
    try {
      const response = await this.k8sApi.listNamespacedPod(namespace)
      return response.body.items.map(pod => ({
        name: pod.metadata?.name || '',
        namespace: pod.metadata?.namespace || '',
        status: pod.status?.phase || 'Unknown',
        ready: this.getPodReadyStatus(pod),
        restarts: this.getPodRestartCount(pod),
        age: this.calculateAge(pod.metadata?.creationTimestamp),
        containers: pod.spec?.containers.map(c => c.name) || [],
        labels: pod.metadata?.labels || {}
      }))
    } catch (error: any) {
      console.error('❌ Failed to list pods:', error.message)
      throw error
    }
  }

  // Stream pod logs
  async streamLogs(
    namespace: string,
    podName: string,
    container: string | undefined,
    onData: (data: string) => void,
    onError: (error: Error) => void,
    onEnd: () => void
  ): Promise<string> {
    try {
      const log = new k8s.Log(this.kc)
      const streamId = `${namespace}/${podName}/${container || 'default'}`

      console.log('🔄 Starting log stream for:', streamId, 'container:', container)

      // Create a PassThrough stream to capture logs
      const logStream = new PassThrough()

      // Set up stream handlers
      logStream.on('data', (chunk: Buffer) => {
        const logLine = chunk.toString()
        // Don't log the actual content to avoid EPIPE errors
        onData(logLine)
      })

      logStream.on('error', (streamErr: Error) => {
        console.error('❌ Stream error:', streamErr)
        onError(streamErr)
      })

      logStream.on('end', () => {
        console.log('📡 Log stream ended:', streamId)
        this.activeStreams.delete(streamId)
        onEnd()
      })

      // Start the log stream using the Log API
      const req = await log.log(
        namespace,
        podName,
        container || undefined, // Use undefined instead of empty string
        logStream,
        {
          follow: true,
          tailLines: 100,
          pretty: false,
          timestamps: true
        }
      )

      // Store stream for cleanup
      this.activeStreams.set(streamId, req)

      console.log('✅ Log stream started:', streamId)
      return streamId
    } catch (error: any) {
      console.error('❌ Failed to start log stream:', error)
      console.error('Error details:', error.message, error.response?.body)
      throw error
    }
  }

  // Stop log stream
  stopLogStream(streamId: string) {
    const stream = this.activeStreams.get(streamId)
    if (stream && stream.destroy) {
      stream.destroy()
      this.activeStreams.delete(streamId)
      console.log('🛑 Stopped log stream:', streamId)
    }
  }

  // Stop all log streams
  stopAllLogStreams() {
    this.activeStreams.forEach((stream, streamId) => {
      if (stream && stream.destroy) {
        stream.destroy()
      }
      console.log('🛑 Stopped log stream:', streamId)
    })
    this.activeStreams.clear()
  }

  // Helper methods
  private getPodReadyStatus(pod: k8s.V1Pod): string {
    const total = pod.spec?.containers.length || 0
    const ready = pod.status?.containerStatuses?.filter(c => c.ready).length || 0
    return `${ready}/${total}`
  }

  private getPodRestartCount(pod: k8s.V1Pod): number {
    return pod.status?.containerStatuses?.reduce((sum, c) => sum + c.restartCount, 0) || 0
  }

  private calculateAge(timestamp: Date | undefined): string {
    if (!timestamp) return 'Unknown'
    const now = new Date()
    const diff = now.getTime() - new Date(timestamp).getTime()
    const days = Math.floor(diff / (1000 * 60 * 60 * 24))
    const hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60))
    const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60))
    
    if (days > 0) return `${days}d${hours}h`
    if (hours > 0) return `${hours}h${minutes}m`
    if (minutes > 0) return `${minutes}m`
    return `${Math.floor(diff / 1000)}s`
  }
}

export default KubernetesClient

