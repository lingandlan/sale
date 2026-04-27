import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { ElMessage } from 'element-plus'
import CardManage from '../card/CardManage.vue'
import { getCardList, getCardStats } from '@/api/card'

// Mock vue-router
const mockPush = vi.fn()
vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: mockPush
  })
}))

// Mock element-plus message
vi.mock('element-plus', () => ({
  ElMessage: {
    warning: vi.fn(),
    success: vi.fn(),
    error: vi.fn(),
    info: vi.fn()
  }
}))

// Mock API
vi.mock('@/api/card', () => ({
  getCardList: vi.fn(),
  getCardStats: vi.fn(),
  verifyCard: vi.fn(),
  consumeCard: vi.fn(),
  getCardDetail: vi.fn()
}))

const globalStubs = {
  global: {
    stubs: {
      'el-input': true,
      'el-select': true,
      'el-option': true,
      'el-button': true,
      'el-table': true,
      'el-table-column': true,
      'el-tag': true,
      'el-pagination': true,
      'el-drawer': true,
      'el-form': true,
      'el-form-item': true,
      'el-dialog': true,
      'el-divider': true
    }
  }
}

describe('CardManage.vue', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    vi.mocked(getCardList).mockResolvedValue({
      data: { list: [], total: 0 }
    } as any)
    vi.mocked(getCardStats).mockResolvedValue({
      data: { totalCards: 100, activeCards: 60, totalBalance: 50000, todayConsume: 0, todayIssue: 0, expireIn7Days: 0 }
    } as any)
  })

  describe('UI渲染', () => {
    it('应该渲染页面标题"门店卡管理"', () => {
      const wrapper = mount(CardManage, globalStubs)
      const pageTitle = wrapper.find('.page-title')
      expect(pageTitle.exists()).toBe(true)
      expect(pageTitle.text()).toBe('门店卡管理')
    })

    it('应该渲染统计卡片区域', () => {
      const wrapper = mount(CardManage, globalStubs)
      const statsRow = wrapper.find('.stats-row')
      expect(statsRow.exists()).toBe(true)

      const statCards = wrapper.findAll('.stat-card')
      expect(statCards.length).toBe(4)
    })

    it('应该渲染筛选区域（状态选择和查询按钮）', () => {
      const wrapper = mount(CardManage, globalStubs)
      const filterCard = wrapper.find('.filter-card')
      expect(filterCard.exists()).toBe(true)

      const selects = filterCard.findAllComponents({ name: 'el-select' })
      expect(selects.length).toBeGreaterThanOrEqual(1)

      // 验证查询按钮存在
      const buttons = filterCard.findAllComponents({ name: 'el-button' })
      expect(buttons.length).toBeGreaterThanOrEqual(1)
    })

    it('应该渲染门店卡列表表格', () => {
      const wrapper = mount(CardManage, globalStubs)
      const listCard = wrapper.find('.list-card')
      expect(listCard.exists()).toBe(true)

      // el-table 被 stub 后不渲染子组件，验证表格组件存在
      const table = wrapper.findComponent({ name: 'el-table' })
      expect(table.exists()).toBe(true)

      // 验证列表区域标题存在
      const listHeader = wrapper.find('.list-header')
      expect(listHeader.exists()).toBe(true)
      expect(listHeader.text()).toContain('门店卡列表')
    })
  })

  describe('交互功能', () => {
    it('handleRefresh 应该调用 loadData 并显示成功消息', async () => {
      const wrapper = mount(CardManage, globalStubs)
      const vm = wrapper.vm as any

      // 清除 onMounted 中的调用记录
      vi.clearAllMocks()
      vi.mocked(getCardList).mockResolvedValue({
        data: { list: [], total: 0 }
      } as any)
      vi.mocked(getCardStats).mockResolvedValue({
        data: { totalCards: 100, activeCards: 60, totalBalance: 50000, todayConsume: 0, todayIssue: 0, expireIn7Days: 0 }
      } as any)

      await vm.handleRefresh()

      expect(getCardList).toHaveBeenCalledTimes(1)
      expect(getCardStats).toHaveBeenCalledTimes(1)
      expect(ElMessage.success).toHaveBeenCalledWith('刷新成功')
    })

    it('handleSearch 应该调用 loadData', async () => {
      const wrapper = mount(CardManage, globalStubs)
      const vm = wrapper.vm as any

      // 清除 onMounted 中的调用记录
      vi.clearAllMocks()
      vi.mocked(getCardList).mockResolvedValue({
        data: { list: [], total: 0 }
      } as any)
      vi.mocked(getCardStats).mockResolvedValue({
        data: { totalCards: 100, activeCards: 60, totalBalance: 50000, todayConsume: 0, todayIssue: 0, expireIn7Days: 0 }
      } as any)

      await vm.handleSearch()

      expect(getCardList).toHaveBeenCalledTimes(1)
    })
  })

  describe('辅助函数', () => {
    it('getStatusText 应该返回正确的状态标签', () => {
      const wrapper = mount(CardManage, globalStubs)
      const vm = wrapper.vm as any

      expect(vm.getStatusText('active')).toBe('已发放')
      expect(vm.getStatusText('inactive')).toBe('已冻结')
      expect(vm.getStatusText('expired')).toBe('已过期')
      expect(vm.getStatusText('unknown')).toBe('未知')
    })

    it('getStatusType 应该返回正确的标签类型', () => {
      const wrapper = mount(CardManage, globalStubs)
      const vm = wrapper.vm as any

      expect(vm.getStatusType('active')).toBe('success')
      expect(vm.getStatusType('inactive')).toBe('warning')
      expect(vm.getStatusType('expired')).toBe('danger')
      expect(vm.getStatusType('unknown')).toBe('info')
    })
  })
})
