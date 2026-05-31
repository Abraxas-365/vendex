import { Link } from '@tanstack/react-router'
import { ArrowRight, Sparkles, Package, Truck, ShieldCheck } from 'lucide-react'
import { useStoreProducts } from '../../lib/store-hooks'
import ProductCard from '../../components/store/ProductCard'

// ─── Static data ─────────────────────────────────────────────────────────────

const CATEGORIES = [
  { name: 'Skincare', emoji: '✨', description: 'Glow naturally' },
  { name: 'Haircare', emoji: '💆', description: 'Nourish every strand' },
  { name: 'Body', emoji: '🌿', description: 'From head to toe' },
  { name: 'Wellness', emoji: '🌸', description: 'Mind & body balance' },
]

const PERKS = [
  { icon: Truck, title: 'Free Shipping', body: 'On orders over $50' },
  { icon: ShieldCheck, title: 'Natural Ingredients', body: '100% plant-based formulas' },
  { icon: Sparkles, title: 'Cruelty Free', body: 'Certified & ethical' },
  { icon: Package, title: 'Easy Returns', body: '30-day hassle-free returns' },
]

// ─── Component ───────────────────────────────────────────────────────────────

export default function Home() {
  const { data: productsPage } = useStoreProducts({ page: 1, page_size: 6 })
  const products = productsPage?.items ?? []

  return (
    <div className="min-h-screen bg-gray-50">
      {/* ── Hero ─────────────────────────────────────────────────────────── */}
      <section className="relative overflow-hidden bg-gradient-to-br from-indigo-50 via-white to-purple-50">
        {/* Decorative blobs */}
        <div className="absolute -top-40 -right-40 w-96 h-96 rounded-full bg-indigo-100 opacity-50 blur-3xl pointer-events-none" />
        <div className="absolute -bottom-32 -left-32 w-80 h-80 rounded-full bg-purple-100 opacity-40 blur-3xl pointer-events-none" />

        <div className="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-24 md:py-36 text-center">
          <span className="inline-flex items-center gap-2 bg-indigo-100 text-indigo-700 text-xs font-semibold px-4 py-1.5 rounded-full mb-6">
            <Sparkles size={13} />
            New arrivals this season
          </span>

          <h1 className="text-4xl sm:text-5xl md:text-6xl font-bold text-gray-900 leading-tight tracking-tight mb-6">
            Welcome to{' '}
            <span className="text-transparent bg-clip-text bg-gradient-to-r from-indigo-600 to-purple-600">
              Vendex Store
            </span>
          </h1>

          <p className="max-w-xl mx-auto text-lg text-gray-500 mb-10">
            Discover our curated collection of natural beauty and wellness
            products — crafted with care for you and the planet.
          </p>

          <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
            <Link
              to="/products"
              className="inline-flex items-center gap-2 bg-indigo-600 hover:bg-indigo-700 text-white font-semibold px-8 py-3.5 rounded-2xl transition-colors shadow-lg shadow-indigo-200"
            >
              Shop Now
              <ArrowRight size={18} />
            </Link>
            <a
              href="#featured"
              className="inline-flex items-center gap-2 bg-white border border-gray-200 hover:border-indigo-300 text-gray-700 hover:text-indigo-600 font-medium px-8 py-3.5 rounded-2xl transition-colors"
            >
              See featured
            </a>
          </div>
        </div>
      </section>

      {/* ── Perks bar ────────────────────────────────────────────────────── */}
      <section className="bg-white border-y border-gray-100">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
          <div className="grid grid-cols-2 md:grid-cols-4 gap-6">
            {PERKS.map(({ icon: Icon, title, body }) => (
              <div key={title} className="flex items-start gap-3">
                <div className="w-9 h-9 rounded-xl bg-indigo-50 flex items-center justify-center shrink-0">
                  <Icon size={18} className="text-indigo-600" />
                </div>
                <div>
                  <p className="text-sm font-semibold text-gray-800">{title}</p>
                  <p className="text-xs text-gray-500">{body}</p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ── Featured products ─────────────────────────────────────────────── */}
      <section id="featured" className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        <div className="flex items-end justify-between mb-8">
          <div>
            <p className="text-xs font-semibold text-indigo-600 uppercase tracking-widest mb-1">
              Hand-picked for you
            </p>
            <h2 className="text-2xl sm:text-3xl font-bold text-gray-900">
              Featured Products
            </h2>
          </div>
          <Link
            to="/products"
            className="hidden sm:inline-flex items-center gap-1.5 text-sm font-medium text-indigo-600 hover:text-indigo-700 transition-colors"
          >
            View all <ArrowRight size={15} />
          </Link>
        </div>

        {products.length > 0 ? (
          <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-6 gap-4">
            {products.map((product) => (
              <ProductCard key={product.id} product={product} />
            ))}
          </div>
        ) : (
          /* Empty / loading state */
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

        {/* Mobile "view all" link */}
        <div className="sm:hidden mt-6 text-center">
          <Link
            to="/products"
            className="inline-flex items-center gap-1.5 text-sm font-medium text-indigo-600"
          >
            View all products <ArrowRight size={15} />
          </Link>
        </div>
      </section>

      {/* ── Categories ───────────────────────────────────────────────────── */}
      <section className="bg-white border-t border-gray-100">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
          <div className="text-center mb-10">
            <p className="text-xs font-semibold text-indigo-600 uppercase tracking-widest mb-1">
              Browse by
            </p>
            <h2 className="text-2xl sm:text-3xl font-bold text-gray-900">
              Shop Categories
            </h2>
          </div>

          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            {CATEGORIES.map((cat) => (
              <Link
                key={cat.name}
                to="/products"
                className="group relative bg-gradient-to-br from-indigo-50 to-purple-50 hover:from-indigo-100 hover:to-purple-100 border border-indigo-100 rounded-2xl p-6 text-center transition-all duration-200 hover:shadow-md hover:-translate-y-0.5"
              >
                <div className="text-4xl mb-3">{cat.emoji}</div>
                <h3 className="font-semibold text-gray-800 group-hover:text-indigo-700 transition-colors">
                  {cat.name}
                </h3>
                <p className="text-xs text-gray-500 mt-1">{cat.description}</p>
              </Link>
            ))}
          </div>
        </div>
      </section>

      {/* ── CTA banner ───────────────────────────────────────────────────── */}
      <section className="bg-gradient-to-r from-indigo-600 to-purple-600">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-14 text-center">
          <h2 className="text-2xl sm:text-3xl font-bold text-white mb-3">
            Ready to start your wellness journey?
          </h2>
          <p className="text-indigo-200 mb-8">
            Join thousands of happy customers. Free shipping on your first order.
          </p>
          <Link
            to="/products"
            className="inline-flex items-center gap-2 bg-white text-indigo-700 font-semibold px-8 py-3.5 rounded-2xl hover:bg-indigo-50 transition-colors shadow-lg"
          >
            Start Shopping <ArrowRight size={18} />
          </Link>
        </div>
      </section>
    </div>
  )
}
