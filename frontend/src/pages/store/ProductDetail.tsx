import { useState } from 'react'
import { Link, useParams } from '@tanstack/react-router'
import {
  ShoppingCart,
  ChevronLeft,
  Minus,
  Plus,
  Package,
  Tag,
  CheckCircle,
  XCircle,
  AlertCircle,
} from 'lucide-react'
import { useStoreProduct } from '../../lib/store-hooks'
import { useCart } from '../../lib/cart'

// ─── Helpers ─────────────────────────────────────────────────────────────────

function formatPrice(price: unknown): string {
  if (!price) return '$0.00'
  if (typeof price === 'object' && price !== null) {
    const p = price as { amount: number; currency?: string }
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: p.currency ?? 'USD',
    }).format(p.amount)
  }
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
  }).format(Number(price))
}

// ─── Stock indicator ─────────────────────────────────────────────────────────

function StockIndicator({ stock }: { stock?: number }) {
  if (stock === undefined) return null
  if (stock === 0) {
    return (
      <span className="inline-flex items-center gap-1.5 text-sm text-red-600 font-medium">
        <XCircle size={15} /> Out of stock
      </span>
    )
  }
  if (stock <= 5) {
    return (
      <span className="inline-flex items-center gap-1.5 text-sm text-amber-600 font-medium">
        <AlertCircle size={15} /> Only {stock} left
      </span>
    )
  }
  return (
    <span className="inline-flex items-center gap-1.5 text-sm text-green-600 font-medium">
      <CheckCircle size={15} /> In stock
    </span>
  )
}

// ─── Component ───────────────────────────────────────────────────────────────

export default function ProductDetail() {
  const { id } = useParams({ from: '/_store/products/$id' })
  const { data: product, isLoading, isError } = useStoreProduct(id)
  const { addItem } = useCart()

  const [selectedImage, setSelectedImage] = useState(0)
  const [quantity, setQuantity] = useState(1)
  const [addedToCart, setAddedToCart] = useState(false)

  const images: string[] =
    Array.isArray(product?.images) && product.images.length > 0
      ? product.images
      : []

  const handleAddToCart = () => {
    if (!product) return
    addItem(product, quantity)
    setAddedToCart(true)
    setTimeout(() => setAddedToCart(false), 2000)
  }

  // ── Loading state ────────────────────────────────────────────────────────
  if (isLoading) {
    return (
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
        <div className="grid md:grid-cols-2 gap-10 animate-pulse">
          <div className="aspect-square bg-gray-100 rounded-2xl" />
          <div className="space-y-4 py-4">
            <div className="h-6 bg-gray-100 rounded w-2/3" />
            <div className="h-10 bg-gray-100 rounded w-1/3" />
            <div className="h-4 bg-gray-100 rounded w-full" />
            <div className="h-4 bg-gray-100 rounded w-5/6" />
            <div className="h-4 bg-gray-100 rounded w-4/6" />
          </div>
        </div>
      </div>
    )
  }

  // ── Error / not found ────────────────────────────────────────────────────
  if (isError || !product) {
    return (
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-24 text-center">
        <div className="text-6xl mb-4">🌿</div>
        <h1 className="text-2xl font-bold text-gray-800 mb-2">Product not found</h1>
        <p className="text-gray-500 mb-6">
          The product you're looking for doesn't exist or has been removed.
        </p>
        <Link
          to="/products"
          className="inline-flex items-center gap-2 bg-indigo-600 hover:bg-indigo-700 text-white font-medium px-6 py-3 rounded-xl transition-colors"
        >
          <ChevronLeft size={16} /> Back to products
        </Link>
      </div>
    )
  }

  const tags: string[] = Array.isArray(product.tags) ? product.tags : []
  const outOfStock = product.stock === 0

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Breadcrumb */}
        <nav className="flex items-center gap-2 text-sm text-gray-400 mb-8">
          <Link to="/" className="hover:text-indigo-600 transition-colors">
            Home
          </Link>
          <span>/</span>
          <Link to="/products" className="hover:text-indigo-600 transition-colors">
            Products
          </Link>
          <span>/</span>
          <span className="text-gray-700 font-medium line-clamp-1">{product.name}</span>
        </nav>

        <div className="grid md:grid-cols-2 gap-10 lg:gap-16">
          {/* ── Image gallery ────────────────────────────────────────── */}
          <div>
            {/* Main image */}
            <div className="aspect-square bg-white rounded-2xl border border-gray-100 overflow-hidden mb-3 shadow-sm">
              {images.length > 0 ? (
                <img
                  src={images[selectedImage]}
                  alt={product.name}
                  className="w-full h-full object-cover"
                />
              ) : (
                <div className="w-full h-full flex items-center justify-center text-gray-200">
                  <Package size={96} strokeWidth={1} />
                </div>
              )}
            </div>

            {/* Thumbnails */}
            {images.length > 1 && (
              <div className="flex gap-2 overflow-x-auto pb-1">
                {images.map((src, i) => (
                  <button
                    key={i}
                    onClick={() => setSelectedImage(i)}
                    className={`shrink-0 w-16 h-16 rounded-xl overflow-hidden border-2 transition-colors ${
                      selectedImage === i
                        ? 'border-indigo-500'
                        : 'border-gray-100 hover:border-indigo-300'
                    }`}
                  >
                    <img
                      src={src}
                      alt={`${product.name} view ${i + 1}`}
                      className="w-full h-full object-cover"
                    />
                  </button>
                ))}
              </div>
            )}
          </div>

          {/* ── Product info ─────────────────────────────────────────── */}
          <div className="flex flex-col">
            {/* Category ID badge */}
            {product.category_id && (
              <span className="text-xs font-semibold text-indigo-600 uppercase tracking-widest mb-2">
                {product.category_id}
              </span>
            )}

            <h1 className="text-2xl sm:text-3xl font-bold text-gray-900 leading-snug mb-3">
              {product.name}
            </h1>

            {/* Price */}
            <p className="text-3xl font-bold text-gray-900 mb-4">
              {formatPrice(product.price)}
            </p>

            {/* Stock */}
            <div className="mb-4">
              <StockIndicator stock={product.stock} />
            </div>

            {/* Description */}
            {product.description && (
              <p className="text-gray-600 leading-relaxed mb-6">
                {product.description}
              </p>
            )}

            {/* SKU */}
            {product.sku && (
              <p className="text-sm text-gray-400 mb-5">
                SKU: <span className="font-medium text-gray-600">{product.sku}</span>
              </p>
            )}

            {/* Quantity selector */}
            <div className="flex items-center gap-4 mb-5">
              <span className="text-sm font-medium text-gray-700">Qty</span>
              <div className="flex items-center gap-2 bg-white border border-gray-200 rounded-xl overflow-hidden">
                <button
                  onClick={() => setQuantity((q) => Math.max(1, q - 1))}
                  className="px-3 py-2.5 hover:bg-gray-50 text-gray-600 transition-colors"
                  aria-label="Decrease quantity"
                >
                  <Minus size={15} />
                </button>
                <span className="w-10 text-center font-semibold text-gray-800 text-sm">
                  {quantity}
                </span>
                <button
                  onClick={() =>
                    setQuantity((q) =>
                      product.stock !== undefined ? Math.min(product.stock, q + 1) : q + 1
                    )
                  }
                  disabled={product.stock !== undefined && quantity >= product.stock}
                  className="px-3 py-2.5 hover:bg-gray-50 text-gray-600 disabled:opacity-40 transition-colors"
                  aria-label="Increase quantity"
                >
                  <Plus size={15} />
                </button>
              </div>
            </div>

            {/* Add to cart */}
            <button
              onClick={handleAddToCart}
              disabled={outOfStock}
              className={`flex items-center justify-center gap-2 w-full py-4 rounded-2xl font-semibold text-base transition-all duration-200 mb-4 shadow-sm ${
                addedToCart
                  ? 'bg-green-500 text-white shadow-green-200'
                  : outOfStock
                    ? 'bg-gray-100 text-gray-400 cursor-not-allowed'
                    : 'bg-indigo-600 hover:bg-indigo-700 text-white shadow-indigo-200'
              }`}
            >
              {addedToCart ? (
                <>
                  <CheckCircle size={18} /> Added to cart!
                </>
              ) : (
                <>
                  <ShoppingCart size={18} />
                  {outOfStock ? 'Out of Stock' : 'Add to Cart'}
                </>
              )}
            </button>

            <Link
              to="/cart"
              className="text-center text-sm text-gray-500 hover:text-indigo-600 transition-colors"
            >
              View cart →
            </Link>

            {/* Tags */}
            {tags.length > 0 && (
              <div className="mt-8 pt-6 border-t border-gray-100">
                <div className="flex items-center gap-2 flex-wrap">
                  <Tag size={14} className="text-gray-400" />
                  {tags.map((tag) => (
                    <span
                      key={tag}
                      className="bg-gray-100 text-gray-600 text-xs font-medium px-2.5 py-1 rounded-full"
                    >
                      {tag}
                    </span>
                  ))}
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
