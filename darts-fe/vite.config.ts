import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    port: 3000,
    // Optional: if you want to call /api/* instead of http://localhost:8080
    // proxy: { '/api': 'http://localhost:8080' }
  }
})