import { useState } from 'react'
import { Plus, Pencil, Trash2, X, Package, MinusCircle } from 'lucide-react'
import type { Bundle } from '../../types'
import {
  useBundles,
  useCreateBundle,
  useUpdateBundle,
  useDeleteBundle,
  useAddBundleItem,
  useRemoveBundleItem,
  useBundlePrice,
} from '../../lib/hooks'

interface BundleFormData {
  name: string
  description: string
  active: boolean
  discount_percent: number
}

const emptyBundleForm: BundleFormData = {
  name: '',
  description: '',
  active: true,
  discount_percent: 0,
}

interface AddItemFormData {
  product_id: string
  quantity: number
  discount_percent: number
}

const emptyItemForm: AddItemFormData = {
  product_id: '',
  quantity: 1,
  discount_percent: 0,
}

function BundlePriceInfo({ bundleId }: { bundleId: string }) {
  const { data: price } = useBundlePrice(bundleId)
  if (!price) return null
  return (
    <div className="rounded-lg bg-gray-50 px-3 py-2 text-xs text-gray-600">
      <span className="mr-3">Total: <strong>{(price.original_total.amount / 100).toFixed(2)} {price.original_total.currency}</strong></span>
      <span className="mr-3">Final: <strong className="text-green-700">{(price.final_price.amount / 100).toFixed(2)} {price.final_price.currency}</strong></span>
      <span>You save: <strong className="text-indigo-700">{price.savings_percent.toFixed(1)}%</strong></span>
    </div>
  )
}

function BundleItemsSection({ bundle }: { bundle: Bundle }) {
  const [showAddItem, setShowAddItem] = useState(false)
  const [itemForm, setItemForm] = useState<AddItemFormData>(emptyItemForm)

  const addItem = useAddBundleItem()
  const removeItem = useRemoveBundleItem()

  function handleAddItem(e: React.FormEvent) {
    e.preventDefault()
    addItem.mutate({
      bundleId: bundle.id,
      product_id: itemForm.product_id,
      quantity: itemForm.quantity,
      discount_percent: itemForm.discount_percent,
    })
    setShowAddItem(false)
    setItemForm(emptyItemForm)
  }

  return (
    <div className="border-t border-gray-50 bg-gray-50 px-4 py-3">
      <BundlePriceInfo bundleId={bundle.id} />
      <div className="mt-3 flex items-center justify-between">
        <p className="text-xs font-semibold uppercase tracking-wide text-gray-400">Items</p>
        <button
          onClick={() => setShowAddItem(!showAddItem)}
          className="inline-flex items-center gap-1 rounded-md border border-gray-200 px-2 py-1 text-xs font-medium text-gray-700 hover:bg-white transition-colors"
        >
          <Plus className="h-3 w-3" />
          Add Item
        </button>
      </div>
      {showAddItem && (
        <form onSubmit={handleAddItem} className="mt-2 flex items-end gap-2">
          <div>
            <label className="block text-xs text-gray-600">Product ID</label>
            <input type="text" required value={itemForm.product_id} onChange={(e) => setItemForm({ ...itemForm, product_id: e.target.value })}
              className="mt-0.5 w-48 rounded-lg border border-gray-200 bg-white px-2 py-1.5 text-xs focus:border-gray-400 focus:outline-none" placeholder="Product ID" />
          </div>
          <div>
            <label className="block text-xs text-gray-600">Qty</label>
            <input type="number" min="1" required value={itemForm.quantity} onChange={(e) => setItemForm({ ...itemForm, quantity: parseInt(e.target.value) })}
              className="mt-0.5 w-16 rounded-lg border border-gray-200 bg-white px-2 py-1.5 text-xs focus:border-gray-400 focus:outline-none" />
          </div>
          <div>
            <label className="block text-xs text-gray-600">Discount %</label>
            <input type="number" min="0" max="100" value={itemForm.discount_percent} onChange={(e) => setItemForm({ ...itemForm, discount_percent: parseFloat(e.target.value) })}
              className="mt-0.5 w-16 rounded-lg border border-gray-200 bg-white px-2 py-1.5 text-xs focus:border-gray-400 focus:outline-none" />
          </div>
          <button type="submit" className="rounded-lg bg-gray-900 px-2.5 py-1.5 text-xs font-medium text-white hover:bg-gray-800 transition-colors">Add</button>
          <button type="button" onClick={() => setShowAddItem(false)} className="rounded-lg border border-gray-200 bg-white px-2.5 py-1.5 text-xs font-medium text-gray-700 hover:bg-gray-50 transition-colors">Cancel</button>
        </form>
      )}
      {(bundle.items ?? []).length > 0 && (
        <div className="mt-2 space-y-1">
          {(bundle.items ?? []).map((item) => (
            <div key={item.id} className="flex items-center gap-2 rounded-md bg-white border border-gray-100 px-3 py-1.5 text-xs">
              <span className="flex-1 font-medium text-gray-900">{item.product_name ?? item.product_id}</span>
              <span className="text-gray-500">Qty: {item.quantity}</span>
              {item.discount_percent > 0 && <span className="text-green-700">-{item.discount_percent}%</span>}
              <button
                onClick={() => removeItem.mutate({ bundleId: bundle.id, itemId: item.id })}
                className="rounded-md p-0.5 text-red-400 hover:bg-red-50 hover:text-red-600 transition-colors"
              >
                <MinusCircle className="h-3.5 w-3.5" />
              </button>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}

export default function Bundles() {
  const [showForm, setShowForm] = useState(false)
  const [formData, setFormData] = useState<BundleFormData>(emptyBundleForm)
  const [editingId, setEditingId] = useState<string | null>(null)
  const [expandedId, setExpandedId] = useState<string | null>(null)
  const [page] = useState(1)

  const { data, isLoading } = useBundles({ page, page_size: 20 })
  const bundles: Bundle[] = data?.items ?? []

  const createBundle = useCreateBundle()
  const updateBundle = useUpdateBundle()
  const deleteBundle = useDeleteBundle()

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (editingId) {
      updateBundle.mutate({ id: editingId, data: formData })
    } else {
      createBundle.mutate(formData)
    }
    setShowForm(false)
    setFormData(emptyBundleForm)
    setEditingId(null)
  }

  function handleEdit(b: Bundle) {
    setFormData({ name: b.name, description: b.description, active: b.active, discount_percent: b.discount_percent })
    setEditingId(b.id)
    setShowForm(true)
  }

  function handleDelete(id: string) {
    if (confirm('Delete this bundle?')) deleteBundle.mutate(id)
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Product Bundles</h1>
          <p className="mt-1 text-sm text-gray-500">Grouped products with bundle pricing</p>
        </div>
        <button
          onClick={() => { setShowForm(true); setEditingId(null); setFormData(emptyBundleForm) }}
          className="inline-flex items-center gap-2 rounded-lg bg-gray-900 px-4 py-2.5 text-sm font-medium text-white shadow-sm hover:bg-gray-800 transition-colors"
        >
          <Plus className="h-4 w-4" />
          New Bundle
        </button>
      </div>

      {/* Form */}
      {showForm && (
        <div className="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
          <div className="mb-4 flex items-center justify-between">
            <h2 className="text-lg font-semibold text-gray-900">{editingId ? 'Edit Bundle' : 'New Bundle'}</h2>
            <button onClick={() => setShowForm(false)} className="text-gray-400 hover:text-gray-600">
              <X className="h-5 w-5" />
            </button>
          </div>
          <form onSubmit={handleSubmit} className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <div>
              <label className="block text-sm font-medium text-gray-700">Name</label>
              <input type="text" required value={formData.name} onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                className="mt-1 w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-gray-400 focus:outline-none" />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700">Bundle Discount %</label>
              <input type="number" min="0" max="100" value={formData.discount_percent} onChange={(e) => setFormData({ ...formData, discount_percent: parseFloat(e.target.value) })}
                className="mt-1 w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-gray-400 focus:outline-none" />
            </div>
            <div className="sm:col-span-2">
              <label className="block text-sm font-medium text-gray-700">Description</label>
              <input type="text" value={formData.description} onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                className="mt-1 w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-gray-400 focus:outline-none" />
            </div>
            <div className="flex items-center gap-2">
              <input type="checkbox" id="active" checked={formData.active} onChange={(e) => setFormData({ ...formData, active: e.target.checked })} className="rounded border-gray-300" />
              <label htmlFor="active" className="text-sm font-medium text-gray-700">Active</label>
            </div>
            <div className="flex items-end gap-3 sm:col-span-2">
              <button type="submit" className="rounded-lg bg-gray-900 px-4 py-2.5 text-sm font-medium text-white hover:bg-gray-800 transition-colors">
                {editingId ? 'Update Bundle' : 'Create Bundle'}
              </button>
              <button type="button" onClick={() => setShowForm(false)} className="rounded-lg border border-gray-200 bg-white px-4 py-2.5 text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors">
                Cancel
              </button>
            </div>
          </form>
        </div>
      )}

      {/* Bundles List */}
      <div className="rounded-xl border border-gray-200 bg-white shadow-sm">
        {isLoading ? (
          <div className="flex items-center justify-center py-12">
            <div className="h-6 w-6 animate-spin rounded-full border-2 border-gray-300 border-t-gray-900" />
          </div>
        ) : bundles.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-16 text-gray-400">
            <Package className="mb-3 h-10 w-10" />
            <p className="text-sm font-medium">No bundles yet</p>
          </div>
        ) : (
          <div className="divide-y divide-gray-50">
            {bundles.map((b) => (
              <div key={b.id}>
                <div className="flex items-center gap-4 px-6 py-4">
                  <div className="flex-1 min-w-0">
                    <div className="text-sm font-medium text-gray-900">{b.name}</div>
                    <div className="text-xs text-gray-500">{b.description}</div>
                    <div className="mt-0.5 text-xs text-gray-400">{(b.items ?? []).length} items · {b.discount_percent}% discount</div>
                  </div>
                  <span className={`inline-flex rounded-full px-2.5 py-0.5 text-xs font-medium ${b.active ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-600'}`}>
                    {b.active ? 'Active' : 'Draft'}
                  </span>
                  <div className="flex items-center gap-2">
                    <button onClick={() => setExpandedId(expandedId === b.id ? null : b.id)}
                      className="rounded-lg border border-gray-200 px-2.5 py-1 text-xs font-medium text-gray-700 hover:bg-gray-50 transition-colors">
                      {expandedId === b.id ? 'Collapse' : 'Items'}
                    </button>
                    <button onClick={() => handleEdit(b)} className="rounded-md p-1.5 text-gray-400 hover:bg-gray-100 hover:text-gray-600 transition-colors">
                      <Pencil className="h-4 w-4" />
                    </button>
                    <button onClick={() => handleDelete(b.id)} className="rounded-md p-1.5 text-gray-400 hover:bg-red-50 hover:text-red-600 transition-colors">
                      <Trash2 className="h-4 w-4" />
                    </button>
                  </div>
                </div>
                {expandedId === b.id && <BundleItemsSection bundle={b} />}
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
