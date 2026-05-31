import { useState } from 'react'
import { SlidersHorizontal, ChevronLeft, ChevronRight, X } from 'lucide-react'
import { useStoreProducts } from '../../lib/store-hooks'
import ProductCard from '../../components/store/ProductCard'

// ─── Types ───────────────────────────────────────────────────────────────────

type SortOption = 'newest' | 'price_asc' | 'price_desc' | 'name_asc'

interface Filters {
  category: string
  minPrice: string
  maxPrice: string
  sort: SortOption
}

const CATEGORIES = ['All', 'Skincare', 'Haircare', 'Body', 'Wellness']

const SORT_OPTIONS: { value: SortOption; label: string }[] = [
  { value: 'newest', label: 'Newest' },
  { value: 'price_asc', label: 'Price: Low → High' },
  { value: 'price_desc', label: 'Price: High → Low' },
  { value: 'name_asc', label: 'Name A–Z' },
]

const PAGE_SIZE = 12

// ─── Component ───────────────────────────────────────────────────────────────

export default function ProductList() {
  const [page, setPage] = useState(1)
  const [filters, setFilters] = useState<Filters>({
    category: '',
    minPrice: '',
    maxPrice: '',
    sort: 'newest',
  })
  const [sidebarOpen, setSidebarOpen] = useState(false)

  const { data, isLoading, isError } = useStoreProducts({
    page,
    page_size: PAGE_SIZE,
  })

  const products = data?.items ?? []
  const totalPages = data?.total_pages ?? 1

  const updateFilter = <K extends keyof Filters>(key: K, value: Filters[K]) => {
    setFilters((prev) => ({ ...prev, [key]: value }))
    setPage(1)
  }

  const clearFilters = () => {
    setFilters({ category: '', minPrice: '', maxPrice: '', sort: 'newest' })
    setPage(1)
  }

  const hasActiveFilters =
    filters.category !== '' ||
    filters.minPrice !== '' ||
    filters.maxPrice !== ''

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* ── Page header ──────────────────────────────────────────────── */}
        <div className="mb-6">
          <h1 className="text-2xl sm:text-3xl font-bold text-gray-900">Products</h1>
          {data && (
            <p className="text-sm text-gray-500 mt-1">
              {data.total} product{data.total !== 1 ? 's' : ''} found
            </p>
          )}
        </div>

        {/* ── Top bar (mobile filter toggle + sort) ────────────────────── */}
        <div className="flex items-center justify-between gap-3 mb-6 lg:hidden">
          <button
            onClick={() => setSidebarOpen(true)}
            className="flex items-center gap-2 bg-white border border-gray-200 text-sm font-medium text-gray-700 px-4 py-2.5 rounded-xl hover:border-indigo-300 transition-colors"
          >
            <SlidersHorizontal size={15} />
            Filters
            {hasActiveFilters && (
              <span className="bg-indigo-600 text-white text-xs w-4 h-4 rounded-full flex items-center justify-center">
                !
              </span>
            )}
          </button>

          <select
            value={filters.sort}
            onChange={(e) => updateFilter('sort', e.target.value as SortOption)}
            className="bg-white border border-gray-200 text-sm text-gray-700 px-3 py-2.5 rounded-xl outline-none focus:border-indigo-400 cursor-pointer"
          >
            {SORT_OPTIONS.map((opt) => (
              <option key={opt.value} value={opt.value}>
                {opt.label}
              </option>
            ))}
          </select>
        </div>

        <div className="flex gap-8">
          {/* ── Sidebar (desktop) / Drawer (mobile) ───────────────────── */}

          {/* Mobile overlay */}
          {sidebarOpen && (
            <div
              className="fixed inset-0 bg-black/40 z-40 lg:hidden"
              onClick={() => setSidebarOpen(false)}
            />
          )}

          <aside
            className={`
              fixed inset-y-0 left-0 z-50 w-72 bg-white shadow-2xl p-6 overflow-y-auto transition-transform duration-300 lg:static lg:block lg:w-56 lg:shadow-none lg:p-0 lg:translate-x-0 lg:bg-transparent lg:z-auto lg:shrink-0
              ${sidebarOpen ? 'translate-x-0' : '-translate-x-full'}
            `}
          >
            <div className="flex items-center justify-between mb-6 lg:hidden">
              <h2 className="font-semibold text-gray-900">Filters</h2>
              <button
                onClick={() => setSidebarOpen(false)}
                className="p-1.5 hover:bg-gray-100 rounded-lg text-gray-500"
              >
                <X size={18} />
              </button>
            </div>

            <div className="space-y-6">
              {/* Category */}
              <div>
                <h3 className="text-xs font-semibold text-gray-500 uppercase tracking-widest mb-3">
                  Category
                </h3>
                <ul className="space-y-1">
                  {CATEGORIES.map((cat) => {
                    const value = cat === 'All' ? '' : cat.toLowerCase()
                    const active = filters.category === value
                    return (
                      <li key={cat}>
                        <button
                          onClick={() => {
                            updateFilter('category', value)
                            setSidebarOpen(false)
                          }}
                          className={`w-full text-left px-3 py-2 rounded-xl text-sm transition-colors ${
                            active
                              ? 'bg-indigo-50 text-indigo-700 font-medium'
                              : 'text-gray-600 hover:bg-gray-100'
                          }`}
                        >
                          {cat}
                        </button>
                      </li>
                    )
                  })}
                </ul>
              </div>

              {/* Price range */}
              <div>
                <h3 className="text-xs font-semibold text-gray-500 uppercase tracking-widest mb-3">
                  Price Range
                </h3>
                <div className="flex items-center gap-2">
                  <input
                    type="number"
                    placeholder="Min"
                    value={filters.minPrice}
                    onChange={(e) => updateFilter('minPrice', e.target.value)}
                    className="w-full border border-gray-200 rounded-xl px-3 py-2 text-sm outline-none focus:border-indigo-400"
                  />
                  <span className="text-gray-400 text-sm shrink-0">–</span>
                  <input
                    type="number"
                    placeholder="Max"
                    value={filters.maxPrice}
                    onChange={(e) => updateFilter('maxPrice', e.target.value)}
                    className="w-full border border-gray-200 rounded-xl px-3 py-2 text-sm outline-none focus:border-indigo-400"
                  />
                </div>
              </div>

              {/* Clear filters */}
              {hasActiveFilters && (
                <button
                  onClick={clearFilters}
                  className="w-full text-sm text-red-500 hover:text-red-600 font-medium py-2 border border-red-100 hover:border-red-200 rounded-xl transition-colors"
                >
                  Clear filters
                </button>
              )}
            </div>
          </aside>

          {/* ── Main content ──────────────────────────────────────────── */}
          <div className="flex-1 min-w-0">
            {/* Desktop sort + results count */}
            <div className="hidden lg:flex items-center justify-between mb-5">
              {data && (
                <p className="text-sm text-gray-500">
                  {data.total} result{data.total !== 1 ? 's' : ''}
                </p>
              )}
              <select
                value={filters.sort}
                onChange={(e) => updateFilter('sort', e.target.value as SortOption)}
                className="bg-white border border-gray-200 text-sm text-gray-700 px-3 py-2 rounded-xl outline-none focus:border-indigo-400 cursor-pointer"
              >
                {SORT_OPTIONS.map((opt) => (
                  <option key={opt.value} value={opt.value}>
                    {opt.label}
                  </option>
                ))}
              </select>
            </div>

            {/* Product grid */}
            {isError ? (
              <div className="text-center py-20">
                <p className="text-gray-400 text-lg mb-2">Failed to load products</p>
                <p className="text-gray-300 text-sm">Please try again later.</p>
              </div>
            ) : isLoading ? (
              <div className="grid grid-cols-2 md:grid-cols-3 xl:grid-cols-4 gap-4">
                {Array.from({ length: PAGE_SIZE }).map((_, i) => (
                  <div
                    key={i}
                    className="bg-white rounded-2xl border border-gray-100 overflow-hidden animate-pulse"
                  >
                    <div className="aspect-square bg-gray-100" />
                    <div className="p-4 space-y-2">
                      <div className="h-3 bg-gray-100 rounded w-3/4" />
                      <div className="h-3 bg-gray-100 rounded w-1/2" />
                      <div className="h-8 bg-gray-100 rounded w-full mt-3" />
                    </div>
                  </div>
                ))}
              </div>
            ) : products.length === 0 ? (
              <div className="text-center py-24">
                <div className="text-6xl mb-4">🌿</div>
                <p className="text-gray-500 text-lg font-medium mb-1">No products found</p>
                <p className="text-gray-400 text-sm">
                  Try adjusting your filters or{' '}
                  <button
                    onClick={clearFilters}
                    className="text-indigo-600 hover:underline"
                  >
                    clear all
                  </button>
                </p>
              </div>
            ) : (
              <div className="grid grid-cols-2 md:grid-cols-3 xl:grid-cols-4 gap-4">
                {products.map((product) => (
                  <ProductCard key={product.id} product={product} />
                ))}
              </div>
            )}

            {/* Pagination */}
            {totalPages > 1 && (
              <div className="flex items-center justify-center gap-2 mt-10">
                <button
                  onClick={() => setPage((p) => Math.max(1, p - 1))}
                  disabled={page === 1}
                  className="flex items-center gap-1 px-4 py-2 text-sm font-medium text-gray-600 bg-white border border-gray-200 rounded-xl disabled:opacity-40 hover:border-indigo-300 hover:text-indigo-600 transition-colors"
                >
                  <ChevronLeft size={15} /> Prev
                </button>

                <div className="flex items-center gap-1">
                  {Array.from({ length: totalPages }, (_, i) => i + 1)
                    .filter(
                      (p) =>
                        p === 1 ||
                        p === totalPages ||
                        Math.abs(p - page) <= 1
                    )
                    .reduce<(number | '…')[]>((acc, p, idx, arr) => {
                      if (idx > 0 && p - (arr[idx - 1] as number) > 1) {
                        acc.push('…')
                      }
                      acc.push(p)
                      return acc
                    }, [])
                    .map((p, i) =>
                      p === '…' ? (
                        <span key={`ellipsis-${i}`} className="px-1 text-gray-400 text-sm">
                          …
                        </span>
                      ) : (
                        <button
                          key={p}
                          onClick={() => setPage(p as number)}
                          className={`w-9 h-9 rounded-xl text-sm font-medium transition-colors ${
                            page === p
                              ? 'bg-indigo-600 text-white'
                              : 'bg-white border border-gray-200 text-gray-600 hover:border-indigo-300 hover:text-indigo-600'
                          }`}
                        >
                          {p}
                        </button>
                      )
                    )}
                </div>

                <button
                  onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
                  disabled={page === totalPages}
                  className="flex items-center gap-1 px-4 py-2 text-sm font-medium text-gray-600 bg-white border border-gray-200 rounded-xl disabled:opacity-40 hover:border-indigo-300 hover:text-indigo-600 transition-colors"
                >
                  Next <ChevronRight size={15} />
                </button>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
