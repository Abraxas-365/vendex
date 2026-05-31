import { useState } from 'react'
import { FlaskConical, Plus, Trash2, Edit2, Play, Pause, CheckCircle } from 'lucide-react'
import type { Experiment } from '../../types'
import {
  useExperiments,
  useCreateExperiment,
  useUpdateExperiment,
  useDeleteExperiment,
  useStartExperiment,
  usePauseExperiment,
  useCompleteExperiment,
  useExperimentResults,
} from '../../lib/hooks'

const emptyForm: Partial<Experiment> = {
  name: '',
  description: '',
  type: 'price',
  traffic_percentage: 50,
  variants: [],
}

const statusColors: Record<string, string> = {
  draft: 'bg-yellow-100 text-yellow-700',
  running: 'bg-green-100 text-green-700',
  paused: 'bg-orange-100 text-orange-700',
  completed: 'bg-blue-100 text-blue-700',
}

function ResultsPanel({ experimentId }: { experimentId: string }) {
  const { data, isLoading } = useExperimentResults(experimentId)

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-6">
        <div className="h-5 w-5 animate-spin rounded-full border-2 border-gray-300 border-t-indigo-600" />
      </div>
    )
  }

  if (!data) return <p className="text-sm text-gray-400 py-4">No results available.</p>

  return (
    <div className="space-y-3">
      <p className="text-sm text-gray-500">Total Impressions: <span className="font-medium text-gray-900">{data.total_impressions.toLocaleString()}</span></p>
      <div className="overflow-x-auto">
        <table className="w-full text-sm">
          <thead className="border-b border-gray-100 bg-gray-50">
            <tr>
              <th className="px-4 py-2 text-left font-medium text-gray-500">Variant</th>
              <th className="px-4 py-2 text-right font-medium text-gray-500">Impressions</th>
              <th className="px-4 py-2 text-right font-medium text-gray-500">Conversions</th>
              <th className="px-4 py-2 text-right font-medium text-gray-500">Conv. Rate</th>
              <th className="px-4 py-2 text-right font-medium text-gray-500">Revenue/Visitor</th>
              <th className="px-4 py-2 text-left font-medium text-gray-500">Winner</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-50">
            {(data.variants ?? []).map((v) => (
              <tr key={v.variant_id} className={v.is_winner ? 'bg-green-50/40' : ''}>
                <td className="px-4 py-2 font-medium text-gray-900">{v.name}</td>
                <td className="px-4 py-2 text-right text-gray-600">{v.impressions.toLocaleString()}</td>
                <td className="px-4 py-2 text-right text-gray-600">{v.conversions.toLocaleString()}</td>
                <td className="px-4 py-2 text-right text-gray-600">{(v.conversion_rate * 100).toFixed(2)}%</td>
                <td className="px-4 py-2 text-right text-gray-600">${v.revenue_per_visitor.toFixed(2)}</td>
                <td className="px-4 py-2">
                  {v.is_winner && (
                    <span className="inline-flex rounded-full bg-green-100 px-2 py-0.5 text-xs font-medium text-green-700">Winner</span>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}

export default function ABTesting() {
  const [page] = useState(1)
  const [showDialog, setShowDialog] = useState(false)
  const [editing, setEditing] = useState<Experiment | null>(null)
  const [form, setForm] = useState<Partial<Experiment>>(emptyForm)
  const [selectedId, setSelectedId] = useState<string | null>(null)

  const { data, isLoading } = useExperiments({ page, page_size: 20 })
  const experiments = data?.items ?? []

  const create = useCreateExperiment()
  const update = useUpdateExperiment()
  const remove = useDeleteExperiment()
  const start = useStartExperiment()
  const pause = usePauseExperiment()
  const complete = useCompleteExperiment()

  function openCreate() {
    setEditing(null)
    setForm(emptyForm)
    setShowDialog(true)
  }

  function openEdit(exp: Experiment) {
    setEditing(exp)
    setForm({
      name: exp.name,
      description: exp.description,
      type: exp.type,
      traffic_percentage: exp.traffic_percentage,
    })
    setShowDialog(true)
  }

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (editing) {
      update.mutate(
        { id: editing.id, data: form },
        { onSuccess: () => setShowDialog(false) },
      )
    } else {
      create.mutate(form, { onSuccess: () => setShowDialog(false) })
    }
  }

  const isPending = create.isPending || update.isPending

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <FlaskConical className="h-6 w-6 text-indigo-600" />
          <div>
            <h1 className="text-2xl font-bold text-gray-900">A/B Testing</h1>
            <p className="mt-0.5 text-sm text-gray-500">Run experiments to optimise conversion</p>
          </div>
        </div>
        <button
          onClick={openCreate}
          className="inline-flex items-center gap-2 rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700 transition-colors"
        >
          <Plus className="h-4 w-4" />
          New Experiment
        </button>
      </div>

      {/* Table */}
      <div className="rounded-xl border border-gray-200 bg-white shadow-sm overflow-hidden">
        {isLoading ? (
          <div className="flex items-center justify-center py-12">
            <div className="h-6 w-6 animate-spin rounded-full border-2 border-gray-300 border-t-indigo-600" />
          </div>
        ) : experiments.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-16 text-gray-400">
            <FlaskConical className="mb-3 h-10 w-10" />
            <p className="text-sm font-medium">No experiments yet</p>
          </div>
        ) : (
          <table className="w-full text-sm">
            <thead className="border-b border-gray-100 bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left font-medium text-gray-500">Name</th>
                <th className="px-6 py-3 text-left font-medium text-gray-500">Type</th>
                <th className="px-6 py-3 text-left font-medium text-gray-500">Status</th>
                <th className="px-6 py-3 text-left font-medium text-gray-500">Traffic %</th>
                <th className="px-6 py-3 text-left font-medium text-gray-500">Variants</th>
                <th className="px-6 py-3 text-right font-medium text-gray-500">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-50">
              {experiments.map((exp) => (
                <>
                  <tr key={exp.id} className="hover:bg-gray-50/50 transition-colors">
                    <td className="px-6 py-4">
                      <button
                        className="text-left font-medium text-indigo-700 hover:underline"
                        onClick={() => setSelectedId(selectedId === exp.id ? null : exp.id)}
                      >
                        {exp.name}
                      </button>
                      {exp.description && (
                        <p className="text-xs text-gray-400 mt-0.5 truncate max-w-xs">{exp.description}</p>
                      )}
                    </td>
                    <td className="px-6 py-4 text-gray-600 capitalize">{exp.type}</td>
                    <td className="px-6 py-4">
                      <span className={`inline-flex rounded-full px-2 py-0.5 text-xs font-medium ${statusColors[exp.status] ?? 'bg-gray-100 text-gray-600'}`}>
                        {exp.status}
                      </span>
                    </td>
                    <td className="px-6 py-4 text-gray-600">{exp.traffic_percentage}%</td>
                    <td className="px-6 py-4 text-gray-600">{(exp.variants ?? []).length}</td>
                    <td className="px-6 py-4">
                      <div className="flex items-center justify-end gap-1">
                        {exp.status === 'draft' && (
                          <button
                            onClick={() => start.mutate(exp.id)}
                            className="rounded-md p-1.5 text-gray-400 hover:bg-green-50 hover:text-green-600 transition-colors"
                            title="Start"
                          >
                            <Play className="h-4 w-4" />
                          </button>
                        )}
                        {exp.status === 'running' && (
                          <>
                            <button
                              onClick={() => pause.mutate(exp.id)}
                              className="rounded-md p-1.5 text-gray-400 hover:bg-orange-50 hover:text-orange-600 transition-colors"
                              title="Pause"
                            >
                              <Pause className="h-4 w-4" />
                            </button>
                            <button
                              onClick={() => complete.mutate(exp.id)}
                              className="rounded-md p-1.5 text-gray-400 hover:bg-blue-50 hover:text-blue-600 transition-colors"
                              title="Complete"
                            >
                              <CheckCircle className="h-4 w-4" />
                            </button>
                          </>
                        )}
                        {exp.status === 'paused' && (
                          <button
                            onClick={() => start.mutate(exp.id)}
                            className="rounded-md p-1.5 text-gray-400 hover:bg-green-50 hover:text-green-600 transition-colors"
                            title="Resume"
                          >
                            <Play className="h-4 w-4" />
                          </button>
                        )}
                        <button
                          onClick={() => openEdit(exp)}
                          className="rounded-md p-1.5 text-gray-400 hover:bg-blue-50 hover:text-blue-600 transition-colors"
                          title="Edit"
                        >
                          <Edit2 className="h-4 w-4" />
                        </button>
                        <button
                          onClick={() => remove.mutate(exp.id)}
                          className="rounded-md p-1.5 text-gray-400 hover:bg-red-50 hover:text-red-600 transition-colors"
                          title="Delete"
                        >
                          <Trash2 className="h-4 w-4" />
                        </button>
                      </div>
                    </td>
                  </tr>
                  {selectedId === exp.id && (
                    <tr key={`${exp.id}-results`}>
                      <td colSpan={6} className="px-6 py-4 bg-gray-50/50">
                        <p className="text-xs font-semibold uppercase tracking-widest text-gray-400 mb-3">Results</p>
                        <ResultsPanel experimentId={exp.id} />
                      </td>
                    </tr>
                  )}
                </>
              ))}
            </tbody>
          </table>
        )}
      </div>

      {/* Dialog */}
      {showDialog && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4">
          <div className="w-full max-w-lg rounded-2xl bg-white shadow-2xl">
            <div className="border-b border-gray-100 px-6 py-4">
              <h2 className="text-lg font-semibold text-gray-900">
                {editing ? 'Edit Experiment' : 'New Experiment'}
              </h2>
            </div>
            <form onSubmit={handleSubmit} className="space-y-4 px-6 py-5">
              <div>
                <label className="block text-xs font-medium text-gray-700 mb-1">Name *</label>
                <input
                  required
                  value={form.name ?? ''}
                  onChange={(e) => setForm((f) => ({ ...f, name: e.target.value }))}
                  className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
                />
              </div>
              <div>
                <label className="block text-xs font-medium text-gray-700 mb-1">Description</label>
                <textarea
                  rows={2}
                  value={form.description ?? ''}
                  onChange={(e) => setForm((f) => ({ ...f, description: e.target.value }))}
                  className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
                />
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-xs font-medium text-gray-700 mb-1">Type</label>
                  <input
                    value={form.type ?? ''}
                    onChange={(e) => setForm((f) => ({ ...f, type: e.target.value }))}
                    placeholder="price, layout, copy..."
                    className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
                  />
                </div>
                <div>
                  <label className="block text-xs font-medium text-gray-700 mb-1">Traffic % *</label>
                  <input
                    required
                    type="number"
                    min={1}
                    max={100}
                    value={form.traffic_percentage ?? 50}
                    onChange={(e) => setForm((f) => ({ ...f, traffic_percentage: parseInt(e.target.value) || 50 }))}
                    className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
                  />
                </div>
              </div>
              {(create.error ?? update.error) && (
                <p className="text-sm text-red-600">{(create.error ?? update.error)?.message}</p>
              )}
              <div className="flex justify-end gap-3 pt-2">
                <button
                  type="button"
                  onClick={() => setShowDialog(false)}
                  className="rounded-lg border border-gray-200 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={isPending}
                  className="rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700 disabled:opacity-50 transition-colors"
                >
                  {isPending ? 'Saving…' : editing ? 'Save Changes' : 'Create'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}
