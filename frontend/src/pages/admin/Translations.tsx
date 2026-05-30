import { useState } from 'react'
import { Languages, Save } from 'lucide-react'
import { useTranslationBundle, useEntityLocales, useSupportedLocales, useSetTranslations } from '../../lib/hooks'

// ---------------------------------------------------------------------------
// Translatable fields for each entity type
// ---------------------------------------------------------------------------

const TRANSLATABLE_FIELDS: Record<string, { key: string; label: string; multiline?: boolean }[]> = {
  product: [
    { key: 'name', label: 'Name' },
    { key: 'description', label: 'Description', multiline: true },
    { key: 'meta_title', label: 'Meta Title' },
    { key: 'meta_description', label: 'Meta Description', multiline: true },
  ],
  category: [
    { key: 'name', label: 'Name' },
    { key: 'description', label: 'Description', multiline: true },
    { key: 'meta_title', label: 'Meta Title' },
    { key: 'meta_description', label: 'Meta Description', multiline: true },
  ],
  collection: [
    { key: 'name', label: 'Name' },
    { key: 'description', label: 'Description', multiline: true },
  ],
  page: [
    { key: 'title', label: 'Title' },
    { key: 'meta_title', label: 'Meta Title' },
    { key: 'meta_description', label: 'Meta Description', multiline: true },
  ],
}

const ENTITY_TYPES = ['product', 'category', 'collection', 'page']

// ---------------------------------------------------------------------------
// Translation Form
// ---------------------------------------------------------------------------

interface TranslationFormProps {
  entityType: string
  entityId: string
  locale: string
}

function TranslationForm({ entityType, entityId, locale }: TranslationFormProps) {
  const fields = TRANSLATABLE_FIELDS[entityType] ?? []
  const { data: bundle, isPending } = useTranslationBundle(entityType, entityId, locale)
  const setTranslations = useSetTranslations()

  const [formValues, setFormValues] = useState<Record<string, string>>({})
  const [initialized, setInitialized] = useState(false)

  // Initialize form values from bundle when loaded
  if (bundle && !initialized) {
    setFormValues(bundle.fields ?? {})
    setInitialized(true)
  }

  // Reset when locale/entity changes
  const formKey = `${entityType}:${entityId}:${locale}`
  const [lastKey, setLastKey] = useState(formKey)
  if (formKey !== lastKey) {
    setLastKey(formKey)
    setFormValues(bundle?.fields ?? {})
    setInitialized(Boolean(bundle))
  }

  function handleSave(e: React.FormEvent) {
    e.preventDefault()
    // Filter out empty values
    const nonEmpty = Object.fromEntries(
      Object.entries(formValues).filter(([, v]) => v.trim() !== '')
    )
    setTranslations.mutate({ entityType, entityId, locale, fields: nonEmpty })
  }

  if (isPending) return <p className="text-sm text-slate-500">Loading translations…</p>

  return (
    <form onSubmit={handleSave} className="space-y-4">
      {fields.map((field) => (
        <div key={field.key}>
          <label className="mb-1 block text-xs font-medium text-slate-600">{field.label}</label>
          {field.multiline ? (
            <textarea
              value={formValues[field.key] ?? ''}
              onChange={(e) => setFormValues((prev) => ({ ...prev, [field.key]: e.target.value }))}
              rows={3}
              className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm outline-none focus:border-indigo-400 focus:ring-1 focus:ring-indigo-300 resize-none"
            />
          ) : (
            <input
              type="text"
              value={formValues[field.key] ?? ''}
              onChange={(e) => setFormValues((prev) => ({ ...prev, [field.key]: e.target.value }))}
              className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm outline-none focus:border-indigo-400 focus:ring-1 focus:ring-indigo-300"
            />
          )}
        </div>
      ))}

      <div className="flex items-center gap-3 pt-2">
        <button
          type="submit"
          disabled={setTranslations.isPending}
          className="flex items-center gap-1.5 rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700 disabled:opacity-60"
        >
          <Save size={14} />
          {setTranslations.isPending ? 'Saving…' : 'Save Translations'}
        </button>
        {setTranslations.isSuccess && (
          <p className="text-xs text-green-600">Saved!</p>
        )}
        {setTranslations.error && (
          <p className="text-xs text-red-600">{setTranslations.error.message}</p>
        )}
      </div>
    </form>
  )
}

// ---------------------------------------------------------------------------
// Existing Locales Panel
// ---------------------------------------------------------------------------

interface ExistingLocalesProps {
  entityType: string
  entityId: string
  selectedLocale: string
  onSelectLocale: (locale: string) => void
}

function ExistingLocales({ entityType, entityId, selectedLocale, onSelectLocale }: ExistingLocalesProps) {
  const { data: locales, isPending } = useEntityLocales(entityType, entityId)

  if (isPending) return <p className="text-xs text-slate-400">Loading…</p>
  if (!locales || locales.length === 0) return <p className="text-xs text-slate-400">No translations yet.</p>

  return (
    <div className="flex flex-wrap gap-2">
      {locales.map((locale) => (
        <button
          key={locale}
          onClick={() => onSelectLocale(locale)}
          className={`rounded-full px-3 py-1 text-xs font-medium transition-colors ${
            locale === selectedLocale
              ? 'bg-indigo-600 text-white'
              : 'bg-slate-100 text-slate-600 hover:bg-slate-200'
          }`}
        >
          {locale}
        </button>
      ))}
    </div>
  )
}

// ---------------------------------------------------------------------------
// Main Page
// ---------------------------------------------------------------------------

export default function Translations() {
  const [entityType, setEntityType] = useState('product')
  const [entityId, setEntityId] = useState('')
  const [locale, setLocale] = useState('')

  const { data: supportedLocales } = useSupportedLocales()
  const localeOptions = supportedLocales ?? []

  const canShow = Boolean(entityType) && Boolean(entityId) && Boolean(locale)

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center gap-2">
        <Languages size={20} className="text-indigo-600" />
        <h1 className="text-lg font-semibold text-slate-800">Translations</h1>
      </div>

      {/* Selectors */}
      <div className="rounded-xl border border-slate-200 bg-white p-5">
        <h2 className="mb-4 text-sm font-semibold text-slate-700">Select Entity</h2>
        <div className="flex flex-wrap gap-4 items-end">
          <div>
            <label className="mb-1 block text-xs font-medium text-slate-600">Entity Type</label>
            <select
              value={entityType}
              onChange={(e) => setEntityType(e.target.value)}
              className="rounded-lg border border-slate-200 px-3 py-2 text-sm outline-none focus:border-indigo-400"
            >
              {ENTITY_TYPES.map((t) => (
                <option key={t} value={t}>{t}</option>
              ))}
            </select>
          </div>
          <div>
            <label className="mb-1 block text-xs font-medium text-slate-600">Entity ID</label>
            <input
              type="text"
              value={entityId}
              onChange={(e) => setEntityId(e.target.value)}
              placeholder="Enter entity ID"
              className="w-64 rounded-lg border border-slate-200 px-3 py-2 text-sm outline-none focus:border-indigo-400"
            />
          </div>
          <div>
            <label className="mb-1 block text-xs font-medium text-slate-600">Locale</label>
            {localeOptions.length > 0 ? (
              <select
                value={locale}
                onChange={(e) => setLocale(e.target.value)}
                className="rounded-lg border border-slate-200 px-3 py-2 text-sm outline-none focus:border-indigo-400"
              >
                <option value="">Select locale…</option>
                {localeOptions.map((l) => (
                  <option key={l} value={l}>{l}</option>
                ))}
              </select>
            ) : (
              <input
                type="text"
                value={locale}
                onChange={(e) => setLocale(e.target.value)}
                placeholder="en, es, fr, de…"
                className="w-32 rounded-lg border border-slate-200 px-3 py-2 text-sm outline-none focus:border-indigo-400"
              />
            )}
          </div>
        </div>
      </div>

      {/* Existing Locales */}
      {entityType && entityId && (
        <div className="rounded-xl border border-slate-200 bg-white p-5">
          <h2 className="mb-3 text-sm font-semibold text-slate-700">Existing Locales for This Entity</h2>
          <ExistingLocales
            entityType={entityType}
            entityId={entityId}
            selectedLocale={locale}
            onSelectLocale={setLocale}
          />
        </div>
      )}

      {/* Translation Form */}
      {canShow && (
        <div className="rounded-xl border border-slate-200 bg-white p-5">
          <div className="mb-4 flex items-center justify-between">
            <h2 className="text-sm font-semibold text-slate-700">
              Translations — <span className="text-indigo-600">{locale}</span>
            </h2>
            <span className="rounded-full bg-slate-100 px-3 py-0.5 text-xs font-medium text-slate-500">
              {entityType} · {entityId}
            </span>
          </div>
          <TranslationForm
            key={`${entityType}:${entityId}:${locale}`}
            entityType={entityType}
            entityId={entityId}
            locale={locale}
          />
        </div>
      )}

      {!canShow && entityType && entityId && !locale && (
        <div className="rounded-xl border border-dashed border-slate-200 bg-slate-50 p-8 text-center">
          <Languages size={32} className="mx-auto mb-2 text-slate-300" />
          <p className="text-sm text-slate-400">Select a locale above to edit translations.</p>
        </div>
      )}

      {!entityId && (
        <div className="rounded-xl border border-dashed border-slate-200 bg-slate-50 p-8 text-center">
          <Languages size={32} className="mx-auto mb-2 text-slate-300" />
          <p className="text-sm text-slate-400">Enter an entity ID to manage its translations.</p>
        </div>
      )}
    </div>
  )
}
