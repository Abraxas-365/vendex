import { ShoppingCart, Package } from 'lucide-react'
import type { Product } from '../../types/index'
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

// ─── Component ───────────────────────────────────────────────────────────────

interface ProductCardProps {
  product: Product
  onNavigate?: (id: string) => void
}

export default function ProductCard({ product, onNavigate }: ProductCardProps) {
  const { addItem } = useCart()

  const handleAddToCart = (e: React.MouseEvent) => {
    e.stopPropagation()
    addItem(product, 1)
  }

  const imageUrl =
    Array.isArray(product.images) && product.images.length > 0
      ? product.images[0]
      : null

  return (
    <div
      className="group bg-white rounded-2xl border border-gray-100 overflow-hidden shadow-sm hover:shadow-md transition-all duration-200 cursor-pointer flex flex-col"
      onClick={() => onNavigate?.(product.id)}
    >
      {/* Image */}
      <div className="relative aspect-square bg-gray-50 overflow-hidden">
        {imageUrl ? (
          <img
            src={imageUrl}
            alt={product.name}
            className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
          />
        ) : (
          <div className="w-full h-full flex items-center justify-center text-gray-300">
            <Package size={64} strokeWidth={1} />
          </div>
        )}

        {/* Out of stock badge */}
        {product.stock !== undefined && product.stock === 0 && (
          <div className="absolute top-3 left-3 bg-gray-800 text-white text-xs font-medium px-2 py-1 rounded-full">
            Out of stock
          </div>
        )}
      </div>

      {/* Info */}
      <div className="p-4 flex flex-col flex-1">
        <h3 className="font-medium text-gray-900 text-sm leading-snug line-clamp-2 mb-1 group-hover:text-indigo-600 transition-colors">
          {product.name}
        </h3>

        {product.sku && (
          <p className="text-xs text-gray-400 mb-2">SKU: {product.sku}</p>
        )}

        <div className="mt-auto flex items-center justify-between gap-2 pt-3">
          <span className="text-lg font-semibold text-gray-900">
            {formatPrice(product.price)}
          </span>

          <button
            onClick={handleAddToCart}
            disabled={product.stock === 0}
            className="flex items-center gap-1.5 bg-indigo-600 hover:bg-indigo-700 disabled:bg-gray-200 disabled:text-gray-400 text-white text-xs font-medium px-3 py-2 rounded-xl transition-colors"
          >
            <ShoppingCart size={14} />
            <span className="hidden sm:inline">Add</span>
          </button>
        </div>
      </div>
    </div>
  )
}
