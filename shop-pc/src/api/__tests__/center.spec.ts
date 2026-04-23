import { describe, it, expect, vi, beforeEach } from 'vitest'
import { getCenterList, createCenter, updateCenter, deleteCenter } from '../center'
import request from '@/utils/request'

vi.mock('@/utils/request')

describe('Center API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('getCenterList应该发送GET请求到/center', async () => {
    const mockResponse = { code: 0, data: [{ id: '1', name: '北京中心' }] }
    vi.mocked(request).get!.mockResolvedValue(mockResponse)

    const result = await getCenterList()
    expect(request.get).toHaveBeenCalledWith('/center')
    expect(result).toEqual(mockResponse)
  })

  it('createCenter应该发送POST请求到/center', async () => {
    const mockResponse = { code: 0, data: { id: '1', name: '新中心' } }
    vi.mocked(request).post!.mockResolvedValue(mockResponse)

    const data = { name: '新中心', code: 'BJ001', address: '北京市', phone: '010-12345678' }
    const result = await createCenter(data)
    expect(request.post).toHaveBeenCalledWith('/center', data)
    expect(result).toEqual(mockResponse)
  })

  it('updateCenter应该发送PUT请求到/center/:id', async () => {
    const mockResponse = { code: 0, data: { id: '1', name: '更新中心' } }
    vi.mocked(request).put!.mockResolvedValue(mockResponse)

    const data = { name: '更新中心' }
    const result = await updateCenter('1', data)
    expect(request.put).toHaveBeenCalledWith('/center/1', data)
    expect(result).toEqual(mockResponse)
  })

  it('deleteCenter应该发送DELETE请求到/center/:id', async () => {
    const mockResponse = { code: 0, data: { success: true } }
    vi.mocked(request).delete!.mockResolvedValue(mockResponse)

    const result = await deleteCenter('1')
    expect(request.delete).toHaveBeenCalledWith('/center/1')
    expect(result).toEqual(mockResponse)
  })
})
