import { useState, useRef } from 'react'
import {
  ArrowLeft,
  Plus,
  Pencil,
  Trash2,
  Copy,
  Globe,
  CheckCircle2,
  Loader2,
  Palette,
  ChevronDown,
  ChevronRight,
  Save,
} from 'lucide-react'
import type { Theme, ThemeTokens, ThemeColors, ThemeTypography, ThemeSpacing, ThemeBorders, ThemeShadows } from '../../types'
import {
  useThemes,
  useCreateTheme,
  useUpdateTheme,
  useActivateTheme,
  useDuplicateTheme,
  useDeleteTheme,
} from '../../lib/hooks'

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

const inputClass =
  'w-full rounded-lg border border-gray-300 px-3 py-2 text-sm text-gray-900 placeholder-gray-400 focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500'

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  })
}

const DEFAULT_TOKENS: ThemeTokens = {
  colors: {
    primary: '#6366f1',
    secondary: '#f3f4f6',
    background: '#ffffff',
    surface: '#f9fafb',
    text: '#111827',
    text_muted: '#6b7280',
    border: '#e5e7eb',
    success: '#22c55e',
    error: '#ef4444',
    warning: '#f59e0b',
    info: '#3b82f6',
  },
  typography: {
    font_heading: 'Inter, sans-serif',
    font_body: 'Inter, sans-serif',
    base_size: '16px',
    scale_ratio: 1.25,
  },
  spacing: {
    unit: '4px',
    section_padding: '80px',
  },
  borders: {
    radius_sm: '4px',
    radius_md: '8px',
    radius_lg: '12px',
    radius_full: '9999px',
  },
  shadows: {
    sm: '0 1px 2px 0 rgb(0 0 0 / 0.05)',
    md: '0 4px 6px -1px rgb(0 0 0 / 0.1), 0 2px 4px -2px rgb(0 0 0 / 0.1)',
    lg: '0 10px 15px -3px rgb(0 0 0 / 0.1), 0 4px 6px -4px rgb(0 0 0 / 0.1)',
  },
}

// ---------------------------------------------------------------------------
// SectionCard (collapsible)
// ---------------------------------------------------------------------------

function SectionCard({
  title,
  children,
  defaultOpen = true,
}: {
  title: string
  children: React.ReactNode
  defaultOpen?: boolean
}) {
  const [open, setOpen] = useState(defaultOpen)
  return (
    <div className="rounded-xl border border-gray-200 bg-white shadow-sm">
      <button
        type="button"
        onClick={() => setOpen(!open)}
        className="flex w-full items-center justify-between border-b border-gray-200 px-6 py-4 text-left"
      >
        <h2 className="text-base font-semibold text-gray-900">{title}</h2>
        {open ? (
          <ChevronDown className="h-4 w-4 text-gray-400" />
        ) : (
          <ChevronRight className="h-4 w-4 text-gray-400" />
        )}
      </button>
      {open && <div className="space-y-4 px-6 py-5">{children}</div>}
    </div>
  )
}

function Field({ label, children }: { label: string; children: React.ReactNode }) {
  return (
    <div>
      <label className="mb-1 block text-sm font-medium text-gray-700">{label}</label>
      {children}
    </div>
  )
}

// ---------------------------------------------------------------------------
// Color swatch input
// ---------------------------------------------------------------------------

interface ColorInputProps {
  label: string
  value: string
  onChange: (v: string) => void
}

function ColorInput({ label, value, onChange }: ColorInputProps) {
  const inputRef = useRef<HTMLInputElement>(null)

  return (
    <div className="flex items-center justify-between gap-3">
      <span className="min-w-0 flex-1 truncate text-sm text-gray-700">{label}</span>
      <div className="flex items-center gap-2">
        <button
          type="button"
          onClick={() => inputRef.current?.click()}
          className="h-6 w-6 shrink-0 rounded border border-gray-300 shadow-sm hover:ring-2 hover:ring-indigo-400 transition-all"
          style={{ backgroundColor: value }}
          title="Click to pick color"
        />
        <input
          type="color"
          ref={inputRef}
          value={value}
          onChange={(e) => onChange(e.target.value)}
          className="sr-only"
        />
        <input
          type="text"
          value={value}
          onChange={(e) => onChange(e.target.value)}
          className="w-24 rounded border border-gray-300 px-2 py-1 font-mono text-xs text-gray-900 focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
          placeholder="#000000"
        />
      </div>
    </div>
  )
}

// ---------------------------------------------------------------------------
// Live Preview
// ---------------------------------------------------------------------------

interface LivePreviewProps {
  tokens: ThemeTokens
}

function LivePreview({ tokens }: LivePreviewProps) {
  const { colors, typography, borders, shadows } = tokens

  return (
    <div
      className="rounded-xl border border-gray-200 overflow-hidden shadow-sm"
      style={{ fontFamily: typography.font_body, backgroundColor: colors.background, color: colors.text }}
    >
      <div
        className="border-b px-4 py-3"
        style={{ borderColor: colors.border, backgroundColor: colors.surface }}
      >
        <p className="text-xs font-semibold uppercase tracking-wider text-gray-400">
          Live Preview
        </p>
      </div>
      <div className="space-y-4 p-5">
        {/* Headings */}
        <div>
          <p
            style={{
              fontFamily: typography.font_heading,
              fontSize: '1.375rem',
              fontWeight: 700,
              color: colors.text,
              lineHeight: 1.3,
            }}
          >
            Heading Text
          </p>
          <p style={{ fontSize: '0.875rem', color: colors.text_muted, marginTop: '4px' }}>
            Body text in {typography.font_body.split(',')[0]}
          </p>
        </div>

        {/* Buttons */}
        <div className="flex flex-wrap gap-2">
          <button
            style={{
              backgroundColor: colors.primary,
              color: '#fff',
              borderRadius: borders.radius_md,
              padding: '8px 16px',
              fontSize: '0.875rem',
              fontWeight: 600,
              border: 'none',
              boxShadow: shadows.sm,
              cursor: 'default',
            }}
          >
            Primary Button
          </button>
          <button
            style={{
              backgroundColor: 'transparent',
              color: colors.primary,
              borderRadius: borders.radius_md,
              padding: '7px 15px',
              fontSize: '0.875rem',
              fontWeight: 600,
              border: `1.5px solid ${colors.primary}`,
              cursor: 'default',
            }}
          >
            Secondary
          </button>
        </div>

        {/* Surface card */}
        <div
          style={{
            backgroundColor: colors.surface,
            border: `1px solid ${colors.border}`,
            borderRadius: borders.radius_lg,
            padding: '12px',
            boxShadow: shadows.md,
          }}
        >
          <p style={{ fontSize: '0.875rem', fontWeight: 600, color: colors.text }}>
            Surface Card
          </p>
          <p style={{ fontSize: '0.75rem', color: colors.text_muted, marginTop: '4px' }}>
            Card content with muted text
          </p>
        </div>

        {/* Status badges */}
        <div className="flex flex-wrap gap-2">
          <span
            style={{
              backgroundColor: colors.success + '20',
              color: colors.success,
              borderRadius: borders.radius_full,
              padding: '2px 10px',
              fontSize: '0.75rem',
              fontWeight: 600,
            }}
          >
            ✓ Success
          </span>
          <span
            style={{
              backgroundColor: colors.error + '20',
              color: colors.error,
              borderRadius: borders.radius_full,
              padding: '2px 10px',
              fontSize: '0.75rem',
              fontWeight: 600,
            }}
          >
            ✕ Error
          </span>
          <span
            style={{
              backgroundColor: colors.warning + '20',
              color: colors.warning,
              borderRadius: borders.radius_full,
              padding: '2px 10px',
              fontSize: '0.75rem',
              fontWeight: 600,
            }}
          >
            ⚠ Warning
          </span>
        </div>

        {/* Mono sample */}
        <code
          style={{
            fontFamily: 'monospace',
            fontSize: '0.75rem',
            color: colors.text_muted,
            backgroundColor: colors.surface,
            border: `1px solid ${colors.border}`,
            borderRadius: borders.radius_sm,
            padding: '2px 6px',
          }}
        >
          monospace sample
        </code>
      </div>
    </div>
  )
}

// ---------------------------------------------------------------------------
// Theme editor form (token editing)
// ---------------------------------------------------------------------------

interface ThemeEditorFormProps {
  theme: Theme
  onBack: () => void
}

function ThemeEditorForm({ theme, onBack }: ThemeEditorFormProps) {
  const updateTheme = useUpdateTheme()
  const activateTheme = useActivateTheme()

  const [name, setName] = useState(theme.name)
  const [tokens, setTokens] = useState<ThemeTokens>(theme.tokens ?? DEFAULT_TOKENS)
  const [saved, setSaved] = useState(false)

  // Helpers
  function setColors(patch: Partial<ThemeColors>) {
    setTokens((t) => ({ ...t, colors: { ...t.colors, ...patch } }))
  }
  function setTypography(patch: Partial<ThemeTypography>) {
    setTokens((t) => ({ ...t, typography: { ...t.typography, ...patch } }))
  }
  function setSpacing(patch: Partial<ThemeSpacing>) {
    setTokens((t) => ({ ...t, spacing: { ...t.spacing, ...patch } }))
  }
  function setBorders(patch: Partial<ThemeBorders>) {
    setTokens((t) => ({ ...t, borders: { ...t.borders, ...patch } }))
  }
  function setShadows(patch: Partial<ThemeShadows>) {
    setTokens((t) => ({ ...t, shadows: { ...t.shadows, ...patch } }))
  }

  function handleSave() {
    updateTheme.mutate(
      { id: theme.id, data: { name, tokens } },
      {
        onSuccess: () => {
          setSaved(true)
          setTimeout(() => setSaved(false), 2000)
        },
      },
    )
  }

  function handleActivate() {
    activateTheme.mutate(theme.id)
  }

  const isSaving = updateTheme.isPending || activateTheme.isPending

  const colorLabels: Record<keyof ThemeColors, string> = {
    primary: 'Primary',
    secondary: 'Secondary',
    background: 'Background',
    surface: 'Surface',
    text: 'Text',
    text_muted: 'Text Muted',
    border: 'Border',
    success: 'Success',
    error: 'Error',
    warning: 'Warning',
    info: 'Info',
  }

  return (
    <div className="space-y-6">
      {/* Top bar */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <button
            onClick={onBack}
            className="flex items-center gap-1.5 rounded-lg px-3 py-1.5 text-sm font-medium text-gray-600 hover:bg-gray-100 transition-colors"
          >
            <ArrowLeft className="h-4 w-4" />
            Themes
          </button>
          <span className="text-gray-300">|</span>
          <input
            className="rounded border border-transparent px-2 py-1 text-base font-semibold text-gray-900 hover:border-gray-300 focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
            value={name}
            onChange={(e) => setName(e.target.value)}
          />
          {theme.is_active && (
            <span className="inline-flex rounded-full bg-green-100 px-2.5 py-0.5 text-xs font-medium text-green-700">
              Active
            </span>
          )}
        </div>
        <div className="flex items-center gap-2">
          {updateTheme.isError && (
            <span className="text-xs text-red-600">{updateTheme.error.message}</span>
          )}
          {!theme.is_active && (
            <button
              onClick={handleActivate}
              disabled={isSaving}
              className="inline-flex items-center gap-1.5 rounded-lg border border-green-300 bg-green-50 px-4 py-2 text-sm font-medium text-green-700 hover:bg-green-100 disabled:opacity-50 transition-colors"
            >
              <Globe className="h-4 w-4" />
              Activate
            </button>
          )}
          <button
            onClick={handleSave}
            disabled={isSaving}
            className="inline-flex items-center gap-1.5 rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700 disabled:opacity-50 transition-colors"
          >
            {isSaving ? (
              <Loader2 className="h-4 w-4 animate-spin" />
            ) : saved ? (
              <CheckCircle2 className="h-4 w-4" />
            ) : (
              <Save className="h-4 w-4" />
            )}
            {saved ? 'Saved!' : 'Save'}
          </button>
        </div>
      </div>

      {/* Two-panel layout */}
      <div className="flex gap-6 items-start">
        {/* Left: token forms (60%) */}
        <div className="w-3/5 space-y-4">
          {/* Colors */}
          <SectionCard title="Colors">
            <div className="space-y-3">
              {(Object.keys(tokens.colors) as (keyof ThemeColors)[]).map((key) => (
                <ColorInput
                  key={key}
                  label={colorLabels[key]}
                  value={tokens.colors[key]}
                  onChange={(v) => setColors({ [key]: v })}
                />
              ))}
            </div>
          </SectionCard>

          {/* Typography */}
          <SectionCard title="Typography" defaultOpen={false}>
            <div className="grid grid-cols-1 gap-4">
              <Field label="Heading Font">
                <input
                  className={inputClass}
                  value={tokens.typography.font_heading}
                  onChange={(e) => setTypography({ font_heading: e.target.value })}
                  placeholder="Inter, sans-serif"
                />
              </Field>
              <Field label="Body Font">
                <input
                  className={inputClass}
                  value={tokens.typography.font_body}
                  onChange={(e) => setTypography({ font_body: e.target.value })}
                  placeholder="Inter, sans-serif"
                />
              </Field>
              <div className="grid grid-cols-2 gap-4">
                <Field label="Base Size">
                  <input
                    className={inputClass}
                    value={tokens.typography.base_size}
                    onChange={(e) => setTypography({ base_size: e.target.value })}
                    placeholder="16px"
                  />
                </Field>
                <Field label="Scale Ratio">
                  <input
                    type="number"
                    step="0.05"
                    min="1"
                    max="2"
                    className={inputClass}
                    value={tokens.typography.scale_ratio}
                    onChange={(e) => setTypography({ scale_ratio: parseFloat(e.target.value) })}
                  />
                </Field>
              </div>
            </div>
          </SectionCard>

          {/* Spacing */}
          <SectionCard title="Spacing" defaultOpen={false}>
            <div className="grid grid-cols-2 gap-4">
              <Field label="Base Unit">
                <input
                  className={inputClass}
                  value={tokens.spacing.unit}
                  onChange={(e) => setSpacing({ unit: e.target.value })}
                  placeholder="4px"
                />
              </Field>
              <Field label="Section Padding">
                <input
                  className={inputClass}
                  value={tokens.spacing.section_padding}
                  onChange={(e) => setSpacing({ section_padding: e.target.value })}
                  placeholder="80px"
                />
              </Field>
            </div>
          </SectionCard>

          {/* Borders */}
          <SectionCard title="Borders" defaultOpen={false}>
            <div className="grid grid-cols-2 gap-4">
              <Field label="Radius SM">
                <input
                  className={inputClass}
                  value={tokens.borders.radius_sm}
                  onChange={(e) => setBorders({ radius_sm: e.target.value })}
                  placeholder="4px"
                />
              </Field>
              <Field label="Radius MD">
                <input
                  className={inputClass}
                  value={tokens.borders.radius_md}
                  onChange={(e) => setBorders({ radius_md: e.target.value })}
                  placeholder="8px"
                />
              </Field>
              <Field label="Radius LG">
                <input
                  className={inputClass}
                  value={tokens.borders.radius_lg}
                  onChange={(e) => setBorders({ radius_lg: e.target.value })}
                  placeholder="12px"
                />
              </Field>
              <Field label="Radius Full">
                <input
                  className={inputClass}
                  value={tokens.borders.radius_full}
                  onChange={(e) => setBorders({ radius_full: e.target.value })}
                  placeholder="9999px"
                />
              </Field>
            </div>
          </SectionCard>

          {/* Shadows */}
          <SectionCard title="Shadows" defaultOpen={false}>
            <div className="space-y-4">
              <Field label="Shadow SM">
                <textarea
                  rows={2}
                  className={`${inputClass} font-mono text-xs`}
                  value={tokens.shadows.sm}
                  onChange={(e) => setShadows({ sm: e.target.value })}
                />
              </Field>
              <Field label="Shadow MD">
                <textarea
                  rows={2}
                  className={`${inputClass} font-mono text-xs`}
                  value={tokens.shadows.md}
                  onChange={(e) => setShadows({ md: e.target.value })}
                />
              </Field>
              <Field label="Shadow LG">
                <textarea
                  rows={2}
                  className={`${inputClass} font-mono text-xs`}
                  value={tokens.shadows.lg}
                  onChange={(e) => setShadows({ lg: e.target.value })}
                />
              </Field>
            </div>
          </SectionCard>
        </div>

        {/* Right: live preview (40%) */}
        <div className="w-2/5 sticky top-0">
          <LivePreview tokens={tokens} />
        </div>
      </div>
    </div>
  )
}

// ---------------------------------------------------------------------------
// Theme list (default view)
// ---------------------------------------------------------------------------

interface ThemeListProps {
  themes: Theme[]
  onEdit: (theme: Theme) => void
}

function ThemeList({ themes, onEdit }: ThemeListProps) {
  const createTheme = useCreateTheme()
  const activateTheme = useActivateTheme()
  const duplicateTheme = useDuplicateTheme()
  const deleteTheme = useDeleteTheme()

  function handleCreate() {
    createTheme.mutate(
      {
        name: 'New Theme',
        tokens: DEFAULT_TOKENS,
      },
      {
        onSuccess: (created) => onEdit(created),
      },
    )
  }

  function handleDelete(theme: Theme) {
    if (confirm(`Delete theme "${theme.name}"?`)) {
      deleteTheme.mutate(theme.id)
    }
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Themes</h1>
          <p className="mt-1 text-sm text-gray-500">
            Manage storefront themes and design tokens
          </p>
        </div>
        <button
          onClick={handleCreate}
          disabled={createTheme.isPending}
          className="inline-flex items-center gap-2 rounded-lg bg-indigo-600 px-4 py-2.5 text-sm font-medium text-white shadow-sm hover:bg-indigo-700 disabled:opacity-50 transition-colors"
        >
          {createTheme.isPending ? (
            <Loader2 className="h-4 w-4 animate-spin" />
          ) : (
            <Plus className="h-4 w-4" />
          )}
          Create Theme
        </button>
      </div>

      {createTheme.isError && (
        <div className="rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
          Failed to create theme: {createTheme.error.message}
        </div>
      )}

      {/* Theme grid */}
      {themes.length === 0 ? (
        <div className="flex flex-col items-center justify-center rounded-xl border-2 border-dashed border-gray-300 bg-white py-16 text-center">
          <Palette className="mb-3 h-10 w-10 text-gray-300" />
          <p className="text-sm font-medium text-gray-500">No themes yet</p>
          <p className="mt-1 text-xs text-gray-400">Create your first theme to get started.</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {themes.map((theme) => (
            <div
              key={theme.id}
              className="rounded-xl border border-gray-200 bg-white p-5 shadow-sm hover:shadow-md transition-shadow"
            >
              {/* Color swatches preview */}
              {theme.tokens?.colors && (
                <div className="mb-4 flex gap-1.5">
                  {(['primary', 'secondary', 'background', 'surface', 'text'] as const).map(
                    (key) => (
                      <div
                        key={key}
                        className="h-6 flex-1 rounded border border-gray-100"
                        style={{ backgroundColor: theme.tokens.colors[key] }}
                        title={key}
                      />
                    ),
                  )}
                </div>
              )}

              <div className="flex items-start justify-between gap-2">
                <div className="min-w-0">
                  <h3 className="truncate text-sm font-semibold text-gray-900">{theme.name}</h3>
                  <p className="text-xs text-gray-400">{formatDate(theme.created_at)}</p>
                </div>
                {theme.is_active && (
                  <span className="shrink-0 inline-flex items-center gap-1 rounded-full bg-green-100 px-2.5 py-0.5 text-xs font-medium text-green-700">
                    <CheckCircle2 className="h-3 w-3" />
                    Active
                  </span>
                )}
              </div>

              <div className="mt-4 flex items-center gap-1.5">
                <button
                  onClick={() => onEdit(theme)}
                  className="flex-1 rounded-lg border border-gray-200 px-3 py-1.5 text-xs font-medium text-gray-700 hover:bg-gray-50 transition-colors"
                >
                  <span className="flex items-center justify-center gap-1">
                    <Pencil className="h-3 w-3" />
                    Edit
                  </span>
                </button>
                <button
                  onClick={() => duplicateTheme.mutate({ id: theme.id, name: `${theme.name} (Copy)` })}
                  disabled={duplicateTheme.isPending}
                  className="flex-1 rounded-lg border border-gray-200 px-3 py-1.5 text-xs font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-50 transition-colors"
                  title="Duplicate"
                >
                  <span className="flex items-center justify-center gap-1">
                    <Copy className="h-3 w-3" />
                    Copy
                  </span>
                </button>
                {!theme.is_active && (
                  <button
                    onClick={() => activateTheme.mutate(theme.id)}
                    disabled={activateTheme.isPending}
                    className="flex-1 rounded-lg border border-green-200 bg-green-50 px-3 py-1.5 text-xs font-medium text-green-700 hover:bg-green-100 disabled:opacity-50 transition-colors"
                    title="Activate"
                  >
                    <span className="flex items-center justify-center gap-1">
                      <Globe className="h-3 w-3" />
                      Use
                    </span>
                  </button>
                )}
                {!theme.is_active && (
                  <button
                    onClick={() => handleDelete(theme)}
                    disabled={deleteTheme.isPending}
                    className="rounded-lg border border-red-100 bg-red-50 p-1.5 text-red-400 hover:bg-red-100 hover:text-red-600 disabled:opacity-50 transition-colors"
                    title="Delete"
                  >
                    <Trash2 className="h-3 w-3" />
                  </button>
                )}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}

// ---------------------------------------------------------------------------
// Main ThemeEditor page
// ---------------------------------------------------------------------------

export default function ThemeEditor() {
  const { data: themesData, isLoading } = useThemes()
  const [editingTheme, setEditingTheme] = useState<Theme | null>(null)

  const themes = themesData ?? []

  if (isLoading) {
    return (
      <div className="flex h-64 items-center justify-center">
        <Loader2 className="h-6 w-6 animate-spin text-gray-400" />
      </div>
    )
  }

  if (editingTheme) {
    return (
      <ThemeEditorForm
        theme={editingTheme}
        onBack={() => setEditingTheme(null)}
      />
    )
  }

  return <ThemeList themes={themes} onEdit={setEditingTheme} />
}
