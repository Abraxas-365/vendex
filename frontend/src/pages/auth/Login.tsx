import { useState } from 'react'
import { Box, Card, Flex, Heading, Text, Separator } from '@radix-ui/themes'
import { Store } from 'lucide-react'
import { useAuth } from '../../lib/auth'

// Google SVG icon
function GoogleIcon() {
  return (
    <svg width="20" height="20" viewBox="0 0 18 18" fill="none" xmlns="http://www.w3.org/2000/svg" style={{ flexShrink: 0 }}>
      <path
        d="M17.64 9.20455C17.64 8.56636 17.5827 7.95273 17.4764 7.36364H9V10.845H13.8436C13.635 11.97 13.0009 12.9232 12.0477 13.5614V15.8195H14.9564C16.6582 14.2527 17.64 11.9455 17.64 9.20455Z"
        fill="#4285F4"
      />
      <path
        d="M9 18C11.43 18 13.4673 17.1941 14.9564 15.8195L12.0477 13.5614C11.2418 14.1014 10.2109 14.4205 9 14.4205C6.65591 14.4205 4.67182 12.8373 3.96409 10.71H0.957275V13.0418C2.43818 15.9832 5.48182 18 9 18Z"
        fill="#34A853"
      />
      <path
        d="M3.96409 10.71C3.78409 10.17 3.68182 9.59318 3.68182 9C3.68182 8.40682 3.78409 7.83 3.96409 7.29V4.95818H0.957275C0.347727 6.17318 0 7.54773 0 9C0 10.4523 0.347727 11.8268 0.957275 13.0418L3.96409 10.71Z"
        fill="#FBBC05"
      />
      <path
        d="M9 3.57955C10.3214 3.57955 11.5077 4.03364 12.4405 4.92545L15.0218 2.34409C13.4632 0.891818 11.4259 0 9 0C5.48182 0 2.43818 2.01682 0.957275 4.95818L3.96409 7.29C4.67182 5.16273 6.65591 3.57955 9 3.57955Z"
        fill="#EA4335"
      />
    </svg>
  )
}

// Microsoft SVG icon
function MicrosoftIcon() {
  return (
    <svg width="20" height="20" viewBox="0 0 18 18" fill="none" xmlns="http://www.w3.org/2000/svg" style={{ flexShrink: 0 }}>
      <path d="M0 0H8.57143V8.57143H0V0Z" fill="#F25022" />
      <path d="M9.42857 0H18V8.57143H9.42857V0Z" fill="#7FBA00" />
      <path d="M0 9.42857H8.57143V18H0V9.42857Z" fill="#00A4EF" />
      <path d="M9.42857 9.42857H18V18H9.42857V9.42857Z" fill="#FFB900" />
    </svg>
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
      // login() redirects the browser — this line is only reached on error
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to initiate login. Please try again.')
      setLoadingProvider(null)
    }
  }

  const isLoading = loadingProvider !== null

  return (
    <Box
      style={{
        minHeight: '100vh',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        background: 'linear-gradient(135deg, #f8fafc 0%, #e2e8f0 100%)',
      }}
    >
      <Box style={{ width: '100%', maxWidth: 420, padding: '0 16px' }}>
        <Card
          size="4"
          style={{
            boxShadow: '0 8px 32px rgba(0,0,0,0.10), 0 1.5px 6px rgba(0,0,0,0.06)',
            borderRadius: 16,
          }}
        >
          <Flex direction="column" align="center" gap="6" style={{ padding: '8px 0' }}>
            {/* Branding */}
            <Flex direction="column" align="center" gap="3">
              <Box
                style={{
                  width: 60,
                  height: 60,
                  borderRadius: 16,
                  background: 'linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%)',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  boxShadow: '0 4px 14px rgba(99,102,241,0.35)',
                }}
              >
                <Store size={30} color="white" strokeWidth={2} />
              </Box>
              <Flex direction="column" align="center" gap="1">
                <Heading size="6" align="center" weight="bold" style={{ color: '#1e293b', letterSpacing: '-0.02em' }}>
                  Hada Commerce
                </Heading>
                <Text size="2" color="gray" align="center" style={{ maxWidth: 260 }}>
                  Sign in to access your admin dashboard
                </Text>
              </Flex>
            </Flex>

            {/* Error message */}
            {error && (
              <Box
                style={{
                  width: '100%',
                  padding: '10px 14px',
                  background: '#fef2f2',
                  border: '1px solid #fecaca',
                  borderRadius: 8,
                }}
              >
                <Text size="2" style={{ color: '#dc2626' }}>
                  {error}
                </Text>
              </Box>
            )}

            {/* OAuth buttons */}
            <Flex direction="column" gap="3" style={{ width: '100%' }}>
              {/* Google */}
              <button
                onClick={() => void handleLogin('GOOGLE')}
                disabled={isLoading}
                style={{
                  width: '100%',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  gap: 12,
                  padding: '12px 20px',
                  borderRadius: 10,
                  border: '1.5px solid #e2e8f0',
                  background: isLoading && loadingProvider === 'GOOGLE' ? '#f8fafc' : '#ffffff',
                  cursor: isLoading ? 'not-allowed' : 'pointer',
                  opacity: isLoading && loadingProvider !== 'GOOGLE' ? 0.6 : 1,
                  fontFamily: 'inherit',
                  fontSize: 15,
                  fontWeight: 500,
                  color: '#1e293b',
                  transition: 'background 0.15s, border-color 0.15s, box-shadow 0.15s',
                  boxShadow: '0 1px 3px rgba(0,0,0,0.06)',
                  outline: 'none',
                }}
                onMouseEnter={e => {
                  if (!isLoading) {
                    const btn = e.currentTarget
                    btn.style.background = '#f8fafc'
                    btn.style.borderColor = '#c7d2fe'
                    btn.style.boxShadow = '0 2px 8px rgba(99,102,241,0.10)'
                  }
                }}
                onMouseLeave={e => {
                  const btn = e.currentTarget
                  btn.style.background = isLoading && loadingProvider === 'GOOGLE' ? '#f8fafc' : '#ffffff'
                  btn.style.borderColor = '#e2e8f0'
                  btn.style.boxShadow = '0 1px 3px rgba(0,0,0,0.06)'
                }}
              >
                <GoogleIcon />
                <span>{loadingProvider === 'GOOGLE' ? 'Redirecting…' : 'Continue with Google'}</span>
              </button>

              {/* Microsoft */}
              <button
                onClick={() => void handleLogin('MICROSOFT')}
                disabled={isLoading}
                style={{
                  width: '100%',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  gap: 12,
                  padding: '12px 20px',
                  borderRadius: 10,
                  border: '1.5px solid #e2e8f0',
                  background: isLoading && loadingProvider === 'MICROSOFT' ? '#f8fafc' : '#ffffff',
                  cursor: isLoading ? 'not-allowed' : 'pointer',
                  opacity: isLoading && loadingProvider !== 'MICROSOFT' ? 0.6 : 1,
                  fontFamily: 'inherit',
                  fontSize: 15,
                  fontWeight: 500,
                  color: '#1e293b',
                  transition: 'background 0.15s, border-color 0.15s, box-shadow 0.15s',
                  boxShadow: '0 1px 3px rgba(0,0,0,0.06)',
                  outline: 'none',
                }}
                onMouseEnter={e => {
                  if (!isLoading) {
                    const btn = e.currentTarget
                    btn.style.background = '#f8fafc'
                    btn.style.borderColor = '#c7d2fe'
                    btn.style.boxShadow = '0 2px 8px rgba(99,102,241,0.10)'
                  }
                }}
                onMouseLeave={e => {
                  const btn = e.currentTarget
                  btn.style.background = isLoading && loadingProvider === 'MICROSOFT' ? '#f8fafc' : '#ffffff'
                  btn.style.borderColor = '#e2e8f0'
                  btn.style.boxShadow = '0 1px 3px rgba(0,0,0,0.06)'
                }}
              >
                <MicrosoftIcon />
                <span>{loadingProvider === 'MICROSOFT' ? 'Redirecting…' : 'Continue with Microsoft'}</span>
              </button>
            </Flex>

            <Separator size="4" style={{ opacity: 0.5 }} />

            {/* Disclaimer */}
            <Text size="1" color="gray" align="center" style={{ lineHeight: 1.6, maxWidth: 300 }}>
              By signing in, you agree to our{' '}
              <span style={{ textDecoration: 'underline', cursor: 'pointer' }}>terms of service</span>
              {' '}and{' '}
              <span style={{ textDecoration: 'underline', cursor: 'pointer' }}>privacy policy</span>.
            </Text>
          </Flex>
        </Card>
      </Box>
    </Box>
  )
}
