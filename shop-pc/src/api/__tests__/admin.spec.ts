import { describe, it, expect, vi, beforeEach } from 'vitest'
import { getAdminUsers, createAdminUser, updateAdminUser, resetUserPassword, toggleUserStatus, deleteAdminUser } from '../admin'
import request from '@/utils/request'

vi.mock('@/utils/request')

describe('Admin API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('getAdminUsers', () => {
    it('应该发送GET请求到/admin/users并传递参数', async () => {
      const mockResponse = { code: 0, data: { items: [], total: 0 } }
      vi.mocked(request).get!.mockResolvedValue(mockResponse)

      const params = { keyword: 'test', role: 'operator', status: '1', page: 1, page_size: 10 }
      await getAdminUsers(params)

      expect(request.get).toHaveBeenCalledWith('/admin/users', { params })
    })

    it('应该返回用户列表数据', async () => {
      const mockResponse = { code: 0, data: { items: [{ id: 1, username: 'admin' }], total: 1, page: 1, page_size: 10 } }
      vi.mocked(request).get!.mockResolvedValue(mockResponse)

      const result = await getAdminUsers({ page: 1, page_size: 10 })
      expect(result).toEqual(mockResponse)
    })
  })

  describe('createAdminUser', () => {
    it('应该发送POST请求到/admin/users', async () => {
      const mockResponse = { code: 0, data: { id: 1, username: 'newuser' } }
      vi.mocked(request).post!.mockResolvedValue(mockResponse)

      const data = { username: 'newuser', phone: '13800138000', name: 'New User', role: 'operator', password: '123456' }
      const result = await createAdminUser(data)

      expect(request.post).toHaveBeenCalledWith('/admin/users', data)
      expect(result).toEqual(mockResponse)
    })
  })

  describe('updateAdminUser', () => {
    it('应该发送PUT请求到/admin/users/:id', async () => {
      const mockResponse = { code: 0, data: { id: 1 } }
      vi.mocked(request).put!.mockResolvedValue(mockResponse)

      const data = { name: 'Updated' }
      const result = await updateAdminUser(1, data)

      expect(request.put).toHaveBeenCalledWith('/admin/users/1', data)
      expect(result).toEqual(mockResponse)
    })
  })

  describe('resetUserPassword', () => {
    it('应该发送POST请求到/admin/users/:id/reset-password', async () => {
      const mockResponse = { code: 0, data: { success: true } }
      vi.mocked(request).post!.mockResolvedValue(mockResponse)

      const result = await resetUserPassword(1, 'newpass123')

      expect(request.post).toHaveBeenCalledWith('/admin/users/1/reset-password', { new_password: 'newpass123' })
      expect(result).toEqual(mockResponse)
    })
  })

  describe('toggleUserStatus', () => {
    it('应该发送PUT请求到/admin/users/:id/status', async () => {
      const mockResponse = { code: 0, data: { success: true } }
      vi.mocked(request).put!.mockResolvedValue(mockResponse)

      const result = await toggleUserStatus(1, 0)

      expect(request.put).toHaveBeenCalledWith('/admin/users/1/status', { status: 0 })
      expect(result).toEqual(mockResponse)
    })
  })

  describe('deleteAdminUser', () => {
    it('应该发送DELETE请求到/admin/users/:id', async () => {
      const mockResponse = { code: 0, data: { success: true } }
      vi.mocked(request).delete!.mockResolvedValue(mockResponse)

      const result = await deleteAdminUser(1)

      expect(request.delete).toHaveBeenCalledWith('/admin/users/1')
      expect(result).toEqual(mockResponse)
    })
  })
})
