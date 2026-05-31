import { useState } from 'react'
import {
  Zap,
  Plus,
  Trash2,
  ChevronLeft,
  CheckCircle2,
  XCircle,
  Clock,
  X,
  ChevronDown,
} from 'lucide-react'
import {
  useAgentTriggers,
  useTriggerEventTypes,
  useTriggerLogs,
  useCreateAgentTrigger,
  useDeleteAgentTrigger,
  useEnableAgentTrigger,
  useDisableAgentTrigger,
} from '../../lib/hooks'
import type { AgentTrigger, CreateTriggerRequest, TriggerLog } from '../../types'

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

function formatDate(dateStr: string): string {
  try {
    return new Date(dateStr).toLocaleString()
  } catch {
    return dateStr
  }
}

function formatDateShort(dateStr: string): string {
  try {
    return new Date(dateStr).toLocaleDateString()
  } catch {
    return dateStr
  }
}

function LogStatusBadge({ status }: { status: TriggerLog['status'] }) {
  const classes = {
    success: 'bg-green-100 text-green-700',
    error: 'bg-red-100 text-red-700',
    skipped_cooldown: 'bg-slate-100 text-slate-600',
  }
  const icons = {
    success: <CheckCircle2 size={11} className="shrink-0" />,
    error: <XCircle size={11} className="shrink-0" />,
    skipped_cooldown: <Clock size={11} className="shrink-0" />,
  }
  const labels = {
    success: 'Success',
    error: 'Error',
    skipped_cooldown: 'Cooldown',
  }
  return (
    <span
      className={`inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-xs font-medium ${classes[status]}`}
    >
      {icons[status]}
      {labels[status]}
    </span>
  )
}

// ---------------------------------------------------------------------------
// Payload viewer (collapsible)
// ---------------------------------------------------------------------------

function PayloadViewer({ data, label }: { data: Record<string, unknown>; label: string }) {
  const [expanded, setExpanded] = useState(false)
  return (
    <div>
      <button
        onClick={() => setExpanded(!expanded)}
        className="inline-flex items-center gap-1 text-xs text-indigo-600 hover:text-indigo-800 transition-colors"
      >
        <ChevronDown
          size={12}
          className={`transition-transform ${expanded ? '' : '-rotate-90'}`}
        />
        {label}
      </button>
      {expanded && (
        <pre className="mt-1 rounded-md bg-slate-900 p-2 text-xs text-slate-100 overflow-x-auto max-h-32 max-w-xs">
          {JSON.stringify(data, null, 2)}
        </pre>
      )}
    </div>
  )
}

// ---------------------------------------------------------------------------
// Create trigger dialog
// ---------------------------------------------------------------------------

interface CreateTriggerDialogProps {
  eventTypes: string[]
  onClose: () => void
}

function CreateTriggerDialog({ eventTypes, onClose }: CreateTriggerDialogProps) {
  const [name, setName] = useState('')
  const [eventType, setEventType] = useState(eventTypes[0] ?? '')
  const [prompt, setPrompt] = useState('')
  const [cooldown, setCooldown] = useState('0')

  const createMutation = useCreateAgentTrigger()

  function handleSubmit() {
    const data: CreateTriggerRequest = {
      name,
      event_type: eventType,
      prompt,
      cooldown: parseInt(cooldown, 10) || 0,
    }
    createMutation.mutate(data, { onSuccess: onClose })
  }

  const isValid = name.trim().length > 0 && eventType.length > 0 && prompt.trim().length > 0

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40">
      <div className="w-full max-w-lg rounded-xl bg-white shadow-xl p-6 mx-4 max-h-[90vh] overflow-y-auto">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-base font-semibold text-slate-800">New Trigger</h2>
          <button
            onClick={onClose}
            className="text-slate-400 hover:text-slate-600 transition-colors"
          >
            <X size={18} />
          </button>
        </div>

        <div className="space-y-4">
          {/* Name */}
          <div>
            <label className="block text-xs font-medium text-slate-600 mb-1">Name *</label>
            <input
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="Trigger name…"
              className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm text-slate-700 placeholder:text-slate-400 focus:border-indigo-400 focus:outline-none focus:ring-2 focus:ring-indigo-100"
            />
          </div>

          {/* Event Type */}
          <div>
            <label className="block text-xs font-medium text-slate-600 mb-1">Event Type *</label>
            {eventTypes.length > 0 ? (
              <select
                value={eventType}
                onChange={(e) => setEventType(e.target.value)}
                className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm text-slate-700 focus:border-indigo-400 focus:outline-none focus:ring-2 focus:ring-indigo-100 bg-white"
              >
                {eventTypes.map((et) => (
                  <option key={et} value={et}>
                    {et}
                  </option>
                ))}
              </select>
            ) : (
              <input
                type="text"
                value={eventType}
                onChange={(e) => setEventType(e.target.value)}
                placeholder="e.g. order.created"
                className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm text-slate-700 placeholder:text-slate-400 focus:border-indigo-400 focus:outline-none focus:ring-2 focus:ring-indigo-100"
              />
            )}
          </div>

          {/* Prompt */}
          <div>
            <label className="block text-xs font-medium text-slate-600 mb-1">
              Prompt *
              <span className="ml-1 font-normal text-slate-400">
                — use <code className="bg-slate-100 px-1 rounded text-xs">{`{{.EventPayload}}`}</code> to inject event data
              </span>
            </label>
            <textarea
              value={prompt}
              onChange={(e) => setPrompt(e.target.value)}
              placeholder={`Analyze this event and take appropriate action:\n{{.EventPayload}}`}
              rows={5}
              className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm text-slate-700 placeholder:text-slate-400 focus:border-indigo-400 focus:outline-none focus:ring-2 focus:ring-indigo-100 resize-none font-mono"
            />
          </div>

          {/* Cooldown */}
          <div>
            <label className="block text-xs font-medium text-slate-600 mb-1">
              Cooldown (seconds)
              <span className="ml-1 font-normal text-slate-400">— minimum time between firings</span>
            </label>
            <input
              type="number"
              value={cooldown}
              onChange={(e) => setCooldown(e.target.value)}
              min="0"
              className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm text-slate-700 focus:border-indigo-400 focus:outline-none focus:ring-2 focus:ring-indigo-100"
            />
          </div>
        </div>

        {createMutation.error && (
          <p className="mt-3 text-xs text-red-600">{createMutation.error.message}</p>
        )}

        <div className="mt-5 flex items-center justify-end gap-2">
          <button
            onClick={onClose}
            disabled={createMutation.isPending}
            className="rounded-lg border border-slate-200 px-4 py-2 text-sm font-medium text-slate-600 hover:bg-slate-50 disabled:opacity-50 transition-colors"
          >
            Cancel
          </button>
          <button
            onClick={handleSubmit}
            disabled={createMutation.isPending || !isValid}
            className="rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            {createMutation.isPending ? 'Creating…' : 'Create Trigger'}
          </button>
        </div>
      </div>
    </div>
  )
}

// ---------------------------------------------------------------------------
// Delete confirmation dialog
// ---------------------------------------------------------------------------

interface DeleteTriggerDialogProps {
  trigger: AgentTrigger
  onClose: () => void
}

function DeleteTriggerDialog({ trigger, onClose }: DeleteTriggerDialogProps) {
  const deleteMutation = useDeleteAgentTrigger()

  function handleDelete() {
    deleteMutation.mutate(trigger.id, { onSuccess: onClose })
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40">
      <div className="w-full max-w-sm rounded-xl bg-white shadow-xl p-6 mx-4">
        <h2 className="text-base font-semibold text-slate-800 mb-2">Delete Trigger</h2>
        <p className="text-sm text-slate-500 mb-1">
          Are you sure you want to delete{' '}
          <span className="font-medium text-slate-700">&ldquo;{trigger.name}&rdquo;</span>?
        </p>
        <p className="text-xs text-slate-400">This action cannot be undone.</p>

        {deleteMutation.error && (
          <p className="mt-2 text-xs text-red-600">{deleteMutation.error.message}</p>
        )}

        <div className="mt-5 flex items-center justify-end gap-2">
          <button
            onClick={onClose}
            disabled={deleteMutation.isPending}
            className="rounded-lg border border-slate-200 px-4 py-2 text-sm font-medium text-slate-600 hover:bg-slate-50 disabled:opacity-50 transition-colors"
          >
            Cancel
          </button>
          <button
            onClick={handleDelete}
            disabled={deleteMutation.isPending}
            className="rounded-lg bg-red-600 px-4 py-2 text-sm font-medium text-white hover:bg-red-700 disabled:opacity-50 transition-colors"
          >
            {deleteMutation.isPending ? 'Deleting…' : 'Delete'}
          </button>
        </div>
      </div>
    </div>
  )
}

// ---------------------------------------------------------------------------
// Toggle switch
// ---------------------------------------------------------------------------

interface ToggleProps {
  enabled: boolean
  onChange: (enabled: boolean) => void
  disabled?: boolean
}

function Toggle({ enabled, onChange, disabled }: ToggleProps) {
  return (
    <button
      role="switch"
      aria-checked={enabled}
      disabled={disabled}
      onClick={() => onChange(!enabled)}
      className={`relative inline-flex h-5 w-9 shrink-0 cursor-pointer rounded-full transition-colors focus:outline-none disabled:opacity-50 disabled:cursor-not-allowed ${
        enabled ? 'bg-indigo-600' : 'bg-slate-200'
      }`}
    >
      <span
        className={`pointer-events-none inline-block h-4 w-4 rounded-full bg-white shadow-sm transition-transform mt-0.5 ${
          enabled ? 'translate-x-4.5' : 'translate-x-0.5'
        }`}
      />
    </button>
  )
}

// ---------------------------------------------------------------------------
// Trigger detail view (with logs)
// ---------------------------------------------------------------------------

interface TriggerDetailProps {
  trigger: AgentTrigger
  onBack: () => void
  onDelete: (trigger: AgentTrigger) => void
}

function TriggerDetail({ trigger, onBack, onDelete }: TriggerDetailProps) {
  const { data: logsPage, isLoading: loadingLogs } = useTriggerLogs(trigger.id, { page: 1, page_size: 20 })
  const enableMutation = useEnableAgentTrigger()
  const disableMutation = useDisableAgentTrigger()

  const logs = logsPage?.items ?? []
  const isToggling = enableMutation.isPending || disableMutation.isPending

  function handleToggle(enabled: boolean) {
    if (enabled) {
      enableMutation.mutate(trigger.id)
    } else {
      disableMutation.mutate(trigger.id)
    }
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center gap-3">
        <button
          onClick={onBack}
          className="inline-flex items-center gap-1 text-sm text-slate-500 hover:text-slate-700 transition-colors"
        >
          <ChevronLeft size={16} />
          Back
        </button>
      </div>

      <div className="rounded-xl border border-slate-200 bg-white p-5">
        <div className="flex items-start justify-between gap-4">
          <div className="min-w-0 flex-1">
            <div className="flex items-center gap-3">
              <h2 className="text-lg font-semibold text-slate-800">{trigger.name}</h2>
              <span className="rounded-full bg-slate-100 px-2.5 py-0.5 text-xs font-mono text-slate-600">
                {trigger.event_type}
              </span>
            </div>
            <div className="mt-2 flex items-center gap-4 text-xs text-slate-500">
              <span>Cooldown: {trigger.cooldown}s</span>
              {trigger.last_fired_at && (
                <span>Last fired: {formatDate(trigger.last_fired_at)}</span>
              )}
              <span>Created: {formatDateShort(trigger.created_at)}</span>
            </div>
          </div>
          <div className="flex items-center gap-3 shrink-0">
            <div className="flex items-center gap-2">
              <span className="text-xs text-slate-500">{trigger.enabled ? 'Enabled' : 'Disabled'}</span>
              <Toggle enabled={trigger.enabled} onChange={handleToggle} disabled={isToggling} />
            </div>
            <button
              onClick={() => onDelete(trigger)}
              className="rounded-lg border border-red-200 px-3 py-1.5 text-xs font-medium text-red-600 hover:bg-red-50 transition-colors"
            >
              <Trash2 size={12} />
            </button>
          </div>
        </div>

        <div className="mt-4">
          <label className="block text-xs font-medium text-slate-500 mb-1 uppercase tracking-wider">Prompt</label>
          <pre className="rounded-lg bg-slate-50 border border-slate-200 p-3 text-xs text-slate-700 overflow-x-auto whitespace-pre-wrap font-mono">
            {trigger.prompt}
          </pre>
        </div>
      </div>

      {/* Execution logs */}
      <div>
        <h3 className="text-sm font-semibold text-slate-700 mb-3">Execution Logs</h3>
        <div className="rounded-xl border border-slate-200 bg-white overflow-hidden">
          {loadingLogs ? (
            <div className="p-8 text-center text-sm text-slate-500">Loading logs…</div>
          ) : logs.length === 0 ? (
            <div className="p-12 text-center">
              <Zap size={28} className="mx-auto mb-3 text-slate-300" />
              <p className="text-sm text-slate-500">No executions yet</p>
            </div>
          ) : (
            <table className="w-full">
              <thead>
                <tr className="border-b border-slate-200 bg-slate-50">
                  <th className="py-3 px-4 text-left text-xs font-semibold uppercase tracking-wider text-slate-500">
                    Status
                  </th>
                  <th className="py-3 px-4 text-left text-xs font-semibold uppercase tracking-wider text-slate-500">
                    Event Payload
                  </th>
                  <th className="py-3 px-4 text-left text-xs font-semibold uppercase tracking-wider text-slate-500">
                    Agent Response
                  </th>
                  <th className="py-3 px-4 text-left text-xs font-semibold uppercase tracking-wider text-slate-500">
                    Time
                  </th>
                </tr>
              </thead>
              <tbody>
                {logs.map((log) => (
                  <tr
                    key={log.id}
                    className="border-b border-slate-100 last:border-0 hover:bg-slate-50 transition-colors"
                  >
                    <td className="py-3 px-4">
                      <LogStatusBadge status={log.status} />
                    </td>
                    <td className="py-3 px-4">
                      <PayloadViewer data={log.event_payload} label="View payload" />
                    </td>
                    <td className="py-3 px-4">
                      <p className="text-xs text-slate-600 max-w-xs truncate" title={log.agent_response}>
                        {log.agent_response || <span className="text-slate-300 italic">—</span>}
                      </p>
                    </td>
                    <td className="py-3 px-4 text-xs text-slate-500 whitespace-nowrap">
                      {formatDate(log.created_at)}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      </div>
    </div>
  )
}

// ---------------------------------------------------------------------------
// Trigger row
// ---------------------------------------------------------------------------

interface TriggerRowProps {
  trigger: AgentTrigger
  onClick: () => void
  onDelete: (trigger: AgentTrigger) => void
}

function TriggerRow({ trigger, onClick, onDelete }: TriggerRowProps) {
  const enableMutation = useEnableAgentTrigger()
  const disableMutation = useDisableAgentTrigger()
  const isToggling = enableMutation.isPending || disableMutation.isPending

  function handleToggle(enabled: boolean) {
    if (enabled) {
      enableMutation.mutate(trigger.id)
    } else {
      disableMutation.mutate(trigger.id)
    }
  }

  return (
    <tr className="border-b border-slate-100 last:border-0 hover:bg-slate-50 transition-colors">
      <td
        className="py-3 px-4 cursor-pointer"
        onClick={onClick}
      >
        <p className="text-sm font-medium text-slate-800">{trigger.name}</p>
        <p className="text-xs text-slate-400 mt-0.5 font-mono">{trigger.event_type}</p>
      </td>
      <td className="py-3 px-4">
        <Toggle
          enabled={trigger.enabled}
          onChange={handleToggle}
          disabled={isToggling}
        />
      </td>
      <td className="py-3 px-4 text-sm text-slate-500">
        {trigger.cooldown > 0 ? `${trigger.cooldown}s` : <span className="text-slate-300">—</span>}
      </td>
      <td className="py-3 px-4 text-sm text-slate-500">
        {trigger.last_fired_at ? formatDateShort(trigger.last_fired_at) : <span className="text-slate-300">Never</span>}
      </td>
      <td className="py-3 px-4">
        <div className="flex items-center justify-end gap-2">
          <button
            onClick={onClick}
            className="rounded-lg border border-slate-200 px-2.5 py-1.5 text-xs font-medium text-slate-600 hover:bg-slate-50 transition-colors"
          >
            View logs
          </button>
          <button
            onClick={(e) => {
              e.stopPropagation()
              onDelete(trigger)
            }}
            className="rounded-lg border border-red-200 px-2.5 py-1.5 text-xs font-medium text-red-600 hover:bg-red-50 transition-colors"
          >
            <Trash2 size={12} />
          </button>
        </div>
      </td>
    </tr>
  )
}

// ---------------------------------------------------------------------------
// Main page
// ---------------------------------------------------------------------------

export default function AgentTriggers() {
  const [selectedTrigger, setSelectedTrigger] = useState<AgentTrigger | null>(null)
  const [showCreateDialog, setShowCreateDialog] = useState(false)
  const [deletingTrigger, setDeletingTrigger] = useState<AgentTrigger | null>(null)

  const { data: triggersPage, isLoading, error } = useAgentTriggers()
  const { data: eventTypesData } = useTriggerEventTypes()

  const triggers = triggersPage?.items ?? []
  const eventTypes = eventTypesData?.event_types ?? []

  // If a selected trigger was deleted or updated, sync from list
  const syncedSelectedTrigger = selectedTrigger
    ? (triggers.find((t) => t.id === selectedTrigger.id) ?? selectedTrigger)
    : null

  function handleDelete(trigger: AgentTrigger) {
    setDeletingTrigger(trigger)
    if (selectedTrigger?.id === trigger.id) {
      setSelectedTrigger(null)
    }
  }

  if (syncedSelectedTrigger) {
    return (
      <>
        <TriggerDetail
          trigger={syncedSelectedTrigger}
          onBack={() => setSelectedTrigger(null)}
          onDelete={handleDelete}
        />
        {deletingTrigger && (
          <DeleteTriggerDialog
            trigger={deletingTrigger}
            onClose={() => setDeletingTrigger(null)}
          />
        )}
      </>
    )
  }

  return (
    <div className="space-y-6">
      {/* Page header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <Zap size={24} className="text-indigo-600" />
          <div>
            <h1 className="text-xl font-semibold text-slate-800">Agent Triggers</h1>
            <p className="text-sm text-slate-500">Automate agent actions in response to store events</p>
          </div>
        </div>
        <button
          onClick={() => setShowCreateDialog(true)}
          className="inline-flex items-center gap-2 rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700 transition-colors"
        >
          <Plus size={16} />
          New Trigger
        </button>
      </div>

      {/* Table */}
      <div className="rounded-xl border border-slate-200 bg-white overflow-hidden">
        {isLoading ? (
          <div className="p-8 text-center text-sm text-slate-500">Loading…</div>
        ) : error ? (
          <div className="rounded-xl border border-red-200 bg-red-50 p-6 text-center text-sm text-red-600 m-4">
            Failed to load triggers. Please try again.
          </div>
        ) : triggers.length === 0 ? (
          <div className="p-12 text-center">
            <Zap size={32} className="mx-auto mb-3 text-slate-300" />
            <p className="text-sm text-slate-500">No triggers configured yet</p>
            <button
              onClick={() => setShowCreateDialog(true)}
              className="mt-3 text-xs font-medium text-indigo-600 hover:underline"
            >
              Create your first trigger →
            </button>
          </div>
        ) : (
          <table className="w-full">
            <thead>
              <tr className="border-b border-slate-200 bg-slate-50">
                <th className="py-3 px-4 text-left text-xs font-semibold uppercase tracking-wider text-slate-500">
                  Name / Event Type
                </th>
                <th className="py-3 px-4 text-left text-xs font-semibold uppercase tracking-wider text-slate-500">
                  Enabled
                </th>
                <th className="py-3 px-4 text-left text-xs font-semibold uppercase tracking-wider text-slate-500">
                  Cooldown
                </th>
                <th className="py-3 px-4 text-left text-xs font-semibold uppercase tracking-wider text-slate-500">
                  Last Fired
                </th>
                <th className="py-3 px-4 text-right text-xs font-semibold uppercase tracking-wider text-slate-500">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody>
              {triggers.map((trigger) => (
                <TriggerRow
                  key={trigger.id}
                  trigger={trigger}
                  onClick={() => setSelectedTrigger(trigger)}
                  onDelete={handleDelete}
                />
              ))}
            </tbody>
          </table>
        )}
      </div>

      {/* Create dialog */}
      {showCreateDialog && (
        <CreateTriggerDialog
          eventTypes={eventTypes}
          onClose={() => setShowCreateDialog(false)}
        />
      )}

      {/* Delete dialog */}
      {deletingTrigger && (
        <DeleteTriggerDialog
          trigger={deletingTrigger}
          onClose={() => setDeletingTrigger(null)}
        />
      )}
    </div>
  )
}
