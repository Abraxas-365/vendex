import { useEffect, useRef, useState } from 'react'
import { Box, Flex, Heading, Text, Spinner } from '@radix-ui/themes'
import { handleAuthCallback, setTokens } from '../../lib/api'

// This page handles the OAuth callback redirect.
// The backend redirects here after OAuth completes.
// URL pattern: /auth/callback/:provider?code=...&state=...
//
// The backend may also set cookies (access_token / refresh_token).
// We attempt to parse URL params first, then fall back to reading cookies.

interface CallbackProps {
  provider: string
}

function getCookie(name: string): string | null {
  const match = document.cookie.match(new RegExp(`(?:^|; )${name}=([^;]*)`))
  return match ? decodeURIComponent(match[1]) : null
}

export default function Callback({ provider }: CallbackProps) {
  const [error, setError] = useState<string | null>(null)
  const attempted = useRef(false)

  useEffect(() => {
    if (attempted.current) return
    attempted.current = true

    const run = async () => {
      const params = new URLSearchParams(window.location.search)
      const code = params.get('code')
      const state = params.get('state')

      // Try query-param based flow (backend redirects with code+state)
      if (code && state) {
        try {
          const tokenResponse = await handleAuthCallback(provider, code, state)
          setTokens(tokenResponse.access_token, tokenResponse.refresh_token)
          window.location.replace('/admin')
          return
        } catch (err) {
          setError(err instanceof Error ? err.message : 'Authentication failed. Please try again.')
          return
        }
      }

      // Fallback: check if backend already set cookies with the tokens
      const cookieToken = getCookie('access_token')
      const cookieRefresh = getCookie('refresh_token')
      if (cookieToken) {
        setTokens(cookieToken, cookieRefresh ?? undefined)
        window.location.replace('/admin')
        return
      }

      // No tokens found
      setError('No authentication data received. Please try logging in again.')
    }

    void run()
  }, [provider])

  if (error) {
    return (
      <Box
        style={{
          minHeight: '100vh',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          background: '#f8fafc',
        }}
      >
        <Flex direction="column" align="center" gap="4" style={{ maxWidth: 400, padding: '0 16px' }}>
          <Heading size="5" style={{ color: '#dc2626' }}>
            Authentication Failed
          </Heading>
          <Text size="2" color="gray" align="center">
            {error}
          </Text>
          <a
            href="/login"
            style={{
              color: '#6366f1',
              fontWeight: 500,
              textDecoration: 'underline',
              fontSize: 14,
            }}
          >
            Back to Login
          </a>
        </Flex>
      </Box>
    )
  }

  return (
    <Box
      style={{
        minHeight: '100vh',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        background: '#f8fafc',
      }}
    >
      <Flex direction="column" align="center" gap="4">
        <Spinner size="3" />
        <Text size="2" color="gray">
          Completing sign in…
        </Text>
      </Flex>
    </Box>
  )
}
