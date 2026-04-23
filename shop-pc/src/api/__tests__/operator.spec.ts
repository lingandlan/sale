import { describe, it, expect, vi, beforeEach } from 'vitest'
import { getOperatorList, createOperator, updateOperator, deleteOperator } from '../operator'
import request from '@/utils/request'

vi.mock('@/utils/request')

describe('Operator API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('getOperatorList应该发送GET请求到/operator', async () => {
    const mockResponse = { code: 0, data: [{ id: '1', name: '王五' }] }
    vi.mocked(request).get!.mockResolvedValue(mockResponse)

    const result = await getOperatorList()
    expect(request.get).toHaveBeenCalledWith('/operator')
    expect(result).toEqual(mockResponse)
  })

  it('createOperator应该发送POST请求到/operator', async () => {
    const mockResponse = { code: 0, data: { id: '1', name: '新操作员' } }
    vi.mocked(request).post!.mockResolvedValue(mockResponse)

    const data = { name: '新操作员', phone: '13800001111', center_id: '1', role: 'operator', password: '123456' }
    const result = await createOperator(data)
    expect(request.post).toHaveBeenCalledWith('/operator', data)
    expect(result).toEqual(mockResponse)
  })

  it('updateOperator应该发送PUT请求到/operator/:id', async () => {
    const mockResponse = { code: 0, data: { id: '1', name: '更新操作员' } }
    vi.mocked(request).put!.mockResolvedValue(mockResponse)

    const data = { name: '更新操作员' }
    const result = await updateOperator('1', data)
    expect(request.put).toHaveBeenCalledWith('/operator/1', data)
    expect(result).toEqual(mockResponse)
  })

  it('deleteOperator应该发送DELETE请求到/operator/:id', async () => {
    const mockResponse = { code: 0, data: { success: true } }
    vi.mocked(request).delete!.mockResolvedValue(mockResponse)

    const result = await deleteOperator('1')
    expect(request.delete).toHaveBeenCalledWith('/operator/1')
    expect(result).toEqual(mockResponse)
  })
})
