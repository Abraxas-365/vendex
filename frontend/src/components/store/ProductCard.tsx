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
    }).format(p.amount / 100)
  }
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
  }).format(Number(price))
}

// ─── Component ───────────────────────────────────────────────────────────────

interface ProductCardProps {
  product: Product
  accent?: string
  onNavigate?: (id: string) => void
}

export default function ProductCard({ product, accent = '#6366f1', onNavigate }: ProductCardProps) {
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
          <div className="w-full h-full flex flex-col items-center justify-center bg-gradient-to-br from-gray-50 to-gray-100">
            <Package size={48} strokeWidth={1} className="text-gray-300 mb-2" />
            <span className="text-[10px] font-medium text-gray-300 uppercase tracking-wider">No image</span>
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
        <h3 className="font-medium text-gray-900 text-sm leading-snug line-clamp-2 mb-1 transition-colors group-hover:opacity-80">
          {product.name}
        </h3>

        <div className="mt-auto flex items-center justify-between gap-2 pt-3">
          <span className="text-lg font-semibold text-gray-900">
            {formatPrice(product.price)}
          </span>

          <button
            onClick={handleAddToCart}
            disabled={product.stock === 0}
            className="flex items-center gap-1.5 disabled:bg-gray-200 disabled:text-gray-400 text-white text-xs font-medium px-3 py-2 rounded-xl transition-colors hover:opacity-90"
            style={{ backgroundColor: product.stock === 0 ? undefined : accent }}
          >
            <ShoppingCart size={14} />
            <span className="hidden sm:inline">Add</span>
          </button>
        </div>
      </div>
    </div>
  )
}
