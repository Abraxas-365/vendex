import { useState } from 'react'
import { Plus, DollarSign, X, RefreshCw } from 'lucide-react'
import type { CurrencyRate } from '../../types'
import {
  useCurrencyRates,
  useSupportedCurrencies,
  useSetCurrencyRate,
  useDeleteCurrencyRate,
  useConvertCurrency,
} from '../../lib/hooks'

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

// ---------------------------------------------------------------------------
// Set Rate Dialog
// ---------------------------------------------------------------------------

interface RateDialogProps {
  initial?: CurrencyRate
  currencies: string[]
  onClose: () => void
  onSave: (data: { base_currency: string; target_currency: string; rate: number; auto_update: boolean }) => void
  saving: boolean
  error?: string
}

function RateDialog({ initial, currencies, onClose, onSave, saving, error }: RateDialogProps) {
  const [base, setBase] = useState(initial?.base_currency ?? 'USD')
  const [target, setTarget] = useState(initial?.target_currency ?? 'EUR')
  const [rate, setRate] = useState(String(initial?.rate ?? ''))
  const [autoUpdate, setAutoUpdate] = useState(initial?.auto_update ?? false)

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    const parsedRate = parseFloat(rate)
    if (isNaN(parsedRate) || parsedRate <= 0) return
    onSave({ base_currency: base, target_currency: target, rate: parsedRate, auto_update: autoUpdate })
  }

  const currencyOptions = currencies.length > 0 ? currencies : ['USD', 'EUR', 'GBP', 'JPY', 'CAD', 'AUD']

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4">
      <div className="w-full max-w-md rounded-xl bg-white shadow-xl">
        <div className="flex items-center justify-between border-b border-slate-100 px-6 py-4">
          <h2 className="text-base font-semibold text-slate-800">
            {initial ? 'Edit Rate' : 'Set Currency Rate'}
          </h2>
          <button onClick={onClose} className="text-slate-400 hover:text-slate-600">
            <X size={18} />
          </button>
        </div>
        <form onSubmit={handleSubmit} className="px-6 py-4 space-y-4">
          <div className="flex gap-3">
            <div className="flex-1">
              <label className="mb-1 block text-xs font-medium text-slate-600">Base Currency</label>
              <select
                value={base}
                onChange={(e) => setBase(e.target.value)}
                className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm outline-none focus:border-indigo-400"
              >
                {currencyOptions.map((c) => (
                  <option key={c} value={c}>{c}</option>
                ))}
              </select>
            </div>
            <div className="flex-1">
              <label className="mb-1 block text-xs font-medium text-slate-600">Target Currency</label>
              <select
                value={target}
                onChange={(e) => setTarget(e.target.value)}
                className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm outline-none focus:border-indigo-400"
              >
                {currencyOptions.map((c) => (
                  <option key={c} value={c}>{c}</option>
                ))}
              </select>
            </div>
          </div>
          <div>
            <label className="mb-1 block text-xs font-medium text-slate-600">Rate (1 {base} = ? {target})</label>
            <input
              type="number"
              step="0.000001"
              min="0.000001"
              required
              value={rate}
              onChange={(e) => setRate(e.target.value)}
              placeholder="1.08"
              className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm outline-none focus:border-indigo-400 focus:ring-1 focus:ring-indigo-300"
            />
          </div>
          <label className="flex items-center gap-2 cursor-pointer">
            <input
              type="checkbox"
              checked={autoUpdate}
              onChange={(e) => setAutoUpdate(e.target.checked)}
              className="rounded border-slate-300 text-indigo-600"
            />
            <span className="text-sm text-slate-600">Auto-update rate</span>
          </label>
          {error && <p className="text-xs text-red-600">{error}</p>}
          <div className="flex justify-end gap-2 pt-2">
            <button
              type="button"
              onClick={onClose}
              className="rounded-lg border border-slate-200 px-4 py-2 text-sm font-medium text-slate-600 hover:bg-slate-50"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={saving}
              className="rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700 disabled:opacity-60"
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
// Currency Converter
// ---------------------------------------------------------------------------

function CurrencyConverter({ currencies }: { currencies: string[] }) {
  const [amount, setAmount] = useState('')
  const [fromCurrency, setFromCurrency] = useState('USD')
  const [toCurrency, setToCurrency] = useState('EUR')
  const convertCurrency = useConvertCurrency()

  const currencyOptions = currencies.length > 0 ? currencies : ['USD', 'EUR', 'GBP', 'JPY', 'CAD', 'AUD']

  function handleConvert(e: React.FormEvent) {
    e.preventDefault()
    const parsedAmount = Math.round(parseFloat(amount) * 100)
    if (isNaN(parsedAmount) || parsedAmount <= 0) return
    convertCurrency.mutate({ amount: parsedAmount, currency: fromCurrency, target_currency: toCurrency })
  }

  return (
    <div className="rounded-xl border border-slate-200 bg-white p-5">
      <h2 className="mb-4 text-sm font-semibold text-slate-700">Currency Converter</h2>
      <form onSubmit={handleConvert} className="flex flex-wrap gap-3 items-end">
        <div>
          <label className="mb-1 block text-xs font-medium text-slate-600">Amount</label>
          <input
            type="number"
            step="0.01"
            min="0.01"
            required
            value={amount}
            onChange={(e) => setAmount(e.target.value)}
            placeholder="100.00"
            className="w-32 rounded-lg border border-slate-200 px-3 py-2 text-sm outline-none focus:border-indigo-400"
          />
        </div>
        <div>
          <label className="mb-1 block text-xs font-medium text-slate-600">From</label>
          <select
            value={fromCurrency}
            onChange={(e) => setFromCurrency(e.target.value)}
            className="rounded-lg border border-slate-200 px-3 py-2 text-sm outline-none focus:border-indigo-400"
          >
            {currencyOptions.map((c) => (
              <option key={c} value={c}>{c}</option>
            ))}
          </select>
        </div>
        <div>
          <label className="mb-1 block text-xs font-medium text-slate-600">To</label>
          <select
            value={toCurrency}
            onChange={(e) => setToCurrency(e.target.value)}
            className="rounded-lg border border-slate-200 px-3 py-2 text-sm outline-none focus:border-indigo-400"
          >
            {currencyOptions.map((c) => (
              <option key={c} value={c}>{c}</option>
            ))}
          </select>
        </div>
        <button
          type="submit"
          disabled={convertCurrency.isPending}
          className="flex items-center gap-1.5 rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700 disabled:opacity-60"
        >
          <RefreshCw size={14} />
          Convert
        </button>
        {convertCurrency.data && (
          <div className="ml-2 rounded-lg bg-indigo-50 px-4 py-2">
            <p className="text-sm font-semibold text-indigo-700">
              {(convertCurrency.data.from_amount / 100).toFixed(2)} {convertCurrency.data.from_currency}
              {' → '}
              {(convertCurrency.data.to_amount / 100).toFixed(2)} {convertCurrency.data.to_currency}
            </p>
            <p className="text-xs text-indigo-400">Rate: {convertCurrency.data.rate}</p>
          </div>
        )}
        {convertCurrency.error && (
          <p className="ml-2 text-sm text-red-600">{convertCurrency.error.message}</p>
        )}
      </form>
    </div>
  )
}

// ---------------------------------------------------------------------------
// Main Page
// ---------------------------------------------------------------------------

export default function CurrencyRates() {
  const [showDialog, setShowDialog] = useState(false)
  const [editRate, setEditRate] = useState<CurrencyRate | null>(null)

  const { data: rates, isPending, error } = useCurrencyRates()
  const { data: currencies } = useSupportedCurrencies()
  const setRate = useSetCurrencyRate()
  const deleteRate = useDeleteCurrencyRate()

  const ratesList = rates ?? []
  const currencyList = currencies ?? []

  function handleSave(data: { base_currency: string; target_currency: string; rate: number; auto_update: boolean }) {
    setRate.mutate(data, {
      onSuccess: () => {
        setShowDialog(false)
        setEditRate(null)
      },
    })
  }

  function handleDelete(rate: CurrencyRate) {
    if (!confirm(`Delete rate ${rate.base_currency} → ${rate.target_currency}?`)) return
    deleteRate.mutate(rate.id)
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <DollarSign size={20} className="text-indigo-600" />
          <h1 className="text-lg font-semibold text-slate-800">Currency Rates</h1>
        </div>
        <button
          onClick={() => { setEditRate(null); setShowDialog(true) }}
          className="flex items-center gap-1.5 rounded-lg bg-indigo-600 px-3 py-2 text-sm font-medium text-white hover:bg-indigo-700"
        >
          <Plus size={15} />
          Set Rate
        </button>
      </div>

      {/* Converter */}
      <CurrencyConverter currencies={currencyList} />

      {/* Rates Table */}
      <div className="rounded-xl border border-slate-200 bg-white">
        <div className="border-b border-slate-100 px-5 py-3">
          <h2 className="text-sm font-semibold text-slate-700">Exchange Rates</h2>
        </div>

        {isPending ? (
          <p className="p-6 text-sm text-slate-500">Loading…</p>
        ) : error ? (
          <p className="p-6 text-sm text-red-600">{error.message}</p>
        ) : ratesList.length === 0 ? (
          <p className="p-6 text-sm text-slate-400">No rates configured yet.</p>
        ) : (
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-slate-100">
                <th className="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide text-slate-400">Base</th>
                <th className="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide text-slate-400">Target</th>
                <th className="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide text-slate-400">Rate</th>
                <th className="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide text-slate-400">Auto Update</th>
                <th className="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide text-slate-400">Updated At</th>
                <th className="px-4 py-3" />
              </tr>
            </thead>
            <tbody>
              {ratesList.map((rate) => (
                <tr key={rate.id} className="border-b border-slate-50 hover:bg-slate-50">
                  <td className="px-4 py-3 font-semibold text-slate-700">{rate.base_currency}</td>
                  <td className="px-4 py-3 text-slate-700">{rate.target_currency}</td>
                  <td className="px-4 py-3 font-mono text-slate-700">{rate.rate.toFixed(6)}</td>
                  <td className="px-4 py-3">
                    <span
                      className={`rounded-full px-2 py-0.5 text-xs font-medium ${
                        rate.auto_update ? 'bg-green-50 text-green-700' : 'bg-slate-100 text-slate-500'
                      }`}
                    >
                      {rate.auto_update ? 'Yes' : 'No'}
                    </span>
                  </td>
                  <td className="px-4 py-3 text-slate-400 text-xs">{formatDate(rate.updated_at)}</td>
                  <td className="px-4 py-3">
                    <div className="flex items-center gap-1.5">
                      <button
                        onClick={() => { setEditRate(rate); setShowDialog(true) }}
                        className="rounded px-2 py-1 text-xs text-indigo-600 hover:bg-indigo-50"
                      >
                        Edit
                      </button>
                      <button
                        onClick={() => handleDelete(rate)}
                        className="rounded p-1 text-slate-400 hover:bg-red-50 hover:text-red-600"
                      >
                        <X size={14} />
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
      {showDialog && (
        <RateDialog
          initial={editRate ?? undefined}
          currencies={currencyList}
          onClose={() => { setShowDialog(false); setEditRate(null) }}
          onSave={handleSave}
          saving={setRate.isPending}
          error={setRate.error?.message}
        />
      )}
    </div>
  )
}
