import { Link, useNavigate } from '@tanstack/react-router'
import { useCallback } from 'react'
import { ArrowRight, Sparkles, Truck, ShieldCheck, Package, Leaf, Gem, RefreshCw, Headphones } from 'lucide-react'
import { useStoreProducts, useStoreInfo } from '../../lib/store-hooks'
import ProductCard from '../../components/store/ProductCard'

// ─── Icon map for trust badges ──────────────────────────────────────────────

const ICON_MAP: Record<string, React.ComponentType<{ size?: number; className?: string; style?: React.CSSProperties }>> = {
  truck: Truck,
  shield: ShieldCheck,
  refresh: RefreshCw,
  headphones: Headphones,
  leaf: Leaf,
  gem: Gem,
  sparkles: Sparkles,
  package: Package,
}

// ─── Component ───────────────────────────────────────────────────────────────

export default function Home() {
  const navigate = useNavigate()
  const { data: productsPage } = useStoreProducts({ page: 1, page_size: 6 })
  const { data: info } = useStoreInfo()
  const products = productsPage?.items ?? []

  const handleProductClick = useCallback(
    (id: string) => void navigate({ to: '/products/$id', params: { id } }),
    [navigate],
  )

  const accent = info?.accent_color ?? '#6366f1'
  const isMinimal = info?.bg_style === 'minimal'
  const storeName = info?.store_name ?? 'Store'

  return (
    <div className="min-h-screen bg-gray-50">
      {/* ── Hero ─────────────────────────────────────────────────────────── */}
      <section
        className="relative overflow-hidden"
        style={{
          background: isMinimal
            ? '#fafaf9'
            : `linear-gradient(135deg, ${accent}08 0%, #ffffff 50%, ${accent}06 100%)`,
        }}
      >
        {!isMinimal && (
          <>
            <div
              className="absolute -top-40 -right-40 w-96 h-96 rounded-full opacity-30 blur-3xl pointer-events-none"
              style={{ background: accent }}
            />
            <div
              className="absolute -bottom-32 -left-32 w-80 h-80 rounded-full opacity-20 blur-3xl pointer-events-none"
              style={{ background: accent }}
            />
          </>
        )}

        <div className="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-24 md:py-36 text-center">
          {info?.announcement && (
            <span
              className="inline-flex items-center gap-2 text-xs font-semibold px-4 py-1.5 rounded-full mb-6"
              style={{
                backgroundColor: isMinimal ? '#f5f5f4' : `${accent}15`,
                color: isMinimal ? '#44403c' : accent,
              }}
            >
              <Sparkles size={13} />
              {info.announcement}
            </span>
          )}

          <h1
            className={`text-4xl sm:text-5xl md:text-6xl font-bold leading-tight tracking-tight mb-6 ${
              isMinimal ? 'text-stone-900' : 'text-gray-900'
            }`}
          >
            {info?.hero_title ? (
              info.hero_title
            ) : (
              <>
                Welcome to{' '}
                <span
                  className="text-transparent bg-clip-text"
                  style={{
                    backgroundImage: `linear-gradient(to right, ${accent}, ${accent}cc)`,
                  }}
                >
                  {storeName}
                </span>
              </>
            )}
          </h1>

          <p className={`max-w-xl mx-auto text-lg mb-10 ${isMinimal ? 'text-stone-500' : 'text-gray-500'}`}>
            {info?.hero_subtitle ?? 'Discover our curated collection.'}
          </p>

          <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
            <Link
              to="/products"
              className="inline-flex items-center gap-2 text-white font-semibold px-8 py-3.5 rounded-2xl transition-all shadow-lg hover:opacity-90"
              style={{ backgroundColor: accent }}
            >
              Shop Now
              <ArrowRight size={18} />
            </Link>
            <a
              href="#featured"
              className={`inline-flex items-center gap-2 bg-white font-medium px-8 py-3.5 rounded-2xl transition-colors ${
                isMinimal
                  ? 'border border-stone-200 text-stone-700 hover:border-stone-400'
                  : 'border border-gray-200 text-gray-700 hover:border-gray-400'
              }`}
            >
              See featured
            </a>
          </div>
        </div>
      </section>

      {/* ── Trust badges ─────────────────────────────────────────────────── */}
      {info?.trust_badges && info.trust_badges.length > 0 && (
        <section className={`border-y ${isMinimal ? 'bg-stone-50 border-stone-100' : 'bg-white border-gray-100'}`}>
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
            <div className="grid grid-cols-2 md:grid-cols-4 gap-6">
              {info.trust_badges.map((badge) => {
                const Icon = ICON_MAP[badge.icon] ?? Package
                return (
                  <div key={badge.title} className="flex items-start gap-3">
                    <div
                      className="w-9 h-9 rounded-xl flex items-center justify-center shrink-0"
                      style={{ backgroundColor: `${accent}12` }}
                    >
                      <Icon size={18} style={{ color: accent }} />
                    </div>
                    <div>
                      <p className={`text-sm font-semibold ${isMinimal ? 'text-stone-800' : 'text-gray-800'}`}>
                        {badge.title}
                      </p>
                      <p className={`text-xs ${isMinimal ? 'text-stone-500' : 'text-gray-500'}`}>{badge.desc}</p>
                    </div>
                  </div>
                )
              })}
            </div>
          </div>
        </section>
      )}

      {/* ── Featured products ─────────────────────────────────────────────── */}
      <section id="featured" className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        <div className="flex items-end justify-between mb-8">
          <div>
            <p
              className="text-xs font-semibold uppercase tracking-widest mb-1"
              style={{ color: accent }}
            >
              Hand-picked for you
            </p>
            <h2 className={`text-2xl sm:text-3xl font-bold ${isMinimal ? 'text-stone-900' : 'text-gray-900'}`}>
              Featured Products
            </h2>
          </div>
          <Link
            to="/products"
            className="hidden sm:inline-flex items-center gap-1.5 text-sm font-medium transition-colors hover:opacity-80"
            style={{ color: accent }}
          >
            View all <ArrowRight size={15} />
          </Link>
        </div>

        {products.length > 0 ? (
          <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-6 gap-4">
            {products.map((product) => (
              <ProductCard key={product.id} product={product} accent={accent} onNavigate={handleProductClick} />
            ))}
          </div>
        ) : (
          <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-6 gap-4">
            {Array.from({ length: 6 }).map((_, i) => (
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
        )}

        <div className="sm:hidden mt-6 text-center">
          <Link
            to="/products"
            className="inline-flex items-center gap-1.5 text-sm font-medium"
            style={{ color: accent }}
          >
            View all products <ArrowRight size={15} />
          </Link>
        </div>
      </section>

      {/* ── CTA banner ───────────────────────────────────────────────────── */}
      <section style={{ background: isMinimal ? '#1c1917' : accent }}>
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-14 text-center">
          <h2 className="text-2xl sm:text-3xl font-bold text-white mb-3">
            {isMinimal ? 'Discover your new favorites' : 'Ready to explore?'}
          </h2>
          <p className="text-white/70 mb-8">
            {isMinimal
              ? 'Curated with care. Free shipping on your first order.'
              : 'Join thousands of happy customers. Free shipping on your first order.'}
          </p>
          <Link
            to="/products"
            className="inline-flex items-center gap-2 bg-white font-semibold px-8 py-3.5 rounded-2xl hover:bg-gray-100 transition-colors shadow-lg"
            style={{ color: isMinimal ? '#1c1917' : accent }}
          >
            Start Shopping <ArrowRight size={18} />
          </Link>
        </div>
      </section>
    </div>
  )
}
