import { useState } from 'react'
import {
  Sparkles,
  Download,
  Trash2,
  Search,
  Bot,
  CheckCircle2,
} from 'lucide-react'
import {
  useMarketplacePresets,
  useInstalledPresets,
  useInstallPreset,
  useUninstallPreset,
} from '../../lib/hooks'
import type { Preset, PresetCategory, PresetInstall } from '../../types'

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

type ViewTab = 'browse' | 'installed'
type CategoryFilter = 'all' | PresetCategory

function categoryBadgeClass(category: PresetCategory): string {
  switch (category) {
    case 'webdev':
      return 'bg-blue-100 text-blue-700'
    case 'research':
      return 'bg-purple-100 text-purple-700'
    case 'content':
      return 'bg-green-100 text-green-700'
    case 'analytics':
      return 'bg-amber-100 text-amber-700'
    case 'store-manager':
      return 'bg-indigo-100 text-indigo-700'
    default:
      return 'bg-slate-100 text-slate-600'
  }
}

function categoryIconColor(category: PresetCategory): string {
  switch (category) {
    case 'webdev':
      return 'bg-blue-500'
    case 'research':
      return 'bg-purple-500'
    case 'content':
      return 'bg-green-500'
    case 'analytics':
      return 'bg-amber-500'
    case 'store-manager':
      return 'bg-indigo-500'
    default:
      return 'bg-slate-400'
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
// Preset card (Browse view)
// ---------------------------------------------------------------------------

interface PresetCardProps {
  preset: Preset
  isInstalled: boolean
  isInstalling: boolean
  onInstall: (id: string) => void
}

function PresetCard({ preset, isInstalled, isInstalling, onInstall }: PresetCardProps) {
  return (
    <div className="flex flex-col rounded-xl border border-slate-200 bg-white p-5 shadow-sm hover:shadow-md transition-shadow">
      <div className="flex items-start gap-4">
        {/* Icon */}
        {preset.icon_url ? (
          <img
            src={preset.icon_url}
            alt={preset.name}
            className="h-12 w-12 shrink-0 rounded-xl object-cover"
          />
        ) : (
          <div
            className={`flex h-12 w-12 shrink-0 items-center justify-center rounded-xl text-white font-bold text-lg ${categoryIconColor(preset.category)}`}
          >
            <Bot size={22} />
          </div>
        )}

        <div className="min-w-0 flex-1">
          <div className="flex items-center gap-2 flex-wrap">
            <h3 className="text-sm font-semibold text-slate-800 truncate">{preset.name}</h3>
            <span
              className={`inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium capitalize ${categoryBadgeClass(preset.category)}`}
            >
              {preset.category}
            </span>
          </div>
          <p className="mt-0.5 text-xs text-slate-500">by {preset.author} · v{preset.version}</p>
        </div>
      </div>

      <p className="mt-3 text-sm text-slate-600 line-clamp-2 flex-1">{preset.description}</p>

      {/* Docker image */}
      <div className="mt-3">
        <span className="inline-flex items-center rounded-md bg-slate-100 px-2 py-0.5 font-mono text-xs text-slate-500 truncate max-w-full">
          {preset.image}
        </span>
      </div>

      {/* Action */}
      <div className="mt-4 flex items-center justify-end">
        {isInstalled ? (
          <span className="inline-flex items-center gap-1.5 rounded-full bg-green-100 px-3 py-1.5 text-xs font-medium text-green-700">
            <CheckCircle2 size={12} />
            Installed
          </span>
        ) : (
          <button
            onClick={() => onInstall(preset.id)}
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
  install: PresetInstall
  preset?: Preset
  isUninstalling: boolean
  onUninstall: (id: string) => void
}

function InstalledRow({ install, preset, isUninstalling, onUninstall }: InstalledRowProps) {
  return (
    <tr className="border-b border-slate-100 last:border-0 hover:bg-slate-50 transition-colors">
      <td className="py-3 px-4">
        <div className="flex items-center gap-3">
          <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-indigo-100 text-indigo-600">
            <Bot size={14} />
          </div>
          <div>
            <p className="text-sm font-medium text-slate-800">
              {preset?.name ?? install.preset_id}
            </p>
            {preset && (
              <p className="text-xs text-slate-400">v{preset.version} · {preset.category}</p>
            )}
          </div>
        </div>
      </td>
      <td className="py-3 px-4 text-sm text-slate-500">{formatDate(install.installed_at)}</td>
      <td className="py-3 px-4">
        <div className="flex items-center justify-end gap-2">
          <button
            onClick={() => onUninstall(install.preset_id)}
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

export default function Presets() {
  const [view, setView] = useState<ViewTab>('browse')
  const [categoryFilter, setCategoryFilter] = useState<CategoryFilter>('all')
  const [search, setSearch] = useState('')
  const [installingId, setInstallingId] = useState<string | null>(null)
  const [uninstallingId, setUninstallingId] = useState<string | null>(null)

  const { data: presetsPage, isLoading: loadingPresets, error: presetsError } = useMarketplacePresets()
  const { data: installedPage, isLoading: loadingInstalled } = useInstalledPresets()
  const installPreset = useInstallPreset()
  const uninstallPreset = useUninstallPreset()

  const installed = installedPage?.items ?? []
  const installedPresetIds = new Set(installed.map((i) => i.preset_id))

  // Filter presets
  const allPresets: Preset[] = presetsPage?.items ?? []
  const filteredPresets = allPresets.filter((p) => {
    if (categoryFilter !== 'all' && p.category !== categoryFilter) return false
    if (search.trim()) {
      const q = search.toLowerCase()
      return (
        p.name.toLowerCase().includes(q) ||
        p.description.toLowerCase().includes(q) ||
        p.author.toLowerCase().includes(q) ||
        p.category.toLowerCase().includes(q)
      )
    }
    return true
  })

  // Build a map of preset_id → Preset for the installed tab
  const presetMap = new Map<string, Preset>(allPresets.map((p) => [p.id, p]))

  function handleInstall(presetId: string) {
    setInstallingId(presetId)
    installPreset.mutate(presetId, {
      onSettled: () => setInstallingId(null),
    })
  }

  function handleUninstall(presetId: string) {
    setUninstallingId(presetId)
    uninstallPreset.mutate(presetId, {
      onSettled: () => setUninstallingId(null),
    })
  }

  return (
    <div className="space-y-6">
      {/* Page header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <Sparkles size={24} className="text-indigo-600" />
          <div>
            <h1 className="text-xl font-semibold text-slate-800">Preset Marketplace</h1>
            <p className="text-sm text-slate-500">Browse and install AI agent presets for your store</p>
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
                placeholder="Search presets…"
                className="w-full rounded-lg border border-slate-200 bg-white pl-8 pr-3 py-2 text-sm text-slate-700 placeholder:text-slate-400 focus:border-indigo-400 focus:outline-none focus:ring-2 focus:ring-indigo-100"
              />
            </div>

            {/* Category tabs */}
            <div className="flex gap-1 rounded-lg border border-slate-200 bg-white p-1 flex-wrap">
              {(['all', 'webdev', 'research', 'content', 'analytics', 'store-manager'] as const).map((cat) => (
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

          {/* Preset grid */}
          {loadingPresets ? (
            <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
              {[...Array(6)].map((_, i) => (
                <div key={i} className="h-48 rounded-xl border border-slate-200 bg-white animate-pulse" />
              ))}
            </div>
          ) : presetsError ? (
            <div className="rounded-xl border border-red-200 bg-red-50 p-6 text-center text-sm text-red-600">
              Failed to load presets. Please try again.
            </div>
          ) : filteredPresets.length === 0 ? (
            <div className="rounded-xl border border-slate-200 bg-white p-12 text-center">
              <Bot size={32} className="mx-auto mb-3 text-slate-300" />
              <p className="text-sm text-slate-500">No presets found</p>
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
              {filteredPresets.map((preset) => (
                <PresetCard
                  key={preset.id}
                  preset={preset}
                  isInstalled={installedPresetIds.has(preset.id)}
                  isInstalling={installingId === preset.id}
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
              <Bot size={32} className="mx-auto mb-3 text-slate-300" />
              <p className="text-sm text-slate-500">No presets installed yet</p>
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
                    Preset
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
                    install={inst}
                    preset={presetMap.get(inst.preset_id)}
                    isUninstalling={uninstallingId === inst.preset_id}
                    onUninstall={handleUninstall}
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
