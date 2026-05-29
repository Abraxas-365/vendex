import { useState } from 'react'
import { useParams, useNavigate } from '@tanstack/react-router'
import { ArrowLeft, Puzzle, AlertCircle, Loader2 } from 'lucide-react'
import { useQuery } from '@tanstack/react-query'
import * as api from '../../lib/api'

// ---------------------------------------------------------------------------
// Route: /admin/plugins/:name
// ---------------------------------------------------------------------------

export default function PluginView() {
  const { name } = useParams({ from: '/_admin/admin/plugins/$name' })
  const navigate = useNavigate()

  const {
    data: manifest,
    isLoading,
    error,
  } = useQuery({
    queryKey: ['plugins', name, 'manifest'],
    queryFn: () => api.getPluginManifest(name),
    enabled: Boolean(name),
  })

  const tabs = manifest?.ui?.tabs ?? []
  const [activeTab, setActiveTab] = useState(0)

  function handleBack() {
    void navigate({ to: '/admin/marketplace' })
  }

  // ── Loading state ──
  if (isLoading) {
    return (
      <div className="flex h-full items-center justify-center">
        <div className="flex flex-col items-center gap-3 text-slate-500">
          <Loader2 size={32} className="animate-spin" />
          <p className="text-sm">Loading plugin…</p>
        </div>
      </div>
    )
  }

  // ── Error state ──
  if (error || !manifest) {
    return (
      <div className="space-y-4">
        <button
          onClick={handleBack}
          className="inline-flex items-center gap-2 text-sm text-slate-500 hover:text-slate-700 transition-colors"
        >
          <ArrowLeft size={16} />
          Back to Marketplace
        </button>
        <div className="flex flex-col items-center gap-3 rounded-xl border border-red-200 bg-red-50 p-12 text-center">
          <AlertCircle size={32} className="text-red-400" />
          <p className="text-sm font-medium text-red-700">Failed to load plugin manifest</p>
          <p className="text-xs text-red-500">
            {error instanceof Error ? error.message : 'Unknown error'}
          </p>
        </div>
      </div>
    )
  }

  // ── No UI tabs ──
  if (tabs.length === 0) {
    return (
      <div className="space-y-4">
        <button
          onClick={handleBack}
          className="inline-flex items-center gap-2 text-sm text-slate-500 hover:text-slate-700 transition-colors"
        >
          <ArrowLeft size={16} />
          Back to Marketplace
        </button>
        <div className="rounded-xl border border-slate-200 bg-white">
          {/* Header */}
          <div className="flex items-center gap-4 border-b border-slate-200 px-6 py-4">
            <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-indigo-100 text-indigo-600">
              <Puzzle size={20} />
            </div>
            <div>
              <h2 className="text-base font-semibold text-slate-800">{manifest.display_name}</h2>
              <p className="text-xs text-slate-500">v{manifest.version} · by {manifest.author}</p>
            </div>
          </div>
          <div className="flex flex-col items-center gap-2 p-12 text-center">
            <Puzzle size={28} className="text-slate-300" />
            <p className="text-sm text-slate-500">This plugin has no UI tabs.</p>
          </div>
        </div>
      </div>
    )
  }

  const activeEntry = tabs[activeTab]?.entry ?? ''

  return (
    <div className="flex h-full flex-col space-y-0 -m-6">
      {/* Top bar */}
      <div className="flex h-14 shrink-0 items-center gap-4 border-b border-slate-200 bg-white px-6">
        <button
          onClick={handleBack}
          className="inline-flex items-center gap-2 text-sm text-slate-500 hover:text-slate-700 transition-colors"
        >
          <ArrowLeft size={16} />
          Marketplace
        </button>

        <div className="h-4 w-px bg-slate-200" />

        <div className="flex items-center gap-2">
          <div className="flex h-7 w-7 items-center justify-center rounded-lg bg-indigo-100 text-indigo-600">
            <Puzzle size={14} />
          </div>
          <span className="text-sm font-semibold text-slate-800">{manifest.display_name}</span>
          <span className="rounded-full bg-slate-100 px-2 py-0.5 text-xs text-slate-500">
            v{manifest.version}
          </span>
        </div>

        {/* Tabs */}
        {tabs.length > 1 && (
          <>
            <div className="h-4 w-px bg-slate-200" />
            <nav className="flex gap-1">
              {tabs.map((tab, idx) => (
                <button
                  key={idx}
                  onClick={() => setActiveTab(idx)}
                  className={`rounded-md px-3 py-1.5 text-xs font-medium transition-colors ${
                    activeTab === idx
                      ? 'bg-indigo-50 text-indigo-700'
                      : 'text-slate-500 hover:text-slate-700 hover:bg-slate-100'
                  }`}
                >
                  {tab.label}
                </button>
              ))}
            </nav>
          </>
        )}
      </div>

      {/* iframe container — fills remaining viewport height */}
      <div className="flex-1 overflow-hidden">
        {activeEntry ? (
          <iframe
            key={activeEntry}
            src={activeEntry}
            sandbox="allow-scripts allow-same-origin"
            className="h-full w-full border-0"
            title={tabs[activeTab]?.label ?? manifest.display_name}
          />
        ) : (
          <div className="flex h-full items-center justify-center">
            <p className="text-sm text-slate-400">No entry point configured for this tab.</p>
          </div>
        )}
      </div>
    </div>
  )
}
