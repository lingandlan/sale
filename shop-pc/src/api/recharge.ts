import request from '@/utils/request'

// B端充值申请
export interface BRechargeApplyData {
  centerId: string
  centerName: string
  amount: number
  lastMonthConsumption: number
  transactionNo: string
  screenshot: string
  remark: string
}

export interface BRechargeApplyResponse {
  id: string
  status: 'pending' | 'approved' | 'rejected'
  createdAt: string
}

// 充值审批列表
export interface BRechargeApprovalItem {
  id: string
  centerName: string
  amount: number
  points: number
  applicantName: string
  createdAt: string
  status: 'pending' | 'approved' | 'rejected'
  transactionNo?: string
  screenshot?: string
}

export interface BRechargeApprovalListResponse {
  list: BRechargeApprovalItem[]
  total: number
}

export interface BRechargeApprovalListParams {
  status?: string
  centerId?: string
  page: number
  pageSize: number
}

// 充值审批详情
export interface BRechargeApprovalDetail {
  id: string
  centerName: string
  centerId: string
  amount: number
  points: {
    base: number
    rebate: number
    rebateRate: number
    total: number
  }
  applicant: {
    id: string
    name: string
    phone: string
  }
  transactionNo?: string
  screenshot?: string
  status: 'pending' | 'approved' | 'rejected'
  remark?: string
  createdAt: string
  updatedAt?: string
  approvedBy?: string
  approvedAt?: string
}

// 充值审批操作
export interface ApprovalActionData {
  id: string
  action: 'approve' | 'reject'
  remark?: string
}

// C端充值录入
export interface CRechargeEntryData {
  memberId: string
  memberName: string
  memberPhone: string
  centerId: string
  centerName: string
  amount: number
  paymentMethod: 'cash' | 'wechat' | 'alipay' | 'card'
  remark: string
}

export interface CRechargeEntryResponse {
  id: string
  transactionNo: string
  balanceBefore: number
  balanceAfter: number
  createdAt: string
}

// 充值记录列表
export interface RechargeRecordItem {
  id: string
  transactionNo: string
  memberName: string
  memberPhone: string
  centerName: string
  amount: number
  operatorId: string
  operatorName: string
  createdAt: string
}

export interface RechargeRecordListParams {
  memberPhone?: string
  centerId?: string
  startDate?: string
  endDate?: string
  page: number
  pageSize: number
}

export interface RechargeRecordListResponse {
  list: RechargeRecordItem[]
  total: number
}

// 充值记录详情
export interface RechargeRecordDetail {
  id: string
  transactionNo: string
  memberId: string
  memberName: string
  memberPhone: string
  centerId: string
  centerName: string
  amount: number
  points: number
  operatorId: string
  operatorName: string
  remark?: string
  balanceBefore: number
  balanceAfter: number
  createdAt: string
}

// 上传文件（用 fetch 绕开 axios 对 FormData 的处理）
export async function uploadFile(file: File) {
  const formData = new FormData()
  formData.append('file', file)
  const token = localStorage.getItem('access_token')
  const res = await fetch('/api/v1/upload', {
    method: 'POST',
    headers: { Authorization: `Bearer ${token}` },
    body: formData
  })
  const json = await res.json()
  if (json.code !== 0) {
    throw new Error(json.message || '上传失败')
  }
  return json
}

// 提交B端充值申请
export function submitBRechargeApply(data: BRechargeApplyData) {
  return request.post('/recharge/b-apply', data)
}

// 获取充值审批列表
export function getBRechargeApprovalList(params: BRechargeApprovalListParams) {
  return request.get('/recharge/b-approval', { params })
}

// 获取充值审批详情
export function getBRechargeApprovalDetail(id: string) {
  return request.get(`/recharge/b-approval/${id}`)
}

// 审批操作
export function approvalAction(data: ApprovalActionData) {
  return request.post<{ data: { success: boolean } }>('/recharge/b-approval/action', data)
}

// C端充值录入
export function submitCRechargeEntry(data: CRechargeEntryData) {
  return request.post('/recharge/c-entry', data)
}

// 会员查询
export interface MemberInfo {
  userId: string
  name: string
  phone: string
  balance: number
  level: string
  nickName: string
}

export function searchMember(phone: string) {
  return request.get('/recharge/c-entry/search-member', { params: { phone } })
}

// 获取充值中心列表
export interface CenterItem {
  id: string
  name: string
  balance: number
  status: string
}

export function getCenterList() {
  return request.get('/center')
}

// 获取充值中心详情（含余额）
export function getCenterDetail(id: string) {
  return request.get<{ data: { id: string; name: string; balance: number } }>(`/center/${id}`)
}

// 获取充值记录列表
export function getRechargeRecordList(params: RechargeRecordListParams) {
  return request.get('/recharge/records', { params })
}

// 获取充值记录详情
export function getRechargeRecordDetail(id: string) {
  return request.get(`/recharge/records/${id}`)
}

// 获取充值中心上月消费
export interface LastMonthConsumption {
  consumption: number
  rebateRate: number
  month: string
}

export function getCenterLastMonthConsumption(centerId: string) {
  return request.get(`/center/${centerId}/last-month-consumption`)
}

// 录入月度消费
export interface MonthlyConsumptionData {
  centerId: string
  month: string
  consumption: number
}

export function upsertMonthlyConsumption(data: MonthlyConsumptionData) {
  return request.post<{ data: { success: boolean } }>('/center-monthly-consumption', data)
}

// 查询月度消费列表
export interface MonthlyConsumptionRecord {
  id: string
  centerId: string
  consumption: number
  month: string
  createdAt: string
}

export function listMonthlyConsumption(month?: string) {
  return request.get('/center-monthly-consumption', { params: { month } })
}
