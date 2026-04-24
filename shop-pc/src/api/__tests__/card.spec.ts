import { describe, it, expect, vi, beforeEach } from 'vitest'
import { verifyCard, consumeCard, getCardList, issueCard, getCardDetail, toggleCardStatus, getCardStats } from '../card'
import request from '@/utils/request'

vi.mock('@/utils/request')

describe('Card API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('verifyCard应该发送GET请求到/card/verify/:cardNo', async () => {
    const mockResponse = { code: 0, data: { cardNo: 'TJ26001', holder: '张三', balance: 5000 } }
    vi.mocked(request).get!.mockResolvedValue(mockResponse)

    const result = await verifyCard('TJ26001')
    expect(request.get).toHaveBeenCalledWith('/card/verify/TJ26001')
    expect(result).toEqual(mockResponse)
  })

  it('consumeCard应该发送POST请求到/card/consume', async () => {
    const mockResponse = { code: 0, data: { transactionNo: 'TX001', balanceAfter: 3000 } }
    vi.mocked(request).post!.mockResolvedValue(mockResponse)

    const data = { cardNo: 'TJ26001', amount: 2000, remark: '消费' }
    const result = await consumeCard(data)
    expect(request.post).toHaveBeenCalledWith('/card/consume', data)
    expect(result).toEqual(mockResponse)
  })

  it('getCardList应该发送GET请求到/card/list并传递参数', async () => {
    const mockResponse = { code: 0, data: { list: [], total: 0 } }
    vi.mocked(request).get!.mockResolvedValue(mockResponse)

    const params = { status: 'active', page: 1, pageSize: 10 }
    await getCardList(params)
    expect(request.get).toHaveBeenCalledWith('/card/list', { params })
  })

  it('issueCard应该发送POST请求到/card/issue', async () => {
    const mockResponse = { code: 0, data: { cardNo: 'TJ26001', holder: '李四', balance: 5000 } }
    vi.mocked(request).post!.mockResolvedValue(mockResponse)

    const data = { holderName: '李四', holderPhone: '13800001111', amount: 5000 }
    const result = await issueCard(data)
    expect(request.post).toHaveBeenCalledWith('/card/issue', data)
    expect(result).toEqual(mockResponse)
  })

  it('getCardDetail应该发送GET请求到/card/detail/:cardNo', async () => {
    const mockResponse = { code: 0, data: { cardNo: 'TJ26001', balance: 5000 } }
    vi.mocked(request).get!.mockResolvedValue(mockResponse)

    const result = await getCardDetail('TJ26001')
    expect(request.get).toHaveBeenCalledWith('/card/detail/TJ26001')
    expect(result).toEqual(mockResponse)
  })

  it('toggleCardStatus应该发送POST请求到/card/:cardNo/status', async () => {
    const mockResponse = { code: 0, data: { success: true } }
    vi.mocked(request).post!.mockResolvedValue(mockResponse)

    const result = await toggleCardStatus('TJ26001', 'inactive')
    expect(request.post).toHaveBeenCalledWith('/card/TJ26001/status', { status: 'inactive' })
    expect(result).toEqual(mockResponse)
  })

  it('getCardStats应该发送GET请求到/card/stats', async () => {
    const mockResponse = { code: 0, data: { totalCards: 100, activeCards: 80 } }
    vi.mocked(request).get!.mockResolvedValue(mockResponse)

    const result = await getCardStats()
    expect(request.get).toHaveBeenCalledWith('/card/stats')
    expect(result).toEqual(mockResponse)
  })
})
