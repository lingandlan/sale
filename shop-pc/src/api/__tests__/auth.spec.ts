import { describe, it, expect, vi, beforeEach } from 'vitest'
import { login } from '../auth'
import request from '@/utils/request'

// Mock request模块
vi.mock('@/utils/request')

describe('登录API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('login函数', () => {
    it('应该发送登录请求到正确的接口', async () => {
      const mockResponse = {
        code: 0,
        message: 'success',
        data: {
          access_token: 'test_access_token',
          refresh_token: 'test_refresh_token'
        }
      }

      vi.mocked(request).mockResolvedValue(mockResponse)

      const loginData = {
        phone: '13800138000',
        password: 'Test123456'
      }

      await login(loginData)

      expect(request).toHaveBeenCalledWith({
        url: '/auth/login',
        method: 'POST',
        data: loginData
      })
    })

    it('应该返回登录响应数据', async () => {
      const mockResponse = {
        code: 0,
        message: 'success',
        data: {
          access_token: 'test_access_token',
          refresh_token: 'test_refresh_token',
          expires_in: 86400
        }
      }

      vi.mocked(request).mockResolvedValue(mockResponse)

      const result = await login({
        phone: '13800138000',
        password: 'Test123456'
      })

      expect(result).toEqual(mockResponse)
      expect(result.data.access_token).toBe('test_access_token')
      expect(result.data.refresh_token).toBe('test_refresh_token')
    })

    it('应该传递正确的请求方法', async () => {
      vi.mocked(request).mockResolvedValue({
        code: 0,
        message: 'success',
        data: {}
      })

      await login({
        phone: '13800138000',
        password: 'Test123456'
      })

      expect(request).toHaveBeenCalledWith(
        expect.objectContaining({
          method: 'POST'
        })
      )
    })

    it('应该传递完整的登录数据', async () => {
      vi.mocked(request).mockResolvedValue({
        code: 0,
        message: 'success',
        data: {}
      })

      const loginData = {
        phone: '13900139000',
        password: 'MyPassword123'
      }

      await login(loginData)

      expect(request).toHaveBeenCalledWith({
        url: '/auth/login',
        method: 'POST',
        data: loginData
      })
    })
  })

  describe('错误处理', () => {
    it('登录失败时应该抛出错误', async () => {
      const mockError = new Error('手机号或密码错误')
      vi.mocked(request).mockRejectedValue(mockError)

      await expect(
        login({
          phone: '13800138000',
          password: 'WrongPassword'
        })
      ).rejects.toThrow('手机号或密码错误')
    })

    it('网络错误时应该抛出错误', async () => {
      const mockError = new Error('网络连接失败')
      vi.mocked(request).mockRejectedValue(mockError)

      await expect(
        login({
          phone: '13800138000',
          password: 'Test123456'
        })
      ).rejects.toThrow('网络连接失败')
    })
  })
})
