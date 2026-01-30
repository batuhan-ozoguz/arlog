import { useState, useEffect, useRef } from 'react'
import { useKubernetes } from '../contexts/KubernetesContext'
import { ChevronDown, Check, RefreshCw } from 'lucide-react'

export default function ContextSwitcher() {
  const { contexts, currentContext, switchContext } = useKubernetes()
  const [isOpen, setIsOpen] = useState(false)
  const [switching, setSwitching] = useState(false)
  const dropdownRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsOpen(false)
      }
    }

    if (isOpen) {
      document.addEventListener('mousedown', handleClickOutside)
      return () => document.removeEventListener('mousedown', handleClickOutside)
    }
  }, [isOpen])

  const handleContextSwitch = async (contextName: string) => {
    try {
      setSwitching(true)
      await switchContext(contextName)
      setIsOpen(false)
    } catch (error) {
      console.error('Failed to switch context:', error)
    } finally {
      setSwitching(false)
    }
  }

  const currentCtx = contexts.find(ctx => ctx.name === currentContext)

  return (
    <div className="relative" ref={dropdownRef}>
      <button
        onClick={() => setIsOpen(!isOpen)}
        disabled={switching}
        className="flex items-center gap-2 px-3 py-2 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 hover:border-primary transition-all disabled:opacity-50"
      >
        <div className="flex items-center gap-2">
          <div className="w-2 h-2 bg-green-500 rounded-full animate-pulse"></div>
          <span className="font-medium text-sm text-gray-900">
            {currentCtx?.name || 'Select Context'}
          </span>
        </div>
        {switching ? (
          <RefreshCw size={14} className="animate-spin text-gray-400" />
        ) : (
          <ChevronDown 
            size={16} 
            className={`text-gray-400 transition-transform ${isOpen ? 'rotate-180' : ''}`}
          />
        )}
      </button>

      {isOpen && (
        <div className="absolute top-full left-0 mt-2 w-96 bg-white border border-gray-200 rounded-lg shadow-xl py-2 z-50 animate-in fade-in slide-in-from-top-2 duration-200">
          <div className="px-4 py-2 border-b border-gray-100">
            <p className="text-xs font-semibold text-gray-500 uppercase tracking-wide">
              Kubernetes Contexts ({contexts.length})
            </p>
          </div>

          <div className="max-h-96 overflow-y-auto">
            {contexts.map((ctx) => {
              const isActive = ctx.name === currentContext

              return (
                <button
                  key={ctx.name}
                  onClick={() => handleContextSwitch(ctx.name)}
                  disabled={switching}
                  className={`w-full px-4 py-3 text-left hover:bg-gray-50 transition-colors flex items-center justify-between group ${
                    isActive ? 'bg-blue-50 border-l-2 border-primary' : ''
                  }`}
                >
                  <div className="flex-1 min-w-0">
                    <div className="font-semibold text-sm text-gray-900 truncate flex items-center gap-2">
                      {ctx.name}
                      {isActive && (
                        <span className="px-2 py-0.5 text-xs bg-primary text-white rounded">Current</span>
                      )}
                    </div>
                    <div className="text-xs text-gray-500 truncate mt-1">
                      <span className="font-medium">Cluster:</span> {ctx.cluster}
                    </div>
                    <div className="text-xs text-gray-400 truncate">
                      <span className="font-medium">Default NS:</span> {ctx.namespace}
                    </div>
                  </div>
                  {isActive && (
                    <Check size={18} className="text-primary flex-shrink-0 ml-3" strokeWidth={3} />
                  )}
                </button>
              )
            })}
          </div>

          {contexts.length === 0 && (
            <div className="px-4 py-8 text-center text-gray-500 text-sm">
              <p className="font-medium">No Kubernetes contexts found</p>
              <p className="text-xs mt-1">Check your ~/.kube/config file</p>
            </div>
          )}
        </div>
      )}
    </div>
  )
}



