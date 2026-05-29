import { useState, useEffect } from 'react'
import { Save, Loader2, CheckCircle2 } from 'lucide-react'
import { useSettings, useUpdateSettings } from '../../lib/hooks'
import type { StoreSettings, StoreAddress, SocialLinks, CheckoutConfig } from '../../types'

function SectionCard({
  title,
  children,
}: {
  title: string
  children: React.ReactNode
}) {
  return (
    <div className="rounded-xl border border-gray-200 bg-white shadow-sm">
      <div className="border-b border-gray-200 px-6 py-4">
        <h2 className="text-base font-semibold text-gray-900">{title}</h2>
      </div>
      <div className="space-y-4 px-6 py-5">{children}</div>
    </div>
  )
}

function Field({
  label,
  children,
}: {
  label: string
  children: React.ReactNode
}) {
  return (
    <div>
      <label className="mb-1 block text-sm font-medium text-gray-700">{label}</label>
      {children}
    </div>
  )
}

const inputClass =
  'w-full rounded-lg border border-gray-300 px-3 py-2 text-sm text-gray-900 placeholder-gray-400 focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500'

const CURRENCIES = ['USD', 'EUR', 'GBP', 'BRL', 'CAD', 'AUD', 'JPY', 'MXN']
const TIMEZONES = [
  'UTC',
  'America/New_York',
  'America/Chicago',
  'America/Denver',
  'America/Los_Angeles',
  'America/Sao_Paulo',
  'America/Mexico_City',
  'Europe/London',
  'Europe/Paris',
  'Europe/Berlin',
  'Asia/Tokyo',
  'Asia/Shanghai',
  'Australia/Sydney',
]

export default function Settings() {
  const settings = useSettings()
  const updateSettings = useUpdateSettings()

  const [storeName, setStoreName] = useState('')
  const [storeEmail, setStoreEmail] = useState('')
  const [storePhone, setStorePhone] = useState('')
  const [currency, setCurrency] = useState('USD')
  const [timezone, setTimezone] = useState('UTC')
  const [logoUrl, setLogoUrl] = useState('')
  const [faviconUrl, setFaviconUrl] = useState('')
  const [address, setAddress] = useState<StoreAddress>({
    street: '',
    city: '',
    state: '',
    country: '',
    zip: '',
  })
  const [social, setSocial] = useState<SocialLinks>({
    instagram: '',
    twitter: '',
    facebook: '',
  })
  const [checkout, setCheckout] = useState<CheckoutConfig>({
    guest_checkout: true,
    require_phone: false,
  })

  const [saved, setSaved] = useState(false)

  // Populate form when data loads
  useEffect(() => {
    if (settings.data) {
      const d = settings.data
      setStoreName(d.store_name)
      setStoreEmail(d.store_email)
      setStorePhone(d.store_phone)
      setCurrency(d.currency)
      setTimezone(d.timezone)
      setLogoUrl(d.logo_url)
      setFaviconUrl(d.favicon_url)
      if (d.address) setAddress(d.address)
      if (d.social_links) setSocial(d.social_links)
      if (d.checkout_config) setCheckout(d.checkout_config)
    }
  }, [settings.data])

  function handleSave() {
    const payload: Partial<StoreSettings> = {
      store_name: storeName,
      store_email: storeEmail,
      store_phone: storePhone,
      currency,
      timezone,
      logo_url: logoUrl,
      favicon_url: faviconUrl,
      address,
      social_links: social,
      checkout_config: checkout,
    }
    updateSettings.mutate(payload, {
      onSuccess: () => {
        setSaved(true)
        setTimeout(() => setSaved(false), 2000)
      },
    })
  }

  if (settings.isLoading) {
    return (
      <div className="flex h-64 items-center justify-center">
        <Loader2 className="h-6 w-6 animate-spin text-gray-400" />
      </div>
    )
  }

  return (
    <div className="space-y-8">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Settings</h1>
          <p className="mt-1 text-sm text-gray-500">
            Configure your store settings and preferences.
          </p>
        </div>
        <button
          onClick={handleSave}
          disabled={updateSettings.isPending}
          className="inline-flex items-center gap-2 rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white shadow-sm transition-colors hover:bg-indigo-700 disabled:opacity-50"
        >
          {updateSettings.isPending ? (
            <Loader2 className="h-4 w-4 animate-spin" />
          ) : saved ? (
            <CheckCircle2 className="h-4 w-4" />
          ) : (
            <Save className="h-4 w-4" />
          )}
          {saved ? 'Saved!' : 'Save Changes'}
        </button>
      </div>

      {updateSettings.isError && (
        <div className="rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
          Failed to save settings: {updateSettings.error.message}
        </div>
      )}

      {/* General */}
      <SectionCard title="General">
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <Field label="Store Name">
            <input
              className={inputClass}
              value={storeName}
              onChange={(e) => setStoreName(e.target.value)}
              placeholder="My Store"
            />
          </Field>
          <Field label="Store Email">
            <input
              type="email"
              className={inputClass}
              value={storeEmail}
              onChange={(e) => setStoreEmail(e.target.value)}
              placeholder="hello@mystore.com"
            />
          </Field>
          <Field label="Store Phone">
            <input
              className={inputClass}
              value={storePhone}
              onChange={(e) => setStorePhone(e.target.value)}
              placeholder="+1 (555) 000-0000"
            />
          </Field>
        </div>
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <Field label="Logo URL">
            <input
              className={inputClass}
              value={logoUrl}
              onChange={(e) => setLogoUrl(e.target.value)}
              placeholder="https://..."
            />
          </Field>
          <Field label="Favicon URL">
            <input
              className={inputClass}
              value={faviconUrl}
              onChange={(e) => setFaviconUrl(e.target.value)}
              placeholder="https://..."
            />
          </Field>
        </div>
      </SectionCard>

      {/* Localization */}
      <SectionCard title="Localization">
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <Field label="Currency">
            <select
              className={inputClass}
              value={currency}
              onChange={(e) => setCurrency(e.target.value)}
            >
              {CURRENCIES.map((c) => (
                <option key={c} value={c}>
                  {c}
                </option>
              ))}
            </select>
          </Field>
          <Field label="Timezone">
            <select
              className={inputClass}
              value={timezone}
              onChange={(e) => setTimezone(e.target.value)}
            >
              {TIMEZONES.map((tz) => (
                <option key={tz} value={tz}>
                  {tz}
                </option>
              ))}
            </select>
          </Field>
        </div>
      </SectionCard>

      {/* Address */}
      <SectionCard title="Store Address">
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <div className="sm:col-span-2">
            <Field label="Street">
              <input
                className={inputClass}
                value={address.street}
                onChange={(e) => setAddress({ ...address, street: e.target.value })}
                placeholder="123 Main St"
              />
            </Field>
          </div>
          <Field label="City">
            <input
              className={inputClass}
              value={address.city}
              onChange={(e) => setAddress({ ...address, city: e.target.value })}
              placeholder="San Francisco"
            />
          </Field>
          <Field label="State / Province">
            <input
              className={inputClass}
              value={address.state}
              onChange={(e) => setAddress({ ...address, state: e.target.value })}
              placeholder="CA"
            />
          </Field>
          <Field label="Country">
            <input
              className={inputClass}
              value={address.country}
              onChange={(e) => setAddress({ ...address, country: e.target.value })}
              placeholder="US"
            />
          </Field>
          <Field label="ZIP / Postal Code">
            <input
              className={inputClass}
              value={address.zip}
              onChange={(e) => setAddress({ ...address, zip: e.target.value })}
              placeholder="94102"
            />
          </Field>
        </div>
      </SectionCard>

      {/* Social Links */}
      <SectionCard title="Social Links">
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-3">
          <Field label="Instagram">
            <input
              className={inputClass}
              value={social.instagram}
              onChange={(e) => setSocial({ ...social, instagram: e.target.value })}
              placeholder="https://instagram.com/..."
            />
          </Field>
          <Field label="Twitter / X">
            <input
              className={inputClass}
              value={social.twitter}
              onChange={(e) => setSocial({ ...social, twitter: e.target.value })}
              placeholder="https://twitter.com/..."
            />
          </Field>
          <Field label="Facebook">
            <input
              className={inputClass}
              value={social.facebook}
              onChange={(e) => setSocial({ ...social, facebook: e.target.value })}
              placeholder="https://facebook.com/..."
            />
          </Field>
        </div>
      </SectionCard>

      {/* Checkout */}
      <SectionCard title="Checkout">
        <div className="space-y-4">
          <label className="flex items-center gap-3">
            <input
              type="checkbox"
              checked={checkout.guest_checkout}
              onChange={(e) =>
                setCheckout({ ...checkout, guest_checkout: e.target.checked })
              }
              className="h-4 w-4 rounded border-gray-300 text-indigo-600 focus:ring-indigo-500"
            />
            <div>
              <span className="text-sm font-medium text-gray-900">Allow Guest Checkout</span>
              <p className="text-xs text-gray-500">
                Customers can check out without creating an account.
              </p>
            </div>
          </label>
          <label className="flex items-center gap-3">
            <input
              type="checkbox"
              checked={checkout.require_phone}
              onChange={(e) =>
                setCheckout({ ...checkout, require_phone: e.target.checked })
              }
              className="h-4 w-4 rounded border-gray-300 text-indigo-600 focus:ring-indigo-500"
            />
            <div>
              <span className="text-sm font-medium text-gray-900">Require Phone Number</span>
              <p className="text-xs text-gray-500">
                Customers must provide a phone number at checkout.
              </p>
            </div>
          </label>
        </div>
      </SectionCard>
    </div>
  )
}
