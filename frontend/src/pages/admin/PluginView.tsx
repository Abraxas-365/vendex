import { useState } from 'react'
import { useParams, useNavigate } from '@tanstack/react-router'
import { ArrowLeft, Puzzle, AlertCircle, Loader2, Power, PowerOff, Settings, Layout, ExternalLink, Circle } from 'lucide-react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import * as api from '../../lib/api'

// ---------------------------------------------------------------------------
// Route: /admin/plugins/:name  (name is used as plugin ID)
// ---------------------------------------------------------------------------

type PluginViewTab = 'ui' | 'widget' | 'settings'

export default function PluginView() {
  const { name } = useParams({ from: '/_admin/admin/plugins/$name' })
  const navigate = useNavigate()
  const qc = useQueryClient()

  const [activeUiTab, setActiveUiTab] = useState(0)
  const [viewTab, setViewTab] = useState<PluginViewTab>('ui')
  const [settingsJson, setSettingsJson] = useState('')
  const [settingsError, setSettingsError] = useState('')
  const [savingSettings, setSavingSettings] = useState(false)

  const {
    data: manifestData,
    isLoading,
    error,
  } = useQuery({
    queryKey: ['plugins', name, 'manifest'],
    queryFn: () => api.getPluginManifest(name),
    enabled: Boolean(name),
  })

  const { data: pluginData } = useQuery({
    queryKey: ['plugins', name, 'detail'],
    queryFn: () => api.getMarketplacePlugin(name),
    enabled: Boolean(name),
  })

  const manifest = manifestData?.manifest
  const latestVersion = pluginData?.latest_version

  // Enable/disable mutations
  const enableMutation = useMutation({
    mutationFn: () => api.enablePlugin(name),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: ['marketplace'] })
      void qc.invalidateQueries({ queryKey: ['plugins'] })
    },
  })

  const disableMutation = useMutation({
    mutationFn: () => api.disablePlugin(name),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: ['marketplace'] })
      void qc.invalidateQueries({ queryKey: ['plugins'] })
    },
  })

  const tabs = manifest?.ui?.tabs ?? []

  function handleBack() {
    void navigate({ to: '/admin/marketplace' })
  }

  async function handleSaveSettings() {
    setSettingsError('')
    let parsed: Record<string, unknown>
    try {
      parsed = JSON.parse(settingsJson) as Record<string, unknown>
    } catch {
      setSettingsError('Invalid JSON — please check your syntax.')
      return
    }
    setSavingSettings(true)
    try {
      await api.updatePluginSettings(name, parsed)
      void qc.invalidateQueries({ queryKey: ['marketplace'] })
    } catch (err) {
      setSettingsError(err instanceof Error ? err.message : 'Failed to save settings.')
    } finally {
      setSavingSettings(false)
    }
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

  const activeEntry = tabs[activeUiTab]?.entry ?? ''

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
          <span className="text-xs text-slate-400">by {manifest.author}</span>
        </div>

        {/* View tabs: UI / Settings */}
        <div className="h-4 w-px bg-slate-200" />
        <nav className="flex gap-1">
          <button
            onClick={() => setViewTab('ui')}
            className={`rounded-md px-3 py-1.5 text-xs font-medium transition-colors ${
              viewTab === 'ui'
                ? 'bg-indigo-50 text-indigo-700'
                : 'text-slate-500 hover:text-slate-700 hover:bg-slate-100'
            }`}
          >
            UI
          </button>
          <button
            onClick={() => setViewTab('widget')}
            className={`inline-flex items-center gap-1.5 rounded-md px-3 py-1.5 text-xs font-medium transition-colors ${
              viewTab === 'widget'
                ? 'bg-indigo-50 text-indigo-700'
                : 'text-slate-500 hover:text-slate-700 hover:bg-slate-100'
            }`}
          >
            <Layout size={12} />
            Widget
          </button>
          <button
            onClick={() => setViewTab('settings')}
            className={`inline-flex items-center gap-1.5 rounded-md px-3 py-1.5 text-xs font-medium transition-colors ${
              viewTab === 'settings'
                ? 'bg-indigo-50 text-indigo-700'
                : 'text-slate-500 hover:text-slate-700 hover:bg-slate-100'
            }`}
          >
            <Settings size={12} />
            Settings
          </button>
        </nav>

        {/* UI sub-tabs */}
        {viewTab === 'ui' && tabs.length > 1 && (
          <>
            <div className="h-4 w-px bg-slate-200" />
            <nav className="flex gap-1">
              {tabs.map((tab, idx) => (
                <button
                  key={idx}
                  onClick={() => setActiveUiTab(idx)}
                  className={`rounded-md px-3 py-1.5 text-xs font-medium transition-colors ${
                    activeUiTab === idx
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

        {/* Spacer */}
        <div className="flex-1" />

        {/* Enable/Disable */}
        <button
          onClick={() => disableMutation.mutate()}
          disabled={disableMutation.isPending || enableMutation.isPending}
          className="inline-flex items-center gap-1.5 rounded-lg border border-yellow-200 px-3 py-1.5 text-xs font-medium text-yellow-700 hover:bg-yellow-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          <PowerOff size={12} />
          {disableMutation.isPending ? 'Disabling…' : 'Disable'}
        </button>
        <button
          onClick={() => enableMutation.mutate()}
          disabled={enableMutation.isPending || disableMutation.isPending}
          className="inline-flex items-center gap-1.5 rounded-lg border border-green-200 px-3 py-1.5 text-xs font-medium text-green-700 hover:bg-green-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          <Power size={12} />
          {enableMutation.isPending ? 'Enabling…' : 'Enable'}
        </button>
      </div>

      {/* Content area */}
      <div className="flex-1 overflow-hidden">
        {viewTab === 'widget' ? (
          /* ── Frontend Widget section ── */
          <div className="h-full overflow-y-auto p-6">
            <div className="mx-auto max-w-2xl space-y-5">
              <div>
                <h2 className="text-sm font-semibold text-slate-800">Frontend Widget</h2>
                <p className="mt-1 text-xs text-slate-500">
                  Frontend widgets are JavaScript bundles that get automatically injected into the
                  storefront when this plugin is active.
                </p>
              </div>

              {/* Widget bundle URL */}
              <div className="rounded-xl border border-slate-200 bg-white p-5 space-y-4">
                <div className="flex items-start gap-3">
                  <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-indigo-100 text-indigo-600">
                    <Layout size={15} />
                  </div>
                  <div className="flex-1">
                    <p className="text-sm font-semibold text-slate-800">Frontend Bundle</p>
                    {latestVersion?.frontend_url ? (
                      <div className="mt-2 space-y-2">
                        <div className="flex items-center gap-2 rounded-lg border border-slate-200 bg-slate-50 px-3 py-2">
                          <code className="flex-1 truncate font-mono text-xs text-slate-700">
                            {latestVersion.frontend_url}
                          </code>
                          <a
                            href={latestVersion.frontend_url}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="shrink-0 text-indigo-500 hover:text-indigo-700 transition-colors"
                            title="Open bundle URL"
                          >
                            <ExternalLink size={13} />
                          </a>
                        </div>
                        <p className="text-xs text-slate-500">
                          Version {latestVersion.version} · JS bundle loaded by the storefront
                          renderer
                        </p>
                      </div>
                    ) : (
                      <p className="mt-1 text-xs text-slate-400">
                        No frontend bundle — this plugin runs server-side only.
                      </p>
                    )}
                  </div>
                </div>
              </div>

              {/* Widget widgets (slots) */}
              {(manifest?.ui?.widgets?.length ?? 0) > 0 && (
                <div className="rounded-xl border border-slate-200 bg-white p-5 space-y-3">
                  <p className="text-sm font-semibold text-slate-800">Widget Slots</p>
                  <p className="text-xs text-slate-500">
                    These components will be rendered in the specified storefront slots.
                  </p>
                  <div className="divide-y divide-slate-100 rounded-lg border border-slate-200">
                    {manifest!.ui.widgets.map((w, i) => (
                      <div key={i} className="flex items-center gap-3 px-4 py-3">
                        <code className="rounded bg-indigo-50 px-2 py-0.5 font-mono text-xs text-indigo-700">
                          {w.slot}
                        </code>
                        <span className="flex-1 truncate text-xs text-slate-600">{w.component}</span>
                        <code className="truncate max-w-[180px] font-mono text-xs text-slate-400">
                          {w.entry}
                        </code>
                      </div>
                    ))}
                  </div>
                </div>
              )}

              {/* Injection status */}
              <div className="rounded-xl border border-slate-200 bg-white p-5">
                <div className="flex items-center gap-3">
                  {enableMutation.isPending || disableMutation.isPending ? (
                    <Loader2 size={16} className="animate-spin text-slate-400" />
                  ) : (
                    <div className="text-slate-500">
                      {/* We don't have live installation status here, show based on latest version */}
                      <Circle size={16} className="text-slate-300" />
                    </div>
                  )}
                  <div>
                    <p className="text-sm font-semibold text-slate-800">Widget Injection Status</p>
                    <p className="mt-0.5 text-xs text-slate-500">
                      {latestVersion?.frontend_url
                        ? 'This plugin has a frontend widget. It will be automatically loaded on the storefront when the plugin is enabled.'
                        : 'No frontend widget — enabling this plugin will not inject any scripts into the storefront.'}
                    </p>
                  </div>
                </div>
              </div>

              {/* Note */}
              <div className="rounded-lg border border-blue-100 bg-blue-50 px-4 py-3 text-xs text-blue-700">
                <strong>Note:</strong> Frontend widgets are automatically loaded by the storefront
                renderer when this plugin is active. No manual configuration is required.
              </div>
            </div>
          </div>
        ) : viewTab === 'settings' ? (
          /* ── Settings editor ── */
          <div className="h-full overflow-y-auto p-6">
            <div className="mx-auto max-w-2xl space-y-4">
              <div>
                <h2 className="text-sm font-semibold text-slate-800">Plugin Settings</h2>
                <p className="mt-1 text-xs text-slate-500">
                  Edit the JSON configuration for this plugin. Changes are applied immediately on save.
                </p>
              </div>
              <div className="rounded-xl border border-slate-200 bg-white p-4">
                <label className="mb-2 block text-xs font-medium text-slate-600">
                  Settings (JSON)
                </label>
                <textarea
                  rows={16}
                  value={settingsJson}
                  onChange={(e) => {
                    setSettingsJson(e.target.value)
                    setSettingsError('')
                  }}
                  placeholder='{\n  "key": "value"\n}'
                  className="w-full rounded-lg border border-slate-200 bg-slate-50 px-3 py-2 font-mono text-xs text-slate-800 focus:border-indigo-400 focus:outline-none focus:ring-2 focus:ring-indigo-100"
                />
                {settingsError && (
                  <p className="mt-2 text-xs text-red-600">{settingsError}</p>
                )}
              </div>
              <div className="flex justify-end">
                <button
                  onClick={() => void handleSaveSettings()}
                  disabled={savingSettings}
                  className="inline-flex items-center gap-1.5 rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                >
                  {savingSettings ? (
                    <Loader2 size={14} className="animate-spin" />
                  ) : null}
                  {savingSettings ? 'Saving…' : 'Save Settings'}
                </button>
              </div>
            </div>
          </div>
        ) : tabs.length === 0 ? (
          /* ── No UI tabs ── */
          <div className="flex h-full flex-col items-center justify-center gap-3">
            <Puzzle size={28} className="text-slate-300" />
            <p className="text-sm text-slate-500">This plugin has no UI tabs.</p>
            <button
              onClick={() => setViewTab('settings')}
              className="text-xs font-medium text-indigo-600 hover:underline"
            >
              Open Settings →
            </button>
          </div>
        ) : activeEntry ? (
          /* ── iframe for plugin UI ── */
          <iframe
            key={activeEntry}
            src={activeEntry}
            sandbox="allow-scripts allow-same-origin"
            className="h-full w-full border-0"
            title={tabs[activeUiTab]?.label ?? manifest.display_name}
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
