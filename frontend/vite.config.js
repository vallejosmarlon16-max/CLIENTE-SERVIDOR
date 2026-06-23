import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    port: 3000,
    host: true, // Permite mapear el host en contenedores Docker
    watch: {
      usePolling: true // Asegura que los cambios en archivos se recarguen en Docker
    }
  }
})
