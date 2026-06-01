import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

export default defineConfig({
  plugins: [react(), tailwindcss()],
  server: {
    port: 5173,
    proxy: {
      '/api/v1': {
        target: 'http://localhost:8080',
      },
      '/auth': {
        target: 'http://localhost:8080',
      },
      '/store': {
        target: 'http://localhost:8080',
        // Preserve the original Host header (e.g. fashion.localhost:5174)
        // so the backend tenant resolver can extract the subdomain
        changeOrigin: false,
      },
    },
  },
})
