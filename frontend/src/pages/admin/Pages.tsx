import { useState, useRef, useCallback } from 'react'
import {
  Plus,
  FileText,
  Eye,
  Pencil,
  Trash2,
  Globe,
  Archive,
  ArrowLeft,
  Search,
  X,
} from 'lucide-react'
import type { Page, PageStatus, PageMeta } from '../../types'
import {
  usePages,
  useCreatePage,
  useUpdatePage,
  useDeletePage,
} from '../../lib/hooks'

// ---------------------------------------------------------------------------
// Status badge colors
// ---------------------------------------------------------------------------

const statusColors: Record<PageStatus, string> = {
  draft: 'bg-gray-100 text-gray-800',
  pending_review: 'bg-yellow-100 text-yellow-800',
  published: 'bg-green-100 text-green-800',
  archived: 'bg-red-100 text-red-800',
}

const statusLabels: Record<PageStatus, string> = {
  draft: 'Draft',
  pending_review: 'Pending Review',
  published: 'Published',
  archived: 'Archived',
}

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  })
}

// ---------------------------------------------------------------------------
// Page Editor Component
// ---------------------------------------------------------------------------

interface PageEditorProps {
  page?: Page
  onSave: (data: PageEditorData) => void
  onCancel: () => void
}

interface PageEditorData {
  title: string
  slug: string
  html: string
  css: string
  meta: PageMeta
  status: PageStatus
}

function PageEditor({ page, onSave, onCancel }: PageEditorProps) {
  const [activeTab, setActiveTab] = useState<'html' | 'css'>('html')
  const [showPreview, setShowPreview] = useState(false)
  const iframeRef = useRef<HTMLIFrameElement>(null)

  const [form, setForm] = useState<PageEditorData>({
    title: page?.title ?? '',
    slug: page?.slug ?? '',
    html: page?.html ?? '',
    css: page?.css ?? '',
    meta: {
      description: page?.meta.description ?? '',
      og_title: page?.meta.og_title ?? '',
      og_image: page?.meta.og_image ?? '',
      keywords: page?.meta.keywords ?? [],
    },
    status: page?.status ?? 'draft',
  })

  const [keywordsInput, setKeywordsInput] = useState(
    page?.meta.keywords?.join(', ') ?? '',
  )

  // Auto-generate slug from title
  function handleTitleChange(title: string) {
    const slug = title
      .toLowerCase()
      .replace(/[^a-z0-9]+/g, '-')
      .replace(/^-|-$/g, '')
    setForm((f) => ({ ...f, title, slug: page ? f.slug : slug }))
  }

  function handleKeywordsChange(value: string) {
    setKeywordsInput(value)
    setForm((f) => ({
      ...f,
      meta: {
        ...f.meta,
        keywords: value.split(',').map((k) => k.trim()).filter(Boolean),
      },
    }))
  }

  const renderPreview = useCallback(() => {
    if (!iframeRef.current) return
    const doc = iframeRef.current.contentDocument
    if (!doc) return
    doc.open()
    doc.write(`
      <!DOCTYPE html>
      <html>
        <head>
          <meta charset="utf-8" />
          <meta name="viewport" content="width=device-width, initial-scale=1" />
          <style>${form.css}</style>
        </head>
        <body>${form.html}</body>
      </html>
    `)
    doc.close()
  }, [form.html, form.css])

  function handlePreviewToggle() {
    const next = !showPreview
    setShowPreview(next)
    if (next) {
      // Wait for iframe to mount
      requestAnimationFrame(() => renderPreview())
    }
  }

  function handleSave(status: PageStatus) {
    onSave({ ...form, status })
  }

  const isPendingReview = page?.status === 'pending_review'

  return (
    <div className="space-y-6">
      {/* Editor Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <button
            onClick={onCancel}
            className="rounded-lg p-2 text-gray-400 hover:bg-gray-100 hover:text-gray-600 transition-colors"
          >
            <ArrowLeft className="h-5 w-5" />
          </button>
          <div>
            <h2 className="text-xl font-bold text-gray-900">
              {page ? 'Edit Page' : 'Create Page'}
            </h2>
            {page && (
              <span
                className={`mt-1 inline-flex rounded-full px-2 py-0.5 text-xs font-medium ${statusColors[page.status]}`}
              >
                {statusLabels[page.status]}
              </span>
            )}
          </div>
        </div>
        <div className="flex items-center gap-3">
          <button
            onClick={handlePreviewToggle}
            className="inline-flex items-center gap-2 rounded-lg border border-gray-200 bg-white px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors"
          >
            <Eye className="h-4 w-4" />
            {showPreview ? 'Close Preview' : 'Preview'}
          </button>
          <button
            onClick={() => handleSave('draft')}
            className="rounded-lg border border-gray-200 bg-white px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors"
          >
            Save Draft
          </button>
          <button
            onClick={() => handleSave('pending_review')}
            className="rounded-lg bg-yellow-500 px-4 py-2 text-sm font-medium text-white hover:bg-yellow-600 transition-colors"
          >
            Submit for Review
          </button>
          {isPendingReview && (
            <button
              onClick={() => handleSave('published')}
              className="rounded-lg bg-green-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-green-700 transition-colors"
            >
              Publish
            </button>
          )}
        </div>
      </div>

      {/* Preview iframe */}
      {showPreview && (
        <div className="rounded-xl border border-gray-200 bg-white shadow-sm overflow-hidden">
          <div className="flex items-center justify-between border-b border-gray-200 bg-gray-50 px-4 py-2">
            <span className="text-xs font-medium text-gray-500 uppercase tracking-wider">
              Live Preview
            </span>
            <button
              onClick={() => setShowPreview(false)}
              className="text-gray-400 hover:text-gray-600"
            >
              <X className="h-4 w-4" />
            </button>
          </div>
          <iframe
            ref={iframeRef}
            title="Page preview"
            className="h-[500px] w-full border-0"
            sandbox="allow-same-origin"
            onLoad={renderPreview}
          />
        </div>
      )}

      {/* Title & Slug */}
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
        <div>
          <label className="block text-sm font-medium text-gray-700">Title</label>
          <input
            type="text"
            required
            value={form.title}
            onChange={(e) => handleTitleChange(e.target.value)}
            placeholder="My Awesome Page"
            className="mt-1 w-full rounded-lg border border-gray-200 px-3 py-2.5 text-sm focus:border-gray-400 focus:outline-none focus:ring-1 focus:ring-gray-400"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700">Slug</label>
          <div className="mt-1 flex items-center rounded-lg border border-gray-200 overflow-hidden">
            <span className="bg-gray-50 px-3 py-2.5 text-sm text-gray-400 border-r border-gray-200">
              /
            </span>
            <input
              type="text"
              required
              value={form.slug}
              onChange={(e) => setForm({ ...form, slug: e.target.value })}
              placeholder="my-awesome-page"
              className="w-full px-3 py-2.5 text-sm focus:outline-none"
            />
          </div>
        </div>
      </div>

      {/* Code Editor Tabs */}
      <div className="rounded-xl border border-gray-200 bg-white shadow-sm overflow-hidden">
        <div className="flex border-b border-gray-200">
          <button
            onClick={() => setActiveTab('html')}
            className={`px-6 py-3 text-sm font-medium transition-colors ${
              activeTab === 'html'
                ? 'border-b-2 border-gray-900 text-gray-900 bg-white'
                : 'text-gray-500 hover:text-gray-700 bg-gray-50'
            }`}
          >
            HTML
          </button>
          <button
            onClick={() => setActiveTab('css')}
            className={`px-6 py-3 text-sm font-medium transition-colors ${
              activeTab === 'css'
                ? 'border-b-2 border-gray-900 text-gray-900 bg-white'
                : 'text-gray-500 hover:text-gray-700 bg-gray-50'
            }`}
          >
            CSS
          </button>
        </div>
        <div className="relative">
          {activeTab === 'html' ? (
            <textarea
              value={form.html}
              onChange={(e) => setForm({ ...form, html: e.target.value })}
              placeholder="<div class='hero'>&#10;  <h1>Welcome</h1>&#10;  <p>Your content here...</p>&#10;</div>"
              spellCheck={false}
              className="h-80 w-full resize-y bg-gray-950 p-4 font-mono text-sm text-green-400 placeholder:text-gray-600 focus:outline-none"
            />
          ) : (
            <textarea
              value={form.css}
              onChange={(e) => setForm({ ...form, css: e.target.value })}
              placeholder=".hero {&#10;  text-align: center;&#10;  padding: 4rem 2rem;&#10;}"
              spellCheck={false}
              className="h-80 w-full resize-y bg-gray-950 p-4 font-mono text-sm text-blue-400 placeholder:text-gray-600 focus:outline-none"
            />
          )}
        </div>
      </div>

      {/* SEO Section */}
      <div className="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
        <h3 className="mb-4 text-sm font-semibold uppercase tracking-wider text-gray-500">
          SEO &amp; Social
        </h3>
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <div className="sm:col-span-2">
            <label className="block text-sm font-medium text-gray-700">Meta Description</label>
            <textarea
              rows={2}
              value={form.meta.description}
              onChange={(e) =>
                setForm({ ...form, meta: { ...form.meta, description: e.target.value } })
              }
              placeholder="A brief description of this page for search engines..."
              className="mt-1 w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-gray-400 focus:outline-none focus:ring-1 focus:ring-gray-400"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700">OG Title</label>
            <input
              type="text"
              value={form.meta.og_title}
              onChange={(e) =>
                setForm({ ...form, meta: { ...form.meta, og_title: e.target.value } })
              }
              placeholder="Open Graph title"
              className="mt-1 w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-gray-400 focus:outline-none focus:ring-1 focus:ring-gray-400"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700">OG Image URL</label>
            <input
              type="url"
              value={form.meta.og_image}
              onChange={(e) =>
                setForm({ ...form, meta: { ...form.meta, og_image: e.target.value } })
              }
              placeholder="https://example.com/image.jpg"
              className="mt-1 w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-gray-400 focus:outline-none focus:ring-1 focus:ring-gray-400"
            />
          </div>
          <div className="sm:col-span-2">
            <label className="block text-sm font-medium text-gray-700">
              Keywords (comma-separated)
            </label>
            <input
              type="text"
              value={keywordsInput}
              onChange={(e) => handleKeywordsChange(e.target.value)}
              placeholder="e-commerce, landing page, sale"
              className="mt-1 w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-gray-400 focus:outline-none focus:ring-1 focus:ring-gray-400"
            />
          </div>
        </div>
      </div>
    </div>
  )
}

// ---------------------------------------------------------------------------
// Pages List (main component)
// ---------------------------------------------------------------------------

export default function Pages() {
  const [view, setView] = useState<'list' | 'editor'>('list')
  const [editingPage, setEditingPage] = useState<Page | undefined>(undefined)
  const [searchQuery, setSearchQuery] = useState('')
  const [statusFilter, setStatusFilter] = useState<PageStatus | 'all'>('all')

  const { data, isLoading } = usePages()
  const createPage = useCreatePage()
  const updatePage = useUpdatePage()
  const deletePage = useDeletePage()

  const pages: Page[] = data?.items ?? []

  const filtered = pages.filter((p) => {
    const matchesStatus = statusFilter === 'all' || p.status === statusFilter
    const matchesSearch =
      searchQuery === '' ||
      p.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
      p.slug.toLowerCase().includes(searchQuery.toLowerCase())
    return matchesStatus && matchesSearch
  })

  function handleCreate() {
    setEditingPage(undefined)
    setView('editor')
  }

  function handleEdit(page: Page) {
    setEditingPage(page)
    setView('editor')
  }

  function handleSave(data: PageEditorData) {
    if (editingPage) {
      updatePage.mutate({ id: editingPage.id, data })
    } else {
      createPage.mutate(data)
    }
    setView('list')
    setEditingPage(undefined)
  }

  function handlePublish(page: Page) {
    updatePage.mutate({ id: page.id, data: { status: 'published' } })
  }

  function handleUnpublish(page: Page) {
    updatePage.mutate({ id: page.id, data: { status: 'draft' } })
  }

  function handleArchive(page: Page) {
    updatePage.mutate({ id: page.id, data: { status: 'archived' } })
  }

  function handleDelete(page: Page) {
    if (confirm(`Are you sure you want to delete "${page.title}"?`)) {
      deletePage.mutate(page.id)
    }
  }

  // ------ Editor view ------
  if (view === 'editor') {
    return (
      <PageEditor
        page={editingPage}
        onSave={handleSave}
        onCancel={() => {
          setView('list')
          setEditingPage(undefined)
        }}
      />
    )
  }

  // ------ List view ------
  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Pages</h1>
          <p className="mt-1 text-sm text-gray-500">
            Manage storefront pages &mdash; create, review, and publish content
          </p>
        </div>
        <button
          onClick={handleCreate}
          className="inline-flex items-center gap-2 rounded-lg bg-gray-900 px-4 py-2.5 text-sm font-medium text-white shadow-sm hover:bg-gray-800 transition-colors"
        >
          <Plus className="h-4 w-4" />
          Create Page
        </button>
      </div>

      {/* Filters */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center">
        <div className="relative flex-1">
          <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
          <input
            type="text"
            placeholder="Search by title or slug..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="w-full rounded-lg border border-gray-200 bg-white py-2.5 pl-10 pr-4 text-sm text-gray-900 placeholder:text-gray-400 focus:border-gray-400 focus:outline-none focus:ring-1 focus:ring-gray-400"
          />
        </div>
        <select
          value={statusFilter}
          onChange={(e) => setStatusFilter(e.target.value as PageStatus | 'all')}
          className="rounded-lg border border-gray-200 bg-white px-3 py-2.5 text-sm text-gray-700 focus:border-gray-400 focus:outline-none focus:ring-1 focus:ring-gray-400"
        >
          <option value="all">All Statuses</option>
          <option value="draft">Draft</option>
          <option value="pending_review">Pending Review</option>
          <option value="published">Published</option>
          <option value="archived">Archived</option>
        </select>
      </div>

      {/* Pages Table */}
      <div className="rounded-xl border border-gray-200 bg-white shadow-sm">
        {isLoading ? (
          <div className="flex items-center justify-center py-12">
            <div className="h-6 w-6 animate-spin rounded-full border-2 border-gray-300 border-t-gray-900" />
          </div>
        ) : filtered.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-16 text-gray-400">
            <FileText className="mb-3 h-10 w-10" />
            <p className="text-sm font-medium">No pages found</p>
            <p className="mt-1 text-xs">Create your first page to get started.</p>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b border-gray-100 text-left text-sm text-gray-500">
                  <th className="px-6 py-3 font-medium">Title</th>
                  <th className="px-6 py-3 font-medium">Slug</th>
                  <th className="px-6 py-3 font-medium">Status</th>
                  <th className="px-6 py-3 font-medium">Created By</th>
                  <th className="px-6 py-3 font-medium">Last Updated</th>
                  <th className="px-6 py-3 font-medium text-right">Actions</th>
                </tr>
              </thead>
              <tbody>
                {filtered.map((page) => (
                  <tr key={page.id} className="border-b border-gray-50 last:border-0">
                    <td className="px-6 py-4">
                      <span className="text-sm font-medium text-gray-900">{page.title}</span>
                    </td>
                    <td className="px-6 py-4 text-sm font-mono text-gray-500">/{page.slug}</td>
                    <td className="px-6 py-4">
                      <span
                        className={`inline-flex rounded-full px-2.5 py-0.5 text-xs font-medium ${statusColors[page.status]}`}
                      >
                        {statusLabels[page.status]}
                      </span>
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-600">{page.created_by}</td>
                    <td className="px-6 py-4 text-sm text-gray-500">
                      {formatDate(page.updated_at)}
                    </td>
                    <td className="px-6 py-4">
                      <div className="flex items-center justify-end gap-1">
                        {/* Preview — open editor in preview mode (future) */}
                        <button
                          onClick={() => handleEdit(page)}
                          className="rounded-md p-1.5 text-gray-400 hover:bg-gray-100 hover:text-gray-600 transition-colors"
                          title="Preview"
                        >
                          <Eye className="h-4 w-4" />
                        </button>
                        {/* Edit */}
                        <button
                          onClick={() => handleEdit(page)}
                          className="rounded-md p-1.5 text-gray-400 hover:bg-gray-100 hover:text-gray-600 transition-colors"
                          title="Edit"
                        >
                          <Pencil className="h-4 w-4" />
                        </button>
                        {/* Publish / Unpublish */}
                        {page.status === 'published' ? (
                          <button
                            onClick={() => handleUnpublish(page)}
                            className="rounded-md p-1.5 text-gray-400 hover:bg-yellow-50 hover:text-yellow-600 transition-colors"
                            title="Unpublish"
                          >
                            <Globe className="h-4 w-4" />
                          </button>
                        ) : (
                          <button
                            onClick={() => handlePublish(page)}
                            className="rounded-md p-1.5 text-gray-400 hover:bg-green-50 hover:text-green-600 transition-colors"
                            title="Publish"
                          >
                            <Globe className="h-4 w-4" />
                          </button>
                        )}
                        {/* Archive */}
                        <button
                          onClick={() => handleArchive(page)}
                          className="rounded-md p-1.5 text-gray-400 hover:bg-orange-50 hover:text-orange-600 transition-colors"
                          title="Archive"
                        >
                          <Archive className="h-4 w-4" />
                        </button>
                        {/* Delete */}
                        <button
                          onClick={() => handleDelete(page)}
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
          </div>
        )}
      </div>
    </div>
  )
}
