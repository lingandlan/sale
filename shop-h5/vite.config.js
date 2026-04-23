import { defineConfig } from 'vite'
import uni from '@dcloudio/vite-plugin-uni'

export default defineConfig({
  plugins: [uni.default ? uni.default() : uni()],
  resolve: {
    alias: {
      '@': '/src'
    }
  },
  css: {
    preprocessorOptions: {
      scss: {
        additionalData: `
          $u-primary: #C00000;
          $u-warning: #FAAD14;
          $u-success: #52C41A;
          $u-error: #FF4D4F;
          $u-info: #1677FF;
          $u-main-color: #262626;
          $u-content-color: #595959;
          $u-tips-color: #8C8C8C;
          $u-light-color: #BFBFBF;
          $u-border-color: #E5E5E5;
          $u-bg-color: #F5F5F5;
        `
      }
    }
  },
  server: {
    port: 5174,
    host: '0.0.0.0',
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true
      }
    }
  }
})
