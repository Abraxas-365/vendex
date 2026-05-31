import { useEffect, useRef } from 'react'
import { useParams, Link } from '@tanstack/react-router'
import { Loader2, AlertCircle, ArrowLeft } from 'lucide-react'
import { useStorePageBySlug, useStoreInfo } from '../../lib/store-hooks'

// ─── Component ───────────────────────────────────────────────────────────────

export default function DynamicPage() {
  const { slug } = useParams({ from: '/_store/pages/$slug' })
  const { data: page, isLoading, isError } = useStorePageBySlug(slug)
  const { data: storeInfo } = useStoreInfo()
  const accent = storeInfo?.accent_color ?? '#6366f1'
  const styleRef = useRef<HTMLStyleElement | null>(null)

  // Inject / remove the page's custom CSS as a <style> tag
  useEffect(() => {
    // Remove any previous page style
    if (styleRef.current) {
      styleRef.current.remove()
      styleRef.current = null
    }

    if (page?.css) {
      const el = document.createElement('style')
      el.setAttribute('data-dynamic-page', slug)
      el.textContent = page.css
      document.head.appendChild(el)
      styleRef.current = el
    }

    return () => {
      styleRef.current?.remove()
      styleRef.current = null
    }
  }, [page?.css, slug])

  // ── Loading ──────────────────────────────────────────────────────────────
  if (isLoading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="flex flex-col items-center gap-3 text-gray-400">
          <Loader2 size={32} className="animate-spin" />
          <p className="text-sm">Loading page…</p>
        </div>
      </div>
    )
  }

  // ── Error / 404 ──────────────────────────────────────────────────────────
  if (isError || !page) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center px-4">
        <div className="text-center max-w-sm">
          <div className="w-16 h-16 bg-red-50 rounded-full flex items-center justify-center mx-auto mb-4">
            <AlertCircle size={28} className="text-red-400" />
          </div>
          <h1 className="text-2xl font-bold text-gray-900 mb-2">Page not found</h1>
          <p className="text-gray-500 text-sm mb-6">
            The page{' '}
            <code className="bg-gray-100 px-1.5 py-0.5 rounded text-gray-700">/{slug}</code>{' '}
            doesn't exist or hasn't been published yet.
          </p>
          <Link
            to="/"
            className="inline-flex items-center gap-2 text-white font-medium px-6 py-3 rounded-xl transition-opacity hover:opacity-90"
            style={{ backgroundColor: accent }}
          >
            <ArrowLeft size={15} /> Go home
          </Link>
        </div>
      </div>
    )
  }

  // ── Unpublished guard ────────────────────────────────────────────────────
  if (page.status !== 'published') {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center px-4">
        <div className="text-center max-w-sm">
          <div className="w-16 h-16 bg-amber-50 rounded-full flex items-center justify-center mx-auto mb-4">
            <AlertCircle size={28} className="text-amber-400" />
          </div>
          <h1 className="text-2xl font-bold text-gray-900 mb-2">Not available</h1>
          <p className="text-gray-500 text-sm mb-6">
            This page is not published yet.
          </p>
          <Link
            to="/"
            className="inline-flex items-center gap-2 text-white font-medium px-6 py-3 rounded-xl transition-opacity hover:opacity-90"
            style={{ backgroundColor: accent }}
          >
            <ArrowLeft size={15} /> Go home
          </Link>
        </div>
      </div>
    )
  }

  // ── Render CMS page ──────────────────────────────────────────────────────
  return (
    <div className="min-h-screen bg-white">
      {/* Page header */}
      {page.title && (
        <header className="border-b border-gray-100 py-6 px-4 sm:px-6 lg:px-8">
          <div className="max-w-4xl mx-auto">
            <h1 className="text-2xl sm:text-3xl font-bold text-gray-900">{page.title}</h1>
            {page.meta?.description && (
              <p className="text-gray-500 mt-1 text-sm">{page.meta.description}</p>
            )}
          </div>
        </header>
      )}

      {/* CMS HTML content — page.html contains admin-authored markup */}
      {page.html && (
        <main className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          {/*
            dangerouslySetInnerHTML is intentional — this renders HTML
            authored by authenticated admins in the CMS editor.
            CSS is scoped via the <style> tag injected above.
          */}
          <div
            className="dynamic-page-content prose prose-gray max-w-none"
            // eslint-disable-next-line react/no-danger
            dangerouslySetInnerHTML={{ __html: page.html }}
          />
        </main>
      )}
    </div>
  )
}
