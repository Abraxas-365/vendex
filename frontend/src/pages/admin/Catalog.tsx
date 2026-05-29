import { useState } from 'react'
import { Plus, Trash2, X, Search, Layers, FolderTree } from 'lucide-react'
import type { Category, Collection } from '../../types'
import {
  useCategories,
  useCreateCategory,
  useDeleteCategory,
  useCollections,
  useCreateCollection,
  useDeleteCollection,
} from '../../lib/hooks'

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  })
}

function slugify(value: string): string {
  return value
    .toLowerCase()
    .trim()
    .replace(/[^\w\s-]/g, '')
    .replace(/[\s_-]+/g, '-')
    .replace(/^-+|-+$/g, '')
}

// ---------------------------------------------------------------------------
// Category form
// ---------------------------------------------------------------------------

interface CategoryFormData {
  name: string
  slug: string
  parent_id: string
  description: string
}

const emptyCategoryForm: CategoryFormData = {
  name: '',
  slug: '',
  parent_id: '',
  description: '',
}

interface CategoryFormProps {
  categories: Category[]
  onSubmit: (data: CategoryFormData) => void
  onCancel: () => void
  isSubmitting: boolean
}

function CategoryForm({ categories, onSubmit, onCancel, isSubmitting }: CategoryFormProps) {
  const [form, setForm] = useState<CategoryFormData>(emptyCategoryForm)

  function handleNameChange(name: string) {
    setForm((prev) => ({
      ...prev,
      name,
      slug: slugify(name),
    }))
  }

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    onSubmit(form)
  }

  return (
    <div className="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
      <div className="mb-4 flex items-center justify-between">
        <h2 className="text-lg font-semibold text-gray-900">New Category</h2>
        <button onClick={onCancel} className="text-gray-400 hover:text-gray-600">
          <X className="h-5 w-5" />
        </button>
      </div>
      <form onSubmit={handleSubmit} className="grid grid-cols-1 gap-4 sm:grid-cols-2">
        <div>
          <label className="block text-sm font-medium text-gray-700">Name</label>
          <input
            type="text"
            required
            value={form.name}
            onChange={(e) => handleNameChange(e.target.value)}
            className="mt-1 w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-gray-400 focus:outline-none focus:ring-1 focus:ring-gray-400"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700">Slug</label>
          <input
            type="text"
            required
            value={form.slug}
            onChange={(e) => setForm({ ...form, slug: e.target.value })}
            className="mt-1 w-full rounded-lg border border-gray-200 px-3 py-2 text-sm font-mono focus:border-gray-400 focus:outline-none focus:ring-1 focus:ring-gray-400"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700">Parent Category</label>
          <select
            value={form.parent_id}
            onChange={(e) => setForm({ ...form, parent_id: e.target.value })}
            className="mt-1 w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-gray-400 focus:outline-none focus:ring-1 focus:ring-gray-400"
          >
            <option value="">None (Root category)</option>
            {categories.map((cat) => (
              <option key={cat.id} value={cat.id}>
                {cat.name}
              </option>
            ))}
          </select>
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700">Description</label>
          <input
            type="text"
            value={form.description}
            onChange={(e) => setForm({ ...form, description: e.target.value })}
            className="mt-1 w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-gray-400 focus:outline-none focus:ring-1 focus:ring-gray-400"
            placeholder="Optional description"
          />
        </div>
        <div className="flex items-end gap-3 sm:col-span-2">
          <button
            type="submit"
            disabled={isSubmitting}
            className="rounded-lg bg-gray-900 px-4 py-2.5 text-sm font-medium text-white hover:bg-gray-800 transition-colors disabled:opacity-60"
          >
            {isSubmitting ? 'Creating…' : 'Create Category'}
          </button>
          <button
            type="button"
            onClick={onCancel}
            className="rounded-lg border border-gray-200 bg-white px-4 py-2.5 text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors"
          >
            Cancel
          </button>
        </div>
      </form>
    </div>
  )
}

// ---------------------------------------------------------------------------
// Collection form
// ---------------------------------------------------------------------------

interface CollectionFormData {
  name: string
  slug: string
  description: string
}

const emptyCollectionForm: CollectionFormData = {
  name: '',
  slug: '',
  description: '',
}

interface CollectionFormProps {
  onSubmit: (data: CollectionFormData) => void
  onCancel: () => void
  isSubmitting: boolean
}

function CollectionForm({ onSubmit, onCancel, isSubmitting }: CollectionFormProps) {
  const [form, setForm] = useState<CollectionFormData>(emptyCollectionForm)

  function handleNameChange(name: string) {
    setForm((prev) => ({
      ...prev,
      name,
      slug: slugify(name),
    }))
  }

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    onSubmit(form)
  }

  return (
    <div className="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
      <div className="mb-4 flex items-center justify-between">
        <h2 className="text-lg font-semibold text-gray-900">New Collection</h2>
        <button onClick={onCancel} className="text-gray-400 hover:text-gray-600">
          <X className="h-5 w-5" />
        </button>
      </div>
      <form onSubmit={handleSubmit} className="grid grid-cols-1 gap-4 sm:grid-cols-2">
        <div>
          <label className="block text-sm font-medium text-gray-700">Name</label>
          <input
            type="text"
            required
            value={form.name}
            onChange={(e) => handleNameChange(e.target.value)}
            className="mt-1 w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-gray-400 focus:outline-none focus:ring-1 focus:ring-gray-400"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700">Slug</label>
          <input
            type="text"
            required
            value={form.slug}
            onChange={(e) => setForm({ ...form, slug: e.target.value })}
            className="mt-1 w-full rounded-lg border border-gray-200 px-3 py-2 text-sm font-mono focus:border-gray-400 focus:outline-none focus:ring-1 focus:ring-gray-400"
          />
        </div>
        <div className="sm:col-span-2">
          <label className="block text-sm font-medium text-gray-700">Description</label>
          <textarea
            rows={3}
            value={form.description}
            onChange={(e) => setForm({ ...form, description: e.target.value })}
            className="mt-1 w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-gray-400 focus:outline-none focus:ring-1 focus:ring-gray-400"
            placeholder="Optional description"
          />
        </div>
        <div className="flex items-end gap-3 sm:col-span-2">
          <button
            type="submit"
            disabled={isSubmitting}
            className="rounded-lg bg-gray-900 px-4 py-2.5 text-sm font-medium text-white hover:bg-gray-800 transition-colors disabled:opacity-60"
          >
            {isSubmitting ? 'Creating…' : 'Create Collection'}
          </button>
          <button
            type="button"
            onClick={onCancel}
            className="rounded-lg border border-gray-200 bg-white px-4 py-2.5 text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors"
          >
            Cancel
          </button>
        </div>
      </form>
    </div>
  )
}

// ---------------------------------------------------------------------------
// Categories tab
// ---------------------------------------------------------------------------

function CategoriesTab() {
  const [showForm, setShowForm] = useState(false)
  const [searchQuery, setSearchQuery] = useState('')

  const { data, isLoading } = useCategories()
  const createCategory = useCreateCategory()
  const deleteCategory = useDeleteCategory()

  const categories: Category[] = data?.items ?? []

  // Build hierarchy: root categories first, then children indented underneath
  const rootCategories = categories.filter((c) => !c.parent_id)
  const childMap = new Map<string, Category[]>()
  for (const cat of categories) {
    if (cat.parent_id) {
      const children = childMap.get(cat.parent_id) ?? []
      children.push(cat)
      childMap.set(cat.parent_id, children)
    }
  }

  // Flatten into display order: root → its children → next root → …
  const ordered: Array<{ category: Category; depth: number }> = []
  for (const root of rootCategories) {
    ordered.push({ category: root, depth: 0 })
    for (const child of childMap.get(root.id) ?? []) {
      ordered.push({ category: child, depth: 1 })
    }
  }
  // Also include any orphaned children (parent not in current page)
  for (const cat of categories) {
    if (cat.parent_id && !categories.find((c) => c.id === cat.parent_id)) {
      ordered.push({ category: cat, depth: 0 })
    }
  }

  const filtered = ordered.filter(({ category: c }) =>
    c.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    c.slug.toLowerCase().includes(searchQuery.toLowerCase()),
  )

  function handleSubmit(formData: CategoryFormData) {
    createCategory.mutate(
      {
        name: formData.name,
        slug: formData.slug,
        parent_id: formData.parent_id || null,
        description: formData.description,
      },
      {
        onSuccess: () => setShowForm(false),
      },
    )
  }

  function handleDelete(id: string) {
    if (confirm('Are you sure you want to delete this category?')) {
      deleteCategory.mutate(id)
    }
  }

  return (
    <div className="space-y-4">
      {/* Tab header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-lg font-semibold text-gray-900">Categories</h2>
          <p className="mt-0.5 text-sm text-gray-500">
            Organise your products into a hierarchy
          </p>
        </div>
        <button
          onClick={() => setShowForm(true)}
          className="inline-flex items-center gap-2 rounded-lg bg-gray-900 px-4 py-2.5 text-sm font-medium text-white shadow-sm hover:bg-gray-800 transition-colors"
        >
          <Plus className="h-4 w-4" />
          Add Category
        </button>
      </div>

      {/* Search */}
      <div className="relative">
        <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
        <input
          type="text"
          placeholder="Search categories…"
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="w-full rounded-lg border border-gray-200 bg-white py-2.5 pl-10 pr-4 text-sm text-gray-900 placeholder:text-gray-400 focus:border-gray-400 focus:outline-none focus:ring-1 focus:ring-gray-400"
        />
      </div>

      {/* Inline form */}
      {showForm && (
        <CategoryForm
          categories={categories}
          onSubmit={handleSubmit}
          onCancel={() => setShowForm(false)}
          isSubmitting={createCategory.isPending}
        />
      )}

      {/* Table */}
      <div className="rounded-xl border border-gray-200 bg-white shadow-sm">
        {isLoading ? (
          <div className="flex items-center justify-center py-12">
            <div className="h-6 w-6 animate-spin rounded-full border-2 border-gray-300 border-t-gray-900" />
          </div>
        ) : filtered.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-16 text-gray-400">
            <FolderTree className="mb-3 h-10 w-10" />
            <p className="text-sm font-medium">No categories found</p>
            <p className="mt-1 text-xs">Create your first category to get started.</p>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b border-gray-100 text-left text-sm text-gray-500">
                  <th className="px-6 py-3 font-medium">Name</th>
                  <th className="px-6 py-3 font-medium">Slug</th>
                  <th className="px-6 py-3 font-medium">Description</th>
                  <th className="px-6 py-3 font-medium">Created</th>
                  <th className="px-6 py-3 font-medium text-right">Actions</th>
                </tr>
              </thead>
              <tbody>
                {filtered.map(({ category, depth }) => (
                  <tr key={category.id} className="border-b border-gray-50 last:border-0">
                    <td className="px-6 py-4">
                      <div
                        className="flex items-center gap-2"
                        style={{ paddingLeft: depth * 20 }}
                      >
                        {depth > 0 && (
                          <span className="text-gray-300 text-xs">└</span>
                        )}
                        <span className="text-sm font-medium text-gray-900">{category.name}</span>
                      </div>
                    </td>
                    <td className="px-6 py-4 text-sm font-mono text-gray-500">{category.slug}</td>
                    <td className="px-6 py-4 text-sm text-gray-500">
                      <span className="line-clamp-1 max-w-xs">
                        {category.description || <span className="text-gray-300">—</span>}
                      </span>
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-500">
                      {formatDate(category.created_at)}
                    </td>
                    <td className="px-6 py-4 text-right">
                      <button
                        onClick={() => handleDelete(category.id)}
                        className="rounded-md p-1.5 text-gray-400 hover:bg-red-50 hover:text-red-600 transition-colors"
                        title="Delete"
                      >
                        <Trash2 className="h-4 w-4" />
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  )
}

// ---------------------------------------------------------------------------
// Collections tab
// ---------------------------------------------------------------------------

function CollectionsTab() {
  const [showForm, setShowForm] = useState(false)
  const [searchQuery, setSearchQuery] = useState('')

  const { data, isLoading } = useCollections()
  const createCollection = useCreateCollection()
  const deleteCollection = useDeleteCollection()

  const collections: Collection[] = data?.items ?? []

  const filtered = collections.filter(
    (c) =>
      c.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      c.slug.toLowerCase().includes(searchQuery.toLowerCase()),
  )

  function handleSubmit(formData: CollectionFormData) {
    createCollection.mutate(
      {
        name: formData.name,
        slug: formData.slug,
        description: formData.description,
        product_ids: [],
        is_automatic: false,
        rules: {},
      },
      {
        onSuccess: () => setShowForm(false),
      },
    )
  }

  function handleDelete(id: string) {
    if (confirm('Are you sure you want to delete this collection?')) {
      deleteCollection.mutate(id)
    }
  }

  return (
    <div className="space-y-4">
      {/* Tab header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-lg font-semibold text-gray-900">Collections</h2>
          <p className="mt-0.5 text-sm text-gray-500">
            Group products into curated or automatic collections
          </p>
        </div>
        <button
          onClick={() => setShowForm(true)}
          className="inline-flex items-center gap-2 rounded-lg bg-gray-900 px-4 py-2.5 text-sm font-medium text-white shadow-sm hover:bg-gray-800 transition-colors"
        >
          <Plus className="h-4 w-4" />
          Add Collection
        </button>
      </div>

      {/* Search */}
      <div className="relative">
        <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
        <input
          type="text"
          placeholder="Search collections…"
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="w-full rounded-lg border border-gray-200 bg-white py-2.5 pl-10 pr-4 text-sm text-gray-900 placeholder:text-gray-400 focus:border-gray-400 focus:outline-none focus:ring-1 focus:ring-gray-400"
        />
      </div>

      {/* Inline form */}
      {showForm && (
        <CollectionForm
          onSubmit={handleSubmit}
          onCancel={() => setShowForm(false)}
          isSubmitting={createCollection.isPending}
        />
      )}

      {/* Table */}
      <div className="rounded-xl border border-gray-200 bg-white shadow-sm">
        {isLoading ? (
          <div className="flex items-center justify-center py-12">
            <div className="h-6 w-6 animate-spin rounded-full border-2 border-gray-300 border-t-gray-900" />
          </div>
        ) : filtered.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-16 text-gray-400">
            <Layers className="mb-3 h-10 w-10" />
            <p className="text-sm font-medium">No collections found</p>
            <p className="mt-1 text-xs">Create your first collection to get started.</p>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b border-gray-100 text-left text-sm text-gray-500">
                  <th className="px-6 py-3 font-medium">Name</th>
                  <th className="px-6 py-3 font-medium">Slug</th>
                  <th className="px-6 py-3 font-medium">Products</th>
                  <th className="px-6 py-3 font-medium">Type</th>
                  <th className="px-6 py-3 font-medium">Created</th>
                  <th className="px-6 py-3 font-medium text-right">Actions</th>
                </tr>
              </thead>
              <tbody>
                {filtered.map((collection) => (
                  <tr key={collection.id} className="border-b border-gray-50 last:border-0">
                    <td className="px-6 py-4">
                      <div className="text-sm font-medium text-gray-900">{collection.name}</div>
                      {collection.description && (
                        <div className="mt-0.5 text-xs text-gray-500 line-clamp-1 max-w-xs">
                          {collection.description}
                        </div>
                      )}
                    </td>
                    <td className="px-6 py-4 text-sm font-mono text-gray-500">{collection.slug}</td>
                    <td className="px-6 py-4 text-sm text-gray-600">
                      {collection.product_ids.length}
                    </td>
                    <td className="px-6 py-4">
                      {collection.is_automatic ? (
                        <span className="inline-flex rounded-full px-2.5 py-0.5 text-xs font-medium bg-indigo-100 text-indigo-800">
                          Automatic
                        </span>
                      ) : (
                        <span className="inline-flex rounded-full px-2.5 py-0.5 text-xs font-medium bg-gray-100 text-gray-700">
                          Manual
                        </span>
                      )}
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-500">
                      {formatDate(collection.created_at)}
                    </td>
                    <td className="px-6 py-4 text-right">
                      <button
                        onClick={() => handleDelete(collection.id)}
                        className="rounded-md p-1.5 text-gray-400 hover:bg-red-50 hover:text-red-600 transition-colors"
                        title="Delete"
                      >
                        <Trash2 className="h-4 w-4" />
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  )
}

// ---------------------------------------------------------------------------
// Catalog page (tabs)
// ---------------------------------------------------------------------------

type Tab = 'categories' | 'collections'

export default function Catalog() {
  const [activeTab, setActiveTab] = useState<Tab>('categories')

  return (
    <div className="space-y-6">
      {/* Page header */}
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Catalog</h1>
        <p className="mt-1 text-sm text-gray-500">
          Manage categories and collections for your store
        </p>
      </div>

      {/* Tabs */}
      <div className="flex gap-1 rounded-lg border border-gray-200 bg-gray-50 p-1 w-fit">
        <button
          onClick={() => setActiveTab('categories')}
          className={`inline-flex items-center gap-2 rounded-md px-4 py-2 text-sm font-medium transition-colors ${
            activeTab === 'categories'
              ? 'bg-white text-gray-900 shadow-sm'
              : 'text-gray-500 hover:text-gray-700'
          }`}
        >
          <FolderTree className="h-4 w-4" />
          Categories
        </button>
        <button
          onClick={() => setActiveTab('collections')}
          className={`inline-flex items-center gap-2 rounded-md px-4 py-2 text-sm font-medium transition-colors ${
            activeTab === 'collections'
              ? 'bg-white text-gray-900 shadow-sm'
              : 'text-gray-500 hover:text-gray-700'
          }`}
        >
          <Layers className="h-4 w-4" />
          Collections
        </button>
      </div>

      {/* Tab content */}
      {activeTab === 'categories' ? <CategoriesTab /> : <CollectionsTab />}
    </div>
  )
}
