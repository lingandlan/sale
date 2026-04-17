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

          <el-upload
            class="upload-area"
            drag
            :auto-upload="false"
            :on-change="handleFileChange"
            :limit="1"
            accept="image/*"
          >
            <div class="upload-content">
              <el-icon class="upload-icon"><Upload /></el-icon>
              <div class="upload-text">银行流水截图上传区域</div>
            </div>
          </el-upload>
        </div>

        <div class="button-group">
          <el-button class="cancel-btn" @click="handleCancel">取消</el-button>
          <el-button type="primary" class="submit-btn" @click="handleSubmit">提交申请</el-button>
        </div>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules, type UploadFile } from 'element-plus'
import { extractErrorMessage } from '@/utils/request'
import { Upload } from '@element-plus/icons-vue'
import { submitBRechargeApply } from '@/api/recharge'

interface RechargeCenter {
  id: string
  name: string
}

const router = useRouter()
const formRef = ref<FormInstance>()

const centers = ref<RechargeCenter[]>([
  { id: '1', name: '北京朝阳中心' },
  { id: '2', name: '北京海淀中心' },
  { id: '3', name: '上海浦东中心' }
])

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

const calculatePoints = () => {
  const amount = formData.value.amount || 0

  // 计算规则：≥100000返2%，否则返1%
  const rebateRate = amount >= 100000 ? 2 : 1
  const base = Math.floor(amount)
  const rebate = Math.floor(base * (rebateRate / 100))

  calculatedPoints.value = {
    base,
    rebate,
    rebateRate,
    total: base + rebate
  }
}

// 实时监听金额变化，自动计算积分
watch(() => formData.value.amount, () => {
  calculatePoints()
}, { immediate: true })

const handleFileChange = (file: UploadFile) => {
  formData.value.screenshot = file.name
  ElMessage.success('已选择文件: ' + file.name)
}

const handleCancel = () => {
  formRef.value?.resetFields()
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

  // 验证至少填写一项付款凭证
  if (!formData.value.transactionNo && !formData.value.screenshot) {
    ElMessage.warning('银行流水单号或截图至少需要填写一项')
    return
  }

  try {
    const center = centers.value.find(c => c.id === formData.value.centerId)
    await submitBRechargeApply({
      centerId: formData.value.centerId,
      centerName: center?.name || '',
      amount: formData.value.amount,
      lastMonthConsumption: 0,
      transactionNo: formData.value.transactionNo || '',
      screenshot: formData.value.screenshot || '',
      remark: ''
    })
    ElMessage.success('充值申请已提交，等待审核')
    router.push('/recharge/b-approval')
  } catch (err: any) {
    ElMessage.error(extractErrorMessage(err, '提交充值申请失败'))
  }
}
</script>

<style scoped>
.brecharge-apply {
  display: flex;
  justify-content: center;
  padding: 24px;
  min-height: calc(100vh - 64px);
  background-color: #F5F5F5;
}

.page-card {
  width: 720px;
  background-color: #FFFFFF;
  border-radius: 8px;
  padding: 32px;
}

.page-title {
  font-family: 'Inter', sans-serif;
  font-size: 24px;
  font-weight: 600;
  color: #262626;
  margin: 0 0 24px 0;
}

.apply-form {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.apply-form :deep(.el-form-item__label) {
  font-family: 'Inter', sans-serif;
  font-size: 14px;
  color: #262626;
}

.points-box {
  background-color: #FFF9E6;
  border-radius: 4px;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.points-label {
  font-family: 'Inter', sans-serif;
  font-size: 20px;
  font-weight: 600;
  color: #C00000;
}

.points-detail {
  font-family: 'Inter', sans-serif;
  font-size: 12px;
  color: #595959;
}

.payment-section {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.payment-label {
  font-family: 'Inter', sans-serif;
  font-size: 14px;
  color: #262626;
  margin-bottom: -8px;
}

.upload-area {
  width: 100%;
}

.upload-area :deep(.el-upload-dragger) {
  width: 100%;
  height: 100px;
  background-color: #F5F5F5;
  border: none;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.upload-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.upload-icon {
  font-size: 32px;
  color: #8C8C8C;
}

.upload-text {
  font-family: 'Inter', sans-serif;
  font-size: 13px;
  color: #8C8C8C;
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
  border-radius: 4px;
}

.submit-btn {
  width: 120px;
  height: 40px;
  border-radius: 4px;
  background-color: #C00000;
  border-color: #C00000;
}

.submit-btn:hover {
  background-color: #A00000;
  border-color: #A00000;
}
</style>
