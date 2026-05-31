import { Link } from '@tanstack/react-router'
import {
  ShoppingCart,
  Trash2,
  Minus,
  Plus,
  ArrowRight,
  Package,
  ShoppingBag,
} from 'lucide-react'
import { useCart } from '../../lib/cart'
import type { Money } from '../../types/index'

// ─── Helpers ─────────────────────────────────────────────────────────────────

function formatMoney(price: Money): string {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: price.currency ?? 'USD',
  }).format(price.amount / 100)
}

function formatAmount(amount: number, currency = 'USD'): string {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency,
  }).format(amount)
}

const FREE_SHIPPING_THRESHOLD = 50   // $50.00
const SHIPPING_COST = 7.99           // $7.99

// ─── Component ───────────────────────────────────────────────────────────────

export default function Cart() {
  const { items, removeItem, updateQuantity, clearCart, total, itemCount } = useCart()

  const shippingCost = total >= FREE_SHIPPING_THRESHOLD ? 0 : SHIPPING_COST
  const orderTotal = total + shippingCost
  const freeShippingProgress = Math.min((total / FREE_SHIPPING_THRESHOLD) * 100, 100)
  const remainingForFreeShipping = FREE_SHIPPING_THRESHOLD - total

  // ── Empty state ──────────────────────────────────────────────────────────
  if (items.length === 0) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center px-4">
        <div className="text-center max-w-sm">
          <div className="w-24 h-24 bg-indigo-50 rounded-full flex items-center justify-center mx-auto mb-6">
            <ShoppingCart size={40} className="text-indigo-400" />
          </div>
          <h1 className="text-2xl font-bold text-gray-900 mb-2">Your cart is empty</h1>
          <p className="text-gray-500 mb-8">
            Looks like you haven't added anything yet. Let's change that!
          </p>
          <Link
            to="/products"
            className="inline-flex items-center gap-2 bg-indigo-600 hover:bg-indigo-700 text-white font-semibold px-8 py-3.5 rounded-2xl transition-colors shadow-lg shadow-indigo-200"
          >
            <ShoppingBag size={18} />
            Continue Shopping
          </Link>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <h1 className="text-2xl sm:text-3xl font-bold text-gray-900 mb-2">
          Shopping Cart
        </h1>
        <p className="text-gray-500 text-sm mb-8">
          {itemCount} item{itemCount !== 1 ? 's' : ''} in your cart
        </p>

        {/* Free shipping progress */}
        {total < FREE_SHIPPING_THRESHOLD && (
          <div className="bg-indigo-50 border border-indigo-100 rounded-2xl p-4 mb-6">
            <p className="text-sm text-indigo-700 mb-2">
              Add{' '}
              <span className="font-semibold">{formatAmount(remainingForFreeShipping)}</span>{' '}
              more for free shipping!
            </p>
            <div className="h-2 bg-indigo-100 rounded-full overflow-hidden">
              <div
                className="h-full bg-indigo-500 rounded-full transition-all duration-500"
                style={{ width: `${freeShippingProgress}%` }}
              />
            </div>
          </div>
        )}
        {total >= FREE_SHIPPING_THRESHOLD && (
          <div className="bg-green-50 border border-green-100 rounded-2xl p-4 mb-6">
            <p className="text-sm text-green-700 font-medium">
              🎉 You've unlocked free shipping!
            </p>
          </div>
        )}

        <div className="grid lg:grid-cols-3 gap-8">
          {/* ── Cart items ───────────────────────────────────────────── */}
          <div className="lg:col-span-2 space-y-3">
            {/* Clear all */}
            <div className="flex justify-end">
              <button
                onClick={clearCart}
                className="text-xs text-gray-400 hover:text-red-500 transition-colors"
              >
                Clear all
              </button>
            </div>

            {items.map(({ product, quantity }) => {
              const lineTotal = product.price.amount * quantity

              return (
                <div
                  key={product.id}
                  className="bg-white rounded-2xl border border-gray-100 p-4 flex gap-4 shadow-sm"
                >
                  {/* Image */}
                  <div className="w-20 h-20 sm:w-24 sm:h-24 rounded-xl bg-gray-50 border border-gray-100 overflow-hidden shrink-0">
                    {product.images.length > 0 ? (
                      <img
                        src={product.images[0]}
                        alt={product.name}
                        className="w-full h-full object-cover"
                      />
                    ) : (
                      <div className="w-full h-full flex items-center justify-center text-gray-200">
                        <Package size={32} strokeWidth={1} />
                      </div>
                    )}
                  </div>

                  {/* Info */}
                  <div className="flex-1 min-w-0">
                    <div className="flex items-start justify-between gap-2">
                      <div>
                        <p className="font-medium text-gray-800 text-sm line-clamp-2">
                          {product.name}
                        </p>
                        {product.sku && (
                          <p className="text-xs text-gray-400 mt-0.5">SKU: {product.sku}</p>
                        )}
                      </div>

                      {/* Remove */}
                      <button
                        onClick={() => removeItem(product.id)}
                        className="p-1.5 text-gray-300 hover:text-red-500 hover:bg-red-50 rounded-lg transition-colors shrink-0"
                        aria-label="Remove item"
                      >
                        <Trash2 size={15} />
                      </button>
                    </div>

                    <div className="flex items-center justify-between gap-4 mt-3">
                      {/* Quantity controls */}
                      <div className="flex items-center gap-1 bg-gray-50 border border-gray-100 rounded-xl overflow-hidden">
                        <button
                          onClick={() => updateQuantity(product.id, quantity - 1)}
                          className="px-2.5 py-1.5 hover:bg-gray-100 text-gray-500 transition-colors"
                          aria-label="Decrease"
                        >
                          <Minus size={13} />
                        </button>
                        <span className="w-8 text-center text-sm font-semibold text-gray-800">
                          {quantity}
                        </span>
                        <button
                          onClick={() => updateQuantity(product.id, quantity + 1)}
                          className="px-2.5 py-1.5 hover:bg-gray-100 text-gray-500 transition-colors"
                          aria-label="Increase"
                        >
                          <Plus size={13} />
                        </button>
                      </div>

                      {/* Price */}
                      <div className="text-right">
                        <p className="font-semibold text-gray-900 text-sm">
                          {formatAmount(lineTotal, product.price.currency)}
                        </p>
                        {quantity > 1 && (
                          <p className="text-xs text-gray-400">
                            {formatMoney(product.price)} each
                          </p>
                        )}
                      </div>
                    </div>
                  </div>
                </div>
              )
            })}
          </div>

          {/* ── Order summary ────────────────────────────────────────── */}
          <div className="lg:col-span-1">
            <div className="bg-white rounded-2xl border border-gray-100 shadow-sm p-6 sticky top-24">
              <h2 className="font-semibold text-gray-900 text-lg mb-5">Order Summary</h2>

              <div className="space-y-3 mb-5">
                <div className="flex justify-between text-sm">
                  <span className="text-gray-500">
                    Subtotal ({itemCount} item{itemCount !== 1 ? 's' : ''})
                  </span>
                  <span className="font-medium text-gray-800">{formatAmount(total)}</span>
                </div>
                <div className="flex justify-between text-sm">
                  <span className="text-gray-500">Shipping</span>
                  {shippingCost === 0 ? (
                    <span className="text-green-600 font-medium">Free</span>
                  ) : (
                    <span className="font-medium text-gray-800">{formatAmount(shippingCost)}</span>
                  )}
                </div>
                <div className="flex justify-between text-sm">
                  <span className="text-gray-500">Taxes</span>
                  <span className="text-gray-400 text-xs">Calculated at checkout</span>
                </div>
              </div>

              <div className="border-t border-gray-100 pt-4 mb-6">
                <div className="flex justify-between">
                  <span className="font-semibold text-gray-900">Total</span>
                  <span className="font-bold text-xl text-gray-900">
                    {formatAmount(orderTotal)}
                  </span>
                </div>
              </div>

              <Link
                to="/checkout"
                className="flex items-center justify-center gap-2 w-full bg-indigo-600 hover:bg-indigo-700 text-white font-semibold py-4 rounded-2xl transition-colors shadow-lg shadow-indigo-200 mb-3"
              >
                Proceed to Checkout <ArrowRight size={18} />
              </Link>

              <Link
                to="/products"
                className="flex items-center justify-center w-full text-sm text-gray-500 hover:text-indigo-600 py-2 transition-colors"
              >
                Continue Shopping
              </Link>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
