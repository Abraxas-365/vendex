import { useState } from 'react'
import { Plus, X, Receipt } from 'lucide-react'
import type { TaxRate } from '../../types'
import { useTaxRates, useCreateTaxRate, useUpdateTaxRate, useDeleteTaxRate } from '../../lib/hooks'

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

function formatRate(rate: number): string {
  return `${(rate * 100).toFixed(4).replace(/\.?0+$/, '')}%`
}

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  })
}

// ---------------------------------------------------------------------------
// Tax Rate Form
// ---------------------------------------------------------------------------

interface TaxFormData {
  name: string
  rate: string // Stored as percentage string, e.g. "8.25" for 8.25%
  country: string
  state: string
  city: string
  zip_code: string
  priority: string
  compound: boolean
  includes_shipping: boolean
  active: boolean
}

const emptyForm: TaxFormData = {
  name: '',
  rate: '',
  country: '',
  state: '',
  city: '',
  zip_code: '',
  priority: '0',
  compound: false,
  includes_shipping: false,
  active: true,
}

function rateToForm(rate: TaxRate): TaxFormData {
  return {
    name: rate.name,
    rate: (rate.rate * 100).toFixed(4).replace(/\.?0+$/, ''),
    country: rate.country,
    state: rate.state ?? '',
    city: rate.city ?? '',
    zip_code: rate.zip_code ?? '',
    priority: String(rate.priority),
    compound: rate.compound,
    includes_shipping: rate.includes_shipping,
    active: rate.active,
  }
}

function formToData(form: TaxFormData): Partial<TaxRate> {
  return {
    name: form.name,
    rate: parseFloat(form.rate) / 100,
    country: form.country,
    state: form.state || undefined,
    city: form.city || undefined,
    zip_code: form.zip_code || undefined,
    priority: parseInt(form.priority) || 0,
    compound: form.compound,
    includes_shipping: form.includes_shipping,
    active: form.active,
  }
}

interface TaxDialogProps {
  initial?: TaxRate
  onClose: () => void
  onSave: (data: Partial<TaxRate>) => void
  saving: boolean
  error?: string
}

function TaxDialog({ initial, onClose, onSave, saving, error }: TaxDialogProps) {
  const [form, setForm] = useState<TaxFormData>(initial ? rateToForm(initial) : emptyForm)

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    onSave(formToData(form))
  }

  function toggle(field: 'compound' | 'includes_shipping' | 'active') {
    return (
      <button
        type="button"
        onClick={() => setForm({ ...form, [field]: !form[field] })}
        className={`relative inline-flex h-5 w-9 items-center rounded-full transition-colors ${
          form[field] ? 'bg-indigo-600' : 'bg-slate-300'
        }`}
      >
        <span
          className={`inline-block h-3.5 w-3.5 transform rounded-full bg-white shadow transition-transform ${
            form[field] ? 'translate-x-4.5' : 'translate-x-0.5'
          }`}
        />
      </button>
    )
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4">
      <div className="w-full max-w-lg rounded-xl bg-white shadow-xl">
        <div className="flex items-center justify-between border-b border-slate-100 px-6 py-4">
          <h2 className="text-base font-semibold text-slate-800">
            {initial ? 'Edit Tax Rate' : 'New Tax Rate'}
          </h2>
          <button onClick={onClose} className="text-slate-400 hover:text-slate-600">
            <X size={18} />
          </button>
        </div>
        <form onSubmit={handleSubmit} className="space-y-4 overflow-y-auto px-6 py-4 max-h-[75vh]">
          <div className="grid grid-cols-2 gap-4">
            <div className="col-span-2">
              <label className="mb-1 block text-sm font-medium text-slate-700">Name</label>
              <input
                type="text"
                required
                value={form.name}
                onChange={(e) => setForm({ ...form, name: e.target.value })}
                className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm focus:border-indigo-400 focus:outline-none"
                placeholder="e.g. California State Tax"
              />
            </div>

            <div>
              <label className="mb-1 block text-sm font-medium text-slate-700">
                Rate <span className="font-normal text-slate-400">(%)</span>
              </label>
              <input
                type="number"
                required
                step="any"
                min="0"
                max="100"
                value={form.rate}
                onChange={(e) => setForm({ ...form, rate: e.target.value })}
                className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm focus:border-indigo-400 focus:outline-none"
                placeholder="8.25"
              />
            </div>

            <div>
              <label className="mb-1 block text-sm font-medium text-slate-700">Priority</label>
              <input
                type="number"
                min="0"
                value={form.priority}
                onChange={(e) => setForm({ ...form, priority: e.target.value })}
                className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm focus:border-indigo-400 focus:outline-none"
                placeholder="0"
              />
            </div>

            <div>
              <label className="mb-1 block text-sm font-medium text-slate-700">Country</label>
              <input
                type="text"
                required
                value={form.country}
                onChange={(e) => setForm({ ...form, country: e.target.value })}
                className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm focus:border-indigo-400 focus:outline-none"
                placeholder="US"
              />
            </div>

            <div>
              <label className="mb-1 block text-sm font-medium text-slate-700">
                State <span className="font-normal text-slate-400">(optional)</span>
              </label>
              <input
                type="text"
                value={form.state}
                onChange={(e) => setForm({ ...form, state: e.target.value })}
                className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm focus:border-indigo-400 focus:outline-none"
                placeholder="CA"
              />
            </div>

            <div>
              <label className="mb-1 block text-sm font-medium text-slate-700">
                City <span className="font-normal text-slate-400">(optional)</span>
              </label>
              <input
                type="text"
                value={form.city}
                onChange={(e) => setForm({ ...form, city: e.target.value })}
                className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm focus:border-indigo-400 focus:outline-none"
                placeholder="San Francisco"
              />
            </div>

            <div>
              <label className="mb-1 block text-sm font-medium text-slate-700">
                ZIP Code <span className="font-normal text-slate-400">(optional)</span>
              </label>
              <input
                type="text"
                value={form.zip_code}
                onChange={(e) => setForm({ ...form, zip_code: e.target.value })}
                className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm focus:border-indigo-400 focus:outline-none"
                placeholder="94102"
              />
            </div>
          </div>

          {/* Toggles */}
          <div className="space-y-3 border-t border-slate-100 pt-3">
            <div className="flex items-center gap-3">
              {toggle('compound')}
              <div>
                <span className="text-sm font-medium text-slate-700">Compound tax</span>
                <p className="text-xs text-slate-400">Applied on top of other taxes</p>
              </div>
            </div>
            <div className="flex items-center gap-3">
              {toggle('includes_shipping')}
              <div>
                <span className="text-sm font-medium text-slate-700">Includes shipping</span>
                <p className="text-xs text-slate-400">Tax applies to shipping costs</p>
              </div>
            </div>
            <div className="flex items-center gap-3">
              {toggle('active')}
              <span className="text-sm font-medium text-slate-700">
                {form.active ? 'Active' : 'Inactive'}
              </span>
            </div>
          </div>

          {error && <p className="text-sm text-red-600">{error}</p>}

          <div className="flex gap-2 pt-2">
            <button
              type="button"
              onClick={onClose}
              className="flex-1 rounded-lg border border-slate-200 px-4 py-2 text-sm font-medium text-slate-600 hover:bg-slate-50"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={saving}
              className="flex-1 rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700 disabled:opacity-60"
            >
              {saving ? 'Saving…' : 'Save Tax Rate'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}

// ---------------------------------------------------------------------------
// Main Tax page
// ---------------------------------------------------------------------------

export default function Tax() {
  const [showForm, setShowForm] = useState(false)
  const [editingRate, setEditingRate] = useState<TaxRate | null>(null)

  const { data: rates = [], isLoading } = useTaxRates()
  const createTaxRate = useCreateTaxRate()
  const updateTaxRate = useUpdateTaxRate()
  const deleteTaxRate = useDeleteTaxRate()

  function handleSave(data: Partial<TaxRate>) {
    if (editingRate) {
      updateTaxRate.mutate(
        { id: editingRate.id, data },
        {
          onSuccess: () => setEditingRate(null),
        },
      )
    } else {
      createTaxRate.mutate(data, {
        onSuccess: () => setShowForm(false),
      })
    }
  }

  function handleDelete(rate: TaxRate) {
    if (!confirm(`Delete tax rate "${rate.name}"?`)) return
    deleteTaxRate.mutate(rate.id)
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-xl font-semibold text-slate-800">Tax Rates</h1>
          <p className="mt-0.5 text-sm text-slate-500">
            Configure tax rates by country, state, city, or ZIP code.
          </p>
        </div>
        <button
          onClick={() => setShowForm(true)}
          className="flex items-center gap-2 rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700"
        >
          <Plus size={16} />
          New Tax Rate
        </button>
      </div>

      {/* Table */}
      <div className="rounded-xl border border-slate-200 bg-white overflow-hidden">
        {isLoading && (
          <div className="py-12 text-center text-sm text-slate-400">Loading…</div>
        )}

        {!isLoading && rates.length === 0 && (
          <div className="py-16 text-center">
            <Receipt size={40} className="mx-auto mb-3 text-slate-300" />
            <p className="text-sm font-medium text-slate-500">No tax rates configured</p>
            <p className="text-xs text-slate-400 mt-1">
              Add a tax rate to automatically apply taxes at checkout.
            </p>
          </div>
        )}

        {rates.length > 0 && (
          <table className="w-full text-sm">
            <thead className="bg-slate-50 text-xs font-semibold uppercase tracking-wide text-slate-500">
              <tr>
                <th className="px-5 py-3 text-left">Name</th>
                <th className="px-5 py-3 text-left">Rate</th>
                <th className="px-5 py-3 text-left">Country</th>
                <th className="px-5 py-3 text-left">State</th>
                <th className="px-5 py-3 text-center">Priority</th>
                <th className="px-5 py-3 text-center">Compound</th>
                <th className="px-5 py-3 text-center">Incl. Shipping</th>
                <th className="px-5 py-3 text-left">Status</th>
                <th className="px-5 py-3 text-left">Created</th>
                <th className="px-5 py-3" />
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100">
              {rates.map((rate) => (
                <tr key={rate.id} className="hover:bg-slate-50">
                  <td className="px-5 py-3 font-medium text-slate-800">{rate.name}</td>
                  <td className="px-5 py-3 text-slate-700 font-mono">{formatRate(rate.rate)}</td>
                  <td className="px-5 py-3 text-slate-600">{rate.country}</td>
                  <td className="px-5 py-3 text-slate-500">{rate.state ?? '—'}</td>
                  <td className="px-5 py-3 text-center text-slate-600">{rate.priority}</td>
                  <td className="px-5 py-3 text-center">
                    <span
                      className={`inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium ${
                        rate.compound
                          ? 'bg-indigo-50 text-indigo-700'
                          : 'bg-slate-100 text-slate-400'
                      }`}
                    >
                      {rate.compound ? 'Yes' : 'No'}
                    </span>
                  </td>
                  <td className="px-5 py-3 text-center">
                    <span
                      className={`inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium ${
                        rate.includes_shipping
                          ? 'bg-indigo-50 text-indigo-700'
                          : 'bg-slate-100 text-slate-400'
                      }`}
                    >
                      {rate.includes_shipping ? 'Yes' : 'No'}
                    </span>
                  </td>
                  <td className="px-5 py-3">
                    <span
                      className={`inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium ${
                        rate.active
                          ? 'bg-green-50 text-green-700'
                          : 'bg-slate-100 text-slate-500'
                      }`}
                    >
                      {rate.active ? 'Active' : 'Inactive'}
                    </span>
                  </td>
                  <td className="px-5 py-3 text-slate-400 text-xs">{formatDate(rate.created_at)}</td>
                  <td className="px-5 py-3">
                    <div className="flex items-center gap-3 justify-end">
                      <button
                        onClick={() => setEditingRate(rate)}
                        className="text-xs text-indigo-600 hover:underline"
                      >
                        Edit
                      </button>
                      <button
                        onClick={() => handleDelete(rate)}
                        className="text-xs text-red-500 hover:underline"
                      >
                        Delete
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>

      {/* Dialog */}
      {(showForm || editingRate) && (
        <TaxDialog
          initial={editingRate ?? undefined}
          onClose={() => {
            setShowForm(false)
            setEditingRate(null)
          }}
          onSave={handleSave}
          saving={createTaxRate.isPending || updateTaxRate.isPending}
          error={createTaxRate.error?.message ?? updateTaxRate.error?.message}
        />
      )}
    </div>
  )
}
