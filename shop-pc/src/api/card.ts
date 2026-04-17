import request from '@/utils/request'

// 卡状态枚举
export const CardStatusMap: Record<number, string> = {
  1: '已入库',
  2: '已发放',
  3: '已激活',
  4: '已冻结',
  5: '已过期',
  6: '已作废'
}

export const CardStatusTagType: Record<number, string> = {
  1: 'info',
  2: '',
  3: 'success',
  4: 'warning',
  5: 'danger',
  6: 'info'
}

export const CardTypeMap: Record<number, string> = {
  1: '实体卡',
  2: '虚拟卡'
}

// 门店卡核销
export interface CardVerifyResponse {
  cardNo: string
  cardType: number
  balance: number
  status: number
  activatedAt: string | null
  expiredAt: string | null
  issuedAt: string | null
}

export interface CardConsumeData {
  cardNo: string
  amount: number
  remark?: string
}

// 门店卡列表
export interface CardListItem {
  id: string
  cardNo: string
  cardType: number
  balance: number
  status: number
  rechargeCenterId: string
  userId: string
  batchNo: string
  issueReason: string
  issuedAt: string | null
  activatedAt: string | null
  expiredAt: string | null
  createdAt: string
}

export interface CardListParams {
  status?: number
  cardNo?: string
  centerId?: string
  page: number
  pageSize: number
}

export interface CardListResponse {
  list: CardListItem[]
  total: number
}

// 门店卡详情
export interface CardDetail {
  id: string
  cardNo: string
  cardType: number
  balance: number
  status: number
  rechargeCenterId: string
  userId: string
  batchNo: string
  issueReason: string
  issuedAt: string | null
  activatedAt: string | null
  expiredAt: string | null
  createdAt: string
  transactions: {
    id: string
    type: string
    amount: number
    balanceBefore: number
    balanceAfter: number
    remark: string
    operatorId: string
    createdAt: string
  }[]
}

// 门店卡统计
export interface CardStatsResponse {
  totalCards: number
  inStockCards: number
  issuedCards: number
  activeCards: number
  frozenCards: number
  expiredCards: number
  voidedCards: number
  totalBalance: number
  todayConsume: number
  expireIn7Days: number
}

// 总卡库统计
export interface CardInventoryResponse {
  totalCards: number
  issuedCards: number
  inStockCards: number
}

// 批量入库（Excel上传）
export function batchImportCards(file: File) {
  const formData = new FormData()
  formData.append('file', file)
  return request.post('/card/batch-import', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}

// 划拨到充值中心（按数量）
export function allocateCards(data: { centerId: string; quantity: number }) {
  return request.post('/card/allocate', data)
}

// 绑定卡号和用户
export function bindCard(data: {
  cardNo: string
  userPhone: string
  issueReason: string
  issueType: number
  rechargeCenterId: string
  relatedUserPhone?: string
  remark?: string
}) {
  return request.post('/card/bind', data)
}

// 查询卡信息
export function verifyCard(cardNo: string) {
  return request.get(`/card/verify/${cardNo}`)
}

// 核销
export function consumeCard(data: CardConsumeData) {
  return request.post('/card/consume', data)
}

// 获取卡列表
export function getCardList(params: CardListParams) {
  return request.get('/card/list', { params })
}

// 获取卡详情
export function getCardDetail(cardNo: string) {
  return request.get(`/card/detail/${cardNo}`)
}

// 冻结卡
export function freezeCard(cardNo: string) {
  return request.post(`/card/${cardNo}/freeze`)
}

// 解冻卡
export function unfreezeCard(cardNo: string) {
  return request.post(`/card/${cardNo}/unfreeze`)
}

// 作废卡
export function voidCard(cardNo: string) {
  return request.post(`/card/${cardNo}/void`)
}

// 获取卡统计
export function getCardStats() {
  return request.get('/card/stats')
}

// 获取可发放卡号列表（已入库、未绑定的卡）
export function getAvailableCards(centerId: string, keyword?: string) {
  return request.get('/card/available', { params: { centerId, keyword } })
}

// 获取总卡库统计
export function getCardInventoryStats() {
  return request.get('/card/inventory-stats')
}
