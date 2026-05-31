import { Link } from '@tanstack/react-router'
import { useStoreInfo } from '../../lib/store-hooks'

// ─── Component ───────────────────────────────────────────────────────────────

export default function Footer() {
  const { data: storeInfo } = useStoreInfo()
  const storeName = storeInfo?.store_name ?? 'Our Store'
  const accent = storeInfo?.accent_color ?? '#6366f1'
  const year = new Date().getFullYear()

  return (
    <footer className="bg-white border-t border-gray-100 mt-auto">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        {/* Grid columns */}
        <div className="grid grid-cols-2 sm:grid-cols-4 gap-8 mb-10">
          {/* Brand */}
          <div className="col-span-2 sm:col-span-1">
            <div className="flex items-center gap-2 mb-3">
              <div
                className="w-7 h-7 rounded-lg flex items-center justify-center text-white text-sm font-bold"
                style={{ backgroundColor: accent }}
              >
                {storeName.charAt(0).toUpperCase()}
              </div>
              <span className="font-semibold text-gray-900">{storeName}</span>
            </div>
            <p className="text-sm text-gray-500 leading-relaxed">
              {storeInfo?.tagline ?? 'Quality products delivered to your door.'}
            </p>
          </div>

          {/* Shop */}
          <div>
            <h4 className="text-xs font-semibold uppercase tracking-widest mb-4" style={{ color: accent }}>
              Shop
            </h4>
            <ul className="space-y-2.5">
              <li>
                <Link
                  to="/"
                  className="text-sm text-gray-500 hover:text-gray-900 transition-colors"
                >
                  Home
                </Link>
              </li>
              <li>
                <Link
                  to="/products"
                  className="text-sm text-gray-500 hover:text-gray-900 transition-colors"
                >
                  All Products
                </Link>
              </li>
            </ul>
          </div>

          {/* Info */}
          <div>
            <h4 className="text-xs font-semibold uppercase tracking-widest mb-4" style={{ color: accent }}>
              Info
            </h4>
            <ul className="space-y-2.5">
              <li>
                <Link
                  to="/pages/$slug"
                  params={{ slug: 'about' }}
                  className="text-sm text-gray-500 hover:text-gray-900 transition-colors"
                >
                  About
                </Link>
              </li>
              <li>
                <Link
                  to="/pages/$slug"
                  params={{ slug: 'contact' }}
                  className="text-sm text-gray-500 hover:text-gray-900 transition-colors"
                >
                  Contact
                </Link>
              </li>
              <li>
                <Link
                  to="/pages/$slug"
                  params={{ slug: 'faq' }}
                  className="text-sm text-gray-500 hover:text-gray-900 transition-colors"
                >
                  FAQ
                </Link>
              </li>
            </ul>
          </div>

          {/* Policies */}
          <div>
            <h4 className="text-xs font-semibold uppercase tracking-widest mb-4" style={{ color: accent }}>
              Policies
            </h4>
            <ul className="space-y-2.5">
              <li>
                <Link
                  to="/pages/$slug"
                  params={{ slug: 'privacy' }}
                  className="text-sm text-gray-500 hover:text-gray-900 transition-colors"
                >
                  Privacy Policy
                </Link>
              </li>
              <li>
                <Link
                  to="/pages/$slug"
                  params={{ slug: 'shipping-returns' }}
                  className="text-sm text-gray-500 hover:text-gray-900 transition-colors"
                >
                  Shipping & Returns
                </Link>
              </li>
            </ul>
          </div>
        </div>

        {/* Divider + copyright */}
        <div className="border-t border-gray-100 pt-6 flex flex-col sm:flex-row items-center justify-between gap-3">
          <p className="text-sm text-gray-400">
            © {year} {storeName}. All rights reserved.
          </p>
          {storeInfo?.store_email && (
            <a
              href={`mailto:${storeInfo.store_email}`}
              className="text-sm text-gray-400 hover:text-gray-600 transition-colors"
            >
              {storeInfo.store_email}
            </a>
          )}
        </div>
      </div>
    </footer>
  )
}
