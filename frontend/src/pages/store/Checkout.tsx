import { useState } from 'react'
import { Link } from '@tanstack/react-router'
import {
  Package,
  Tag,
  CheckCircle,
  ChevronRight,
  Loader2,
  ArrowLeft,
} from 'lucide-react'
import { useCart } from '../../lib/cart'
import { useStoreInfo } from '../../lib/store-hooks'

// ─── Types ───────────────────────────────────────────────────────────────────

interface ShippingForm {
  firstName: string
  lastName: string
  email: string
  phone: string
  street: string
  city: string
  state: string
  postalCode: string
  country: string
}

const EMPTY_FORM: ShippingForm = {
  firstName: '',
  lastName: '',
  email: '',
  phone: '',
  street: '',
  city: '',
  state: '',
  postalCode: '',
  country: 'US',
}

// ─── Helpers ─────────────────────────────────────────────────────────────────

function formatAmount(amount: number, currency = 'USD'): string {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency,
  }).format(amount)
}

const SHIPPING_COST = 7.99           // $7.99
const FREE_SHIPPING_THRESHOLD = 50   // $50.00

// ─── Input component ─────────────────────────────────────────────────────────

function Field({
  label,
  required,
  children,
}: {
  label: string
  required?: boolean
  children: React.ReactNode
}) {
  return (
    <div>
      <label className="block text-sm font-medium text-gray-700 mb-1.5">
        {label}
        {required && <span className="text-red-400 ml-0.5">*</span>}
      </label>
      {children}
    </div>
  )
}

const inputCls =
  'w-full border border-gray-200 rounded-xl px-3.5 py-2.5 text-sm text-gray-800 placeholder-gray-400 outline-none focus:border-gray-400 focus:ring-2 focus:ring-gray-100 transition-all'

// ─── Success state ────────────────────────────────────────────────────────────

function OrderSuccess({ orderNumber }: { orderNumber: string }) {
  return (
    <div className="min-h-screen bg-gray-50 flex items-center justify-center px-4">
      <div className="bg-white rounded-3xl border border-gray-100 shadow-sm p-10 text-center max-w-md w-full">
        <div className="w-20 h-20 bg-green-50 rounded-full flex items-center justify-center mx-auto mb-6">
          <CheckCircle size={40} className="text-green-500" />
        </div>
        <h1 className="text-2xl font-bold text-gray-900 mb-2">Order Placed!</h1>
        <p className="text-gray-500 mb-1 text-sm">Thank you for your order.</p>
        <p className="text-xs text-gray-400 mb-6">
          Order #{orderNumber}
        </p>
        <p className="text-sm text-gray-600 mb-8">
          We'll send you a confirmation email shortly with your order details and tracking
          information.
        </p>
        <Link
          to="/products"
          className="inline-flex items-center gap-2 bg-gray-900 hover:bg-gray-800 text-white font-semibold px-8 py-3.5 rounded-2xl transition-colors"
        >
          Continue Shopping
        </Link>
      </div>
    </div>
  )
}

// ─── Component ───────────────────────────────────────────────────────────────

export default function Checkout() {
  const { items, total, clearCart } = useCart()
  const { data: storeInfo } = useStoreInfo()
  const accent = storeInfo?.accent_color ?? '#6366f1'
  const [form, setForm] = useState<ShippingForm>(EMPTY_FORM)
  const [promoCode, setPromoCode] = useState('')
  const [promoApplied, setPromoApplied] = useState(false)
  const [promoError, setPromoError] = useState('')
  const [isPlacing, setIsPlacing] = useState(false)
  const [orderNumber, setOrderNumber] = useState('')
  const [errors, setErrors] = useState<Partial<ShippingForm>>({})

  const shippingCost = total >= FREE_SHIPPING_THRESHOLD ? 0 : SHIPPING_COST
  const discount = promoApplied ? total * 0.1 : 0
  const orderTotal = total + shippingCost - discount

  // ── Promo ────────────────────────────────────────────────────────────────
  const handleApplyPromo = () => {
    setPromoError('')
    if (promoCode.trim().toUpperCase() === 'HADA10') {
      setPromoApplied(true)
    } else {
      setPromoError('Invalid promo code. Try HADA10 for 10% off.')
    }
  }

  // ── Validation ───────────────────────────────────────────────────────────
  const validate = (): boolean => {
    const errs: Partial<ShippingForm> = {}
    if (!form.firstName) errs.firstName = 'Required'
    if (!form.lastName) errs.lastName = 'Required'
    if (!form.email || !form.email.includes('@')) errs.email = 'Valid email required'
    if (!form.street) errs.street = 'Required'
    if (!form.city) errs.city = 'Required'
    if (!form.postalCode) errs.postalCode = 'Required'
    setErrors(errs)
    return Object.keys(errs).length === 0
  }

  // ── Place order ──────────────────────────────────────────────────────────
  const handlePlaceOrder = async () => {
    if (!validate()) return
    setIsPlacing(true)
    // TODO: POST to /api/v1/orders once backend is ready
    await new Promise((resolve) => setTimeout(resolve, 1800))
    const num = `ORD-${Date.now().toString(36).toUpperCase()}`
    setOrderNumber(num)
    clearCart()
    setIsPlacing(false)
  }

  // ── Success ──────────────────────────────────────────────────────────────
  if (orderNumber) {
    return <OrderSuccess orderNumber={orderNumber} />
  }

  const update = (key: keyof ShippingForm) => (
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>
  ) => {
    setForm((prev) => ({ ...prev, [key]: e.target.value }))
    setErrors((prev) => ({ ...prev, [key]: undefined }))
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Header */}
        <div className="flex items-center gap-3 mb-8">
          <Link
            to="/cart"
            className="p-2 hover:bg-gray-100 rounded-xl text-gray-500 hover:text-gray-700 transition-colors"
          >
            <ArrowLeft size={18} />
          </Link>
          <h1 className="text-2xl sm:text-3xl font-bold text-gray-900">Checkout</h1>
        </div>

        <div className="grid lg:grid-cols-5 gap-8">
          {/* ── Shipping form ────────────────────────────────────────── */}
          <div className="lg:col-span-3 space-y-6">
            {/* Contact */}
            <div className="bg-white rounded-2xl border border-gray-100 shadow-sm p-6">
              <h2 className="font-semibold text-gray-900 mb-5 flex items-center gap-2">
                <span className="w-6 h-6 text-xs font-bold rounded-full flex items-center justify-center" style={{ backgroundColor: `${accent}15`, color: accent }}>
                  1
                </span>
                Contact Information
              </h2>
              <div className="grid sm:grid-cols-2 gap-4">
                <Field label="First name" required>
                  <input
                    className={inputCls}
                    value={form.firstName}
                    onChange={update('firstName')}
                    placeholder="Jane"
                  />
                  {errors.firstName && (
                    <p className="text-xs text-red-500 mt-1">{errors.firstName}</p>
                  )}
                </Field>
                <Field label="Last name" required>
                  <input
                    className={inputCls}
                    value={form.lastName}
                    onChange={update('lastName')}
                    placeholder="Doe"
                  />
                  {errors.lastName && (
                    <p className="text-xs text-red-500 mt-1">{errors.lastName}</p>
                  )}
                </Field>
                <Field label="Email" required>
                  <input
                    type="email"
                    className={inputCls}
                    value={form.email}
                    onChange={update('email')}
                    placeholder="jane@example.com"
                  />
                  {errors.email && (
                    <p className="text-xs text-red-500 mt-1">{errors.email}</p>
                  )}
                </Field>
                <Field label="Phone">
                  <input
                    type="tel"
                    className={inputCls}
                    value={form.phone}
                    onChange={update('phone')}
                    placeholder="+1 555 000 0000"
                  />
                </Field>
              </div>
            </div>

            {/* Shipping address */}
            <div className="bg-white rounded-2xl border border-gray-100 shadow-sm p-6">
              <h2 className="font-semibold text-gray-900 mb-5 flex items-center gap-2">
                <span className="w-6 h-6 text-xs font-bold rounded-full flex items-center justify-center" style={{ backgroundColor: `${accent}15`, color: accent }}>
                  2
                </span>
                Shipping Address
              </h2>
              <div className="space-y-4">
                <Field label="Street address" required>
                  <input
                    className={inputCls}
                    value={form.street}
                    onChange={update('street')}
                    placeholder="123 Main St, Apt 4B"
                  />
                  {errors.street && (
                    <p className="text-xs text-red-500 mt-1">{errors.street}</p>
                  )}
                </Field>
                <div className="grid sm:grid-cols-2 gap-4">
                  <Field label="City" required>
                    <input
                      className={inputCls}
                      value={form.city}
                      onChange={update('city')}
                      placeholder="New York"
                    />
                    {errors.city && (
                      <p className="text-xs text-red-500 mt-1">{errors.city}</p>
                    )}
                  </Field>
                  <Field label="State / Province">
                    <input
                      className={inputCls}
                      value={form.state}
                      onChange={update('state')}
                      placeholder="NY"
                    />
                  </Field>
                </div>
                <div className="grid sm:grid-cols-2 gap-4">
                  <Field label="Postal code" required>
                    <input
                      className={inputCls}
                      value={form.postalCode}
                      onChange={update('postalCode')}
                      placeholder="10001"
                    />
                    {errors.postalCode && (
                      <p className="text-xs text-red-500 mt-1">{errors.postalCode}</p>
                    )}
                  </Field>
                  <Field label="Country">
                    <select
                      className={inputCls}
                      value={form.country}
                      onChange={update('country')}
                    >
                      <option value="US">United States</option>
                      <option value="CA">Canada</option>
                      <option value="GB">United Kingdom</option>
                      <option value="AU">Australia</option>
                      <option value="DE">Germany</option>
                      <option value="FR">France</option>
                      <option value="JP">Japan</option>
                    </select>
                  </Field>
                </div>
              </div>
            </div>

            {/* Payment placeholder */}
            <div className="bg-white rounded-2xl border border-gray-100 shadow-sm p-6">
              <h2 className="font-semibold text-gray-900 mb-3 flex items-center gap-2">
                <span className="w-6 h-6 text-xs font-bold rounded-full flex items-center justify-center" style={{ backgroundColor: `${accent}15`, color: accent }}>
                  3
                </span>
                Payment
              </h2>
              <div className="bg-amber-50 border border-amber-100 rounded-xl p-4 text-sm text-amber-700">
                Payment integration coming soon — orders placed here are for demo purposes only.
              </div>
            </div>
          </div>

          {/* ── Order summary ────────────────────────────────────────── */}
          <div className="lg:col-span-2">
            <div className="bg-white rounded-2xl border border-gray-100 shadow-sm p-6 sticky top-24">
              <h2 className="font-semibold text-gray-900 mb-5">Order Summary</h2>

              {/* Items */}
              <div className="space-y-3 mb-5 max-h-60 overflow-y-auto pr-1">
                {items.map(({ product, quantity }) => (
                  <div key={product.id} className="flex items-center gap-3">
                    <div className="w-12 h-12 rounded-xl bg-gray-50 border border-gray-100 overflow-hidden shrink-0">
                      {product.images.length > 0 ? (
                        <img
                          src={product.images[0]}
                          alt={product.name}
                          className="w-full h-full object-cover"
                        />
                      ) : (
                        <div className="w-full h-full flex items-center justify-center text-gray-200">
                          <Package size={20} strokeWidth={1.5} />
                        </div>
                      )}
                    </div>
                    <div className="flex-1 min-w-0">
                      <p className="text-sm font-medium text-gray-800 line-clamp-1">
                        {product.name}
                      </p>
                      <p className="text-xs text-gray-400">Qty: {quantity}</p>
                    </div>
                    <span className="text-sm font-medium text-gray-800 shrink-0">
                      {formatAmount(product.price.amount * quantity, product.price.currency)}
                    </span>
                  </div>
                ))}
              </div>

              {/* Promo code */}
              <div className="mb-5">
                <label className="block text-xs font-medium text-gray-500 uppercase tracking-widest mb-2">
                  Promo Code
                </label>
                <div className="flex gap-2">
                  <div className="flex-1 flex items-center gap-2 border border-gray-200 rounded-xl px-3 focus-within:border-gray-400 transition-colors">
                    <Tag size={14} className="text-gray-400 shrink-0" />
                    <input
                      type="text"
                      placeholder="Enter code"
                      value={promoCode}
                      onChange={(e) => {
                        setPromoCode(e.target.value.toUpperCase())
                        setPromoError('')
                        setPromoApplied(false)
                      }}
                      className="flex-1 py-2.5 text-sm outline-none bg-transparent text-gray-700 placeholder-gray-400"
                      disabled={promoApplied}
                    />
                  </div>
                  <button
                    onClick={handleApplyPromo}
                    disabled={promoApplied || !promoCode.trim()}
                    className="px-4 py-2.5 bg-gray-100 hover:bg-gray-200 disabled:opacity-40 text-gray-700 text-sm font-medium rounded-xl transition-colors"
                  >
                    Apply
                  </button>
                </div>
                {promoApplied && (
                  <p className="text-xs text-green-600 mt-1.5 flex items-center gap-1">
                    <CheckCircle size={12} /> 10% discount applied!
                  </p>
                )}
                {promoError && (
                  <p className="text-xs text-red-500 mt-1.5">{promoError}</p>
                )}
              </div>

              {/* Totals */}
              <div className="space-y-2.5 mb-5 pt-4 border-t border-gray-100">
                <div className="flex justify-between text-sm">
                  <span className="text-gray-500">Subtotal</span>
                  <span className="font-medium text-gray-800">{formatAmount(total)}</span>
                </div>
                {promoApplied && (
                  <div className="flex justify-between text-sm">
                    <span className="text-green-600">Promo (HADA10)</span>
                    <span className="text-green-600 font-medium">
                      -{formatAmount(discount)}
                    </span>
                  </div>
                )}
                <div className="flex justify-between text-sm">
                  <span className="text-gray-500">Shipping</span>
                  {shippingCost === 0 ? (
                    <span className="text-green-600 font-medium">Free</span>
                  ) : (
                    <span className="font-medium text-gray-800">
                      {formatAmount(shippingCost)}
                    </span>
                  )}
                </div>
                <div className="flex justify-between font-bold text-base pt-2 border-t border-gray-100">
                  <span>Total</span>
                  <span>{formatAmount(orderTotal)}</span>
                </div>
              </div>

              {/* Place order button */}
              <button
                onClick={handlePlaceOrder}
                disabled={isPlacing || items.length === 0}
                className="flex items-center justify-center gap-2 w-full disabled:opacity-60 text-white font-semibold py-4 rounded-2xl transition-opacity hover:opacity-90 shadow-lg"
                style={{ backgroundColor: accent }}
              >
                {isPlacing ? (
                  <>
                    <Loader2 size={18} className="animate-spin" /> Placing order…
                  </>
                ) : (
                  <>
                    Place Order <ChevronRight size={18} />
                  </>
                )}
              </button>

              <p className="text-xs text-gray-400 text-center mt-3">
                By placing your order you agree to our terms of service.
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
