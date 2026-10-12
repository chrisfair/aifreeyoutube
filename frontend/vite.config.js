import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  server: {
    port: 5173,
    host: '0.0.0.0',
    proxy: {
      // Route all backend API requests to the Go HTTP server on port 8484
      '/api': {
        target: 'http://localhost:8484',
        changeOrigin: true,
      },
    },
  },
});
