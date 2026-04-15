import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')

  return {
    plugins: [vue()],
    resolve: {
      alias: {
        '@': resolve(__dirname, 'src')
      }
    },
    server: {
      port: parseInt(env.VITE_PORT || '5175'),
      // Mock 模式下不使用代理（直接请求 Apifox 本地服务）
      proxy: env.VITE_USE_MOCK === 'true' ? {} : {
        '/api': {
          target: env.VITE_API_BASE_URL ? env.VITE_API_BASE_URL.replace(/\/api\/v1$/, '') : 'http://localhost:8080',
          changeOrigin: true
        }
      }
    },
    test: {
      globals: true,
      environment: 'jsdom',
      setupFiles: ['./src/test/setup.ts'],
      coverage: {
        provider: 'v8',
        reporter: ['text', 'json', 'html'],
        exclude: [
          'node_modules/',
          'src/test/',
          '**/*.d.ts',
          '**/*.config.*',
          '**/mockData',
          'dist/',
        ]
      }
    }
  }
})
