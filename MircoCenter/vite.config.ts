import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    // 设置代理
    proxy: {
      '/api': {
        target: 'http://127.0.0.1:8080/wcenter/ServiceCenter/api/v1',
        changeOrigin: true, //是否跨域
        rewrite: (path) => path.replace('/api', '')
      }
    }
  }
})
