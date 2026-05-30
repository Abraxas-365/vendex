import { useState } from 'react'
import { Plus, Truck, X, ChevronRight, Globe, Trash2 } from 'lucide-react'
import type { ShippingZone, ShippingRate } from '../../types'
import {
  useShippingZones,
  useCreateShippingZone,
  useUpdateShippingZone,
  useDeleteShippingZone,
  useShippingRates,
  useCreateShippingRate,
  useUpdateShippingRate,
  useDeleteShippingRate,
} from '../../lib/hooks'

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

const rateTypeLabels: Record<ShippingRate['type'], string> = {
  flat: 'Flat Rate',
  weight_based: 'Weight Based',
  price_based: 'Price Based',
  free: 'Free',
}

function formatMoney(amount: number, currency: string): string {
  return new Intl.NumberFormat('en-US', { style: 'currency', currency }).format(amount / 100)
}

// ---------------------------------------------------------------------------
// Zone Form
// ---------------------------------------------------------------------------

interface ZoneFormData {
  name: string
  countries: string
  states: string
}

interface ZoneDialogProps {
  initial?: ShippingZone
  onClose: () => void
  onSave: (data: Partial<ShippingZone>) => void
  saving: boolean
  error?: string
}

function ZoneDialog({ initial, onClose, onSave, saving, error }: ZoneDialogProps) {
  const [form, setForm] = useState<ZoneFormData>({
    name: initial?.name ?? '',
    countries: initial?.countries.join(', ') ?? '',
    states: initial?.states.join(', ') ?? '',
  })

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    onSave({
      name: form.name,
      countries: form.countries
        .split(',')
        .map((c) => c.trim())
        .filter(Boolean),
      states: form.states
        .split(',')
        .map((s) => s.trim())
        .filter(Boolean),
    })
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4">
      <div className="w-full max-w-md rounded-xl bg-white shadow-xl">
        <div className="flex items-center justify-between border-b border-slate-100 px-6 py-4">
          <h2 className="text-base font-semibold text-slate-800">
            {initial ? 'Edit Zone' : 'New Shipping Zone'}
          </h2>
          <button onClick={onClose} className="text-slate-400 hover:text-slate-600">
            <X size={18} />
          </button>
        </div>
        <form onSubmit={handleSubmit} className="space-y-4 px-6 py-4">
          <div>
            <label className="mb-1 block text-sm font-medium text-slate-700">Zone Name</label>
            <input
              type="text"
              required
              value={form.name}
              onChange={(e) => setForm({ ...form, name: e.target.value })}
              className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm focus:border-indigo-400 focus:outline-none"
              placeholder="e.g. North America"
            />
          </div>
          <div>
            <label className="mb-1 block text-sm font-medium text-slate-700">
              Countries <span className="font-normal text-slate-400">(comma-separated ISO codes)</span>
            </label>
            <input
              type="text"
              value={form.countries}
              onChange={(e) => setForm({ ...form, countries: e.target.value })}
              className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm focus:border-indigo-400 focus:outline-none"
              placeholder="US, CA, MX"
            />
          </div>
          <div>
            <label className="mb-1 block text-sm font-medium text-slate-700">
              States / Provinces <span className="font-normal text-slate-400">(comma-separated, optional)</span>
            </label>
            <input
              type="text"
              value={form.states}
              onChange={(e) => setForm({ ...form, states: e.target.value })}
              className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm focus:border-indigo-400 focus:outline-none"
              placeholder="CA, NY, TX"
            />
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
              {saving ? 'Saving…' : 'Save Zone'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}

// ---------------------------------------------------------------------------
// Rate Form
// ---------------------------------------------------------------------------

interface RateFormData {
  name: string
  type: ShippingRate['type']
  price_amount: string
  min_weight: string
  max_weight: string
  min_order_amount: string
  max_order_amount: string
  est_days_min: string
  est_days_max: string
  active: boolean
}

const emptyRateForm: RateFormData = {
  name: '',
  type: 'flat',
  price_amount: '0',
  min_weight: '',
  max_weight: '',
  min_order_amount: '',
  max_order_amount: '',
  est_days_min: '',
  est_days_max: '',
  active: true,
}

function rateToForm(rate: ShippingRate): RateFormData {
  return {
    name: rate.name,
    type: rate.type,
    price_amount: (rate.price.amount / 100).toFixed(2),
    min_weight: rate.min_weight != null ? String(rate.min_weight) : '',
    max_weight: rate.max_weight != null ? String(rate.max_weight) : '',
    min_order_amount: rate.min_order_amount != null ? (rate.min_order_amount / 100).toFixed(2) : '',
    max_order_amount: rate.max_order_amount != null ? (rate.max_order_amount / 100).toFixed(2) : '',
    est_days_min: rate.est_days_min != null ? String(rate.est_days_min) : '',
    est_days_max: rate.est_days_max != null ? String(rate.est_days_max) : '',
    active: rate.active,
  }
}

interface RateDialogProps {
  zoneId: string
  initial?: ShippingRate
  onClose: () => void
  onSave: (data: Partial<ShippingRate>) => void
  saving: boolean
  error?: string
}

function RateDialog({ initial, onClose, onSave, saving, error }: RateDialogProps) {
  const [form, setForm] = useState<RateFormData>(initial ? rateToForm(initial) : emptyRateForm)

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    const priceAmountCents = Math.round(parseFloat(form.price_amount || '0') * 100)
    onSave({
      name: form.name,
      type: form.type,
      price: { amount: priceAmountCents, currency: 'USD' },
      min_weight: form.min_weight ? parseFloat(form.min_weight) : undefined,
      max_weight: form.max_weight ? parseFloat(form.max_weight) : undefined,
      min_order_amount: form.min_order_amount
        ? Math.round(parseFloat(form.min_order_amount) * 100)
        : undefined,
      max_order_amount: form.max_order_amount
        ? Math.round(parseFloat(form.max_order_amount) * 100)
        : undefined,
      est_days_min: form.est_days_min ? parseInt(form.est_days_min) : undefined,
      est_days_max: form.est_days_max ? parseInt(form.est_days_max) : undefined,
      active: form.active,
    })
  }

  function field(label: string, key: keyof RateFormData, type = 'text', placeholder = '') {
    const val = form[key]
    return (
      <div>
        <label className="mb-1 block text-sm font-medium text-slate-700">{label}</label>
        <input
          type={type}
          value={val as string}
          onChange={(e) => setForm({ ...form, [key]: e.target.value })}
          className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm focus:border-indigo-400 focus:outline-none"
          placeholder={placeholder}
          min={type === 'number' ? '0' : undefined}
          step={type === 'number' ? 'any' : undefined}
        />
      </div>
    )
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4">
      <div className="w-full max-w-lg rounded-xl bg-white shadow-xl">
        <div className="flex items-center justify-between border-b border-slate-100 px-6 py-4">
          <h2 className="text-base font-semibold text-slate-800">
            {initial ? 'Edit Rate' : 'New Shipping Rate'}
          </h2>
          <button onClick={onClose} className="text-slate-400 hover:text-slate-600">
            <X size={18} />
          </button>
        </div>
        <form onSubmit={handleSubmit} className="space-y-4 overflow-y-auto px-6 py-4 max-h-[70vh]">
          <div className="grid grid-cols-2 gap-4">
            <div className="col-span-2">
              <label className="mb-1 block text-sm font-medium text-slate-700">Rate Name</label>
              <input
                type="text"
                required
                value={form.name}
                onChange={(e) => setForm({ ...form, name: e.target.value })}
                className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm focus:border-indigo-400 focus:outline-none"
                placeholder="e.g. Standard Shipping"
              />
            </div>
            <div>
              <label className="mb-1 block text-sm font-medium text-slate-700">Type</label>
              <select
                value={form.type}
                onChange={(e) => setForm({ ...form, type: e.target.value as ShippingRate['type'] })}
                className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm focus:border-indigo-400 focus:outline-none"
              >
                {Object.entries(rateTypeLabels).map(([val, label]) => (
                  <option key={val} value={val}>
                    {label}
                  </option>
                ))}
              </select>
            </div>
            {field('Price ($)', 'price_amount', 'number', '0.00')}
            {field('Min Weight (g)', 'min_weight', 'number', 'Optional')}
            {field('Max Weight (g)', 'max_weight', 'number', 'Optional')}
            {field('Min Order Amount ($)', 'min_order_amount', 'number', 'Optional')}
            {field('Max Order Amount ($)', 'max_order_amount', 'number', 'Optional')}
            {field('Est. Days Min', 'est_days_min', 'number', 'Optional')}
            {field('Est. Days Max', 'est_days_max', 'number', 'Optional')}
          </div>
          <div className="flex items-center gap-3">
            <button
              type="button"
              onClick={() => setForm({ ...form, active: !form.active })}
              className={`relative inline-flex h-5 w-9 items-center rounded-full transition-colors ${
                form.active ? 'bg-indigo-600' : 'bg-slate-300'
              }`}
            >
              <span
                className={`inline-block h-3.5 w-3.5 transform rounded-full bg-white shadow transition-transform ${
                  form.active ? 'translate-x-4.5' : 'translate-x-0.5'
                }`}
              />
            </button>
            <span className="text-sm text-slate-700">{form.active ? 'Active' : 'Inactive'}</span>
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
              {saving ? 'Saving…' : 'Save Rate'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}

// ---------------------------------------------------------------------------
// Rates Panel
// ---------------------------------------------------------------------------

interface RatesPanelProps {
  zone: ShippingZone
}

function RatesPanel({ zone }: RatesPanelProps) {
  const [showForm, setShowForm] = useState(false)
  const [editingRate, setEditingRate] = useState<ShippingRate | null>(null)

  const { data: rates = [], isLoading } = useShippingRates(zone.id)
  const createRate = useCreateShippingRate()
  const updateRate = useUpdateShippingRate()
  const deleteRate = useDeleteShippingRate()

  function handleSaveRate(data: Partial<ShippingRate>) {
    if (editingRate) {
      updateRate.mutate(
        { id: editingRate.id, data },
        {
          onSuccess: () => {
            setEditingRate(null)
          },
        },
      )
    } else {
      createRate.mutate(
        { zoneId: zone.id, data },
        {
          onSuccess: () => {
            setShowForm(false)
          },
        },
      )
    }
  }

  return (
    <div className="flex flex-col gap-4">
      <div className="flex items-center justify-between">
        <div>
          <h3 className="text-sm font-semibold text-slate-800">Shipping Rates</h3>
          <p className="text-xs text-slate-500 mt-0.5">For zone: {zone.name}</p>
        </div>
        <button
          onClick={() => setShowForm(true)}
          className="flex items-center gap-1.5 rounded-lg bg-indigo-600 px-3 py-1.5 text-sm font-medium text-white hover:bg-indigo-700"
        >
          <Plus size={15} />
          Add Rate
        </button>
      </div>

      {isLoading && (
        <div className="py-8 text-center text-sm text-slate-400">Loading rates…</div>
      )}

      {!isLoading && rates.length === 0 && (
        <div className="rounded-lg border border-dashed border-slate-200 py-8 text-center text-sm text-slate-400">
          No shipping rates yet. Add one above.
        </div>
      )}

      {rates.length > 0 && (
        <div className="overflow-hidden rounded-lg border border-slate-200">
          <table className="w-full text-sm">
            <thead className="bg-slate-50 text-xs font-semibold uppercase tracking-wide text-slate-500">
              <tr>
                <th className="px-4 py-2.5 text-left">Name</th>
                <th className="px-4 py-2.5 text-left">Type</th>
                <th className="px-4 py-2.5 text-left">Price</th>
                <th className="px-4 py-2.5 text-left">Delivery</th>
                <th className="px-4 py-2.5 text-left">Status</th>
                <th className="px-4 py-2.5" />
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100">
              {rates.map((rate) => (
                <tr key={rate.id} className="hover:bg-slate-50">
                  <td className="px-4 py-3 font-medium text-slate-800">{rate.name}</td>
                  <td className="px-4 py-3 text-slate-600">{rateTypeLabels[rate.type]}</td>
                  <td className="px-4 py-3 text-slate-600">
                    {rate.type === 'free' ? 'Free' : formatMoney(rate.price.amount, rate.price.currency)}
                  </td>
                  <td className="px-4 py-3 text-slate-500">
                    {rate.est_days_min != null && rate.est_days_max != null
                      ? `${rate.est_days_min}–${rate.est_days_max} days`
                      : '—'}
                  </td>
                  <td className="px-4 py-3">
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
                  <td className="px-4 py-3">
                    <div className="flex items-center gap-2 justify-end">
                      <button
                        onClick={() => setEditingRate(rate)}
                        className="text-xs text-indigo-600 hover:underline"
                      >
                        Edit
                      </button>
                      <button
                        onClick={() =>
                          deleteRate.mutate({ id: rate.id, zoneId: zone.id })
                        }
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
        </div>
      )}

      {(showForm || editingRate) && (
        <RateDialog
          zoneId={zone.id}
          initial={editingRate ?? undefined}
          onClose={() => {
            setShowForm(false)
            setEditingRate(null)
          }}
          onSave={handleSaveRate}
          saving={createRate.isPending || updateRate.isPending}
          error={createRate.error?.message ?? updateRate.error?.message}
        />
      )}
    </div>
  )
}

// ---------------------------------------------------------------------------
// Main Shipping page
// ---------------------------------------------------------------------------

export default function Shipping() {
  const [showZoneForm, setShowZoneForm] = useState(false)
  const [editingZone, setEditingZone] = useState<ShippingZone | null>(null)
  const [selectedZone, setSelectedZone] = useState<ShippingZone | null>(null)

  const { data: zones = [], isLoading } = useShippingZones()
  const createZone = useCreateShippingZone()
  const updateZone = useUpdateShippingZone()
  const deleteZone = useDeleteShippingZone()

  function handleSaveZone(data: Partial<ShippingZone>) {
    if (editingZone) {
      updateZone.mutate(
        { id: editingZone.id, data },
        {
          onSuccess: (updated) => {
            setEditingZone(null)
            if (selectedZone?.id === updated.id) {
              setSelectedZone(updated)
            }
          },
        },
      )
    } else {
      createZone.mutate(data, {
        onSuccess: () => {
          setShowZoneForm(false)
        },
      })
    }
  }

  function handleDeleteZone(zone: ShippingZone) {
    if (!confirm(`Delete zone "${zone.name}"? All rates will be removed.`)) return
    deleteZone.mutate(zone.id, {
      onSuccess: () => {
        if (selectedZone?.id === zone.id) setSelectedZone(null)
      },
    })
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-xl font-semibold text-slate-800">Shipping</h1>
          <p className="mt-0.5 text-sm text-slate-500">
            Manage shipping zones and rates for your store.
          </p>
        </div>
        <button
          onClick={() => setShowZoneForm(true)}
          className="flex items-center gap-2 rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700"
        >
          <Plus size={16} />
          New Zone
        </button>
      </div>

      {/* Two-panel layout */}
      <div className="flex gap-6">
        {/* Left: Zones list */}
        <div className="w-72 shrink-0">
          <div className="rounded-xl border border-slate-200 bg-white overflow-hidden">
            <div className="border-b border-slate-100 px-4 py-3">
              <p className="text-xs font-semibold uppercase tracking-widest text-slate-400">
                Shipping Zones
              </p>
            </div>

            {isLoading && (
              <div className="py-8 text-center text-sm text-slate-400">Loading…</div>
            )}

            {!isLoading && zones.length === 0 && (
              <div className="py-8 text-center">
                <Globe size={32} className="mx-auto mb-2 text-slate-300" />
                <p className="text-sm text-slate-400">No zones yet</p>
              </div>
            )}

            {zones.length > 0 && (
              <ul className="divide-y divide-slate-100">
                {zones.map((zone) => (
                  <li key={zone.id}>
                    <button
                      onClick={() => setSelectedZone(zone)}
                      className={`flex w-full items-center gap-3 px-4 py-3 text-left transition-colors hover:bg-slate-50 ${
                        selectedZone?.id === zone.id
                          ? 'bg-indigo-50 text-indigo-700'
                          : 'text-slate-700'
                      }`}
                    >
                      <Truck size={16} className="shrink-0 text-slate-400" />
                      <div className="flex-1 min-w-0">
                        <p className="text-sm font-medium truncate">{zone.name}</p>
                        <p className="text-xs text-slate-400 truncate">
                          {zone.countries.length > 0
                            ? zone.countries.join(', ')
                            : 'No countries'}
                        </p>
                      </div>
                      <ChevronRight size={14} className="shrink-0 text-slate-300" />
                    </button>
                  </li>
                ))}
              </ul>
            )}
          </div>

          {/* Zone actions */}
          {selectedZone && (
            <div className="mt-3 flex gap-2">
              <button
                onClick={() => setEditingZone(selectedZone)}
                className="flex-1 rounded-lg border border-slate-200 px-3 py-1.5 text-xs font-medium text-slate-600 hover:bg-slate-50"
              >
                Edit Zone
              </button>
              <button
                onClick={() => handleDeleteZone(selectedZone)}
                className="flex items-center gap-1 rounded-lg border border-red-200 px-3 py-1.5 text-xs font-medium text-red-600 hover:bg-red-50"
              >
                <Trash2 size={13} />
                Delete
              </button>
            </div>
          )}
        </div>

        {/* Right: Rates panel */}
        <div className="flex-1">
          {selectedZone ? (
            <div className="rounded-xl border border-slate-200 bg-white p-6">
              <RatesPanel zone={selectedZone} />
            </div>
          ) : (
            <div className="flex h-full min-h-48 items-center justify-center rounded-xl border border-dashed border-slate-200 bg-white text-center">
              <div>
                <Truck size={32} className="mx-auto mb-2 text-slate-300" />
                <p className="text-sm text-slate-400">Select a shipping zone to manage its rates</p>
              </div>
            </div>
          )}
        </div>
      </div>

      {/* Zone dialog */}
      {(showZoneForm || editingZone) && (
        <ZoneDialog
          initial={editingZone ?? undefined}
          onClose={() => {
            setShowZoneForm(false)
            setEditingZone(null)
          }}
          onSave={handleSaveZone}
          saving={createZone.isPending || updateZone.isPending}
          error={createZone.error?.message ?? updateZone.error?.message}
        />
      )}
    </div>
  )
}
