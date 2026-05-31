import { useState } from 'react'
import { Link } from '@tanstack/react-router'
import { ShoppingCart, Search, Store, Menu, X } from 'lucide-react'
import { useCart } from '../../lib/cart'
import { useStoreInfo } from '../../lib/store-hooks'

export default function Navbar() {
  const { itemCount } = useCart()
  const { data: info } = useStoreInfo()
  const [mobileOpen, setMobileOpen] = useState(false)
  const [searchQuery, setSearchQuery] = useState('')

  const accent = info?.accent_color ?? '#6366f1'
  const storeName = info?.store_name ?? 'Store'

  return (
    <nav className="sticky top-0 z-50 bg-white/95 backdrop-blur border-b border-gray-100 shadow-sm">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          {/* Logo */}
          <Link
            to="/"
            className="flex items-center gap-2 font-bold text-xl transition-colors hover:opacity-80"
            style={{ color: accent }}
          >
            <Store size={22} />
            <span>{storeName}</span>
          </Link>

          {/* Desktop nav links */}
          <div className="hidden md:flex items-center gap-8">
            <Link
              to="/"
              className="text-sm font-medium text-gray-600 hover:opacity-80 transition-colors"
            >
              Home
            </Link>
            <Link
              to="/products"
              className="text-sm font-medium text-gray-600 hover:opacity-80 transition-colors"
            >
              Products
            </Link>
          </div>

          {/* Search + Cart */}
          <div className="flex items-center gap-3">
            <div className="hidden sm:flex items-center gap-2 bg-gray-50 border border-gray-200 rounded-xl px-3 py-1.5 transition-all"
              style={{ ['--tw-ring-color' as string]: `${accent}30` }}
            >
              <Search size={15} className="text-gray-400 shrink-0" />
              <input
                type="search"
                placeholder="Search products..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="bg-transparent text-sm text-gray-700 placeholder-gray-400 outline-none w-40 lg:w-56"
              />
            </div>

            <Link
              to="/cart"
              className="relative flex items-center justify-center w-10 h-10 rounded-xl text-gray-600 hover:opacity-80 transition-colors"
              aria-label={`Cart (${itemCount} items)`}
            >
              <ShoppingCart size={20} />
              {itemCount > 0 && (
                <span
                  className="absolute -top-1 -right-1 text-white text-[10px] font-bold w-5 h-5 rounded-full flex items-center justify-center"
                  style={{ backgroundColor: accent }}
                >
                  {itemCount > 99 ? '99+' : itemCount}
                </span>
              )}
            </Link>

            <button
              className="md:hidden flex items-center justify-center w-10 h-10 rounded-xl hover:bg-gray-100 text-gray-600 transition-colors"
              onClick={() => setMobileOpen((v) => !v)}
              aria-label="Toggle menu"
            >
              {mobileOpen ? <X size={20} /> : <Menu size={20} />}
            </button>
          </div>
        </div>
      </div>

      {mobileOpen && (
        <div className="md:hidden border-t border-gray-100 bg-white px-4 pb-4">
          <div className="flex items-center gap-2 bg-gray-50 border border-gray-200 rounded-xl px-3 py-2 mt-3 mb-3">
            <Search size={15} className="text-gray-400 shrink-0" />
            <input
              type="search"
              placeholder="Search products..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="bg-transparent text-sm text-gray-700 placeholder-gray-400 outline-none flex-1"
            />
          </div>
          <nav className="flex flex-col gap-1">
            <Link to="/" onClick={() => setMobileOpen(false)}
              className="px-3 py-2.5 rounded-xl text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors">
              Home
            </Link>
            <Link to="/products" onClick={() => setMobileOpen(false)}
              className="px-3 py-2.5 rounded-xl text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors">
              Products
            </Link>
            <Link to="/cart" onClick={() => setMobileOpen(false)}
              className="px-3 py-2.5 rounded-xl text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors">
              Cart {itemCount > 0 && `(${itemCount})`}
            </Link>
          </nav>
        </div>
      )}
    </nav>
  )
}
