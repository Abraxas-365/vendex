import { createRouter, RouterProvider } from '@tanstack/react-router'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { Theme } from '@radix-ui/themes'
import { CartProvider } from './lib/cart'
import { AuthProvider } from './lib/auth'
import { routeTree } from './routeTree'

// ─── TanStack Query client ────────────────────────────────────────────────────

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 30_000,
      retry: 1,
      refetchOnWindowFocus: false,
    },
  },
})

// ─── Router ───────────────────────────────────────────────────────────────────

const router = createRouter({ routeTree })

declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router
  }
}

// ─── App ─────────────────────────────────────────────────────────────────────
// AuthProvider wraps the whole app so auth state is accessible everywhere.
// CartProvider wraps the store so the cart is accessible from both
// store pages (via useCart) and the admin Navbar "View Store" link context.

export default function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <AuthProvider>
        <CartProvider>
          <Theme>
            <RouterProvider router={router} />
          </Theme>
        </CartProvider>
      </AuthProvider>
    </QueryClientProvider>
  )
}
