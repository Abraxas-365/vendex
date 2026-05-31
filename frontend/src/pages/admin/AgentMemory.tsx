import { useState } from 'react'
import {
  Brain,
  Plus,
  Search,
  Trash2,
  Pencil,
  X,
} from 'lucide-react'
import {
  useAgentMemories,
  useSearchAgentMemories,
  useCreateAgentMemory,
  useUpdateAgentMemory,
  useDeleteAgentMemory,
} from '../../lib/hooks'
import type { AgentMemory, CreateMemoryRequest, UpdateMemoryRequest } from '../../types'

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

type MemoryCategory = 'brand' | 'product' | 'seo' | 'general' | 'decision'
type CategoryFilter = 'all' | MemoryCategory

const CATEGORIES: MemoryCategory[] = ['brand', 'product', 'seo', 'general', 'decision']

function categoryBadgeClass(category: string): string {
  switch (category) {
    case 'brand':
      return 'bg-purple-100 text-purple-700'
    case 'product':
      return 'bg-blue-100 text-blue-700'
    case 'seo':
      return 'bg-green-100 text-green-700'
    case 'general':
      return 'bg-slate-100 text-slate-600'
    case 'decision':
      return 'bg-amber-100 text-amber-700'
    default:
      return 'bg-slate-100 text-slate-600'
  }
}

function sourceBadgeClass(source: string): string {
  return source === 'agent'
    ? 'bg-indigo-100 text-indigo-700'
    : 'bg-teal-100 text-teal-700'
}

function formatDate(dateStr: string): string {
  try {
    return new Date(dateStr).toLocaleDateString()
  } catch {
    return dateStr
  }
}

// ---------------------------------------------------------------------------
// Memory form dialog (create / edit)
// ---------------------------------------------------------------------------

interface MemoryFormDialogProps {
  initial?: AgentMemory
  onClose: () => void
}

function MemoryFormDialog({ initial, onClose }: MemoryFormDialogProps) {
  const [title, setTitle] = useState(initial?.title ?? '')
  const [category, setCategory] = useState<MemoryCategory>(
    (initial?.category as MemoryCategory) ?? 'general',
  )
  const [content, setContent] = useState(initial?.content ?? '')
  const [tagsInput, setTagsInput] = useState((initial?.tags ?? []).join(', '))

  const createMutation = useCreateAgentMemory()
  const updateMutation = useUpdateAgentMemory()

  const isPending = createMutation.isPending || updateMutation.isPending
  const error = createMutation.error ?? updateMutation.error

  function handleSubmit() {
    const tags = tagsInput
      .split(',')
      .map((t) => t.trim())
      .filter(Boolean)

    if (initial) {
      const data: UpdateMemoryRequest = { title, category, content, tags }
      updateMutation.mutate({ id: initial.id, data }, { onSuccess: onClose })
    } else {
      const data: CreateMemoryRequest = { title, category, content, tags }
      createMutation.mutate(data, { onSuccess: onClose })
    }
  }

  const isValid = title.trim().length > 0 && content.trim().length > 0

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40">
      <div className="w-full max-w-lg rounded-xl bg-white shadow-xl p-6 mx-4 max-h-[90vh] overflow-y-auto">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-base font-semibold text-slate-800">
            {initial ? 'Edit Memory' : 'Add Memory'}
          </h2>
          <button
            onClick={onClose}
            className="text-slate-400 hover:text-slate-600 transition-colors"
          >
            <X size={18} />
          </button>
        </div>

        <div className="space-y-4">
          {/* Title */}
          <div>
            <label className="block text-xs font-medium text-slate-600 mb-1">Title *</label>
            <input
              type="text"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder="Memory title…"
              className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm text-slate-700 placeholder:text-slate-400 focus:border-indigo-400 focus:outline-none focus:ring-2 focus:ring-indigo-100"
            />
          </div>

          {/* Category */}
          <div>
            <label className="block text-xs font-medium text-slate-600 mb-1">Category *</label>
            <select
              value={category}
              onChange={(e) => setCategory(e.target.value as MemoryCategory)}
              className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm text-slate-700 focus:border-indigo-400 focus:outline-none focus:ring-2 focus:ring-indigo-100 bg-white"
            >
              {CATEGORIES.map((cat) => (
                <option key={cat} value={cat}>
                  {cat.charAt(0).toUpperCase() + cat.slice(1)}
                </option>
              ))}
            </select>
          </div>

          {/* Content */}
          <div>
            <label className="block text-xs font-medium text-slate-600 mb-1">Content *</label>
            <textarea
              value={content}
              onChange={(e) => setContent(e.target.value)}
              placeholder="Memory content…"
              rows={5}
              className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm text-slate-700 placeholder:text-slate-400 focus:border-indigo-400 focus:outline-none focus:ring-2 focus:ring-indigo-100 resize-none"
            />
          </div>

          {/* Tags */}
          <div>
            <label className="block text-xs font-medium text-slate-600 mb-1">
              Tags (comma-separated)
            </label>
            <input
              type="text"
              value={tagsInput}
              onChange={(e) => setTagsInput(e.target.value)}
              placeholder="e.g. brand, tone, voice"
              className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm text-slate-700 placeholder:text-slate-400 focus:border-indigo-400 focus:outline-none focus:ring-2 focus:ring-indigo-100"
            />
          </div>
        </div>

        {error && (
          <p className="mt-3 text-xs text-red-600">{error.message}</p>
        )}

        <div className="mt-5 flex items-center justify-end gap-2">
          <button
            onClick={onClose}
            disabled={isPending}
            className="rounded-lg border border-slate-200 px-4 py-2 text-sm font-medium text-slate-600 hover:bg-slate-50 disabled:opacity-50 transition-colors"
          >
            Cancel
          </button>
          <button
            onClick={handleSubmit}
            disabled={isPending || !isValid}
            className="rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            {isPending ? 'Saving…' : initial ? 'Save Changes' : 'Add Memory'}
          </button>
        </div>
      </div>
    </div>
  )
}

// ---------------------------------------------------------------------------
// Delete confirmation dialog
// ---------------------------------------------------------------------------

interface DeleteDialogProps {
  memory: AgentMemory
  onClose: () => void
}

function DeleteDialog({ memory, onClose }: DeleteDialogProps) {
  const deleteMutation = useDeleteAgentMemory()

  function handleDelete() {
    deleteMutation.mutate(memory.id, { onSuccess: onClose })
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40">
      <div className="w-full max-w-sm rounded-xl bg-white shadow-xl p-6 mx-4">
        <h2 className="text-base font-semibold text-slate-800 mb-2">Delete Memory</h2>
        <p className="text-sm text-slate-500 mb-1">
          Are you sure you want to delete{' '}
          <span className="font-medium text-slate-700">&ldquo;{memory.title}&rdquo;</span>?
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
// Memory card
// ---------------------------------------------------------------------------

interface MemoryCardProps {
  memory: AgentMemory
  onEdit: (memory: AgentMemory) => void
  onDelete: (memory: AgentMemory) => void
}

function MemoryCard({ memory, onEdit, onDelete }: MemoryCardProps) {
  const [expanded, setExpanded] = useState(false)

  return (
    <div
      className="rounded-xl border border-slate-200 bg-white p-4 shadow-sm hover:shadow-md transition-shadow cursor-pointer"
      onClick={() => setExpanded(!expanded)}
    >
      <div className="flex items-start justify-between gap-2">
        <div className="min-w-0 flex-1">
          <div className="flex items-center gap-2 flex-wrap">
            <span
              className={`inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium capitalize ${categoryBadgeClass(memory.category)}`}
            >
              {memory.category}
            </span>
            <span
              className={`inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium capitalize ${sourceBadgeClass(memory.source)}`}
            >
              {memory.source}
            </span>
          </div>
          <h3 className="mt-1 text-sm font-semibold text-slate-800 truncate">{memory.title}</h3>
        </div>
        <div className="flex items-center gap-1 shrink-0" onClick={(e) => e.stopPropagation()}>
          <button
            onClick={() => onEdit(memory)}
            className="rounded-md p-1.5 text-slate-400 hover:text-indigo-600 hover:bg-indigo-50 transition-colors"
            title="Edit"
          >
            <Pencil size={13} />
          </button>
          <button
            onClick={() => onDelete(memory)}
            className="rounded-md p-1.5 text-slate-400 hover:text-red-600 hover:bg-red-50 transition-colors"
            title="Delete"
          >
            <Trash2 size={13} />
          </button>
        </div>
      </div>

      <p
        className={`mt-2 text-sm text-slate-600 ${expanded ? '' : 'line-clamp-2'}`}
      >
        {memory.content}
      </p>

      {memory.tags && memory.tags.length > 0 && (
        <div className="mt-2 flex flex-wrap gap-1">
          {memory.tags.map((tag) => (
            <span
              key={tag}
              className="rounded-full bg-slate-100 px-2 py-0.5 text-xs text-slate-500"
            >
              {tag}
            </span>
          ))}
        </div>
      )}

      <p className="mt-2 text-xs text-slate-400">{formatDate(memory.updated_at)}</p>
    </div>
  )
}

// ---------------------------------------------------------------------------
// Main page
// ---------------------------------------------------------------------------

export default function AgentMemory() {
  const [search, setSearch] = useState('')
  const [categoryFilter, setCategoryFilter] = useState<CategoryFilter>('all')
  const [showCreateDialog, setShowCreateDialog] = useState(false)
  const [editingMemory, setEditingMemory] = useState<AgentMemory | null>(null)
  const [deletingMemory, setDeletingMemory] = useState<AgentMemory | null>(null)

  const isSearching = search.trim().length > 0 || categoryFilter !== 'all'

  const { data: listPage, isLoading: loadingList } = useAgentMemories(
    { page: 1, page_size: 50 },
  )
  const { data: searchPage, isLoading: loadingSearch } = useSearchAgentMemories(
    {
      q: search.trim() || undefined,
      category: categoryFilter !== 'all' ? categoryFilter : undefined,
    },
    isSearching,
  )

  const isLoading = isSearching ? loadingSearch : loadingList
  const memories = (isSearching ? searchPage?.items : listPage?.items) ?? []

  return (
    <div className="space-y-6">
      {/* Page header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <Brain size={24} className="text-indigo-600" />
          <div>
            <h1 className="text-xl font-semibold text-slate-800">Agent Memory</h1>
            <p className="text-sm text-slate-500">Manage knowledge entries for your AI agents</p>
          </div>
        </div>
        <button
          onClick={() => setShowCreateDialog(true)}
          className="inline-flex items-center gap-2 rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700 transition-colors"
        >
          <Plus size={16} />
          Add Memory
        </button>
      </div>

      {/* Filters */}
      <div className="flex flex-col sm:flex-row gap-3">
        {/* Search */}
        <div className="relative flex-1 max-w-sm">
          <Search size={14} className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
          <input
            type="text"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Search memories…"
            className="w-full rounded-lg border border-slate-200 bg-white pl-8 pr-3 py-2 text-sm text-slate-700 placeholder:text-slate-400 focus:border-indigo-400 focus:outline-none focus:ring-2 focus:ring-indigo-100"
          />
        </div>

        {/* Category filter */}
        <div className="flex gap-1 rounded-lg border border-slate-200 bg-white p-1 flex-wrap">
          {(['all', ...CATEGORIES] as const).map((cat) => (
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

      {/* Memory grid */}
      {isLoading ? (
        <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
          {[...Array(6)].map((_, i) => (
            <div key={i} className="h-36 rounded-xl border border-slate-200 bg-white animate-pulse" />
          ))}
        </div>
      ) : memories.length === 0 ? (
        <div className="rounded-xl border border-slate-200 bg-white p-12 text-center">
          <Brain size={32} className="mx-auto mb-3 text-slate-300" />
          <p className="text-sm text-slate-500">
            {isSearching ? 'No memories match your search' : 'No memories yet'}
          </p>
          {!isSearching && (
            <button
              onClick={() => setShowCreateDialog(true)}
              className="mt-3 text-xs font-medium text-indigo-600 hover:underline"
            >
              Add your first memory →
            </button>
          )}
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
          {memories.map((memory) => (
            <MemoryCard
              key={memory.id}
              memory={memory}
              onEdit={setEditingMemory}
              onDelete={setDeletingMemory}
            />
          ))}
        </div>
      )}

      {/* Create dialog */}
      {showCreateDialog && (
        <MemoryFormDialog onClose={() => setShowCreateDialog(false)} />
      )}

      {/* Edit dialog */}
      {editingMemory && (
        <MemoryFormDialog
          initial={editingMemory}
          onClose={() => setEditingMemory(null)}
        />
      )}

      {/* Delete dialog */}
      {deletingMemory && (
        <DeleteDialog
          memory={deletingMemory}
          onClose={() => setDeletingMemory(null)}
        />
      )}
    </div>
  )
}
