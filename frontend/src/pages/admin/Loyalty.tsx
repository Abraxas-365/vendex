import { useState } from 'react'
import { Plus, Pencil, X, Award, ChevronDown, ChevronUp } from 'lucide-react'
import type { LoyaltyReward, LoyaltyAccount, RewardType } from '../../types'
import {
  useLoyaltyRewards,
  useCreateLoyaltyReward,
  useUpdateLoyaltyReward,
  useLoyaltyAccounts,
  useAdjustLoyaltyPoints,
  useLoyaltyTransactions,
} from '../../lib/hooks'

interface RewardFormData {
  name: string
  description: string
  type: RewardType
  points_cost: number
  value: number
  active: boolean
}

const emptyRewardForm: RewardFormData = {
  name: '',
  description: '',
  type: 'fixed_discount',
  points_cost: 0,
  value: 0,
  active: true,
}

function TransactionList({ accountId }: { accountId: string }) {
  const { data: transactions = [], isLoading } = useLoyaltyTransactions(accountId)
  if (isLoading) return <div className="flex justify-center py-3"><div className="h-4 w-4 animate-spin rounded-full border-2 border-gray-300 border-t-gray-900" /></div>
  if (transactions.length === 0) return <p className="py-2 text-sm text-gray-400 text-center">No transactions</p>
  return (
    <div className="space-y-1 max-h-40 overflow-y-auto">
      {transactions.map((t) => (
        <div key={t.id} className="flex items-center justify-between rounded-md px-3 py-1.5 text-xs hover:bg-gray-50">
          <span className="font-mono text-gray-700">{t.type}</span>
          <span className={t.points >= 0 ? 'text-green-700 font-semibold' : 'text-red-700 font-semibold'}>
            {t.points >= 0 ? '+' : ''}{t.points} pts
          </span>
          <span className="text-gray-400">{t.note}</span>
          <span className="text-gray-400">{new Date(t.created_at).toLocaleDateString()}</span>
        </div>
      ))}
    </div>
  )
}

const tierColors = {
  bronze: 'bg-orange-100 text-orange-800',
  silver: 'bg-gray-100 text-gray-700',
  gold: 'bg-yellow-100 text-yellow-800',
  platinum: 'bg-blue-100 text-blue-800',
}

export default function Loyalty() {
  const [showRewardForm, setShowRewardForm] = useState(false)
  const [rewardForm, setRewardForm] = useState<RewardFormData>(emptyRewardForm)
  const [editingRewardId, setEditingRewardId] = useState<string | null>(null)
  const [adjustAccountId, setAdjustAccountId] = useState<string | null>(null)
  const [adjustPoints, setAdjustPoints] = useState(0)
  const [adjustNote, setAdjustNote] = useState('')
  const [expandedAccountId, setExpandedAccountId] = useState<string | null>(null)

  const { data: rewards = [], isLoading: rewardsLoading } = useLoyaltyRewards()
  const { data: accountsData, isLoading: accountsLoading } = useLoyaltyAccounts()
  const accounts: LoyaltyAccount[] = accountsData?.items ?? []

  const createReward = useCreateLoyaltyReward()
  const updateReward = useUpdateLoyaltyReward()
  const adjustPoints_ = useAdjustLoyaltyPoints()

  function handleRewardSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (editingRewardId) {
      updateReward.mutate({ id: editingRewardId, data: rewardForm })
    } else {
      createReward.mutate(rewardForm)
    }
    setShowRewardForm(false)
    setRewardForm(emptyRewardForm)
    setEditingRewardId(null)
  }

  function handleEditReward(r: LoyaltyReward) {
    setRewardForm({ name: r.name, description: r.description, type: r.type, points_cost: r.points_cost, value: r.value, active: r.active })
    setEditingRewardId(r.id)
    setShowRewardForm(true)
  }

  function handleAdjustSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (!adjustAccountId) return
    adjustPoints_.mutate({ id: adjustAccountId, points: adjustPoints, note: adjustNote })
    setAdjustAccountId(null)
    setAdjustPoints(0)
    setAdjustNote('')
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Loyalty & Rewards</h1>
        <p className="mt-1 text-sm text-gray-500">Manage points system, tiers, and rewards</p>
      </div>

      <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
        {/* Rewards Panel */}
        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <h2 className="text-base font-semibold text-gray-900">Rewards</h2>
            <button
              onClick={() => { setShowRewardForm(true); setEditingRewardId(null); setRewardForm(emptyRewardForm) }}
              className="inline-flex items-center gap-1.5 rounded-lg bg-gray-900 px-3 py-1.5 text-xs font-medium text-white hover:bg-gray-800 transition-colors"
            >
              <Plus className="h-3.5 w-3.5" />
              Add Reward
            </button>
          </div>

          {showRewardForm && (
            <div className="rounded-xl border border-gray-200 bg-white p-4 shadow-sm">
              <div className="mb-3 flex items-center justify-between">
                <h3 className="text-sm font-semibold text-gray-900">{editingRewardId ? 'Edit Reward' : 'New Reward'}</h3>
                <button onClick={() => setShowRewardForm(false)} className="text-gray-400 hover:text-gray-600">
                  <X className="h-4 w-4" />
                </button>
              </div>
              <form onSubmit={handleRewardSubmit} className="space-y-3">
                <div>
                  <label className="block text-xs font-medium text-gray-700">Name</label>
                  <input type="text" required value={rewardForm.name} onChange={(e) => setRewardForm({ ...rewardForm, name: e.target.value })}
                    className="mt-1 w-full rounded-lg border border-gray-200 px-3 py-1.5 text-sm focus:border-gray-400 focus:outline-none" />
                </div>
                <div>
                  <label className="block text-xs font-medium text-gray-700">Description</label>
                  <input type="text" value={rewardForm.description} onChange={(e) => setRewardForm({ ...rewardForm, description: e.target.value })}
                    className="mt-1 w-full rounded-lg border border-gray-200 px-3 py-1.5 text-sm focus:border-gray-400 focus:outline-none" />
                </div>
                <div className="grid grid-cols-2 gap-3">
                  <div>
                    <label className="block text-xs font-medium text-gray-700">Type</label>
                    <select value={rewardForm.type} onChange={(e) => setRewardForm({ ...rewardForm, type: e.target.value as RewardType })}
                      className="mt-1 w-full rounded-lg border border-gray-200 px-3 py-1.5 text-sm focus:border-gray-400 focus:outline-none">
                      <option value="fixed_discount">Fixed Discount</option>
                      <option value="points_multiplier">Points Multiplier</option>
                      <option value="free_shipping">Free Shipping</option>
                      <option value="free_product">Free Product</option>
                    </select>
                  </div>
                  <div>
                    <label className="block text-xs font-medium text-gray-700">Points Cost</label>
                    <input type="number" min="0" required value={rewardForm.points_cost}
                      onChange={(e) => setRewardForm({ ...rewardForm, points_cost: parseInt(e.target.value) })}
                      className="mt-1 w-full rounded-lg border border-gray-200 px-3 py-1.5 text-sm focus:border-gray-400 focus:outline-none" />
                  </div>
                </div>
                <div className="grid grid-cols-2 gap-3">
                  <div>
                    <label className="block text-xs font-medium text-gray-700">Value</label>
                    <input type="number" step="0.01" min="0" value={rewardForm.value}
                      onChange={(e) => setRewardForm({ ...rewardForm, value: parseFloat(e.target.value) })}
                      className="mt-1 w-full rounded-lg border border-gray-200 px-3 py-1.5 text-sm focus:border-gray-400 focus:outline-none" />
                  </div>
                  <div className="flex items-end pb-1">
                    <label className="flex items-center gap-2 text-xs font-medium text-gray-700">
                      <input type="checkbox" checked={rewardForm.active} onChange={(e) => setRewardForm({ ...rewardForm, active: e.target.checked })} className="rounded" />
                      Active
                    </label>
                  </div>
                </div>
                <div className="flex gap-2">
                  <button type="submit" className="rounded-lg bg-gray-900 px-3 py-1.5 text-xs font-medium text-white hover:bg-gray-800 transition-colors">
                    {editingRewardId ? 'Update' : 'Create'}
                  </button>
                  <button type="button" onClick={() => setShowRewardForm(false)} className="rounded-lg border border-gray-200 px-3 py-1.5 text-xs font-medium text-gray-700 hover:bg-gray-50 transition-colors">
                    Cancel
                  </button>
                </div>
              </form>
            </div>
          )}

          <div className="rounded-xl border border-gray-200 bg-white shadow-sm">
            {rewardsLoading ? (
              <div className="flex items-center justify-center py-10">
                <div className="h-5 w-5 animate-spin rounded-full border-2 border-gray-300 border-t-gray-900" />
              </div>
            ) : rewards.length === 0 ? (
              <div className="flex flex-col items-center justify-center py-10 text-gray-400">
                <Award className="mb-2 h-8 w-8" />
                <p className="text-sm">No rewards yet</p>
              </div>
            ) : (
              <div className="divide-y divide-gray-50">
                {rewards.map((r) => (
                  <div key={r.id} className="flex items-center gap-3 px-4 py-3">
                    <div className="flex-1 min-w-0">
                      <div className="text-sm font-medium text-gray-900">{r.name}</div>
                      <div className="text-xs text-gray-500">{r.type} · {r.points_cost} pts</div>
                    </div>
                    <span className={`inline-flex rounded-full px-2 py-0.5 text-xs font-medium ${r.active ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-600'}`}>
                      {r.active ? 'Active' : 'Off'}
                    </span>
                    <button onClick={() => handleEditReward(r)} className="rounded-md p-1 text-gray-400 hover:bg-gray-100 hover:text-gray-600 transition-colors">
                      <Pencil className="h-3.5 w-3.5" />
                    </button>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>

        {/* Accounts Panel */}
        <div className="space-y-4">
          <h2 className="text-base font-semibold text-gray-900">Loyalty Accounts</h2>
          <div className="rounded-xl border border-gray-200 bg-white shadow-sm">
            {accountsLoading ? (
              <div className="flex items-center justify-center py-10">
                <div className="h-5 w-5 animate-spin rounded-full border-2 border-gray-300 border-t-gray-900" />
              </div>
            ) : accounts.length === 0 ? (
              <div className="flex flex-col items-center justify-center py-10 text-gray-400">
                <p className="text-sm">No accounts found</p>
              </div>
            ) : (
              <div className="divide-y divide-gray-50">
                {accounts.map((acc) => (
                  <div key={acc.id}>
                    <div className="flex items-center gap-3 px-4 py-3">
                      <div className="flex-1 min-w-0">
                        <div className="text-sm font-medium text-gray-900">{acc.customer_name}</div>
                        <div className="text-xs text-gray-500">{acc.customer_email}</div>
                      </div>
                      <span className={`inline-flex rounded-full px-2 py-0.5 text-xs font-medium ${tierColors[acc.tier]}`}>
                        {acc.tier}
                      </span>
                      <span className="text-sm font-bold text-indigo-700">{acc.points.toLocaleString()} pts</span>
                      <button
                        onClick={() => { setAdjustAccountId(acc.id); setAdjustPoints(0); setAdjustNote('') }}
                        className="rounded-lg border border-gray-200 px-2 py-1 text-xs font-medium text-gray-700 hover:bg-gray-50 transition-colors"
                      >
                        Adjust
                      </button>
                      <button
                        onClick={() => setExpandedAccountId(expandedAccountId === acc.id ? null : acc.id)}
                        className="rounded-md p-1 text-gray-400 hover:bg-gray-100 transition-colors"
                      >
                        {expandedAccountId === acc.id ? <ChevronUp className="h-3.5 w-3.5" /> : <ChevronDown className="h-3.5 w-3.5" />}
                      </button>
                    </div>
                    {expandedAccountId === acc.id && (
                      <div className="border-t border-gray-50 bg-gray-50 px-4 py-2">
                        <TransactionList accountId={acc.id} />
                      </div>
                    )}
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Adjust Points Dialog */}
      {adjustAccountId && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40">
          <div className="w-full max-w-sm rounded-xl bg-white p-6 shadow-xl">
            <div className="mb-4 flex items-center justify-between">
              <h2 className="text-lg font-semibold text-gray-900">Adjust Points</h2>
              <button onClick={() => setAdjustAccountId(null)} className="text-gray-400 hover:text-gray-600">
                <X className="h-5 w-5" />
              </button>
            </div>
            <form onSubmit={handleAdjustSubmit} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700">Points (positive = add, negative = deduct)</label>
                <input type="number" required value={adjustPoints} onChange={(e) => setAdjustPoints(parseInt(e.target.value))}
                  className="mt-1 w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-gray-400 focus:outline-none" />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700">Note</label>
                <input type="text" required value={adjustNote} onChange={(e) => setAdjustNote(e.target.value)}
                  className="mt-1 w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-gray-400 focus:outline-none" />
              </div>
              <div className="flex gap-3">
                <button type="submit" className="flex-1 rounded-lg bg-gray-900 py-2.5 text-sm font-medium text-white hover:bg-gray-800 transition-colors">
                  Apply
                </button>
                <button type="button" onClick={() => setAdjustAccountId(null)} className="flex-1 rounded-lg border border-gray-200 py-2.5 text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors">
                  Cancel
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}
