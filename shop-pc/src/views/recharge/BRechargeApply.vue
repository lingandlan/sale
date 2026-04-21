<template>
  <div class="brecharge-apply">
    <div class="page-card">
      <h1 class="page-title">B端充值申请</h1>

      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="auto"
        class="apply-form"
      >
        <el-form-item label="目标充值中心" prop="centerId">
          <el-select
            v-model="formData.centerId"
            placeholder="请选择充值中心"
            @change="handleCenterChange"
            style="width: 100%"
          >
            <el-option
              v-for="center in centers"
              :key="center.id"
              :label="center.name"
              :value="center.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="充值金额（元）" prop="amount">
          <el-input-number
            v-model="formData.amount"
            :min="0"
            :precision="2"
            :step="100"
            style="width: 100%"
            placeholder="50000.00"
          />
        </el-form-item>

        <div class="points-box">
          <div class="points-label">预计获得积分: {{ formatPoints }}</div>
          <div class="points-detail">{{ pointsDetail }}</div>
        </div>

        <div class="payment-section">
          <div class="payment-label">
            付款凭证（银行流水单号或截图至少填一项）
          </div>

          <el-form-item prop="transactionNo">
            <el-input
              v-model="formData.transactionNo"
              placeholder="银行流水单号"
              clearable
            />
          </el-form-item>

          <div class="upload-area" @click="triggerFileInput" @dragover.prevent @drop.prevent="handleDrop">
            <input ref="fileInputRef" type="file" accept="image/*" style="display: none" @change="handleFileSelect" />
            <template v-if="selectedFileName">
              <div class="upload-preview">
                <span class="file-name">{{ selectedFileName }}</span>
                <span class="file-remove" @click.stop="clearFile">&times;</span>
              </div>
            </template>
            <template v-else>
              <div class="upload-content">
                <el-icon class="upload-icon"><Upload /></el-icon>
                <div class="upload-text">点击或拖拽上传银行流水截图</div>
              </div>
            </template>
          </div>
        </div>

        <div class="button-group">
          <el-button class="cancel-btn" @click="handleCancel">取消</el-button>
          <el-button type="primary" class="submit-btn" :loading="uploading" @click="handleSubmit">提交申请</el-button>
        </div>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { extractErrorMessage } from '@/utils/request'
import { Upload } from '@element-plus/icons-vue'
import { submitBRechargeApply, getCenterList, getCenterLastMonthConsumption, uploadFile } from '@/api/recharge'

interface RechargeCenter {
  id: string
  name: string
}

const router = useRouter()
const formRef = ref<FormInstance>()
const fileInputRef = ref<HTMLInputElement>()

const centers = ref<RechargeCenter[]>([])

const loadCenters = async () => {
  try {
    const res = await getCenterList()
    if (res?.data) {
      centers.value = res.data.map((c: any) => ({ id: c.id, name: c.name }))
    }
  } catch {}
}
loadCenters()

const formData = ref({
  centerId: '',
  amount: 0,
  transactionNo: '',
  screenshot: ''
})

const calculatedPoints = ref({
  base: 0,
  rebate: 0,
  rebateRate: 0,
  total: 0
})

const formRules: FormRules = {
  centerId: [{ required: true, message: '请选择充值中心', trigger: 'change' }],
  amount: [
    { required: true, message: '请输入充值金额', trigger: 'blur' },
    { type: 'number', min: 0.01, message: '充值金额必须大于0', trigger: 'blur' }
  ]
}

const formatPoints = computed(() => {
  return calculatedPoints.value.total.toLocaleString()
})

const pointsDetail = computed(() => {
  const { base, rebate, rebateRate } = calculatedPoints.value
  return `基础 ${base.toLocaleString()} + 返还 ${rebate.toLocaleString()} (${rebateRate}%)`
})

const lastMonthConsumption = ref(0)
const rebateRate = ref(1)
const lastMonth = ref('')

const uploading = ref(false)
const pendingFile = ref<File | null>(null)
const selectedFileName = ref('')

// 选中充值中心时查询上月消费
const handleCenterChange = async () => {
  if (!formData.value.centerId) return
  try {
    const res = await getCenterLastMonthConsumption(formData.value.centerId)
    if (res?.data) {
      lastMonthConsumption.value = res.data.consumption
      rebateRate.value = res.data.rebateRate
      lastMonth.value = res.data.month
    }
  } catch {}
  calculatePoints()
}

const calculatePoints = () => {
  const amount = formData.value.amount || 0
  const rate = rebateRate.value
  const base = Math.floor(amount)
  const rebate = Math.floor(base * (rate / 100))

  calculatedPoints.value = {
    base,
    rebate,
    rebateRate: rate,
    total: base + rebate
  }
}

watch(() => formData.value.amount, () => {
  calculatePoints()
}, { immediate: true })

const triggerFileInput = () => {
  fileInputRef.value?.click()
}

const handleFileSelect = async (e: Event) => {
  const input = e.target as HTMLInputElement
  if (input.files && input.files[0]) {
    const file = input.files[0]
    // 立即读入内存，避免 Chrome ERR_UPLOAD_FILE_CHANGED
    try {
      const buffer = await file.arrayBuffer()
      pendingFile.value = new File([buffer], file.name, { type: file.type })
      selectedFileName.value = file.name
      ElMessage.success('已选择文件: ' + file.name)
    } catch {
      ElMessage.error('无法读取文件，请重试')
    }
  }
}

const handleDrop = async (e: DragEvent) => {
  if (e.dataTransfer?.files && e.dataTransfer.files[0]) {
    const file = e.dataTransfer.files[0]
    try {
      const buffer = await file.arrayBuffer()
      pendingFile.value = new File([buffer], file.name, { type: file.type })
      selectedFileName.value = file.name
      ElMessage.success('已选择文件: ' + file.name)
    } catch {
      ElMessage.error('无法读取文件，请重试')
    }
  }
}

const clearFile = () => {
  pendingFile.value = null
  selectedFileName.value = ''
  if (fileInputRef.value) fileInputRef.value.value = ''
}

const handleCancel = () => {
  formRef.value?.resetFields()
  clearFile()
  calculatedPoints.value = {
    base: 0,
    rebate: 0,
    rebateRate: 0,
    total: 0
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
  } catch {
    return
  }

  if (!formData.value.transactionNo && !pendingFile.value) {
    ElMessage.warning('银行流水单号或截图至少需要填写一项')
    return
  }

  uploading.value = true
  try {
    let screenshotUrl = ''
    if (pendingFile.value) {
      const uploadRes = await uploadFile(pendingFile.value)
      screenshotUrl = uploadRes?.data?.url || ''
    }

    const center = centers.value.find(c => c.id === formData.value.centerId)
    await submitBRechargeApply({
      centerId: formData.value.centerId,
      centerName: center?.name || '',
      amount: formData.value.amount,
      lastMonthConsumption: lastMonthConsumption.value,
      transactionNo: formData.value.transactionNo || '',
      screenshot: screenshotUrl,
      remark: ''
    })
    ElMessage.success('充值申请已提交，等待审核')
    router.push('/recharge/b-approval')
  } catch (err: any) {
    ElMessage.error(extractErrorMessage(err, '提交充值申请失败'))
  } finally {
    uploading.value = false
  }
}
</script>

<style scoped>
.brecharge-apply {
  display: flex;
  justify-content: center;
  padding: 24px;
  min-height: calc(100vh - 64px);
  background-color: var(--color-bg);
}

.page-card {
  width: 720px;
  background-color: var(--color-bg-card);
  border-radius: var(--radius-md);
  padding: 32px;
}

.page-title {
  font-family: var(--font-family);
  font-size: 24px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0 0 24px 0;
}

.apply-form {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.apply-form :deep(.el-form-item__label) {
  font-family: var(--font-family);
  font-size: 14px;
  color: var(--color-text-primary);
}

.points-box {
  background-color: #FFF9E6;
  border-radius: var(--radius-sm);
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.points-label {
  font-family: var(--font-family);
  font-size: 20px;
  font-weight: 600;
  color: var(--color-primary);
}

.points-detail {
  font-family: var(--font-family);
  font-size: 12px;
  color: var(--color-text-secondary);
}

.payment-section {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.payment-label {
  font-family: var(--font-family);
  font-size: 14px;
  color: var(--color-text-primary);
  margin-bottom: -8px;
}

.upload-area {
  width: 100%;
  height: 100px;
  background-color: var(--color-bg);
  border: 1px dashed var(--color-border);
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: border-color 0.3s;
}

.upload-area:hover {
  border-color: var(--color-primary);
}

.upload-preview {
  display: flex;
  gap: 12px;
  align-items: center;
}

.file-name {
  font-family: var(--font-family);
  font-size: 14px;
  color: var(--color-primary);
}

.file-remove {
  font-size: 18px;
  color: var(--color-text-secondary);
  cursor: pointer;
}

.file-remove:hover {
  color: var(--color-danger, #f56c6c);
}

.upload-content {
  display: flex;
  gap: 8px;
  align-items: center;
}

.upload-icon {
  font-size: 32px;
  color: var(--color-text-muted);
}

.upload-text {
  font-family: var(--font-family);
  font-size: 13px;
  color: var(--color-text-muted);
}

.button-group {
  display: flex;
  justify-content: flex-end;
  gap: 16px;
  margin-top: 8px;
}

.cancel-btn {
  width: 100px;
  height: 40px;
  border-radius: var(--radius-sm);
}

.submit-btn {
  width: 120px;
  height: 40px;
  border-radius: var(--radius-sm);
  background-color: var(--color-primary);
  border-color: var(--color-primary);
}

.submit-btn:hover {
  background-color: var(--color-primary-hover);
  border-color: var(--color-primary-hover);
}
</style>
