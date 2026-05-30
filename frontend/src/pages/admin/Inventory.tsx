import { useState } from 'react'
import { Plus, Pencil, Trash2, X, Warehouse as WarehouseIcon, AlertTriangle, BarChart2 } from 'lucide-react'
import type { Warehouse } from '../../types'
import {
  useWarehouses,
  useCreateWarehouse,
  useUpdateWarehouse,
  useDeleteWarehouse,
  useLowStockAlerts,
  useStockLevels,
  useAdjustStock,
} from '../../lib/hooks'

type Tab = 'warehouses' | 'low-stock' | 'stock-levels'

interface WarehouseFormData {
  name: string
  address: string
  is_default: boolean
}

const emptyWarehouseForm: WarehouseFormData = {
  name: '',
  address: '',
  is_default: false,
}

interface AdjustFormData {
  warehouse_id: string
  quantity: number
  note: string
}

export default function Inventory() {
  const [activeTab, setActiveTab] = useState<Tab>('warehouses')
  const [showWarehouseForm, setShowWarehouseForm] = useState(false)
  const [warehouseForm, setWarehouseForm] = useState<WarehouseFormData>(emptyWarehouseForm)
  const [editingWarehouseId, setEditingWarehouseId] = useState<string | null>(null)
  const [showAdjustDialog, setShowAdjustDialog] = useState(false)
  const [adjustProductId, setAdjustProductId] = useState('')
  const [adjustProductName, setAdjustProductName] = useState('')
  const [adjustForm, setAdjustForm] = useState<AdjustFormData>({ warehouse_id: '', quantity: 0, note: '' })
  const [stockProductId, setStockProductId] = useState('')

  const { data: warehouses = [], isLoading: warehousesLoading } = useWarehouses()
  const { data: lowStockAlerts = [], isLoading: lowStockLoading } = useLowStockAlerts()
  const { data: stockLevels = [], isLoading: stockLoading } = useStockLevels(stockProductId)

  const createWarehouse = useCreateWarehouse()
  const updateWarehouse = useUpdateWarehouse()
  const deleteWarehouse = useDeleteWarehouse()
  const adjustStock = useAdjustStock()

  function handleWarehouseSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (editingWarehouseId) {
      updateWarehouse.mutate({ id: editingWarehouseId, data: warehouseForm })
    } else {
      createWarehouse.mutate(warehouseForm)
    }
    setShowWarehouseForm(false)
    setWarehouseForm(emptyWarehouseForm)
    setEditingWarehouseId(null)
  }

  function handleEditWarehouse(w: Warehouse) {
    setWarehouseForm({ name: w.name, address: w.address, is_default: w.is_default })
    setEditingWarehouseId(w.id)
    setShowWarehouseForm(true)
  }

  function handleDeleteWarehouse(id: string) {
    if (confirm('Delete this warehouse?')) deleteWarehouse.mutate(id)
  }

  function handleOpenAdjust(productId: string, productName: string) {
    setAdjustProductId(productId)
    setAdjustProductName(productName)
    setAdjustForm({ warehouse_id: warehouses[0]?.id ?? '', quantity: 0, note: '' })
    setShowAdjustDialog(true)
  }

  function handleAdjustSubmit(e: React.FormEvent) {
    e.preventDefault()
    adjustStock.mutate({
      product_id: adjustProductId,
      warehouse_id: adjustForm.warehouse_id,
      quantity: adjustForm.quantity,
      note: adjustForm.note,
    })
    setShowAdjustDialog(false)
  }

  const tabs: { key: Tab; label: string }[] = [
    { key: 'warehouses', label: 'Warehouses' },
    { key: 'low-stock', label: 'Low Stock Alerts' },
    { key: 'stock-levels', label: 'Stock Levels' },
  ]

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Inventory</h1>
          <p className="mt-1 text-sm text-gray-500">Manage warehouses, stock levels and adjustments</p>
        </div>
        {activeTab === 'warehouses' && (
          <button
            onClick={() => { setShowWarehouseForm(true); setEditingWarehouseId(null); setWarehouseForm(emptyWarehouseForm) }}
            className="inline-flex items-center gap-2 rounded-lg bg-gray-900 px-4 py-2.5 text-sm font-medium text-white shadow-sm hover:bg-gray-800 transition-colors"
          >
            <Plus className="h-4 w-4" />
            Add Warehouse
          </button>
        )}
      </div>

      {/* Tabs */}
      <div className="flex gap-1 border-b border-gray-200">
        {tabs.map((t) => (
          <button
            key={t.key}
            onClick={() => setActiveTab(t.key)}
            className={`px-4 py-2 text-sm font-medium border-b-2 transition-colors ${
              activeTab === t.key
                ? 'border-gray-900 text-gray-900'
                : 'border-transparent text-gray-500 hover:text-gray-700'
            }`}
          >
            {t.label}
          </button>
        ))}
      </div>

      {/* Warehouse Form */}
      {activeTab === 'warehouses' && showWarehouseForm && (
        <div className="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
          <div className="mb-4 flex items-center justify-between">
            <h2 className="text-lg font-semibold text-gray-900">
              {editingWarehouseId ? 'Edit Warehouse' : 'New Warehouse'}
            </h2>
            <button onClick={() => setShowWarehouseForm(false)} className="text-gray-400 hover:text-gray-600">
              <X className="h-5 w-5" />
            </button>
          </div>
          <form onSubmit={handleWarehouseSubmit} className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <div>
              <label className="block text-sm font-medium text-gray-700">Name</label>
              <input
                type="text"
                required
                value={warehouseForm.name}
                onChange={(e) => setWarehouseForm({ ...warehouseForm, name: e.target.value })}
                className="mt-1 w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-gray-400 focus:outline-none focus:ring-1 focus:ring-gray-400"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700">Address</label>
              <input
                type="text"
                value={warehouseForm.address}
                onChange={(e) => setWarehouseForm({ ...warehouseForm, address: e.target.value })}
                className="mt-1 w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-gray-400 focus:outline-none focus:ring-1 focus:ring-gray-400"
              />
            </div>
            <div className="flex items-center gap-2 sm:col-span-2">
              <input
                type="checkbox"
                id="is_default"
                checked={warehouseForm.is_default}
                onChange={(e) => setWarehouseForm({ ...warehouseForm, is_default: e.target.checked })}
                className="rounded border-gray-300"
              />
              <label htmlFor="is_default" className="text-sm font-medium text-gray-700">Default warehouse</label>
            </div>
            <div className="flex items-end gap-3 sm:col-span-2">
              <button
                type="submit"
                className="rounded-lg bg-gray-900 px-4 py-2.5 text-sm font-medium text-white hover:bg-gray-800 transition-colors"
              >
                {editingWarehouseId ? 'Update Warehouse' : 'Create Warehouse'}
              </button>
              <button
                type="button"
                onClick={() => setShowWarehouseForm(false)}
                className="rounded-lg border border-gray-200 bg-white px-4 py-2.5 text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors"
              >
                Cancel
              </button>
            </div>
          </form>
        </div>
      )}

      {/* Warehouses Tab */}
      {activeTab === 'warehouses' && (
        <div className="rounded-xl border border-gray-200 bg-white shadow-sm">
          {warehousesLoading ? (
            <div className="flex items-center justify-center py-12">
              <div className="h-6 w-6 animate-spin rounded-full border-2 border-gray-300 border-t-gray-900" />
            </div>
          ) : warehouses.length === 0 ? (
            <div className="flex flex-col items-center justify-center py-16 text-gray-400">
              <WarehouseIcon className="mb-3 h-10 w-10" />
              <p className="text-sm font-medium">No warehouses yet</p>
            </div>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead>
                  <tr className="border-b border-gray-100 text-left text-sm text-gray-500">
                    <th className="px-6 py-3 font-medium">Name</th>
                    <th className="px-6 py-3 font-medium">Address</th>
                    <th className="px-6 py-3 font-medium">Default</th>
                    <th className="px-6 py-3 font-medium text-right">Actions</th>
                  </tr>
                </thead>
                <tbody>
                  {warehouses.map((w) => (
                    <tr key={w.id} className="border-b border-gray-50 last:border-0">
                      <td className="px-6 py-4 text-sm font-medium text-gray-900">{w.name}</td>
                      <td className="px-6 py-4 text-sm text-gray-600">{w.address || '—'}</td>
                      <td className="px-6 py-4">
                        {w.is_default && (
                          <span className="inline-flex rounded-full px-2.5 py-0.5 text-xs font-medium bg-green-100 text-green-800">
                            Default
                          </span>
                        )}
                      </td>
                      <td className="px-6 py-4 text-right">
                        <div className="flex items-center justify-end gap-2">
                          <button
                            onClick={() => handleEditWarehouse(w)}
                            className="rounded-md p-1.5 text-gray-400 hover:bg-gray-100 hover:text-gray-600 transition-colors"
                          >
                            <Pencil className="h-4 w-4" />
                          </button>
                          <button
                            onClick={() => handleDeleteWarehouse(w.id)}
                            className="rounded-md p-1.5 text-gray-400 hover:bg-red-50 hover:text-red-600 transition-colors"
                          >
                            <Trash2 className="h-4 w-4" />
                          </button>
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </div>
      )}

      {/* Low Stock Tab */}
      {activeTab === 'low-stock' && (
        <div className="rounded-xl border border-gray-200 bg-white shadow-sm">
          {lowStockLoading ? (
            <div className="flex items-center justify-center py-12">
              <div className="h-6 w-6 animate-spin rounded-full border-2 border-gray-300 border-t-gray-900" />
            </div>
          ) : lowStockAlerts.length === 0 ? (
            <div className="flex flex-col items-center justify-center py-16 text-gray-400">
              <AlertTriangle className="mb-3 h-10 w-10" />
              <p className="text-sm font-medium">No low stock alerts</p>
            </div>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead>
                  <tr className="border-b border-gray-100 text-left text-sm text-gray-500">
                    <th className="px-6 py-3 font-medium">Product</th>
                    <th className="px-6 py-3 font-medium">SKU</th>
                    <th className="px-6 py-3 font-medium">Warehouse</th>
                    <th className="px-6 py-3 font-medium">Qty</th>
                    <th className="px-6 py-3 font-medium">Threshold</th>
                    <th className="px-6 py-3 font-medium text-right">Actions</th>
                  </tr>
                </thead>
                <tbody>
                  {lowStockAlerts.map((alert, i) => (
                    <tr key={i} className="border-b border-gray-50 last:border-0">
                      <td className="px-6 py-4 text-sm font-medium text-gray-900">{alert.product_name}</td>
                      <td className="px-6 py-4 text-sm font-mono text-gray-600">{alert.sku}</td>
                      <td className="px-6 py-4 text-sm text-gray-600">{alert.warehouse_name}</td>
                      <td className="px-6 py-4">
                        <span className="text-sm font-semibold text-red-600">{alert.quantity}</span>
                      </td>
                      <td className="px-6 py-4 text-sm text-gray-600">{alert.threshold}</td>
                      <td className="px-6 py-4 text-right">
                        <button
                          onClick={() => handleOpenAdjust(alert.product_id, alert.product_name)}
                          className="rounded-lg border border-gray-200 px-3 py-1.5 text-xs font-medium text-gray-700 hover:bg-gray-50 transition-colors"
                        >
                          Adjust
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </div>
      )}

      {/* Stock Levels Tab */}
      {activeTab === 'stock-levels' && (
        <div className="space-y-4">
          <div className="flex items-center gap-3">
            <input
              type="text"
              placeholder="Enter product ID to view stock levels..."
              value={stockProductId}
              onChange={(e) => setStockProductId(e.target.value)}
              className="w-full max-w-sm rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-gray-400 focus:outline-none focus:ring-1 focus:ring-gray-400"
            />
          </div>
          {stockProductId && (
            <div className="rounded-xl border border-gray-200 bg-white shadow-sm">
              {stockLoading ? (
                <div className="flex items-center justify-center py-12">
                  <div className="h-6 w-6 animate-spin rounded-full border-2 border-gray-300 border-t-gray-900" />
                </div>
              ) : stockLevels.length === 0 ? (
                <div className="flex flex-col items-center justify-center py-16 text-gray-400">
                  <BarChart2 className="mb-3 h-10 w-10" />
                  <p className="text-sm font-medium">No stock data found</p>
                </div>
              ) : (
                <div className="overflow-x-auto">
                  <table className="w-full">
                    <thead>
                      <tr className="border-b border-gray-100 text-left text-sm text-gray-500">
                        <th className="px-6 py-3 font-medium">Warehouse</th>
                        <th className="px-6 py-3 font-medium">Quantity</th>
                        <th className="px-6 py-3 font-medium">Reserved</th>
                        <th className="px-6 py-3 font-medium">Available</th>
                        <th className="px-6 py-3 font-medium text-right">Actions</th>
                      </tr>
                    </thead>
                    <tbody>
                      {stockLevels.map((s, i) => (
                        <tr key={i} className="border-b border-gray-50 last:border-0">
                          <td className="px-6 py-4 text-sm font-medium text-gray-900">{s.warehouse_name}</td>
                          <td className="px-6 py-4 text-sm text-gray-600">{s.quantity}</td>
                          <td className="px-6 py-4 text-sm text-gray-600">{s.reserved}</td>
                          <td className="px-6 py-4 text-sm font-medium text-green-700">{s.available}</td>
                          <td className="px-6 py-4 text-right">
                            <button
                              onClick={() => handleOpenAdjust(s.product_id, s.product_id)}
                              className="rounded-lg border border-gray-200 px-3 py-1.5 text-xs font-medium text-gray-700 hover:bg-gray-50 transition-colors"
                            >
                              Adjust
                            </button>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              )}
            </div>
          )}
        </div>
      )}

      {/* Adjust Stock Dialog */}
      {showAdjustDialog && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40">
          <div className="w-full max-w-md rounded-xl bg-white p-6 shadow-xl">
            <div className="mb-4 flex items-center justify-between">
              <h2 className="text-lg font-semibold text-gray-900">Adjust Stock: {adjustProductName}</h2>
              <button onClick={() => setShowAdjustDialog(false)} className="text-gray-400 hover:text-gray-600">
                <X className="h-5 w-5" />
              </button>
            </div>
            <form onSubmit={handleAdjustSubmit} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700">Warehouse</label>
                <select
                  required
                  value={adjustForm.warehouse_id}
                  onChange={(e) => setAdjustForm({ ...adjustForm, warehouse_id: e.target.value })}
                  className="mt-1 w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-gray-400 focus:outline-none"
                >
                  <option value="">Select warehouse</option>
                  {warehouses.map((w) => (
                    <option key={w.id} value={w.id}>{w.name}</option>
                  ))}
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700">Quantity (positive = add, negative = remove)</label>
                <input
                  type="number"
                  required
                  value={adjustForm.quantity}
                  onChange={(e) => setAdjustForm({ ...adjustForm, quantity: parseInt(e.target.value) })}
                  className="mt-1 w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-gray-400 focus:outline-none"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700">Note</label>
                <input
                  type="text"
                  value={adjustForm.note}
                  onChange={(e) => setAdjustForm({ ...adjustForm, note: e.target.value })}
                  className="mt-1 w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-gray-400 focus:outline-none"
                />
              </div>
              <div className="flex gap-3">
                <button
                  type="submit"
                  className="flex-1 rounded-lg bg-gray-900 py-2.5 text-sm font-medium text-white hover:bg-gray-800 transition-colors"
                >
                  Apply Adjustment
                </button>
                <button
                  type="button"
                  onClick={() => setShowAdjustDialog(false)}
                  className="flex-1 rounded-lg border border-gray-200 bg-white py-2.5 text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors"
                >
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
