import React, { createContext, useContext, useEffect, useState, useCallback } from 'react'
import type { AuthUser, AuthTenant } from '../types'
import { getMe, logout as apiLogout, clearTokens, getAccessToken, setTokens, setTenantId, initiateLogin as apiInitiateLogin } from './api'

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

export interface AuthContextValue {
  user: AuthUser | null
  tenant: AuthTenant | null
  isAuthenticated: boolean
  isLoading: boolean
  login: (provider: 'GOOGLE' | 'MICROSOFT') => Promise<void>
  logout: () => Promise<void>
}

// ---------------------------------------------------------------------------
// Context
// ---------------------------------------------------------------------------

const AuthContext = createContext<AuthContextValue | null>(null)

// ---------------------------------------------------------------------------
// AuthProvider
// ---------------------------------------------------------------------------

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<AuthUser | null>(null)
  const [tenant, setTenant] = useState<AuthTenant | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  // On mount, check if we have a valid token and fetch the current user
  useEffect(() => {
    const token = getAccessToken()
    if (!token) {
      setIsLoading(false)
      return
    }

    getMe()
      .then(({ user: u, tenant: t }) => {
        setUser(u)
        setTenant(t)
        if (t?.id) setTenantId(t.id)
      })
      .catch(() => {
        // Token invalid or expired and refresh failed — clear everything
        clearTokens()
        setUser(null)
        setTenant(null)
      })
      .finally(() => {
        setIsLoading(false)
      })
  }, [])

  const login = useCallback(async (provider: 'GOOGLE' | 'MICROSOFT') => {
    const { auth_url } = await apiInitiateLogin(provider)
    // Redirect browser to OAuth provider
    window.location.href = auth_url
  }, [])

  const logout = useCallback(async () => {
    try {
      await apiLogout()
    } catch {
      // Ignore errors — we clear local state regardless
    }
    clearTokens()
    setUser(null)
    setTenant(null)
  }, [])

  return (
    <AuthContext.Provider
      value={{
        user,
        tenant,
        isAuthenticated: Boolean(user),
        isLoading,
        login,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

// ---------------------------------------------------------------------------
// useAuth hook
// ---------------------------------------------------------------------------

export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext)
  if (!ctx) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return ctx
}

// ---------------------------------------------------------------------------
// Helper: store tokens after OAuth callback
// ---------------------------------------------------------------------------

export function storeAuthTokens(accessToken: string, refreshToken?: string): void {
  setTokens(accessToken, refreshToken)
}
