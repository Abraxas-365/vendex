import { useState } from 'react'
import { Store, Plus, Trash2, Edit2, Check, Globe } from 'lucide-react'
import type { Storefront } from '../../types'
import {
  useStorefronts,
  useCreateStorefront,
  useUpdateStorefront,
  useDeleteStorefront,
  useSetDefaultStorefront,
} from '../../lib/hooks'

const emptyForm: Partial<Storefront> = {
  name: '',
  slug: '',
  domain: '',
  description: '',
  logo_url: '',
  favicon_url: '',
  theme: '',
  default_locale: 'en',
  default_currency: 'USD',
  is_active: true,
}

export default function Multistores() {
  const [page] = useState(1)
  const [showDialog, setShowDialog] = useState(false)
  const [editing, setEditing] = useState<Storefront | null>(null)
  const [form, setForm] = useState<Partial<Storefront>>(emptyForm)

  const { data, isLoading } = useStorefronts({ page, page_size: 20 })
  const storefronts = data?.items ?? []

  const create = useCreateStorefront()
  const update = useUpdateStorefront()
  const remove = useDeleteStorefront()
  const setDefault = useSetDefaultStorefront()

  function openCreate() {
    setEditing(null)
    setForm(emptyForm)
    setShowDialog(true)
  }

  function openEdit(sf: Storefront) {
    setEditing(sf)
    setForm({
      name: sf.name,
      slug: sf.slug,
      domain: sf.domain,
      description: sf.description,
      logo_url: sf.logo_url,
      favicon_url: sf.favicon_url,
      theme: sf.theme,
      default_locale: sf.default_locale,
      default_currency: sf.default_currency,
      is_active: sf.is_active,
    })
    setShowDialog(true)
  }

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (editing) {
      update.mutate(
        { id: editing.id, data: form },
        { onSuccess: () => setShowDialog(false) },
      )
    } else {
      create.mutate(form, { onSuccess: () => setShowDialog(false) })
    }
  }

  const isPending = create.isPending || update.isPending

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <Store className="h-6 w-6 text-indigo-600" />
          <div>
            <h1 className="text-2xl font-bold text-gray-900">Storefronts</h1>
            <p className="mt-0.5 text-sm text-gray-500">Manage multiple storefront identities</p>
          </div>
        </div>
        <button
          onClick={openCreate}
          className="inline-flex items-center gap-2 rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700 transition-colors"
        >
          <Plus className="h-4 w-4" />
          New Storefront
        </button>
      </div>

      {/* Table */}
      <div className="rounded-xl border border-gray-200 bg-white shadow-sm overflow-hidden">
        {isLoading ? (
          <div className="flex items-center justify-center py-12">
            <div className="h-6 w-6 animate-spin rounded-full border-2 border-gray-300 border-t-indigo-600" />
          </div>
        ) : storefronts.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-16 text-gray-400">
            <Store className="mb-3 h-10 w-10" />
            <p className="text-sm font-medium">No storefronts yet</p>
            <p className="mt-1 text-xs">Create your first storefront to get started</p>
          </div>
        ) : (
          <table className="w-full text-sm">
            <thead className="border-b border-gray-100 bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left font-medium text-gray-500">Name</th>
                <th className="px-6 py-3 text-left font-medium text-gray-500">Slug</th>
                <th className="px-6 py-3 text-left font-medium text-gray-500">Domain</th>
                <th className="px-6 py-3 text-left font-medium text-gray-500">Locale / Currency</th>
                <th className="px-6 py-3 text-left font-medium text-gray-500">Status</th>
                <th className="px-6 py-3 text-right font-medium text-gray-500">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-50">
              {storefronts.map((sf) => (
                <tr key={sf.id} className="hover:bg-gray-50/50 transition-colors">
                  <td className="px-6 py-4">
                    <div className="flex items-center gap-2">
                      <span className="font-medium text-gray-900">{sf.name}</span>
                      {sf.is_default && (
                        <span className="inline-flex items-center rounded-full bg-indigo-100 px-2 py-0.5 text-xs font-medium text-indigo-700">
                          Default
                        </span>
                      )}
                    </div>
                  </td>
                  <td className="px-6 py-4 text-gray-600">{sf.slug}</td>
                  <td className="px-6 py-4">
                    {sf.domain ? (
                      <a
                        href={`https://${sf.domain}`}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="flex items-center gap-1 text-indigo-600 hover:underline"
                      >
                        <Globe className="h-3 w-3" />
                        {sf.domain}
                      </a>
                    ) : (
                      <span className="text-gray-400">—</span>
                    )}
                  </td>
                  <td className="px-6 py-4 text-gray-600">
                    {sf.default_locale} / {sf.default_currency}
                  </td>
                  <td className="px-6 py-4">
                    <span
                      className={`inline-flex rounded-full px-2 py-0.5 text-xs font-medium ${
                        sf.is_active
                          ? 'bg-green-100 text-green-700'
                          : 'bg-gray-100 text-gray-600'
                      }`}
                    >
                      {sf.is_active ? 'Active' : 'Inactive'}
                    </span>
                  </td>
                  <td className="px-6 py-4">
                    <div className="flex items-center justify-end gap-1">
                      {!sf.is_default && (
                        <button
                          onClick={() => setDefault.mutate(sf.id)}
                          className="rounded-md p-1.5 text-gray-400 hover:bg-indigo-50 hover:text-indigo-600 transition-colors"
                          title="Set as default"
                        >
                          <Check className="h-4 w-4" />
                        </button>
                      )}
                      <button
                        onClick={() => openEdit(sf)}
                        className="rounded-md p-1.5 text-gray-400 hover:bg-blue-50 hover:text-blue-600 transition-colors"
                        title="Edit"
                      >
                        <Edit2 className="h-4 w-4" />
                      </button>
                      <button
                        onClick={() => remove.mutate(sf.id)}
                        className="rounded-md p-1.5 text-gray-400 hover:bg-red-50 hover:text-red-600 transition-colors"
                        title="Delete"
                      >
                        <Trash2 className="h-4 w-4" />
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
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4">
          <div className="w-full max-w-lg rounded-2xl bg-white shadow-2xl">
            <div className="border-b border-gray-100 px-6 py-4">
              <h2 className="text-lg font-semibold text-gray-900">
                {editing ? 'Edit Storefront' : 'New Storefront'}
              </h2>
            </div>
            <form onSubmit={handleSubmit} className="space-y-4 px-6 py-5">
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-xs font-medium text-gray-700 mb-1">Name *</label>
                  <input
                    required
                    value={form.name ?? ''}
                    onChange={(e) => setForm((f) => ({ ...f, name: e.target.value }))}
                    className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
                  />
                </div>
                <div>
                  <label className="block text-xs font-medium text-gray-700 mb-1">Slug *</label>
                  <input
                    required
                    value={form.slug ?? ''}
                    onChange={(e) => setForm((f) => ({ ...f, slug: e.target.value }))}
                    className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
                  />
                </div>
              </div>
              <div>
                <label className="block text-xs font-medium text-gray-700 mb-1">Domain</label>
                <input
                  value={form.domain ?? ''}
                  onChange={(e) => setForm((f) => ({ ...f, domain: e.target.value }))}
                  placeholder="example.com"
                  className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
                />
              </div>
              <div>
                <label className="block text-xs font-medium text-gray-700 mb-1">Description</label>
                <textarea
                  rows={2}
                  value={form.description ?? ''}
                  onChange={(e) => setForm((f) => ({ ...f, description: e.target.value }))}
                  className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
                />
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-xs font-medium text-gray-700 mb-1">Default Locale</label>
                  <input
                    value={form.default_locale ?? ''}
                    onChange={(e) => setForm((f) => ({ ...f, default_locale: e.target.value }))}
                    placeholder="en"
                    className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
                  />
                </div>
                <div>
                  <label className="block text-xs font-medium text-gray-700 mb-1">Default Currency</label>
                  <input
                    value={form.default_currency ?? ''}
                    onChange={(e) => setForm((f) => ({ ...f, default_currency: e.target.value }))}
                    placeholder="USD"
                    className="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
                  />
                </div>
              </div>
              <div className="flex items-center gap-2">
                <input
                  type="checkbox"
                  id="is_active"
                  checked={form.is_active ?? true}
                  onChange={(e) => setForm((f) => ({ ...f, is_active: e.target.checked }))}
                  className="rounded border-gray-300"
                />
                <label htmlFor="is_active" className="text-sm text-gray-700">Active</label>
              </div>
              {(create.error ?? update.error) && (
                <p className="text-sm text-red-600">{(create.error ?? update.error)?.message}</p>
              )}
              <div className="flex justify-end gap-3 pt-2">
                <button
                  type="button"
                  onClick={() => setShowDialog(false)}
                  className="rounded-lg border border-gray-200 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={isPending}
                  className="rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700 disabled:opacity-50 transition-colors"
                >
                  {isPending ? 'Saving…' : editing ? 'Save Changes' : 'Create'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}
