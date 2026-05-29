import { useState } from 'react'
import { Box, Card, Flex, Heading, Text, Button } from '@radix-ui/themes'
import { Store } from 'lucide-react'
import { useAuth } from '../../lib/auth'

// Google and Microsoft SVG icons as inline components
function GoogleIcon() {
  return (
    <svg width="18" height="18" viewBox="0 0 18 18" fill="none" xmlns="http://www.w3.org/2000/svg">
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

function MicrosoftIcon() {
  return (
    <svg width="18" height="18" viewBox="0 0 18 18" fill="none" xmlns="http://www.w3.org/2000/svg">
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
      <Box style={{ width: '100%', maxWidth: 400, padding: '0 16px' }}>
        <Card size="4" style={{ boxShadow: '0 4px 24px rgba(0,0,0,0.08)' }}>
          <Flex direction="column" align="center" gap="6">
            {/* Logo area */}
            <Flex direction="column" align="center" gap="2">
              <Box
                style={{
                  width: 56,
                  height: 56,
                  borderRadius: 14,
                  background: 'linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%)',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                }}
              >
                <Store size={28} color="white" />
              </Box>
              <Heading size="5" align="center" style={{ color: '#1e293b' }}>
                Hada Commerce
              </Heading>
              <Text size="2" color="gray" align="center">
                Sign in to access your admin dashboard
              </Text>
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
              <Button
                size="3"
                variant="outline"
                onClick={() => void handleLogin('GOOGLE')}
                disabled={loadingProvider !== null}
                style={{
                  width: '100%',
                  cursor: loadingProvider !== null ? 'not-allowed' : 'pointer',
                  justifyContent: 'center',
                  gap: 10,
                  fontWeight: 500,
                  color: '#1e293b',
                  borderColor: '#e2e8f0',
                }}
              >
                <GoogleIcon />
                {loadingProvider === 'GOOGLE' ? 'Redirecting…' : 'Continue with Google'}
              </Button>

              <Button
                size="3"
                variant="outline"
                onClick={() => void handleLogin('MICROSOFT')}
                disabled={loadingProvider !== null}
                style={{
                  width: '100%',
                  cursor: loadingProvider !== null ? 'not-allowed' : 'pointer',
                  justifyContent: 'center',
                  gap: 10,
                  fontWeight: 500,
                  color: '#1e293b',
                  borderColor: '#e2e8f0',
                }}
              >
                <MicrosoftIcon />
                {loadingProvider === 'MICROSOFT' ? 'Redirecting…' : 'Continue with Microsoft'}
              </Button>
            </Flex>

            <Text size="1" color="gray" align="center">
              By signing in, you agree to our terms of service and privacy policy.
            </Text>
          </Flex>
        </Card>
      </Box>
    </Box>
  )
}
