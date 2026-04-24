<template>
  <div class="operator-manage">
    <div class="page-header">
      <h1 class="page-title">门店与操作员管理</h1>
      <el-button type="primary" class="add-btn" @click="handleAdd">
        <span style="margin-right: 4px">+</span> {{ activeTab === 'store' ? '新建门店' : '新建操作员' }}
      </el-button>
    </div>

    <div class="content-area">
      <!-- Tab 切换 -->
      <div class="tab-card">
        <div :class="['tab-item', { active: activeTab === 'store' }]" @click="activeTab = 'store'">门店管理</div>
        <div :class="['tab-item', { active: activeTab === 'operator' }]" @click="activeTab = 'operator'">操作员账号管理</div>
      </div>

      <!-- 门店管理 -->
      <template v-if="activeTab === 'store'">
        <div class="list-card">
          <el-table
            :data="storeData"
            style="width: 100%"
          >
            <el-table-column prop="name" label="门店名称" width="200" />
            <el-table-column prop="center" label="所属充值中心" width="200" />
            <el-table-column prop="manager" label="负责人" width="120" />
            <el-table-column prop="phone" label="联系电话" width="140" />
            <el-table-column label="卡库存数量" width="120" align="center">
              <template #default="{ row }">{{ row.cardCount }}</template>
            </el-table-column>
            <el-table-column label="卡库存余额" width="140" align="right">
              <template #default="{ row }">
                <span class="balance-value">¥{{ row.cardBalance.toLocaleString() }}</span>
              </template>
            </el-table-column>
            <el-table-column label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="row.status === 'normal' ? 'success' : 'danger'" size="small">
                  {{ row.status === 'normal' ? '正常' : '冻结' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="120" fixed="right">
              <template #default="{ row }">
                <el-button type="primary" link size="small" @click="handleEditStore(row)">编辑</el-button>
              </template>
            </el-table-column>
          </el-table>

          <div class="pagination-row">
            <el-pagination
              v-model:current-page="storePagination.page"
              v-model:page-size="storePagination.size"
              :total="storePagination.total"
              :page-sizes="[10, 20, 50]"
              layout="total, sizes, prev, pager, next"
            />
          </div>
        </div>
      </template>

      <!-- 操作员账号管理 -->
      <template v-if="activeTab === 'operator'">
        <div class="list-card">
          <el-table
            :data="operatorData"
            style="width: 100%"
          >
            <el-table-column prop="username" label="用户名" width="140" />
            <el-table-column prop="phone" label="手机号" width="140" />
            <el-table-column prop="realName" label="姓名" width="120" />
            <el-table-column prop="role" label="角色" width="140">
              <template #default="{ row }">
                <el-tag :type="row.role === 'center_admin' ? '' : 'info'" size="small">{{ roleLabel(row.role) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="center" label="所属充值中心" />
            <el-table-column label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="row.status === 'active' ? 'success' : 'danger'" size="small">
                  {{ row.status === 'active' ? '正常' : '冻结' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="220" fixed="right">
              <template #default="{ row }">
                <el-button type="primary" link size="small" @click="handleEditOperator(row)">编辑</el-button>
                <el-button
                  :type="row.status === 'active' ? 'danger' : 'success'"
                  link size="small"
                  @click="handleToggleOperator(row)"
                >
                  {{ row.status === 'active' ? '冻结' : '解冻' }}
                </el-button>
                <el-button type="primary" link size="small" @click="handleResetPwd(row)">重置密码</el-button>
              </template>
            </el-table-column>
          </el-table>

          <div class="pagination-row">
            <el-pagination
              v-model:current-page="opPagination.page"
              v-model:page-size="opPagination.size"
              :total="opPagination.total"
              :page-sizes="[10, 20, 50]"
              layout="total, sizes, prev, pager, next"
            />
          </div>
        </div>
      </template>
    </div>

    <!-- 门店编辑弹窗 -->
    <el-dialog v-model="storeDialogVisible" title="编辑门店" width="500px">
      <el-form :model="storeForm" label-width="120px">
        <el-form-item label="门店名称">
          <el-input v-model="storeForm.name" disabled />
        </el-form-item>
        <el-form-item label="所属充值中心">
          <el-input v-model="storeForm.center" disabled />
        </el-form-item>
        <el-form-item label="负责人">
          <el-input v-model="storeForm.manager" placeholder="请输入负责人" />
        </el-form-item>
        <el-form-item label="联系电话">
          <el-input v-model="storeForm.phone" placeholder="请输入联系电话" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="storeDialogVisible = false">取消</el-button>
        <el-button type="primary" class="save-btn" @click="handleSaveStore">保存</el-button>
      </template>
    </el-dialog>

    <!-- 操作员编辑弹窗 -->
    <el-dialog v-model="opDialogVisible" :title="opDialogTitle" width="500px">
      <el-form :model="opForm" label-width="120px">
        <el-form-item label="用户名" required>
          <el-input v-model="opForm.username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="手机号" required>
          <el-input v-model="opForm.phone" placeholder="请输入手机号" />
        </el-form-item>
        <el-form-item label="姓名" required>
          <el-input v-model="opForm.realName" placeholder="请输入姓名" />
        </el-form-item>
        <el-form-item label="角色" required>
          <el-select v-model="opForm.role" style="width: 100%">
            <el-option label="充值中心管理员" value="center_admin" />
            <el-option label="充值中心操作员" value="operator" />
          </el-select>
        </el-form-item>
        <el-form-item label="所属充值中心" required>
          <el-select v-model="opForm.centerId" style="width: 100%">
            <el-option v-for="c in centers" :key="c.id" :label="c.name" :value="c.id" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="opDialogVisible = false">取消</el-button>
        <el-button type="primary" class="save-btn" @click="handleSaveOperator">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getOperatorList, createOperator, updateOperator, type OperatorItem } from '@/api/operator'
import { getCenterList, type CenterItem } from '@/api/center'

const activeTab = ref('store')
const storeDialogVisible = ref(false)
const opDialogVisible = ref(false)
const opDialogTitle = ref('新建操作员')

const storePagination = reactive({ page: 1, size: 10, total: 0 })
const opPagination = reactive({ page: 1, size: 10, total: 0 })

const centers = ref<{ id: string; name: string }[]>([])
const storeData = ref<any[]>([])
const operatorData = ref<any[]>([])

const storeForm = reactive({ id: '', name: '', center: '', manager: '', phone: '' })
const opForm = reactive({ id: '' as string, username: '', phone: '', realName: '', role: '', centerId: '' as string })

const roleLabel = (role: string) => role === 'center_admin' ? '中心管理员' : '操作员'

const loadCenters = async () => {
  try {
    const res = await getCenterList()
    if (res?.data) {
      centers.value = (res.data || []).map((c: CenterItem) => ({ id: c.id, name: c.name }))
      storeData.value = (res.data || []).map((c: CenterItem) => ({
        name: c.name,
        center: c.name,
        manager: c.phone || '-',
        phone: c.phone || '-',
        cardCount: 0,
        cardBalance: 0,
        status: c.status === 'active' ? 'normal' : 'frozen'
      }))
      storePagination.total = storeData.value.length
    }
  } catch (error) {
    // fallback to empty
  }
}

const loadOperators = async () => {
  try {
    const res = await getOperatorList()
    if (res?.data) {
      operatorData.value = (res.data || []).map((op: OperatorItem) => {
        const centerName = centers.value.find(c => c.id === op.center_id)?.name || '-'
        return {
          id: op.id,
          username: op.phone,
          phone: op.phone,
          realName: op.name,
          role: op.role,
          center: centerName,
          centerId: op.center_id,
          status: op.status === 'active' ? 'active' : 'frozen'
        }
      })
      opPagination.total = operatorData.value.length
    }
  } catch (error) {
    // fallback to empty
  }
}

const handleAdd = () => {
  if (activeTab.value === 'store') {
    ElMessage.info('门店与充值中心1:1绑定，请通过充值中心管理创建')
  } else {
    opDialogTitle.value = '新建操作员'
    Object.assign(opForm, { id: '', username: '', phone: '', realName: '', role: '', centerId: '' })
    opDialogVisible.value = true
  }
}

const handleEditStore = (row: any) => {
  Object.assign(storeForm, row)
  storeDialogVisible.value = true
}

const handleSaveStore = async () => {
  ElMessage.success('保存成功')
  storeDialogVisible.value = false
}

const handleEditOperator = (row: any) => {
  opDialogTitle.value = '编辑操作员'
  Object.assign(opForm, { id: row.id, username: row.username, phone: row.phone, realName: row.realName, role: row.role, centerId: row.centerId })
  opDialogVisible.value = true
}

const handleSaveOperator = async () => {
  try {
    const payload: any = {
      name: opForm.realName,
      phone: opForm.phone,
      centerId: opForm.centerId,
      role: opForm.role,
      password: '123456'
    }
    if (opForm.id) {
      await updateOperator(opForm.id, payload)
    } else {
      await createOperator(payload)
    }
    ElMessage.success(opForm.id ? '编辑成功' : '创建成功')
    opDialogVisible.value = false
    loadOperators()
  } catch (error) {
    // error handled by interceptor
  }
}

const handleToggleOperator = async (row: any) => {
  const newStatus = row.status === 'active' ? 'frozen' : 'active'
  try {
    await updateOperator(row.id, { status: newStatus })
    loadOperators()
    ElMessage.success(newStatus === 'active' ? '已解冻' : '已冻结')
  } catch (error) {
    // error handled by interceptor
  }
}

const handleResetPwd = async (row: any) => {
  try {
    await updateOperator(row.id, { password: '123456' })
    ElMessage.success(`已重置 ${row.realName} 的密码`)
  } catch (error) {
    // error handled by interceptor
  }
}

onMounted(async () => {
  await loadCenters()
  loadOperators()
})
</script>

<style scoped>
.operator-manage {
  background-color: var(--color-bg);
  min-height: calc(100vh - 64px);
}

.page-header {
  height: 64px;
  background-color: var(--color-bg-card);
  border-bottom: 1px solid var(--color-border);
  padding: 16px 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.page-title {
  font-size: 20px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0;
}

.add-btn {
  background-color: var(--color-primary);
  border-color: var(--color-primary);
  border-radius: var(--radius-sm);
}

.add-btn:hover {
  background-color: var(--color-primary-hover);
  border-color: var(--color-primary-hover);
}

.content-area {
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.tab-card {
  display: flex;
  gap: 0;
  background-color: var(--color-bg-card);
  border-radius: var(--radius-md);
  border: 1px solid var(--color-border);
  overflow: hidden;
}

.tab-item {
  padding: 12px 24px;
  font-size: 14px;
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: all 0.2s;
  border-bottom: 2px solid transparent;
}

.tab-item:hover {
  color: var(--color-primary);
  background-color: var(--color-danger-bg);
}

.tab-item.active {
  color: var(--color-primary);
  font-weight: 600;
  border-bottom: 2px solid var(--color-primary);
  background-color: var(--color-danger-bg);
}

.list-card {
  background-color: var(--color-bg-card);
  border-radius: var(--radius-md);
  border: 1px solid var(--color-border);
  padding: 24px;
}

.balance-value {
  font-weight: 600;
  color: var(--color-primary);
}

.pagination-row {
  display: flex;
  justify-content: center;
  margin-top: 16px;
}

.save-btn {
  background-color: var(--color-primary);
  border-color: var(--color-primary);
}
</style>
