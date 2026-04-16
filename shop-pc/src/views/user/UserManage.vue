<template>
  <div class="user-manage">
    <div class="page-header">
      <h1 class="page-title">用户管理</h1>
      <el-button type="primary" class="add-btn" @click="handleAdd">
        <span style="margin-right: 4px">+</span> 新建用户
      </el-button>
    </div>

    <div class="content-area">
      <!-- 筛选栏 -->
      <div class="filter-card">
        <el-input v-model="filters.keyword" placeholder="搜索用户名/手机号/姓名" style="width: 240px" clearable />
        <el-select v-model="filters.role" placeholder="角色" style="width: 160px" clearable>
          <el-option label="超级管理员" value="super_admin" />
          <el-option label="总部管理员" value="hq_admin" />
          <el-option label="财务运营" value="finance" />
          <el-option label="充值中心管理员" value="center_admin" />
          <el-option label="充值中心操作员" value="operator" />
        </el-select>
        <el-select v-model="filters.status" placeholder="状态" style="width: 120px" clearable>
          <el-option label="启用" value="1" />
          <el-option label="停用" value="0" />
        </el-select>
        <el-button type="primary" plain @click="handleSearch">查询</el-button>
        <el-button @click="handleResetFilter">重置</el-button>
      </div>

      <!-- 表格 -->
      <div class="list-card">
        <el-table
          :data="tableData"
          style="width: 100%"
          :header-cell-style="{ backgroundColor: '#FAFAFA', color: '#262626', fontWeight: '600' }"
        >
          <el-table-column prop="username" label="用户名" width="140" />
          <el-table-column prop="phone" label="手机号" width="140" />
          <el-table-column prop="realName" label="姓名" width="120" />
          <el-table-column prop="role" label="角色" width="150">
            <template #default="{ row }">
              <el-tag :type="roleTagType(row.role)" size="small">{{ roleLabel(row.role) }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="center" label="所属充值中心" />
          <el-table-column label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="row.status === 'active' ? 'success' : 'danger'" size="small">
                {{ row.status === 'active' ? '启用' : '停用' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="lastLogin" label="最后登录" width="180" />
          <el-table-column label="操作" width="180" fixed="right">
            <template #default="{ row }">
              <el-button type="primary" link size="small" @click="handleEdit(row)">编辑</el-button>
              <el-button
                :type="row.status === 'active' ? 'danger' : 'success'"
                link size="small"
                @click="handleToggleStatus(row)"
              >
                {{ row.status === 'active' ? '禁用' : '启用' }}
              </el-button>
              <el-button type="primary" link size="small" @click="handleResetPwd(row)">重置密码</el-button>
            </template>
          </el-table-column>
        </el-table>

        <div class="pagination-row">
          <el-pagination
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.size"
            :total="pagination.total"
            :page-sizes="[10, 20, 50]"
            layout="total, sizes, prev, pager, next"
          />
        </div>
      </div>
    </div>

    <!-- 编辑抽屉 -->
    <el-drawer v-model="drawerVisible" :title="drawerTitle" size="500px" :close-on-click-modal="false">
      <el-form :model="formData" label-width="120px" style="padding: 0 16px">
        <el-form-item label="用户名" required>
          <el-input v-model="formData.username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="手机号" required>
          <el-input v-model="formData.phone" placeholder="请输入手机号" maxlength="11" />
        </el-form-item>
        <el-form-item label="姓名" required>
          <el-input v-model="formData.realName" placeholder="请输入姓名" />
        </el-form-item>
        <el-form-item label="角色" required>
          <el-select v-model="formData.role" placeholder="请选择角色" style="width: 100%">
            <el-option label="超级管理员" value="super_admin" />
            <el-option label="总部管理员" value="hq_admin" />
            <el-option label="财务运营" value="finance" />
            <el-option label="充值中心管理员" value="center_admin" />
            <el-option label="充值中心操作员" value="operator" />
          </el-select>
        </el-form-item>
        <el-form-item label="所属充值中心" v-if="formData.role !== 'super_admin' && formData.role !== 'hq_admin'">
          <el-select v-model="formData.centerId" placeholder="请选择充值中心" style="width: 100%">
            <el-option v-for="c in centers" :key="c.id" :label="c.name" :value="c.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="初始密码" v-if="!formData.id">
          <el-input v-model="formData.password" type="password" show-password placeholder="默认密码 123456" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="drawerVisible = false">取消</el-button>
        <el-button type="primary" class="save-btn" @click="handleSaveUser">保存</el-button>
      </template>
    </el-drawer>

    <!-- 禁用确认弹窗 -->
    <el-dialog v-model="confirmVisible" :title="confirmTitle" width="400px">
      <p>{{ confirmMessage }}</p>
      <template #footer>
        <el-button @click="confirmVisible = false">取消</el-button>
        <el-button type="primary" :class="confirmAction === 'disable' ? 'danger-btn' : 'save-btn'" @click="handleConfirmAction">
          确认
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getAdminUsers, createAdminUser, updateAdminUser, toggleUserStatus, resetUserPassword } from '@/api/admin'
import { getCenterList } from '@/api/center'

const filters = reactive({ keyword: '', role: '', status: '' })
const pagination = reactive({ page: 1, size: 10, total: 0 })

const drawerVisible = ref(false)
const drawerTitle = ref('新建用户')
const confirmVisible = ref(false)
const confirmTitle = ref('')
const confirmMessage = ref('')
const confirmAction = ref<'disable' | 'enable'>('disable')
const editingUser = ref<any>(null)

const centers = ref<{ id: number; name: string }[]>([])

const formData = reactive({
  id: null as number | null,
  username: '',
  phone: '',
  realName: '',
  role: '',
  centerId: null as number | null,
  password: ''
})

const tableData = ref<any[]>([])

const roleLabel = (role: string) => {
  const map: Record<string, string> = { super_admin: '超级管理员', hq_admin: '总部管理员', finance: '财务运营', center_admin: '中心管理员', operator: '操作员' }
  return map[role] || role
}

const roleTagType = (role: string) => {
  const map: Record<string, string> = { super_admin: 'danger', hq_admin: 'warning', finance: '', center_admin: 'success', operator: 'info' }
  return map[role] || ''
}

const loadCenters = async () => {
  try {
    const res = await getCenterList()
    if (res?.data) {
      const items = Array.isArray(res.data) ? res.data : []
      centers.value = items.map((c: any) => ({ id: Number(c.id), name: c.name }))
    }
  } catch {
    // fallback to empty
  }
}

const loadData = async () => {
  try {
    const res = await getAdminUsers({
      keyword: filters.keyword || undefined,
      role: filters.role || undefined,
      status: filters.status ? Number(filters.status) : undefined,
      page: pagination.page,
      page_size: pagination.size
    })
    if (res?.data) {
      const items = res.data.items || []
      tableData.value = items.map((u: any) => ({
        id: u.id,
        username: u.username,
        phone: u.phone,
        realName: u.name,
        role: u.role,
        centerId: u.center_id,
        center: centers.value.find(c => c.id === Number(u.center_id))?.name || '-',
        status: u.status === 1 ? 'active' : 'disabled',
        lastLogin: u.last_login_at || '-'
      }))
      pagination.total = res.data.total || 0
    }
  } catch {
    // fallback to empty
  }
}

const handleSearch = () => { loadData() }
const handleResetFilter = () => { filters.keyword = ''; filters.role = ''; filters.status = '' }

const handleAdd = () => {
  drawerTitle.value = '新建用户'
  Object.assign(formData, { id: null, username: '', phone: '', realName: '', role: '', centerId: null, password: '' })
  drawerVisible.value = true
}

const handleEdit = (row: any) => {
  drawerTitle.value = '编辑用户'
  Object.assign(formData, {
    id: row.id,
    username: row.username,
    phone: row.phone,
    realName: row.realName,
    role: row.role,
    centerId: row.centerId,
    password: ''
  })
  drawerVisible.value = true
}

const handleSaveUser = async () => {
  // 前端校验
  if (!formData.username?.trim()) { ElMessage.warning('请输入用户名'); return }
  if (!formData.phone || formData.phone.length !== 11) { ElMessage.warning('请输入正确的11位手机号'); return }
  if (!formData.realName?.trim()) { ElMessage.warning('请输入姓名'); return }
  if (!formData.role) { ElMessage.warning('请选择角色'); return }

  try {
    if (formData.id) {
      await updateAdminUser(formData.id, {
        username: formData.username,
        name: formData.realName,
        phone: formData.phone,
        role: formData.role,
        center_id: formData.centerId ? Number(formData.centerId) : undefined
      })
    } else {
      await createAdminUser({
        username: formData.username,
        phone: formData.phone,
        name: formData.realName,
        role: formData.role,
        center_id: formData.centerId ? Number(formData.centerId) : undefined,
        password: formData.password || '123456'
      })
    }
    ElMessage.success(formData.id ? '编辑成功' : '创建成功')
    drawerVisible.value = false
    loadData()
  } catch {
    // error handled by interceptor
  }
}

const handleToggleStatus = (row: any) => {
  editingUser.value = row
  if (row.status === 'active') {
    confirmTitle.value = '禁用用户'
    confirmMessage.value = `确认禁用用户 "${row.realName}"？禁用后该用户将无法登录系统。`
    confirmAction.value = 'disable'
  } else {
    confirmTitle.value = '启用用户'
    confirmMessage.value = `确认启用用户 "${row.realName}"？`
    confirmAction.value = 'enable'
  }
  confirmVisible.value = true
}

const handleConfirmAction = async () => {
  if (editingUser.value) {
    try {
      await toggleUserStatus(editingUser.value.id, confirmAction.value === 'disable' ? 0 : 1)
      loadData()
      ElMessage.success(confirmAction.value === 'disable' ? '已禁用' : '已启用')
    } catch {
      // error handled by interceptor
    }
  }
  confirmVisible.value = false
}

const handleResetPwd = async (row: any) => {
  try {
    await resetUserPassword(row.id, '123456')
    ElMessage.success(`已重置 ${row.realName} 的密码为默认密码`)
  } catch {
    // error handled by interceptor
  }
}

onMounted(async () => {
  await loadCenters()
  loadData()
})
</script>

<style scoped>
.user-manage {
  background-color: #F5F5F5;
  min-height: calc(100vh - 64px);
}

.page-header {
  height: 64px;
  background-color: #FFFFFF;
  border-bottom: 1px solid #E5E5E5;
  padding: 16px 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.page-title {
  font-size: 20px;
  font-weight: 600;
  color: #262626;
  margin: 0;
}

.add-btn {
  background-color: #C00000;
  border-color: #C00000;
  border-radius: 4px;
}

.add-btn:hover {
  background-color: #A00000;
  border-color: #A00000;
}

.content-area {
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.filter-card {
  background-color: #FFFFFF;
  border-radius: 8px;
  border: 1px solid #E5E5E5;
  padding: 16px;
  display: flex;
  gap: 12px;
  align-items: center;
}

.list-card {
  background-color: #FFFFFF;
  border-radius: 8px;
  border: 1px solid #E5E5E5;
  padding: 24px;
}

.pagination-row {
  display: flex;
  justify-content: center;
  margin-top: 16px;
}

.save-btn {
  background-color: #C00000;
  border-color: #C00000;
}

.danger-btn {
  background-color: #FF4D4F;
  border-color: #FF4D4F;
}
</style>
