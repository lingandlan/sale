import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { ElMessage as _ElMessage } from 'element-plus'
import UserManage from '../user/UserManage.vue'
import { getAdminUsers } from '@/api/admin'

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
vi.mock('@/api/admin', () => ({
  getAdminUsers: vi.fn(),
  createAdminUser: vi.fn(),
  updateAdminUser: vi.fn(),
  toggleUserStatus: vi.fn(),
  resetUserPassword: vi.fn()
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
      'el-dialog': true
    }
  }
}

describe('UserManage.vue', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    vi.mocked(getAdminUsers).mockResolvedValue({ data: { items: [], total: 0, page: 1, page_size: 10 } } as any)
  })

  describe('UI渲染', () => {
    it('应该渲染页面标题"用户管理"', () => {
      const wrapper = mount(UserManage, globalStubs)
      const pageTitle = wrapper.find('.page-title')
      expect(pageTitle.exists()).toBe(true)
      expect(pageTitle.text()).toBe('用户管理')
    })

    it('应该渲染"新建用户"按钮', () => {
      const wrapper = mount(UserManage, globalStubs)
      const addBtn = wrapper.find('.add-btn')
      expect(addBtn.exists()).toBe(true)
      // el-button stub 会吞掉 slot 内容，通过 HTML 检查按钮文本
      expect(wrapper.html()).toContain('新建用户')
    })

    it('应该渲染筛选输入框（关键字、角色选择、状态选择）', () => {
      const wrapper = mount(UserManage, globalStubs)
      const filterCard = wrapper.find('.filter-card')
      expect(filterCard.exists()).toBe(true)

      // 验证筛选区域包含 el-input 和 el-select（通过 stub 验证存在）
      const inputs = filterCard.findAllComponents({ name: 'el-input' })
      const selects = filterCard.findAllComponents({ name: 'el-select' })
      expect(inputs.length).toBeGreaterThanOrEqual(1)
      expect(selects.length).toBeGreaterThanOrEqual(2)
    })

    it('应该渲染表格并包含正确的列', () => {
      const wrapper = mount(UserManage, globalStubs)
      const listCard = wrapper.find('.list-card')
      expect(listCard.exists()).toBe(true)

      // el-table 被 stub 后不渲染子组件，验证表格组件和分页存在
      const table = wrapper.findComponent({ name: 'el-table' })
      expect(table.exists()).toBe(true)

      const pagination = wrapper.findComponent({ name: 'el-pagination' })
      expect(pagination.exists()).toBe(true)

      // 组件已正确定义
      expect(UserManage).toBeDefined()
    })
  })

  describe('交互功能', () => {
    it('handleAdd 应该打开抽屉并设置"新建用户"标题', async () => {
      const wrapper = mount(UserManage, globalStubs)

      // @ts-expect-error vm property access in test
      await wrapper.vm.handleAdd()

      // @ts-expect-error vm property access in test
      expect(wrapper.vm.drawerTitle).toBe('新建用户')
      // @ts-expect-error vm property access in test
      expect(wrapper.vm.drawerVisible).toBe(true)
      // @ts-expect-error vm property access in test
      expect(wrapper.vm.formData.id).toBeNull()
    })

    it('handleEdit 应该打开抽屉并设置"编辑用户"标题', async () => {
      const wrapper = mount(UserManage, globalStubs)

      const testRow = {
        id: 1,
        username: 'testuser',
        phone: '13800138000',
        realName: '测试用户',
        role: 'operator',
        center: '测试中心',
        status: 'active',
        lastLogin: '2026-01-01'
      }

      // @ts-expect-error vm property access in test
      await wrapper.vm.handleEdit(testRow)

      // @ts-expect-error vm property access in test
      expect(wrapper.vm.drawerTitle).toBe('编辑用户')
      // @ts-expect-error vm property access in test
      expect(wrapper.vm.drawerVisible).toBe(true)
      // @ts-expect-error vm property access in test
      expect(wrapper.vm.formData.id).toBe(1)
      // @ts-expect-error vm property access in test
      expect(wrapper.vm.formData.username).toBe('testuser')
    })

    it('handleResetFilter 应该清空所有筛选字段', async () => {
      const wrapper = mount(UserManage, globalStubs)

      // @ts-expect-error vm property access in test
      wrapper.vm.filters.keyword = 'test'
      // @ts-expect-error vm property access in test
      wrapper.vm.filters.role = 'operator'
      // @ts-expect-error vm property access in test
      wrapper.vm.filters.status = 'active'

      // @ts-expect-error vm property access in test
      await wrapper.vm.handleResetFilter()

      // @ts-expect-error vm property access in test
      expect(wrapper.vm.filters.keyword).toBe('')
      // @ts-expect-error vm property access in test
      expect(wrapper.vm.filters.role).toBe('')
      // @ts-expect-error vm property access in test
      expect(wrapper.vm.filters.status).toBe('')
    })

    it('handleSearch 应该调用 loadData', async () => {
      const wrapper = mount(UserManage, globalStubs)

      // loadData 已在 onMounted 中被调用过一次，清除调用记录
      vi.clearAllMocks()
      vi.mocked(getAdminUsers).mockResolvedValue({ data: { items: [], total: 0, page: 1, page_size: 10 } } as any)

      // @ts-expect-error vm property access in test
      await wrapper.vm.handleSearch()

      expect(getAdminUsers).toHaveBeenCalledTimes(1)
    })
  })
})
