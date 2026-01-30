import { useState, useEffect, useRef } from 'react'
import { ArrowLeft, Play, Pause, Download, Trash2, Settings } from 'lucide-react'

interface LogViewerProps {
  namespace: string
  podName: string
  containers: string[]
  onBack: () => void
}

export default function LogViewer({ namespace, podName, containers, onBack }: LogViewerProps) {
  const [logs, setLogs] = useState<string[]>([])
  const [selectedContainer, setSelectedContainer] = useState(containers[0] || '')
  const [isPaused, setIsPaused] = useState(false)
  const [autoScroll, setAutoScroll] = useState(true)
  const [connected, setConnected] = useState(false)
  const logsContainerRef = useRef<HTMLDivElement>(null)
  const streamIdRef = useRef<string | null>(null)

  useEffect(() => {
    startLogStream()
    return () => {
      cleanup()
    }
  }, [namespace, podName, selectedContainer])

  useEffect(() => {
    // Auto-scroll to bottom when new logs arrive
    if (autoScroll && !isPaused && logsContainerRef.current) {
      logsContainerRef.current.scrollTop = logsContainerRef.current.scrollHeight
    }
  }, [logs, autoScroll, isPaused])

  const startLogStream = async () => {
    try {
      console.log('🔄 LogViewer: Setting up log stream...', { namespace, podName, container: selectedContainer })
      
      // Setup log listeners BEFORE starting the stream
      const removeDataListener = window.k8s.onLogData((data) => {
        // Received log data - add to state
        if (!isPaused) {
          setLogs(prev => [...prev, data.data])
        }
      })

      const removeErrorListener = window.k8s.onLogError((error) => {
        console.error('❌ LogViewer: Log stream error:', error)
        setConnected(false)
      })

      const removeEndListener = window.k8s.onLogEnd((data) => {
        console.log('🏁 LogViewer: Log stream ended:', data)
        setConnected(false)
      })

      // Start stream
      console.log('📡 LogViewer: Starting stream with params:', { namespace, podName, container: selectedContainer })
      const result = await window.k8s.startLogStream({
        namespace,
        podName,
        container: selectedContainer || undefined
      })

      console.log('✅ LogViewer: Stream start result:', result)

      if (result.success) {
        streamIdRef.current = result.streamId
        setConnected(true)
        console.log('✅ LogViewer: Connected, streamId:', result.streamId)
      } else {
        console.error('❌ LogViewer: Stream start failed:', result.error)
        setConnected(false)
      }

      // Cleanup function
      return () => {
        removeDataListener()
        removeErrorListener()
        removeEndListener()
      }
    } catch (error: any) {
      console.error('❌ LogViewer: Failed to start log stream:', error)
      setConnected(false)
    }
  }

  const cleanup = () => {
    if (streamIdRef.current) {
      window.k8s.stopLogStream(streamIdRef.current)
      streamIdRef.current = null
    }
    window.k8s.removeAllLogListeners()
    setConnected(false)
  }

  const handleTogglePause = () => {
    setIsPaused(!isPaused)
  }

  const handleClear = () => {
    setLogs([])
  }

  const handleDownload = () => {
    const logText = logs.join('')
    const blob = new Blob([logText], { type: 'text/plain' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `${podName}-${new Date().toISOString()}.txt`
    a.click()
    URL.revokeObjectURL(url)
  }

  const handleContainerChange = (container: string) => {
    cleanup()
    setLogs([])
    setSelectedContainer(container)
  }

  return (
    <div className="h-full flex flex-col">
      {/* Header */}
      <div className="bg-white border-b border-gray-200 px-6 py-4">
        <button
          onClick={onBack}
          className="flex items-center gap-2 text-primary hover:text-primary-hover mb-3 transition-colors"
        >
          <ArrowLeft size={20} />
          <span className="font-medium">Back to Pods</span>
        </button>

        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-2xl font-bold text-gray-900 font-mono mb-1">{podName}</h2>
            <p className="text-sm text-gray-500">Namespace: {namespace}</p>
          </div>

          <div className="flex items-center gap-4">
            {containers.length > 1 && (
              <select
                value={selectedContainer}
                onChange={(e) => handleContainerChange(e.target.value)}
                className="px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary"
              >
                {containers.map(container => (
                  <option key={container} value={container}>
                    Container: {container}
                  </option>
                ))}
              </select>
            )}

            <span className={`flex items-center gap-2 px-3 py-1 rounded-full text-xs font-medium ${
              connected ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'
            }`}>
              <span className={`w-2 h-2 rounded-full ${connected ? 'bg-green-500 animate-pulse' : 'bg-red-500'}`}></span>
              {connected ? 'Connected' : 'Disconnected'}
            </span>
          </div>
        </div>
      </div>

      {/* Controls */}
      <div className="bg-white border-b border-gray-200 px-6 py-3 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <button
            onClick={handleTogglePause}
            className={`flex items-center gap-2 px-4 py-2 rounded-lg border transition-colors ${
              isPaused
                ? 'bg-green-50 border-green-300 text-green-700 hover:bg-green-100'
                : 'bg-white border-gray-300 text-gray-700 hover:bg-gray-50'
            }`}
          >
            {isPaused ? (
              <>
                <Play size={16} />
                Resume
              </>
            ) : (
              <>
                <Pause size={16} />
                Pause
              </>
            )}
          </button>

          <button
            onClick={handleClear}
            className="flex items-center gap-2 px-4 py-2 bg-white border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors"
          >
            <Trash2 size={16} />
            Clear
          </button>

          <button
            onClick={handleDownload}
            disabled={logs.length === 0}
            className="flex items-center gap-2 px-4 py-2 bg-white border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors disabled:opacity-50"
          >
            <Download size={16} />
            Download
          </button>

          <label className="flex items-center gap-2 px-3 py-2 text-sm text-gray-700">
            <input
              type="checkbox"
              checked={autoScroll}
              onChange={(e) => setAutoScroll(e.target.checked)}
              className="rounded border-gray-300 text-primary focus:ring-primary"
            />
            Auto-scroll
          </label>
        </div>

        <div className="flex items-center gap-3 text-sm text-gray-500">
          <span>{logs.length} lines</span>
        </div>
      </div>

      {/* Logs Display */}
      <div 
        ref={logsContainerRef}
        className="flex-1 bg-[#0d1117] p-4 overflow-auto font-mono text-sm"
      >
        {logs.length === 0 ? (
          <div className="text-center py-16 text-gray-400">
            {connected ? 'Waiting for logs...' : 'Connecting to pod...'}
          </div>
        ) : (
          <pre className="text-gray-300">
            {logs.map((log, index) => (
              <div 
                key={index} 
                className="hover:bg-gray-800/30 transition-colors border-l-2 border-transparent hover:border-primary pl-2"
              >
                {log}
              </div>
            ))}
          </pre>
        )}
      </div>
    </div>
  )
}

