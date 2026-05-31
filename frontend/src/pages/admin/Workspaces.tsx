import { useState } from 'react'
import { useNavigate } from '@tanstack/react-router'
import {
  Layers,
  Plus,
  PlayCircle,
  StopCircle,
  Clock,
  AlertCircle,
  Bot,
  X,
  Loader2,
} from 'lucide-react'
import {
  useAgentSessions,
  useCreateAgentSession,
  useStopAgentSession,
  useInstalledPresets,
  useMarketplacePresets,
} from '../../lib/hooks'
import type { AgentSession, SessionStatus, Preset } from '../../types'

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

function statusBadgeClass(status: SessionStatus): string {
  switch (status) {
    case 'creating':
      return 'bg-yellow-100 text-yellow-700'
    case 'running':
      return 'bg-green-100 text-green-700'
    case 'stopped':
      return 'bg-slate-100 text-slate-600'
    case 'failed':
      return 'bg-red-100 text-red-700'
  }
}

function StatusIcon({ status }: { status: SessionStatus }) {
  switch (status) {
    case 'creating':
      return <Loader2 size={12} className="animate-spin" />
    case 'running':
      return <PlayCircle size={12} />
    case 'stopped':
      return <StopCircle size={12} />
    case 'failed':
      return <AlertCircle size={12} />
  }
}

function formatRelativeTime(dateStr: string): string {
  const diff = Date.now() - new Date(dateStr).getTime()
  const minutes = Math.floor(diff / 60_000)
  if (minutes < 1) return 'just now'
  if (minutes < 60) return `${minutes}m ago`
  const hours = Math.floor(minutes / 60)
  if (hours < 24) return `${hours}h ago`
  return `${Math.floor(hours / 24)}d ago`
}

// ---------------------------------------------------------------------------
// Session card
// ---------------------------------------------------------------------------

interface SessionCardProps {
  session: AgentSession
  preset?: Preset
  onOpen: (id: string) => void
  onStop: (id: string) => void
  isStopping: boolean
}

function SessionCard({ session, preset, onOpen, onStop, isStopping }: SessionCardProps) {
  const canOpen = session.status === 'running'
  const canStop = session.status === 'running' || session.status === 'creating'

  return (
    <div
      className={`flex flex-col rounded-xl border bg-white p-5 shadow-sm transition-shadow ${
        canOpen ? 'border-slate-200 hover:shadow-md cursor-pointer' : 'border-slate-200 opacity-75'
      }`}
      onClick={() => { if (canOpen) onOpen(session.id) }}
    >
      {/* Header */}
      <div className="flex items-start justify-between gap-3">
        <div className="flex items-center gap-3 min-w-0">
          <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-indigo-100 text-indigo-600">
            <Bot size={18} />
          </div>
          <div className="min-w-0">
            <p className="text-sm font-semibold text-slate-800 truncate">{session.name || 'Unnamed Workspace'}</p>
            <p className="text-xs text-slate-400 truncate">{preset?.name ?? session.preset_id}</p>
          </div>
        </div>
        <span
          className={`inline-flex shrink-0 items-center gap-1.5 rounded-full px-2.5 py-1 text-xs font-medium capitalize ${statusBadgeClass(session.status)}`}
        >
          <StatusIcon status={session.status} />
          {session.status}
        </span>
      </div>

      {/* Metadata */}
      <div className="mt-3 flex items-center gap-1.5 text-xs text-slate-400">
        <Clock size={11} />
        <span>Created {formatRelativeTime(session.created_at)}</span>
        {session.stopped_at && (
          <>
            <span>·</span>
            <span>Stopped {formatRelativeTime(session.stopped_at)}</span>
          </>
        )}
      </div>

      {/* Actions */}
      <div className="mt-4 flex items-center justify-end gap-2" onClick={(e) => e.stopPropagation()}>
        {canOpen && (
          <button
            onClick={() => onOpen(session.id)}
            className="inline-flex items-center gap-1.5 rounded-lg bg-indigo-600 px-3 py-1.5 text-xs font-medium text-white hover:bg-indigo-700 transition-colors"
          >
            <PlayCircle size={12} />
            Open
          </button>
        )}
        {canStop && (
          <button
            onClick={() => onStop(session.id)}
            disabled={isStopping}
            className="inline-flex items-center gap-1.5 rounded-lg border border-red-200 px-3 py-1.5 text-xs font-medium text-red-600 hover:bg-red-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            <StopCircle size={12} />
            {isStopping ? 'Stopping…' : 'Stop'}
          </button>
        )}
      </div>
    </div>
  )
}

// ---------------------------------------------------------------------------
// New Workspace Modal
// ---------------------------------------------------------------------------

interface NewWorkspaceModalProps {
  onClose: () => void
  presets: Preset[]
  isCreating: boolean
  onCreate: (presetId: string, name: string) => void
}

function NewWorkspaceModal({ onClose, presets, isCreating, onCreate }: NewWorkspaceModalProps) {
  const [selectedPresetId, setSelectedPresetId] = useState('')
  const [name, setName] = useState('')

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (!selectedPresetId) return
    onCreate(selectedPresetId, name)
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40">
      <div className="w-full max-w-md rounded-2xl bg-white p-6 shadow-2xl">
        <div className="flex items-center justify-between mb-5">
          <h2 className="text-base font-semibold text-slate-800">New Workspace</h2>
          <button
            onClick={onClose}
            className="flex h-7 w-7 items-center justify-center rounded-lg text-slate-400 hover:bg-slate-100 hover:text-slate-600 transition-colors"
          >
            <X size={16} />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          {/* Preset selection */}
          <div>
            <label className="mb-1.5 block text-xs font-medium text-slate-700">
              Preset <span className="text-red-500">*</span>
            </label>
            {presets.length === 0 ? (
              <div className="rounded-lg border border-amber-200 bg-amber-50 p-3 text-xs text-amber-700">
                No presets installed yet.{' '}
                <a href="/admin/presets" className="font-medium underline">
                  Browse the marketplace
                </a>{' '}
                to install one first.
              </div>
            ) : (
              <div className="grid grid-cols-1 gap-2 max-h-48 overflow-y-auto">
                {presets.map((preset) => (
                  <label
                    key={preset.id}
                    className={`flex cursor-pointer items-center gap-3 rounded-lg border p-3 transition-colors ${
                      selectedPresetId === preset.id
                        ? 'border-indigo-300 bg-indigo-50'
                        : 'border-slate-200 hover:bg-slate-50'
                    }`}
                  >
                    <input
                      type="radio"
                      name="preset"
                      value={preset.id}
                      checked={selectedPresetId === preset.id}
                      onChange={(e) => setSelectedPresetId(e.target.value)}
                      className="accent-indigo-600"
                    />
                    <div className="min-w-0">
                      <p className="text-sm font-medium text-slate-800">{preset.name}</p>
                      <p className="text-xs text-slate-400 truncate">{preset.description}</p>
                    </div>
                  </label>
                ))}
              </div>
            )}
          </div>

          {/* Name */}
          <div>
            <label className="mb-1.5 block text-xs font-medium text-slate-700">
              Workspace name <span className="text-slate-400">(optional)</span>
            </label>
            <input
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="e.g. Q4 Campaign Planning"
              className="w-full rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm text-slate-700 placeholder:text-slate-400 focus:border-indigo-400 focus:outline-none focus:ring-2 focus:ring-indigo-100"
            />
          </div>

          {/* Actions */}
          <div className="flex items-center justify-end gap-2 pt-1">
            <button
              type="button"
              onClick={onClose}
              className="rounded-lg border border-slate-200 px-4 py-2 text-sm font-medium text-slate-600 hover:bg-slate-50 transition-colors"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={!selectedPresetId || isCreating}
              className="inline-flex items-center gap-2 rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              {isCreating && <Loader2 size={14} className="animate-spin" />}
              {isCreating ? 'Creating…' : 'Create Workspace'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}

// ---------------------------------------------------------------------------
// Main page
// ---------------------------------------------------------------------------

export default function Workspaces() {
  const navigate = useNavigate()
  const [showModal, setShowModal] = useState(false)
  const [stoppingId, setStoppingId] = useState<string | null>(null)

  const { data: sessionsPage, isLoading: loadingSessions, error: sessionsError } = useAgentSessions()
  const { data: installedPage } = useInstalledPresets()
  const { data: presetsPage } = useMarketplacePresets()
  const createSession = useCreateAgentSession()
  const stopSession = useStopAgentSession()

  const sessions = sessionsPage?.items ?? []
  const installed = installedPage?.items ?? []
  const allPresets = presetsPage?.items ?? []

  // Get presets that the tenant has installed
  const installedPresetIds = new Set(installed.map((i) => i.preset_id))
  const installedPresets = allPresets.filter((p) => installedPresetIds.has(p.id))
  const presetMap = new Map<string, Preset>(allPresets.map((p) => [p.id, p]))

  function handleOpen(sessionId: string) {
    void navigate({ to: '/admin/workspaces/$id', params: { id: sessionId } })
  }

  function handleStop(sessionId: string) {
    setStoppingId(sessionId)
    stopSession.mutate(sessionId, {
      onSettled: () => setStoppingId(null),
    })
  }

  function handleCreate(presetId: string, name: string) {
    createSession.mutate(
      { preset_id: presetId, name: name || undefined },
      {
        onSuccess: (session) => {
          setShowModal(false)
          void navigate({ to: '/admin/workspaces/$id', params: { id: session.id } })
        },
      },
    )
  }

  // Separate running from stopped
  const activeSessions = sessions.filter((s) => s.status === 'running' || s.status === 'creating')
  const pastSessions = sessions.filter((s) => s.status === 'stopped' || s.status === 'failed')

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <Layers size={24} className="text-indigo-600" />
          <div>
            <h1 className="text-xl font-semibold text-slate-800">Workspaces</h1>
            <p className="text-sm text-slate-500">Launch and manage AI agent workspace sessions</p>
          </div>
        </div>
        <button
          onClick={() => setShowModal(true)}
          className="inline-flex items-center gap-2 rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700 transition-colors shadow-sm"
        >
          <Plus size={16} />
          New Workspace
        </button>
      </div>

      {/* Content */}
      {loadingSessions ? (
        <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
          {[...Array(3)].map((_, i) => (
            <div key={i} className="h-40 rounded-xl border border-slate-200 bg-white animate-pulse" />
          ))}
        </div>
      ) : sessionsError ? (
        <div className="rounded-xl border border-red-200 bg-red-50 p-6 text-center text-sm text-red-600">
          Failed to load workspaces. Please try again.
        </div>
      ) : sessions.length === 0 ? (
        <div className="rounded-xl border border-slate-200 bg-white p-16 text-center">
          <Layers size={36} className="mx-auto mb-4 text-slate-300" />
          <p className="text-sm font-medium text-slate-600">No workspaces yet</p>
          <p className="mt-1 text-xs text-slate-400">Create a workspace to start working with an AI agent.</p>
          <button
            onClick={() => setShowModal(true)}
            className="mt-4 inline-flex items-center gap-2 rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700 transition-colors"
          >
            <Plus size={16} />
            New Workspace
          </button>
        </div>
      ) : (
        <div className="space-y-6">
          {/* Active sessions */}
          {activeSessions.length > 0 && (
            <div>
              <h2 className="mb-3 text-xs font-semibold uppercase tracking-wider text-slate-400">Active</h2>
              <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
                {activeSessions.map((session) => (
                  <SessionCard
                    key={session.id}
                    session={session}
                    preset={presetMap.get(session.preset_id)}
                    onOpen={handleOpen}
                    onStop={handleStop}
                    isStopping={stoppingId === session.id}
                  />
                ))}
              </div>
            </div>
          )}

          {/* Past sessions */}
          {pastSessions.length > 0 && (
            <div>
              <h2 className="mb-3 text-xs font-semibold uppercase tracking-wider text-slate-400">Recent</h2>
              <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
                {pastSessions.map((session) => (
                  <SessionCard
                    key={session.id}
                    session={session}
                    preset={presetMap.get(session.preset_id)}
                    onOpen={handleOpen}
                    onStop={handleStop}
                    isStopping={stoppingId === session.id}
                  />
                ))}
              </div>
            </div>
          )}
        </div>
      )}

      {/* New Workspace Modal */}
      {showModal && (
        <NewWorkspaceModal
          onClose={() => setShowModal(false)}
          presets={installedPresets}
          isCreating={createSession.isPending}
          onCreate={handleCreate}
        />
      )}
    </div>
  )
}
