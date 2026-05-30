import { useState } from 'react'
import { useNavigate } from '@tanstack/react-router'
import {
  Puzzle,
  Download,
  Trash2,
  Settings,
  PowerOff,
  Power,
  Search,
} from 'lucide-react'
import {
  useMarketplacePlugins,
  useInstalledPlugins,
  useInstallPlugin,
  useUninstallPlugin,
  useEnablePlugin,
  useDisablePlugin,
} from '../../lib/hooks'
import type { Plugin, PluginCategory, InstallationStatus } from '../../types'

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

type ViewTab = 'browse' | 'installed'
type CategoryFilter = 'all' | PluginCategory

function categoryBadgeClass(category: PluginCategory): string {
  switch (category) {
    case 'official':
      return 'bg-indigo-100 text-indigo-700'
    case 'community':
      return 'bg-green-100 text-green-700'
    case 'custom':
      return 'bg-gray-100 text-gray-600'
  }
}

function statusBadgeClass(status: InstallationStatus): string {
  switch (status) {
    case 'active':
      return 'bg-green-100 text-green-700'
    case 'inactive':
      return 'bg-yellow-100 text-yellow-700'
    case 'failed':
      return 'bg-red-100 text-red-700'
  }
}

function pluginIconColor(category: PluginCategory): string {
  switch (category) {
    case 'official':
      return 'bg-indigo-500'
    case 'community':
      return 'bg-green-500'
    case 'custom':
      return 'bg-gray-400'
  }
}

function formatDate(dateStr: string): string {
  try {
    return new Date(dateStr).toLocaleDateString()
  } catch {
    return dateStr
  }
}

// ---------------------------------------------------------------------------
// Plugin card (Browse view)
// ---------------------------------------------------------------------------

interface PluginCardProps {
  plugin: Plugin
  isInstalled: boolean
  isInstalling: boolean
  onInstall: (id: string) => void
}

function PluginCard({ plugin, isInstalled, isInstalling, onInstall }: PluginCardProps) {
  return (
    <div className="flex flex-col rounded-xl border border-slate-200 bg-white p-5 shadow-sm hover:shadow-md transition-shadow">
      <div className="flex items-start gap-4">
        {/* Icon */}
        <div
          className={`flex h-12 w-12 shrink-0 items-center justify-center rounded-xl text-white font-bold text-lg ${pluginIconColor(plugin.category)}`}
        >
          {plugin.display_name.charAt(0).toUpperCase()}
        </div>

        <div className="min-w-0 flex-1">
          <div className="flex items-center gap-2 flex-wrap">
            <h3 className="text-sm font-semibold text-slate-800 truncate">{plugin.display_name}</h3>
            <span
              className={`inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium capitalize ${categoryBadgeClass(plugin.category)}`}
            >
              {plugin.category}
            </span>
          </div>
          <p className="mt-0.5 text-xs text-slate-500">by {plugin.author}</p>
        </div>
      </div>

      <p className="mt-3 text-sm text-slate-600 line-clamp-2 flex-1">{plugin.description}</p>

      {/* Tags */}
      {plugin.tags.length > 0 && (
        <div className="mt-3 flex flex-wrap gap-1">
          {plugin.tags.slice(0, 4).map((tag) => (
            <span key={tag} className="rounded-full bg-slate-100 px-2 py-0.5 text-xs text-slate-500">
              {tag}
            </span>
          ))}
        </div>
      )}

      {/* Action */}
      <div className="mt-4 flex items-center justify-end">
        {isInstalled ? (
          <span className="inline-flex items-center gap-1.5 rounded-full bg-green-100 px-3 py-1.5 text-xs font-medium text-green-700">
            <Puzzle size={12} />
            Installed
          </span>
        ) : (
          <button
            onClick={() => onInstall(plugin.id)}
            disabled={isInstalling}
            className="inline-flex items-center gap-1.5 rounded-lg bg-indigo-600 px-3 py-1.5 text-xs font-medium text-white hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            <Download size={12} />
            {isInstalling ? 'Installing…' : 'Install'}
          </button>
        )}
      </div>
    </div>
  )
}

// ---------------------------------------------------------------------------
// Installed row
// ---------------------------------------------------------------------------

interface InstalledRowProps {
  pluginId: string
  versionId: string
  status: InstallationStatus
  installedAt: string
  isUninstalling: boolean
  isTogglingStatus: boolean
  onUninstall: (id: string) => void
  onOpen: (id: string) => void
  onEnable: (id: string) => void
  onDisable: (id: string) => void
}

function InstalledRow({
  pluginId,
  versionId,
  status,
  installedAt,
  isUninstalling,
  isTogglingStatus,
  onUninstall,
  onOpen,
  onEnable,
  onDisable,
}: InstalledRowProps) {
  return (
    <tr className="border-b border-slate-100 last:border-0 hover:bg-slate-50 transition-colors">
      <td className="py-3 px-4">
        <div className="flex items-center gap-3">
          <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-indigo-100 text-indigo-600">
            <Puzzle size={14} />
          </div>
          <div>
            <p className="text-sm font-medium text-slate-800 font-mono">{pluginId}</p>
            <p className="text-xs text-slate-400">v{versionId}</p>
          </div>
        </div>
      </td>
      <td className="py-3 px-4">
        <span
          className={`inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium capitalize ${statusBadgeClass(status)}`}
        >
          {status}
        </span>
      </td>
      <td className="py-3 px-4 text-sm text-slate-500">{formatDate(installedAt)}</td>
      <td className="py-3 px-4">
        <div className="flex items-center justify-end gap-2">
          <button
            onClick={() => onOpen(pluginId)}
            className="inline-flex items-center gap-1.5 rounded-lg border border-slate-200 px-2.5 py-1.5 text-xs font-medium text-slate-600 hover:bg-slate-100 transition-colors"
          >
            <Settings size={12} />
            Open
          </button>
          {status === 'active' ? (
            <button
              onClick={() => onDisable(pluginId)}
              disabled={isTogglingStatus}
              className="inline-flex items-center gap-1.5 rounded-lg border border-yellow-200 px-2.5 py-1.5 text-xs font-medium text-yellow-700 hover:bg-yellow-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              <PowerOff size={12} />
              {isTogglingStatus ? 'Disabling…' : 'Disable'}
            </button>
          ) : (
            <button
              onClick={() => onEnable(pluginId)}
              disabled={isTogglingStatus}
              className="inline-flex items-center gap-1.5 rounded-lg border border-green-200 px-2.5 py-1.5 text-xs font-medium text-green-700 hover:bg-green-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              <Power size={12} />
              {isTogglingStatus ? 'Enabling…' : 'Enable'}
            </button>
          )}
          <button
            onClick={() => onUninstall(pluginId)}
            disabled={isUninstalling}
            className="inline-flex items-center gap-1.5 rounded-lg border border-red-200 px-2.5 py-1.5 text-xs font-medium text-red-600 hover:bg-red-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            <Trash2 size={12} />
            {isUninstalling ? 'Removing…' : 'Uninstall'}
          </button>
        </div>
      </td>
    </tr>
  )
}

// ---------------------------------------------------------------------------
// Main page
// ---------------------------------------------------------------------------

export default function Marketplace() {
  const navigate = useNavigate()
  const [view, setView] = useState<ViewTab>('browse')
  const [categoryFilter, setCategoryFilter] = useState<CategoryFilter>('all')
  const [search, setSearch] = useState('')
  const [installingId, setInstallingId] = useState<string | null>(null)
  const [uninstallingId, setUninstallingId] = useState<string | null>(null)
  const [togglingId, setTogglingId] = useState<string | null>(null)

  const { data: pluginsPage, isLoading: loadingPlugins, error: pluginsError } = useMarketplacePlugins()
  const { data: installedPage, isLoading: loadingInstalled } = useInstalledPlugins()
  const installPlugin = useInstallPlugin()
  const uninstallPlugin = useUninstallPlugin()
  const enablePlugin = useEnablePlugin()
  const disablePlugin = useDisablePlugin()

  const installed = installedPage?.items ?? []
  const installedIds = new Set(installed.map((i) => i.plugin_id))

  // Filter plugins
  const allPlugins: Plugin[] = pluginsPage?.items ?? []
  const filteredPlugins = allPlugins.filter((p) => {
    if (categoryFilter !== 'all' && p.category !== categoryFilter) return false
    if (search.trim()) {
      const q = search.toLowerCase()
      return (
        p.display_name.toLowerCase().includes(q) ||
        p.description.toLowerCase().includes(q) ||
        p.author.toLowerCase().includes(q) ||
        p.tags.some((t) => t.toLowerCase().includes(q))
      )
    }
    return true
  })

  function handleInstall(pluginId: string) {
    setInstallingId(pluginId)
    installPlugin.mutate(pluginId, {
      onSettled: () => setInstallingId(null),
    })
  }

  function handleUninstall(pluginId: string) {
    setUninstallingId(pluginId)
    uninstallPlugin.mutate(pluginId, {
      onSettled: () => setUninstallingId(null),
    })
  }

  function handleEnable(pluginId: string) {
    setTogglingId(pluginId)
    enablePlugin.mutate(pluginId, {
      onSettled: () => setTogglingId(null),
    })
  }

  function handleDisable(pluginId: string) {
    setTogglingId(pluginId)
    disablePlugin.mutate(pluginId, {
      onSettled: () => setTogglingId(null),
    })
  }

  function handleOpen(pluginId: string) {
    void navigate({ to: '/admin/plugins/$name', params: { name: pluginId } })
  }

  return (
    <div className="space-y-6">
      {/* Page header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <Puzzle size={24} className="text-indigo-600" />
          <div>
            <h1 className="text-xl font-semibold text-slate-800">Marketplace</h1>
            <p className="text-sm text-slate-500">Browse and manage plugins for your store</p>
          </div>
        </div>
      </div>

      {/* View tabs */}
      <div className="flex gap-1 rounded-lg border border-slate-200 bg-white p-1 w-fit">
        {(['browse', 'installed'] as const).map((tab) => (
          <button
            key={tab}
            onClick={() => setView(tab)}
            className={`rounded-md px-4 py-1.5 text-sm font-medium capitalize transition-colors ${
              view === tab
                ? 'bg-indigo-600 text-white shadow-sm'
                : 'text-slate-600 hover:text-slate-800 hover:bg-slate-100'
            }`}
          >
            {tab}
            {tab === 'installed' && installed.length > 0 && (
              <span className="ml-1.5 rounded-full bg-indigo-500/20 px-1.5 text-xs">
                {installedPage?.total ?? installed.length}
              </span>
            )}
          </button>
        ))}
      </div>

      {/* ── Browse view ── */}
      {view === 'browse' && (
        <div className="space-y-4">
          {/* Filters row */}
          <div className="flex flex-col sm:flex-row gap-3">
            {/* Search */}
            <div className="relative flex-1 max-w-sm">
              <Search size={14} className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
              <input
                type="text"
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                placeholder="Search plugins…"
                className="w-full rounded-lg border border-slate-200 bg-white pl-8 pr-3 py-2 text-sm text-slate-700 placeholder:text-slate-400 focus:border-indigo-400 focus:outline-none focus:ring-2 focus:ring-indigo-100"
              />
            </div>

            {/* Category tabs */}
            <div className="flex gap-1 rounded-lg border border-slate-200 bg-white p-1">
              {(['all', 'official', 'community'] as const).map((cat) => (
                <button
                  key={cat}
                  onClick={() => setCategoryFilter(cat)}
                  className={`rounded-md px-3 py-1.5 text-xs font-medium capitalize transition-colors ${
                    categoryFilter === cat
                      ? 'bg-slate-800 text-white'
                      : 'text-slate-500 hover:text-slate-700 hover:bg-slate-100'
                  }`}
                >
                  {cat}
                </button>
              ))}
            </div>
          </div>

          {/* Plugin grid */}
          {loadingPlugins ? (
            <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
              {[...Array(6)].map((_, i) => (
                <div key={i} className="h-48 rounded-xl border border-slate-200 bg-white animate-pulse" />
              ))}
            </div>
          ) : pluginsError ? (
            <div className="rounded-xl border border-red-200 bg-red-50 p-6 text-center text-sm text-red-600">
              Failed to load plugins. Please try again.
            </div>
          ) : filteredPlugins.length === 0 ? (
            <div className="rounded-xl border border-slate-200 bg-white p-12 text-center">
              <Puzzle size={32} className="mx-auto mb-3 text-slate-300" />
              <p className="text-sm text-slate-500">No plugins found</p>
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
              {filteredPlugins.map((plugin) => (
                <PluginCard
                  key={plugin.id}
                  plugin={plugin}
                  isInstalled={installedIds.has(plugin.id)}
                  isInstalling={installingId === plugin.id}
                  onInstall={handleInstall}
                />
              ))}
            </div>
          )}
        </div>
      )}

      {/* ── Installed view ── */}
      {view === 'installed' && (
        <div className="rounded-xl border border-slate-200 bg-white overflow-hidden">
          {loadingInstalled ? (
            <div className="p-8 text-center text-sm text-slate-500">Loading…</div>
          ) : installed.length === 0 ? (
            <div className="p-12 text-center">
              <Puzzle size={32} className="mx-auto mb-3 text-slate-300" />
              <p className="text-sm text-slate-500">No plugins installed yet</p>
              <button
                onClick={() => setView('browse')}
                className="mt-3 text-xs font-medium text-indigo-600 hover:underline"
              >
                Browse marketplace →
              </button>
            </div>
          ) : (
            <table className="w-full">
              <thead>
                <tr className="border-b border-slate-200 bg-slate-50">
                  <th className="py-3 px-4 text-left text-xs font-semibold uppercase tracking-wider text-slate-500">
                    Plugin
                  </th>
                  <th className="py-3 px-4 text-left text-xs font-semibold uppercase tracking-wider text-slate-500">
                    Status
                  </th>
                  <th className="py-3 px-4 text-left text-xs font-semibold uppercase tracking-wider text-slate-500">
                    Installed
                  </th>
                  <th className="py-3 px-4 text-right text-xs font-semibold uppercase tracking-wider text-slate-500">
                    Actions
                  </th>
                </tr>
              </thead>
              <tbody>
                {installed.map((inst) => (
                  <InstalledRow
                    key={inst.id}
                    pluginId={inst.plugin_id}
                    versionId={inst.version_id}
                    status={inst.status}
                    installedAt={inst.installed_at}
                    isUninstalling={uninstallingId === inst.plugin_id}
                    isTogglingStatus={togglingId === inst.plugin_id}
                    onUninstall={handleUninstall}
                    onOpen={handleOpen}
                    onEnable={handleEnable}
                    onDisable={handleDisable}
                  />
                ))}
              </tbody>
            </table>
          )}
        </div>
      )}
    </div>
  )
}
