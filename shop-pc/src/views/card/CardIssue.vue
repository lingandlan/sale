<template>
  <div class="card-issue">
    <div class="page-header">
      <h1 class="page-title">门店卡发放</h1>
    </div>

    <div class="content-area">
      <!-- 左侧表单区 -->
      <div class="form-card">
        <!-- 查询会员 -->
        <div class="form-section">
          <div class="section-title">查询会员</div>
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
            <div class="member-row">
              <div class="member-item">
                <span class="info-label">手机号</span>
                <span class="info-value">{{ memberInfo.phone }}</span>
              </div>
              <div class="member-item">
                <span class="info-label">等级</span>
                <span class="info-value">{{ memberInfo.level || '-' }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- 选择充值中心 -->
        <div class="form-section">
          <div class="section-title">选择充值中心</div>
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
        </div>

        <!-- 选择卡号 -->
        <div class="form-section">
          <div class="section-title">选择卡号</div>
          <div v-if="centerCardCount !== null" class="card-count-hint">
            该中心剩余 <span class="count-num">{{ centerCardCount }}</span> 张可用卡
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
            @change="handleCardChange"
          >
            <el-option
              v-for="cardNo in cardOptions"
              :key="cardNo"
              :label="cardNo"
              :value="cardNo"
            />
          </el-select>
        </div>

        <!-- 发放信息 -->
        <div class="form-section">
          <div class="section-title">发放信息</div>
          <el-divider />
          <el-form :model="issueData" label-width="100px" class="issue-form">
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
            <el-form-item label="备注">
              <el-input
                v-model="issueData.remark"
                type="textarea"
                :rows="3"
                placeholder="请输入备注信息（可选）"
              />
            </el-form-item>
          </el-form>
        </div>
      </div>

      <!-- 右侧汇总卡片 -->
      <div class="summary-card">
        <div class="summary-title">发放汇总</div>
        <el-divider />

        <div class="summary-list">
          <div class="summary-row">
            <span class="summary-label">会员手机号</span>
            <span class="summary-value">{{ memberPhone || '-' }}</span>
          </div>
          <div class="summary-row">
            <span class="summary-label">充值中心</span>
            <span class="summary-value">{{ selectedCenterName || '-' }}</span>
          </div>
          <div class="summary-row">
            <span class="summary-label">卡号</span>
            <span class="summary-value summary-value--highlight">{{ selectedCardNo || '-' }}</span>
          </div>
          <div class="summary-row">
            <span class="summary-label">面值</span>
            <span class="summary-value summary-value--highlight">{{ cardDenomination || '-' }}</span>
          </div>
          <div class="summary-row">
            <span class="summary-label">发放原因</span>
            <span class="summary-value">{{ issueData.issueReason || '-' }}</span>
          </div>
          <div class="summary-row">
            <span class="summary-label">发放方式</span>
            <span class="summary-value">{{ issueTypeLabel }}</span>
          </div>
          <div class="summary-row">
            <span class="summary-label">备注</span>
            <span class="summary-value">{{ issueData.remark || '-' }}</span>
          </div>
        </div>

        <div class="summary-status">
          <el-tag v-if="isFormComplete" type="success" effect="light">信息已完整，可以发放</el-tag>
          <el-tag v-else type="warning" effect="light">信息不完整，请填写所有必填项</el-tag>
        </div>
      </div>
    </div>

    <!-- 底部操作栏 -->
    <div class="action-bar">
      <el-button @click="handleReset">重置</el-button>
      <el-button type="primary" class="submit-btn" :loading="submitting" :disabled="!isFormComplete" @click="handleSubmit">
        确认发放
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { extractErrorMessage } from '@/utils/request'
import { bindCard, getAvailableCards, getCardDetail, getAvailableCardCount } from '@/api/card'
import { searchMember } from '@/api/recharge'
import { getCenterList } from '@/api/center'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()

// 会员查询
const memberPhone = ref('')
const searchingMember = ref(false)
const memberInfo = ref<{ userId: string; name: string; phone: string; level: string } | null>(null)

// 充值中心
const selectedCenterId = ref('')
const selectedCenterName = ref('')
const centerOptions = ref<{ id: string; name: string }[]>([])

// 卡号
const selectedCardNo = ref('')
const cardDenomination = ref('')
const cardOptions = ref<string[]>([])
const loadingCards = ref(false)
const centerCardCount = ref<number | null>(null)

// 发放信息
const submitting = ref(false)
const issueData = reactive({
  issueReason: '',
  issueType: 1 as number,
  remark: ''
})

// 汇总是否完整
const isFormComplete = computed(() => {
  return !!(
    memberPhone.value &&
    selectedCenterId.value &&
    selectedCardNo.value &&
    issueData.issueReason &&
    issueData.issueType
  )
})

// 发放方式标签
const issueTypeLabel = computed(() => {
  if (!issueData.issueType) return '-'
  return issueData.issueType === 1 ? '实体卡' : '虚拟卡'
})

// ==================== 查询会员 ====================

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

// ==================== 充值中心 ====================

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
    const cid = userStore.userCenterId
    const cname = userStore.userCenterName
    if (cid) {
      centerOptions.value = [{ id: String(cid), name: cname || '' }]
      selectedCenterId.value = String(cid)
      selectedCenterName.value = cname || ''
    }
  }
}

const handleCenterChange = async (id: string) => {
  const center = centerOptions.value.find(c => c.id === id)
  selectedCenterName.value = center?.name || ''
  selectedCardNo.value = ''
  cardDenomination.value = ''
  cardOptions.value = []

  if (!id) {
    centerCardCount.value = null
    return
  }

  // 加载可用卡数量
  try {
    const res = await getAvailableCardCount(id)
    centerCardCount.value = res.data?.count ?? 0
  } catch {
    centerCardCount.value = null
  }

  // 触发卡号搜索
  handleCardSearch('')
}

// ==================== 卡号 ====================

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

const handleCardChange = async (cardNo: string) => {
  if (!cardNo) {
    cardDenomination.value = ''
    return
  }
  try {
    const res = await getCardDetail(cardNo)
    const balance = res.data?.card?.balance
    cardDenomination.value = balance != null && balance !== '' ? `¥${Number(balance).toFixed(2)}` : ''
  } catch {
    cardDenomination.value = ''
  }
}

// ==================== 重置 ====================

const handleReset = () => {
  memberPhone.value = ''
  memberInfo.value = null
  selectedCenterId.value = ''
  selectedCenterName.value = ''
  selectedCardNo.value = ''
  cardDenomination.value = ''
  cardOptions.value = []
  centerCardCount.value = null
  issueData.issueReason = ''
  issueData.issueType = 1
  issueData.remark = ''
}

// ==================== 提交 ====================

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
      remark: issueData.remark || undefined
    })
    ElMessage.success('发放成功')
    memberPhone.value = ''
    memberInfo.value = null
    selectedCardNo.value = ''
    cardDenomination.value = ''
    cardOptions.value = []
    issueData.issueReason = ''
    issueData.issueType = 1
    issueData.remark = ''
    // 刷新当前中心的剩余卡数量
    if (selectedCenterId.value) {
      handleCenterChange(selectedCenterId.value)
    }
  } catch (err: any) {
    ElMessage.error(extractErrorMessage(err, '发放失败'))
  } finally {
    submitting.value = false
  }
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
  background-color: var(--color-bg);
  min-height: calc(100vh - 64px);
  display: flex;
  flex-direction: column;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 64px;
  background-color: var(--color-bg-card);
  border-bottom: 1px solid var(--color-border);
  padding: 16px 24px;
}

.page-title {
  font-family: var(--font-family);
  font-size: 20px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0;
}

.content-area {
  padding: 24px;
  display: flex;
  gap: 24px;
  flex: 1;
}

.form-card {
  flex: 1;
  background-color: var(--color-bg-card);
  border-radius: var(--radius-md);
  border: 1px solid var(--color-border);
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 0;
}

.form-section {
  margin-bottom: 24px;
}

.section-title {
  font-family: var(--font-family);
  font-size: 16px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.search-row {
  display: flex;
  gap: 12px;
  align-items: center;
  margin-top: 12px;
}

.search-btn {
  background-color: var(--color-primary);
  border-color: var(--color-primary);
  border-radius: var(--radius-sm);
}

.search-btn:hover {
  background-color: var(--color-primary-hover);
  border-color: var(--color-primary-hover);
}

.member-info-box {
  background-color: var(--color-bg-section);
  border-radius: var(--radius-sm);
  padding: 16px;
  margin-top: 16px;
}

.member-row {
  display: flex;
  gap: 32px;
}

.member-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.info-label {
  font-family: var(--font-family);
  font-size: 12px;
  color: var(--color-text-muted);
}

.info-value {
  font-family: var(--font-family);
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.card-count-hint {
  font-family: var(--font-family);
  font-size: 13px;
  color: var(--color-text-muted);
  margin-top: 8px;
}

.count-num {
  color: var(--color-primary);
  font-weight: 600;
}

.issue-form {
  margin-top: 12px;
}

.action-bar {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 24px;
  background-color: var(--color-bg-card);
  border-top: 1px solid var(--color-border);
}

.submit-btn {
  background-color: var(--color-primary);
  border-color: var(--color-primary);
  border-radius: var(--radius-sm);
  min-width: 100px;
}

.submit-btn:hover {
  background-color: var(--color-primary-hover);
  border-color: var(--color-primary-hover);
}

.submit-btn:disabled {
  background-color: var(--color-border);
  border-color: var(--color-border);
}

/* 汇总卡片 */
.summary-card {
  width: 360px;
  background-color: var(--color-bg-card);
  border-radius: var(--radius-md);
  border: 1px solid var(--color-border);
  padding: 24px;
  flex-shrink: 0;
}

.summary-title {
  font-family: var(--font-family);
  font-size: 16px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.summary-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-top: 16px;
}

.summary-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.summary-label {
  font-family: var(--font-family);
  font-size: 14px;
  color: var(--color-text-muted);
}

.summary-value {
  font-family: var(--font-family);
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.summary-value--highlight {
  color: var(--color-primary);
}

.summary-status {
  margin-top: 24px;
}
</style>
