import { describe, it, expect, vi, beforeEach } from 'vitest'
import { submitBRechargeApply, getBRechargeApprovalList, getBRechargeApprovalDetail, approvalAction, submitCRechargeEntry, getRechargeRecordList, getRechargeRecordDetail } from '../recharge'
import request from '@/utils/request'

vi.mock('@/utils/request')

describe('Recharge API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('submitBRechargeApply应该发送POST请求到/recharge/b-apply', async () => {
    const mockResponse = { code: 0, data: { id: '1', status: 'pending' } }
    vi.mocked(request).post!.mockResolvedValue(mockResponse)

    const data = { centerId: '1', amount: 10000 }
    const result = await submitBRechargeApply(data)
    expect(request.post).toHaveBeenCalledWith('/recharge/b-apply', data)
    expect(result).toEqual(mockResponse)
  })

  it('getBRechargeApprovalList应该发送GET请求到/recharge/b-approval', async () => {
    const mockResponse = { code: 0, data: { list: [], total: 0 } }
    vi.mocked(request).get!.mockResolvedValue(mockResponse)

    const params = { status: 'pending', page: 1, pageSize: 10 }
    await getBRechargeApprovalList(params)
    expect(request.get).toHaveBeenCalledWith('/recharge/b-approval', { params })
  })

  it('getBRechargeApprovalDetail应该发送GET请求到/recharge/b-approval/:id', async () => {
    const mockResponse = { code: 0, data: { id: '1', status: 'pending' } }
    vi.mocked(request).get!.mockResolvedValue(mockResponse)

    const result = await getBRechargeApprovalDetail('1')
    expect(request.get).toHaveBeenCalledWith('/recharge/b-approval/1')
    expect(result).toEqual(mockResponse)
  })

  it('approvalAction应该发送POST请求到/recharge/b-approval/action', async () => {
    const mockResponse = { code: 0, data: { success: true } }
    vi.mocked(request).post!.mockResolvedValue(mockResponse)

    const data = { id: '1', action: 'approve' as const }
    const result = await approvalAction(data)
    expect(request.post).toHaveBeenCalledWith('/recharge/b-approval/action', data)
    expect(result).toEqual(mockResponse)
  })

  it('submitCRechargeEntry应该发送POST请求到/recharge/c-entry', async () => {
    const mockResponse = { code: 0, data: { id: '1', transactionNo: 'TX001' } }
    vi.mocked(request).post!.mockResolvedValue(mockResponse)

    const data = { memberId: '1', memberName: '张三', memberPhone: '13800001111', amount: 5000, paymentMethod: 'cash' as const }
    const result = await submitCRechargeEntry(data)
    expect(request.post).toHaveBeenCalledWith('/recharge/c-entry', data)
    expect(result).toEqual(mockResponse)
  })

  it('getRechargeRecordList应该发送GET请求到/recharge/records', async () => {
    const mockResponse = { code: 0, data: { list: [], total: 0 } }
    vi.mocked(request).get!.mockResolvedValue(mockResponse)

    const params = { memberPhone: '13800001111', page: 1, pageSize: 10 }
    await getRechargeRecordList(params)
    expect(request.get).toHaveBeenCalledWith('/recharge/records', { params })
  })

  it('getRechargeRecordDetail应该发送GET请求到/recharge/records/:id', async () => {
    const mockResponse = { code: 0, data: { id: '1', amount: 5000 } }
    vi.mocked(request).get!.mockResolvedValue(mockResponse)

    const result = await getRechargeRecordDetail('1')
    expect(request.get).toHaveBeenCalledWith('/recharge/records/1')
    expect(result).toEqual(mockResponse)
  })

  it('API调用失败时应该抛出错误', async () => {
    vi.mocked(request).post!.mockRejectedValue(new Error('服务器错误'))

    await expect(submitBRechargeApply({ centerId: '1', amount: 100 })).rejects.toThrow('服务器错误')
  })
})
