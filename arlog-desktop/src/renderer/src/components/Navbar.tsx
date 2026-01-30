import { useKubernetes } from '../contexts/KubernetesContext'
import ContextSwitcher from './ContextSwitcher'

export default function Navbar() {
  const { currentContext } = useKubernetes()

  return (
    <nav className="h-[60px] bg-white border-b border-gray-200 px-6 flex items-center justify-between sticky top-0 z-50">
      <div className="flex items-center gap-8">
        <div className="flex items-center gap-3">
          <h1 className="text-xl font-bold text-gray-900 tracking-tight">ArLOG</h1>
          <span className="px-2 py-1 text-xs font-medium bg-blue-100 text-blue-700 rounded border border-blue-200">
            DESKTOP
          </span>
        </div>

        {currentContext && <ContextSwitcher />}
      </div>

      <div className="flex items-center gap-4">
        <div className="text-sm text-gray-500">
          {currentContext && (
            <span className="font-medium">Context: {currentContext}</span>
          )}
        </div>
      </div>
    </nav>
  )
}



