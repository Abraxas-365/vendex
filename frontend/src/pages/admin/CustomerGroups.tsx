import { useState } from 'react'
import { Plus, Users, X, Trash2, ChevronRight, UserPlus } from 'lucide-react'
import type { CustomerGroup, GroupMembership, GroupRules } from '../../types'
import {
  useCustomerGroups,
  useCreateCustomerGroup,
  useUpdateCustomerGroup,
  useDeleteCustomerGroup,
  useGroupMembers,
  useAddGroupMember,
  useRemoveGroupMember,
} from '../../lib/hooks'

// ---------------------------------------------------------------------------
// Group form dialog
// ---------------------------------------------------------------------------

interface GroupFormData {
  name: string
  description: string
  min_order_count: string
  min_total_spent: string
  tags: string
}

function groupToForm(group: CustomerGroup): GroupFormData {
  return {
    name: group.name,
    description: group.description,
    min_order_count: group.rules.min_order_count != null ? String(group.rules.min_order_count) : '',
    min_total_spent: group.rules.min_total_spent != null ? (group.rules.min_total_spent / 100).toFixed(2) : '',
    tags: group.rules.tags ? group.rules.tags.join(', ') : '',
  }
}

const emptyGroupForm: GroupFormData = {
  name: '',
  description: '',
  min_order_count: '',
  min_total_spent: '',
  tags: '',
}

function formToRules(form: GroupFormData): GroupRules {
  const rules: GroupRules = {}
  if (form.min_order_count) rules.min_order_count = parseInt(form.min_order_count)
  if (form.min_total_spent) rules.min_total_spent = Math.round(parseFloat(form.min_total_spent) * 100)
  if (form.tags.trim()) {
    rules.tags = form.tags
      .split(',')
      .map((t) => t.trim())
      .filter(Boolean)
  }
  return rules
}

interface GroupDialogProps {
  initial?: CustomerGroup
  onClose: () => void
  onSave: (data: { name: string; description: string; rules: GroupRules }) => void
  saving: boolean
  error?: string
}

function GroupDialog({ initial, onClose, onSave, saving, error }: GroupDialogProps) {
  const [form, setForm] = useState<GroupFormData>(initial ? groupToForm(initial) : emptyGroupForm)

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    onSave({
      name: form.name,
      description: form.description,
      rules: formToRules(form),
    })
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4">
      <div className="w-full max-w-md rounded-xl bg-white shadow-xl">
        <div className="flex items-center justify-between border-b border-slate-100 px-6 py-4">
          <h2 className="text-base font-semibold text-slate-800">
            {initial ? 'Edit Group' : 'New Customer Group'}
          </h2>
          <button onClick={onClose} className="text-slate-400 hover:text-slate-600">
            <X size={18} />
          </button>
        </div>
        <form onSubmit={handleSubmit} className="space-y-4 px-6 py-4">
          <div>
            <label className="mb-1 block text-sm font-medium text-slate-700">Group Name</label>
            <input
              type="text"
              required
              value={form.name}
              onChange={(e) => setForm({ ...form, name: e.target.value })}
              className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm focus:border-indigo-400 focus:outline-none"
              placeholder="e.g. VIP Customers"
            />
          </div>
          <div>
            <label className="mb-1 block text-sm font-medium text-slate-700">Description</label>
            <textarea
              rows={2}
              value={form.description}
              onChange={(e) => setForm({ ...form, description: e.target.value })}
              className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm focus:border-indigo-400 focus:outline-none resize-none"
              placeholder="Optional description"
            />
          </div>

          <div className="rounded-lg border border-slate-100 bg-slate-50 p-4 space-y-3">
            <p className="text-xs font-semibold uppercase tracking-wide text-slate-400">
              Auto-qualify Rules <span className="font-normal normal-case">(optional)</span>
            </p>
            <div>
              <label className="mb-1 block text-sm font-medium text-slate-700">
                Min. Order Count
              </label>
              <input
                type="number"
                min="0"
                value={form.min_order_count}
                onChange={(e) => setForm({ ...form, min_order_count: e.target.value })}
                className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm focus:border-indigo-400 focus:outline-none"
                placeholder="e.g. 5"
              />
            </div>
            <div>
              <label className="mb-1 block text-sm font-medium text-slate-700">
                Min. Total Spent ($)
              </label>
              <input
                type="number"
                min="0"
                step="0.01"
                value={form.min_total_spent}
                onChange={(e) => setForm({ ...form, min_total_spent: e.target.value })}
                className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm focus:border-indigo-400 focus:outline-none"
                placeholder="e.g. 500.00"
              />
            </div>
            <div>
              <label className="mb-1 block text-sm font-medium text-slate-700">
                Tags <span className="font-normal text-slate-400">(comma-separated)</span>
              </label>
              <input
                type="text"
                value={form.tags}
                onChange={(e) => setForm({ ...form, tags: e.target.value })}
                className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm focus:border-indigo-400 focus:outline-none"
                placeholder="wholesale, partner"
              />
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
              {saving ? 'Saving…' : 'Save Group'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}

// ---------------------------------------------------------------------------
// Add member dialog
// ---------------------------------------------------------------------------

interface AddMemberDialogProps {
  onClose: () => void
  onAdd: (customerId: string) => void
  saving: boolean
  error?: string
}

function AddMemberDialog({ onClose, onAdd, saving, error }: AddMemberDialogProps) {
  const [customerId, setCustomerId] = useState('')

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    onAdd(customerId.trim())
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4">
      <div className="w-full max-w-sm rounded-xl bg-white shadow-xl">
        <div className="flex items-center justify-between border-b border-slate-100 px-6 py-4">
          <h2 className="text-base font-semibold text-slate-800">Add Member</h2>
          <button onClick={onClose} className="text-slate-400 hover:text-slate-600">
            <X size={18} />
          </button>
        </div>
        <form onSubmit={handleSubmit} className="space-y-4 px-6 py-4">
          <div>
            <label className="mb-1 block text-sm font-medium text-slate-700">Customer ID</label>
            <input
              type="text"
              required
              value={customerId}
              onChange={(e) => setCustomerId(e.target.value)}
              className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm focus:border-indigo-400 focus:outline-none"
              placeholder="Paste customer UUID"
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
              disabled={saving || !customerId.trim()}
              className="flex-1 rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700 disabled:opacity-60"
            >
              {saving ? 'Adding…' : 'Add Member'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}

// ---------------------------------------------------------------------------
// Members panel
// ---------------------------------------------------------------------------

interface MembersPanelProps {
  group: CustomerGroup
}

function MembersPanel({ group }: MembersPanelProps) {
  const [showAddDialog, setShowAddDialog] = useState(false)

  const { data: members = [], isLoading } = useGroupMembers(group.id)
  const addMember = useAddGroupMember()
  const removeMember = useRemoveGroupMember()

  function handleAddMember(customerId: string) {
    addMember.mutate(
      { groupId: group.id, customerId },
      {
        onSuccess: () => {
          setShowAddDialog(false)
        },
      },
    )
  }

  function handleRemoveMember(membership: GroupMembership) {
    if (!confirm(`Remove customer ${membership.customer_id} from this group?`)) return
    removeMember.mutate({ groupId: group.id, customerId: membership.customer_id })
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <div>
          <h3 className="text-sm font-semibold text-slate-800">Members</h3>
          <p className="mt-0.5 text-xs text-slate-500">Group: {group.name}</p>
        </div>
        <button
          onClick={() => setShowAddDialog(true)}
          className="flex items-center gap-1.5 rounded-lg bg-indigo-600 px-3 py-1.5 text-sm font-medium text-white hover:bg-indigo-700"
        >
          <UserPlus size={15} />
          Add Member
        </button>
      </div>

      {isLoading && (
        <div className="py-8 text-center text-sm text-slate-400">Loading members…</div>
      )}

      {!isLoading && members.length === 0 && (
        <div className="rounded-lg border border-dashed border-slate-200 py-8 text-center text-sm text-slate-400">
          No members yet. Add one above.
        </div>
      )}

      {members.length > 0 && (
        <div className="overflow-hidden rounded-lg border border-slate-200">
          <table className="w-full text-sm">
            <thead className="bg-slate-50 text-xs font-semibold uppercase tracking-wide text-slate-500">
              <tr>
                <th className="px-4 py-2.5 text-left">Customer ID</th>
                <th className="px-4 py-2.5 text-left">Joined</th>
                <th className="px-4 py-2.5" />
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100">
              {members.map((m) => (
                <tr key={m.id} className="hover:bg-slate-50">
                  <td className="px-4 py-3 font-mono text-xs text-slate-600">{m.customer_id}</td>
                  <td className="px-4 py-3 text-slate-500">
                    {new Date(m.joined_at).toLocaleDateString()}
                  </td>
                  <td className="px-4 py-3 text-right">
                    <button
                      onClick={() => handleRemoveMember(m)}
                      disabled={removeMember.isPending}
                      className="text-xs text-red-500 hover:underline disabled:opacity-50"
                    >
                      Remove
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {showAddDialog && (
        <AddMemberDialog
          onClose={() => setShowAddDialog(false)}
          onAdd={handleAddMember}
          saving={addMember.isPending}
          error={addMember.error?.message}
        />
      )}
    </div>
  )
}

// ---------------------------------------------------------------------------
// Rules summary helper
// ---------------------------------------------------------------------------

function RulesSummary({ rules }: { rules: GroupRules }) {
  const parts: string[] = []
  if (rules.min_order_count != null) parts.push(`≥${rules.min_order_count} orders`)
  if (rules.min_total_spent != null) {
    parts.push(`≥$${(rules.min_total_spent / 100).toFixed(0)} spent`)
  }
  if (rules.tags && rules.tags.length > 0) parts.push(`tags: ${rules.tags.join(', ')}`)
  if (parts.length === 0) return <span className="text-slate-400 text-xs">No rules</span>
  return <span className="text-xs text-slate-600">{parts.join(' · ')}</span>
}

// ---------------------------------------------------------------------------
// Main page
// ---------------------------------------------------------------------------

export default function CustomerGroups() {
  const [showGroupForm, setShowGroupForm] = useState(false)
  const [editingGroup, setEditingGroup] = useState<CustomerGroup | null>(null)
  const [selectedGroup, setSelectedGroup] = useState<CustomerGroup | null>(null)

  const { data: groups = [], isLoading } = useCustomerGroups()
  const createGroup = useCreateCustomerGroup()
  const updateGroup = useUpdateCustomerGroup()
  const deleteGroup = useDeleteCustomerGroup()

  function handleSaveGroup(data: { name: string; description: string; rules: GroupRules }) {
    if (editingGroup) {
      updateGroup.mutate(
        { id: editingGroup.id, ...data },
        {
          onSuccess: (updated) => {
            setEditingGroup(null)
            if (selectedGroup?.id === updated.id) setSelectedGroup(updated)
          },
        },
      )
    } else {
      createGroup.mutate(data, {
        onSuccess: () => {
          setShowGroupForm(false)
        },
      })
    }
  }

  function handleDeleteGroup(group: CustomerGroup) {
    if (!confirm(`Delete group "${group.name}"? All memberships will be removed.`)) return
    deleteGroup.mutate(group.id, {
      onSuccess: () => {
        if (selectedGroup?.id === group.id) setSelectedGroup(null)
      },
    })
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-xl font-semibold text-slate-800">Customer Groups</h1>
          <p className="mt-0.5 text-sm text-slate-500">
            Organise customers into groups for targeted promotions and pricing.
          </p>
        </div>
        <button
          onClick={() => setShowGroupForm(true)}
          className="flex items-center gap-2 rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700"
        >
          <Plus size={16} />
          New Group
        </button>
      </div>

      {/* Two-panel layout */}
      <div className="flex gap-6">
        {/* Left: Groups list */}
        <div className="w-72 shrink-0">
          <div className="rounded-xl border border-slate-200 bg-white overflow-hidden">
            <div className="border-b border-slate-100 px-4 py-3">
              <p className="text-xs font-semibold uppercase tracking-widest text-slate-400">
                Groups
              </p>
            </div>

            {isLoading && (
              <div className="py-8 text-center text-sm text-slate-400">Loading…</div>
            )}

            {!isLoading && groups.length === 0 && (
              <div className="py-8 text-center">
                <Users size={32} className="mx-auto mb-2 text-slate-300" />
                <p className="text-sm text-slate-400">No groups yet</p>
              </div>
            )}

            {groups.length > 0 && (
              <ul className="divide-y divide-slate-100">
                {groups.map((group) => (
                  <li key={group.id}>
                    <button
                      onClick={() => setSelectedGroup(group)}
                      className={`flex w-full items-center gap-3 px-4 py-3 text-left transition-colors hover:bg-slate-50 ${
                        selectedGroup?.id === group.id
                          ? 'bg-indigo-50 text-indigo-700'
                          : 'text-slate-700'
                      }`}
                    >
                      <Users size={16} className="shrink-0 text-slate-400" />
                      <div className="flex-1 min-w-0">
                        <p className="text-sm font-medium truncate">{group.name}</p>
                        <p className="text-xs text-slate-400 truncate">
                          {group.member_count} member{group.member_count !== 1 ? 's' : ''}
                        </p>
                      </div>
                      <ChevronRight size={14} className="shrink-0 text-slate-300" />
                    </button>
                  </li>
                ))}
              </ul>
            )}
          </div>

          {/* Group actions */}
          {selectedGroup && (
            <div className="mt-3 flex gap-2">
              <button
                onClick={() => setEditingGroup(selectedGroup)}
                className="flex-1 rounded-lg border border-slate-200 px-3 py-1.5 text-xs font-medium text-slate-600 hover:bg-slate-50"
              >
                Edit Group
              </button>
              <button
                onClick={() => handleDeleteGroup(selectedGroup)}
                disabled={deleteGroup.isPending}
                className="flex items-center gap-1 rounded-lg border border-red-200 px-3 py-1.5 text-xs font-medium text-red-600 hover:bg-red-50 disabled:opacity-50"
              >
                <Trash2 size={13} />
                Delete
              </button>
            </div>
          )}
        </div>

        {/* Right: Members panel (or group detail) */}
        <div className="flex-1">
          {selectedGroup ? (
            <div className="rounded-xl border border-slate-200 bg-white p-6 space-y-6">
              {/* Group info header */}
              <div className="border-b border-slate-100 pb-4">
                <h2 className="text-base font-semibold text-slate-800">{selectedGroup.name}</h2>
                {selectedGroup.description && (
                  <p className="mt-1 text-sm text-slate-500">{selectedGroup.description}</p>
                )}
                <div className="mt-2">
                  <RulesSummary rules={selectedGroup.rules} />
                </div>
              </div>
              <MembersPanel group={selectedGroup} />
            </div>
          ) : (
            <div className="flex h-full min-h-48 items-center justify-center rounded-xl border border-dashed border-slate-200 bg-white text-center">
              <div>
                <Users size={32} className="mx-auto mb-2 text-slate-300" />
                <p className="text-sm text-slate-400">Select a group to manage its members</p>
              </div>
            </div>
          )}
        </div>
      </div>

      {/* Group dialog */}
      {(showGroupForm || editingGroup) && (
        <GroupDialog
          initial={editingGroup ?? undefined}
          onClose={() => {
            setShowGroupForm(false)
            setEditingGroup(null)
          }}
          onSave={handleSaveGroup}
          saving={createGroup.isPending || updateGroup.isPending}
          error={createGroup.error?.message ?? updateGroup.error?.message}
        />
      )}
    </div>
  )
}
