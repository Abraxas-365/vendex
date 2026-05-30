import { useState, useEffect } from 'react'
import { useNavigate, useParams } from '@tanstack/react-router'
import {
  ArrowLeft,
  ChevronUp,
  ChevronDown,
  Pencil,
  Trash2,
  Plus,
  Save,
  Globe,
  Loader2,
  Layers,
} from 'lucide-react'
import type { Section, BlockType, PageStatus, PageMeta } from '../../types'
import {
  useBlockTypes,
  usePage,
  useCreatePage,
  useUpdatePage,
  usePublishPage,
} from '../../lib/hooks'

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

function toTitleCase(str: string): string {
  return str.replace(/_/g, ' ').replace(/\b\w/g, (c) => c.toUpperCase())
}

function getSectionPreview(section: Section, blockType?: BlockType): string {
  const s = section.settings
  if (!blockType) return ''
  const name = blockType.name

  if (name.includes('hero')) {
    const heading = s.heading ?? s.title ?? ''
    return heading ? `"${String(heading)}"` : 'Hero section'
  }
  if (name.includes('rich') || name.includes('text') || name.includes('content')) {
    const content = s.content ?? s.text ?? ''
    const str = String(content)
    return str.length > 60 ? str.slice(0, 60) + '…' : str || 'Rich text section'
  }
  if (name.includes('product') || name.includes('grid')) {
    const cols = s.columns ?? s.num_columns ?? 4
    const count = s.num_products ?? s.limit ?? 8
    return `${String(cols)} columns, ${String(count)} products`
  }
  if (name.includes('banner')) {
    return String(s.text ?? s.message ?? 'Banner section')
  }
  if (name.includes('image') || name.includes('media')) {
    return String(s.alt ?? s.caption ?? 'Image section')
  }

  // Fallback: first string property
  for (const val of Object.values(s)) {
    if (typeof val === 'string' && val.length > 0) {
      return val.length > 60 ? val.slice(0, 60) + '…' : val
    }
  }
  return blockType.display_name
}

const categoryColors: Record<string, string> = {
  content: 'bg-blue-50 text-blue-700 border-blue-200',
  commerce: 'bg-green-50 text-green-700 border-green-200',
  media: 'bg-purple-50 text-purple-700 border-purple-200',
  layout: 'bg-orange-50 text-orange-700 border-orange-200',
}

const statusColors: Record<PageStatus, string> = {
  draft: 'bg-gray-100 text-gray-700',
  pending_review: 'bg-yellow-100 text-yellow-700',
  published: 'bg-green-100 text-green-700',
  archived: 'bg-red-100 text-red-700',
}

const inputClass =
  'w-full rounded-lg border border-gray-300 px-3 py-2 text-sm text-gray-900 placeholder-gray-400 focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500'

// ---------------------------------------------------------------------------
// Dynamic settings form for a selected section
// ---------------------------------------------------------------------------

interface SectionSettingsFormProps {
  section: Section
  blockType: BlockType
  onChange: (settings: Record<string, unknown>) => void
}

function SectionSettingsForm({ section, blockType, onChange }: SectionSettingsFormProps) {
  const schema = blockType.schema
  const properties = schema?.properties ?? {}

  function handleChange(key: string, value: unknown) {
    onChange({ ...section.settings, [key]: value })
  }

  return (
    <div className="space-y-4">
      <div className="border-b border-gray-200 pb-3">
        <h3 className="text-sm font-semibold text-gray-900">{blockType.display_name}</h3>
      </div>
      {Object.entries(properties).map(([key, prop]) => {
        const p = prop as Record<string, unknown>
        const label = (p.title as string) ?? toTitleCase(key)
        const val = section.settings[key] ?? p.default ?? ''
        const isTextarea =
          p.type === 'string' &&
          (key.includes('content') || key.includes('description') || key.includes('text') || key.includes('html'))

        if (p.type === 'boolean') {
          return (
            <div key={key}>
              <label className="flex items-center gap-2">
                <input
                  type="checkbox"
                  checked={Boolean(val)}
                  onChange={(e) => handleChange(key, e.target.checked)}
                  className="h-4 w-4 rounded border-gray-300 text-indigo-600 focus:ring-indigo-500"
                />
                <span className="text-sm font-medium text-gray-700">{label}</span>
              </label>
            </div>
          )
        }

        if (p.type === 'integer' || p.type === 'number') {
          return (
            <div key={key}>
              <label className="mb-1 block text-sm font-medium text-gray-700">{label}</label>
              <input
                type="number"
                className={inputClass}
                value={typeof val === 'number' ? val : Number(val)}
                onChange={(e) =>
                  handleChange(
                    key,
                    p.type === 'integer' ? parseInt(e.target.value, 10) : parseFloat(e.target.value),
                  )
                }
              />
            </div>
          )
        }

        if (p.type === 'string' && Array.isArray(p.enum)) {
          return (
            <div key={key}>
              <label className="mb-1 block text-sm font-medium text-gray-700">{label}</label>
              <select
                className={inputClass}
                value={String(val)}
                onChange={(e) => handleChange(key, e.target.value)}
              >
                {(p.enum as string[]).map((opt: string) => (
                  <option key={opt} value={opt}>
                    {toTitleCase(opt)}
                  </option>
                ))}
              </select>
            </div>
          )
        }

        if (p.type === 'array' || p.type === 'object') {
          return (
            <div key={key}>
              <label className="mb-1 block text-sm font-medium text-gray-700">{label}</label>
              <textarea
                rows={4}
                className={`${inputClass} font-mono text-xs`}
                value={typeof val === 'string' ? val : JSON.stringify(val, null, 2)}
                onChange={(e) => {
                  try {
                    handleChange(key, JSON.parse(e.target.value))
                  } catch {
                    handleChange(key, e.target.value)
                  }
                }}
              />
              <p className="mt-1 text-xs text-gray-400">JSON format</p>
            </div>
          )
        }

        if (isTextarea) {
          return (
            <div key={key}>
              <label className="mb-1 block text-sm font-medium text-gray-700">{label}</label>
              <textarea
                rows={4}
                className={inputClass}
                value={String(val)}
                onChange={(e) => handleChange(key, e.target.value)}
              />
            </div>
          )
        }

        // Default: text input
        return (
          <div key={key}>
            <label className="mb-1 block text-sm font-medium text-gray-700">{label}</label>
            <input
              type="text"
              className={inputClass}
              value={String(val)}
              onChange={(e) => handleChange(key, e.target.value)}
            />
          </div>
        )
      })}

      {Object.keys(properties).length === 0 && (
        <p className="text-sm text-gray-400">No configurable settings for this block type.</p>
      )}
    </div>
  )
}

// ---------------------------------------------------------------------------
// Main PageEditor component
// ---------------------------------------------------------------------------

export default function PageEditor() {
  const navigate = useNavigate()
  // pageId is optional — only present on /admin/pages/$pageId/edit
  const params = useParams({ strict: false }) as { pageId?: string }
  const pageId = params.pageId

  const { data: blockTypesData, isLoading: blockTypesLoading } = useBlockTypes()
  const { data: existingPage, isLoading: pageLoading } = usePage(pageId ?? '')
  const createPage = useCreatePage()
  const updatePage = useUpdatePage()
  const publishPage = usePublishPage()

  // Page state
  const [title, setTitle] = useState('')
  const [slug, setSlug] = useState('')
  const [status, setStatus] = useState<PageStatus>('draft')
  const [meta, setMeta] = useState<PageMeta>({
    description: '',
    og_title: '',
    og_image: '',
    keywords: [],
  })
  const [sections, setSections] = useState<Section[]>([])
  const [selectedSectionId, setSelectedSectionId] = useState<string | null>(null)

  // Load existing page data
  useEffect(() => {
    if (existingPage) {
      setTitle(existingPage.title)
      setSlug(existingPage.slug)
      setStatus(existingPage.status)
      if (existingPage.meta) setMeta(existingPage.meta)
      // Load sections from page data if content_type is 'blocks'
      const rawSections = (existingPage as unknown as { sections?: Section[] }).sections
      if (rawSections && Array.isArray(rawSections)) {
        setSections(rawSections)
      }
    }
  }, [existingPage])

  const blockTypes = blockTypesData ?? []

  // Group by category
  const categories = ['content', 'commerce', 'media', 'layout'] as const
  const byCategory = categories.reduce<Record<string, BlockType[]>>((acc, cat) => {
    acc[cat] = blockTypes.filter((bt) => bt.category === cat)
    return acc
  }, {})

  const selectedSection = sections.find((s) => s.id === selectedSectionId) ?? null
  const blockTypeByName = Object.fromEntries(blockTypes.map((bt) => [bt.name, bt]))
  const selectedBlockType = selectedSection ? blockTypeByName[selectedSection.type] : undefined

  // Auto-generate slug from title (only for new pages)
  function handleTitleChange(newTitle: string) {
    setTitle(newTitle)
    if (!pageId) {
      const generated = newTitle
        .toLowerCase()
        .replace(/[^a-z0-9]+/g, '-')
        .replace(/^-|-$/g, '')
      setSlug(generated)
    }
  }

  function addSection(blockType: BlockType) {
    const newSection: Section = {
      id: crypto.randomUUID(),
      type: blockType.name,
      settings: { ...(blockType.default_settings ?? {}) },
      blocks: [],
    }
    setSections((prev) => [...prev, newSection])
    setSelectedSectionId(newSection.id)
  }

  function moveSection(index: number, direction: 'up' | 'down') {
    const newSections = [...sections]
    const swapIndex = direction === 'up' ? index - 1 : index + 1
    if (swapIndex < 0 || swapIndex >= newSections.length) return
    ;[newSections[index], newSections[swapIndex]] = [newSections[swapIndex], newSections[index]]
    setSections(newSections)
  }

  function deleteSection(id: string) {
    setSections((prev) => prev.filter((s) => s.id !== id))
    if (selectedSectionId === id) setSelectedSectionId(null)
  }

  function updateSectionSettings(id: string, settings: Record<string, unknown>) {
    setSections((prev) => prev.map((s) => (s.id === id ? { ...s, settings } : s)))
  }

  function buildPayload(targetStatus?: PageStatus) {
    return {
      title,
      slug,
      status: targetStatus ?? status,
      meta,
      content_type: 'blocks' as const,
      sections,
    }
  }

  function handleSave() {
    const payload = buildPayload()
    if (pageId) {
      updatePage.mutate(
        { id: pageId, data: payload },
        { onSuccess: () => void navigate({ to: '/admin/pages' }) },
      )
    } else {
      createPage.mutate(payload, {
        onSuccess: () => void navigate({ to: '/admin/pages' }),
      })
    }
  }

  function handlePublish() {
    if (pageId) {
      publishPage.mutate(pageId, {
        onSuccess: () => void navigate({ to: '/admin/pages' }),
      })
    } else {
      createPage.mutate(buildPayload('published'), {
        onSuccess: () => void navigate({ to: '/admin/pages' }),
      })
    }
  }

  const isLoading = blockTypesLoading || (!!pageId && pageLoading)
  const isSaving = createPage.isPending || updatePage.isPending || publishPage.isPending
  const saveError = createPage.error ?? updatePage.error ?? publishPage.error

  if (isLoading) {
    return (
      <div className="flex h-96 items-center justify-center">
        <Loader2 className="h-6 w-6 animate-spin text-gray-400" />
      </div>
    )
  }

  return (
    <div className="flex h-full flex-col overflow-hidden">
      {/* Top bar */}
      <div className="flex shrink-0 items-center justify-between border-b border-gray-200 bg-white px-4 py-3 shadow-sm">
        <div className="flex items-center gap-3">
          <button
            onClick={() => void navigate({ to: '/admin/pages' })}
            className="flex items-center gap-1.5 rounded-lg px-3 py-1.5 text-sm font-medium text-gray-600 hover:bg-gray-100 transition-colors"
          >
            <ArrowLeft className="h-4 w-4" />
            Pages
          </button>
          <span className="text-gray-300">|</span>
          <input
            className="rounded border border-transparent px-2 py-1 text-base font-semibold text-gray-900 hover:border-gray-300 focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
            value={title}
            onChange={(e) => handleTitleChange(e.target.value)}
            placeholder="Untitled Page"
          />
          <span
            className={`inline-flex rounded-full px-2.5 py-0.5 text-xs font-medium ${statusColors[status]}`}
          >
            {status.replace('_', ' ')}
          </span>
        </div>
        <div className="flex items-center gap-2">
          {saveError && (
            <span className="text-xs text-red-600">{saveError.message}</span>
          )}
          <button
            onClick={handleSave}
            disabled={isSaving}
            className="inline-flex items-center gap-1.5 rounded-lg border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-50 transition-colors"
          >
            {isSaving ? (
              <Loader2 className="h-4 w-4 animate-spin" />
            ) : (
              <Save className="h-4 w-4" />
            )}
            Save
          </button>
          <button
            onClick={handlePublish}
            disabled={isSaving}
            className="inline-flex items-center gap-1.5 rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700 disabled:opacity-50 transition-colors"
          >
            <Globe className="h-4 w-4" />
            Publish
          </button>
        </div>
      </div>

      {/* Three-column layout */}
      <div className="flex flex-1 overflow-hidden">
        {/* Left sidebar: Block palette */}
        <aside className="flex w-60 shrink-0 flex-col overflow-y-auto border-r border-gray-200 bg-white">
          <div className="border-b border-gray-200 px-4 py-3">
            <h2 className="text-xs font-semibold uppercase tracking-wider text-gray-500">
              Block Palette
            </h2>
          </div>
          {blockTypes.length === 0 ? (
            <div className="flex flex-1 flex-col items-center justify-center p-4 text-center">
              <Layers className="h-8 w-8 text-gray-300" />
              <p className="mt-2 text-xs text-gray-400">No block types available</p>
            </div>
          ) : (
            <div className="space-y-4 p-3">
              {categories.map((cat) => {
                const items = byCategory[cat] ?? []
                if (items.length === 0) return null
                return (
                  <div key={cat}>
                    <p className="mb-1.5 px-1 text-xs font-semibold uppercase tracking-wider text-gray-400">
                      {cat}
                    </p>
                    <div className="space-y-1">
                      {items.map((bt) => (
                        <button
                          key={bt.id}
                          onClick={() => addSection(bt)}
                          className={`flex w-full items-center gap-2.5 rounded-lg border px-3 py-2 text-left text-xs font-medium transition-colors hover:opacity-90 ${categoryColors[cat] ?? 'bg-gray-50 text-gray-700 border-gray-200'}`}
                        >
                          <span className="text-base leading-none">
                            {bt.icon || '📦'}
                          </span>
                          <span>{bt.display_name}</span>
                          <Plus className="ml-auto h-3 w-3 shrink-0 opacity-60" />
                        </button>
                      ))}
                    </div>
                  </div>
                )
              })}
            </div>
          )}
        </aside>

        {/* Center canvas */}
        <main className="flex-1 overflow-y-auto bg-slate-50 p-6">
          <div className="mx-auto max-w-3xl space-y-6">
            {/* Page settings card */}
            <div className="rounded-xl border border-gray-200 bg-white shadow-sm">
              <div className="border-b border-gray-200 px-6 py-4">
                <h2 className="text-sm font-semibold text-gray-900">Page Settings</h2>
              </div>
              <div className="space-y-4 px-6 py-5">
                <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
                  <div>
                    <label className="mb-1 block text-sm font-medium text-gray-700">Title</label>
                    <input
                      className={inputClass}
                      value={title}
                      onChange={(e) => handleTitleChange(e.target.value)}
                      placeholder="My Page"
                    />
                  </div>
                  <div>
                    <label className="mb-1 block text-sm font-medium text-gray-700">Slug</label>
                    <div className="flex items-center rounded-lg border border-gray-300 overflow-hidden focus-within:border-indigo-500 focus-within:ring-1 focus-within:ring-indigo-500">
                      <span className="bg-gray-50 px-3 py-2 text-sm text-gray-400 border-r border-gray-300">
                        /
                      </span>
                      <input
                        className="w-full px-3 py-2 text-sm text-gray-900 focus:outline-none"
                        value={slug}
                        onChange={(e) => setSlug(e.target.value)}
                        placeholder="my-page"
                      />
                    </div>
                  </div>
                </div>
                <div>
                  <label className="mb-1 block text-sm font-medium text-gray-700">
                    Meta Description
                  </label>
                  <textarea
                    rows={2}
                    className={inputClass}
                    value={meta.description}
                    onChange={(e) => setMeta({ ...meta, description: e.target.value })}
                    placeholder="Brief description for search engines..."
                  />
                </div>
                <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
                  <div>
                    <label className="mb-1 block text-sm font-medium text-gray-700">OG Title</label>
                    <input
                      className={inputClass}
                      value={meta.og_title ?? ''}
                      onChange={(e) => setMeta({ ...meta, og_title: e.target.value })}
                      placeholder="Open Graph title"
                    />
                  </div>
                  <div>
                    <label className="mb-1 block text-sm font-medium text-gray-700">OG Image URL</label>
                    <input
                      className={inputClass}
                      value={meta.og_image ?? ''}
                      onChange={(e) => setMeta({ ...meta, og_image: e.target.value })}
                      placeholder="https://..."
                    />
                  </div>
                </div>
              </div>
            </div>

            {/* Sections */}
            <div className="space-y-3">
              {sections.length === 0 ? (
                <div className="flex flex-col items-center justify-center rounded-xl border-2 border-dashed border-gray-300 bg-white py-16 text-center">
                  <Layers className="mb-3 h-10 w-10 text-gray-300" />
                  <p className="text-sm font-medium text-gray-500">No sections yet</p>
                  <p className="mt-1 text-xs text-gray-400">
                    Click a block in the left panel to add it here.
                  </p>
                </div>
              ) : (
                sections.map((section, index) => {
                  const bt = blockTypeByName[section.type]
                  const isSelected = selectedSectionId === section.id
                  const preview = getSectionPreview(section, bt)

                  return (
                    <div
                      key={section.id}
                      className={`rounded-xl border bg-white shadow-sm transition-all ${
                        isSelected
                          ? 'border-indigo-400 ring-2 ring-indigo-200'
                          : 'border-gray-200 hover:border-gray-300'
                      }`}
                    >
                      <div className="flex items-center gap-3 px-4 py-3">
                        {/* Drag handle area */}
                        <div className="flex-shrink-0 text-lg leading-none">
                          {bt?.icon ?? '📦'}
                        </div>
                        <div className="min-w-0 flex-1">
                          <p className="text-sm font-medium text-gray-900">
                            {bt?.display_name ?? 'Unknown Block'}
                          </p>
                          {preview && (
                            <p className="mt-0.5 truncate text-xs text-gray-400">{preview}</p>
                          )}
                        </div>
                        {/* Controls */}
                        <div className="flex items-center gap-1">
                          <button
                            onClick={() => moveSection(index, 'up')}
                            disabled={index === 0}
                            className="rounded p-1 text-gray-400 hover:bg-gray-100 hover:text-gray-600 disabled:opacity-30 transition-colors"
                            title="Move up"
                          >
                            <ChevronUp className="h-4 w-4" />
                          </button>
                          <button
                            onClick={() => moveSection(index, 'down')}
                            disabled={index === sections.length - 1}
                            className="rounded p-1 text-gray-400 hover:bg-gray-100 hover:text-gray-600 disabled:opacity-30 transition-colors"
                            title="Move down"
                          >
                            <ChevronDown className="h-4 w-4" />
                          </button>
                          <button
                            onClick={() =>
                              setSelectedSectionId(isSelected ? null : section.id)
                            }
                            className={`rounded p-1 transition-colors ${
                              isSelected
                                ? 'bg-indigo-100 text-indigo-600'
                                : 'text-gray-400 hover:bg-gray-100 hover:text-gray-600'
                            }`}
                            title="Edit settings"
                          >
                            <Pencil className="h-4 w-4" />
                          </button>
                          <button
                            onClick={() => deleteSection(section.id)}
                            className="rounded p-1 text-gray-400 hover:bg-red-50 hover:text-red-600 transition-colors"
                            title="Delete section"
                          >
                            <Trash2 className="h-4 w-4" />
                          </button>
                        </div>
                      </div>
                    </div>
                  )
                })
              )}
            </div>
          </div>
        </main>

        {/* Right sidebar: Settings panel */}
        <aside className="flex w-80 shrink-0 flex-col overflow-y-auto border-l border-gray-200 bg-white">
          <div className="border-b border-gray-200 px-4 py-3">
            <h2 className="text-xs font-semibold uppercase tracking-wider text-gray-500">
              Section Settings
            </h2>
          </div>
          <div className="flex-1 p-4">
            {selectedSection && selectedBlockType ? (
              <SectionSettingsForm
                section={selectedSection}
                blockType={selectedBlockType}
                onChange={(settings) => updateSectionSettings(selectedSection.id, settings)}
              />
            ) : (
              <div className="flex h-full flex-col items-center justify-center text-center">
                <Pencil className="mb-3 h-8 w-8 text-gray-300" />
                <p className="text-sm text-gray-400">
                  Select a section to edit its settings
                </p>
              </div>
            )}
          </div>
        </aside>
      </div>
    </div>
  )
}
