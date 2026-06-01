import { useState } from 'react'
import { Plus, Gift, X, ChevronRight, Slash } from 'lucide-react'
import type { GiftCard, GiftCardTransaction } from '../../types'
import {
  useGiftCards,
  useGiftCardTransactions,
  useCreateGiftCard,
  useDisableGiftCard,
  useDeleteGiftCard,
} from '../../lib/hooks'

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

function formatMoney(amount: number, currency: string): string {
  return new Intl.NumberFormat('en-US', { style: 'currency', currency: currency || 'USD' }).format(amount / 100)
}

function formatDate(dateStr: string | null): string {
  if (!dateStr) return '—'
  return new Date(dateStr).toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' })
}

// ---------------------------------------------------------------------------
// Create Dialog
// ---------------------------------------------------------------------------

interface CreateGiftCardDialogProps {
  onClose: () => void
  onSave: (data: { code?: string; initial_amount: number; currency: string; expires_at?: string }) => void
  saving: boolean
  error?: string
}

function CreateGiftCardDialog({ onClose, onSave, saving, error }: CreateGiftCardDialogProps) {
  const [code, setCode] = useState('')
  const [amount, setAmount] = useState('')
  const [currency, setCurrency] = useState('USD')
  const [expiresAt, setExpiresAt] = useState('')

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    const parsedAmount = Math.round(parseFloat(amount) * 100)
    if (isNaN(parsedAmount) || parsedAmount <= 0) return
    onSave({
      ...(code.trim() ? { code: code.trim() } : {}),
      initial_amount: parsedAmount,
      currency,
      ...(expiresAt ? { expires_at: expiresAt } : {}),
    })
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4">
      <div className="w-full max-w-md rounded-xl bg-white shadow-xl">
        <div className="flex items-center justify-between border-b border-slate-100 px-6 py-4">
          <h2 className="text-base font-semibold text-slate-800">New Gift Card</h2>
          <button onClick={onClose} className="text-slate-400 hover:text-slate-600">
            <X size={18} />
          </button>
        </div>
        <form onSubmit={handleSubmit} className="px-6 py-4 space-y-4">
          <div>
            <label className="mb-1 block text-xs font-medium text-slate-600">
              Code <span className="text-slate-400">(leave empty to auto-generate)</span>
            </label>
            <input
              type="text"
              value={code}
              onChange={(e) => setCode(e.target.value)}
              placeholder="e.g. SUMMER2024"
              className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm outline-none focus:border-indigo-400 focus:ring-1 focus:ring-indigo-300"
            />
          </div>
          <div className="flex gap-3">
            <div className="flex-1">
              <label className="mb-1 block text-xs font-medium text-slate-600">Initial Amount ($)</label>
              <input
                type="number"
                min="0.01"
                step="0.01"
                required
                value={amount}
                onChange={(e) => setAmount(e.target.value)}
                placeholder="50.00"
                className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm outline-none focus:border-indigo-400 focus:ring-1 focus:ring-indigo-300"
              />
            </div>
            <div>
              <label className="mb-1 block text-xs font-medium text-slate-600">Currency</label>
              <input
                type="text"
                maxLength={3}
                value={currency}
                onChange={(e) => setCurrency(e.target.value.toUpperCase())}
                className="w-20 rounded-lg border border-slate-200 px-3 py-2 text-sm uppercase outline-none focus:border-indigo-400 focus:ring-1 focus:ring-indigo-300"
              />
            </div>
          </div>
          <div>
            <label className="mb-1 block text-xs font-medium text-slate-600">
              Expires At <span className="text-slate-400">(optional)</span>
            </label>
            <input
              type="date"
              value={expiresAt}
              onChange={(e) => setExpiresAt(e.target.value)}
              className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm outline-none focus:border-indigo-400 focus:ring-1 focus:ring-indigo-300"
            />
          </div>
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
              {saving ? 'Creating…' : 'Create Gift Card'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}

// ---------------------------------------------------------------------------
// Transaction History Panel
// ---------------------------------------------------------------------------

interface TransactionPanelProps {
  card: GiftCard
  onClose: () => void
}

function TransactionPanel({ card, onClose }: TransactionPanelProps) {
  const { data: transactions, isPending } = useGiftCardTransactions(card.id)

  return (
    <div className="flex flex-col rounded-xl border border-slate-200 bg-white">
      <div className="flex items-center justify-between border-b border-slate-100 px-5 py-3">
        <div>
          <p className="text-sm font-semibold text-slate-800">{card.code}</p>
          <p className="text-xs text-slate-400">
            Balance: {formatMoney(card.balance, card.currency)} / {formatMoney(card.initial_amount, card.currency)}
          </p>
        </div>
        <button onClick={onClose} className="text-slate-400 hover:text-slate-600">
          <X size={18} />
        </button>
      </div>
      <div className="flex-1 overflow-auto">
        {isPending ? (
          <p className="p-4 text-sm text-slate-500">Loading transactions…</p>
        ) : (transactions ?? []).length === 0 ? (
          <p className="p-4 text-sm text-slate-400">No transactions yet.</p>
        ) : (
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-slate-100">
                <th className="px-4 py-2 text-left text-xs font-medium uppercase text-slate-400">Type</th>
                <th className="px-4 py-2 text-left text-xs font-medium uppercase text-slate-400">Amount</th>
                <th className="px-4 py-2 text-left text-xs font-medium uppercase text-slate-400">Note</th>
                <th className="px-4 py-2 text-left text-xs font-medium uppercase text-slate-400">Date</th>
              </tr>
            </thead>
            <tbody>
              {(transactions as GiftCardTransaction[]).map((tx) => (
                <tr key={tx.id} className="border-b border-slate-50 hover:bg-slate-50">
                  <td className="px-4 py-2">
                    <span
                      className={`rounded-full px-2 py-0.5 text-xs font-medium ${
                        tx.type === 'credit'
                          ? 'bg-green-50 text-green-700'
                          : 'bg-red-50 text-red-700'
                      }`}
                    >
                      {tx.type}
                    </span>
                  </td>
                  <td className="px-4 py-2 font-medium">{formatMoney(tx.amount, tx.currency)}</td>
                  <td className="px-4 py-2 text-slate-500">{tx.note || '—'}</td>
                  <td className="px-4 py-2 text-slate-400">{formatDate(tx.created_at)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  )
}

// ---------------------------------------------------------------------------
// Main Page
// ---------------------------------------------------------------------------

export default function GiftCards() {
  const [page, setPage] = useState(1)
  const [showCreate, setShowCreate] = useState(false)
  const [selectedCard, setSelectedCard] = useState<GiftCard | null>(null)

  const { data, isPending, error } = useGiftCards({ page, page_size: 20 })
  const createGiftCard = useCreateGiftCard()
  const disableGiftCard = useDisableGiftCard()
  const deleteGiftCard = useDeleteGiftCard()

  const cards = data?.items ?? []

  function handleCreate(input: { code?: string; initial_amount: number; currency: string; expires_at?: string }) {
    createGiftCard.mutate(input, {
      onSuccess: () => setShowCreate(false),
    })
  }

  function handleDisable(card: GiftCard, e: React.MouseEvent) {
    e.stopPropagation()
    if (!confirm(`Disable gift card ${card.code}?`)) return
    disableGiftCard.mutate(card.id)
  }

  function handleDelete(card: GiftCard, e: React.MouseEvent) {
    e.stopPropagation()
    if (!confirm(`Delete gift card ${card.code}? This cannot be undone.`)) return
    deleteGiftCard.mutate(card.id, {
      onSuccess: () => {
        if (selectedCard?.id === card.id) setSelectedCard(null)
      },
    })
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <Gift size={20} className="text-indigo-600" />
          <h1 className="text-lg font-semibold text-slate-800">Gift Cards</h1>
        </div>
        <button
          onClick={() => setShowCreate(true)}
          className="flex items-center gap-1.5 rounded-lg bg-indigo-600 px-3 py-2 text-sm font-medium text-white hover:bg-indigo-700"
        >
          <Plus size={15} />
          New Gift Card
        </button>
      </div>

      <div className={`grid gap-6 ${selectedCard ? 'grid-cols-[1fr_420px]' : 'grid-cols-1'}`}>
        {/* Table */}
        <div className="rounded-xl border border-slate-200 bg-white">
          {isPending ? (
            <p className="p-6 text-sm text-slate-500">Loading…</p>
          ) : error ? (
            <p className="p-6 text-sm text-red-600">{error.message}</p>
          ) : cards.length === 0 ? (
            <p className="p-6 text-sm text-slate-400">No gift cards yet.</p>
          ) : (
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-slate-100">
                  <th className="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide text-slate-400">Code</th>
                  <th className="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide text-slate-400">Initial</th>
                  <th className="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide text-slate-400">Balance</th>
                  <th className="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide text-slate-400">Status</th>
                  <th className="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide text-slate-400">Expires</th>
                  <th className="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide text-slate-400">Created</th>
                  <th className="px-4 py-3" />
                </tr>
              </thead>
              <tbody>
                {cards.map((card) => (
                  <tr
                    key={card.id}
                    onClick={() => setSelectedCard(selectedCard?.id === card.id ? null : card)}
                    className={`cursor-pointer border-b border-slate-50 hover:bg-slate-50 ${
                      selectedCard?.id === card.id ? 'bg-indigo-50' : ''
                    }`}
                  >
                    <td className="px-4 py-3 font-mono text-xs font-semibold text-slate-700">{card.code}</td>
                    <td className="px-4 py-3 text-slate-700">{formatMoney(card.initial_amount, card.currency)}</td>
                    <td className="px-4 py-3 text-slate-700">{formatMoney(card.balance, card.currency)}</td>
                    <td className="px-4 py-3">
                      <span
                        className={`rounded-full px-2 py-0.5 text-xs font-medium ${
                          card.active ? 'bg-green-50 text-green-700' : 'bg-slate-100 text-slate-500'
                        }`}
                      >
                        {card.active ? 'Active' : 'Inactive'}
                      </span>
                    </td>
                    <td className="px-4 py-3 text-slate-500">{formatDate(card.expires_at)}</td>
                    <td className="px-4 py-3 text-slate-400">{formatDate(card.created_at)}</td>
                    <td className="px-4 py-3">
                      <div className="flex items-center gap-1.5">
                        {card.active && (
                          <button
                            onClick={(e) => handleDisable(card, e)}
                            title="Disable"
                            className="rounded p-1 text-slate-400 hover:bg-amber-50 hover:text-amber-600"
                          >
                            <Slash size={14} />
                          </button>
                        )}
                        <button
                          onClick={(e) => handleDelete(card, e)}
                          title="Delete"
                          className="rounded p-1 text-slate-400 hover:bg-red-50 hover:text-red-600"
                        >
                          <X size={14} />
                        </button>
                        <ChevronRight size={14} className="text-slate-300" />
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}

          {/* Pagination */}
          {data && data.total_pages > 1 && (
            <div className="flex items-center justify-between border-t border-slate-100 px-4 py-3">
              <p className="text-xs text-slate-400">
                Page {data.page} of {data.total_pages} · {data.total} cards
              </p>
              <div className="flex gap-2">
                <button
                  disabled={page <= 1}
                  onClick={() => setPage((p) => p - 1)}
                  className="rounded border border-slate-200 px-3 py-1 text-xs disabled:opacity-40"
                >
                  Prev
                </button>
                <button
                  disabled={page >= data.total_pages}
                  onClick={() => setPage((p) => p + 1)}
                  className="rounded border border-slate-200 px-3 py-1 text-xs disabled:opacity-40"
                >
                  Next
                </button>
              </div>
            </div>
          )}
        </div>

        {/* Transaction Panel */}
        {selectedCard && (
          <TransactionPanel card={selectedCard} onClose={() => setSelectedCard(null)} />
        )}
      </div>

      {/* Create Dialog */}
      {showCreate && (
        <CreateGiftCardDialog
          onClose={() => setShowCreate(false)}
          onSave={handleCreate}
          saving={createGiftCard.isPending}
          error={createGiftCard.error?.message}
        />
      )}
    </div>
  )
}
