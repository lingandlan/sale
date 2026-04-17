<template>
  <div class="card-issue">
    <div class="page-header">
      <h1 class="page-title">门店卡发放</h1>
    </div>

    <div class="content-area">
      <!-- 步骤条 -->
      <div class="steps-card">
        <el-steps :active="currentStep" finish-status="success" align-center>
          <el-step title="查询会员" />
          <el-step title="选择充值中心" />
          <el-step title="选择卡号" />
          <el-step title="发放确认" />
        </el-steps>
      </div>

      <!-- Step 1: 查询会员 -->
      <div v-if="currentStep === 0" class="step-card">
        <div class="card-header">
          <h3 class="header-title">查询会员</h3>
        </div>
        <el-divider />
        <div class="search-row">
          <el-input
            v-model="memberPhone"
            placeholder="请输入会员手机号"
            maxlength="11"
            style="width: 300px"
            @keyup.enter="handleSearchMember"
          />
          <el-button type="primary" class="search-btn" :loading="searchingMember" @click="handleSearchMember">
            查询
          </el-button>
        </div>

        <div v-if="memberInfo" class="member-info-box">
          <div class="info-row">
            <span class="info-label">姓名</span>
            <span class="info-value">{{ memberInfo.name || '-' }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">手机号</span>
            <span class="info-value">{{ memberInfo.phone }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">等级</span>
            <span class="info-value">{{ memberInfo.level || '-' }}</span>
          </div>
        </div>

        <div class="step-actions">
          <el-button type="primary" class="next-btn" :disabled="!memberInfo" @click="nextStep">
            下一步
          </el-button>
        </div>
      </div>

      <!-- Step 2: 选择充值中心 -->
      <div v-if="currentStep === 1" class="step-card">
        <div class="card-header">
          <h3 class="header-title">选择充值中心</h3>
        </div>
        <el-divider />
        <el-select
          v-model="selectedCenterId"
          placeholder="请选择充值中心"
          style="width: 100%"
          :disabled="!userStore.canSelectAllCenters"
          @change="handleCenterChange"
        >
          <el-option
            v-for="center in centerOptions"
            :key="center.id"
            :label="center.name"
            :value="center.id"
          />
        </el-select>

        <div class="step-actions">
          <el-button @click="prevStep">上一步</el-button>
          <el-button type="primary" class="next-btn" :disabled="!selectedCenterId" @click="nextStep">
            下一步
          </el-button>
        </div>
      </div>

      <!-- Step 3: 选择卡号 -->
      <div v-if="currentStep === 2" class="step-card">
        <div class="card-header">
          <h3 class="header-title">选择卡号</h3>
        </div>
        <el-divider />
        <el-select
          v-model="selectedCardNo"
          placeholder="输入卡号搜索"
          filterable
          remote
          :remote-method="handleCardSearch"
          :loading="loadingCards"
          style="width: 100%"
        >
          <el-option
            v-for="cardNo in cardOptions"
            :key="cardNo"
            :label="cardNo"
            :value="cardNo"
          />
        </el-select>

        <div class="step-actions">
          <el-button @click="prevStep">上一步</el-button>
          <el-button type="primary" class="next-btn" :disabled="!selectedCardNo" @click="nextStep">
            下一步
          </el-button>
        </div>
      </div>

      <!-- Step 4: 发放确认 -->
      <div v-if="currentStep === 3" class="step-card">
        <div class="card-header">
          <h3 class="header-title">发放信息确认</h3>
        </div>
        <el-divider />

        <el-form :model="issueData" label-width="120px" class="issue-form">
          <el-form-item label="会员手机号">
            <span class="info-value">{{ memberPhone }}</span>
          </el-form-item>
          <el-form-item label="充值中心">
            <span class="info-value">{{ selectedCenterName }}</span>
          </el-form-item>
          <el-form-item label="卡号">
            <span class="info-value">{{ selectedCardNo }}</span>
          </el-form-item>
          <el-form-item label="发放原因" required>
            <el-select v-model="issueData.issueReason" placeholder="请选择发放原因" style="width: 100%">
              <el-option value="购买套餐包" label="购买套餐包" />
              <el-option value="推荐奖励" label="推荐奖励" />
              <el-option value="其他" label="其他" />
            </el-select>
          </el-form-item>
          <el-form-item label="发放方式" required>
            <el-select v-model="issueData.issueType" placeholder="请选择发放方式" style="width: 100%">
              <el-option :value="1" label="实体卡" />
              <el-option :value="2" label="虚拟卡" />
            </el-select>
          </el-form-item>
          <el-form-item v-if="issueData.issueReason === '推荐奖励'" label="关联购买人手机号">
            <el-input v-model="issueData.relatedUserPhone" placeholder="被推荐人手机号" maxlength="11" />
          </el-form-item>
          <el-form-item label="备注">
            <el-input
              v-model="issueData.remark"
              type="textarea"
              :rows="3"
              placeholder="请输入备注信息（可选）"
            />
          </el-form-item>
        </el-form>

        <div class="step-actions">
          <el-button @click="prevStep">上一步</el-button>
          <el-button type="primary" class="next-btn" :loading="submitting" @click="handleSubmit">
            确认发放
          </el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { extractErrorMessage } from '@/utils/request'
import { bindCard, getAvailableCards } from '@/api/card'
import { searchMember } from '@/api/recharge'
import { getCenterList } from '@/api/center'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()

// 步骤控制
const currentStep = ref(0)

// Step 1: 会员查询
const memberPhone = ref('')
const searchingMember = ref(false)
const memberInfo = ref<{ userId: string; name: string; phone: string; level: string } | null>(null)

// Step 2: 充值中心
const selectedCenterId = ref('')
const selectedCenterName = ref('')
const centerOptions = ref<{ id: string; name: string }[]>([])

// Step 3: 卡号选择
const selectedCardNo = ref('')
const cardOptions = ref<string[]>([])
const loadingCards = ref(false)

// Step 4: 发放信息
const submitting = ref(false)
const issueData = reactive({
  issueReason: '',
  issueType: 1 as number,
  relatedUserPhone: '',
  remark: ''
})

// ==================== Step 1: 查询会员 ====================

const handleSearchMember = async () => {
  const phone = memberPhone.value.trim()
  if (!phone) {
    ElMessage.warning('请输入手机号')
    return
  }
  if (!/^1[3-9]\d{9}$/.test(phone)) {
    ElMessage.warning('请输入正确的11位手机号')
    return
  }

  searchingMember.value = true
  try {
    const res = await searchMember(phone)
    const d = res.data
    memberInfo.value = {
      userId: d.userId,
      name: d.name,
      phone: d.phone,
      level: d.level
    }
    ElMessage.success('查询成功')
  } catch (err: any) {
    memberInfo.value = null
    ElMessage.error(extractErrorMessage(err, '未找到该会员'))
  } finally {
    searchingMember.value = false
  }
}

// ==================== Step 2: 充值中心 ====================

const loadCenterOptions = async () => {
  if (userStore.canSelectAllCenters) {
    try {
      const res = await getCenterList()
      centerOptions.value = (res.data || []).map((c: any) => ({ id: c.id, name: c.name }))
    } catch (err: any) {
      ElMessage.error(extractErrorMessage(err, '加载充值中心列表失败'))
      centerOptions.value = []
    }
  } else {
    // center_admin / operator: 只能选自己所属中心
    const cid = userStore.userCenterId
    const cname = userStore.userCenterName
    if (cid) {
      centerOptions.value = [{ id: String(cid), name: cname || '' }]
      selectedCenterId.value = String(cid)
      selectedCenterName.value = cname || ''
    }
  }
}

const handleCenterChange = (id: string) => {
  const center = centerOptions.value.find(c => c.id === id)
  selectedCenterName.value = center?.name || ''
  // 切换中心后清空已选卡号
  selectedCardNo.value = ''
  cardOptions.value = []
}

// ==================== Step 3: 卡号搜索 ====================

const handleCardSearch = async (query: string) => {
  if (!selectedCenterId.value) {
    ElMessage.warning('请先选择充值中心')
    return
  }
  loadingCards.value = true
  try {
    const res = await getAvailableCards(selectedCenterId.value, query || undefined)
    cardOptions.value = res.data?.cardNos || []
  } catch (err: any) {
    ElMessage.error(extractErrorMessage(err, '加载卡号列表失败'))
    cardOptions.value = []
  } finally {
    loadingCards.value = false
  }
}

// ==================== 步骤切换 ====================

const nextStep = () => {
  if (currentStep.value < 3) {
    currentStep.value++
  }
}

const prevStep = () => {
  if (currentStep.value > 0) {
    currentStep.value--
  }
}

// ==================== Step 4: 提交发放 ====================

const handleSubmit = async () => {
  if (!issueData.issueReason) {
    ElMessage.warning('请选择发放原因')
    return
  }
  if (!issueData.issueType) {
    ElMessage.warning('请选择发放方式')
    return
  }

  try {
    await ElMessageBox.confirm(
      `确认为手机号 ${memberPhone.value} 发放卡号 ${selectedCardNo.value}？`,
      '确认发放',
      {
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
  } catch {
    return
  }

  submitting.value = true
  try {
    await bindCard({
      cardNo: selectedCardNo.value,
      userPhone: memberPhone.value,
      issueReason: issueData.issueReason,
      issueType: issueData.issueType,
      rechargeCenterId: selectedCenterId.value,
      relatedUserPhone: issueData.relatedUserPhone || undefined,
      remark: issueData.remark || undefined
    })
    ElMessage.success('发放成功')
    resetAll()
  } catch (err: any) {
    ElMessage.error(extractErrorMessage(err, '发放失败'))
  } finally {
    submitting.value = false
  }
}

// ==================== 重置 ====================

const resetAll = () => {
  currentStep.value = 0
  memberPhone.value = ''
  memberInfo.value = null
  selectedCardNo.value = ''
  cardOptions.value = []
  issueData.issueReason = ''
  issueData.issueType = 1
  issueData.relatedUserPhone = ''
  issueData.remark = ''
}

// ==================== 初始化 ====================

onMounted(() => {
  if (!userStore.userInfo) {
    userStore.fetchUserInfo().then(() => loadCenterOptions())
  } else {
    loadCenterOptions()
  }
})
</script>

<style scoped>
.card-issue {
  background-color: #F5F5F5;
  min-height: calc(100vh - 64px);
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 64px;
  background-color: #FFFFFF;
  border-bottom: 1px solid #E5E5E5;
  padding: 16px 24px;
}

.page-title {
  font-family: 'Inter', sans-serif;
  font-size: 20px;
  font-weight: 600;
  color: #262626;
  margin: 0;
}

.content-area {
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.steps-card {
  background-color: #FFFFFF;
  border-radius: 8px;
  border: 1px solid #E5E5E5;
  padding: 24px 48px;
}

.step-card {
  background-color: #FFFFFF;
  border-radius: 8px;
  border: 1px solid #E5E5E5;
  padding: 24px;
}

.card-header {
  display: flex;
  gap: 8px;
  align-items: center;
}

.header-title {
  font-family: 'Inter', sans-serif;
  font-size: 16px;
  font-weight: 600;
  color: #262626;
  margin: 0;
}

.search-row {
  display: flex;
  gap: 12px;
  align-items: center;
}

.search-btn {
  background-color: #C00000;
  border-color: #C00000;
  border-radius: 4px;
}

.search-btn:hover {
  background-color: #A00000;
  border-color: #A00000;
}

.member-info-box {
  background-color: #F9F9F9;
  border-radius: 4px;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-top: 16px;
}

.info-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.info-label {
  font-family: 'Inter', sans-serif;
  font-size: 14px;
  color: #595959;
}

.info-value {
  font-family: 'Inter', sans-serif;
  font-size: 14px;
  font-weight: 600;
  color: #262626;
}

.step-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 24px;
}

.next-btn {
  background-color: #C00000;
  border-color: #C00000;
  border-radius: 4px;
  min-width: 100px;
}

.next-btn:hover {
  background-color: #A00000;
  border-color: #A00000;
}

.next-btn:disabled {
  background-color: #D9D9D9;
  border-color: #D9D9D9;
}

.issue-form {
  max-width: 600px;
}
</style>
