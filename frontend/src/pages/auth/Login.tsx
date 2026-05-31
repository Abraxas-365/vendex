import { useState } from 'react'
import { Store, ShoppingBag, BarChart3, Sparkles } from 'lucide-react'
import { useAuth } from '../../lib/auth'

function GoogleIcon() {
  return (
    <svg width="20" height="20" viewBox="0 0 18 18" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path d="M17.64 9.20455C17.64 8.56636 17.5827 7.95273 17.4764 7.36364H9V10.845H13.8436C13.635 11.97 13.0009 12.9232 12.0477 13.5614V15.8195H14.9564C16.6582 14.2527 17.64 11.9455 17.64 9.20455Z" fill="#4285F4" />
      <path d="M9 18C11.43 18 13.4673 17.1941 14.9564 15.8195L12.0477 13.5614C11.2418 14.1014 10.2109 14.4205 9 14.4205C6.65591 14.4205 4.67182 12.8373 3.96409 10.71H0.957275V13.0418C2.43818 15.9832 5.48182 18 9 18Z" fill="#34A853" />
      <path d="M3.96409 10.71C3.78409 10.17 3.68182 9.59318 3.68182 9C3.68182 8.40682 3.78409 7.83 3.96409 7.29V4.95818H0.957275C0.347727 6.17318 0 7.54773 0 9C0 10.4523 0.347727 11.8268 0.957275 13.0418L3.96409 10.71Z" fill="#FBBC05" />
      <path d="M9 3.57955C10.3214 3.57955 11.5077 4.03364 12.4405 4.92545L15.0218 2.34409C13.4632 0.891818 11.4259 0 9 0C5.48182 0 2.43818 2.01682 0.957275 4.95818L3.96409 7.29C4.67182 5.16273 6.65591 3.57955 9 3.57955Z" fill="#EA4335" />
    </svg>
  )
}

function MicrosoftIcon() {
  return (
    <svg width="20" height="20" viewBox="0 0 18 18" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path d="M0 0H8.57143V8.57143H0V0Z" fill="#F25022" />
      <path d="M9.42857 0H18V8.57143H9.42857V0Z" fill="#7FBA00" />
      <path d="M0 9.42857H8.57143V18H0V9.42857Z" fill="#00A4EF" />
      <path d="M9.42857 9.42857H18V18H9.42857V9.42857Z" fill="#FFB900" />
    </svg>
  )
}

function FeatureItem({ icon: Icon, title, desc }: { icon: React.ElementType; title: string; desc: string }) {
  return (
    <div className="flex items-start gap-3">
      <div className="mt-0.5 flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-white/10">
        <Icon className="h-4 w-4 text-indigo-300" />
      </div>
      <div>
        <p className="text-sm font-medium text-white">{title}</p>
        <p className="text-xs text-indigo-200/70">{desc}</p>
      </div>
    </div>
  )
}

export default function Login() {
  const { login } = useAuth()
  const [loadingProvider, setLoadingProvider] = useState<'GOOGLE' | 'MICROSOFT' | null>(null)
  const [error, setError] = useState<string | null>(null)

  const handleLogin = async (provider: 'GOOGLE' | 'MICROSOFT') => {
    setLoadingProvider(provider)
    setError(null)
    try {
      await login(provider)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to initiate login. Please try again.')
      setLoadingProvider(null)
    }
  }

  const isLoading = loadingProvider !== null

  return (
    <div className="flex min-h-screen">
      {/* Left panel — branding + features */}
      <div className="hidden lg:flex lg:w-[480px] xl:w-[520px] flex-col justify-between bg-gradient-to-br from-indigo-600 via-indigo-700 to-purple-800 p-10 relative overflow-hidden">
        {/* Subtle pattern overlay */}
        <div className="absolute inset-0 opacity-[0.07]" style={{
          backgroundImage: `radial-gradient(circle at 1px 1px, white 1px, transparent 0)`,
          backgroundSize: '24px 24px',
        }} />

        <div className="relative z-10">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-white/15 backdrop-blur-sm">
              <Store className="h-5 w-5 text-white" />
            </div>
            <span className="text-xl font-bold text-white tracking-tight">Vendex</span>
          </div>

          <div className="mt-16">
            <h1 className="text-3xl font-bold text-white leading-tight">
              Your AI-powered<br />e-commerce platform
            </h1>
            <p className="mt-3 text-sm text-indigo-200/80 leading-relaxed max-w-sm">
              Manage your entire store with an intelligent assistant that understands your brand, products, and customers.
            </p>
          </div>

          <div className="mt-10 space-y-5">
            <FeatureItem
              icon={Sparkles}
              title="AI Store Assistant"
              desc="92+ tools for products, orders, promos, and content"
            />
            <FeatureItem
              icon={ShoppingBag}
              title="Full Commerce Suite"
              desc="Products, orders, inventory, shipping, tax, and payments"
            />
            <FeatureItem
              icon={BarChart3}
              title="Smart Analytics"
              desc="Revenue dashboards, conversion funnels, and trend insights"
            />
          </div>
        </div>

        <div className="relative z-10">
          <p className="text-xs text-indigo-300/60">
            Trusted by merchants worldwide
          </p>
        </div>

        {/* Decorative gradient orbs */}
        <div className="absolute -bottom-32 -right-32 h-64 w-64 rounded-full bg-purple-500/20 blur-3xl" />
        <div className="absolute -top-16 -right-16 h-48 w-48 rounded-full bg-indigo-400/15 blur-3xl" />
      </div>

      {/* Right panel — login form */}
      <div className="flex flex-1 flex-col items-center justify-center bg-gray-50 px-6">
        <div className="w-full max-w-[380px]">
          {/* Mobile logo (hidden on desktop) */}
          <div className="mb-8 flex items-center justify-center gap-2 lg:hidden">
            <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-indigo-600">
              <Store className="h-4.5 w-4.5 text-white" />
            </div>
            <span className="text-lg font-bold text-gray-900 tracking-tight">Vendex</span>
          </div>

          <div className="mb-8">
            <h2 className="text-2xl font-bold text-gray-900 tracking-tight">Welcome back</h2>
            <p className="mt-1.5 text-sm text-gray-500">Sign in to your admin dashboard</p>
          </div>

          {error && (
            <div className="mb-5 rounded-lg border border-red-200 bg-red-50 px-4 py-3">
              <p className="text-sm text-red-700">{error}</p>
            </div>
          )}

          <div className="space-y-3">
            <button
              onClick={() => void handleLogin('GOOGLE')}
              disabled={isLoading}
              className="group flex w-full items-center justify-center gap-3 rounded-xl border border-gray-200 bg-white px-5 py-3 text-sm font-medium text-gray-700 shadow-sm transition-all hover:border-gray-300 hover:bg-gray-50 hover:shadow-md focus:outline-none focus:ring-2 focus:ring-indigo-500/20 focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
            >
              <GoogleIcon />
              <span>{loadingProvider === 'GOOGLE' ? 'Redirecting...' : 'Continue with Google'}</span>
            </button>

            <button
              onClick={() => void handleLogin('MICROSOFT')}
              disabled={isLoading}
              className="group flex w-full items-center justify-center gap-3 rounded-xl border border-gray-200 bg-white px-5 py-3 text-sm font-medium text-gray-700 shadow-sm transition-all hover:border-gray-300 hover:bg-gray-50 hover:shadow-md focus:outline-none focus:ring-2 focus:ring-indigo-500/20 focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
            >
              <MicrosoftIcon />
              <span>{loadingProvider === 'MICROSOFT' ? 'Redirecting...' : 'Continue with Microsoft'}</span>
            </button>
          </div>

          <div className="mt-8 flex items-center gap-3">
            <div className="h-px flex-1 bg-gray-200" />
            <span className="text-xs text-gray-400">Secure authentication</span>
            <div className="h-px flex-1 bg-gray-200" />
          </div>

          <p className="mt-6 text-center text-xs text-gray-400 leading-relaxed">
            By signing in, you agree to our{' '}
            <a href="#" className="text-gray-500 underline underline-offset-2 hover:text-gray-700">terms of service</a>
            {' '}and{' '}
            <a href="#" className="text-gray-500 underline underline-offset-2 hover:text-gray-700">privacy policy</a>.
          </p>
        </div>

        <p className="mt-12 text-xs text-gray-300">
          Vendex v1.0
        </p>
      </div>
    </div>
  )
}
