<script setup lang="ts">
import { onLaunch, onShow, onHide } from '@dcloudio/uni-app'
import { isAuthenticated } from '@/utils/auth'

onLaunch(() => {
  console.log('App Launch')
  // 检查登录状态
  if (!isAuthenticated()) {
    uni.reLaunch({
      url: '/pages/login/index'
    })
  }
})

onShow(() => {
  console.log('App Show')
})

onHide(() => {
  console.log('App Hide')
})
</script>

<style lang="scss">
/* 全局样式 */
page {
  background-color: #F5F5F5;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', 'Helvetica Neue', Helvetica, Arial, sans-serif;
}

/* 清除默认样式 */
* {
  box-sizing: border-box;
}

view, text {
  margin: 0;
  padding: 0;
}
</style>
