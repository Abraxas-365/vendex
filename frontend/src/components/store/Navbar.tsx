import { useState } from 'react'
import { Link } from '@tanstack/react-router'
import { ShoppingCart, Search, Leaf, Menu, X } from 'lucide-react'
import { useCart } from '../../lib/cart'

export default function Navbar() {
  const { itemCount } = useCart()
  const [mobileOpen, setMobileOpen] = useState(false)
  const [searchQuery, setSearchQuery] = useState('')

  return (
    <nav className="sticky top-0 z-50 bg-white/95 backdrop-blur border-b border-gray-100 shadow-sm">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          {/* Logo */}
          <Link
            to="/"
            className="flex items-center gap-2 font-bold text-xl text-indigo-600 hover:text-indigo-700 transition-colors"
          >
            <Leaf size={22} />
            <span>Hada Store</span>
          </Link>

          {/* Desktop nav links */}
          <div className="hidden md:flex items-center gap-8">
            <Link
              to="/"
              className="text-sm font-medium text-gray-600 hover:text-indigo-600 transition-colors"
              activeProps={{ className: 'text-indigo-600' }}
            >
              Home
            </Link>
            <Link
              to="/products"
              className="text-sm font-medium text-gray-600 hover:text-indigo-600 transition-colors"
              activeProps={{ className: 'text-indigo-600' }}
            >
              Products
            </Link>
          </div>

          {/* Search + Cart */}
          <div className="flex items-center gap-3">
            {/* Search */}
            <div className="hidden sm:flex items-center gap-2 bg-gray-50 border border-gray-200 rounded-xl px-3 py-1.5 focus-within:border-indigo-400 focus-within:ring-2 focus-within:ring-indigo-100 transition-all">
              <Search size={15} className="text-gray-400 shrink-0" />
              <input
                type="search"
                placeholder="Search products…"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="bg-transparent text-sm text-gray-700 placeholder-gray-400 outline-none w-40 lg:w-56"
              />
            </div>

            {/* Cart */}
            <Link
              to="/cart"
              className="relative flex items-center justify-center w-10 h-10 rounded-xl hover:bg-indigo-50 text-gray-600 hover:text-indigo-600 transition-colors"
              aria-label={`Cart (${itemCount} items)`}
            >
              <ShoppingCart size={20} />
              {itemCount > 0 && (
                <span className="absolute -top-1 -right-1 bg-indigo-600 text-white text-[10px] font-bold w-5 h-5 rounded-full flex items-center justify-center">
                  {itemCount > 99 ? '99+' : itemCount}
                </span>
              )}
            </Link>

            {/* Mobile menu toggle */}
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

      {/* Mobile menu */}
      {mobileOpen && (
        <div className="md:hidden border-t border-gray-100 bg-white px-4 pb-4">
          {/* Mobile search */}
          <div className="flex items-center gap-2 bg-gray-50 border border-gray-200 rounded-xl px-3 py-2 mt-3 mb-3">
            <Search size={15} className="text-gray-400 shrink-0" />
            <input
              type="search"
              placeholder="Search products…"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="bg-transparent text-sm text-gray-700 placeholder-gray-400 outline-none flex-1"
            />
          </div>
          <nav className="flex flex-col gap-1">
            <Link
              to="/"
              onClick={() => setMobileOpen(false)}
              className="px-3 py-2.5 rounded-xl text-sm font-medium text-gray-700 hover:bg-indigo-50 hover:text-indigo-600 transition-colors"
            >
              Home
            </Link>
            <Link
              to="/products"
              onClick={() => setMobileOpen(false)}
              className="px-3 py-2.5 rounded-xl text-sm font-medium text-gray-700 hover:bg-indigo-50 hover:text-indigo-600 transition-colors"
            >
              Products
            </Link>
            <Link
              to="/cart"
              onClick={() => setMobileOpen(false)}
              className="px-3 py-2.5 rounded-xl text-sm font-medium text-gray-700 hover:bg-indigo-50 hover:text-indigo-600 transition-colors"
            >
              Cart {itemCount > 0 && `(${itemCount})`}
            </Link>
          </nav>
        </div>
      )}
    </nav>
  )
}
