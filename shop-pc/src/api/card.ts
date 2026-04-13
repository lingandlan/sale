import request from '@/utils/request'

// 门店卡核销
export interface CardVerifyResponse {
  cardNo: string
  holder: string
  holderPhone: string
  balance: number
  status: 'active' | 'inactive' | 'expired'
  issueDate: string
  expiryDate: string
}

export interface CardConsumeData {
  cardNo: string
  amount: number
  remark?: string
}

export interface CardConsumeResponse {
  transactionNo: string
  balanceBefore: number
  balanceAfter: number
  consumeTime: string
}

// 门店卡管理
export interface CardListItem {
  id: string
  cardNo: string
  holder: string
  holderPhone: string
  balance: number
  status: 'active' | 'inactive' | 'expired'
  issueDate: string
  expiryDate: string
}

export interface CardListParams {
  status?: string
  holderPhone?: string
  page: number
  pageSize: number
}

export interface CardListResponse {
  list: CardListItem[]
  total: number
}

// 门店卡发放
export interface CardIssueData {
  holderName: string
  holderPhone: string
  amount: number
  remark?: string
}

export interface CardIssueResponse {
  cardNo: string
  holder: string
  balance: number
  issueDate: string
}

// 门店卡详情
export interface CardDetail {
  id: string
  cardNo: string
  holder: {
    name: string
    phone: string
  }
  balance: number
  status: 'active' | 'inactive' | 'expired'
  issueDate: string
  expiryDate: string
  issueCenter: string
  transactions: {
    id: string
    type: 'issue' | 'consume' | 'recharge'
    amount: number
    balanceAfter: number
    time: string
    remark?: string
  }[]
}

// 门店卡统计
export interface CardStatsResponse {
  totalCards: number
  activeCards: number
  totalBalance: number
  todayConsume: number
  todayIssue: number
  expireIn7Days: number
}

// 查询卡信息
export function verifyCard(cardNo: string) {
  return request.get<{ data: CardVerifyResponse }>(`/card/verify/${cardNo}`)
}

// 核销
export function consumeCard(data: CardConsumeData) {
  return request.post<{ data: CardConsumeResponse }>('/card/consume', data)
}

// 获取卡列表
export function getCardList(params: CardListParams) {
  return request.get<{ data: CardListResponse }>('/card/list', { params })
}

// 发放卡
export function issueCard(data: CardIssueData) {
  return request.post<{ data: CardIssueResponse }>('/card/issue', data)
}

// 获取卡详情
export function getCardDetail(cardNo: string) {
  return request.get<{ data: CardDetail }>(`/card/detail/${cardNo}`)
}

// 停用/启用卡
export function toggleCardStatus(cardNo: string, status: 'active' | 'inactive') {
  return request.post<{ data: { success: boolean } }>(`/card/${cardNo}/status`, { status })
}

// 获取卡统计
export function getCardStats() {
  return request.get<{ data: CardStatsResponse }>('/card/stats')
}
