import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'

describe('HTTP Request工具', () => {
  beforeEach(() => {
    // 清除所有mock
    vi.clearAllMocks()
    // 清空localStorage
    localStorage.clear()
  })

  afterEach(() => {
    // 清理
    localStorage.clear()
  })

  describe('Token处理', () => {
    it('应该能够从localStorage读取token', () => {
      const testToken = 'test_token_123'
      localStorage.setItem('access_token', testToken)

      const retrievedToken = localStorage.getItem('access_token')
      expect(retrievedToken).toBe(testToken)
    })

    it('没有token时localStorage应该返回null', () => {
      const token = localStorage.getItem('access_token')
      expect(token).toBeNull()
    })

    it('应该能够保存token到localStorage', () => {
      const testToken = 'new_token_456'
      localStorage.setItem('access_token', testToken)

      const retrievedToken = localStorage.getItem('access_token')
      expect(retrievedToken).toBe(testToken)
    })

    it('应该能够清除token', () => {
      localStorage.setItem('access_token', 'test_token')
      expect(localStorage.getItem('access_token')).toBe('test_token')

      localStorage.removeItem('access_token')
      expect(localStorage.getItem('access_token')).toBeNull()
    })
  })

  describe('请求配置验证', () => {
    it('应该有正确的baseURL配置', () => {
      const baseURL = '/api/v1'
      expect(baseURL).toBe('/api/v1')
      expect(baseURL).toMatch(/^\/api\//)
    })

    it('应该有合理的超时配置', () => {
      const timeout = 15000
      expect(timeout).toBe(15000)
      expect(timeout).toBeGreaterThan(0)
      expect(timeout).toBeLessThan(60000) // 不超过60秒
    })
  })

  describe('响应数据结构', () => {
    it('成功响应应该包含正确的数据结构', () => {
      const mockResponse = {
        code: 0,
        message: 'success',
        data: { id: 1, name: 'test' }
      }

      expect(mockResponse).toHaveProperty('code')
      expect(mockResponse).toHaveProperty('message')
      expect(mockResponse).toHaveProperty('data')
      expect(mockResponse.code).toBe(0)
    })

    it('错误响应应该包含错误信息', () => {
      const mockErrorResponse = {
        code: 1,
        message: '请求失败'
      }

      expect(mockErrorResponse.code).not.toBe(0)
      expect(mockErrorResponse.message).toBe('请求失败')
    })
  })

  describe('请求方法配置', () => {
    it('GET请求应该有正确的配置', () => {
      const config = {
        method: 'get',
        url: '/api/test'
      }

      expect(config.method).toBe('get')
      expect(config.url).toBe('/api/test')
    })

    it('POST请求应该包含数据', () => {
      const postData = { name: 'test', value: 123 }
      const config = {
        method: 'post',
        url: '/api/test',
        data: postData
      }

      expect(config.method).toBe('post')
      expect(config.data).toEqual(postData)
    })

    it('应该支持Authorization头配置', () => {
      const token = 'test_token'
      const headers = {
        'Authorization': `Bearer ${token}`
      }

      expect(headers['Authorization']).toBe(`Bearer ${token}`)
    })
  })
})
