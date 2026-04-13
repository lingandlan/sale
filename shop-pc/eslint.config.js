import parserTs from '@typescript-eslint/parser'
import parserVue from 'vue-eslint-parser'

export default [
  {
    ignores: ['dist/**', 'node_modules/**'],
  },

  // Vue 文件
  {
    files: ['**/*.vue'],
    languageOptions: {
      parser: parserVue,
      parserOptions: {
        parser: parserTs,
        sourceType: 'module',
      },
    },
  },

  // TS 文件
  {
    files: ['**/*.ts'],
    languageOptions: {
      parser: parserTs,
    },
  },

  // 自定义防护规则
  {
    files: ['**/*.{ts,vue}'],
    rules: {
      'no-restricted-syntax': [
        'error',
        {
          selector: 'CallExpression[callee.property.name="then"] > ArrowFunctionExpression[async=true]',
          message: '禁止在 .then() 中使用 async 回调。请改用 async/await + try/catch。',
        },
        {
          selector: 'CallExpression[callee.property.name="then"] > FunctionExpression[async=true]',
          message: '禁止在 .then() 中使用 async 回调。请改用 async/await + try/catch。',
        },
        {
          selector: 'CallExpression[callee.property.name="validate"] > ArrowFunctionExpression[async=true]',
          message: '禁止在 validate() 中使用 async 回调。请改用 try { await validate() } catch { return }。',
        },
        {
          selector: 'CallExpression[callee.property.name="validate"] > FunctionExpression[async=true]',
          message: '禁止在 validate() 中使用 async 回调。请改用 try { await validate() } catch { return }。',
        },
      ],
    },
  },
]
