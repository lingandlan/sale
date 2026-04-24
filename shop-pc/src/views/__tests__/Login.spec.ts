import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { ElMessage } from 'element-plus'
import Login from '../Login.vue'
import { login } from '@/api/auth'

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
vi.mock('@/api/auth', () => ({
  login: vi.fn()
}))

describe('Login.vue', () => {
  beforeEach(() => {
    // 清除所有mock调用记录
    vi.clearAllMocks()
    // 清空localStorage
    localStorage.clear()
  })

  describe('UI渲染', () => {
    it('应该渲染品牌区域', () => {
      const wrapper = mount(Login, {
        global: {
          stubs: {
            'el-input': true,
            'el-button': true,
            'el-checkbox': true
          }
        }
      })
      const brandSection = wrapper.find('.brand-section')
      expect(brandSection.exists()).toBe(true)
    })

    it('应该显示太积堂品牌名称', () => {
      const wrapper = mount(Login, {
        global: {
          stubs: {
            'el-input': true,
            'el-button': true,
            'el-checkbox': true
          }
        }
      })
      const brandName = wrapper.find('.brand-name')
      expect(brandName.text()).toBe('太积堂')
    })

    it('应该显示系统标题', () => {
      const wrapper = mount(Login, {
        global: {
          stubs: {
            'el-input': true,
            'el-button': true,
            'el-checkbox': true
          }
        }
      })
      const systemTitle = wrapper.find('.system-title')
      expect(systemTitle.text()).toBe('充值与门店管理系统')
    })

    it('应该显示登录表单卡片', () => {
      const wrapper = mount(Login, {
        global: {
          stubs: {
            'el-input': true,
            'el-button': true,
            'el-checkbox': true
          }
        }
      })
      const loginCard = wrapper.find('.login-card')
      expect(loginCard.exists()).toBe(true)
    })

    it('应该显示用户协议文本', () => {
      const wrapper = mount(Login, {
        global: {
          stubs: {
            'el-input': true,
            'el-button': true,
            'el-checkbox': true
          }
        }
      })
      const agreement = wrapper.find('.agreement-text')
      expect(agreement.text()).toContain('用户协议')
      expect(agreement.text()).toContain('隐私政策')
    })

    it('应该显示表单标题', () => {
      const wrapper = mount(Login, {
        global: {
          stubs: {
            'el-input': true,
            'el-button': true,
            'el-checkbox': true
          }
        }
      })
      const formHeader = wrapper.find('.form-header h2')
      expect(formHeader.exists()).toBe(true)
      expect(formHeader.text()).toBe('登录')
    })
  })

  describe('登录功能', () => {
    it('登录成功后应该保存token并跳转', async () => {
      const mockToken = {
        access_token: 'test_access_token',
        refresh_token: 'test_refresh_token'
      }

      vi.mocked(login).mockResolvedValue({
        data: mockToken
      } as any)

      const wrapper = mount(Login, {
        global: {
          stubs: {
            'el-input': true,
            'el-button': true,
            'el-checkbox': true
          }
        }
      })

      // 设置表单数据
      wrapper.vm.phone = '13800138000'
      wrapper.vm.password = 'Test123456'

      // 直接调用组件方法
      await wrapper.vm.handleLogin()

      expect(ElMessage.success).toHaveBeenCalledWith('登录成功')
      expect(mockPush).toHaveBeenCalledWith('/dashboard')
    })

    it('手机号为空时应该显示警告', async () => {
      const wrapper = mount(Login, {
        global: {
          stubs: {
            'el-input': true,
            'el-button': true,
            'el-checkbox': true
          }
        }
      })

      wrapper.vm.phone = ''
      wrapper.vm.password = 'Test123456'

      await wrapper.vm.handleLogin()

      expect(ElMessage.warning).toHaveBeenCalledWith('请输入手机号')
    })

    it('密码为空时应该显示警告', async () => {
      const wrapper = mount(Login, {
        global: {
          stubs: {
            'el-input': true,
            'el-button': true,
            'el-checkbox': true
          }
        }
      })

      wrapper.vm.phone = '13800138000'
      wrapper.vm.password = ''

      await wrapper.vm.handleLogin()

      expect(ElMessage.warning).toHaveBeenCalledWith('请输入密码')
    })
  })

  describe('忘记密码功能', () => {
    it('点击忘记密码应该显示提示信息', async () => {
      const wrapper = mount(Login, {
        global: {
          stubs: {
            'el-input': true,
            'el-button': true,
            'el-checkbox': true
          }
        }
      })

      await wrapper.vm.handleForgot()

      expect(ElMessage.info).toHaveBeenCalledWith('忘记密码功能开发中')
    })
  })
})
