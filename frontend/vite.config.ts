import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    port: 2800,
    proxy: {
      '/api': {
        target: 'http://localhost:3800',
        changeOrigin: true,
      },
    },
  },
})
